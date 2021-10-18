package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/types"
)

var _ types.MsgServer = Keeper{}

func (k Keeper) NftTransfer(goCtx context.Context, msg *types.MsgNftTransfer) (*types.MsgNftTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	if err := k.SendNftTransfer(
		ctx, msg.Class, msg.Id, sender, msg.Receiver, msg.DestChain, msg.RealayChain, msg.DestContract); err != nil {
		return nil, err
	}

	k.Logger(ctx).Info("TIBC non fungible token transfer", "nft", msg.Id, "sender", msg.Sender, "receiver", msg.Receiver)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeNftTransfer,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyReceiver, msg.Receiver),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	})

	return &types.MsgNftTransferResponse{}, nil
}
