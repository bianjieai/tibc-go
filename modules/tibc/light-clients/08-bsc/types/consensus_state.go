package types

import (
	commitmenttypes "github.com/bianjieai/tibc-go/modules/tibc/core/23-commitment/types"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

var _ exported.ConsensusState = (*ConsensusState)(nil)

func (m *ConsensusState) ClientType() string {
	return exported.BSC
}

func (m *ConsensusState) GetRoot() exported.Root {
	return commitmenttypes.MerkleRoot{
		Hash: m.Root,
	}
}

func (m *ConsensusState) GetTimestamp() uint64 {
	return m.Timestamp
}

func (m *ConsensusState) ValidateBasic() error {
	return nil
}
