package mttransfer

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	"github.com/cosmos/cosmos-sdk/codec"
	store "github.com/cosmos/cosmos-sdk/store/types"

	modulev1 "github.com/bianjieai/tibc-go/api/tibc/apps/mt_transfer/module/v1"
	"github.com/bianjieai/tibc-go/modules/tibc/apps/mt_transfer/keeper"
	"github.com/bianjieai/tibc-go/modules/tibc/apps/mt_transfer/types"
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
	mk         types.MtKeeper
	pk         types.PacketKeeper
	ck         types.ClientKeeper
}

// Outputs define the module outputs for the depinject.
type Outputs struct {
	depinject.Out

	MtTransferKeeper keeper.Keeper
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
		in.ak,
		in.mk,
		in.pk,
		in.ck,
	)
	return Outputs{
		MtTransferKeeper: keeper,
		Module:           NewAppModule(keeper),
	}
}