package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
)

const moduleName = host.ModuleName + "-" + SubModuleName

// TIBC packet sentinel errors
var (
	ErrSequenceSendNotFound     = sdkerrors.Register(moduleName, 2, "sequence send not found")
	ErrSequenceReceiveNotFound  = sdkerrors.Register(moduleName, 3, "sequence receive not found")
	ErrSequenceAckNotFound      = sdkerrors.Register(moduleName, 4, "sequence acknowledgement not found")
	ErrInvalidPacket            = sdkerrors.Register(moduleName, 5, "invalid packet")
	ErrInvalidAcknowledgement   = sdkerrors.Register(moduleName, 6, "invalid acknowledgement")
	ErrPacketCommitmentNotFound = sdkerrors.Register(moduleName, 7, "packet commitment not found")
	ErrPacketReceived           = sdkerrors.Register(moduleName, 8, "packet already received")
	ErrAcknowledgementExists    = sdkerrors.Register(moduleName, 9, "acknowledgement for packet already exists")
	ErrInvalidCleanPacket       = sdkerrors.Register(moduleName, 10, "invalid clean packet")
)
