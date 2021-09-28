package tibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	client "github.com/bianjieai/tibc-go/modules/tibc/core/02-client"
	packet "github.com/bianjieai/tibc-go/modules/tibc/core/04-packet"
	routing "github.com/bianjieai/tibc-go/modules/tibc/core/26-routing"
	"github.com/bianjieai/tibc-go/modules/tibc/core/keeper"
	"github.com/bianjieai/tibc-go/modules/tibc/core/types"
)

// InitGenesis initializes the tibc state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, createLocalhost bool, gs *types.GenesisState) {
	client.InitGenesis(ctx, k.ClientKeeper, gs.ClientGenesis)
	packet.InitGenesis(ctx, k.PacketKeeper, gs.PacketGenesis)
	routing.InitGenesis(ctx, k.RoutingKeeper, gs.RoutingGenesis)
}

// ExportGenesis returns the tibc exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		ClientGenesis:  client.ExportGenesis(ctx, k.ClientKeeper),
		PacketGenesis:  packet.ExportGenesis(ctx, k.PacketKeeper),
		RoutingGenesis: routing.ExportGenesis(ctx, k.RoutingKeeper),
	}
}
