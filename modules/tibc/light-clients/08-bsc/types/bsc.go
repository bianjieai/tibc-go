package types

import (
	fmt "fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
)

const (
	// Fixed number of extra-data prefix bytes reserved for signer vanity
	extraVanity = 32

	// Fixed number of extra-data suffix bytes reserved for signer seal
	extraSeal = 65

	// AddressLength is the expected length of the address
	addressLength = 20

	// BloomByteLength represents the number of bytes used in a header log bloom.
	bloomByteLength = 256

	// NonceByteLength represents the number of bytes used in a header log nonce.
	nonceByteLength = 8

	// The bound divisor of the gas limit, used in update calculations.
	gasLimitBoundDivisor uint64 = 256
)

var (
	// Always Keccak256(RLP([])) as uncles are meaningless outside of PoW.
	uncleHash = types.CalcUncleHash(nil)

	// Block difficulty for in-turn signatures
	diffInTurn = big.NewInt(2)

	// Block difficulty for out-of-turn signatures
	diffNoTurn = big.NewInt(1)
)

// Bloom represents a 2048 bit bloom filter.
type Bloom [bloomByteLength]byte

// BytesToBloom converts a byte slice to a bloom filter.
// It panics if b is not of suitable size.
func BytesToBloom(b []byte) Bloom {
	var bloom Bloom
	bloom.SetBytes(b)
	return bloom
}

// SetBytes sets the content of b to the given bytes.
// It panics if d is not of suitable size.
func (b *Bloom) SetBytes(d []byte) {
	if len(b) < len(d) {
		panic(fmt.Sprintf("bloom bytes too big %d %d", len(b), len(d)))
	}
	copy(b[bloomByteLength-len(d):], d)
}

// A BlockNonce is a 64-bit hash which proves (combined with the
// mix-hash) that a sufficient amount of computation has been carried
// out on a block.
type BlockNonce [nonceByteLength]byte

// BlockNonce converts a byte slice to a bloom filter.
// It panics if b is not of suitable size.
func BytesToBlockNonce(b []byte) BlockNonce {
	var nonce BlockNonce
	nonce.SetBytes(b)
	return nonce
}

// SetBytes sets the content of b to the given bytes.
// It panics if d is not of suitable size.
func (b *BlockNonce) SetBytes(d []byte) {
	if len(b) < len(d) {
		panic(fmt.Sprintf("bloom bytes too big %d %d", len(b), len(d)))
	}
	copy(b[nonceByteLength-len(d):], d)
}

// BscHeader represents a block header in the Ethereum blockchain.
type BscHeader struct {
	ParentHash  common.Hash    `json:"parentHash"       gencodec:"required"`
	UncleHash   common.Hash    `json:"sha3Uncles"       gencodec:"required"`
	Coinbase    common.Address `json:"miner"            gencodec:"required"`
	Root        common.Hash    `json:"stateRoot"        gencodec:"required"`
	TxHash      common.Hash    `json:"transactionsRoot" gencodec:"required"`
	ReceiptHash common.Hash    `json:"receiptsRoot"     gencodec:"required"`
	Bloom       Bloom          `json:"logsBloom"        gencodec:"required"`
	Difficulty  *big.Int       `json:"difficulty"       gencodec:"required"`
	Number      *big.Int       `json:"number"           gencodec:"required"`
	GasLimit    uint64         `json:"gasLimit"         gencodec:"required"`
	GasUsed     uint64         `json:"gasUsed"          gencodec:"required"`
	Time        uint64         `json:"timestamp"        gencodec:"required"`
	Extra       []byte         `json:"extraData"        gencodec:"required"`
	MixDigest   common.Hash    `json:"mixHash"`
	Nonce       BlockNonce     `json:"nonce"`
}

func (h BscHeader) ToHeader() Header {
	return Header{
		ParentHash:  h.ParentHash[:],
		UncleHash:   h.UncleHash[:],
		Coinbase:    h.Coinbase[:],
		Root:        h.Root[:],
		TxHash:      h.TxHash[:],
		ReceiptHash: h.ReceiptHash[:],
		Bloom:       h.Bloom[:],
		Difficulty:  h.Difficulty.Uint64(),
		Height:      clienttypes.NewHeight(0, h.Number.Uint64()),
		GasLimit:    h.GasLimit,
		GasUsed:     h.GasUsed,
		Time:        h.Time,
		Extra:       h.Extra,
		MixDigest:   h.MixDigest[:],
		Nonce:       h.Nonce[:],
	}
}

// ProofAccount ...
type ProofAccount struct {
	Nonce    *big.Int
	Balance  *big.Int
	Storage  common.Hash
	Codehash common.Hash
}

func ParseValidators(extra []byte) ([][]byte, error) {
	validatorBytes := extra[extraVanity : len(extra)-extraSeal]
	if len(validatorBytes)%addressLength != 0 {
		return nil, sdkerrors.Wrap(ErrInvalidValidatorBytes, "(validatorsBytes % AddressLength) should bz zero")
	}
	n := len(validatorBytes) / addressLength
	result := make([][]byte, n)
	for i := 0; i < n; i++ {
		address := make([]byte, addressLength)
		copy(address, validatorBytes[i*addressLength:(i+1)*addressLength])
		result[i] = address
	}
	return result, nil
}
