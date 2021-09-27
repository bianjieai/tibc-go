package types

import (
	"crypto/sha256"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

// CommitPacket returns the packet commitment bytes. The commitment consists of:
// sha256_hash(timeout_timestamp + timeout_height.RevisionNumber + timeout_height.RevisionHeight + sha256_hash(data))
// from a given packet. This results in a fixed length preimage.
// NOTE: sdk.Uint64ToBigEndian sets the uint64 to a slice of length 8.
func CommitPacket(packet exported.PacketI) []byte {
	dataHash := sha256.Sum256(packet.GetData())
	return dataHash[:]
}

// CommitAcknowledgement returns the hash of commitment bytes
func CommitAcknowledgement(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

var (
	_         exported.PacketI = (*Packet)(nil)
	ModuleCdc                  = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
)

// NewPacket creates a new Packet instance. It panics if the provided
// packet data interface is not registered.
func NewPacket(
	data []byte, sequence uint64, sourceChain, destinationChain, relayChain, port string,
) Packet {
	return Packet{
		Data:             data,
		Sequence:         sequence,
		SourceChain:      sourceChain,
		DestinationChain: destinationChain,
		RelayChain:       relayChain,
		Port:             port,
	}
}

// GetSequence implements PacketI interface
func (p Packet) GetSequence() uint64 { return p.Sequence }

// GetPort implements PacketI interface
func (p Packet) GetPort() string { return p.Port }

// GetSourceChain implements PacketI interface
func (p Packet) GetSourceChain() string { return p.SourceChain }

// GetDestinationChain implements PacketI interface
func (p Packet) GetDestChain() string { return p.DestinationChain }

// GetRelayChain implements PacketI interface
func (p Packet) GetRelayChain() string { return p.RelayChain }

// GetData implements PacketI interface
func (p Packet) GetData() []byte { return p.Data }

// ValidateBasic implements PacketI interface
func (p Packet) ValidateBasic() error {
	if p.Sequence == 0 {
		return sdkerrors.Wrap(ErrInvalidPacket, "packet sequence cannot be 0")
	}
	if len(p.Data) == 0 {
		return sdkerrors.Wrap(ErrInvalidPacket, "packet data bytes cannot be empty")
	}
	return nil
}

var _ exported.CleanPacketI = (*CleanPacket)(nil)

// NewCleanPacket creates a new CleanPacket instance.
func NewCleanPacket(
	sequence uint64, sourceChain, destinationChain, relayChain string,
) CleanPacket {
	return CleanPacket{
		Sequence:         sequence,
		SourceChain:      sourceChain,
		DestinationChain: destinationChain,
		RelayChain:       relayChain,
	}
}

// GetSequence implements PacketI interface
func (p CleanPacket) GetSequence() uint64 { return p.Sequence }

// GetSourceChain implements PacketI interface
func (p CleanPacket) GetSourceChain() string { return p.SourceChain }

// GetDestinationChain implements PacketI interface
func (p CleanPacket) GetDestChain() string { return p.DestinationChain }

// GetRelayChain implements PacketI interface
func (p CleanPacket) GetRelayChain() string { return p.RelayChain }

// ValidateBasic implements PacketI interface
func (p CleanPacket) ValidateBasic() error {
	if p.Sequence == 0 {
		return sdkerrors.Wrap(ErrInvalidPacket, "packet sequence cannot be 0")
	}
	return nil
}

// NewResultAcknowledgement returns a new instance of Acknowledgement using an Acknowledgement_Result
// type in the Response field.
func NewResultAcknowledgement(result []byte) Acknowledgement {
	return Acknowledgement{
		Response: &Acknowledgement_Result{
			Result: result,
		},
	}
}

// NewErrorAcknowledgement returns a new instance of Acknowledgement using an Acknowledgement_Error
// type in the Response field.
func NewErrorAcknowledgement(err string) Acknowledgement {
	return Acknowledgement{
		Response: &Acknowledgement_Error{
			Error: err,
		},
	}
}

// GetBytes is a helper for serialising acknowledgements
func (ack Acknowledgement) GetBytes() []byte {
	return ModuleCdc.MustMarshal(&ack)
}
