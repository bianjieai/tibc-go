package keeper_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/bianjieai/tibc-go/simapp"

	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/keeper"
	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	commitmenttypes "github.com/bianjieai/tibc-go/modules/tibc/core/23-commitment/types"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
	ibctmtypes "github.com/bianjieai/tibc-go/modules/tibc/light-clients/07-tendermint/types"
	ibctesting "github.com/bianjieai/tibc-go/modules/tibc/testing"
	ibctestingmock "github.com/bianjieai/tibc-go/modules/tibc/testing/mock"
)

const (
	testChainID          = "gaiahub-0"
	testChainIDRevision1 = "gaiahub-1"

	testChainName  = "tendermint-0"
	testChainName2 = "tendermint-1"
	testChainName3 = "tendermint-2"

	height = 5

	trustingPeriod time.Duration = time.Hour * 24 * 7 * 2
	ubdPeriod      time.Duration = time.Hour * 24 * 7 * 3
	maxClockDrift  time.Duration = time.Second * 10
)

var (
	testClientHeight          = types.NewHeight(0, 5)
	testClientHeightRevision1 = types.NewHeight(1, 5)
	newClientHeight           = types.NewHeight(1, 1)
)

type KeeperTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain

	cdc            codec.Marshaler
	ctx            sdk.Context
	keeper         *keeper.Keeper
	consensusState *ibctmtypes.ConsensusState
	header         *ibctmtypes.Header
	valSet         *tmtypes.ValidatorSet
	valSetHash     tmbytes.HexBytes
	privVal        tmtypes.PrivValidator
	now            time.Time
	past           time.Time

	queryClient types.QueryClient
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)

	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(1))

	isCheckTx := false
	suite.now = time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
	suite.past = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	now2 := suite.now.Add(time.Hour)
	app := simapp.Setup(isCheckTx)

	suite.cdc = app.AppCodec()
	suite.ctx = app.BaseApp.NewContext(isCheckTx, tmproto.Header{Height: height, ChainID: testChainName, Time: now2})
	suite.keeper = &app.IBCKeeper.ClientKeeper
	suite.privVal = ibctestingmock.NewPV()

	pubKey, err := suite.privVal.GetPubKey()
	suite.Require().NoError(err)

	testClientHeightMinus1 := types.NewHeight(0, height-1)

	validator := tmtypes.NewValidator(pubKey, 1)
	suite.valSet = tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})
	suite.valSetHash = suite.valSet.Hash()
	suite.header = suite.chainA.CreateTMClientHeader(testChainID, int64(testClientHeight.RevisionHeight), testClientHeightMinus1, now2, suite.valSet, suite.valSet, []tmtypes.PrivValidator{suite.privVal})
	suite.consensusState = ibctmtypes.NewConsensusState(suite.now, commitmenttypes.NewMerkleRoot([]byte("hash")), suite.valSetHash)

	var validators stakingtypes.Validators
	for i := 1; i < 11; i++ {
		privVal := ibctestingmock.NewPV()
		tmPk, err := privVal.GetPubKey()
		suite.Require().NoError(err)
		pk, err := cryptocodec.FromTmPubKeyInterface(tmPk)
		suite.Require().NoError(err)
		val, err := stakingtypes.NewValidator(sdk.ValAddress(pk.Address()), pk, stakingtypes.Description{})
		suite.Require().NoError(err)

		val.Status = stakingtypes.Bonded
		val.Tokens = sdk.NewInt(rand.Int63())
		validators = append(validators, val)

		hi := stakingtypes.NewHistoricalInfo(suite.ctx.BlockHeader(), validators)
		app.StakingKeeper.SetHistoricalInfo(suite.ctx, int64(i), &hi)
	}

	// add localhost client
	// revision := types.ParseChainID(suite.chainA.ChainID)
	// localHostClient := localhosttypes.NewClientState(
	// 	suite.chainA.ChainID, types.NewHeight(revision, uint64(suite.chainA.GetContext().BlockHeight())),
	// )
	// suite.chainA.App.IBCKeeper.ClientKeeper.SetClientState(suite.chainA.GetContext(), testClientName, localHostClient)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.IBCKeeper.ClientKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestSetClientState() {
	clientState := ibctmtypes.NewClientState(testChainID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.ZeroHeight(), commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false)
	suite.keeper.SetClientState(suite.ctx, testChainName, clientState)

	retrievedState, found := suite.keeper.GetClientState(suite.ctx, testChainName)
	suite.Require().True(found, "GetClientState failed")
	suite.Require().Equal(clientState, retrievedState, "Client states are not equal")
}

func (suite *KeeperTestSuite) TestSetClientConsensusState() {
	suite.keeper.SetClientConsensusState(suite.ctx, testChainName, testClientHeight, suite.consensusState)

	retrievedConsState, found := suite.keeper.GetClientConsensusState(suite.ctx, testChainName, testClientHeight)
	suite.Require().True(found, "GetConsensusState failed")

	tmConsState, ok := retrievedConsState.(*ibctmtypes.ConsensusState)
	suite.Require().True(ok)
	suite.Require().Equal(suite.consensusState, tmConsState, "ConsensusState not stored correctly")
}

func (suite KeeperTestSuite) TestGetAllGenesisClients() {
	clientIDs := []string{
		testChainName2, testChainName3, testChainName,
	}
	expClients := []exported.ClientState{
		ibctmtypes.NewClientState(testChainID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.ZeroHeight(), commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
		ibctmtypes.NewClientState(testChainID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.ZeroHeight(), commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
		ibctmtypes.NewClientState(testChainID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, types.ZeroHeight(), commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false),
	}

	expGenClients := make(types.IdentifiedClientStates, len(expClients))

	for i := range expClients {
		suite.chainA.App.IBCKeeper.ClientKeeper.SetClientState(suite.chainA.GetContext(), clientIDs[i], expClients[i])
		expGenClients[i] = types.NewIdentifiedClientState(clientIDs[i], expClients[i])
	}

	// add localhost client
	localHostClient, found := suite.chainA.App.IBCKeeper.ClientKeeper.GetClientState(suite.chainA.GetContext(), testChainName)
	suite.Require().True(found)
	expGenClients = append(expGenClients, types.NewIdentifiedClientState(testChainName, localHostClient))

	genClients := suite.chainA.App.IBCKeeper.ClientKeeper.GetAllGenesisClients(suite.chainA.GetContext())

	suite.Require().Equal(expGenClients.Sort(), genClients)
}

func (suite KeeperTestSuite) TestGetAllGenesisMetadata() {
	expectedGenMetadata := []types.IdentifiedGenesisMetadata{
		types.NewIdentifiedGenesisMetadata(
			"clientA",
			[]types.GenesisMetadata{
				types.NewGenesisMetadata(ibctmtypes.ProcessedTimeKey(types.NewHeight(0, 1)), []byte("foo")),
				types.NewGenesisMetadata(ibctmtypes.ProcessedTimeKey(types.NewHeight(0, 2)), []byte("bar")),
				types.NewGenesisMetadata(ibctmtypes.ProcessedTimeKey(types.NewHeight(0, 3)), []byte("baz")),
			},
		),
		types.NewIdentifiedGenesisMetadata(
			"clientB",
			[]types.GenesisMetadata{
				types.NewGenesisMetadata(ibctmtypes.ProcessedTimeKey(types.NewHeight(1, 100)), []byte("val1")),
				types.NewGenesisMetadata(ibctmtypes.ProcessedTimeKey(types.NewHeight(2, 300)), []byte("val2")),
			},
		),
	}

	genClients := []types.IdentifiedClientState{
		types.NewIdentifiedClientState("clientA", &ibctmtypes.ClientState{}), types.NewIdentifiedClientState("clientB", &ibctmtypes.ClientState{}),
	}

	suite.chainA.App.IBCKeeper.ClientKeeper.SetAllClientMetadata(suite.chainA.GetContext(), expectedGenMetadata)

	actualGenMetadata, err := suite.chainA.App.IBCKeeper.ClientKeeper.GetAllClientMetadata(suite.chainA.GetContext(), genClients)
	suite.Require().NoError(err, "get client metadata returned error unexpectedly")
	suite.Require().Equal(expectedGenMetadata, actualGenMetadata, "retrieved metadata is unexpected")
}

func (suite KeeperTestSuite) TestConsensusStateHelpers() {
	// initial setup
	clientState := ibctmtypes.NewClientState(testChainID, ibctmtypes.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, testClientHeight, commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath, false, false)

	suite.keeper.SetClientState(suite.ctx, testChainName, clientState)
	suite.keeper.SetClientConsensusState(suite.ctx, testChainName, testClientHeight, suite.consensusState)

	nextState := ibctmtypes.NewConsensusState(suite.now, commitmenttypes.NewMerkleRoot([]byte("next")), suite.valSetHash)

	testClientHeightPlus5 := types.NewHeight(0, height+5)

	header := suite.chainA.CreateTMClientHeader(testChainName, int64(testClientHeightPlus5.RevisionHeight), testClientHeight, suite.header.Header.Time.Add(time.Minute),
		suite.valSet, suite.valSet, []tmtypes.PrivValidator{suite.privVal})

	// mock update functionality
	clientState.LatestHeight = header.GetHeight().(types.Height)
	suite.keeper.SetClientConsensusState(suite.ctx, testChainName, header.GetHeight(), nextState)
	suite.keeper.SetClientState(suite.ctx, testChainName, clientState)

	latest, ok := suite.keeper.GetLatestClientConsensusState(suite.ctx, testChainName)
	suite.Require().True(ok)
	suite.Require().Equal(nextState, latest, "Latest client not returned correctly")
}

// 2 clients in total are created on chainA. The first client is updated so it contains an initial consensus state
// and a consensus state at the update height.
func (suite KeeperTestSuite) TestGetAllConsensusStates() {
	clientA, _ := suite.coordinator.SetupClients(suite.chainA, suite.chainB, exported.Tendermint)

	clientState := suite.chainA.GetClientState(clientA)
	expConsensusHeight0 := clientState.GetLatestHeight()
	consensusState0, ok := suite.chainA.GetConsensusState(clientA, expConsensusHeight0)
	suite.Require().True(ok)

	// update client to create a second consensus state
	err := suite.coordinator.UpdateClient(suite.chainA, suite.chainB, clientA, exported.Tendermint)
	suite.Require().NoError(err)

	clientState = suite.chainA.GetClientState(clientA)
	expConsensusHeight1 := clientState.GetLatestHeight()
	suite.Require().True(expConsensusHeight1.GT(expConsensusHeight0))
	consensusState1, ok := suite.chainA.GetConsensusState(clientA, expConsensusHeight1)
	suite.Require().True(ok)

	expConsensus := []exported.ConsensusState{
		consensusState0,
		consensusState1,
	}

	// create second client on chainA
	clientA2, _ := suite.coordinator.SetupClients(suite.chainA, suite.chainB, exported.Tendermint)
	clientState = suite.chainA.GetClientState(clientA2)

	expConsensusHeight2 := clientState.GetLatestHeight()
	consensusState2, ok := suite.chainA.GetConsensusState(clientA2, expConsensusHeight2)
	suite.Require().True(ok)

	expConsensus2 := []exported.ConsensusState{consensusState2}

	expConsensusStates := types.ClientsConsensusStates{
		types.NewClientConsensusStates(clientA, []types.ConsensusStateWithHeight{
			types.NewConsensusStateWithHeight(expConsensusHeight0.(types.Height), expConsensus[0]),
			types.NewConsensusStateWithHeight(expConsensusHeight1.(types.Height), expConsensus[1]),
		}),
		types.NewClientConsensusStates(clientA2, []types.ConsensusStateWithHeight{
			types.NewConsensusStateWithHeight(expConsensusHeight2.(types.Height), expConsensus2[0]),
		}),
	}.Sort()

	consStates := suite.chainA.App.IBCKeeper.ClientKeeper.GetAllConsensusStates(suite.chainA.GetContext())
	suite.Require().Equal(expConsensusStates, consStates, "%s \n\n%s", expConsensusStates, consStates)
}
