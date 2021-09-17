package types

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"

	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

var (
	allowedFutureBlockTime = 15 * time.Second
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

func (h *EthHeader) Hash() (hash common.Hash) {
	return rlpHash(h)
}

func (h Header) ValidateBasic() error {
	// Ensure that the header's extra-data section is of a reasonable size
	if uint64(len(h.Extra)) > params.MaximumExtraDataSize {
		return fmt.Errorf("extra-data too long: %d > %d", len(h.Extra), params.MaximumExtraDataSize)
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
		Bloom:       types.BytesToBloom(h.Bloom),
		Difficulty:  big.NewInt(int64(h.Difficulty)),
		Number:      big.NewInt(int64(h.Height.RevisionHeight)),
		GasLimit:    h.GasLimit,
		GasUsed:     h.GasUsed,
		Time:        h.Time,
		Extra:       h.Extra,
		MixDigest:   common.BytesToHash(h.MixDigest),
		Nonce:       types.EncodeNonce(h.Nonce),
		BaseFee:     big.NewInt(int64(h.BaseFee)),
	}
}

func (h Header) ToVerifyHeader() *types.Header {
	return &types.Header{
		ParentHash:  common.BytesToHash(h.ParentHash),
		UncleHash:   common.BytesToHash(h.UncleHash),
		Coinbase:    common.BytesToAddress(h.Coinbase),
		Root:        common.BytesToHash(h.Root),
		TxHash:      common.BytesToHash(h.TxHash),
		ReceiptHash: common.BytesToHash(h.ReceiptHash),
		Bloom:       types.BytesToBloom(h.Bloom),
		Difficulty:  big.NewInt(int64(h.Difficulty)),
		Number:      big.NewInt(int64(h.Height.RevisionHeight)),
		GasLimit:    h.GasLimit,
		GasUsed:     h.GasUsed,
		Time:        h.Time,
		Extra:       h.Extra,
		MixDigest:   common.BytesToHash(h.MixDigest),
		Nonce:       types.EncodeNonce(h.Nonce),
		BaseFee:     big.NewInt(int64(h.BaseFee)),
	}
}

func verifyHeader(
	ctx sdk.Context,
	cdc codec.BinaryMarshaler,
	store sdk.KVStore,
	clientState *ClientState,
	header Header,
) error {
	height := header.Height.RevisionHeight
	exist, err := IsHeaderExist(store, header.Hash())
	if err != nil {
		return sdkerrors.Wrap(ErrUnknownBlock, fmt.Errorf("SyncBlockHeader, check header exist err: %v", err).Error())
	}
	if exist == true {
		return sdkerrors.Wrap(ErrHeaderIsExist, "header"+header.String())
	}

	parentbytes := store.Get(ConsensusStateIndexKey(header.ToEthHeader().ParentHash))
	var parentConsInterface exported.ConsensusState
	if err := cdc.UnmarshalInterface(parentbytes, &parentConsInterface); err != nil {
		return sdkerrors.Wrap(ErrUnmarshalInterface, err.Error())
	}
	parent := parentConsInterface.(*ConsensusState)
	if parent.Header.Height.RevisionHeight != height-1 || parent.Header.Hash() != common.BytesToHash(header.ParentHash) {
		return sdkerrors.Wrap(ErrUnknownAncestor, "")
	}

	//verify whether parent hash validity
	ethHeader := parent.Header.ToEthHeader()
	if !bytes.Equal(ethHeader.Hash().Bytes(), header.ToEthHeader().ParentHash.Bytes()) {
		return fmt.Errorf("SyncBlockHeader, parent header is not right. Header: %s", header.String())
	}
	//verify whether extra size validity
	if uint64(len(header.Extra)) > params.MaximumExtraDataSize {
		return sdkerrors.Wrap(ErrExtraLenth, fmt.Errorf("SyncBlockHeader, SyncBlockHeader extra-data too long: %d > %d, header: %s", len(header.Extra), params.MaximumExtraDataSize, header.String()).Error())
	}
	// Verify the header's timestamp
	if header.Time > uint64(ctx.BlockTime().Add(allowedFutureBlockTime).Unix()) {
		return ErrFutureBlock
	}
	if header.Time <= parent.Header.Time {
		return ErrHeader
	}

	// Verify that the gas limit is <= 2^63-1
	capacity := uint64(0x7fffffffffffffff)
	if header.GasLimit > capacity {
		return sdkerrors.Wrap(ErrInvalidGas, fmt.Errorf("invalid gasLimit: have %v, max %v", header.GasLimit, capacity).Error())

	}
	// Verify that the gasUsed is <= gasLimit
	if header.GasUsed > header.GasLimit {
		return sdkerrors.Wrap(ErrInvalidGas, fmt.Errorf("invalid gasUsed: have %d, gasLimit %d", header.GasUsed, header.GasLimit).Error())
	}
	err = VerifyEip1559Header(&parent.Header, &header)
	if err != nil {
		return sdkerrors.Wrap(ErrHeader, fmt.Errorf("SyncBlockHeader, err:%v", err).Error())
	}
	//verify difficulty
	expected := makeDifficultyCalculator(big.NewInt(9700000))(header.Time, &parent.Header)
	if expected.Cmp(header.ToEthHeader().Difficulty) != 0 {
		return sdkerrors.Wrap(ErrWrongDifficulty, fmt.Errorf("SyncBlockHeader, invalid difficulty: have %v, want %v, header: %s", header.Difficulty, expected, header.String()).Error())
	}

	return verifyCascadingFields(header)
}

// verifyCascadingFields verifies all the header fields that are not standalone,
// rather depend on a batch of previous headers. The caller may optionally pass
// in a batch of parents (ascending order) to avoid looking those up from the
// database. This is useful for concurrently verifying a batch of new headers.
func verifyCascadingFields(header Header) error {
	cachedir, err := ioutil.TempDir("", "")
	if err != nil {
		fmt.Println(err)
		return errEthashStopped
	}
	defer os.RemoveAll(cachedir)
	config := Config{
		CacheDir:     cachedir,
		CachesOnDisk: 1,
	}
	ethash := New(config, nil, false)
	defer ethash.Close()
	if err := ethash.VerifySeal(header.ToVerifyHeader(), false); err != nil {
		return ErrHeader
	}
	// All basic checks passed
	return nil

}
