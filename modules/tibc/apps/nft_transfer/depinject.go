package nfttransfer

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	"github.com/cosmos/cosmos-sdk/codec"
	store "github.com/cosmos/cosmos-sdk/store/types"

	modulev1 "github.com/bianjieai/tibc-go/api/tibc/apps/nft_transfer/module/v1"
	"github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/keeper"
	"github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/types"
)

// App Wiring Setup

func init() {
	appmodule.Register(&modulev1.Module{},
		appmodule.Provide(ProvideModule),
	)
}

var _ appmodule.AppModule = AppModule{}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

// Inputs define the module inputs for the depinject.
type Inputs struct {
	depinject.In

	Config *modulev1.Module
	Cdc    codec.Codec
	Key    *store.KVStoreKey

	ak         types.AccountKeeper
	mk         types.NftKeeper
	pk         types.PacketKeeper
	ck         types.ClientKeeper
}

// Outputs define the module outputs for the depinject.
type Outputs struct {
	depinject.Out

	NftTransferKeeper keeper.Keeper
	Module           appmodule.AppModule
}


// ProvideModule creates and returns the module outputs for the depinject.
//
// It takes Inputs as the parameter, which includes the configuration, codec, key, account keeper, NFT keeper, packet keeper, and client keeper.
// It returns Outputs containing the NFT transfer keeper and the app module.
func ProvideModule(in Inputs) Outputs {
	keeper := keeper.NewKeeper(
		in.Cdc,
		in.Key,
		in.ak,
		in.mk,
		in.pk,
		in.ck,
	)
	return Outputs{
		NftTransferKeeper: keeper,
		Module:           NewAppModule(keeper),
	}
}
