package keeper

import (
	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RegisterRelayers saves the relayers under the specified chainname
func (k Keeper) RegisterRelayers(ctx sdk.Context, chainName string, relayers []string) {
	store := k.RelayerStore(ctx)
	ir := &types.IdentifiedRelayers{
		ChainName: chainName,
		Relayers:  relayers,
	}
	irBz := k.cdc.MustMarshalBinaryBare(ir)
	store.Set([]byte(chainName), irBz)
}

// AuthRelayer asserts whether a relay is already registered
func (k Keeper) AuthRelayer(ctx sdk.Context, chainName string, relayer string) bool {
	for _, r := range k.GetRelayers(ctx, chainName) {
		if r == relayer {
			return true
		}
	}
	return false
}

// GetRelayers returns all registered relay addresses under the specified chain name
func (k Keeper) GetRelayers(ctx sdk.Context, chainName string) (relayers []string) {
	store := k.RelayerStore(ctx)
	bz := store.Get([]byte(chainName))

	var ir = &types.IdentifiedRelayers{}
	k.cdc.MustUnmarshalBinaryBare(bz, ir)
	return ir.Relayers
}

// GetAllRelayers returns all registered relay addresses
func (k Keeper) GetAllRelayers(ctx sdk.Context) (relayers []types.IdentifiedRelayers) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.KeyRelayers))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var ir = &types.IdentifiedRelayers{}
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), ir)
		relayers = append(relayers, *ir)
	}
	return
}