package exported

import (
	proto "github.com/gogo/protobuf/proto"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Status represents the status of a client
type Status string

const (
	// TypeClientMisbehaviour is the shared evidence misbehaviour type
	TypeClientMisbehaviour string = "client_misbehaviour"

	// Tendermint is used to indicate that the client uses the Tendermint Consensus Algorithm.
	Tendermint string = "007-tendermint"

	// BSC is the client type for a bianance smart chain client.
	BSC string = "008-bsc"

	// ETH is the client type for a Ethereum client.
	ETH string = "009-eth"

	// Fabric is the client type for a hyperledge fabric client.
	// Fabric string = "009-fabric"

	// Active is a status type of a client. An active client is allowed to be used.
	Active Status = "Active"

	// Expired is a status type of a client. An expired client is not allowed to be used.
	Expired Status = "Expired"

	// Unknown indicates there was an error in determining the status of a client.
	Unknown Status = "Unknown"
)

// ClientState defines the required common functions for light clients.
type ClientState interface {
	proto.Message

	ClientType() string
	GetLatestHeight() Height
	Validate() error
	GetDelayTime() uint64
	GetDelayBlock() uint64
	GetPrefix() Prefix

	// Initialize function
	// Clients must validate the initial consensus state, and may store any client-specific metadata
	// necessary for correct light client operation
	Initialize(sdk.Context, codec.BinaryCodec, sdk.KVStore, ConsensusState) error

	// Status function
	// Clients must return their status. Only Active clients are allowed to process packets.
	Status(ctx sdk.Context, clientStore sdk.KVStore, cdc codec.BinaryCodec) Status

	// ExportMetadata function
	ExportMetadata(sdk.KVStore) []GenesisMetadata

	// Update and Misbehaviour functions
	CheckHeaderAndUpdateState(sdk.Context, codec.BinaryCodec, sdk.KVStore, Header) (ClientState, ConsensusState, error)

	// State verification functions
	// Verify the commitment of the cross-chain data package
	VerifyPacketCommitment(
		ctx sdk.Context,
		store sdk.KVStore,
		cdc codec.BinaryCodec,
		height Height,
		proof []byte,
		sourceChain,
		destChain string,
		sequence uint64,
		commitmentBytes []byte,
	) error

	// Verify the Acknowledgement of the cross-chain data package
	VerifyPacketAcknowledgement(
		ctx sdk.Context,
		store sdk.KVStore,
		cdc codec.BinaryCodec,
		height Height,
		proof []byte,
		sourceChain,
		destChain string,
		sequence uint64,
		ackBytes []byte,
	) error

	// Verify the CleanCommitment of the cross-chain data package
	VerifyPacketCleanCommitment(
		ctx sdk.Context,
		store sdk.KVStore,
		cdc codec.BinaryCodec,
		height Height,
		proof []byte,
		sourceChain string,
		destChain string,
		sequence uint64,
	) error
}

// ConsensusState is the state of the consensus process
type ConsensusState interface {
	proto.Message

	ClientType() string // Consensus kind

	// GetRoot returns the commitment root of the consensus state,
	// which is used for key-value pair verification.
	GetRoot() Root

	// GetTimestamp returns the timestamp (in nanoseconds) of the consensus state
	GetTimestamp() uint64

	ValidateBasic() error
}

// Header is the consensus state update information
type Header interface {
	proto.Message

	ClientType() string
	GetHeight() Height
	ValidateBasic() error
}

// Height is a wrapper interface over clienttypes.Height
// all clients must use the concrete implementation in types
type Height interface {
	IsZero() bool
	LT(Height) bool
	LTE(Height) bool
	EQ(Height) bool
	GT(Height) bool
	GTE(Height) bool
	GetRevisionNumber() uint64
	GetRevisionHeight() uint64
	Increment() Height
	Decrement() (Height, bool)
	String() string
}

// GenesisMetadata is a wrapper interface over clienttypes.GenesisMetadata
// all clients must use the concrete implementation in types
type GenesisMetadata interface {
	// return store key that contains metadata without chainName-prefix
	GetKey() []byte
	// returns metadata value
	GetValue() []byte
}
