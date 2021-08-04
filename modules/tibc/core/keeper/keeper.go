package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	clientkeeper "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/keeper"
	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	packetkeeper "github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/keeper"
	routingkeeper "github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/keeper"
	routingtypes "github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
	"github.com/bianjieai/tibc-go/modules/tibc/core/types"
)

var _ types.QueryServer = (*Keeper)(nil)

// Keeper defines each ICS keeper for IBC
type Keeper struct {
	// implements gRPC QueryServer interface
	types.QueryServer

	cdc codec.BinaryMarshaler

	ClientKeeper  clientkeeper.Keeper
	Packetkeeper  packetkeeper.Keeper
	RoutingKeeper routingkeeper.Keeper
	Router        *routingtypes.Router
}

// NewKeeper creates a new ibc Keeper
func NewKeeper(
	cdc codec.BinaryMarshaler, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	stakingKeeper clienttypes.StakingKeeper, scopedKeeper capabilitykeeper.ScopedKeeper,
) *Keeper {
	clientKeeper := clientkeeper.NewKeeper(cdc, key, paramSpace, stakingKeeper)
	routingKeeper := routingkeeper.NewKeeper()
	packetkeeper := packetkeeper.NewKeeper(cdc, key, clientKeeper, routingKeeper, scopedKeeper)

	return &Keeper{
		cdc:           cdc,
		ClientKeeper:  clientKeeper,
		Packetkeeper:  packetkeeper,
		RoutingKeeper: routingKeeper,
	}
}

// Codec returns the IBC module codec.
func (k Keeper) Codec() codec.BinaryMarshaler {
	return k.cdc
}

// SetRouter sets the Router in IBC Keeper and seals it. The method panics if
// there is an existing router that's already sealed.
func (k *Keeper) SetRouter(rtr *routingtypes.Router) {
	if k.Router != nil && k.Router.Sealed() {
		panic("cannot reset a sealed router")
	}
	k.Router = rtr
	k.Router.Seal()
}
