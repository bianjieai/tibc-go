package client

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/client/cli"
	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/keeper"
	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

var (
	CreateClientProposalHandler    = govclient.NewProposalHandler(cli.NewCreateClientProposalCmd, nil)
	UpgradeClientProposalHandler   = govclient.NewProposalHandler(cli.NewUpgradeClientProposalCmd, nil)
	RegisterRelayerProposalHandler = govclient.NewProposalHandler(cli.NewRegisterRelayerProposalCmd, nil)
)

// NewClientProposalHandler defines the client manager proposal handler
func NewClientProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.CreateClientProposal:
			return k.HandleCreateClientProposal(ctx, c)
		case *types.UpgradeClientProposal:
			return k.HandleUpgradeClientProposal(ctx, c)
		case *types.RegisterRelayerProposal:
			return k.HandleRegisterRelayerProposal(ctx, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized ibc proposal content type: %T", c)
		}
	}
}
