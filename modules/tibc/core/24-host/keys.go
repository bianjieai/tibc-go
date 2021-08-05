package host

import (
	"fmt"

	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

const (
	// ModuleName is the name of the IBC module
	ModuleName = "tibc"

	// StoreKey is the string store representation
	StoreKey string = ModuleName

	// QuerierRoute is the querier route for the IBC module
	QuerierRoute string = ModuleName

	// RouterKey is the msg router key for the IBC module
	RouterKey string = ModuleName
)

// KVStore key prefixes for IBC
var (
	KeyClientStorePrefix = []byte("clients")
)

// KVStore key prefixes for IBC
const (
	KeyClientState             = "clientState"
	KeyConsensusStatePrefix    = "consensusStates"
	KeyConnectionPrefix        = "connections"
	KeyChannelEndPrefix        = "channelEnds"
	KeyChannelPrefix           = "channels"
	KeyPortPrefix              = "ports"
	KeySequencePrefix          = "sequences"
	KeyChannelCapabilityPrefix = "capabilities"
	KeyNextSeqSendPrefix       = "nextSequenceSend"
	KeyNextSeqRecvPrefix       = "nextSequenceRecv"
	KeyNextSeqAckPrefix        = "nextSequenceAck"
	KeyPacketCommitmentPrefix  = "commitments"
	KeyPacketAckPrefix         = "acks"
	KeyPacketReceiptPrefix     = "receipts"
)

// FullClientPath returns the full path of a specific client path in the format:
// "clients/{clientID}/{path}" as a string.
func FullClientPath(clientID string, path string) string {
	return fmt.Sprintf("%s/%s/%s", KeyClientStorePrefix, clientID, path)
}

// FullClientKey returns the full path of specific client path in the format:
// "clients/{clientID}/{path}" as a byte array.
func FullClientKey(clientID string, path []byte) []byte {
	return []byte(FullClientPath(clientID, string(path)))
}

// ICS02
// The following paths are the keys to the store as defined in https://github.com/cosmos/ics/tree/master/spec/ics-002-client-semantics#path-space

// FullClientStatePath takes a client identifier and returns a Path under which to store a
// particular client state
func FullClientStatePath(clientID string) string {
	return FullClientPath(clientID, KeyClientState)
}

// FullClientStateKey takes a client identifier and returns a Key under which to store a
// particular client state.
func FullClientStateKey(clientID string) []byte {
	return FullClientKey(clientID, []byte(KeyClientState))
}

// ClientStateKey returns a store key under which a particular client state is stored
// in a client prefixed store
func ClientStateKey() []byte {
	return []byte(KeyClientState)
}

// FullConsensusStatePath takes a client identifier and returns a Path under which to
// store the consensus state of a client.
func FullConsensusStatePath(clientID string, height exported.Height) string {
	return FullClientPath(clientID, ConsensusStatePath(height))
}

// FullConsensusStateKey returns the store key for the consensus state of a particular
// client.
func FullConsensusStateKey(clientID string, height exported.Height) []byte {
	return []byte(FullConsensusStatePath(clientID, height))
}

// ConsensusStatePath returns the suffix store key for the consensus state at a
// particular height stored in a client prefixed store.
func ConsensusStatePath(height exported.Height) string {
	return fmt.Sprintf("%s/%s", KeyConsensusStatePrefix, height)
}

// ConsensusStateKey returns the store key for a the consensus state of a particular
// client stored in a client prefixed store.
func ConsensusStateKey(height exported.Height) []byte {
	return []byte(ConsensusStatePath(height))
}

// ICS03
// The following paths are the keys to the store as defined in https://github.com/cosmos/ics/tree/master/spec/ics-003-connection-semantics#store-paths

// ClientConnectionsPath defines a reverse mapping from clients to a set of connections
func ClientConnectionsPath(clientID string) string {
	return FullClientPath(clientID, KeyConnectionPrefix)
}

// ClientConnectionsKey returns the store key for the connections of a given client
func ClientConnectionsKey(clientID string) []byte {
	return []byte(ClientConnectionsPath(clientID))
}

// ConnectionPath defines the path under which connection paths are stored
func ConnectionPath(connectionID string) string {
	return fmt.Sprintf("%s/%s", KeyConnectionPrefix, connectionID)
}

// ConnectionKey returns the store key for a particular connection
func ConnectionKey(connectionID string) []byte {
	return []byte(ConnectionPath(connectionID))
}

// ICS04
// The following paths are the keys to the store as defined in https://github.com/cosmos/ics/tree/master/spec/ics-004-channel-and-packet-semantics#store-paths

// ChannelPath defines the path under which channels are stored
func ChannelPath(sourceChain, destinationChain string) string {
	return fmt.Sprintf("%s/%s", KeyChannelEndPrefix, channelPath(sourceChain, destinationChain))
}

// ChannelKey returns the store key for a particular channel
func ChannelKey(sourceChain, destinationChain string) []byte {
	return []byte(ChannelPath(sourceChain, destinationChain))
}

// ChannelCapabilityPath defines the path under which capability keys associated
// with a channel are stored
func ChannelCapabilityPath(port string) string {
	return fmt.Sprintf("%s/%s", KeyChannelCapabilityPrefix, port)
}

// NextSequenceSendPath defines the next send sequence counter store path
func NextSequenceSendPath(sourceChain, destChain string) string {
	return fmt.Sprintf("%s/%s", KeyNextSeqSendPrefix, channelPath(sourceChain, destChain))
}

// NextSequenceSendKey returns the store key for the send sequence of a particular
// channel binded to a specific port.
func NextSequenceSendKey(sourceChain, destChain string) []byte {
	return []byte(NextSequenceSendPath(sourceChain, destChain))
}

// NextSequenceRecvPath defines the next receive sequence counter store path.
func NextSequenceRecvPath(sourceChain, destinationChain string) string {
	return fmt.Sprintf("%s/%s", KeyNextSeqRecvPrefix, channelPath(sourceChain, destinationChain))
}

// NextSequenceRecvKey returns the store key for the receive sequence of a particular
// channel binded to a specific port
func NextSequenceRecvKey(sourceChain, destinationChain string) []byte {
	return []byte(NextSequenceRecvPath(sourceChain, destinationChain))
}

// NextSequenceAckPath defines the next acknowledgement sequence counter store path
func NextSequenceAckPath(sourceChain, destinationChain string) string {
	return fmt.Sprintf("%s/%s", KeyNextSeqAckPrefix, channelPath(sourceChain, destinationChain))
}

// NextSequenceAckKey returns the store key for the acknowledgement sequence of
// a particular channel binded to a specific port.
func NextSequenceAckKey(sourceChain, destinationChain string) []byte {
	return []byte(NextSequenceAckPath(sourceChain, destinationChain))
}

// PacketCommitmentPath defines the commitments to packet data fields store path
func PacketCommitmentPath(sourceChain, destinationChain string, sequence uint64) string {
	return fmt.Sprintf("%s/%d", PacketCommitmentPrefixPath(sourceChain, destinationChain), sequence)
}

// PacketCommitmentKey returns the store key of under which a packet commitment
// is stored
func PacketCommitmentKey(sourceChain, destinationChain string, sequence uint64) []byte {
	return []byte(PacketCommitmentPath(sourceChain, destinationChain, sequence))
}

// PacketCommitmentPrefixPath defines the prefix for commitments to packet data fields store path.
func PacketCommitmentPrefixPath(sourceChain, destinationChain string) string {
	return fmt.Sprintf("%s/%s/%s", KeyPacketCommitmentPrefix, channelPath(sourceChain, destinationChain), KeySequencePrefix)
}

// PacketAcknowledgementPath defines the packet acknowledgement store path
func PacketAcknowledgementPath(sourceChain, destinationChain string, sequence uint64) string {
	return fmt.Sprintf("%s/%d", PacketAcknowledgementPrefixPath(sourceChain, destinationChain), sequence)
}

// PacketAcknowledgementKey returns the store key of under which a packet
// acknowledgement is stored
func PacketAcknowledgementKey(sourceChain, destinationChain string, sequence uint64) []byte {
	return []byte(PacketAcknowledgementPath(sourceChain, destinationChain, sequence))
}

// PacketAcknowledgementPrefixPath defines the prefix for commitments to packet data fields store path.
func PacketAcknowledgementPrefixPath(sourceChain, destinationChain string) string {
	return fmt.Sprintf("%s/%s/%s", KeyPacketAckPrefix, channelPath(sourceChain, destinationChain), KeySequencePrefix)
}

// PacketReceiptPath defines the packet receipt store path
func PacketReceiptPath(sourceChain, destinationChain string, sequence uint64) string {
	return fmt.Sprintf("%s/%s/%s", KeyPacketReceiptPrefix, channelPath(sourceChain, destinationChain), sequencePath(sequence))
}

// PacketReceiptKey returns the store key of under which a packet
// receipt is stored
func PacketReceiptKey(sourceChain, destinationChain string, sequence uint64) []byte {
	return []byte(PacketReceiptPath(sourceChain, destinationChain, sequence))
}

func channelPath(sourceChain, destinationChain string) string {
	return fmt.Sprintf("%s/%s/%s/%s", KeyPortPrefix, sourceChain, KeyChannelPrefix, destinationChain)
}

func sequencePath(sequence uint64) string {
	return fmt.Sprintf("%s/%d", KeySequencePrefix, sequence)
}

// ICS05
// The following paths are the keys to the store as defined in https://github.com/cosmos/ics/tree/master/spec/ics-026-routing-allocation#store-paths

// PortPath defines the path under which ports paths are stored on the capability module
func PortPath(portID string) string {
	return fmt.Sprintf("%s/%s", KeyPortPrefix, portID)
}
