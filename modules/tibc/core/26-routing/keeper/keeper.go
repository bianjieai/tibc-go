package keeper

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/libs/log"

	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	"github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
)

// Keeper defines the TIBC routing keeper
type Keeper struct {
	Router   *types.Router
	storeKey sdk.StoreKey
}

// NewKeeper creates a new TIBC connection Keeper instance
func NewKeeper(key sdk.StoreKey) Keeper {
	return Keeper{
		storeKey: key,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+host.ModuleName+"/"+types.SubModuleName)
}

// SetRouter sets the Router in TIBC Keeper and seals it. The method panics if
// there is an existing router that's already sealed.
func (k *Keeper) SetRouter(rtr *types.Router) {
	if k.Router != nil && k.Router.Sealed() {
		panic("cannot reset a sealed router")
	}
	k.Router = rtr
	k.Router.Seal()
}

func (k Keeper) SetRoutingRules(ctx sdk.Context, rules []string) error {
	for _, rule := range rules {
		valid, _ := regexp.MatchString(types.RulePattern, rule)
		if !valid {
			return sdkerrors.Wrap(types.ErrInvalidRule, "invalid rule")
		}
	}
	routingBz, err := json.Marshal(rules)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrFailMarshalRules, "failed to marshal rules: %s", err.Error())
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(RoutingRulesKey(), routingBz)
	return nil
}

func (k Keeper) GetRoutingRules(ctx sdk.Context) ([]string, bool) {
	var rules []string
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(RoutingRulesKey())
	if bz == nil {
		return nil, false
	}
	json.Unmarshal(bz, &rules)
	return rules, true
}

func (k Keeper) Authenticate(ctx sdk.Context, sourceChain, destinationChain, port string) bool {
	rules, found := k.GetRoutingRules(ctx)
	if !found {
		return false
	}
	flag := false
	for _, rule := range rules {
		flag, _ = regexp.MatchString(ConvWildcardToRegular(rule), sourceChain+"."+destinationChain+"."+port)
		if flag {
			break
		}
	}
	return flag
}

func ConvWildcardToRegular(wildcard string) string {
	regular := strings.Replace(wildcard, ".", "\\.", -1)
	regular = strings.Replace(regular, "*", ".*", -1)
	regular = strings.Replace(regular, "?", ".", -1)
	regular = "^" + regular + "$"
	return regular
}

// RoutingRulesPath defines the routing rules store path
func RoutingRulesPath() string {
	return fmt.Sprintf("Routing/Rules")
}

func RoutingRulesKey() []byte {
	return []byte(RoutingRulesPath())
}
