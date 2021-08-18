package routing

import (
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/client/cli"
	"github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/keeper"
	"github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
)

var (
	SetRoutingRulesProposalHandler = govclient.NewProposalHandler(cli.NewSetRoutingRulesProposalCmd, EmptyRESTHandler)
)

func EmptyRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "tibc",
		Handler:  nil,
	}
}

// NewSetRoutingProposalHandler defines the routing manager proposal handler
func NewSetRoutingProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.SetRoutingRulesProposal:
			return k.HandleSetRoutingRulesProposal(ctx, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized tibc proposal content type: %T", c)
		}
	}
}
