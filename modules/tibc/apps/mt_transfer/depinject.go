package mttransfer

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	"github.com/cosmos/cosmos-sdk/codec"
	store "github.com/cosmos/cosmos-sdk/store/types"

	modulev1 "github.com/bianjieai/tibc-go/api/tibc/apps/mt_transfer/module/v1"
	"github.com/bianjieai/tibc-go/modules/tibc/apps/mt_transfer/keeper"
	"github.com/bianjieai/tibc-go/modules/tibc/apps/mt_transfer/types"
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
	MtKeeper      types.MtKeeper
	PacketKeeper  types.PacketKeeper
	ClientKeeper  types.ClientKeeper
}

// Outputs define the module outputs for the depinject.
type Outputs struct {
	depinject.Out

	MtTransferKeeper keeper.Keeper
	Route            routingtypes.Route
	Module           appmodule.AppModule
}

// ProvideModule creates and returns the farm module with the specified inputs.
//
// It takes Inputs as the parameter, which includes the configuration, codec, key, account keeper, bank keeper, governance keeper, coinswap keeper, and legacy subspace.
// It returns Outputs containing the farm keeper and the app module.
func ProvideModule(in Inputs) Outputs {
	keeper := keeper.NewKeeper(
		in.Cdc,
		in.Key,
		in.AccountKeeper,
		in.MtKeeper,
		in.PacketKeeper,
		in.ClientKeeper,
	)
	m := NewAppModule(keeper)
	route := routingtypes.Route{
		Port:   string(routingtypes.MT),
		Module: m,
	}
	return Outputs{
		MtTransferKeeper: keeper,
		Route:            route,
		Module:           m,
	}
}
