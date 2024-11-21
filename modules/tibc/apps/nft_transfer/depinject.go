package nfttransfer

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	store "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"

	modulev1 "github.com/bianjieai/tibc-go/api/tibc/apps/nft_transfer/module/v1"
	"github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/keeper"
	"github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/types"
	routingtypes "github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
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

	AccountKeeper types.AccountKeeper
	NftKeeper     types.NftKeeper
	PacketKeeper  types.PacketKeeper
	ClientKeeper  types.ClientKeeper
}

// Outputs define the module outputs for the depinject.
type Outputs struct {
	depinject.Out

	NftTransferKeeper keeper.Keeper
	Route             routingtypes.Route
	Module            appmodule.AppModule
}

// ProvideModule creates and returns the module outputs for the depinject.
//
// It takes Inputs as the parameter, which includes the configuration, codec, key, account keeper, NFT keeper, packet keeper, and client keeper.
// It returns Outputs containing the NFT transfer keeper and the app module.
func ProvideModule(in Inputs) Outputs {
	keeper := keeper.NewKeeper(
		in.Cdc,
		in.Key,
		in.AccountKeeper,
		in.NftKeeper,
		in.PacketKeeper,
		in.ClientKeeper,
	)
	m := NewAppModule(keeper)
	route := routingtypes.Route{
		Port:   string(routingtypes.NFT),
		Module: m,
	}
	return Outputs{
		NftTransferKeeper: keeper,
		Route:             route,
		Module:            m,
	}
}
