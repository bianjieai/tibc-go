package types

import (
	"math/big"

	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	paramsIndex  = 104
	paramsLenght = 32
)

type ProofKeyConstructor struct {
	sourceChain string
	destChain   string
	sequence    uint64
}

func NewProofKeyConstructor(sourceChain string, destChain string, sequence uint64) ProofKeyConstructor {
	return ProofKeyConstructor{
		sourceChain: sourceChain,
		destChain:   destChain,
		sequence:    sequence,
	}
}

func (k ProofKeyConstructor) GetPacketCommitmentProofKey() []byte {
	hash := crypto.Keccak256Hash(
		host.PacketCommitmentKey(k.sourceChain, k.destChain, k.sequence),
		common.LeftPadBytes(big.NewInt(paramsIndex).Bytes(), paramsLenght),
	)
	return hash.Bytes()
}

func (k ProofKeyConstructor) GetAckProofKey() []byte {
	hash := crypto.Keccak256Hash(
		host.PacketAcknowledgementKey(k.sourceChain, k.destChain, k.sequence),
		common.LeftPadBytes(big.NewInt(paramsIndex).Bytes(), paramsLenght),
	)
	return hash.Bytes()
}

func (k ProofKeyConstructor) GetCleanPacketCommitmentProofKey() []byte {
	hash := crypto.Keccak256Hash(
		host.CleanPacketCommitmentKey(k.sourceChain, k.destChain),
		common.LeftPadBytes(big.NewInt(paramsIndex).Bytes(), paramsLenght),
	)
	return hash.Bytes()
}
