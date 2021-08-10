package keeper

import (
	"github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/types"
	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"
)



type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        codec.BinaryMarshaler
	paramSpace paramtypes.Subspace

	ak	types.AccountKeeper
	nk  types.NftKeeper
	pk  types.PacketKeeper
	ck  types.ClientKeeper
}


// NewKeeper creates a new IBC transfer Keeper instance
func NewKeeper(
	cdc codec.BinaryMarshaler, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	ak	types.AccountKeeper, nk  types.NftKeeper,
	pk  types.PacketKeeper, ck  types.ClientKeeper,
) Keeper {

	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the TIBC nft-transfer module account has not been set")
	}

	return Keeper{
		cdc:           cdc,
		storeKey:      key,
		paramSpace:    paramSpace,
		ak:    ak,
		nk:    nk,
		pk:    pk,
		ck:    ck,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+host.ModuleName+"-"+types.ModuleName)
}


// GetNftTransferModuleAddr returns the nft transfer module addr
func (k Keeper) GetNftTransferModuleAddr(name string) sdk.AccAddress {
	return k.ak.GetModuleAddress(name)
}
