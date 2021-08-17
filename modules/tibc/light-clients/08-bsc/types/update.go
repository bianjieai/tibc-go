package types

import (
	fmt "fmt"
	io "io"
	"math/big"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"

	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

func (m ClientState) CheckHeaderAndUpdateState(
	ctx sdk.Context,
	cdc codec.BinaryMarshaler,
	store sdk.KVStore,
	header exported.Header,
) (exported.ClientState, exported.ConsensusState, error) {
	bscHeader, ok := header.(*Header)
	if !ok {
		return nil, nil, sdkerrors.Wrapf(
			clienttypes.ErrInvalidHeader, "expected type %T, got %T", &Header{}, header,
		)
	}

	// get consensus state from clientStore
	bscConsState, err := GetConsensusState(store, cdc, bscHeader.Height)
	if err != nil {
		return nil, nil, sdkerrors.Wrapf(
			err, "could not get consensus state from clientstore at TrustedHeight: %s", bscHeader.Height,
		)
	}

	if err := checkValidity(cdc, store, &m, bscConsState, *bscHeader); err != nil {
		return nil, nil, err
	}
	newClientState, consensusState := update(ctx, store, &m, bscHeader)
	return newClientState, consensusState, nil
}

// checkValidity checks if the bsc header is valid.
func checkValidity(
	cdc codec.BinaryMarshaler,
	store sdk.KVStore,
	clientState *ClientState,
	consState *ConsensusState,
	header Header,

) error {
	if err := header.ValidateBasic(); err != nil {
		return err
	}

	return verifyHeader(cdc, store, clientState, header)
}

// update the RecentSingers and the ConsensusState.
func update(ctx sdk.Context, clientStore sdk.KVStore, clientState *ClientState, header *Header) (*ClientState, *ConsensusState) {
	//height := header.GetHeight().(clienttypes.Height)
	// TODO
	return nil, nil
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
	number := header.Height.RevisionHeight

	parent := clientState.Header
	if parent.Height.RevisionHeight != number-1 || parent.Hash() != common.BytesToHash(header.ParentHash) {
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
	limit := parent.GasLimit / GasLimitBoundDivisor

	if uint64(diff) >= limit || header.GasLimit < params.MinGasLimit {
		return fmt.Errorf("invalid gas limit: have %d, want %d += %d", header.GasLimit, parent.GasLimit, limit)

	}
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

	lastHeight, _ := header.Height.Decrement()
	validatorSet := GetValidators(cdc, store, lastHeight)
	if !validatorSet.Has(signer.Bytes()) {
		return sdkerrors.Wrap(ErrUnauthorizedValidator, signer.Hex())
	}

	recentSingers, err := GetRecentSingers(store)
	if err != nil {
		return err
	}

	for _, recent := range recentSingers {
		seen := recent.Height.RevisionHeight
		if common.BytesToAddress(recent.Validator) == signer {
			// Signer is among recent, only fail if the current block doesn't shift it out
			if limit := uint64(len(validatorSet.Validators)/2 + 1); seen > number-limit {
				return sdkerrors.Wrap(ErrRecentlySigned, signer.Hex())
			}
		}
	}

	// Ensure that the difficulty corresponds to the turn-ness of the signer
	inturn := validatorSet.inturn(lastHeight.GetRevisionHeight(), signer.Bytes())
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
	if len(header.Extra) < ExtraSeal {
		return common.Address{}, sdkerrors.Wrap(ErrMissingSignature, "header.Extra")
	}
	signature := header.Extra[len(header.Extra)-ExtraSeal:]

	// Recover the public key and the Ethereum address
	pubkey, err := crypto.Ecrecover(SealHash(header, chainId).Bytes(), signature)
	if err != nil {
		return common.Address{}, err
	}
	var signer common.Address
	copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])

	return signer, nil
}

// SealHash returns the hash of a block prior to it being sealed.
func SealHash(header Header, chainId *big.Int) (hash common.Hash) {
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
