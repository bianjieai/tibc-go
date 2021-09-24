package types_test

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	commitmenttypes "github.com/bianjieai/tibc-go/modules/tibc/core/23-commitment/types"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
	tibctmtypes "github.com/bianjieai/tibc-go/modules/tibc/light-clients/07-tendermint/types"
	tibctesting "github.com/bianjieai/tibc-go/modules/tibc/testing"
)

type caseAny struct {
	name    string
	any     *codectypes.Any
	expPass bool
}

func (suite *TypesTestSuite) TestPackClientState() {

	testCases := []struct {
		name        string
		clientState exported.ClientState
		expPass     bool
	}{{
		"tendermint client",
		tibctmtypes.NewClientState(
			chainID, tibctesting.DefaultTrustLevel, tibctesting.TrustingPeriod,
			tibctesting.UnbondingPeriod, tibctesting.MaxClockDrift, clientHeight,
			commitmenttypes.GetSDKSpecs(), tibctesting.Prefix, 0,
		),
		true,
	}, {
		"nil",
		nil,
		false,
	}}

	testCasesAny := []caseAny{}

	for _, tc := range testCases {
		clientAny, err := types.PackClientState(tc.clientState)
		if tc.expPass {
			suite.Require().NoError(err, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}

		testCasesAny = append(testCasesAny, caseAny{tc.name, clientAny, tc.expPass})
	}

	for i, tc := range testCasesAny {
		cs, err := types.UnpackClientState(tc.any)
		if tc.expPass {
			suite.Require().NoError(err, tc.name)
			suite.Require().Equal(testCases[i].clientState, cs, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
	}
}

func (suite *TypesTestSuite) TestPackConsensusState() {
	testCases := []struct {
		name           string
		consensusState exported.ConsensusState
		expPass        bool
	}{{
		"tendermint consensus",
		suite.chainA.LastHeader.ConsensusState(),
		true,
	}, {
		"nil",
		nil,
		false,
	}}

	testCasesAny := []caseAny{}

	for _, tc := range testCases {
		clientAny, err := types.PackConsensusState(tc.consensusState)
		if tc.expPass {
			suite.Require().NoError(err, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
		testCasesAny = append(testCasesAny, caseAny{tc.name, clientAny, tc.expPass})
	}

	for i, tc := range testCasesAny {
		cs, err := types.UnpackConsensusState(tc.any)
		if tc.expPass {
			suite.Require().NoError(err, tc.name)
			suite.Require().Equal(testCases[i].consensusState, cs, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
	}
}

func (suite *TypesTestSuite) TestPackHeader() {
	testCases := []struct {
		name    string
		header  exported.Header
		expPass bool
	}{{
		"tendermint header",
		suite.chainA.LastHeader,
		true,
	}, {
		"nil",
		nil,
		false,
	}}

	testCasesAny := []caseAny{}
	for _, tc := range testCases {
		clientAny, err := types.PackHeader(tc.header)
		if tc.expPass {
			suite.Require().NoError(err, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
		testCasesAny = append(testCasesAny, caseAny{tc.name, clientAny, tc.expPass})
	}

	for i, tc := range testCasesAny {
		cs, err := types.UnpackHeader(tc.any)
		if tc.expPass {
			suite.Require().NoError(err, tc.name)
			suite.Require().Equal(testCases[i].header, cs, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
	}
}
