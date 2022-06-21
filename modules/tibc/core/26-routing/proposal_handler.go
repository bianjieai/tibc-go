package routing

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/client/cli"
	"github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/keeper"
	"github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
)

var (
	SetRoutingRulesProposalHandler = govclient.NewProposalHandler(cli.NewSetRoutingRulesProposalCmd)
)

// NewSetRoutingProposalHandler defines the routing manager proposal handler
func NewSetRoutingProposalHandler(k keeper.Keeper) govv1beta1.Handler {
	return func(ctx sdk.Context, content govv1beta1.Content) error {
		switch c := content.(type) {
		case *types.SetRoutingRulesProposal:
			return k.HandleSetRoutingRulesProposal(ctx, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized tibc proposal content type: %T", c)
		}
	}
}
