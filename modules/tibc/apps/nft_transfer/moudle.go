package nfttransfer

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/client/cli"
	"github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/keeper"
	"github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/types"
	packettypes "github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
	porttypes "github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
	"github.com/bianjieai/tibc-go/modules/tibc/core/simulation"
)

var (
	_ module.AppModule      = AppModule{}
	_ porttypes.TIBCModule  = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic is the TIBC nft Transfer AppModuleBasic
type AppModuleBasic struct{}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the tibc-nft-transfer module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	_ = types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))
}

// GetTxCmd returns the transaction commands for the tibc-nft-transfer module.
func (a AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}

// GetQueryCmd returns no root query command for the tibc-nft-transfer module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// Name returns the name of the module.
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

// AppModule represents the AppModule for this module
type AppModule struct {
	AppModuleBasic
	keeper keeper.Keeper
}

// GenerateGenesisState creates a randomized GenState of the tibc-nft-transfer module.
func (am AppModule) GenerateGenesisState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
}

// NewAppModule creates a new tibc-nft-transfer module
func NewAppModule(k keeper.Keeper) AppModule {
	return AppModule{
		keeper: k,
	}
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), am.keeper)
	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (am AppModule) ConsensusVersion() uint64 { return 1 }

// OnRecvPacket implements the TIBCModule interface.
// It processes the incoming packet from counterparty chain using the Keeper.OnRecvPacket method.
// It emits an event with the packet data and acknowledgement status.
// It returns the acknowledgement bytes and error.
func (am AppModule) OnRecvPacket(
	ctx sdk.Context,
	packet packettypes.Packet,
) (*sdk.Result, []byte, error) {

	var data types.NonFungibleTokenPacketData
	if err := data.Unmarshal(packet.GetData()); err != nil {
		return nil, nil, errorsmod.Wrapf(
			sdkerrors.ErrUnknownRequest,
			"cannot unmarshal TICS-30 nft-transfer packet data: %s",
			err.Error(),
		)
	}

	acknowledgement := packettypes.NewResultAcknowledgement([]byte{byte(1)})

	err := am.keeper.OnRecvPacket(ctx, packet, data)
	if err != nil {
		acknowledgement = packettypes.NewErrorAcknowledgement(err.Error())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePacket,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyReceiver, data.Receiver),
			sdk.NewAttribute(types.AttributeKeyClass, data.Class),
			sdk.NewAttribute(types.AttributeKeyId, data.Id),
			sdk.NewAttribute(types.AttributeKeyUri, data.Uri),
			sdk.NewAttribute(types.AttributeKeyAckSuccess, fmt.Sprintf("%t", err == nil)),
		),
	)

	// NOTE: acknowledgement will be written synchronously during TIBC handler execution.
	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, acknowledgement.GetBytes(), nil
}

// OnAcknowledgementPacket implements the TIBCModule interface.
// It processes the acknowledgement that the counterparty chain sends back in response to a packet that was sent by this chain.
// It emits an event with the packet data and acknowledgement status.
// It returns the acknowledgement bytes and error.
func (am AppModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet packettypes.Packet,
	acknowledgement []byte,
) (*sdk.Result, error) {
	var ack packettypes.Acknowledgement
	if err := ack.Unmarshal(acknowledgement); err != nil {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrUnknownRequest,
			"cannot unmarshal TICS-30 transfer packet acknowledgement: %v",
			err,
		)
	}
	var data types.NonFungibleTokenPacketData
	if err := data.Unmarshal(packet.GetData()); err != nil {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrUnknownRequest,
			"cannot unmarshal TICS-30 transfer packet data: %s",
			err.Error(),
		)
	}

	if err := am.keeper.OnAcknowledgementPacket(ctx, data, ack); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePacket,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyReceiver, data.Receiver),
			sdk.NewAttribute(types.AttributeKeyReceiver, data.Receiver),
			sdk.NewAttribute(types.AttributeKeyClass, data.Class),
			sdk.NewAttribute(types.AttributeKeyId, data.Id),
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
