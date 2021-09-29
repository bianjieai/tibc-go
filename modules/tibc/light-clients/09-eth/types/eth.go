package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
)

//EthHeader represents a block header in the Ethereum blockchain.
type EthHeader struct {
	ParentHash  common.Hash      `json:"parentHash"       gencodec:"required"`
	UncleHash   common.Hash      `json:"sha3Uncles"       gencodec:"required"`
	Coinbase    common.Address   `json:"miner"            gencodec:"required"`
	Root        common.Hash      `json:"stateRoot"        gencodec:"required"`
	TxHash      common.Hash      `json:"transactionsRoot" gencodec:"required"`
	ReceiptHash common.Hash      `json:"receiptsRoot"     gencodec:"required"`
	Bloom       types.Bloom      `json:"logsBloom"        gencodec:"required"`
	Difficulty  *big.Int         `json:"difficulty"       gencodec:"required"`
	Number      *big.Int         `json:"number"           gencodec:"required"`
	GasLimit    uint64           `json:"gasLimit"         gencodec:"required"`
	GasUsed     uint64           `json:"gasUsed"          gencodec:"required"`
	Time        uint64           `json:"timestamp"        gencodec:"required"`
	Extra       []byte           `json:"extraData"        gencodec:"required"`
	MixDigest   common.Hash      `json:"mixHash"`
	Nonce       types.BlockNonce `json:"nonce"`

	// BaseFee was added by EIP-1559 and is ignored in legacy headers.
	BaseFee *big.Int `json:"baseFeePerGas" rlp:"optional"`
}

func (h EthHeader) ToHeader() Header {
	return Header{
		ParentHash:  h.ParentHash[:],
		UncleHash:   h.UncleHash[:],
		Coinbase:    h.Coinbase[:],
		Root:        h.Root[:],
		TxHash:      h.TxHash[:],
		ReceiptHash: h.ReceiptHash[:],
		Bloom:       h.Bloom[:],
		Difficulty:  h.Difficulty.String(),
		Height:      clienttypes.NewHeight(0, h.Number.Uint64()),
		GasLimit:    h.GasLimit,
		GasUsed:     h.GasUsed,
		Time:        h.Time,
		Extra:       h.Extra,
		MixDigest:   h.MixDigest[:],
		Nonce:       h.Nonce.Uint64(),
		BaseFee:     h.BaseFee.String(),
	}
}

// ProofAccount ...
type ProofAccount struct {
	Nonce    *big.Int
	Balance  *big.Int
	Storage  common.Hash
	Codehash common.Hash
}
