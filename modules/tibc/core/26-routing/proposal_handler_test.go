package routing_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	routing "github.com/bianjieai/tibc-go/modules/tibc/core/26-routing"
	routingtypes "github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
	tibctesting "github.com/bianjieai/tibc-go/modules/tibc/testing"
)

type RoutingTestSuite struct {
	suite.Suite

	coordinator *tibctesting.Coordinator

	chain *tibctesting.TestChain
}

func (suite *RoutingTestSuite) SetupTest() {
	suite.coordinator = tibctesting.NewCoordinator(suite.T(), 1)

	suite.chain = suite.coordinator.GetChain(tibctesting.GetChainID(0))
}

func TestRoutingTestSuite(t *testing.T) {
	suite.Run(t, new(RoutingTestSuite))
}

func (suite *RoutingTestSuite) TestNewSetRoutingRulesProposalHandler() {
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
				content, err = routingtypes.NewSetRoutingRulesProposal(tibctesting.Title, tibctesting.Description, []string{"source,dest,dgsbl"})
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
				content = distributiontypes.NewCommunityPoolSpendProposal(tibctesting.Title, tibctesting.Description, suite.chain.SenderAccount.GetAddress(), sdk.NewCoins(sdk.NewCoin("communityfunds", sdk.NewInt(10))))
			}, false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest() // reset

			tc.malleate()

			proposalHandler := routing.NewSetRoutingProposalHandler(suite.chain.App.TIBCKeeper.RoutingKeeper)

			err = proposalHandler(suite.chain.GetContext(), content)

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}

}
