package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
)

const moduleName = host.ModuleName + "-" + SubModuleName

// TIBC packet sentinel errors
var (
	ErrSequenceSendNotFound     = sdkerrors.Register(moduleName, 10, "sequence send not found")
	ErrSequenceReceiveNotFound  = sdkerrors.Register(moduleName, 11, "sequence receive not found")
	ErrSequenceAckNotFound      = sdkerrors.Register(moduleName, 12, "sequence acknowledgement not found")
	ErrInvalidPacket            = sdkerrors.Register(moduleName, 13, "invalid packet")
	ErrInvalidAcknowledgement   = sdkerrors.Register(moduleName, 16, "invalid acknowledgement")
	ErrPacketCommitmentNotFound = sdkerrors.Register(moduleName, 17, "packet commitment not found")
	ErrPacketReceived           = sdkerrors.Register(moduleName, 18, "packet already received")
	ErrAcknowledgementExists    = sdkerrors.Register(moduleName, 19, "acknowledgement for packet already exists")
	ErrInvalidCleanPacket       = sdkerrors.Register(moduleName, 21, "invalid clean packet")
)
