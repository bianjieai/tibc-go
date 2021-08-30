package types

import (
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

var _ exported.ConsensusState = (*ConsensusState)(nil)

func (m *ConsensusState) ClientType() string {
	panic("implement me")
}

func (m *ConsensusState) GetRoot() exported.Root {
	panic("implement me")
}

func (m *ConsensusState) GetTimestamp() uint64 {
	panic("implement me")
}

func (m *ConsensusState) ValidateBasic() error {
	panic("implement me")
}
