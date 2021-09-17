package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/client/cli"
	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/keeper"
	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
)

// TODO: move to cli
var (
	CreateClientProposalHandler    = govclient.NewProposalHandler(cli.NewCreateClientProposalCmd, EmptyRESTHandler)
	UpgradeClientProposalHandler   = govclient.NewProposalHandler(cli.NewUpgradeClientProposalCmd, EmptyRESTHandler)
	RegisterRelayerProposalHandler = govclient.NewProposalHandler(cli.NewRegisterRelayerProposalCmd, EmptyRESTHandler)
)

func EmptyRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "tibc",
		Handler:  nil,
	}
}

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
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized tibc proposal content type: %T", c)
		}
	}
}
