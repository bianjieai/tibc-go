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

// SendPacket is called by a module in order to send an IBC packet on a channel
// end owned by the calling module to the corresponding module on the counterparty
// chain.
func (k Keeper) SendPacket(
	ctx sdk.Context,
	packet exported.PacketI,
) error {
	if err := packet.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "packet failed basic validation")
	}

	if len(packet.GetRelayChain())>0{
		_, found := k.clientKeeper.GetClientState(ctx, packet.GetRelayChain())
		if !found {
			return clienttypes.ErrConsensusStateNotFound
		}
		////prevent accidental sends with clients that cannot be updated
		//if clientState.IsFrozen() {
		//	return sdkerrors.Wrapf(clienttypes.ErrClientFrozen, "cannot send packet on a frozen client with ID %s", packet.GetDestChain())
		//}
	}else{
		_, found := k.clientKeeper.GetClientState(ctx, packet.GetSourceChain())
		if !found {
			return clienttypes.ErrConsensusStateNotFound
		}
		////prevent accidental sends with clients that cannot be updated
		//if clientState.IsFrozen() {
		//	return sdkerrors.Wrapf(clienttypes.ErrClientFrozen, "cannot send packet on a frozen client with ID %s", packet.GetDestChain())
		//}
	}

	nextSequenceSend, found := k.GetNextSequenceSend(ctx, packet.GetSourceChain(), packet.GetDestChain())
	if !found {
		return sdkerrors.Wrapf(
			types.ErrSequenceSendNotFound,
			"source chain: %s, dest chain: %s", packet.GetSourceChain(), packet.GetDestChain(),
		)
	}

	if packet.GetSequence() != nextSequenceSend {
		return sdkerrors.Wrapf(
			types.ErrInvalidPacket,
			"packet sequence ≠ next send sequence (%d ≠ %d)", packet.GetSequence(), nextSequenceSend,
		)
	}

	commitment := types.CommitPacket(k.cdc, packet)

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

// RecvPacket is called by a module in order to receive & process an IBC packet
// sent on the corresponding channel end on the counterparty chain.
func (k Keeper) RecvPacket(
	ctx sdk.Context,
	packet exported.PacketI,
	proof []byte,
	proofHeight exported.Height,
) error {
	commitment := types.CommitPacket(k.cdc, packet)
	var isRelay bool
	var targetClientID string
	if packet.GetDestChain() == k.clientKeeper.GetChainName(ctx){
		if len(packet.GetRelayChain()) > 0 {
			targetClientID = packet.GetRelayChain()
		}else {
			targetClientID = packet.GetSourceChain()
		}
	}else{
		isRelay = true
		targetClientID = packet.GetSourceChain()
	}

	targetClient, found := k.clientKeeper.GetClientState(ctx, targetClientID)
	if !found {
		return sdkerrors.Wrap(clienttypes.ErrClientNotFound, targetClientID)
	}

	// verify that the counterparty did commit to sending this packet
	if err := targetClient.VerifyPacketCommitment(ctx,
		k.clientKeeper.ClientStore(ctx, targetClientID), k.cdc, proofHeight,
		proof, packet.GetSourceChain(), packet.GetDestChain(),
		packet.GetSequence(), commitment,
	); err != nil {
		return sdkerrors.Wrapf(err, "failed packet commitment verification for client (%s)", targetClientID)
	}

	// check if the packet receipt has been received already for unordered channels
	_, found = k.GetPacketReceipt(ctx, packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence())
	if found {
		return sdkerrors.Wrapf(
			types.ErrInvalidPacket,
			"packet sequence (%d) already has been received", packet.GetSequence(),
		)
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

	if isRelay{
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
// 1) For synchronous execution, this function is be called in the IBC handler .
// For async handling, it needs to be called directly by the module which originally
// processed the packet.
//
// 2) Assumes that packet receipt has been written (unordered), or nextSeqRecv was incremented (ordered)
// previously by RecvPacket.
func (k Keeper) WriteAcknowledgement(
	ctx sdk.Context,
	packet exported.PacketI,
	acknowledgement []byte,
) error {
	// NOTE: IBC app modules might have written the acknowledgement synchronously on
	// the OnRecvPacket callback so we need to check if the acknowledgement is already
	// set on the store and return an error if so.
	if k.HasPacketAcknowledgement(ctx, packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence()) {
		return types.ErrAcknowledgementExists
	}

	if len(acknowledgement) == 0 {
		return sdkerrors.Wrap(types.ErrInvalidAcknowledgement, "acknowledgement cannot be empty")
	}

	// set the acknowledgement so that it can be verified on the other side
	k.SetPacketAcknowledgement(
		ctx, packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence(),
		types.CommitAcknowledgement(acknowledgement),
	)

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
// packet previously sent by the calling module on a channel to a counterparty
// module on the counterparty chain. Its intended usage is within the ante
// handler. AcknowledgePacket will clean up the packet commitment,
// which is no longer necessary since the packet has been received and acted upon.
// It will also increment NextSequenceAck in case of ORDERED channels.
func (k Keeper) AcknowledgePacket(
	ctx sdk.Context,
	packet exported.PacketI,
	acknowledgement []byte,
	proof []byte,
	proofHeight exported.Height,
) error {
	commitment := k.GetPacketCommitment(ctx, packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence())

	packetCommitment := types.CommitPacket(k.cdc, packet)

	// verify we sent the packet and haven't cleared it out yet
	if !bytes.Equal(commitment, packetCommitment) {
		return sdkerrors.Wrapf(types.ErrInvalidPacket, "commitment bytes are not equal: got (%v), expected (%v)", packetCommitment, commitment)
	}

	targetClientID := packet.GetDestChain()

	clientState, found := k.clientKeeper.GetClientState(ctx, targetClientID)
	if !found {
		return sdkerrors.Wrap(clienttypes.ErrClientNotFound, targetClientID)
	}

	if err := clientState.VerifyPacketAcknowledgement(ctx,
		k.clientKeeper.ClientStore(ctx, targetClientID), k.cdc, proofHeight,
		proof, packet.GetSourceChain(), packet.GetDestChain(),
		packet.GetSequence(), acknowledgement,
	); err != nil {
		return sdkerrors.Wrapf(err, "failed packet acknowledgement verification for client (%s)", targetClientID)
	}

	// Delete packet commitment, since the packet has been acknowledged, the commitement is no longer necessary
	k.deletePacketCommitment(ctx, packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence())

	// log that a packet has been acknowledged
	k.Logger(ctx).Info("packet acknowledged", "packet", fmt.Sprintf("%v", packet))

	// emit an event marking that we have processed the acknowledgement
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAcknowledgePacket,
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

	return nil
}

// CleanPacket.
func (k Keeper) CleanPacket(
	ctx sdk.Context,
	packet exported.PacketI,
) error {
	if err := packet.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "packet failed basic validation")
	}

	nextSequenceSend, found := k.GetNextSequenceSend(ctx, packet.GetSourceChain(), packet.GetDestChain())
	if !found {
		return sdkerrors.Wrapf(
			types.ErrSequenceSendNotFound,
			"source chain: %s, dest chain: %s", packet.GetSourceChain(), packet.GetDestChain(),
		)
	}

	if packet.GetSequence() != nextSequenceSend {
		return sdkerrors.Wrapf(
			types.ErrInvalidPacket,
			"packet sequence ≠ next send sequence (%d ≠ %d)", packet.GetSequence(), nextSequenceSend,
		)
	}

	commitment := types.CommitPacket(k.cdc, packet)

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

// CleanPacket.
func (k Keeper) RecvCleanPacket(
	ctx sdk.Context,
	packet exported.PacketI,
	proof []byte,
	proofHeight exported.Height,
) error {
	commitment := types.CommitPacket(k.cdc, packet)
	targetClientID := packet.GetSourceChain()
	targetClient, found := k.clientKeeper.GetClientState(ctx, targetClientID)
	if !found {
		return sdkerrors.Wrap(clienttypes.ErrClientNotFound, targetClientID)
	}

	if err := targetClient.VerifyPacketCommitment(ctx,
		k.clientKeeper.ClientStore(ctx, targetClientID), k.cdc, proofHeight,
		proof, packet.GetSourceChain(), packet.GetDestChain(),
		packet.GetSequence(), commitment,
	); err != nil {
		return sdkerrors.Wrapf(err, "failed packet commitment verification for client (%s)", targetClientID)
	}

	_, found = k.GetPacketReceipt(ctx, packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence())
	if !found {
		return sdkerrors.Wrapf(
			types.ErrInvalidPacket,
			"packet sequence does not exist!", packet.GetSequence(),
		)
	}
	found = k.HasPacketAcknowledgement(ctx, packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence())
	if !found {
		return sdkerrors.Wrapf(
			types.ErrInvalidPacket,
			"packet sequence does not exist!", packet.GetSequence(),
		)
	}

	k.deletePacketAcknowledgement(ctx, packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence())
	k.deletePacketReceipt(ctx, packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence())
	return nil
}
