package keeper

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

// Keeper represents a type that grants read and write permissions to any client
// state information
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           codec.BinaryCodec
	paramSpace    paramtypes.Subspace
	stakingKeeper types.StakingKeeper
}

// NewKeeper creates a new NewKeeper instance
func NewKeeper(cdc codec.BinaryCodec, key sdk.StoreKey, paramSpace paramtypes.Subspace, sk types.StakingKeeper) Keeper {
	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		paramSpace:    paramSpace,
		stakingKeeper: sk,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+host.ModuleName+"/"+types.SubModuleName)
}

// GetClientState gets a particular client from the store
func (k Keeper) GetClientState(ctx sdk.Context, chainName string) (exported.ClientState, bool) {
	store := k.ClientStore(ctx, chainName)
	bz := store.Get(host.ClientStateKey())
	if bz == nil {
		return nil, false
	}

	clientState := k.MustUnmarshalClientState(bz)
	return clientState, true
}

// SetClientState sets a particular Client to the store
func (k Keeper) SetClientState(ctx sdk.Context, chainName string, clientState exported.ClientState) {
	store := k.ClientStore(ctx, chainName)
	store.Set(host.ClientStateKey(), k.MustMarshalClientState(clientState))
}

// GetClientConsensusState gets the stored consensus state from a client at a given height.
func (k Keeper) GetClientConsensusState(ctx sdk.Context, chainName string, height exported.Height) (exported.ConsensusState, bool) {
	store := k.ClientStore(ctx, chainName)
	bz := store.Get(host.ConsensusStateKey(height))
	if bz == nil {
		return nil, false
	}

	consensusState := k.MustUnmarshalConsensusState(bz)
	return consensusState, true
}

// SetClientConsensusState sets a ConsensusState to a particular client at the given
// height
func (k Keeper) SetClientConsensusState(ctx sdk.Context, chainName string, height exported.Height, consensusState exported.ConsensusState) {
	store := k.ClientStore(ctx, chainName)
	store.Set(host.ConsensusStateKey(height), k.MustMarshalConsensusState(consensusState))
}

// SetChainName sets a chain name to the tibc module
func (k Keeper) SetChainName(ctx sdk.Context, chainName string) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.KeyClientName), []byte(chainName))
}

// GetChainName return the chain name of the current chain
func (k Keeper) GetChainName(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(types.KeyClientName))
	return string(bz)
}

// IterateConsensusStates provides an iterator over all stored consensus states.
// objects. For each State object, cb will be called. If the cb returns true,
// the iterator will close and stop.
func (k Keeper) IterateConsensusStates(ctx sdk.Context, cb func(chainName string, cs types.ConsensusStateWithHeight) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, host.KeyClientStorePrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()

		keySplit := strings.Split(string(key), "/")
		// consensus key is in the format "clients/<chainName>/consensusStates/<height>"
		if len(keySplit) != 4 || keySplit[2] != string(host.KeyConsensusStatePrefix) {
			continue
		}
		chainName := keySplit[1]
		//revinum := sdk.BigEndianToUint64(key[35:43])
		//revihei := sdk.BigEndianToUint64(key[44:])
		heightBytes := keySplit[3]
		revisionUint64 := binary.BigEndian.Uint64([]byte(heightBytes[:8]))
		heightUint64 := binary.BigEndian.Uint64([]byte(heightBytes[8:]))
		height := types.MustParseHeight(fmt.Sprintf("%d-%d", revisionUint64, heightUint64))
		consensusState := k.MustUnmarshalConsensusState(iterator.Value())

		consensusStateWithHeight := types.NewConsensusStateWithHeight(height, consensusState)

		if cb(chainName, consensusStateWithHeight) {
			break
		}
	}
}

// GetAllGenesisClients returns all the clients in state with their client ids returned as IdentifiedClientState
func (k Keeper) GetAllGenesisClients(ctx sdk.Context) types.IdentifiedClientStates {
	var genClients types.IdentifiedClientStates
	k.IterateClients(ctx, func(chainName string, cs exported.ClientState) bool {
		genClients = append(genClients, types.NewIdentifiedClientState(chainName, cs))
		return false
	})

	return genClients.Sort()
}

// GetAllClientMetadata will take a list of IdentifiedClientState and return a list
// of IdentifiedGenesisMetadata necessary for exporting and importing client metadata
// into the client store.
func (k Keeper) GetAllClientMetadata(ctx sdk.Context, genClients []types.IdentifiedClientState) ([]types.IdentifiedGenesisMetadata, error) {
	genMetadata := make([]types.IdentifiedGenesisMetadata, 0)
	for _, ic := range genClients {
		cs, err := types.UnpackClientState(ic.ClientState)
		if err != nil {
			return nil, err
		}
		gms := cs.ExportMetadata(k.ClientStore(ctx, ic.ChainName))
		if len(gms) == 0 {
			continue
		}
		clientMetadata := make([]types.GenesisMetadata, len(gms))
		for i, metadata := range gms {
			cmd, ok := metadata.(types.GenesisMetadata)
			if !ok {
				return nil, sdkerrors.Wrapf(types.ErrInvalidClientMetadata, "expected metadata type: %T, got: %T",
					types.GenesisMetadata{}, cmd)
			}
			clientMetadata[i] = cmd
		}
		genMetadata = append(genMetadata, types.NewIdentifiedGenesisMetadata(
			ic.ChainName,
			clientMetadata,
		))
	}
	return genMetadata, nil
}

// SetAllClientMetadata takes a list of IdentifiedGenesisMetadata and stores all of the metadata in the client store at the appropriate paths.
func (k Keeper) SetAllClientMetadata(ctx sdk.Context, genMetadata []types.IdentifiedGenesisMetadata) {
	for _, igm := range genMetadata {
		// create client store
		store := k.ClientStore(ctx, igm.ChainName)
		// set all metadata kv pairs in client store
		for _, md := range igm.Metadata {
			store.Set(md.GetKey(), md.GetValue())
		}
	}
}

// GetAllConsensusStates returns all stored client consensus states.
func (k Keeper) GetAllConsensusStates(ctx sdk.Context) types.ClientsConsensusStates {
	clientConsStates := make(types.ClientsConsensusStates, 0)
	mapChainNameToConsStateIdx := make(map[string]int)

	k.IterateConsensusStates(ctx, func(clientName string, cs types.ConsensusStateWithHeight) bool {
		idx, ok := mapChainNameToConsStateIdx[clientName]
		if ok {
			clientConsStates[idx].ConsensusStates = append(clientConsStates[idx].ConsensusStates, cs)
			return false
		}

		clientConsState := types.ClientConsensusStates{
			ChainName:       clientName,
			ConsensusStates: []types.ConsensusStateWithHeight{cs},
		}

		clientConsStates = append(clientConsStates, clientConsState)
		mapChainNameToConsStateIdx[clientName] = len(clientConsStates) - 1
		return false
	})

	return clientConsStates.Sort()
}

// HasClientConsensusState returns if keeper has a ConsensusState for a particular
// client at the given height
func (k Keeper) HasClientConsensusState(ctx sdk.Context, chainName string, height exported.Height) bool {
	store := k.ClientStore(ctx, chainName)
	return store.Has(host.ConsensusStateKey(height))
}

// GetLatestClientConsensusState gets the latest ConsensusState stored for a given client
func (k Keeper) GetLatestClientConsensusState(ctx sdk.Context, chainName string) (exported.ConsensusState, bool) {
	clientState, ok := k.GetClientState(ctx, chainName)
	if !ok {
		return nil, false
	}
	return k.GetClientConsensusState(ctx, chainName, clientState.GetLatestHeight())
}

// IterateClients provides an iterator over all stored light client State
// objects. For each State object, cb will be called. If the cb returns true,
// the iterator will close and stop.
func (k Keeper) IterateClients(ctx sdk.Context, cb func(chainName string, cs exported.ClientState) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, host.KeyClientStorePrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		keySplit := strings.Split(string(iterator.Key()), "/")
		if keySplit[len(keySplit)-1] != host.KeyClientState {
			continue
		}
		clientState := k.MustUnmarshalClientState(iterator.Value())

		// key is tibc/{clientid}/clientState
		// Thus, keySplit[1] is chainName
		if cb(keySplit[1], clientState) {
			break
		}
	}
}

// GetAllClients returns all stored light client State objects.
func (k Keeper) GetAllClients(ctx sdk.Context) (states []exported.ClientState) {
	k.IterateClients(ctx, func(_ string, state exported.ClientState) bool {
		states = append(states, state)
		return false
	})
	return states
}

// ClientStore returns isolated prefix store for each client so they can read/write in separate
// namespace without being able to read/write other client's data
func (k Keeper) ClientStore(ctx sdk.Context, chainName string) sdk.KVStore {
	clientPrefix := []byte(fmt.Sprintf("%s/%s/", host.KeyClientStorePrefix, chainName))
	return prefix.NewStore(ctx.KVStore(k.storeKey), clientPrefix)
}

// RelayerStore returns isolated prefix store for each client so they can read/write in separate
// namespace without being able to read/write other relayer's data
func (k Keeper) RelayerStore(ctx sdk.Context) sdk.KVStore {
	return prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.KeyRelayers))
}
