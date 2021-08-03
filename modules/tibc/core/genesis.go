package ibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	client "github.com/bianjieai/tibc-go/modules/tibc/core/02-client"
	connection "github.com/bianjieai/tibc-go/modules/tibc/core/03-connection"
	channel "github.com/bianjieai/tibc-go/modules/tibc/core/04-channel"
	"github.com/bianjieai/tibc-go/modules/tibc/core/keeper"
	"github.com/bianjieai/tibc-go/modules/tibc/core/types"
)

// InitGenesis initializes the ibc state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, createLocalhost bool, gs *types.GenesisState) {
	client.InitGenesis(ctx, k.ClientKeeper, gs.ClientGenesis)
	connection.InitGenesis(ctx, k.ConnectionKeeper, gs.ConnectionGenesis)
	channel.InitGenesis(ctx, k.ChannelKeeper, gs.ChannelGenesis)
}

// ExportGenesis returns the ibc exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		ClientGenesis:     client.ExportGenesis(ctx, k.ClientKeeper),
		ConnectionGenesis: connection.ExportGenesis(ctx, k.ConnectionKeeper),
		ChannelGenesis:    channel.ExportGenesis(ctx, k.ChannelKeeper),
	}
}
