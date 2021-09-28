package types

import (
	"strings"
	"time"

	ics23 "github.com/confio/ics23/go"

	"github.com/tendermint/tendermint/light"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	commitmenttypes "github.com/bianjieai/tibc-go/modules/tibc/core/23-commitment/types"
	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

var _ exported.ClientState = (*ClientState)(nil)

// NewClientState creates a new ClientState instance
func NewClientState(
	chainID string,
	trustLevel Fraction,
	trustingPeriod time.Duration,
	ubdPeriod time.Duration,
	maxClockDrift time.Duration,
	latestHeight clienttypes.Height,
	specs []*ics23.ProofSpec,
	prefix commitmenttypes.MerklePrefix,
	timeDelay uint64,
) *ClientState {
	return &ClientState{
		ChainId:         chainID,
		TrustLevel:      trustLevel,
		TrustingPeriod:  trustingPeriod,
		UnbondingPeriod: ubdPeriod,
		MaxClockDrift:   maxClockDrift,
		LatestHeight:    latestHeight,
		ProofSpecs:      specs,
		MerklePrefix:    prefix,
		TimeDelay:       timeDelay,
	}
}

// GetChainID returns the chain-id
func (cs ClientState) GetChainID() string {
	return cs.ChainId
}

// ClientType is tendermint.
func (cs ClientState) ClientType() string {
	return exported.Tendermint
}

// GetLatestHeight returns latest block height.
func (cs ClientState) GetLatestHeight() exported.Height {
	return cs.LatestHeight
}

// GetDelayTime returns the period of transaction confirmation delay.
func (cs ClientState) GetDelayTime() uint64 {
	return cs.TimeDelay
}

// GetDelayBlock returns the number of blocks delayed in transaction confirmation.
func (cs ClientState) GetDelayBlock() uint64 {
	return 0
}

// GetPrefix returns the prefix path for proof key.
func (cs ClientState) GetPrefix() exported.Prefix {
	return &cs.MerklePrefix
}

// IsExpired returns whether or not the client has passed the trusting period since the last
// update (in which case no headers are considered valid).
func (cs ClientState) IsExpired(latestTimestamp, now time.Time) bool {
	expirationTime := latestTimestamp.Add(cs.TrustingPeriod)
	return !expirationTime.After(now)
}

// Validate performs a basic validation of the client state fields.
func (cs ClientState) Validate() error {
	if strings.TrimSpace(cs.ChainId) == "" {
		return sdkerrors.Wrap(ErrInvalidChainID, "chain id cannot be empty string")
	}
	if err := light.ValidateTrustLevel(cs.TrustLevel.ToTendermint()); err != nil {
		return err
	}
	if cs.TrustingPeriod == 0 {
		return sdkerrors.Wrap(ErrInvalidTrustingPeriod, "trusting period cannot be zero")
	}
	if cs.UnbondingPeriod == 0 {
		return sdkerrors.Wrap(ErrInvalidUnbondingPeriod, "unbonding period cannot be zero")
	}
	if cs.MaxClockDrift == 0 {
		return sdkerrors.Wrap(ErrInvalidMaxClockDrift, "max clock drift cannot be zero")
	}
	if cs.LatestHeight.RevisionHeight == 0 {
		return sdkerrors.Wrapf(ErrInvalidHeaderHeight, "tendermint revision height cannot be zero")
	}
	if cs.TrustingPeriod >= cs.UnbondingPeriod {
		return sdkerrors.Wrapf(
			ErrInvalidTrustingPeriod,
			"trusting period (%s) should be < unbonding period (%s)",
			cs.TrustingPeriod, cs.UnbondingPeriod,
		)
	}

	if cs.ProofSpecs == nil {
		return sdkerrors.Wrap(ErrInvalidProofSpecs, "proof specs cannot be nil for tm client")
	}
	for i, spec := range cs.ProofSpecs {
		if spec == nil {
			return sdkerrors.Wrapf(ErrInvalidProofSpecs, "proof spec cannot be nil at index: %d", i)
		}
	}
	return nil
}

// GetProofSpecs returns the format the client expects for proof verification
// as a string array specifying the proof type for each position in chained proof
func (cs ClientState) GetProofSpecs() []*ics23.ProofSpec {
	return cs.ProofSpecs
}

// Initialize will check that initial consensus state is a Tendermint consensus state
// and will store ProcessedTime for initial consensus state as ctx.BlockTime()
func (cs ClientState) Initialize(ctx sdk.Context, _ codec.BinaryCodec, clientStore sdk.KVStore, consState exported.ConsensusState) error {
	if _, ok := consState.(*ConsensusState); !ok {
		return sdkerrors.Wrapf(
			clienttypes.ErrInvalidConsensus,
			"invalid initial consensus state. expected type: %T, got: %T",
			&ConsensusState{}, consState,
		)
	}
	// set processed time with initial consensus state height equal to initial client state's latest height
	setConsensusMetadata(ctx, clientStore, cs.GetLatestHeight())
	return nil
}

// Status function
// Clients must return their status. Only Active clients are allowed to process packets.
func (cs ClientState) Status(ctx sdk.Context, clientStore sdk.KVStore, cdc codec.BinaryCodec) exported.Status {
	// get latest consensus state from clientStore to check for expiry
	consState, err := GetConsensusState(clientStore, cdc, cs.GetLatestHeight())
	if err != nil {
		return exported.Unknown
	}

	if cs.IsExpired(consState.Timestamp, ctx.BlockTime()) {
		return exported.Expired
	}

	return exported.Active
}

// VerifyPacketCommitment verifies a proof of an outgoing packet commitment at
// the specified sourceChain, specified destChain, and specified sequence.
func (cs ClientState) VerifyPacketCommitment(
	ctx sdk.Context,
	store sdk.KVStore,
	cdc codec.BinaryCodec,
	height exported.Height,
	proof []byte,
	sourceChain,
	destChain string,
	sequence uint64,
	commitmentBytes []byte,
) error {
	merkleProof, consensusState, err := produceVerificationArgs(store, cdc, cs, height, cs.GetPrefix(), proof)
	if err != nil {
		return err
	}

	// check delay period has passed
	if err := verifyDelayPeriodPassed(ctx, store, height, cs.GetDelayTime()); err != nil {
		return err
	}

	commitmentPath := commitmenttypes.NewMerklePath(host.PacketCommitmentPath(sourceChain, destChain, sequence))
	path, err := commitmenttypes.ApplyPrefix(cs.GetPrefix(), commitmentPath)
	if err != nil {
		return err
	}

	if err := merkleProof.VerifyMembership(cs.ProofSpecs, consensusState.GetRoot(), path, commitmentBytes); err != nil {
		return err
	}

	return nil
}

// VerifyPacketAcknowledgement verifies a proof of an incoming packet
// acknowledgement at the specified sourceChain, specified destChain, and specified sequence.
func (cs ClientState) VerifyPacketAcknowledgement(
	ctx sdk.Context,
	store sdk.KVStore,
	cdc codec.BinaryCodec,
	height exported.Height,
	proof []byte,
	sourceChain,
	destChain string,
	sequence uint64,
	ackBytes []byte,
) error {
	merkleProof, consensusState, err := produceVerificationArgs(store, cdc, cs, height, cs.GetPrefix(), proof)
	if err != nil {
		return err
	}

	// check delay period has passed
	if err := verifyDelayPeriodPassed(ctx, store, height, cs.GetDelayTime()); err != nil {
		return err
	}

	ackPath := commitmenttypes.NewMerklePath(host.PacketAcknowledgementPath(sourceChain, destChain, sequence))
	path, err := commitmenttypes.ApplyPrefix(cs.GetPrefix(), ackPath)
	if err != nil {
		return err
	}

	if err := merkleProof.VerifyMembership(cs.ProofSpecs, consensusState.GetRoot(), path, ackBytes); err != nil {
		return err
	}

	return nil
}

// VerifyPacketCleanCommitment verifies a proof of an incoming packet
// acknowledgement at the specified sourceChain, specified destChain, and specified sequence.
func (cs ClientState) VerifyPacketCleanCommitment(
	ctx sdk.Context,
	store sdk.KVStore,
	cdc codec.BinaryCodec,
	height exported.Height,
	proof []byte,
	sourceChain string,
	destChain string,
	sequence uint64,
) error {
	merkleProof, consensusState, err := produceVerificationArgs(store, cdc, cs, height, cs.GetPrefix(), proof)
	if err != nil {
		return err
	}

	// check delay period has passed
	if err := verifyDelayPeriodPassed(ctx, store, height, cs.GetDelayTime()); err != nil {
		return err
	}

	cleanCommitmentPath := commitmenttypes.NewMerklePath(host.CleanPacketCommitmentPath(sourceChain, destChain))
	path, err := commitmenttypes.ApplyPrefix(cs.GetPrefix(), cleanCommitmentPath)
	if err != nil {
		return err
	}

	if err := merkleProof.VerifyMembership(cs.ProofSpecs, consensusState.GetRoot(), path, sdk.Uint64ToBigEndian(sequence)); err != nil {
		return err
	}

	return nil
}

// verifyDelayPeriodPassed will ensure that at least delayPeriod amount of time has passed since consensus state was submitted
// before allowing verification to continue.
func verifyDelayPeriodPassed(ctx sdk.Context, store sdk.KVStore, proofHeight exported.Height, delayPeriod uint64) error {
	// check that executing chain's timestamp has passed consensusState's processed time + delay period
	processedTime, ok := GetProcessedTime(store, proofHeight)
	if !ok {
		return sdkerrors.Wrapf(ErrProcessedTimeNotFound, "processed time not found for height: %s", proofHeight)
	}
	currentTimestamp := uint64(ctx.BlockTime().UnixNano())
	validTime := processedTime + delayPeriod
	// NOTE: delay period is inclusive, so if currentTimestamp is validTime, then we return no error
	if validTime > currentTimestamp {
		return sdkerrors.Wrapf(
			ErrDelayPeriodNotPassed,
			"cannot verify packet until time: %d, current time: %d",
			validTime, currentTimestamp,
		)
	}
	return nil
}

// produceVerificationArgs performs the basic checks on the arguments that are
// shared between the verification functions and returns the unmarshal
// merkle proof, the consensus state and an error if one occurred.
func produceVerificationArgs(
	store sdk.KVStore,
	cdc codec.BinaryCodec,
	cs ClientState,
	height exported.Height,
	prefix exported.Prefix,
	proof []byte,
) (merkleProof commitmenttypes.MerkleProof, consensusState *ConsensusState, err error) {
	if cs.GetLatestHeight().LT(height) {
		return commitmenttypes.MerkleProof{}, nil, sdkerrors.Wrapf(
			sdkerrors.ErrInvalidHeight,
			"client state height < proof height (%d < %d)",
			cs.GetLatestHeight(), height,
		)
	}

	if prefix == nil {
		return commitmenttypes.MerkleProof{}, nil, sdkerrors.Wrap(commitmenttypes.ErrInvalidPrefix, "prefix cannot be empty")
	}

	_, ok := prefix.(*commitmenttypes.MerklePrefix)
	if !ok {
		return commitmenttypes.MerkleProof{}, nil, sdkerrors.Wrapf(commitmenttypes.ErrInvalidPrefix, "invalid prefix type %T, expected *MerklePrefix", prefix)
	}

	if proof == nil {
		return commitmenttypes.MerkleProof{}, nil, sdkerrors.Wrap(commitmenttypes.ErrInvalidProof, "proof cannot be empty")
	}

	if err = cdc.Unmarshal(proof, &merkleProof); err != nil {
		return commitmenttypes.MerkleProof{}, nil, sdkerrors.Wrap(commitmenttypes.ErrInvalidProof, "failed to unmarshal proof into commitment merkle proof")
	}

	consensusState, err = GetConsensusState(store, cdc, height)
	if err != nil {
		return commitmenttypes.MerkleProof{}, nil, err
	}

	return merkleProof, consensusState, nil
}
