package keeper

import (
	"context"

	"github.com/armon/go-metrics"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	packettypes "github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
	routingtypes "github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
	coretypes "github.com/bianjieai/tibc-go/modules/tibc/core/types"
)

var (
	_ clienttypes.MsgServer  = msgServer{}
	_ packettypes.MsgServer  = msgServer{}
	_ routingtypes.MsgServer = msgServer{}
)

type msgServer struct {
	k Keeper
}

// NewMsgServerImpl returns an implementation of the mint MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) coretypes.MsgServer {
	return &msgServer{k: keeper}
}

// CreateClient defines a rpc handler method for MsgCreateClient.
func (m msgServer) CreateClient(
	goCtx context.Context,
	msg *clienttypes.MsgCreateClient,
) (*clienttypes.MsgCreateClientResponse, error) {
	if m.k.authority != msg.Authority {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrUnauthorized,
			"invalid authority; expected %s, got %s",
			m.k.authority,
			msg.Authority,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	_, has := m.k.ClientKeeper.GetClientState(ctx, msg.ChainName)
	if has {
		return nil, sdkerrors.Wrapf(
			clienttypes.ErrClientExists,
			"chain-name: %s",
			msg.ChainName,
		)
	}

	clientState, err := clienttypes.UnpackClientState(msg.ClientState)
	if err != nil {
		return nil, err
	}

	consensusState, err := clienttypes.UnpackConsensusState(msg.ConsensusState)
	if err != nil {
		return nil, err
	}

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{"tibc", "client", "create"},
			1,
			[]metrics.Label{
				telemetry.NewLabel(clienttypes.LabelClientType, clientState.ClientType()),
				telemetry.NewLabel(clienttypes.LabelChainName, msg.ChainName),
			},
		)
	}()

	// emitting events in the keeper for proposal updates to clients
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			clienttypes.EventTypeCreateClientProposal,
			sdk.NewAttribute(clienttypes.AttributeKeyChainName, msg.ChainName),
			sdk.NewAttribute(clienttypes.AttributeKeyClientType, clientState.ClientType()),
			sdk.NewAttribute(
				clienttypes.AttributeKeyConsensusHeight,
				clientState.GetLatestHeight().String(),
			),
		),
	)
	if err := m.k.ClientKeeper.CreateClient(ctx, msg.ChainName, clientState, consensusState); err != nil {
		return nil, err
	}
	return &clienttypes.MsgCreateClientResponse{}, nil
}

// UpdateClient defines a rpc handler method for MsgUpdateClient.
func (m msgServer) UpdateClient(
	goCtx context.Context,
	msg *clienttypes.MsgUpdateClient,
) (*clienttypes.MsgUpdateClientResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	header, err := clienttypes.UnpackHeader(msg.Header)
	if err != nil {
		return nil, err
	}

	// Verify that the account has permission to update the client
	if !m.k.ClientKeeper.AuthRelayer(ctx, msg.ChainName, msg.Signer) {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrUnauthorized,
			"relayer: %s",
			msg.Signer,
		)
	}

	if err = m.k.ClientKeeper.UpdateClient(ctx, msg.ChainName, header); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, clienttypes.AttributeValueCategory),
		),
	)

	return &clienttypes.MsgUpdateClientResponse{}, nil
}

// UpgradeClient defines a rpc handler method for MsgUpgradeClient.
func (m msgServer) UpgradeClient(
	goCtx context.Context,
	msg *clienttypes.MsgUpgradeClient,
) (*clienttypes.MsgUpgradeClientResponse, error) {
	if m.k.authority != msg.Authority {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrUnauthorized,
			"invalid authority; expected %s, got %s",
			m.k.authority,
			msg.Authority,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	clientState, err := clienttypes.UnpackClientState(msg.ClientState)
	if err != nil {
		return nil, err
	}

	consensusState, err := clienttypes.UnpackConsensusState(msg.ConsensusState)
	if err != nil {
		return nil, err
	}

	if err := m.k.ClientKeeper.UpgradeClient(ctx, msg.ChainName, clientState, consensusState); err != nil {
		return nil, err
	}
	return &clienttypes.MsgUpgradeClientResponse{}, nil
}

// RegisterRelayer defines a rpc handler method for MsgRegisterRelayer.
func (m msgServer) RegisterRelayer(
	goCtx context.Context,
	msg *clienttypes.MsgRegisterRelayer,
) (*clienttypes.MsgRegisterRelayerResponse, error) {
	if m.k.authority != msg.Authority {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrUnauthorized,
			"invalid authority; expected %s, got %s",
			m.k.authority,
			msg.Authority,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	m.k.ClientKeeper.RegisterRelayers(ctx, msg.ChainName, msg.Relayers)
	return &clienttypes.MsgRegisterRelayerResponse{}, nil
}

// SetRoutingRules defines a rpc handler method for MsgSetRoutingRules.
func (m msgServer) SetRoutingRules(
	goCtx context.Context,
	msg *routingtypes.MsgSetRoutingRules,
) (*routingtypes.MsgSetRoutingRulesResponse, error) {
	if m.k.authority != msg.Authority {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrUnauthorized,
			"invalid authority; expected %s, got %s",
			m.k.authority,
			msg.Authority,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.k.RoutingKeeper.SetRoutingRules(ctx, msg.Rules); err != nil {
		return nil, err
	}
	return &routingtypes.MsgSetRoutingRulesResponse{}, nil
}

// RecvPacket defines a rpc handler method for MsgRecvPacket.
func (m msgServer) RecvPacket(
	goCtx context.Context,
	msg *packettypes.MsgRecvPacket,
) (*packettypes.MsgRecvPacketResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Perform TAO verification
	if err := m.k.PacketKeeper.RecvPacket(ctx, msg.Packet, msg.ProofCommitment, msg.ProofHeight); err != nil {
		switch err {
		case sdkerrors.ErrUnauthorized:
			if err2 := m.k.PacketKeeper.WriteAcknowledgement(ctx, msg.Packet, packettypes.NewErrorAcknowledgement(err.Error()).GetBytes()); err2 != nil {
				return nil, err2
			}
			return &packettypes.MsgRecvPacketResponse{}, nil
		default:
			return nil, sdkerrors.Wrap(
				err,
				"receive packet verification failed",
			)
		}
	}

	if msg.Packet.GetDestChain() == m.k.ClientKeeper.GetChainName(ctx) {
		// Retrieve callbacks from router
		cbs, ok := m.k.RoutingKeeper.Router.GetRoute(routingtypes.Port(msg.Packet.Port))
		if !ok {
			return nil, sdkerrors.Wrapf(
				routingtypes.ErrInvalidRoute,
				"route not found to module: %s",
				msg.Packet.Port,
			)
		}

		// Perform application logic callback
		_, ack, err := cbs.OnRecvPacket(ctx, msg.Packet)
		if err != nil {
			return nil, sdkerrors.Wrap(
				err,
				"receive packet callback failed",
			)
		}

		// Set packet acknowledgement only if the acknowledgement is not nil.
		// NOTE: TIBC applications modules may call the WriteAcknowledgement asynchronously if the
		// acknowledgement is nil.
		if ack != nil {
			if err := m.k.PacketKeeper.WriteAcknowledgement(ctx, msg.Packet, ack); err != nil {
				return nil, err
			}
		}
	}

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{"tx", "msg", "tibc", packettypes.EventTypeRecvPacket},
			1,
			[]metrics.Label{
				telemetry.NewLabel(coretypes.LabelPort, msg.Packet.Port),
				telemetry.NewLabel(coretypes.LabelSourceChain, msg.Packet.SourceChain),
				telemetry.NewLabel(coretypes.LabelDestinationChain, msg.Packet.DestinationChain),
				telemetry.NewLabel(coretypes.LabelRelayChain, msg.Packet.RelayChain),
			},
		)
	}()

	return &packettypes.MsgRecvPacketResponse{}, nil
}

// Acknowledgement defines a rpc handler method for MsgAcknowledgement.
func (m msgServer) Acknowledgement(
	goCtx context.Context,
	msg *packettypes.MsgAcknowledgement,
) (*packettypes.MsgAcknowledgementResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Retrieve callbacks from router
	cbs, ok := m.k.RoutingKeeper.Router.GetRoute(routingtypes.Port(msg.Packet.Port))
	if !ok {
		return nil, sdkerrors.Wrapf(
			routingtypes.ErrInvalidRoute,
			"route not found to module: %s",
			msg.Packet.Port,
		)
	}

	// Perform TAO verification
	if err := m.k.PacketKeeper.AcknowledgePacket(ctx, msg.Packet, msg.Acknowledgement, msg.ProofAcked, msg.ProofHeight); err != nil {
		return nil, sdkerrors.Wrap(
			err,
			"acknowledge packet verification failed",
		)
	}

	// Perform application logic callback
	_, err := cbs.OnAcknowledgementPacket(ctx, msg.Packet, msg.Acknowledgement)
	if err != nil {
		return nil, sdkerrors.Wrap(
			err,
			"acknowledge packet callback failed",
		)
	}

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{"tx", "msg", "tibc", packettypes.EventTypeAcknowledgePacket},
			1,
			[]metrics.Label{
				telemetry.NewLabel(coretypes.LabelPort, msg.Packet.Port),
				telemetry.NewLabel(coretypes.LabelSourceChain, msg.Packet.SourceChain),
				telemetry.NewLabel(coretypes.LabelDestinationChain, msg.Packet.DestinationChain),
				telemetry.NewLabel(coretypes.LabelRelayChain, msg.Packet.RelayChain),
			},
		)
	}()

	return &packettypes.MsgAcknowledgementResponse{}, nil
}

// Acknowledgement defines a rpc handler method for MsgAcknowledgement.
func (m msgServer) CleanPacket(
	goCtx context.Context,
	msg *packettypes.MsgCleanPacket,
) (*packettypes.MsgCleanPacketResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.k.PacketKeeper.CleanPacket(ctx, msg.CleanPacket); err != nil {
		return nil, sdkerrors.Wrap(
			err,
			"send clean packet failed",
		)
	}

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{"tx", "msg", "tibc", packettypes.EventTypeSendCleanPacket},
			1,
			[]metrics.Label{
				telemetry.NewLabel(coretypes.LabelSourceChain, msg.CleanPacket.SourceChain),
				telemetry.NewLabel(
					coretypes.LabelDestinationChain,
					msg.CleanPacket.DestinationChain,
				),
				telemetry.NewLabel(coretypes.LabelRelayChain, msg.CleanPacket.RelayChain),
			},
		)
	}()

	return &packettypes.MsgCleanPacketResponse{}, nil
}

// RecvCleanPacket defines a rpc handler method for MsgAcknowledgement.
func (m msgServer) RecvCleanPacket(
	goCtx context.Context,
	msg *packettypes.MsgRecvCleanPacket,
) (*packettypes.MsgRecvCleanPacketResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.k.PacketKeeper.RecvCleanPacket(ctx, msg.CleanPacket, msg.ProofCommitment, msg.ProofHeight); err != nil {
		return nil, sdkerrors.Wrap(
			err,
			"receive clean packet failed",
		)
	}

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{"tx", "msg", "tibc", packettypes.EventTypeRecvCleanPacket},
			1,
			[]metrics.Label{
				telemetry.NewLabel(coretypes.LabelSourceChain, msg.CleanPacket.SourceChain),
				telemetry.NewLabel(
					coretypes.LabelDestinationChain,
					msg.CleanPacket.DestinationChain,
				),
				telemetry.NewLabel(coretypes.LabelRelayChain, msg.CleanPacket.RelayChain),
			},
		)
	}()

	return &packettypes.MsgRecvCleanPacketResponse{}, nil
}
