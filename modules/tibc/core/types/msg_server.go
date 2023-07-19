package types

import (
	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	packettypes "github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
	routingtypes "github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
)

type MsgServer interface {
	clienttypes.MsgServer
	packettypes.MsgServer
	routingtypes.MsgServer
}
