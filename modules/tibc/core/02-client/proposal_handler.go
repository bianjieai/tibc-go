package client

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/client/cli"
	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/keeper"
	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
)

// TODO: move to cli
var (
	CreateClientProposalHandler    = govclient.NewProposalHandler(cli.NewCreateClientProposalCmd)
	UpgradeClientProposalHandler   = govclient.NewProposalHandler(cli.NewUpgradeClientProposalCmd)
	RegisterRelayerProposalHandler = govclient.NewProposalHandler(cli.NewRegisterRelayerProposalCmd)
)

// NewClientProposalHandler defines the client manager proposal handler
func NewClientProposalHandler(k keeper.Keeper) govv1beta1.Handler {
	return func(ctx sdk.Context, content govv1beta1.Content) error {
		switch c := content.(type) {
		case *types.CreateClientProposal:
			return k.HandleCreateClientProposal(ctx, c)
		case *types.UpgradeClientProposal:
			return k.HandleUpgradeClientProposal(ctx, c)
		case *types.RegisterRelayerProposal:
			return k.HandleRegisterRelayerProposal(ctx, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized tibc proposal content type: %T", c)
		}
	}
}
