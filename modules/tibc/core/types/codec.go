package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	packettypes "github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
	commitmenttypes "github.com/bianjieai/tibc-go/modules/tibc/core/23-commitment/types"
	routingtypes "github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
	ibctmtypes "github.com/bianjieai/tibc-go/modules/tibc/light-clients/07-tendermint/types"
	bsctypes "github.com/bianjieai/tibc-go/modules/tibc/light-clients/08-bsc/types"
	ethtypes "github.com/bianjieai/tibc-go/modules/tibc/light-clients/09-eth/types"
)

// RegisterInterfaces registers x/ibc interfaces into protobuf Any.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	clienttypes.RegisterInterfaces(registry)
	packettypes.RegisterInterfaces(registry)
	routingtypes.RegisterInterfaces(registry)
	ibctmtypes.RegisterInterfaces(registry)
	bsctypes.RegisterInterfaces(registry)
	ethtypes.RegisterInterfaces(registry)
	commitmenttypes.RegisterInterfaces(registry)
}
