package types_test

import (
	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
	"github.com/bianjieai/tibc-go/modules/tibc/light-clients/07-tendermint/types"
	tibctesting "github.com/bianjieai/tibc-go/modules/tibc/testing"
)

func (suite *TendermintTestSuite) TestGetConsensusState() {
	var (
		height exported.Height
		path   *tibctesting.Path
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{{
		"success", func() {}, true,
	}, {
		"consensus state not found", func() {
			// use height with no consensus state set
			height = height.(clienttypes.Height).Increment()
		},
		false,
	}, {
		"not a consensus state interface", func() {
			// marshal an empty client state and set as consensus state
			store := suite.chainA.App.TIBCKeeper.ClientKeeper.ClientStore(suite.chainA.GetContext(), path.EndpointB.ChainName)
			clientStateBz := suite.chainA.App.TIBCKeeper.ClientKeeper.MustMarshalClientState(&types.ClientState{})
			store.Set(host.ConsensusStateKey(height), clientStateBz)
		},
		false,
	}}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest()

			path = tibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)

			clientState := path.EndpointA.GetClientState()
			height = clientState.GetLatestHeight()

			tc.malleate() // change vars as necessary

			store := suite.chainA.App.TIBCKeeper.ClientKeeper.ClientStore(suite.chainA.GetContext(), path.EndpointB.ChainName)
			consensusState, err := types.GetConsensusState(store, suite.chainA.Codec, height)

			if tc.expPass {
				suite.Require().NoError(err)
				expConsensusState, found := suite.chainA.GetConsensusState(path.EndpointB.ChainName, height)
				suite.Require().True(found)
				suite.Require().Equal(expConsensusState, consensusState)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(consensusState)
			}
		})
	}
}

func (suite *TendermintTestSuite) TestGetProcessedTime() {
	// setup
	path := tibctesting.NewPath(suite.chainA, suite.chainB)

	suite.coordinator.UpdateTime()
	// coordinator increments time before creating client
	expectedTime := suite.chainA.CurrentHeader.Time.Add(tibctesting.TimeIncrement)

	// Verify ProcessedTime on CreateClient
	err := path.EndpointA.CreateClient()
	suite.Require().NoError(err)

	clientState := path.EndpointA.GetClientState()
	height := clientState.GetLatestHeight()

	store := path.EndpointA.ClientStore()
	actualTime, ok := types.GetProcessedTime(store, height)
	suite.Require().True(ok, "could not retrieve processed time for stored consensus state")
	suite.Require().Equal(uint64(expectedTime.UnixNano()), actualTime, "retrieved processed time is not expected value")

	suite.coordinator.UpdateTime()
	// coordinator increments time before updating client
	expectedTime = suite.chainA.CurrentHeader.Time.Add(tibctesting.TimeIncrement)

	// Verify ProcessedTime on UpdateClient
	err = path.EndpointA.UpdateClient()
	suite.Require().NoError(err)

	clientState = path.EndpointA.GetClientState()
	height = clientState.GetLatestHeight()

	store = path.EndpointA.ClientStore()
	actualTime, ok = types.GetProcessedTime(store, height)
	suite.Require().True(ok, "could not retrieve processed time for stored consensus state")
	suite.Require().Equal(uint64(expectedTime.UnixNano()), actualTime, "retrieved processed time is not expected value")

	// try to get processed time for height that doesn't exist in store
	_, ok = types.GetProcessedTime(store, clienttypes.NewHeight(1, 1))
	suite.Require().False(ok, "retrieved processed time for a non-existent consensus state")
}
