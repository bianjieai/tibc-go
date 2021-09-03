package types

import (
	fmt "fmt"
	io "io"
	"math/big"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

var (
	// Maximum number of uncles allowed in a single block
	allowedFutureBlockTimeSeconds = int64(15)
)
var _ exported.Header = (*Header)(nil)

func (h Header) ClientType() string {
	return exported.ETH
}

func (h Header) GetHeight() exported.Height {
	return h.Height
}

// Hash returns the block hash of the header, which is simply the keccak256 hash of its
// RLP encoding.
func (h *Header) Hash() common.Hash {
	return rlpHash(h.ToEthHeader())
}

func (h Header) ValidateBasic() error {
	// Ensure that the header's extra-data section is of a reasonable size
	if uint64(len(h.Extra)) > params.MaximumExtraDataSize {
		return fmt.Errorf("extra-data too long: %d > %d", len(h.Extra), params.MaximumExtraDataSize)
	}

	// Ensure that the block's difficulty is meaningful (may not be correct at this point)
	number := h.Height.RevisionHeight
	if number > 0 {
		if h.Difficulty == 0 {
			return sdkerrors.Wrap(ErrInvalidDifficulty, "header Difficulty")
		}
	}
	return nil
}

func (h Header) ToEthHeader() EthHeader {
	return EthHeader{
		ParentHash:  common.BytesToHash(h.ParentHash),
		UncleHash:   common.BytesToHash(h.UncleHash),
		Coinbase:    common.BytesToAddress(h.Coinbase),
		Root:        common.BytesToHash(h.Root),
		TxHash:      common.BytesToHash(h.TxHash),
		ReceiptHash: common.BytesToHash(h.ReceiptHash),
		Bloom:       BytesToBloom(h.Bloom),
		Difficulty:  big.NewInt(int64(h.Difficulty)),
		Number:      big.NewInt(int64(h.Height.RevisionHeight)),
		GasLimit:    h.GasLimit,
		GasUsed:     h.GasUsed,
		Time:        h.Time,
		Extra:       h.Extra,
		MixDigest:   common.BytesToHash(h.MixDigest),
		Nonce:       BytesToBlockNonce(h.Nonce),
	}
}

func verifyHeader(
	cdc codec.BinaryMarshaler,
	store sdk.KVStore,
	clientState *ClientState,
	header Header,
) error {
	if err := header.ValidateBasic(); err != nil {
		return err
	}
	//todo
	return verifyCascadingFields(cdc, store, clientState, header)
}

// verifyCascadingFields verifies all the header fields that are not standalone,
// rather depend on a batch of previous headers. The caller may optionally pass
// in a batch of parents (ascending order) to avoid looking those up from the
// database. This is useful for concurrently verifying a batch of new headers.
func verifyCascadingFields(
	cdc codec.BinaryMarshaler,
	store sdk.KVStore,
	clientState *ClientState,
	header Header) error {
	height := header.Height.RevisionHeight

	parent := clientState.Header
	if parent.Height.RevisionHeight != height-1 || parent.Hash() != common.BytesToHash(header.ParentHash) {
		return sdkerrors.Wrap(ErrUnknownAncestor, "")
	}

	// Verify that the gas limit is <= 2^63-1
	capacity := uint64(0x7fffffffffffffff)
	if header.GasLimit > capacity {
		return fmt.Errorf("invalid gasLimit: have %v, max %v", header.GasLimit, capacity)
	}
	// Verify that the gasUsed is <= gasLimit
	if header.GasUsed > header.GasLimit {
		return fmt.Errorf("invalid gasUsed: have %d, gasLimit %d", header.GasUsed, header.GasLimit)
	}

	// Verify that the gas limit remains within allowed bounds
	diff := int64(parent.GasLimit) - int64(header.GasLimit)
	if diff < 0 {
		diff *= -1
	}
	limit := parent.GasLimit / gasLimitBoundDivisor

	if uint64(diff) >= limit || header.GasLimit < params.MinGasLimit {
		return fmt.Errorf("invalid gas limit: have %d, want %d += %d", header.GasLimit, parent.GasLimit, limit)

	}
	// All basic checks passed, verify the seal and return
	return verifySeal(cdc, store, clientState, header)
}

func verifySeal(
	cdc codec.BinaryMarshaler,
	store sdk.KVStore,
	clientState *ClientState,
	header Header) error {

	number := header.Height.RevisionHeight
	// Resolve the authorization key and check against validators
	signer, err := ecrecover(header, big.NewInt(int64(clientState.ChainId)))
	if err != nil {
		return err
	}

	if signer != common.BytesToAddress(header.Coinbase) {
		return sdkerrors.Wrap(ErrCoinBaseMisMatch, "header.Coinbase")
	}

	// Retrieve the snapshot needed to verify this header and cache it
	snap, err := clientState.snapshot(cdc, store)
	if err != nil {
		return err
	}

	if _, ok := snap.Validators[signer]; !ok {
		return sdkerrors.Wrap(ErrUnauthorizedValidator, signer.Hex())
	}

	for seen, recent := range snap.Recents {
		if recent == signer {
			// Signer is among recents, only fail if the current block doesn't shift it out
			if limit := uint64(len(snap.Validators)/2 + 1); seen > number-limit {
				return sdkerrors.Wrap(ErrRecentlySigned, signer.Hex())
			}
		}
	}
	inturn := snap.inturn(signer)
	diff := big.NewInt(int64(header.Difficulty))
	if inturn && diff.Cmp(diffInTurn) != 0 {
		return sdkerrors.Wrap(ErrWrongDifficulty, "header.Difficulty")
	}
	if !inturn && diff.Cmp(diffNoTurn) != 0 {
		return sdkerrors.Wrap(ErrWrongDifficulty, "header.Difficulty")
	}
	return nil
}

// ecrecover extracts the Ethereum account address from a signed header.
func ecrecover(header Header, chainId *big.Int) (common.Address, error) {
	// Retrieve the signature from the header extra-data
	if len(header.Extra) < extraSeal {
		return common.Address{}, sdkerrors.Wrap(ErrMissingSignature, "header.Extra")
	}
	signature := header.Extra[len(header.Extra)-extraSeal:]

	// Recover the public key and the Ethereum address
	pubkey, err := crypto.Ecrecover(sealHash(header, chainId).Bytes(), signature)
	if err != nil {
		return common.Address{}, err
	}
	var signer common.Address
	copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])

	return signer, nil
}

// sealHash returns the hash of a block prior to it being sealed.
func sealHash(header Header, chainId *big.Int) (hash common.Hash) {
	hasher := sha3.NewLegacyKeccak256()
	encodeSigHeader(hasher, header, chainId)
	hasher.Sum(hash[:0])
	return hash
}

func encodeSigHeader(w io.Writer, header Header, chainId *big.Int) {
	err := rlp.Encode(w, []interface{}{
		chainId,
		header.ParentHash,
		header.UncleHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Height.RevisionHeight,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra[:len(header.Extra)-65], // this will panic if extra is too short, should check before calling encodeSigHeader
		header.MixDigest,
		header.Nonce,
	})
	if err != nil {
		panic("can't encode: " + err.Error())
	}
}
