package types

import (
	"math/big"
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
	"github.com/ethereum/go-ethereum/common"
)

var _ exported.Header = (*Header)(nil)

func (h Header) ClientType() string {
	return exported.BSC
}

func (h Header) GetHeight() exported.Height {
	return h.Height
}

// Hash returns the block hash of the header, which is simply the keccak256 hash of its
// RLP encoding.
func (h *Header) Hash() common.Hash {
	return rlpHash(h.ToBscHeader())
}

func (h Header) ValidateBasic() error {
	number := h.Height.RevisionHeight

	// Don't waste time checking blocks from the future
	if h.Time > uint64(time.Now().Unix()) {
		return sdkerrors.Wrap(ErrFutureBlock, "header Height")
	}
	// Check that the extra-data contains the vanity, validators and signature.
	if len(h.Extra) < ExtraVanity {
		return sdkerrors.Wrap(ErrMissingVanity, "header Extra")
	}
	if len(h.Extra) < ExtraVanity+ExtraSeal {
		return sdkerrors.Wrap(ErrMissingSignature, "header Extra")
	}

	// Ensure that the mix digest is zero as we don't have fork protection currently
	if common.BytesToHash(h.MixDigest) != (common.Hash{}) {
		return sdkerrors.Wrap(ErrInvalidMixDigest, "header MixDigest")
	}
	// Ensure that the block doesn't contain any uncles which are meaningless in PoA
	if common.BytesToHash(h.UncleHash) != uncleHash {
		return sdkerrors.Wrap(ErrInvalidUncleHash, "header UncleHash")
	}
	// Ensure that the block's difficulty is meaningful (may not be correct at this point)
	if number > 0 {
		if h.Difficulty == 0 {
			return sdkerrors.Wrap(ErrInvalidDifficulty, "header Difficulty")
		}
	}
	return nil
}

func (h Header) ToBscHeader() bscHeader {
	return bscHeader{
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

// bscHeader represents a block header in the Ethereum blockchain.
type bscHeader struct {
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
