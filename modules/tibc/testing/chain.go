package tibctesting

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/math"

	"github.com/stretchr/testify/require"

	errorsmod "cosmossdk.io/errors"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto/tmhash"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmtprotoversion "github.com/cometbft/cometbft/proto/tendermint/version"
	cmttypes "github.com/cometbft/cometbft/types"
	cmtversion "github.com/cometbft/cometbft/version"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/staking/testutil"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	commitmenttypes "github.com/bianjieai/tibc-go/modules/tibc/core/23-commitment/types"
	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
	"github.com/bianjieai/tibc-go/modules/tibc/core/types"
	ibccmttypes "github.com/bianjieai/tibc-go/modules/tibc/light-clients/07-tendermint/types"
	"github.com/bianjieai/tibc-go/modules/tibc/testing/mock"
	"github.com/bianjieai/tibc-go/simapp"
)

const (
	// Default params constants used to create a TM client
	TrustingPeriod     time.Duration = time.Hour * 24 * 7 * 2
	UnbondingPeriod    time.Duration = time.Hour * 24 * 7 * 3
	MaxClockDrift      time.Duration = time.Second * 10
	DefaultDelayPeriod uint64        = 0

	InvalidID = "IDisInvalid"

	ConnectionIDPrefix = "conn"
	ChannelIDPrefix    = "chan"

	MockPort = mock.ModuleName

	// used for testing UpdateClientProposal
	Title       = "title"
	Description = "description"
)

var (

	// Default params variables used to create a TM client
	DefaultTrustLevel ibccmttypes.Fraction = ibccmttypes.DefaultTrustLevel
	TestHash                              = tmhash.Sum([]byte("TESTING HASH"))
	TestCoin                              = sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100))

	MockAcknowledgement = mock.MockAcknowledgement
	MockCommitment      = mock.MockCommitment
	Prefix              = commitmenttypes.MerklePrefix{KeyPrefix: []byte("tibc")}
	MaxAccounts         = 10
)

type SenderAccount struct {
	SenderPrivKey cryptotypes.PrivKey
	SenderAccount authtypes.AccountI
}

// TestChain is a testing struct that wraps a simapp with the last TM Header, the current ABCI
// header and the validators of the TestChain. It also contains a field called ChainID. This
// is the chainName that *other* chains use to refer to this TestChain. The SenderAccount
// is used for delivering transactions through the application state.
// NOTE: the actual application uses an empty chain-id for ease of testing.
type TestChain struct {
	*testing.T

	Coordinator        *Coordinator
	App                *simapp.SimApp
	ChainID, ChainName string
	LastHeader         *ibccmttypes.Header // header for last block height committed
	CurrentHeader      cmtproto.Header     // header for current block height
	QueryServer        types.QueryServer
	TxConfig           client.TxConfig
	Codec              codec.BinaryCodec

	Vals     *cmttypes.ValidatorSet
	NextVals *cmttypes.ValidatorSet

	// Signers is a map from validator address to the PrivValidator
	// The map is converted into an array that is the same order as the validators right before signing commit
	// This ensures that signers will always be in correct order even as validator powers change.
	// If a test adds a new validator after chain creation, then the signer map must be updated to include
	// the new PrivValidator entry.
	Signers map[string]cmttypes.PrivValidator

	SenderPrivKey cryptotypes.PrivKey
	SenderAccount authtypes.AccountI

	SenderAccounts []SenderAccount
}

// NewTestChainWithValSet initializes a new TestChain instance with the given validator set
// and signer array. It also initializes 10 Sender accounts with a balance of 10000000000000000000 coins of
// bond denom to use for tests.
//
// The first block height is committed to state in order to allow for client creations on
// counterparty chains. The TestChain will return with a block height starting at 2.
//
// Time management is handled by the Coordinator in order to ensure synchrony between chains.
// Each update of any chain increments the block header time for all chains by 5 seconds.
//
// NOTE: to use a custom sender privkey and account for testing purposes, replace and modify this
// constructor function.
//
// CONTRACT: Validator array must be provided in the order expected by Tendermint.
// i.e. sorted first by power and then lexicographically by address.
func NewTestChainWithValSet(
	t *testing.T,
	coord *Coordinator,
	chainID string,
	valSet *cmttypes.ValidatorSet,
	signers map[string]cmttypes.PrivValidator,
) *TestChain {
	genAccs := []authtypes.GenesisAccount{}
	genBals := []banktypes.Balance{}
	senderAccs := []SenderAccount{}

	// generate genesis accounts
	for i := 0; i < MaxAccounts; i++ {
		senderPrivKey := secp256k1.GenPrivKey()
		acc := authtypes.NewBaseAccount(
			senderPrivKey.PubKey().Address().Bytes(),
			senderPrivKey.PubKey(),
			uint64(i),
			0,
		)
		amount, ok := math.NewIntFromString("10000000000000000000")
		require.True(t, ok)

		balance := banktypes.Balance{
			Address: acc.GetAddress().String(),
			Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, amount)),
		}

		genAccs = append(genAccs, acc)
		genBals = append(genBals, balance)

		senderAcc := SenderAccount{
			SenderAccount: acc,
			SenderPrivKey: senderPrivKey,
		}

		senderAccs = append(senderAccs, senderAcc)
	}

	app := SetupWithGenesisValSet(
		t,
		valSet,
		genAccs,
		chainID,
		sdk.DefaultPowerReduction,
		genBals...)

	// create current header and call begin block
	header := cmtproto.Header{
		ChainID: chainID,
		Height:  1,
		Time:    coord.CurrentTime.UTC(),
	}

	txConfig := simapp.MakeTestEncodingConfig().TxConfig

	// create an account to send transactions from
	chain := &TestChain{
		T:              t,
		Coordinator:    coord,
		ChainID:        chainID,
		ChainName:      chainID,
		App:            app,
		CurrentHeader:  header,
		QueryServer:    app.TIBCKeeper,
		TxConfig:       txConfig,
		Codec:          app.AppCodec(),
		Vals:           valSet,
		NextVals:       valSet,
		Signers:        signers,
		SenderPrivKey:  senderAccs[0].SenderPrivKey,
		SenderAccount:  senderAccs[0].SenderAccount,
		SenderAccounts: senderAccs,
	}

	coord.CommitBlock(chain)

	return chain
}

// NewTestChain initializes a new test chain with a default of 4 validators
// Use this function if the tests do not need custom control over the validator set
func NewTestChain(t *testing.T, coord *Coordinator, chainID string) *TestChain {
	// generate validators private/public key
	var (
		validatorsPerChain = 4
		validators         []*cmttypes.Validator
		signersByAddress   = make(map[string]cmttypes.PrivValidator, validatorsPerChain)
	)

	for i := 0; i < validatorsPerChain; i++ {
		privVal := mock.NewPV()
		pubKey, err := privVal.GetPubKey()
		require.NoError(t, err)
		validators = append(validators, cmttypes.NewValidator(pubKey, 1))
		signersByAddress[pubKey.Address().String()] = privVal
	}

	// construct validator set;
	// Note that the validators are sorted by voting power
	// or, if equal, by address lexical order
	valSet := cmttypes.NewValidatorSet(validators)

	return NewTestChainWithValSet(t, coord, chainID, valSet, signersByAddress)
}

// GetContext returns the current context for the application.
func (chain *TestChain) GetContext() sdk.Context {
	return chain.App.BaseApp.NewContextLegacy(false, chain.CurrentHeader)
}

// QueryProof performs an abci query with the given key and returns the proto encoded merkle proof
// for the query and the height at which the proof will succeed on a tendermint verifier.
func (chain *TestChain) QueryProof(key []byte) ([]byte, clienttypes.Height) {
	return chain.QueryProofAtHeight(key, chain.App.LastBlockHeight())
}

// QueryProofAtHeight performs an abci query with the given key and returns the proto encoded merkle proof
// for the query and the height at which the proof will succeed on a tendermint verifier.
func (chain *TestChain) QueryProofAtHeight(key []byte, height int64) ([]byte, clienttypes.Height) {
	res, err := chain.App.Query(context.Background(), &abci.RequestQuery{
		Path:   fmt.Sprintf("store/%s/key", host.StoreKey),
		Height: height - 1,
		Data:   key,
		Prove:  true,
	})
	require.NoError(chain.T, err)

	merkleProof, err := commitmenttypes.ConvertProofs(res.ProofOps)
	require.NoError(chain.T, err)

	proof, err := chain.App.AppCodec().Marshal(&merkleProof)
	require.NoError(chain.T, err)

	revision := clienttypes.ParseChainID(chain.ChainID)

	// proof height + 1 is returned as the proof created corresponds to the height the proof
	// was created in the IAVL tree. Tendermint and subsequently the clients that rely on it
	// have heights 1 above the IAVL tree. Thus we return proof height + 1
	return proof, clienttypes.NewHeight(revision, uint64(res.Height)+1)
}

// QueryUpgradeProof performs an abci query with the given key and returns the proto encoded merkle proof
// for the query and the height at which the proof will succeed on a tendermint verifier.
func (chain *TestChain) QueryUpgradeProof(key []byte, height uint64) ([]byte, clienttypes.Height) {
	res,err := chain.App.Query(context.Background(),&abci.RequestQuery{
		Path:   "store/upgrade/key",
		Height: int64(height - 1),
		Data:   key,
		Prove:  true,
	})
	require.NoError(chain.T, err)

	merkleProof, err := commitmenttypes.ConvertProofs(res.ProofOps)
	require.NoError(chain.T, err)

	proof, err := chain.App.AppCodec().Marshal(&merkleProof)
	require.NoError(chain.T, err)

	revision := clienttypes.ParseChainID(chain.ChainID)

	// proof height + 1 is returned as the proof created corresponds to the height the proof
	// was created in the IAVL tree. Tendermint and subsequently the clients that rely on it
	// have heights 1 above the IAVL tree. Thus we return proof height + 1
	return proof, clienttypes.NewHeight(revision, uint64(res.Height+1))
}

// QueryClientStateProof performs and abci query for a client state
// stored with a given chainName and returns the ClientState along with the proof
func (chain *TestChain) QueryClientStateProof(chainName string) (exported.ClientState, []byte) {
	// retrieve client state to provide proof for
	clientState, found := chain.App.TIBCKeeper.ClientKeeper.GetClientState(
		chain.GetContext(),
		chainName,
	)
	require.True(chain.T, found)

	clientKey := host.FullClientStateKey(chainName)
	proofClient, _ := chain.QueryProof(clientKey)

	return clientState, proofClient
}

// QueryConsensusStateProof performs an abci query for a consensus state
// stored on the given chainName. The proof and consensusHeight are returned.
func (chain *TestChain) QueryConsensusStateProof(chainName string) ([]byte, clienttypes.Height) {
	clientState := chain.GetClientState(chainName)

	consensusHeight := clientState.GetLatestHeight().(clienttypes.Height)
	consensusKey := host.FullConsensusStateKey(chainName, consensusHeight)
	proofConsensus, _ := chain.QueryProof(consensusKey)

	return proofConsensus, consensusHeight
}

// NextBlock sets the last header to the current header and increments the current header to be
// at the next block height. It does not update the time as that is handled by the Coordinator.
//
// CONTRACT: this function must only be called after app.Commit() occurs
func (chain *TestChain) NextBlock() {
	res, err := chain.App.FinalizeBlock(&abci.RequestFinalizeBlock{
		Height:             chain.CurrentHeader.Height,
		Time:               chain.CurrentHeader.GetTime(),
		NextValidatorsHash: chain.NextVals.Hash(),
	})
	require.NoError(chain.T, err)
	chain.commitBlock(res)
}

// sendMsgs delivers a transaction through the application without returning the result.
func (chain *TestChain) sendMsgs(msgs ...sdk.Msg) error {
	_, err := chain.SendMsgs(msgs...)
	return err
}

// SendMsgs delivers a transaction through the application. It updates the senders sequence
// number and updates the TestChain's headers. It returns the result and error if one
// occurred.
func (chain *TestChain) SendMsgs(msgs ...sdk.Msg) (*sdk.Result, error) {
	// ensure the chain has the latest time
	chain.Coordinator.UpdateTimeForChain(chain)

	_, r, err := simapp.SignAndDeliver(
		chain.T,
		chain.TxConfig,
		chain.App.BaseApp,
		chain.GetContext().BlockHeader(),
		msgs,
		chain.ChainID,
		[]uint64{chain.SenderAccount.GetAccountNumber()},
		[]uint64{chain.SenderAccount.GetSequence()},
		true, true, chain.SenderPrivKey,
	)
	if err != nil {
		return nil, err
	}

	// NextBlock calls app.Commit()
	chain.NextBlock()

	// increment sequence for successful transaction execution
	chain.SenderAccount.SetSequence(chain.SenderAccount.GetSequence() + 1)

	chain.Coordinator.IncrementTime()

	return r, nil
}

// GetClientState retrieves the client state for the provided chainName. The client is
// expected to exist otherwise testing will fail.
func (chain *TestChain) GetClientState(chainName string) exported.ClientState {
	clientState, found := chain.App.TIBCKeeper.ClientKeeper.GetClientState(
		chain.GetContext(),
		chainName,
	)
	require.True(chain.T, found)

	return clientState
}

// GetConsensusState retrieves the consensus state for the provided chainName and height.
// It will return a success boolean depending on if consensus state exists or not.
func (chain *TestChain) GetConsensusState(
	chainName string,
	height exported.Height,
) (exported.ConsensusState, bool) {
	return chain.App.TIBCKeeper.ClientKeeper.GetClientConsensusState(
		chain.GetContext(),
		chainName,
		height,
	)
}

// GetValsAtHeight will return the validator set of the chain at a given height. It will return
// a success boolean depending on if the validator set exists or not at that height.
func (chain *TestChain) GetValsAtHeight(height int64) (*cmttypes.ValidatorSet, bool) {
	histInfo, err := chain.App.StakingKeeper.GetHistoricalInfo(chain.GetContext(), height)
	if err != nil {
		return nil, false
	}

	valSet := stakingtypes.Validators{
		Validators: histInfo.Valset,
	}

	tmValidators, err := testutil.ToCmtValidators(valSet, sdk.DefaultPowerReduction)
	if err != nil {
		panic(err)
	}
	return cmttypes.NewValidatorSet(tmValidators), true
}

// GetAcknowledgement retrieves an acknowledgement for the provided packet. If the
// acknowledgement does not exist then testing will fail.
func (chain *TestChain) GetAcknowledgement(packet exported.PacketI) []byte {
	ack, found := chain.App.TIBCKeeper.PacketKeeper.GetPacketAcknowledgement(
		chain.GetContext(),
		packet.GetSourceChain(),
		packet.GetDestChain(),
		packet.GetSequence(),
	)
	require.True(chain.T, found)

	return ack
}

// GetPrefix returns the prefix for used by a chain in connection creation
func (chain *TestChain) GetPrefix() commitmenttypes.MerklePrefix {
	return commitmenttypes.NewMerklePrefix([]byte(""))
}

// ConstructMsgCreateClient constructs a message to create a new client state (tendermint or solomachine).
// NOTE: a solo machine client will be created with an empty diversifier.
func (chain *TestChain) ConstructMsgCreateClient(
	counterparty *TestChain,
	chainName string,
	clientType string,
) error {
	var (
		clientState    exported.ClientState
		consensusState exported.ConsensusState
	)

	switch clientType {
	case exported.Tendermint:
		height := counterparty.LastHeader.GetHeight().(clienttypes.Height)
		clientState = ibccmttypes.NewClientState(
			counterparty.ChainID, DefaultTrustLevel,
			TrustingPeriod, UnbondingPeriod, MaxClockDrift,
			height, commitmenttypes.GetSDKSpecs(), Prefix, 0,
		)
		consensusState = counterparty.LastHeader.ConsensusState()
	default:
		chain.T.Fatalf("unsupported client state type %s", clientType)
	}

	err := chain.App.TIBCKeeper.ClientKeeper.CreateClient(
		chain.GetContext(),
		chainName,
		clientState, consensusState,
	)

	require.NoError(chain.T, err)
	return err
}

// CreateTMClient will construct and execute a 07-tendermint MsgCreateClient. A counterparty
// client will be created on the (target) chain.
func (chain *TestChain) CreateTMClient(counterparty *TestChain, chainName string) error {
	// construct MsgCreateClient using counterparty
	return chain.ConstructMsgCreateClient(counterparty, chainName, exported.Tendermint)
}

// UpdateTMClient will construct and execute a 07-tendermint MsgUpdateClient. The counterparty
// client will be updated on the (target) chain. UpdateTMClient mocks the relayer flow
// necessary for updating a Tendermint client.
func (chain *TestChain) UpdateTMClient(counterparty *TestChain, chainName string) error {
	header, err := chain.ConstructUpdateTMClientHeader(counterparty, chainName)
	require.NoError(chain.T, err)

	msg, err := clienttypes.NewMsgUpdateClient(
		chainName, header,
		chain.SenderAccount.GetAddress(),
	)
	require.NoError(chain.T, err)

	return chain.sendMsgs(msg)
}

// ConstructUpdateTMClientHeader will construct a valid 07-tendermint Header to update the
// light client on the source chain.
func (chain *TestChain) ConstructUpdateTMClientHeader(
	counterparty *TestChain,
	chainName string,
) (*ibccmttypes.Header, error) {
	return chain.ConstructUpdateTMClientHeaderWithTrustedHeight(
		counterparty,
		chainName,
		clienttypes.ZeroHeight(),
	)
}

// ConstructUpdateTMClientHeaderWithTrustedHeight will construct a valid 07-tendermint Header to update the
// light client on the source chain.
func (chain *TestChain) ConstructUpdateTMClientHeaderWithTrustedHeight(
	counterparty *TestChain,
	chainName string,
	trustedHeight clienttypes.Height,
) (*ibccmttypes.Header, error) {
	header := counterparty.LastHeader
	// Relayer must query for LatestHeight on client to get TrustedHeight if the trusted height is not set
	if trustedHeight.IsZero() {
		trustedHeight = chain.GetClientState(chainName).GetLatestHeight().(clienttypes.Height)
	}
	var (
		tmTrustedVals *cmttypes.ValidatorSet
		ok            bool
	)
	// Once we get TrustedHeight from client, we must query the validators from the counterparty chain
	// If the LatestHeight == LastHeader.Height, then TrustedValidators are current validators
	// If LatestHeight < LastHeader.Height, we can query the historical validator set from HistoricalInfo
	if trustedHeight == counterparty.LastHeader.GetHeight() {
		tmTrustedVals = counterparty.Vals
	} else {
		// NOTE: We need to get validators from counterparty at height: trustedHeight+1
		// since the last trusted validators for a header at height h
		// is the NextValidators at h+1 committed to in header h by
		// NextValidatorsHash
		tmTrustedVals, ok = counterparty.GetValsAtHeight(int64(trustedHeight.RevisionHeight + 1))
		if !ok {
			return nil, errorsmod.Wrapf(ibccmttypes.ErrInvalidHeaderHeight, "could not retrieve trusted validators at trustedHeight: %d", trustedHeight)
		}
	}
	// inject trusted fields into last header
	// for now assume revision number is 0
	header.TrustedHeight = trustedHeight

	trustedVals, err := tmTrustedVals.ToProto()
	if err != nil {
		return nil, err
	}
	header.TrustedValidators = trustedVals
	return header, nil
}

// ExpireClient fast forwards the chain's block time by the provided amount of time which will
// expire any clients with a trusting period less than or equal to this amount of time.
func (chain *TestChain) ExpireClient(amount time.Duration) {
	chain.Coordinator.IncrementTimeBy(amount)
}

// CurrentTMClientHeader creates a TM header using the current header parameters
// on the chain. The trusted fields in the header are set to nil.
func (chain *TestChain) CurrentTMClientHeader() *ibccmttypes.Header {
	return chain.CreateTMClientHeader(
		chain.ChainID,
		chain.CurrentHeader.Height,
		clienttypes.Height{},
		chain.CurrentHeader.Time,
		chain.Vals,
		chain.NextVals,
		nil,
		chain.Signers,
	)
}

// CreateTMClientHeader creates a TM header to update the TM client. Args are passed in to allow
// caller flexibility to use params that differ from the chain.
func (chain *TestChain) CreateTMClientHeader(
	chainID string,
	blockHeight int64,
	trustedHeight clienttypes.Height,
	timestamp time.Time,
	tmValSet,
	nextVals,
	tmTrustedVals *cmttypes.ValidatorSet,
	signers map[string]cmttypes.PrivValidator,
) *ibccmttypes.Header {
	var (
		valSet      *cmtproto.ValidatorSet
		trustedVals *cmtproto.ValidatorSet
	)
	require.NotNil(chain.T, tmValSet)

	vsetHash := tmValSet.Hash()
	nextValHash := nextVals.Hash()

	tmHeader := cmttypes.Header{
		Version: cmtprotoversion.Consensus{Block: cmtversion.BlockProtocol, App: 2},
		ChainID: chainID,
		Height:  blockHeight,
		Time:    timestamp,
		LastBlockID: MakeBlockID(
			make([]byte, tmhash.Size),
			10_000,
			make([]byte, tmhash.Size),
		),
		LastCommitHash:     chain.App.LastCommitID().Hash,
		DataHash:           tmhash.Sum([]byte("data_hash")),
		ValidatorsHash:     vsetHash,
		NextValidatorsHash: nextValHash,
		ConsensusHash:      tmhash.Sum([]byte("consensus_hash")),
		AppHash:            chain.CurrentHeader.AppHash,
		LastResultsHash:    tmhash.Sum([]byte("last_results_hash")),
		EvidenceHash:       tmhash.Sum([]byte("evidence_hash")),
		ProposerAddress:    tmValSet.Proposer.Address, //nolint:staticcheck
	}

	hhash := tmHeader.Hash()
	blockID := MakeBlockID(hhash, 3, tmhash.Sum([]byte("part_set")))
	voteSet := cmttypes.NewVoteSet(chainID, blockHeight, 1, cmtproto.PrecommitType, tmValSet)

	// MakeCommit expects a signer array in the same order as the validator array.
	// Thus we iterate over the ordered validator set and construct a signer array
	// from the signer map in the same order.
	var signerArr []cmttypes.PrivValidator
	for _, v := range tmValSet.Validators {
		signerArr = append(signerArr, signers[v.Address.String()])
	}

	commit, err := cmttypes.MakeExtCommit(
		blockID,
		blockHeight,
		1,
		voteSet,
		signerArr,
		timestamp,
		false,
	)
	require.NoError(chain.T, err)

	signedHeader := &cmtproto.SignedHeader{
		Header: tmHeader.ToProto(),
		Commit: commit.ToCommit().ToProto(),
	}

	if tmValSet != nil {
		valSet, err = tmValSet.ToProto()
		require.NoError(chain.T, err)
	}

	if tmTrustedVals != nil {
		trustedVals, err = tmTrustedVals.ToProto()
		require.NoError(chain.T, err)
	}

	// The trusted fields may be nil. They may be filled before relaying messages to a client.
	// The relayer is responsible for querying client and injecting appropriate trusted fields.
	return &ibccmttypes.Header{
		SignedHeader:      signedHeader,
		ValidatorSet:      valSet,
		TrustedHeight:     trustedHeight,
		TrustedValidators: trustedVals,
	}
}

func (chain *TestChain) commitBlock(res *abci.ResponseFinalizeBlock) {
	_, err := chain.App.Commit()
	require.NoError(chain.T, err)

	// set the last header to the current header
	// use nil trusted fields
	chain.LastHeader = chain.CurrentTMClientHeader()

	// val set changes returned from previous block get applied to the next validators
	// of this block. See tendermint spec for details.
	chain.Vals = chain.NextVals
	chain.NextVals = ApplyValSetChanges(chain.T, chain.Vals, res.ValidatorUpdates)

	// increment the current header
	chain.CurrentHeader = cmtproto.Header{
		ChainID: chain.ChainID,
		Height:  chain.App.LastBlockHeight() + 1,
		AppHash: chain.App.LastCommitID().Hash,
		// NOTE: the time is increased by the coordinator to maintain time synchrony amongst
		// chains.
		Time:               chain.CurrentHeader.Time,
		ValidatorsHash:     chain.Vals.Hash(),
		NextValidatorsHash: chain.NextVals.Hash(),
		ProposerAddress:    chain.CurrentHeader.ProposerAddress,
	}
}

// MakeBlockID copied unimported test functions from cmttypes to use them here
func MakeBlockID(hash []byte, partSetSize uint32, partSetHash []byte) cmttypes.BlockID {
	return cmttypes.BlockID{
		Hash: hash,
		PartSetHeader: cmttypes.PartSetHeader{
			Total: partSetSize,
			Hash:  partSetHash,
		},
	}
}

// CreateSortedSignerArray takes two PrivValidators, and the corresponding Validator structs
// (including voting power). It returns a signer array of PrivValidators that matches the
// sorting of ValidatorSet.
// The sorting is first by .VotingPower (descending), with secondary index of .Address (ascending).
func CreateSortedSignerArray(
	altPrivVal, suitePrivVal cmttypes.PrivValidator, altVal, suiteVal *cmttypes.Validator,
) []cmttypes.PrivValidator {
	switch {
	case altVal.VotingPower > suiteVal.VotingPower:
		return []cmttypes.PrivValidator{altPrivVal, suitePrivVal}
	case altVal.VotingPower < suiteVal.VotingPower:
		return []cmttypes.PrivValidator{suitePrivVal, altPrivVal}
	default:
		if bytes.Compare(altVal.Address, suiteVal.Address) == -1 {
			return []cmttypes.PrivValidator{altPrivVal, suitePrivVal}
		}
		return []cmttypes.PrivValidator{suitePrivVal, altPrivVal}
	}
}
