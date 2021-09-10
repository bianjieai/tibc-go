package routing

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/keeper"
	"github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
)

// InitGenesis initializes the tibc routing submodule's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, gs types.GenesisState) {
	err := k.SetRoutingRules(ctx, gs.Rules)
	if err != nil {
		panic(err)
	}
}

// ExportGenesis returns the tibc routing submodule's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {
	rules, _ := k.GetRoutingRules(ctx)
	return types.GenesisState{
		Rules: rules,
	}
}
