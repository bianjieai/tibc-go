package types

import (
	fmt "fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

const (
	// Fixed number of extra-data prefix bytes reserved for signer vanity
	ExtraVanity = 32

	// Fixed number of extra-data suffix bytes reserved for signer seal
	ExtraSeal = 65

	// AddressLength is the expected length of the address
	AddressLength = 20

	// HashLength is the expected length of the hash
	HashLength = 32

	// BloomByteLength represents the number of bytes used in a header log bloom.
	BloomByteLength = 256

	// NonceByteLength represents the number of bytes used in a header log nonce.
	NonceByteLength = 8

	// The bound divisor of the gas limit, used in update calculations.
	GasLimitBoundDivisor uint64 = 256
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
type Bloom [BloomByteLength]byte

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
	copy(b[BloomByteLength-len(d):], d)
}

// A BlockNonce is a 64-bit hash which proves (combined with the
// mix-hash) that a sufficient amount of computation has been carried
// out on a block.
type BlockNonce [NonceByteLength]byte

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
	copy(b[NonceByteLength-len(d):], d)
}
