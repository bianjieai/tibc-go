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

var _ clienttypes.MsgServer = Keeper{}
var _ packettypes.MsgServer = Keeper{}

// UpdateClient defines a rpc handler method for MsgUpdateClient.
func (k Keeper) UpdateClient(goCtx context.Context, msg *clienttypes.MsgUpdateClient) (*clienttypes.MsgUpdateClientResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	header, err := clienttypes.UnpackHeader(msg.Header)
	if err != nil {
		return nil, err
	}

	// Verify that the account has permission to update the client
	if !k.ClientKeeper.AuthRelayer(ctx, msg.ChainName, msg.Signer) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "relayer: %s", msg.Signer)
	}

	if err = k.ClientKeeper.UpdateClient(ctx, msg.ChainName, header); err != nil {
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

// RecvPacket defines a rpc handler method for MsgRecvPacket.
func (k Keeper) RecvPacket(goCtx context.Context, msg *packettypes.MsgRecvPacket) (*packettypes.MsgRecvPacketResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Perform TAO verification
	if err := k.PacketKeeper.RecvPacket(ctx, msg.Packet, msg.ProofCommitment, msg.ProofHeight); err != nil {
		switch err {
		case sdkerrors.ErrUnauthorized:
			if err2 := k.PacketKeeper.WriteAcknowledgement(ctx, msg.Packet, packettypes.NewErrorAcknowledgement(err.Error()).GetBytes()); err2 != nil {
				return nil, err2
			}
			return &packettypes.MsgRecvPacketResponse{}, nil
		default:
			return nil, sdkerrors.Wrap(err, "receive packet verification failed")
		}
	}

	if msg.Packet.GetDestChain() == k.ClientKeeper.GetChainName(ctx) {
		// Retrieve callbacks from router
		cbs, ok := k.RoutingKeeper.Router.GetRoute(routingtypes.Port(msg.Packet.Port))
		if !ok {
			return nil, sdkerrors.Wrapf(routingtypes.ErrInvalidRoute, "route not found to module: %s", msg.Packet.Port)
		}

		// Perform application logic callback
		_, ack, err := cbs.OnRecvPacket(ctx, msg.Packet)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "receive packet callback failed")
		}

		// Set packet acknowledgement only if the acknowledgement is not nil.
		// NOTE: TIBC applications modules may call the WriteAcknowledgement asynchronously if the
		// acknowledgement is nil.
		if ack != nil {
			if err := k.PacketKeeper.WriteAcknowledgement(ctx, msg.Packet, ack); err != nil {
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
func (k Keeper) Acknowledgement(goCtx context.Context, msg *packettypes.MsgAcknowledgement) (*packettypes.MsgAcknowledgementResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Retrieve callbacks from router
	cbs, ok := k.RoutingKeeper.Router.GetRoute(routingtypes.Port(msg.Packet.Port))
	if !ok {
		return nil, sdkerrors.Wrapf(routingtypes.ErrInvalidRoute, "route not found to module: %s", msg.Packet.Port)
	}

	// Perform TAO verification
	if err := k.PacketKeeper.AcknowledgePacket(ctx, msg.Packet, msg.Acknowledgement, msg.ProofAcked, msg.ProofHeight); err != nil {
		return nil, sdkerrors.Wrap(err, "acknowledge packet verification failed")
	}

	// Perform application logic callback
	_, err := cbs.OnAcknowledgementPacket(ctx, msg.Packet, msg.Acknowledgement)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "acknowledge packet callback failed")
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
func (k Keeper) CleanPacket(goCtx context.Context, msg *packettypes.MsgCleanPacket) (*packettypes.MsgCleanPacketResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.PacketKeeper.CleanPacket(ctx, msg.CleanPacket); err != nil {
		return nil, sdkerrors.Wrap(err, "send clean packet failed")
	}

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{"tx", "msg", "tibc", packettypes.EventTypeSendCleanPacket},
			1,
			[]metrics.Label{
				telemetry.NewLabel(coretypes.LabelSourceChain, msg.CleanPacket.SourceChain),
				telemetry.NewLabel(coretypes.LabelDestinationChain, msg.CleanPacket.DestinationChain),
				telemetry.NewLabel(coretypes.LabelRelayChain, msg.CleanPacket.RelayChain),
			},
		)
	}()

	return &packettypes.MsgCleanPacketResponse{}, nil
}

// RecvCleanPacket defines a rpc handler method for MsgAcknowledgement.
func (k Keeper) RecvCleanPacket(goCtx context.Context, msg *packettypes.MsgRecvCleanPacket) (*packettypes.MsgRecvCleanPacketResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.PacketKeeper.RecvCleanPacket(ctx, msg.CleanPacket, msg.ProofCommitment, msg.ProofHeight); err != nil {
		return nil, sdkerrors.Wrap(err, "receive clean packet failed")
	}

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{"tx", "msg", "tibc", packettypes.EventTypeRecvCleanPacket},
			1,
			[]metrics.Label{
				telemetry.NewLabel(coretypes.LabelSourceChain, msg.CleanPacket.SourceChain),
				telemetry.NewLabel(coretypes.LabelDestinationChain, msg.CleanPacket.DestinationChain),
				telemetry.NewLabel(coretypes.LabelRelayChain, msg.CleanPacket.RelayChain),
			},
		)
	}()

	return &packettypes.MsgRecvCleanPacketResponse{}, nil
}
