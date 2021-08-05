package keeper

import (
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"

	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	"github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
)

// Keeper defines the TIBC routing keeper
type Keeper struct {
	Router *types.Router
}

// NewKeeper creates a new TIBC connection Keeper instance
func NewKeeper() Keeper {
	return Keeper{}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+host.ModuleName+"/"+types.SubModuleName)
}

// SetRouter sets the Router in IBC Keeper and seals it. The method panics if
// there is an existing router that's already sealed.
func (k *Keeper) SetRouter(rtr *types.Router) {
	if k.Router != nil && k.Router.Sealed() {
		panic("cannot reset a sealed router")
	}
	k.Router = rtr
	k.Router.Seal()
}

func (k Keeper) SetRoutingRules(ctx sdk.Context, sourceChain, destinationChain string, rules []string) error {
	return nil
}

func (k Keeper) GetRoutingRules(ctx sdk.Context, sourceChain, destinationChain string) []string {
	return nil
}

func (k Keeper) Authenticate(ctx sdk.Context, sourceChain, destinationChain, port string) bool {
	return true
}
