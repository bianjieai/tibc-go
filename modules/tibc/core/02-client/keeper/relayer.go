package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
)

// RegisterRelayers saves the relayers under the specified chainname
func (k Keeper) RegisterRelayers(ctx sdk.Context, chainName string, relayers []string) {
	store := k.RelayerStore(ctx)
	ir := &types.IdentifiedRelayers{
		ChainName: chainName,
		Relayers:  relayers,
	}
	irBz := k.cdc.MustMarshal(ir)
	store.Set([]byte(chainName), irBz)
}

// AuthRelayer asserts whether a relayer is already registered
func (k Keeper) AuthRelayer(ctx sdk.Context, chainName string, relayer string) bool {
	for _, r := range k.GetRelayers(ctx, chainName) {
		if r == relayer {
			return true
		}
	}
	return false
}

// GetRelayers returns all registered relayer addresses under the specified chain name
func (k Keeper) GetRelayers(ctx sdk.Context, chainName string) (relayers []string) {
	store := k.RelayerStore(ctx)
	bz := store.Get([]byte(chainName))

	var ir = &types.IdentifiedRelayers{}
	k.cdc.MustUnmarshal(bz, ir)
	return ir.Relayers
}

// GetAllRelayers returns all registered relayer addresses
func (k Keeper) GetAllRelayers(ctx sdk.Context) (relayers []types.IdentifiedRelayers) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.KeyRelayers))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var ir = &types.IdentifiedRelayers{}
		k.cdc.MustUnmarshal(iterator.Value(), ir)
		relayers = append(relayers, *ir)
	}
	return
}
