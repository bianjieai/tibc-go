package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

// ClientKeeper expected account IBC client keeper
type ClientKeeper interface {
	GetClientState(ctx sdk.Context, clientID string) (exported.ClientState, bool)
	GetClientConsensusState(ctx sdk.Context, clientID string, height exported.Height) (exported.ConsensusState, bool)
	ClientStore(ctx sdk.Context, clientID string) sdk.KVStore
}

// PortKeeper expected account IBC port keeper
type RoutingKeeper interface {
	Authenticate(ctx sdk.Context, sourceChain, destinationChain, port string) bool
}
