package cli

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	clientcli "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/client/cli"
	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	routingcli "github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/client/cli"
	routingtypes "github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
	"github.com/bianjieai/tibc-go/modules/tibc/core/keeper"
)

// GovHandlers defines the client manager proposal handlers
var GovHandlers = []govclient.ProposalHandler{
	govclient.NewProposalHandler(clientcli.NewCreateClientProposalCmd),
	govclient.NewProposalHandler(clientcli.NewUpgradeClientProposalCmd),
	govclient.NewProposalHandler(clientcli.NewRegisterRelayerProposalCmd),
	govclient.NewProposalHandler(routingcli.NewSetRoutingRulesProposalCmd),
}

// NewProposalHandler defines the client manager proposal handler
func NewProposalHandler(k *keeper.Keeper) govv1beta1.Handler {
	return func(ctx sdk.Context, content govv1beta1.Content) error {
		switch c := content.(type) {
		case *clienttypes.CreateClientProposal:
			return k.ClientKeeper.HandleCreateClientProposal(ctx, c)
		case *clienttypes.UpgradeClientProposal:
			return k.ClientKeeper.HandleUpgradeClientProposal(ctx, c)
		case *clienttypes.RegisterRelayerProposal:
			return k.ClientKeeper.HandleRegisterRelayerProposal(ctx, c)
		case *routingtypes.SetRoutingRulesProposal:
			return k.RoutingKeeper.HandleSetRoutingRulesProposal(ctx, c)
		default:
			return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized tibc proposal content type: %T", c)
		}
	}
}
