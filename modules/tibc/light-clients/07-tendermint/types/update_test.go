package types_test

import (
	"fmt"

	"time"

	tmtypes "github.com/tendermint/tendermint/types"

	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	commitmenttypes "github.com/bianjieai/tibc-go/modules/tibc/core/23-commitment/types"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
	types "github.com/bianjieai/tibc-go/modules/tibc/light-clients/07-tendermint/types"
	ibctesting "github.com/bianjieai/tibc-go/modules/tibc/testing"
	ibctestingmock "github.com/bianjieai/tibc-go/modules/tibc/testing/mock"
)

func (suite *TendermintTestSuite) TestCheckHeaderAndUpdateState() {
	var (
		clientState     *types.ClientState
		consensusState  *types.ConsensusState
		consStateHeight clienttypes.Height
		newHeader       *types.Header
		currentTime     time.Time
		bothValSet      *tmtypes.ValidatorSet
		signers         []tmtypes.PrivValidator
		bothSigners     []tmtypes.PrivValidator
	)

	// Setup different validators and signers for testing different types of updates
	altPrivVal := ibctestingmock.NewPV()
	altPubKey, err := altPrivVal.GetPubKey()
	suite.Require().NoError(err)

	revisionHeight := int64(height.RevisionHeight)

	// create modified heights to use for test-cases
	heightPlus1 := clienttypes.NewHeight(height.RevisionNumber, height.RevisionHeight+1)
	heightMinus1 := clienttypes.NewHeight(height.RevisionNumber, height.RevisionHeight-1)
	heightMinus3 := clienttypes.NewHeight(height.RevisionNumber, height.RevisionHeight-3)
	heightPlus5 := clienttypes.NewHeight(height.RevisionNumber, height.RevisionHeight+5)

	altVal := tmtypes.NewValidator(altPubKey, revisionHeight)
	// Create alternative validator set with only altVal, invalid update (too much change in valSet)
	altValSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{altVal})
	altSigners := []tmtypes.PrivValidator{altPrivVal}

	testCases := []struct {
		name    string
		setup   func(*TendermintTestSuite)
		expPass bool
	}{{
		name: "successful update with next height and same validator set",
		setup: func(suite *TendermintTestSuite) {
			clientState = types.NewClientState(chainID, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), prefix, 0)
			consensusState = types.NewConsensusState(suite.clientTime, commitmenttypes.NewMerkleRoot(suite.header.Header.GetAppHash()), suite.valsHash)
			newHeader = suite.chainA.CreateTMClientHeader(chainID, int64(heightPlus1.RevisionHeight), height, suite.headerTime, suite.valSet, suite.valSet, signers)
			currentTime = suite.now
		},
		expPass: true,
	}, {
		name: "successful update with future height and different validator set",
		setup: func(suite *TendermintTestSuite) {
			clientState = types.NewClientState(chainID, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), prefix, 0)
			consensusState = types.NewConsensusState(suite.clientTime, commitmenttypes.NewMerkleRoot(suite.header.Header.GetAppHash()), suite.valsHash)
			newHeader = suite.chainA.CreateTMClientHeader(chainID, int64(heightPlus5.RevisionHeight), height, suite.headerTime, bothValSet, suite.valSet, bothSigners)
			currentTime = suite.now
		},
		expPass: true,
	}, {
		name: "successful update with next height and different validator set",
		setup: func(suite *TendermintTestSuite) {
			clientState = types.NewClientState(chainID, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), prefix, 0)
			consensusState = types.NewConsensusState(suite.clientTime, commitmenttypes.NewMerkleRoot(suite.header.Header.GetAppHash()), bothValSet.Hash())
			newHeader = suite.chainA.CreateTMClientHeader(chainID, int64(heightPlus1.RevisionHeight), height, suite.headerTime, bothValSet, bothValSet, bothSigners)
			currentTime = suite.now
		},
		expPass: true,
	}, {
		name: "successful update for a previous height",
		setup: func(suite *TendermintTestSuite) {
			clientState = types.NewClientState(chainID, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), prefix, 0)
			consensusState = types.NewConsensusState(suite.clientTime, commitmenttypes.NewMerkleRoot(suite.header.Header.GetAppHash()), suite.valsHash)
			consStateHeight = heightMinus3
			newHeader = suite.chainA.CreateTMClientHeader(chainID, int64(heightMinus1.RevisionHeight), heightMinus3, suite.headerTime, bothValSet, suite.valSet, bothSigners)
			currentTime = suite.now
		},
		expPass: true,
	}, {
		name: "successful update for a previous revision",
		setup: func(suite *TendermintTestSuite) {
			clientState = types.NewClientState(chainIDRevision1, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), prefix, 0)
			consensusState = types.NewConsensusState(suite.clientTime, commitmenttypes.NewMerkleRoot(suite.header.Header.GetAppHash()), suite.valsHash)
			consStateHeight = heightMinus3
			newHeader = suite.chainA.CreateTMClientHeader(chainIDRevision0, int64(height.RevisionHeight), heightMinus3, suite.headerTime, bothValSet, suite.valSet, bothSigners)
			currentTime = suite.now
		},
		expPass: true,
	}, {
		name: "successful update with identical header to a previous update",
		setup: func(suite *TendermintTestSuite) {
			clientState = types.NewClientState(chainID, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, heightPlus1, commitmenttypes.GetSDKSpecs(), prefix, 0)
			consensusState = types.NewConsensusState(suite.clientTime, commitmenttypes.NewMerkleRoot(suite.header.Header.GetAppHash()), suite.valsHash)
			newHeader = suite.chainA.CreateTMClientHeader(chainID, int64(heightPlus1.RevisionHeight), height, suite.headerTime, suite.valSet, suite.valSet, signers)
			currentTime = suite.now
			ctx := suite.chainA.GetContext().WithBlockTime(currentTime)
			// Store the header's consensus state in client store before UpdateClient call
			suite.chainA.App.TIBCKeeper.ClientKeeper.SetClientConsensusState(ctx, chainName, heightPlus1, newHeader.ConsensusState())
		},
		expPass: true,
	}, {
		name: "unsuccessful update with incorrect header chain-id",
		setup: func(suite *TendermintTestSuite) {
			clientState = types.NewClientState(chainID, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), prefix, 0)
			consensusState = types.NewConsensusState(suite.clientTime, commitmenttypes.NewMerkleRoot(suite.header.Header.GetAppHash()), suite.valsHash)
			newHeader = suite.chainA.CreateTMClientHeader("ethermint", int64(heightPlus1.RevisionHeight), height, suite.headerTime, suite.valSet, suite.valSet, signers)
			currentTime = suite.now
		},
		expPass: false,
	}, {
		name: "unsuccessful update to a future revision",
		setup: func(suite *TendermintTestSuite) {
			clientState = types.NewClientState(chainIDRevision0, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), prefix, 0)
			consensusState = types.NewConsensusState(suite.clientTime, commitmenttypes.NewMerkleRoot(suite.header.Header.GetAppHash()), suite.valsHash)
			newHeader = suite.chainA.CreateTMClientHeader(chainIDRevision1, 1, height, suite.headerTime, suite.valSet, suite.valSet, signers)
			currentTime = suite.now
		},
		expPass: false,
	}, {
		name: "unsuccessful update: header height revision and trusted height revision mismatch",
		setup: func(suite *TendermintTestSuite) {
			clientState = types.NewClientState(chainIDRevision1, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, clienttypes.NewHeight(1, 1), commitmenttypes.GetSDKSpecs(), prefix, 0)
			consensusState = types.NewConsensusState(suite.clientTime, commitmenttypes.NewMerkleRoot(suite.header.Header.GetAppHash()), suite.valsHash)
			newHeader = suite.chainA.CreateTMClientHeader(chainIDRevision1, 3, height, suite.headerTime, suite.valSet, suite.valSet, signers)
			currentTime = suite.now
		},
		expPass: false,
	}, {
		name: "unsuccessful update with next height: update header mismatches nextValSetHash",
		setup: func(suite *TendermintTestSuite) {
			clientState = types.NewClientState(chainID, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), prefix, 0)
			consensusState = types.NewConsensusState(suite.clientTime, commitmenttypes.NewMerkleRoot(suite.header.Header.GetAppHash()), suite.valsHash)
			newHeader = suite.chainA.CreateTMClientHeader(chainID, int64(heightPlus1.RevisionHeight), height, suite.headerTime, bothValSet, suite.valSet, bothSigners)
			currentTime = suite.now
		},
		expPass: false,
	}, {
		name: "unsuccessful update with next height: update header mismatches different nextValSetHash",
		setup: func(suite *TendermintTestSuite) {
			clientState = types.NewClientState(chainID, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), prefix, 0)
			consensusState = types.NewConsensusState(suite.clientTime, commitmenttypes.NewMerkleRoot(suite.header.Header.GetAppHash()), bothValSet.Hash())
			newHeader = suite.chainA.CreateTMClientHeader(chainID, int64(heightPlus1.RevisionHeight), height, suite.headerTime, suite.valSet, bothValSet, signers)
			currentTime = suite.now
		},
		expPass: false,
	}, {
		name: "unsuccessful update with future height: too much change in validator set",
		setup: func(suite *TendermintTestSuite) {
			clientState = types.NewClientState(chainID, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), prefix, 0)
			consensusState = types.NewConsensusState(suite.clientTime, commitmenttypes.NewMerkleRoot(suite.header.Header.GetAppHash()), suite.valsHash)
			newHeader = suite.chainA.CreateTMClientHeader(chainID, int64(heightPlus5.RevisionHeight), height, suite.headerTime, altValSet, suite.valSet, altSigners)
			currentTime = suite.now
		},
		expPass: false,
	}, {
		name: "unsuccessful updates, passed in incorrect trusted validators for given consensus state",
		setup: func(suite *TendermintTestSuite) {
			clientState = types.NewClientState(chainID, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), prefix, 0)
			consensusState = types.NewConsensusState(suite.clientTime, commitmenttypes.NewMerkleRoot(suite.header.Header.GetAppHash()), suite.valsHash)
			newHeader = suite.chainA.CreateTMClientHeader(chainID, int64(heightPlus5.RevisionHeight), height, suite.headerTime, bothValSet, bothValSet, bothSigners)
			currentTime = suite.now
		},
		expPass: false,
	}, {
		name: "unsuccessful update: trusting period has passed since last client timestamp",
		setup: func(suite *TendermintTestSuite) {
			clientState = types.NewClientState(chainID, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), prefix, 0)
			consensusState = types.NewConsensusState(suite.clientTime, commitmenttypes.NewMerkleRoot(suite.header.Header.GetAppHash()), suite.valsHash)
			newHeader = suite.chainA.CreateTMClientHeader(chainID, int64(heightPlus1.RevisionHeight), height, suite.headerTime, suite.valSet, suite.valSet, signers)
			// make current time pass trusting period from last timestamp on clientstate
			currentTime = suite.now.Add(trustingPeriod)
		},
		expPass: false,
	}, {
		name: "unsuccessful update: header timestamp is past current timestamp",
		setup: func(suite *TendermintTestSuite) {
			clientState = types.NewClientState(chainID, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), prefix, 0)
			consensusState = types.NewConsensusState(suite.clientTime, commitmenttypes.NewMerkleRoot(suite.header.Header.GetAppHash()), suite.valsHash)
			newHeader = suite.chainA.CreateTMClientHeader(chainID, int64(heightPlus1.RevisionHeight), height, suite.now.Add(time.Minute), suite.valSet, suite.valSet, signers)
			currentTime = suite.now
		},
		expPass: false,
	}, {
		name: "unsuccessful update: header timestamp is not past last client timestamp",
		setup: func(suite *TendermintTestSuite) {
			clientState = types.NewClientState(chainID, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), prefix, 0)
			consensusState = types.NewConsensusState(suite.clientTime, commitmenttypes.NewMerkleRoot(suite.header.Header.GetAppHash()), suite.valsHash)
			newHeader = suite.chainA.CreateTMClientHeader(chainID, int64(heightPlus1.RevisionHeight), height, suite.clientTime, suite.valSet, suite.valSet, signers)
			currentTime = suite.now
		},
		expPass: false,
	}, {
		name: "header basic validation failed",
		setup: func(suite *TendermintTestSuite) {
			clientState = types.NewClientState(chainID, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), prefix, 0)
			consensusState = types.NewConsensusState(suite.clientTime, commitmenttypes.NewMerkleRoot(suite.header.Header.GetAppHash()), suite.valsHash)
			newHeader = suite.chainA.CreateTMClientHeader(chainID, int64(heightPlus1.RevisionHeight), height, suite.headerTime, suite.valSet, suite.valSet, signers)
			// cause new header to fail validatebasic by changing commit height to mismatch header height
			newHeader.SignedHeader.Commit.Height = revisionHeight - 1
			currentTime = suite.now
		},
		expPass: false,
	}, {
		name: "header height < consensus height",
		setup: func(suite *TendermintTestSuite) {
			clientState = types.NewClientState(chainID, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, clienttypes.NewHeight(height.RevisionNumber, heightPlus5.RevisionHeight), commitmenttypes.GetSDKSpecs(), prefix, 0)
			consensusState = types.NewConsensusState(suite.clientTime, commitmenttypes.NewMerkleRoot(suite.header.Header.GetAppHash()), suite.valsHash)
			// Make new header at height less than latest client state
			newHeader = suite.chainA.CreateTMClientHeader(chainID, int64(heightMinus1.RevisionHeight), height, suite.headerTime, suite.valSet, suite.valSet, signers)
			currentTime = suite.now
		},
		expPass: false,
	}}

	for i, tc := range testCases {
		tc := tc
		suite.Run(fmt.Sprintf("Case: %s", tc.name), func() {
			suite.SetupTest() // reset metadata writes
			// Create bothValSet with both suite validator and altVal. Would be valid update
			bothValSet = tmtypes.NewValidatorSet(append(suite.valSet.Validators, altVal))
			signers = []tmtypes.PrivValidator{suite.privVal}

			// Create signer array and ensure it is in same order as bothValSet
			_, suiteVal := suite.valSet.GetByIndex(0)
			bothSigners = ibctesting.CreateSortedSignerArray(altPrivVal, suite.privVal, altVal, suiteVal)

			consStateHeight = height // must be explicitly changed
			// setup test
			tc.setup(suite)

			// Set current timestamp in context
			ctx := suite.chainA.GetContext().WithBlockTime(currentTime)

			// Set trusted consensus state in client store
			suite.chainA.App.TIBCKeeper.ClientKeeper.SetClientConsensusState(ctx, chainName, consStateHeight, consensusState)

			height := newHeader.GetHeight()
			expectedConsensus := &types.ConsensusState{
				Timestamp:          newHeader.GetTime(),
				Root:               commitmenttypes.NewMerkleRoot(newHeader.Header.GetAppHash()),
				NextValidatorsHash: newHeader.Header.NextValidatorsHash,
			}

			newClientState, consensusState, err := clientState.CheckHeaderAndUpdateState(
				ctx,
				suite.cdc,
				suite.chainA.App.TIBCKeeper.ClientKeeper.ClientStore(suite.chainA.GetContext(), chainName), // pass in chainName prefixed clientStore
				newHeader,
			)

			if tc.expPass {
				suite.Require().NoError(err, "valid test case %d failed: %s", i, tc.name)

				// Determine if clientState should be updated or not
				// TODO: check the entire Height struct once GetLatestHeight returns clienttypes.Height
				if height.GT(clientState.LatestHeight) {
					// Header Height is greater than clientState latest Height, clientState should be updated with header.GetHeight()
					suite.Require().Equal(height, newClientState.GetLatestHeight(), "clientstate height did not update")
				} else {
					// Update will add past consensus state, clientState should not be updated at all
					suite.Require().Equal(clientState.LatestHeight, newClientState.GetLatestHeight(), "client state height updated for past header")
				}

				suite.Require().Equal(expectedConsensus, consensusState, "valid test case %d failed: %s", i, tc.name)
			} else {
				suite.Require().Error(err, "invalid test case %d passed: %s", i, tc.name)
				suite.Require().Nil(newClientState, "invalid test case %d passed: %s", i, tc.name)
				suite.Require().Nil(consensusState, "invalid test case %d passed: %s", i, tc.name)
			}
		})
	}
}

func (suite *TendermintTestSuite) TestPruneConsensusState() {
	// create path and setup clients
	path := ibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupClients(path)

	// get the first height as it will be pruned first.
	var pruneHeight exported.Height
	getFirstHeightCb := func(height exported.Height) bool {
		pruneHeight = height
		return true
	}
	clientStore := path.EndpointA.ClientStore()
	types.IterateConsensusStateAscending(clientStore, getFirstHeightCb)

	// this height will be expired but not pruned
	err := path.EndpointA.UpdateClient()
	suite.Require().NoError(err)
	expiredHeight := path.EndpointA.GetClientState().GetLatestHeight()

	// expected values that must still remain in store after pruning
	expectedConsState := path.EndpointA.GetConsensusState(expiredHeight)

	clientStore = path.EndpointA.ClientStore()
	expectedProcessTime, ok := types.GetProcessedTime(clientStore, expiredHeight)
	suite.Require().True(ok)

	expectedConsKey := types.GetIterationKey(clientStore, expiredHeight)
	suite.Require().NotNil(expectedConsKey)

	// Increment the time by a week
	suite.coordinator.IncrementTimeBy(7 * 24 * time.Hour)

	// create the consensus state that can be used as trusted height for next update
	err = path.EndpointA.UpdateClient()
	suite.Require().NoError(err)

	// Increment the time by another week, then update the client.
	// This will cause the first two consensus states to become expired.
	suite.coordinator.IncrementTimeBy(7 * 24 * time.Hour)
	err = path.EndpointA.UpdateClient()
	suite.Require().NoError(err)

	clientStore = path.EndpointA.ClientStore()

	// check that the first expired consensus state got deleted along with all associated metadata
	consState, ok := path.EndpointA.Chain.GetConsensusState(path.EndpointB.ChainName, pruneHeight)
	suite.Require().Nil(consState, "expired consensus state not pruned")
	suite.Require().False(ok)
	// check processed time metadata is pruned
	processTime, ok := types.GetProcessedTime(clientStore, pruneHeight)
	suite.Require().Equal(uint64(0), processTime, "processed time metadata not pruned")
	suite.Require().False(ok)

	// check iteration key metadata is pruned
	consKey := types.GetIterationKey(clientStore, pruneHeight)
	suite.Require().Nil(consKey, "iteration key not pruned")

	// check that second expired consensus state doesn't get deleted
	// this ensures that there is a cap on gas cost of UpdateClient
	consState = path.EndpointA.GetConsensusState(expiredHeight)
	suite.Require().Equal(expectedConsState, consState, "consensus state incorrectly pruned")
	// check processed time metadata is not pruned
	processTime, ok = types.GetProcessedTime(clientStore, expiredHeight)
	suite.Require().Equal(expectedProcessTime, processTime, "processed time metadata incorrectly pruned")
	suite.Require().True(ok)

	// check iteration key metadata is not pruned
	consKey = types.GetIterationKey(clientStore, expiredHeight)
	suite.Require().Equal(expectedConsKey, consKey, "iteration key incorrectly pruned")
}
