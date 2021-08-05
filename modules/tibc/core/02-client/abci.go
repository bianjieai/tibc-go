package client

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/keeper"
)

// BeginBlocker updates an existing localhost client with the latest block height.
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	// _, found := k.GetClientState(ctx, exported.Localhost)
	// if !found {
	// 	return
	// }

	// // update the localhost client with the latest block height
	// if err := k.UpdateClient(ctx, exported.Localhost, nil); err != nil {
	// 	panic(err)
	// }
}
