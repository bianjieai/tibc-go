package client

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/keeper"
	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

// InitGenesis initializes the tibc client submodule's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, gs types.GenesisState) {
	// Set all client metadata first. This will allow client keeper to overwrite client and consensus state keys
	// if clients accidentally write to ClientKeeper reserved keys.
	if len(gs.ClientsMetadata) != 0 {
		k.SetAllClientMetadata(ctx, gs.ClientsMetadata)
	}

	for _, client := range gs.Clients {
		cs, ok := client.ClientState.GetCachedValue().(exported.ClientState)
		if !ok {
			panic("invalid client state")
		}

		k.SetClientState(ctx, client.ChainName, cs)
	}

	for _, cs := range gs.ClientsConsensus {
		for _, consState := range cs.ConsensusStates {
			consensusState, ok := consState.ConsensusState.GetCachedValue().(exported.ConsensusState)
			if !ok {
				panic(fmt.Sprintf("invalid consensus state with chain name %s at height %s", cs.ChainName, consState.Height))
			}

			k.SetClientConsensusState(ctx, cs.ChainName, consState.Height, consensusState)
		}
	}

	for _, rs := range gs.Relayers {
		k.RegisterRelayers(ctx, rs.ChainName, rs.Relayers)
	}
	k.SetChainName(ctx, gs.NativeChainName)
}

// ExportGenesis returns the tibc client submodule's exported genesis.
// NOTE: CreateLocalhost should always be false on export since a
// created localhost will be included in the exported clients.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {
	genClients := k.GetAllGenesisClients(ctx)
	clientsMetadata, err := k.GetAllClientMetadata(ctx, genClients)
	if err != nil {
		panic(err)
	}
	return types.GenesisState{
		Clients:          genClients,
		ClientsMetadata:  clientsMetadata,
		ClientsConsensus: k.GetAllConsensusStates(ctx),
		NativeChainName:  k.GetChainName(ctx),
		Relayers:         k.GetAllRelayers(ctx),
	}
}
