package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	packettypes "github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
	commitmenttypes "github.com/bianjieai/tibc-go/modules/tibc/core/23-commitment/types"
	ibctmtypes "github.com/bianjieai/tibc-go/modules/tibc/light-clients/07-tendermint/types"
)

// RegisterInterfaces registers x/ibc interfaces into protobuf Any.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	clienttypes.RegisterInterfaces(registry)
	packettypes.RegisterInterfaces(registry)
	ibctmtypes.RegisterInterfaces(registry)
	commitmenttypes.RegisterInterfaces(registry)
}
