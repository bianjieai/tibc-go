package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bianjieai/tibc-go/modules/tibc/apps/mt_transfer/types"
)

var _ types.MsgServer = Keeper{}

func (k Keeper) MtTransfer(goCtx context.Context, msg *types.MsgMtTransfer) (*types.MsgMtTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	if err := k.SendMtTransfer(
		ctx, msg.Class, msg.Id,
		sender, msg.Receiver,
		msg.DestChain, msg.RealayChain, msg.DestContract,
		msg.Amount,
	); err != nil {
		return nil, err
	}

	k.Logger(ctx).Info("TIBC multi token transfer", "mt", msg.Id, "sender", msg.Sender, "receiver", msg.Receiver)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMtTransfer,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyReceiver, msg.Receiver),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	})

	return &types.MsgMtTransferResponse{}, nil
}
