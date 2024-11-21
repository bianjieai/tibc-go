package types

import (
	errorsmod "cosmossdk.io/errors"
	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
)

const moduleName = host.ModuleName + "-" + SubModuleName

// TIBC packet sentinel errors
var (
	ErrSequenceSendNotFound     = errorsmod.Register(moduleName, 2, "sequence send not found")
	ErrSequenceReceiveNotFound  = errorsmod.Register(moduleName, 3, "sequence receive not found")
	ErrSequenceAckNotFound      = errorsmod.Register(moduleName, 4, "sequence acknowledgement not found")
	ErrInvalidPacket            = errorsmod.Register(moduleName, 5, "invalid packet")
	ErrInvalidAcknowledgement   = errorsmod.Register(moduleName, 6, "invalid acknowledgement")
	ErrPacketCommitmentNotFound = errorsmod.Register(moduleName, 7, "packet commitment not found")
	ErrPacketReceived           = errorsmod.Register(moduleName, 8, "packet already received")
	ErrAcknowledgementExists    = errorsmod.Register(moduleName, 9, "acknowledgement for packet already exists")
	ErrInvalidCleanPacket       = errorsmod.Register(moduleName, 10, "invalid clean packet")
)
