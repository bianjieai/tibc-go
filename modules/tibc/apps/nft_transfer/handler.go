package nft_transfer

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/keeper"
	nfttransfertypes "github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/types"
)

// NewHandler defines the TIBC nft transfer handler
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *nfttransfertypes.MsgNftTransfer:
			res, err := k.NftTransfer(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized TIBC message type: %T", msg)
		}
	}
}
