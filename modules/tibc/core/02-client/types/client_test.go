package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	tibctesting "github.com/bianjieai/tibc-go/modules/tibc/testing"
)

func (suite *TypesTestSuite) TestMarshalConsensusStateWithHeight() {
	var cswh types.ConsensusStateWithHeight

	testCases := []struct {
		name     string
		malleate func()
	}{{
		"tendermint client",
		func() {
			// setup testing conditions
			path := tibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)
			clientState := path.EndpointA.GetClientState()
			consensusState := path.EndpointA.GetConsensusState(clientState.GetLatestHeight())
			cswh = types.NewConsensusStateWithHeight(clientState.GetLatestHeight().(types.Height), consensusState)
		},
	}}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest()

			tc.malleate()

			cdc := suite.chainA.App.AppCodec()

			// marshal message
			bz, err := cdc.MarshalJSON(&cswh)
			suite.Require().NoError(err)

			// unmarshal message
			newCswh := &types.ConsensusStateWithHeight{}
			err = cdc.UnmarshalJSON(bz, newCswh)
			suite.Require().NoError(err)
		})
	}
}

func TestValidateClientType(t *testing.T) {
	testCases := []struct {
		name       string
		clientType string
		expPass    bool
	}{
		{"valid", "tendermint", true},
		{"valid solomachine", "solomachine-v1", true},
		// {"too large", "tenderminttenderminttenderminttenderminttendermintt", false},
		{"too short", "t", false},
		{"blank id", "               ", false},
		{"empty id", "", false},
		// {"ends with dash", "tendermint-", false},
	}

	for _, tc := range testCases {
		err := types.ValidateClientType(tc.clientType)
		if tc.expPass {
			require.NoError(t, err, tc.name)
		} else {
			require.Error(t, err, tc.name)
		}
	}
}
