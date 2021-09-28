package client_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	client "github.com/bianjieai/tibc-go/modules/tibc/core/02-client"
	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	ibctesting "github.com/bianjieai/tibc-go/modules/tibc/testing"
)

type ClientTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain
}

func (suite *ClientTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)

	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(1))
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

func (suite *ClientTestSuite) TestNewClientUpdateProposalHandler() {
	var (
		content govtypes.Content
		err     error
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{{
		"valid create client proposal",
		func() {
			// setup testing conditions
			path := ibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)
			clientState := path.EndpointA.GetClientState()
			consensusState := path.EndpointA.GetConsensusState(clientState.GetLatestHeight())

			content, err = clienttypes.NewCreateClientProposal(ibctesting.Title, ibctesting.Description, "test-chain-name", clientState, consensusState)
			suite.Require().NoError(err)
		},
		true,
	}, {
		"valid create client proposal",
		func() {
			// setup testing conditions
			path := ibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)
			clientState := path.EndpointA.GetClientState()
			consensusState := path.EndpointA.GetConsensusState(clientState.GetLatestHeight())

			content, err = clienttypes.NewUpgradeClientProposal(ibctesting.Title, ibctesting.Description, path.EndpointB.ChainName, clientState, consensusState)
			suite.Require().NoError(err)
		},
		true,
	}, {
		"valid create client proposal",
		func() {

			// setup testing conditions
			path := ibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)

			relayers := []string{
				suite.chainB.SenderAccount.GetAddress().String(),
			}

			content = clienttypes.NewRegisterRelayerProposal(ibctesting.Title, ibctesting.Description, path.EndpointB.ChainName, relayers)
			suite.Require().NoError(err)
		},
		true,
	}, {
		"nil proposal",
		func() {
			content = nil
		},
		false,
	}, {
		"unsupported proposal type",
		func() {
			content = distributiontypes.NewCommunityPoolSpendProposal(ibctesting.Title, ibctesting.Description, suite.chainA.SenderAccount.GetAddress(), sdk.NewCoins(sdk.NewCoin("communityfunds", sdk.NewInt(10))))
		},
		false,
	}}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset

			tc.malleate()

			proposalHandler := client.NewClientProposalHandler(suite.chainA.App.TIBCKeeper.ClientKeeper)

			err = proposalHandler(suite.chainA.GetContext(), content)

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
