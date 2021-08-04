package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	packettypes "github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
)

// TIBCModule defines an interface that implements all the callbacks
// that modules must define as specified in TICS-26
type TIBCModule interface {
	// OnRecvPacket must return the acknowledgement bytes
	// In the case of an asynchronous acknowledgement, nil should be returned.
	OnRecvPacket(
		ctx sdk.Context,
		packet packettypes.Packet,
	) (*sdk.Result, []byte, error)

	OnAcknowledgementPacket(
		ctx sdk.Context,
		packet packettypes.Packet,
		acknowledgement []byte,
	) (*sdk.Result, error)
}
