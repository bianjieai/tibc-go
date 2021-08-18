package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
)

// HandleSetRoutingRulesProposal will try to set routing rules if and only if the proposal passes.
func (k Keeper) HandleSetRoutingRulesProposal(ctx sdk.Context, p *types.SetRoutingRulesProposal) error {
	return k.SetRoutingRules(ctx, p.Rules)
}
