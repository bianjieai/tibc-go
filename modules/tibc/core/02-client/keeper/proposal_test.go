package keeper_test

import (
	"fmt"

	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	commitmenttypes "github.com/bianjieai/tibc-go/modules/tibc/core/23-commitment/types"
	ibctmtypes "github.com/bianjieai/tibc-go/modules/tibc/light-clients/07-tendermint/types"
	ibctesting "github.com/bianjieai/tibc-go/modules/tibc/testing"
)

func (suite KeeperTestSuite) TestHandleCreateClientProposal() {
	header := suite.chainA.CreateTMClientHeader(
		suite.chainA.ChainID, suite.chainA.CurrentHeader.Height,
		types.NewHeight(0, uint64(suite.chainA.CurrentHeader.Height-1)),
		suite.chainA.CurrentHeader.Time, suite.chainA.Vals,
		suite.chainA.Vals, suite.chainA.Signers,
	)
	testCases := []struct {
		msg      string
		malleate func()
	}{
		{
			"success, create new client",
			func() {
				clientState := ibctmtypes.NewClientState("test", ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.ZeroHeight(), commitmenttypes.GetSDKSpecs(), ibctesting.Prefix, 0)
				consensusState := ibctmtypes.NewConsensusState(
					header.GetTime(), commitmenttypes.NewMerkleRoot(header.Header.AppHash), header.Header.NextValidatorsHash,
				)
				var proposal *types.CreateClientProposal
				var err error
				proposal, err = types.NewCreateClientProposal("test", "test", "test", clientState, consensusState)
				suite.NoError(err)
				err = suite.chainA.App.TIBCKeeper.ClientKeeper.HandleCreateClientProposal(suite.chainA.GetContext(), proposal)
				suite.NoError(err)
			},
		},
		{
			"fail, A client for this chainname already exists",
			func() {
				clientState := ibctmtypes.NewClientState("test", ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.ZeroHeight(), commitmenttypes.GetSDKSpecs(), ibctesting.Prefix, 0)
				consensusState := ibctmtypes.NewConsensusState(
					header.GetTime(), commitmenttypes.NewMerkleRoot(header.Header.AppHash), header.Header.NextValidatorsHash,
				)
				var proposal *types.CreateClientProposal
				var err error
				proposal, err = types.NewCreateClientProposal("test", "test", "test", clientState, consensusState)
				suite.NoError(err)
				err = suite.chainA.App.TIBCKeeper.ClientKeeper.HandleCreateClientProposal(suite.chainA.GetContext(), proposal)
				suite.NoError(err)
				err = suite.chainA.App.TIBCKeeper.ClientKeeper.HandleCreateClientProposal(suite.chainA.GetContext(), proposal)
				suite.Error(err)
			},
		},
		{
			"success, get client and compare",
			func() {
				clientState := ibctmtypes.NewClientState("test", ibctmtypes.DefaultTrustLevel,
					trustingPeriod, ubdPeriod, maxClockDrift, types.NewHeight(0, 5),
					commitmenttypes.GetSDKSpecs(), ibctesting.Prefix, 0)
				consensusState := ibctmtypes.NewConsensusState(
					header.GetTime(), commitmenttypes.NewMerkleRoot(header.Header.AppHash), header.Header.NextValidatorsHash,
				)
				var proposal *types.CreateClientProposal
				var err error
				proposal, err = types.NewCreateClientProposal("test", "test", "test", clientState, consensusState)
				suite.NoError(err)
				err = suite.chainA.App.TIBCKeeper.ClientKeeper.HandleCreateClientProposal(suite.chainA.GetContext(), proposal)
				suite.NoError(err)
				client, _ := suite.chainA.App.TIBCKeeper.ClientKeeper.GetClientState(suite.chainA.GetContext(), "test")
				suite.Equal(clientState, client, "clientState not equal")
				consensus, _ := suite.chainA.App.TIBCKeeper.ClientKeeper.GetClientConsensusState(suite.chainA.GetContext(), "test", types.NewHeight(0, 5))
				suite.Equal(consensusState, consensus, "consensusState not equal")
			},
		},
	}
	for i, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s, %d/%d tests", tc.msg, i+1, len(testCases)), func() {
			suite.SetupTest() // reset the context
			tc.malleate()
		})
	}
}

func (suite KeeperTestSuite) TestHandleUpgradeClientProposal() {
	header := suite.chainA.CreateTMClientHeader(
		suite.chainA.ChainID, suite.chainA.CurrentHeader.Height,
		types.NewHeight(0, uint64(suite.chainA.CurrentHeader.Height-1)),
		suite.chainA.CurrentHeader.Time, suite.chainA.Vals,
		suite.chainA.Vals, suite.chainA.Signers,
	)
	testCases := []struct {
		msg      string
		malleate func()
	}{
		{
			"fail, client and consensus are not existing",
			func() {
				clientState := ibctmtypes.NewClientState("test", ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.ZeroHeight(), commitmenttypes.GetSDKSpecs(), ibctesting.Prefix, 0)
				consensusState := ibctmtypes.NewConsensusState(
					header.GetTime(), commitmenttypes.NewMerkleRoot(header.Header.AppHash), header.Header.NextValidatorsHash,
				)
				var proposal *types.UpgradeClientProposal
				var err error
				proposal, err = types.NewUpgradeClientProposal("test", "test", "test", clientState, consensusState)
				suite.NoError(err)
				err = suite.chainA.App.TIBCKeeper.ClientKeeper.HandleUpgradeClientProposal(suite.chainA.GetContext(), proposal)
				suite.Error(err)
			},
		},
		{
			"success, client and consensus are existing",
			func() {
				clientState := ibctmtypes.NewClientState("test", ibctmtypes.DefaultTrustLevel,
					trustingPeriod, ubdPeriod, maxClockDrift, types.NewHeight(0, 5),
					commitmenttypes.GetSDKSpecs(), ibctesting.Prefix, 0)
				consensusState := ibctmtypes.NewConsensusState(
					header.GetTime(), commitmenttypes.NewMerkleRoot(header.Header.AppHash), header.Header.NextValidatorsHash,
				)
				var proposal *types.CreateClientProposal
				var err error
				proposal, err = types.NewCreateClientProposal("test", "test", "test", clientState, consensusState)
				suite.NoError(err)
				err = suite.chainA.App.TIBCKeeper.ClientKeeper.HandleCreateClientProposal(suite.chainA.GetContext(), proposal)
				suite.NoError(err)
				clientState2 := ibctmtypes.NewClientState("test", ibctmtypes.DefaultTrustLevel,
					trustingPeriod*2, ubdPeriod, maxClockDrift, types.NewHeight(0, 6),
					commitmenttypes.GetSDKSpecs(), ibctesting.Prefix, 0)
				consensusState2 := ibctmtypes.NewConsensusState(
					header.GetTime().Add(1), commitmenttypes.NewMerkleRoot(header.Header.AppHash), header.Header.NextValidatorsHash,
				)
				var proposal2 *types.UpgradeClientProposal
				proposal2, err = types.NewUpgradeClientProposal("test", "test", "test", clientState2, consensusState2)
				suite.NoError(err)
				err = suite.chainA.App.TIBCKeeper.ClientKeeper.HandleUpgradeClientProposal(suite.chainA.GetContext(), proposal2)
				suite.NoError(err)

				// Check the consistency of clientState and consensusState
				client, _ := suite.chainA.App.TIBCKeeper.ClientKeeper.GetClientState(suite.chainA.GetContext(), "test")
				suite.Equal(clientState2, client, "clientState not equal")
				consensus, _ := suite.chainA.App.TIBCKeeper.ClientKeeper.GetClientConsensusState(suite.chainA.GetContext(), "test", types.NewHeight(0, 6))
				suite.Equal(consensusState2, consensus, "consensusState not equal")
			},
		},
	}
	for i, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s, %d/%d tests", tc.msg, i+1, len(testCases)), func() {
			suite.SetupTest() // reset the context
			tc.malleate()
		})
	}
}

func (suite KeeperTestSuite) TestHandleRegisterRelayerProposal() {
	header := suite.chainA.CreateTMClientHeader(
		suite.chainA.ChainID, suite.chainA.CurrentHeader.Height,
		types.NewHeight(0, uint64(suite.chainA.CurrentHeader.Height-1)),
		suite.chainA.CurrentHeader.Time, suite.chainA.Vals,
		suite.chainA.Vals, suite.chainA.Signers,
	)
	testCases := []struct {
		msg      string
		malleate func()
	}{
		{
			"success, exist client",
			func() {
				clientState := ibctmtypes.NewClientState("test", ibctmtypes.DefaultTrustLevel,
					trustingPeriod, ubdPeriod, maxClockDrift, types.NewHeight(0, 5),
					commitmenttypes.GetSDKSpecs(), ibctesting.Prefix, 0)
				consensusState := ibctmtypes.NewConsensusState(
					header.GetTime(), commitmenttypes.NewMerkleRoot(header.Header.AppHash), header.Header.NextValidatorsHash,
				)
				var proposalCreate *types.CreateClientProposal
				var err error
				proposalCreate, err = types.NewCreateClientProposal("test", "test", "test", clientState, consensusState)
				suite.NoError(err)
				err = suite.chainA.App.TIBCKeeper.ClientKeeper.HandleCreateClientProposal(suite.chainA.GetContext(), proposalCreate)
				suite.NoError(err)

				// set relayers
				relayers := []string{"xxx", "yyy"}
				relayerProposal := types.NewRegisterRelayerProposal("test", "test", "test", relayers)
				err = suite.chainA.App.TIBCKeeper.ClientKeeper.HandleRegisterRelayerProposal(suite.chainA.GetContext(), relayerProposal)
				suite.NoError(err)

				// get relayers and compare
				relayers2 := suite.chainA.App.TIBCKeeper.ClientKeeper.GetRelayers(suite.chainA.GetContext(), "test")
				suite.Equal(relayers, relayers2)
			},
		},
		{
			"success, no client",
			func() {
				relayers := []string{"xxx", "yyy"}
				relayerProposal := types.NewRegisterRelayerProposal("test", "test", "test", relayers)
				err := suite.chainA.App.TIBCKeeper.ClientKeeper.HandleRegisterRelayerProposal(suite.chainA.GetContext(), relayerProposal)
				suite.NoError(err)
			},
		},
		{
			"fail, no-existing client",
			func() {
				// set client "test"
				clientState := ibctmtypes.NewClientState("test", ibctmtypes.DefaultTrustLevel,
					trustingPeriod, ubdPeriod, maxClockDrift, types.NewHeight(0, 5),
					commitmenttypes.GetSDKSpecs(), ibctesting.Prefix, 0)
				consensusState := ibctmtypes.NewConsensusState(
					header.GetTime(), commitmenttypes.NewMerkleRoot(header.Header.AppHash), header.Header.NextValidatorsHash,
				)
				var proposalCreate *types.CreateClientProposal
				var err error
				proposalCreate, err = types.NewCreateClientProposal("test", "test", "test", clientState, consensusState)
				suite.NoError(err)
				err = suite.chainA.App.TIBCKeeper.ClientKeeper.HandleCreateClientProposal(suite.chainA.GetContext(), proposalCreate)
				suite.NoError(err)
			},
		},
	}
	for i, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s, %d/%d tests", tc.msg, i+1, len(testCases)), func() {
			suite.SetupTest() // reset the context
			tc.malleate()
		})
	}
}
