package keeper

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	"github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

// SendPacket is called by a module to send an TIBC packet on a port owned
// by the calling module to the corresponding module on the counterparty chain.
func (k Keeper) SendPacket(ctx sdk.Context, packet exported.PacketI) error {
	if err := packet.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "packet failed basic validation")
	}
	if packet.GetSourceChain() != k.clientKeeper.GetChainName(ctx) {
		return sdkerrors.Wrap(types.ErrInvalidPacket, "source chain of packet is not this chain")
	}

	targetChain := packet.GetDestChain()
	if len(packet.GetRelayChain()) > 0 {
		targetChain = packet.GetRelayChain()
	}

	_, found := k.clientKeeper.GetClientState(ctx, targetChain)
	if !found {
		return clienttypes.ErrClientNotFound
	}

	nextSequenceSend := k.GetNextSequenceSend(ctx, packet.GetSourceChain(), packet.GetDestChain())

	if packet.GetSequence() != nextSequenceSend {
		return sdkerrors.Wrapf(
			types.ErrInvalidPacket,
			"packet sequence ≠ next send sequence (%d ≠ %d)", packet.GetSequence(), nextSequenceSend,
		)
	}

	commitment := types.CommitPacket(packet)

	nextSequenceSend++
	k.SetNextSequenceSend(ctx, packet.GetSourceChain(), packet.GetDestChain(), nextSequenceSend)
	k.SetPacketCommitment(ctx, packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence(), commitment)

	// Emit Event with Packet data along with other packet information for relayer to pick up
	// and relay to other chain
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSendPacket,
			sdk.NewAttribute(types.AttributeKeyData, string(packet.GetData())),
			sdk.NewAttribute(types.AttributeKeySequence, fmt.Sprintf("%d", packet.GetSequence())),
			sdk.NewAttribute(types.AttributeKeyPort, packet.GetPort()),
			sdk.NewAttribute(types.AttributeKeySrcChain, packet.GetSourceChain()),
			sdk.NewAttribute(types.AttributeKeyDstChain, packet.GetDestChain()),
			sdk.NewAttribute(types.AttributeKeyRelayChain, packet.GetRelayChain()),
			// we only support 1-hop packets now, and that is the most important hop for a relayer
			// (is it going to a chain I am connected to)
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	k.Logger(ctx).Info("packet sent", "packet", fmt.Sprintf("%v", packet))
	return nil
}

// RecvPacket is called by a module to receive & process an TIBC packet
// sent on the corresponding port on the counterparty chain.
func (k Keeper) RecvPacket(
	ctx sdk.Context,
	packet exported.PacketI,
	proof []byte,
	proofHeight exported.Height,
) error {
	if err := k.ValidatePacket(ctx, packet); err != nil {
		return sdkerrors.Wrap(err, "packet failed basic validation")
	}
	// check if the packet receipt has been received already for unordered channels
	_, found := k.GetPacketReceipt(ctx, packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence())
	if found {
		return sdkerrors.Wrapf(
			types.ErrInvalidPacket,
			"packet sequence (%d) already has been received", packet.GetSequence(),
		)
	}
	chainName := k.clientKeeper.GetChainName(ctx)
	fromChain := packet.GetSourceChain()
	if packet.GetDestChain() == chainName && len(packet.GetRelayChain()) > 0 {
		fromChain = packet.GetRelayChain()
	}

	targetClient, found := k.clientKeeper.GetClientState(ctx, fromChain)
	if !found {
		return sdkerrors.Wrap(clienttypes.ErrClientNotFound, fromChain)
	}

	commitment := types.CommitPacket(packet)
	// verify that the counterparty did commit to sending this packet
	if err := targetClient.VerifyPacketCommitment(ctx,
		k.clientKeeper.ClientStore(ctx, fromChain), k.cdc, proofHeight,
		proof, packet.GetSourceChain(), packet.GetDestChain(),
		packet.GetSequence(), commitment,
	); err != nil {
		return sdkerrors.Wrapf(err, "failed packet commitment verification for client (%s)", fromChain)
	}

	// All verification complete, update state
	// For unordered channels we must set the receipt so it can be verified on the other side.
	// This receipt does not contain any data, since the packet has not yet been processed,
	// it's just a single store key set to an empty string to indicate that the packet has been received
	k.SetPacketReceipt(ctx, packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence())

	// log that a packet has been received & executed
	k.Logger(ctx).Info("packet received", "packet", fmt.Sprintf("%v", packet))

	// emit an event that the relayer can query for
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRecvPacket,
			sdk.NewAttribute(types.AttributeKeyData, string(packet.GetData())),
			sdk.NewAttribute(types.AttributeKeySequence, fmt.Sprintf("%d", packet.GetSequence())),
			sdk.NewAttribute(types.AttributeKeyPort, packet.GetPort()),
			sdk.NewAttribute(types.AttributeKeySrcChain, packet.GetSourceChain()),
			sdk.NewAttribute(types.AttributeKeyDstChain, packet.GetDestChain()),
			sdk.NewAttribute(types.AttributeKeyRelayChain, packet.GetRelayChain()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	if packet.GetRelayChain() == chainName {
		if !k.routingKeeper.Authenticate(ctx, packet.GetSourceChain(), packet.GetDestChain(), packet.GetPort()) {
			return sdkerrors.ErrUnauthorized
		}

		if _, found = k.clientKeeper.GetClientState(ctx, packet.GetDestChain()); !found {
			return sdkerrors.Wrap(clienttypes.ErrClientNotFound, fromChain)
		}

		k.SetPacketCommitment(ctx, packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence(), commitment)
		// Emit Event with Packet data along with other packet information for relayer to pick up
		// and relay to other chain
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeSendPacket,
				sdk.NewAttribute(types.AttributeKeyData, string(packet.GetData())),
				sdk.NewAttribute(types.AttributeKeySequence, fmt.Sprintf("%d", packet.GetSequence())),
				sdk.NewAttribute(types.AttributeKeyPort, packet.GetPort()),
				sdk.NewAttribute(types.AttributeKeySrcChain, packet.GetSourceChain()),
				sdk.NewAttribute(types.AttributeKeyDstChain, packet.GetDestChain()),
				sdk.NewAttribute(types.AttributeKeyRelayChain, packet.GetRelayChain()),
				// we only support 1-hop packets now, and that is the most important hop for a relayer
				// (is it going to a chain I am connected to)
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			),
		})
	}

	return nil
}

// WriteAcknowledgement writes the packet execution acknowledgement to the state,
// which will be verified by the counterparty chain using AcknowledgePacket.
//
// CONTRACT:
//
// 1) For synchronous execution, this function is be called in the TIBC handler .
// For async handling, it needs to be called directly by the module which originally
// processed the packet.
//
// 2) Assumes that packet receipt has been written.
// previously by RecvPacket.
func (k Keeper) WriteAcknowledgement(
	ctx sdk.Context,
	packet exported.PacketI,
	acknowledgement []byte,
) error {
	if len(acknowledgement) == 0 {
		return sdkerrors.Wrap(types.ErrInvalidAcknowledgement, "acknowledgement cannot be empty")
	}

	// NOTE: TIBC app modules might have written the acknowledgement synchronously on
	// the OnRecvPacket callback so we need to check if the acknowledgement is already
	// set on the store and return an error if so.
	if k.HasPacketAcknowledgement(ctx, packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence()) {
		return types.ErrAcknowledgementExists
	}

	targetChain := packet.GetSourceChain()
	if len(packet.GetRelayChain()) > 0 && packet.GetDestChain() == k.clientKeeper.GetChainName(ctx) {
		targetChain = packet.GetRelayChain()
	}

	_, found := k.clientKeeper.GetClientState(ctx, targetChain)
	if !found {
		return clienttypes.ErrClientNotFound
	}

	// set the acknowledgement so that it can be verified on the other side
	k.SetPacketAcknowledgement(
		ctx, packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence(),
		types.CommitAcknowledgement(acknowledgement),
	)
	k.SetMaxAckSequence(ctx, packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence())

	// log that a packet acknowledgement has been written
	k.Logger(ctx).Info("acknowledged written", "packet", fmt.Sprintf("%v", packet))

	// emit an event that the relayer can query for
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWriteAck,
			sdk.NewAttribute(types.AttributeKeyData, string(packet.GetData())),
			sdk.NewAttribute(types.AttributeKeySequence, fmt.Sprintf("%d", packet.GetSequence())),
			sdk.NewAttribute(types.AttributeKeyPort, packet.GetPort()),
			sdk.NewAttribute(types.AttributeKeySrcChain, packet.GetSourceChain()),
			sdk.NewAttribute(types.AttributeKeyDstChain, packet.GetDestChain()),
			sdk.NewAttribute(types.AttributeKeyRelayChain, packet.GetRelayChain()),
			sdk.NewAttribute(types.AttributeKeyAck, string(acknowledgement)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return nil
}

// AcknowledgePacket is called by a module to process the acknowledgement of a
// packet previously sent by the calling module on a port to a counterparty
// module on the counterparty chain. Its intended usage is within the ante
// handler. AcknowledgePacket will clean up the packet commitment,
// which is no longer necessary since the packet has been received and acted upon.
func (k Keeper) AcknowledgePacket(
	ctx sdk.Context,
	packet exported.PacketI,
	acknowledgement []byte,
	proof []byte,
	proofHeight exported.Height,
) error {
	if err := k.ValidatePacket(ctx, packet); err != nil {
		return sdkerrors.Wrap(err, "AcknowledgePacket failed basic validation")
	}
	commitment := k.GetPacketCommitment(ctx, packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence())

	packetCommitment := types.CommitPacket(packet)

	// verify we sent the packet and haven't cleared it out yet
	if !bytes.Equal(commitment, packetCommitment) {
		return sdkerrors.Wrapf(types.ErrInvalidPacket, "commitment bytes are not equal: got (%v), expected (%v)", packetCommitment, commitment)
	}

	chainName := k.clientKeeper.GetChainName(ctx)
	fromChain := packet.GetDestChain()
	if packet.GetSourceChain() == chainName && len(packet.GetRelayChain()) > 0 {
		fromChain = packet.GetRelayChain()
	}

	clientState, found := k.clientKeeper.GetClientState(ctx, fromChain)
	if !found {
		return sdkerrors.Wrap(clienttypes.ErrClientNotFound, fromChain)
	}

	ackCommitment := types.CommitAcknowledgement(acknowledgement)
	if err := clientState.VerifyPacketAcknowledgement(ctx,
		k.clientKeeper.ClientStore(ctx, fromChain), k.cdc, proofHeight,
		proof, packet.GetSourceChain(), packet.GetDestChain(),
		packet.GetSequence(), ackCommitment,
	); err != nil {
		return sdkerrors.Wrapf(err, "failed packet acknowledgement verification for client (%s)", fromChain)
	}

	// Delete packet commitment, since the packet has been acknowledged, the commitement is no longer necessary
	k.deletePacketCommitment(ctx, packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence())
	k.SetMaxAckSequence(ctx, packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence())

	// log that a packet has been acknowledged
	k.Logger(ctx).Info("packet acknowledged", "packet", fmt.Sprintf("%v", packet))

	// emit an event marking that we have processed the acknowledgement
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAcknowledgePacket,
			sdk.NewAttribute(types.AttributeKeyData, string(packet.GetData())),
			sdk.NewAttribute(types.AttributeKeySequence, fmt.Sprintf("%d", packet.GetSequence())),
			sdk.NewAttribute(types.AttributeKeyPort, packet.GetPort()),
			sdk.NewAttribute(types.AttributeKeySrcChain, packet.GetSourceChain()),
			sdk.NewAttribute(types.AttributeKeyDstChain, packet.GetDestChain()),
			sdk.NewAttribute(types.AttributeKeyRelayChain, packet.GetRelayChain()),
			sdk.NewAttribute(types.AttributeKeyAck, string(acknowledgement)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	if packet.GetRelayChain() == chainName {
		if _, found = k.clientKeeper.GetClientState(ctx, packet.GetSourceChain()); !found {
			return sdkerrors.Wrap(clienttypes.ErrClientNotFound, fromChain)
		}
		// set the acknowledgement so that it can be verified on the other side
		k.SetPacketAcknowledgement(
			ctx, packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence(),
			ackCommitment,
		)
		// emit an event that the relayer can query for
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeWriteAck,
				sdk.NewAttribute(types.AttributeKeyData, string(packet.GetData())),
				sdk.NewAttribute(types.AttributeKeySequence, fmt.Sprintf("%d", packet.GetSequence())),
				sdk.NewAttribute(types.AttributeKeyPort, packet.GetPort()),
				sdk.NewAttribute(types.AttributeKeySrcChain, packet.GetSourceChain()),
				sdk.NewAttribute(types.AttributeKeyDstChain, packet.GetDestChain()),
				sdk.NewAttribute(types.AttributeKeyRelayChain, packet.GetRelayChain()),
				sdk.NewAttribute(types.AttributeKeyAck, string(acknowledgement)),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			),
		})
	}
	return nil
}

// CleanPacket is called by users to send an TIBC clean packet to clean ack and receipt whose
// sequence smaller than the sequence in clean packet.
func (k Keeper) CleanPacket(ctx sdk.Context, cleanPacket exported.CleanPacketI) error {
	sourceChain := k.clientKeeper.GetChainName(ctx)
	if err := cleanPacket.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "packet failed basic validation")
	}
	if err := k.ValidateCleanPacket(ctx, types.CleanPacket{
		Sequence:         cleanPacket.GetSequence(),
		SourceChain:      sourceChain,
		DestinationChain: cleanPacket.GetDestChain(),
		RelayChain:       cleanPacket.GetRelayChain(),
	}); err != nil {
		return sdkerrors.Wrap(err, "packet failed basic validation")
	}

	targetChain := cleanPacket.GetDestChain()
	if len(cleanPacket.GetRelayChain()) > 0 {
		targetChain = cleanPacket.GetRelayChain()
	}

	_, found := k.clientKeeper.GetClientState(ctx, targetChain)
	if !found {
		return clienttypes.ErrClientNotFound
	}

	k.SetCleanPacketCommitment(ctx, sourceChain, cleanPacket.GetDestChain(), cleanPacket.GetSequence())
	k.cleanAcknowledgementBySeq(ctx, sourceChain, cleanPacket.GetDestChain(), cleanPacket.GetSequence())
	k.cleanReceiptBySeq(ctx, sourceChain, cleanPacket.GetDestChain(), cleanPacket.GetSequence())

	// Emit Event with Packet data along with other packet information for relayer to pick up
	// and relay to other chain
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSendCleanPacket,
			sdk.NewAttribute(types.AttributeKeySequence, fmt.Sprintf("%d", cleanPacket.GetSequence())),
			sdk.NewAttribute(types.AttributeKeySrcChain, sourceChain),
			sdk.NewAttribute(types.AttributeKeyDstChain, cleanPacket.GetDestChain()),
			sdk.NewAttribute(types.AttributeKeyRelayChain, cleanPacket.GetRelayChain()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	k.Logger(ctx).Info("clean packet sent", "packet", fmt.Sprintf("%v", cleanPacket))
	return nil
}

// RecvCleanPacket is called by relayer to receive & process an TIBC clean packet
// sent on the counterparty chain.
func (k Keeper) RecvCleanPacket(
	ctx sdk.Context,
	cleanPacket exported.CleanPacketI,
	proof []byte,
	proofHeight exported.Height,
) error {
	if err := k.ValidateCleanPacket(ctx, cleanPacket); err != nil {
		return sdkerrors.Wrap(err, "packet failed basic validation")
	}
	chainName := k.clientKeeper.GetChainName(ctx)
	fromChain := cleanPacket.GetSourceChain()
	if cleanPacket.GetDestChain() == chainName && len(cleanPacket.GetRelayChain()) > 0 {
		fromChain = cleanPacket.GetRelayChain()
	}
	targetClient, found := k.clientKeeper.GetClientState(ctx, fromChain)

	if !found {
		return sdkerrors.Wrap(clienttypes.ErrClientNotFound, fromChain)
	}

	if err := targetClient.VerifyPacketCleanCommitment(ctx,
		k.clientKeeper.ClientStore(ctx, fromChain), k.cdc, proofHeight,
		proof, cleanPacket.GetSourceChain(), cleanPacket.GetDestChain(),
		cleanPacket.GetSequence(),
	); err != nil {
		return sdkerrors.Wrapf(err, "failed packet commitment verification for client (%s)", fromChain)
	}

	k.cleanAcknowledgementBySeq(ctx, cleanPacket.GetSourceChain(), cleanPacket.GetDestChain(), cleanPacket.GetSequence())
	k.cleanReceiptBySeq(ctx, cleanPacket.GetSourceChain(), cleanPacket.GetDestChain(), cleanPacket.GetSequence())
	k.SetCleanPacketCommitment(ctx, cleanPacket.GetSourceChain(), cleanPacket.GetDestChain(), cleanPacket.GetSequence())

	// emit an event that the relayer can query for
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRecvCleanPacket,
			sdk.NewAttribute(types.AttributeKeySequence, fmt.Sprintf("%d", cleanPacket.GetSequence())),
			sdk.NewAttribute(types.AttributeKeySrcChain, cleanPacket.GetSourceChain()),
			sdk.NewAttribute(types.AttributeKeyDstChain, cleanPacket.GetDestChain()),
			sdk.NewAttribute(types.AttributeKeyRelayChain, cleanPacket.GetRelayChain()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	if cleanPacket.GetRelayChain() == chainName {
		if _, found = k.clientKeeper.GetClientState(ctx, cleanPacket.GetDestChain()); !found {
			return sdkerrors.Wrap(clienttypes.ErrClientNotFound, fromChain)
		}
		// Emit Event with Packet data along with other packet information for relayer to pick up
		// and relay to other chain
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeSendCleanPacket,
				sdk.NewAttribute(types.AttributeKeySequence, fmt.Sprintf("%d", cleanPacket.GetSequence())),
				sdk.NewAttribute(types.AttributeKeySrcChain, cleanPacket.GetSourceChain()),
				sdk.NewAttribute(types.AttributeKeyDstChain, cleanPacket.GetDestChain()),
				sdk.NewAttribute(types.AttributeKeyRelayChain, cleanPacket.GetRelayChain()),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			),
		})
	}
	return nil
}
