package tibc_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"

	tibc "github.com/bianjieai/tibc-go/modules/tibc/core"
	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	commitmenttypes "github.com/bianjieai/tibc-go/modules/tibc/core/23-commitment/types"
	"github.com/bianjieai/tibc-go/modules/tibc/core/types"
	tibctmtypes "github.com/bianjieai/tibc-go/modules/tibc/light-clients/07-tendermint/types"
	tibctesting "github.com/bianjieai/tibc-go/modules/tibc/testing"
	"github.com/bianjieai/tibc-go/simapp"
)

const (
	chainName  = "07-tendermint-0"
	chainName2 = "07-tendermin-1"
)

var clientHeight = clienttypes.NewHeight(0, 10)

type TIBCTestSuite struct {
	suite.Suite

	coordinator *tibctesting.Coordinator

	chainA *tibctesting.TestChain
	chainB *tibctesting.TestChain
}

// SetupTest creates a coordinator with 2 test chains.
func (suite *TIBCTestSuite) SetupTest() {
	suite.coordinator = tibctesting.NewCoordinator(suite.T(), 2)

	suite.chainA = suite.coordinator.GetChain(tibctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(tibctesting.GetChainID(1))
}

func TestIBCTestSuite(t *testing.T) {
	suite.Run(t, new(TIBCTestSuite))
}

func (suite *TIBCTestSuite) TestValidateGenesis() {
	header := suite.chainA.CreateTMClientHeader(
		suite.chainA.ChainID, suite.chainA.CurrentHeader.Height,
		clienttypes.NewHeight(0, uint64(suite.chainA.CurrentHeader.Height-1)),
		suite.chainA.CurrentHeader.Time, suite.chainA.Vals, suite.chainA.Vals,
		suite.chainA.Signers,
	)

	testCases := []struct {
		name     string
		genState *types.GenesisState
		expPass  bool
	}{{
		name:     "default",
		genState: types.DefaultGenesisState(),
		expPass:  true,
	}, {
		name: "valid genesis",
		genState: &types.GenesisState{
			ClientGenesis: clienttypes.NewGenesisState(
				[]clienttypes.IdentifiedClientState{
					clienttypes.NewIdentifiedClientState(
						chainName, tibctmtypes.NewClientState(
							suite.chainA.ChainID, tibctmtypes.DefaultTrustLevel,
							tibctesting.TrustingPeriod, tibctesting.UnbondingPeriod,
							tibctesting.MaxClockDrift, clientHeight,
							commitmenttypes.GetSDKSpecs(), tibctesting.Prefix, 0,
						),
					),
				},
				[]clienttypes.ClientConsensusStates{
					clienttypes.NewClientConsensusStates(
						chainName,
						[]clienttypes.ConsensusStateWithHeight{
							clienttypes.NewConsensusStateWithHeight(
								header.GetHeight().(clienttypes.Height),
								tibctmtypes.NewConsensusState(
									header.GetTime(), commitmenttypes.NewMerkleRoot(header.Header.AppHash), header.Header.NextValidatorsHash,
								),
							),
						},
					),
				},
				[]clienttypes.IdentifiedGenesisMetadata{
					clienttypes.NewIdentifiedGenesisMetadata(
						chainName,
						[]clienttypes.GenesisMetadata{
							clienttypes.NewGenesisMetadata([]byte("key1"), []byte("val1")),
							clienttypes.NewGenesisMetadata([]byte("key2"), []byte("val2")),
						},
					),
				},
				chainName2,
			),
		},
		expPass: true,
	}, {
		name: "invalid client genesis",
		genState: &types.GenesisState{
			ClientGenesis: clienttypes.NewGenesisState(
				[]clienttypes.IdentifiedClientState{
					clienttypes.NewIdentifiedClientState(
						chainName, tibctmtypes.NewClientState(
							suite.chainA.ChainID, tibctmtypes.DefaultTrustLevel,
							tibctesting.TrustingPeriod, tibctesting.UnbondingPeriod,
							tibctesting.MaxClockDrift, clientHeight,
							commitmenttypes.GetSDKSpecs(), tibctesting.Prefix, 0,
						),
					),
				},
				nil,
				[]clienttypes.IdentifiedGenesisMetadata{
					clienttypes.NewIdentifiedGenesisMetadata(
						chainName,
						[]clienttypes.GenesisMetadata{
							clienttypes.NewGenesisMetadata([]byte(""), []byte("val1")),
							clienttypes.NewGenesisMetadata([]byte("key2"), []byte("")),
						},
					),
				},
				chainName2,
			),
		},
		expPass: false,
	},
	}

	for _, tc := range testCases {
		tc := tc
		err := tc.genState.Validate()
		if tc.expPass {
			suite.Require().NoError(err, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
	}
}

func (suite *TIBCTestSuite) TestInitGenesis() {
	header := suite.chainA.CreateTMClientHeader(
		suite.chainA.ChainID, suite.chainA.CurrentHeader.Height,
		clienttypes.NewHeight(0, uint64(suite.chainA.CurrentHeader.Height-1)),
		suite.chainA.CurrentHeader.Time, suite.chainA.Vals,
		suite.chainA.Vals, suite.chainA.Signers,
	)

	testCases := []struct {
		name     string
		genState *types.GenesisState
	}{{
		name:     "default",
		genState: types.DefaultGenesisState(),
	}, {
		name: "valid genesis",
		genState: &types.GenesisState{
			ClientGenesis: clienttypes.NewGenesisState(
				[]clienttypes.IdentifiedClientState{
					clienttypes.NewIdentifiedClientState(
						chainName, tibctmtypes.NewClientState(
							suite.chainA.ChainID, tibctmtypes.DefaultTrustLevel,
							tibctesting.TrustingPeriod, tibctesting.UnbondingPeriod,
							tibctesting.MaxClockDrift, clientHeight,
							commitmenttypes.GetSDKSpecs(), tibctesting.Prefix, 0,
						),
					),
				},
				[]clienttypes.ClientConsensusStates{
					clienttypes.NewClientConsensusStates(
						chainName,
						[]clienttypes.ConsensusStateWithHeight{
							clienttypes.NewConsensusStateWithHeight(
								header.GetHeight().(clienttypes.Height),
								tibctmtypes.NewConsensusState(
									header.GetTime(), commitmenttypes.NewMerkleRoot(header.Header.AppHash), header.Header.NextValidatorsHash,
								),
							),
						},
					),
				},
				[]clienttypes.IdentifiedGenesisMetadata{
					clienttypes.NewIdentifiedGenesisMetadata(
						chainName,
						[]clienttypes.GenesisMetadata{
							clienttypes.NewGenesisMetadata([]byte("key1"), []byte("val1")),
							clienttypes.NewGenesisMetadata([]byte("key2"), []byte("val2")),
						},
					),
				},
				chainName2,
			),
		},
	}}

	for _, tc := range testCases {
		app := simapp.Setup(false)

		suite.NotPanics(func() {
			tibc.InitGenesis(app.BaseApp.NewContext(false, tmproto.Header{Height: 1}), *app.TIBCKeeper, true, tc.genState)
		})
	}
}

func (suite *TIBCTestSuite) TestExportGenesis() {
	testCases := []struct {
		msg      string
		malleate func()
	}{{
		"success",
		func() {
			// creates clients
			suite.coordinator.SetupClients(tibctesting.NewPath(suite.chainA, suite.chainB))
			// create extra clients
			suite.coordinator.SetupClients(tibctesting.NewPath(suite.chainA, suite.chainB))
			suite.coordinator.SetupClients(tibctesting.NewPath(suite.chainA, suite.chainB))
		},
	}}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest()

			tc.malleate()

			var gs *types.GenesisState
			suite.NotPanics(func() {
				gs = tibc.ExportGenesis(suite.chainA.GetContext(), *suite.chainA.App.TIBCKeeper)
			})

			// init genesis based on export
			suite.NotPanics(func() {
				tibc.InitGenesis(suite.chainA.GetContext(), *suite.chainA.App.TIBCKeeper, true, gs)
			})

			suite.NotPanics(func() {
				cdc := codec.NewProtoCodec(suite.chainA.App.InterfaceRegistry())
				genState := cdc.MustMarshalJSON(gs)
				cdc.MustUnmarshalJSON(genState, gs)
			})

			// init genesis based on marshal and unmarshal
			suite.NotPanics(func() {
				tibc.InitGenesis(suite.chainA.GetContext(), *suite.chainA.App.TIBCKeeper, true, gs)
			})
		})
	}
}
