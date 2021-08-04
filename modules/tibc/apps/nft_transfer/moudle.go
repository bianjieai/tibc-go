package nft_transfer

import (
	"encoding/json"
	"fmt"
	"github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/keeper"
	 nftTransferTypes "github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/types"
	"github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
	porttypes "github.com/bianjieai/tibc-go/modules/tibc/core/05-port/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ porttypes.IBCModule   = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic is the TIBC nft Transfer AppModuleBasic
type AppModuleBasic struct{}

func (a AppModuleBasic) Name() string {
	panic("implement me")
}

func (a AppModuleBasic) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {
	panic("implement me")
}

func (a AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	panic("implement me")
}

func (a AppModuleBasic) DefaultGenesis(marshaler codec.JSONMarshaler) json.RawMessage {
	panic("implement me")
}

func (a AppModuleBasic) ValidateGenesis(marshaler codec.JSONMarshaler, config client.TxEncodingConfig, message json.RawMessage) error {
	panic("implement me")
}

func (a AppModuleBasic) RegisterRESTRoutes(context client.Context, m *interface{}) {
	panic("implement me")
}

func (a AppModuleBasic) RegisterGRPCGatewayRoutes(context client.Context, r *interface{}) {
	panic("implement me")
}

func (a AppModuleBasic) GetTxCmd() *interface{} {
	panic("implement me")
}

func (a AppModuleBasic) GetQueryCmd() *interface{} {
	panic("implement me")
}

// AppModule represents the AppModule for this module
type AppModule struct {
	AppModuleBasic
	keeper keeper.Keeper
}

func (a AppModule) Name() string {
	panic("implement me")
}

func (a AppModule) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {
	panic("implement me")
}

func (a AppModule) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	panic("implement me")
}

func (a AppModule) DefaultGenesis(marshaler codec.JSONMarshaler) json.RawMessage {
	panic("implement me")
}

func (a AppModule) ValidateGenesis(marshaler codec.JSONMarshaler, config client.TxEncodingConfig, message json.RawMessage) error {
	panic("implement me")
}

func (a AppModule) RegisterRESTRoutes(context client.Context, m *interface{}) {
	panic("implement me")
}

func (a AppModule) RegisterGRPCGatewayRoutes(context client.Context, r *interface{}) {
	panic("implement me")
}

func (a AppModule) GetTxCmd() *interface{} {
	panic("implement me")
}

func (a AppModule) GetQueryCmd() *interface{} {
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
	panic("implement me")
}

func (a AppModule) LegacyQuerierHandler(amino *codec.LegacyAmino) sdk.Querier {
	panic("implement me")
}

func (a AppModule) RegisterServices(configurator module.Configurator) {
	panic("implement me")
}

func (a AppModule) ConsensusVersion() uint64 {
	panic("implement me")
}

func (a AppModule) BeginBlock(context sdk.Context, block abci.RequestBeginBlock) {
	panic("implement me")
}

func (a AppModule) EndBlock(context sdk.Context, block abci.RequestEndBlock) []abci.ValidatorUpdate {
	panic("implement me")
}

func (a AppModule) OnRecvPacket(ctx sdk.Context, packet types.Packet) (*sdk.Result, []byte, error) {

	// create ack
	ack := channeltypes.NewResultAcknowledgement([]byte{byte(1)})

	var data nftTransferTypes.NonFungibleTokenPacketData

	// unmarshal packet.data
	if err := nftTransferTypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		ack = channeltypes.NewErrorAcknowledgement(fmt.Sprintf("cannot unmarshal ICS-20 transfer packet data: %s", err.Error()))
	}

	// only attempt the application logic if the packet data
	// was successfully decoded
	if ack.Success() {
		err := a.keeper.OnRecvPacket(ctx, packet, data)
		if err != nil {
			ack = channeltypes.NewErrorAcknowledgement(err.Error())
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			nftTransferTypes.EventTypePacket,
			sdk.NewAttribute(sdk.AttributeKeyModule, nftTransferTypes.ModuleName),
			sdk.NewAttribute(nftTransferTypes.AttributeKeyReceiver, data.Receiver),
			sdk.NewAttribute(types.AttributeKeyDenom, data.Denom),
			sdk.NewAttribute(types.AttributeKeyAmount, fmt.Sprintf("%d", data.Amount)),
			sdk.NewAttribute(types.AttributeKeyAckSuccess, fmt.Sprintf("%t", ack.Success())),
		),
	)

	// NOTE: acknowledgement will be written synchronously during IBC handler execution.
	return ack
}

func (a AppModule) OnAcknowledgementPacket(ctx sdk.Context, packet types.Packet, acknowledgement []byte) (*sdk.Result, error) {

}

func (a AppModule) OnTimeoutPacket(ctx sdk.Context, packet types.Packet) (*sdk.Result, error) {
	panic("implement me")
}

