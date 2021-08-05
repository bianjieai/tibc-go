package client_test

// import (
// 	"testing"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
// 	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
// 	"github.com/stretchr/testify/suite"

// 	client "github.com/bianjieai/tibc-go/modules/tibc/core/02-client"
// 	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
// 	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
// 	ibctmtypes "github.com/bianjieai/tibc-go/modules/tibc/light-clients/07-tendermint/types"
// 	ibctesting "github.com/bianjieai/tibc-go/modules/tibc/testing"
// )

// type ClientTestSuite struct {
// 	suite.Suite

// 	coordinator *ibctesting.Coordinator

// 	chainA *ibctesting.TestChain
// 	chainB *ibctesting.TestChain
// }

// func (suite *ClientTestSuite) SetupTest() {
// 	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)

// 	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(0))
// 	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(1))
// }

// func TestClientTestSuite(t *testing.T) {
// 	suite.Run(t, new(ClientTestSuite))
// }

// func (suite *ClientTestSuite) TestNewClientUpdateProposalHandler() {
// 	var (
// 		content govtypes.Content
// 		err     error
// 	)

// 	testCases := []struct {
// 		name     string
// 		malleate func()
// 		expPass  bool
// 	}{
// 		{
// 			"valid update client proposal", func() {
// 				clientA, _ := suite.coordinator.SetupClients(suite.chainA, suite.chainB, exported.Tendermint)
// 				clientState := suite.chainA.GetClientState(clientA)

// 				tmClientState, ok := clientState.(*ibctmtypes.ClientState)
// 				suite.Require().True(ok)
// 				tmClientState.AllowUpdateAfterMisbehaviour = true
// 				tmClientState.FrozenHeight = tmClientState.LatestHeight
// 				suite.chainA.App.IBCKeeper.ClientKeeper.SetClientState(suite.chainA.GetContext(), clientA, tmClientState)

// 				// use next header for chainB to update the client on chainA
// 				// header, err := suite.chainA.ConstructUpdateTMClientHeader(suite.chainB, clientA)
// 				// suite.Require().NoError(err)

// 				content, err = clienttypes.NewCreateClientProposal(ibctesting.Title, ibctesting.Description, clientA, nil, nil)
// 				suite.Require().NoError(err)
// 			}, true,
// 		},
// 		{
// 			"nil proposal", func() {
// 				content = nil
// 			}, false,
// 		},
// 		{
// 			"unsupported proposal type", func() {
// 				content = distributiontypes.NewCommunityPoolSpendProposal(ibctesting.Title, ibctesting.Description, suite.chainA.SenderAccount.GetAddress(), sdk.NewCoins(sdk.NewCoin("communityfunds", sdk.NewInt(10))))
// 			}, false,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		tc := tc

// 		suite.Run(tc.name, func() {
// 			suite.SetupTest() // reset

// 			tc.malleate()

// 			proposalHandler := client.NewClientProposalHandler(suite.chainA.App.IBCKeeper.ClientKeeper)

// 			err = proposalHandler(suite.chainA.GetContext(), content)

// 			if tc.expPass {
// 				suite.Require().NoError(err)
// 			} else {
// 				suite.Require().Error(err)
// 			}
// 		})
// 	}

// }
