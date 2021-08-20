package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"

	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

// RegisterInterfaces register the tibc packet submodule interfaces to protobuf
// Any.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterInterface(
		"tibc.core.packet.v1.PacketI",
		(*exported.PacketI)(nil),
	)
	registry.RegisterInterface(
		"tibc.core.packet.v1.CleanPacketI",
		(*exported.CleanPacketI)(nil),
	)
	registry.RegisterImplementations(
		(*exported.PacketI)(nil),
		&Packet{},
	)
	registry.RegisterImplementations(
		(*exported.CleanPacketI)(nil),
		&CleanPacket{},
	)
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgRecvPacket{},
		&MsgAcknowledgement{},
		&MsgCleanPacket{},
		&MsgRecvCleanPacket{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// SubModuleCdc references the global x/ibc/core/04-channel module codec. Note, the codec should
// ONLY be used in certain instances of tests and for JSON encoding.
//
// The actual codec used for serialization should be provided to x/ibc/core/04-channel and
// defined at the application level.
var SubModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
