package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	clientkeeper "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/keeper"
	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	packetkeeper "github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/keeper"
	routingkeeper "github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/keeper"
	routingtypes "github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
	"github.com/bianjieai/tibc-go/modules/tibc/core/types"
)

var _ types.QueryServer = (*Keeper)(nil)

// Keeper defines each TICS keeper for TIBC
type Keeper struct {
	// implements gRPC QueryServer interface
	types.QueryServer

	cdc codec.BinaryCodec

	ClientKeeper  clientkeeper.Keeper
	PacketKeeper  packetkeeper.Keeper
	RoutingKeeper routingkeeper.Keeper
}

// NewKeeper creates a new tibc Keeper
func NewKeeper(
	cdc codec.BinaryCodec, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	stakingKeeper clienttypes.StakingKeeper,
) *Keeper {
	clientKeeper := clientkeeper.NewKeeper(cdc, key, paramSpace, stakingKeeper)
	routingKeeper := routingkeeper.NewKeeper(key)
	packetkeeper := packetkeeper.NewKeeper(cdc, key, clientKeeper, routingKeeper)

	return &Keeper{
		cdc:           cdc,
		ClientKeeper:  clientKeeper,
		PacketKeeper:  packetkeeper,
		RoutingKeeper: routingKeeper,
	}
}

// Codec returns the TIBC module codec.
func (k Keeper) Codec() codec.BinaryCodec {
	return k.cdc
}

// SetRouter sets the Router in TIBC Keeper and seals it. The method panics if
// there is an existing router that's already sealed.
func (k *Keeper) SetRouter(rtr *routingtypes.Router) {
	k.RoutingKeeper.SetRouter(rtr)
}
