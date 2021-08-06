package nft_transfer

import (
	"encoding/json"
	"fmt"
	"github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/client/cli"
	"github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/keeper"
	"github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/types"
	packettypes "github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
	porttypes "github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ porttypes.TIBCModule   = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic is the TIBC nft Transfer AppModuleBasic
type AppModuleBasic struct{}

func (a AppModuleBasic) RegisterRESTRoutes(context client.Context, m *interface{}) {
}

func (a AppModuleBasic) RegisterGRPCGatewayRoutes(context client.Context, r *interface{}) {
}

func (a AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}

// GetQueryCmd returns no root query command for the ibc module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	//todo
	//return cli.GetQueryCmd()
	return nil
}

func (a AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec implements AppModuleBasic interface
func (a AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers module concrete types into protobuf Any.
func (a AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

func (a AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	//return cdc.MustMarshalJSON(types.DefaultGenesisState())
	//todo
	return nil
}

func (a AppModuleBasic) ValidateGenesis(marshaler codec.JSONMarshaler, config client.TxEncodingConfig, message json.RawMessage) error {
	panic("implement me")
}

// AppModule represents the AppModule for this module
type AppModule struct {
	AppModuleBasic
	keeper keeper.Keeper
}


func (a AppModule) RegisterRESTRoutes(context client.Context, m *interface{}) {
}

func (a AppModule) RegisterGRPCGatewayRoutes(context client.Context, r *interface{}) {
	panic("implement me")
}

func (a AppModule) InitGenesis(context sdk.Context, marshaler codec.JSONMarshaler, message json.RawMessage) []abci.ValidatorUpdate {
	panic("implement me")
}

func (a AppModule) ExportGenesis(context sdk.Context, marshaler codec.JSONMarshaler) json.RawMessage {
	panic("implement me")
}

func (a AppModule) RegisterInvariants(registry sdk.InvariantRegistry) {
	panic("implement me")
}

func (a AppModule) Route() sdk.Route {
	panic("implement me")
}

func (a AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

func (a AppModule) LegacyQuerierHandler(amino *codec.LegacyAmino) sdk.Querier {
	panic("implement me")
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	panic("implement me")
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (am AppModule) ConsensusVersion() uint64 { return 1 }

// BeginBlock implements the AppModule interface
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
}

// EndBlock implements the AppModule interface
func (am AppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

func (a AppModule) OnRecvPacket(ctx sdk.Context, packet packettypes.Packet) (*sdk.Result, []byte, error) {

	var data types.NonFungibleTokenPacketData
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal TICS-30 nft-transfer packet data: %s", err.Error())
	}

	acknowledgement := packettypes.NewResultAcknowledgement([]byte{byte(1)})

	err := a.keeper.OnRecvPacket(ctx, packet, data)
	if err != nil {
		acknowledgement = packettypes.NewErrorAcknowledgement(err.Error())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePacket,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyReceiver, data.Receiver),
			sdk.NewAttribute(types.AttributeKeyClass, data.Class),
			sdk.NewAttribute(types.AttributeKeyId, data.Class),
			sdk.NewAttribute(types.AttributeKeyUri, data.Uri),
			sdk.NewAttribute(types.AttributeKeyAckSuccess, fmt.Sprintf("%t", err == nil)),
		),
	)

	// NOTE: acknowledgement will be written synchronously during IBC handler execution.
	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, acknowledgement.GetBytes(), nil
}

func (a AppModule) OnAcknowledgementPacket(ctx sdk.Context, packet packettypes.Packet, acknowledgement []byte) (*sdk.Result, error) {
	var ack packettypes.Acknowledgement
	if err := types.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet acknowledgement: %v", err)
	}
	var data types.NonFungibleTokenPacketData
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet data: %s", err.Error())
	}

	if err := a.keeper.OnAcknowledgementPacket(ctx, data, ack); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePacket,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyReceiver, data.Receiver),
			sdk.NewAttribute(types.AttributeKeyReceiver, data.Receiver),
			sdk.NewAttribute(types.AttributeKeyClass, data.Class),
			sdk.NewAttribute(types.AttributeKeyId, data.Class),
			sdk.NewAttribute(types.AttributeKeyUri, data.Uri),
			sdk.NewAttribute(types.AttributeKeyAck, fmt.Sprintf("%v", ack)),
		),
	)

	switch resp := ack.Response.(type) {
	case *packettypes.Acknowledgement_Result:
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypePacket,
				sdk.NewAttribute(types.AttributeKeyAckSuccess, string(resp.Result)),
			),
		)
	case *packettypes.Acknowledgement_Error:
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypePacket,
				sdk.NewAttribute(types.AttributeKeyAckError, resp.Error),
			),
		)
	}

	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

