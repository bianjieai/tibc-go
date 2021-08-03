package ibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	packettypes "github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
	"github.com/bianjieai/tibc-go/modules/tibc/core/keeper"
)

// NewHandler defines the IBC handler
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		// IBC client msg interface types
		case *clienttypes.MsgCreateClient:
			res, err := k.CreateClient(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *clienttypes.MsgUpdateClient:
			res, err := k.UpdateClient(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *clienttypes.MsgUpgradeClient:
			res, err := k.UpgradeClient(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *clienttypes.MsgSubmitMisbehaviour:
			res, err := k.SubmitMisbehaviour(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		// IBC packet msgs get routed to the appropriate module callback
		case *packettypes.MsgRecvPacket:
			res, err := k.RecvPacket(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *packettypes.MsgAcknowledgement:
			res, err := k.Acknowledgement(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized IBC message type: %T", msg)
		}
	}
}
