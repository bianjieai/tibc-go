package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

// ClientKeeper expected account TIBC client keeper
type ClientKeeper interface {
	GetClientState(ctx sdk.Context, chainName string) (exported.ClientState, bool)
	GetClientConsensusState(ctx sdk.Context, chainName string, height exported.Height) (exported.ConsensusState, bool)
	ClientStore(ctx sdk.Context, chainName string) sdk.KVStore
	GetChainName(ctx sdk.Context) string
}

// PortKeeper expected account TIBC port keeper
type RoutingKeeper interface {
	Authenticate(ctx sdk.Context, sourceChain, destinationChain, port string) bool
}
