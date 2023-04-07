package cli_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	routingtypes "github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
	clientcli "github.com/bianjieai/tibc-go/modules/tibc/core/client/cli"
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
		content govv1beta1.Content
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

			content, err = clienttypes.NewUpgradeClientProposal(ibctesting.Title, ibctesting.Description, path.EndpointB.Chain.ChainName, clientState, consensusState)
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

			content = clienttypes.NewRegisterRelayerProposal(ibctesting.Title, ibctesting.Description, path.EndpointB.Chain.ChainName, relayers)
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

			proposalHandler := clientcli.NewProposalHandler(suite.chainA.App.TIBCKeeper)

			err = proposalHandler(suite.chainA.GetContext(), content)

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *ClientTestSuite) TestNewSetRoutingRulesProposalHandler() {
	var (
		content govv1beta1.Content
		err     error
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"valid routing rules proposal",
			func() {
				content, err = routingtypes.NewSetRoutingRulesProposal(ibctesting.Title, ibctesting.Description, []string{"source,dest,dgsbl"})
				suite.Require().NoError(err)
			}, true,
		},
		{
			"nil proposal",
			func() {
				content = nil
			}, false,
		},
		{
			"unsupported proposal type",
			func() {
				content = distributiontypes.NewCommunityPoolSpendProposal(ibctesting.Title, ibctesting.Description, suite.chainA.SenderAccount.GetAddress(), sdk.NewCoins(sdk.NewCoin("communityfunds", sdk.NewInt(10))))
			}, false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest() // reset

			tc.malleate()

			proposalHandler := clientcli.NewProposalHandler(suite.chainA.App.TIBCKeeper)

			err = proposalHandler(suite.chainA.GetContext(), content)

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}

}
