package types

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"time"

	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
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

	return rlpHash(h.ToVerifyHeader())
}

//func (h *EthHeader) Hash() (hash common.Hash) {
//	return rlpHash(h)
//}

func (h Header) ValidateBasic() error {
	// Ensure that the header's extra-data section is of a reasonable size
	if uint64(len(h.Extra)) > params.MaximumExtraDataSize {
		return fmt.Errorf("extra-data too long: %d > %d", len(h.Extra), params.MaximumExtraDataSize)
	}
	if h.Time > uint64(time.Now().Unix()+allowedFutureBlockTimeSeconds) {

		return sdkerrors.Wrap(ErrFutureBlock, consensus.ErrFutureBlock.Error())
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

//
//func (h Header) ToEthHeader() EthHeader {
//	return EthHeader{
//		ParentHash:  common.BytesToHash(h.ParentHash),
//		UncleHash:   common.BytesToHash(h.UncleHash),
//		Coinbase:    common.BytesToAddress(h.Coinbase),
//		Root:        common.BytesToHash(h.Root),
//		TxHash:      common.BytesToHash(h.TxHash),
//		ReceiptHash: common.BytesToHash(h.ReceiptHash),
//		Bloom:       types.BytesToBloom(h.Bloom),
//		Difficulty:  big.NewInt(int64(h.Difficulty)),
//		Number:      big.NewInt(int64(h.Height.RevisionHeight)),
//		GasLimit:    h.GasLimit,
//		GasUsed:     h.GasUsed,
//		Time:        h.Time,
//		Extra:       h.Extra,
//		MixDigest:   common.BytesToHash(h.MixDigest),
//		Nonce:       types.EncodeNonce(h.Nonce),
//		BaseFee:     big.NewInt(int64(h.BaseFee)),
//	}
//}
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
	exist, err := IsHeaderExist(store, header.Hash())
	if err != nil {
		return sdkerrors.Wrap(ErrUnknownBlock, fmt.Errorf("SyncBlockHeader, check header exist err: %v", err).Error())
	}
	if exist == true {
		return sdkerrors.Wrap(ErrHeaderIsExist, "header"+header.String())
	}

	parentbytes := store.Get(host.ConsensusStateIndexKey(header.ToVerifyHeader().ParentHash))
	var parentConsInterface exported.ConsensusState
	if err := cdc.UnmarshalInterface(parentbytes, &parentConsInterface); err != nil {
		return sdkerrors.Wrap(ErrUnmarshalInterface, err.Error())
	}
	parent := parentConsInterface.(*ConsensusState)
	if parent.Header.Height.RevisionHeight != height-1 || parent.Header.Hash() != common.BytesToHash(header.ParentHash) {
		return sdkerrors.Wrap(ErrUnknownAncestor, "")
	}

	//verify whether parent hash validity
	ethHeader := parent.Header.ToVerifyHeader()
	if !bytes.Equal(ethHeader.Hash().Bytes(), header.ToVerifyHeader().ParentHash.Bytes()) {
		return fmt.Errorf("SyncBlockHeader, parent header is not right. Header: %s", header.String())
	}
	//verify whether extra size validity
	if uint64(len(header.Extra)) > params.MaximumExtraDataSize {
		return sdkerrors.Wrap(ErrExtraLenth, fmt.Errorf("SyncBlockHeader, SyncBlockHeader extra-data too long: %d > %d, header: %s", len(header.Extra), params.MaximumExtraDataSize, header.String()).Error())
	}
	// Verify the header's timestamp
	if header.Time > uint64(time.Now().Unix()+allowedFutureBlockTimeSeconds) {
		return ErrFutureBlock
	}
	if header.Time <= parent.Header.Time {
		return ErrFutureBlock
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
	if expected.Cmp(header.ToVerifyHeader().Difficulty) != 0 {
		return sdkerrors.Wrap(ErrInvalidDifficult, fmt.Errorf("SyncBlockHeader, invalid difficulty: have %v, want %v, header: %s", header.Difficulty, expected, header.String()).Error())
	}
	// todo !

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
	if err := ethash.verifySeal(header.ToVerifyHeader(), false); err != nil {
		return ErrUnknownBlock
	}
	// All basic checks passed
	return nil

}

func IsHeaderExist(store sdk.KVStore, hash common.Hash) (bool, error) {
	headerStore := store.Get(host.ConsensusStateIndexKey(hash))
	if headerStore == nil {
		return false, nil
	} else {
		return true, nil
	}
}

// Some weird constants to avoid constant memory allocs for them.
var (
	expDiffPeriod = big.NewInt(100000)
	big1          = big.NewInt(1)
	big2          = big.NewInt(2)
	big9          = big.NewInt(9)
	bigMinus99    = big.NewInt(-99)
)

// makeDifficultyCalculator creates a difficultyCalculator with the given bomb-delay.
// the difficulty is calculated with Byzantium rules, which differs from Homestead in
// how uncles affect the calculation
func makeDifficultyCalculator(bombDelay *big.Int) func(time uint64, parent *Header) *big.Int {
	// Note, the calculations below looks at the parent number, which is 1 below
	// the block number. Thus we remove one from the delay given
	bombDelayFromParent := new(big.Int).Sub(bombDelay, big1)
	return func(time uint64, parent *Header) *big.Int {
		// https://github.com/ethereum/EIPs/issues/100.
		// algorithm:
		// diff = (parent_diff +
		//         (parent_diff / 2048 * max((2 if len(parent.uncles) else 1) - ((timestamp - parent.timestamp) // 9), -99))
		//        ) + 2^(periodCount - 2)

		bigTime := new(big.Int).SetUint64(time)
		bigParentTime := new(big.Int).SetUint64(parent.Time)

		// holds intermediate values to make the algo easier to read & audit
		x := new(big.Int)
		y := new(big.Int)

		// (2 if len(parent_uncles) else 1) - (block_timestamp - parent_timestamp) // 9
		x.Sub(bigTime, bigParentTime)
		x.Div(x, big9)
		if parent.ToVerifyHeader().UncleHash == types.EmptyUncleHash {
			x.Sub(big1, x)
		} else {
			x.Sub(big2, x)
		}
		// max((2 if len(parent_uncles) else 1) - (block_timestamp - parent_timestamp) // 9, -99)
		if x.Cmp(bigMinus99) < 0 {
			x.Set(bigMinus99)
		}
		// parent_diff + (parent_diff / 2048 * max((2 if len(parent.uncles) else 1) - ((timestamp - parent.timestamp) // 9), -99))
		y.Div(parent.ToVerifyHeader().Difficulty, params.DifficultyBoundDivisor)
		x.Mul(y, x)
		x.Add(parent.ToVerifyHeader().Difficulty, x)

		// minimum difficulty can ever be (before exponential factor)
		if x.Cmp(params.MinimumDifficulty) < 0 {
			x.Set(params.MinimumDifficulty)
		}
		// calculate a fake block number for the ice-age delay
		// Specification: https://eips.ethereum.org/EIPS/eip-1234
		fakeBlockNumber := new(big.Int)
		if parent.ToVerifyHeader().Number.Cmp(bombDelayFromParent) >= 0 {
			fakeBlockNumber = fakeBlockNumber.Sub(parent.ToVerifyHeader().Number, bombDelayFromParent)
		}
		// for the exponential factor
		periodCount := fakeBlockNumber
		periodCount.Div(periodCount, expDiffPeriod)

		// the exponential factor, commonly referred to as "the bomb"
		// diff = diff + 2^(periodCount - 2)
		if periodCount.Cmp(big1) > 0 {
			y.Sub(periodCount, big2)
			y.Exp(big2, y, nil)
			x.Add(x, y)
		}
		return x
	}
}

const (
	BaseFeeChangeDenominator = 8          // Bounds the amount the base fee can change between blocks.
	ElasticityMultiplier     = 2          // Bounds the maximum gas limit an EIP-1559 block may have.
	InitialBaseFee           = 1000000000 // Initial base fee for EIP-1559 blocks.
)

// VerifyEip1559Header verifies some header attributes which were changed in EIP-1559,
// - gas limit check
// - basefee check
func VerifyEip1559Header(parent, header *Header) error {
	// Verify that the gas limit remains within allowed bounds
	parentGasLimit := parent.ToVerifyHeader().GasLimit

	if err := VerifyGaslimit(parentGasLimit, header.ToVerifyHeader().GasLimit); err != nil {
		return err
	}
	// Verify the header is not malformed
	if header.ToVerifyHeader().BaseFee == nil {
		return fmt.Errorf("header is missing baseFee")
	}
	// Verify the baseFee is correct based on the parent header.
	expectedBaseFee := CalcBaseFee(parent)
	if header.ToVerifyHeader().BaseFee.Cmp(expectedBaseFee) != 0 {
		return fmt.Errorf("invalid baseFee: have %s, want %s, parentBaseFee %s, parentGasUsed %d",
			expectedBaseFee, header.ToVerifyHeader().BaseFee, parent.ToVerifyHeader().BaseFee, parent.ToVerifyHeader().GasUsed)
	}
	return nil
}

// VerifyGaslimit verifies the header gas limit according increase/decrease
// in relation to the parent gas limit.
func VerifyGaslimit(parentGasLimit, headerGasLimit uint64) error {
	// Verify that the gas limit remains within allowed bounds
	diff := int64(parentGasLimit) - int64(headerGasLimit)
	if diff < 0 {
		diff *= -1
	}
	limit := parentGasLimit / params.GasLimitBoundDivisor
	if uint64(diff) >= limit {
		return fmt.Errorf("invalid gas limit: have %d, want %d +-= %d", headerGasLimit, parentGasLimit, limit-1)
	}
	if headerGasLimit < params.MinGasLimit {
		return errors.New("invalid gas limit below 5000")
	}
	return nil
}

// CalcBaseFee calculates the basefee of the header.
func CalcBaseFee(parent *Header) *big.Int {

	var (
		parentGasTarget          = parent.ToVerifyHeader().GasLimit / ElasticityMultiplier
		parentGasTargetBig       = new(big.Int).SetUint64(parentGasTarget)
		baseFeeChangeDenominator = new(big.Int).SetUint64(BaseFeeChangeDenominator)
	)
	// If the parent gasUsed is the same as the target, the baseFee remains unchanged.
	if parent.ToVerifyHeader().GasUsed == parentGasTarget {
		return new(big.Int).Set(parent.ToVerifyHeader().BaseFee)
	}
	if parent.ToVerifyHeader().GasUsed > parentGasTarget {
		// If the parent block used more gas than its target, the baseFee should increase.
		gasUsedDelta := new(big.Int).SetUint64(parent.ToVerifyHeader().GasUsed - parentGasTarget)
		x := new(big.Int).Mul(parent.ToVerifyHeader().BaseFee, gasUsedDelta)
		y := x.Div(x, parentGasTargetBig)
		baseFeeDelta := math.BigMax(
			x.Div(y, baseFeeChangeDenominator),
			common.Big1,
		)

		return x.Add(parent.ToVerifyHeader().BaseFee, baseFeeDelta)
	} else {
		// Otherwise if the parent block used less gas than its target, the baseFee should decrease.
		gasUsedDelta := new(big.Int).SetUint64(parentGasTarget - parent.ToVerifyHeader().GasUsed)
		x := new(big.Int).Mul(parent.ToVerifyHeader().BaseFee, gasUsedDelta)
		y := x.Div(x, parentGasTargetBig)
		baseFeeDelta := x.Div(y, baseFeeChangeDenominator)

		return math.BigMax(
			x.Sub(parent.ToVerifyHeader().BaseFee, baseFeeDelta),
			common.Big0,
		)
	}
}
