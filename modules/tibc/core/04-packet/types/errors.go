package types

import (
	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const moduleName = host.ModuleName + "-" + SubModuleName

// IBC channel sentinel errors
var (
	ErrChannelExists             = sdkerrors.Register(moduleName, 2, "channel already exists")
	ErrChannelNotFound           = sdkerrors.Register(moduleName, 3, "channel not found")
	ErrInvalidChannel            = sdkerrors.Register(moduleName, 4, "invalid channel")
	ErrInvalidChannelState       = sdkerrors.Register(moduleName, 5, "invalid channel state")
	ErrInvalidChannelOrdering    = sdkerrors.Register(moduleName, 6, "invalid channel ordering")
	ErrInvalidCounterparty       = sdkerrors.Register(moduleName, 7, "invalid counterparty channel")
	ErrInvalidChannelCapability  = sdkerrors.Register(moduleName, 8, "invalid channel capability")
	ErrChannelCapabilityNotFound = sdkerrors.Register(moduleName, 9, "channel capability not found")
	ErrSequenceSendNotFound      = sdkerrors.Register(moduleName, 10, "sequence send not found")
	ErrSequenceReceiveNotFound   = sdkerrors.Register(moduleName, 11, "sequence receive not found")
	ErrSequenceAckNotFound       = sdkerrors.Register(moduleName, 12, "sequence acknowledgement not found")
	ErrInvalidPacket             = sdkerrors.Register(moduleName, 13, "invalid packet")
	ErrPacketTimeout             = sdkerrors.Register(moduleName, 14, "packet timeout")
	ErrTooManyConnectionHops     = sdkerrors.Register(moduleName, 15, "too many connection hops")
	ErrInvalidAcknowledgement    = sdkerrors.Register(moduleName, 16, "invalid acknowledgement")
	ErrPacketCommitmentNotFound  = sdkerrors.Register(moduleName, 17, "packet commitment not found")
	ErrPacketReceived            = sdkerrors.Register(moduleName, 18, "packet already received")
	ErrAcknowledgementExists     = sdkerrors.Register(moduleName, 19, "acknowledgement for packet already exists")
	ErrInvalidChannelIdentifier  = sdkerrors.Register(moduleName, 20, "invalid channel identifier")
)
