package types

import (
	fmt "fmt"
	"math/big"
	"time"

	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"

	"github.com/ethereum/go-ethereum/consensus"

	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
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
	if h.Time > uint64(time.Now().Unix()+allowedFutureBlockTimeSeconds) {
		return consensus.ErrFutureBlock
	}
	// Verify that the gas limit is <= 2^63-1
	cap := uint64(0x7fffffffffffffff)
	if h.GasLimit > cap {
		return fmt.Errorf("invalid gasLimit: have %v, max %v", h.GasLimit, cap)
	}
	// Verify that the gasUsed is <= gasLimit
	if h.GasUsed > h.GasLimit {
		return fmt.Errorf("invalid gasUsed: have %d, gasLimit %d", h.GasUsed, h.GasLimit)
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
	exist, err := IsHeaderExist(store, header.Hash().String())
	if err != nil {
		return fmt.Errorf("SyncBlockHeader, check header exist err: %v", err)
	}
	if exist == true {
		return fmt.Errorf("SyncBlockHeader, header has exist. Header: %s", header.String())
	}

	parentbytes := store.Get(host.ConsensusStateIndexKey(string(header.ParentHash)))
	var parent ConsensusState
	err1 := cdc.UnmarshalInterface(parentbytes, parent)
	if err1 != nil {
		return err1
	}
	if parent.Header.Height.RevisionHeight != height-1 || parent.Header.Hash() != common.BytesToHash(header.ParentHash) {
		return sdkerrors.Wrap(ErrUnknownAncestor, "")
	}

	// Verify the header's timestamp
	if header.Time > uint64(time.Now().Unix()+allowedFutureBlockTimeSeconds) {
		return fmt.Errorf("block in the future")
	}
	if header.Time <= parent.Header.Time {
		return fmt.Errorf("timestamp older than parent")
	}

	//todo ? Verify the block's difficulty based on its timestamp and parent's difficulty

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
	diff := int64(parent.Header.GasLimit) - int64(header.GasLimit)
	if diff < 0 {
		diff *= -1
	}
	limit := parent.Header.GasLimit / gasLimitBoundDivisor

	if uint64(diff) >= limit || header.GasLimit < params.MinGasLimit {
		return fmt.Errorf("invalid gas limit: have %d, want %d += %d", header.GasLimit, parent.Header.GasLimit, limit)

	}
	return nil
	// All basic checks passed, verify the seal and return
}

func IsHeaderExist(store sdk.KVStore, hash string) (bool, error) {
	headerStore := store.Get(host.ConsensusStateIndexKey(hash))
	if headerStore == nil {
		return false, nil
	} else {
		return true, nil
	}
}
