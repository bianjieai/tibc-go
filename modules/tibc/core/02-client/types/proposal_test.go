package types_test

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	ibctmtypes "github.com/bianjieai/tibc-go/modules/tibc/light-clients/07-tendermint/types"
	ibctesting "github.com/bianjieai/tibc-go/modules/tibc/testing"
)

func (suite *TypesTestSuite) TestNewCreateClientProposal() {
	p, err := types.NewCreateClientProposal(ibctesting.Title, ibctesting.Description, chainName, &ibctmtypes.ClientState{}, &ibctmtypes.ConsensusState{})
	suite.Require().NoError(err)
	suite.Require().NotNil(p)

	p, err = types.NewCreateClientProposal(ibctesting.Title, ibctesting.Description, chainName, nil, nil)
	suite.Require().Error(err)
	suite.Require().Nil(p)
}

// tests a client update proposal can be marshaled and unmarshaled, and the
// client state can be unpacked
func (suite *TypesTestSuite) TestMarshalCreateClientProposalProposal() {
	// create proposal
	proposal, err := types.NewCreateClientProposal("update IBC client", "description", "client-id", nil, nil)
	suite.Require().NoError(err)

	// create codec
	ir := codectypes.NewInterfaceRegistry()
	types.RegisterInterfaces(ir)
	govtypes.RegisterInterfaces(ir)
	ibctmtypes.RegisterInterfaces(ir)
	cdc := codec.NewProtoCodec(ir)

	// marshal message
	bz, err := cdc.MarshalJSON(proposal)
	suite.Require().NoError(err)

	// unmarshal proposal
	newProposal := &types.CreateClientProposal{}
	err = cdc.UnmarshalJSON(bz, newProposal)
	suite.Require().NoError(err)
}
