package types_test

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	tibctmtypes "github.com/bianjieai/tibc-go/modules/tibc/light-clients/07-tendermint/types"
	tibctesting "github.com/bianjieai/tibc-go/modules/tibc/testing"
)

func (suite *TypesTestSuite) TestNewCreateClientProposal() {
	p, err := types.NewCreateClientProposal(tibctesting.Title, tibctesting.Description, chainName, &tibctmtypes.ClientState{}, &tibctmtypes.ConsensusState{})
	suite.Require().NoError(err)
	suite.Require().NotNil(p)

	p, err = types.NewCreateClientProposal(tibctesting.Title, tibctesting.Description, chainName, nil, nil)
	suite.Require().Error(err)
	suite.Require().Nil(p)
}

// tests a client update proposal can be marshaled and unmarshaled, and the
// client state can be unpacked
func (suite *TypesTestSuite) TestMarshalCreateClientProposalProposal() {
	path := tibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupClients(path)

	clientState := path.EndpointA.GetClientState()
	consensusState := path.EndpointA.GetConsensusState(clientState.GetLatestHeight())
	// create proposal
	proposal, err := types.NewCreateClientProposal("update TIBC client", "description", "chain-name", clientState, consensusState)
	suite.Require().NoError(err)

	// create codec
	ir := codectypes.NewInterfaceRegistry()
	types.RegisterInterfaces(ir)
	govtypes.RegisterInterfaces(ir)
	tibctmtypes.RegisterInterfaces(ir)
	cdc := codec.NewProtoCodec(ir)

	// marshal message
	bz, err := cdc.MarshalJSON(proposal)
	suite.Require().NoError(err)

	// unmarshal proposal
	newProposal := &types.CreateClientProposal{}
	err = cdc.UnmarshalJSON(bz, newProposal)
	suite.Require().NoError(err)
}
