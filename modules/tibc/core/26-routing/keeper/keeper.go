package keeper

import (
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"

	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	"github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
)

// Keeper defines the TIBC routing keeper
type Keeper struct {
}

// NewKeeper creates a new TIBC connection Keeper instance
func NewKeeper() Keeper {
	return Keeper{}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+host.ModuleName+"/"+types.SubModuleName)
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
