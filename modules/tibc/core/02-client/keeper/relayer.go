package keeper

import (
	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// SetRelayers saves the relayers under the specified chainname
func (k Keeper) SetRelayers(ctx sdk.Context, chainName string, relayers []string) error {
	store := k.ClientStore(ctx, chainName)
	for _, relayer := range relayers {
		if k.HasRelayer(ctx, chainName, relayer) {
			return sdkerrors.Wrapf(types.ErrRelayerExists, "relayer %s", relayer)
		}
		store.Set(types.RelayerKey(relayer), []byte(relayer))
	}
	return nil
}

// HasRelayer asserts whether a relay is already registered
func (k Keeper) HasRelayer(ctx sdk.Context, chainName string, relayer string) bool {
	store := k.ClientStore(ctx, chainName)
	bz := store.Get(types.RelayerKey(relayer))
	return len(bz) > 0
}

// GetRelayers returns all registered relay addresses under the specified chain name
func (k Keeper) GetRelayers(ctx sdk.Context, chainName string, relayer string) (relayers []string) {
	store := k.ClientStore(ctx, chainName)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.KeyRelayers))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		relayers = append(relayers, string(iterator.Value()))
	}
	return
}
