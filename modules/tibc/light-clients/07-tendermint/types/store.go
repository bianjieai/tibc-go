package types

import (
	"encoding/binary"
	"strings"

	errorsmod "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

// KeyProcessedTime is appended to consensus state key to store the processed time
var (
	KeyIterateConsensusStatePrefix = "iterateConsensusStates"
	KeyProcessedTime               = []byte("/processedTime")
)

// GetConsensusState retrieves the consensus state from the client prefixed
// store. An error is returned if the consensus state does not exist.
func GetConsensusState(store storetypes.KVStore, cdc codec.BinaryCodec, height exported.Height) (*ConsensusState, error) {
	bz := store.Get(host.ConsensusStateKey(height))
	if bz == nil {
		return nil, errorsmod.Wrapf(
			clienttypes.ErrConsensusStateNotFound,
			"consensus state does not exist for height %s",
			height,
		)
	}

	consensusStateI, err := clienttypes.UnmarshalConsensusState(cdc, bz)
	if err != nil {
		return nil, errorsmod.Wrapf(clienttypes.ErrInvalidConsensus, "unmarshal error: %v", err)
	}

	consensusState, ok := consensusStateI.(*ConsensusState)
	if !ok {
		return nil, errorsmod.Wrapf(
			clienttypes.ErrInvalidConsensus,
			"invalid consensus type %T, expected %T",
			consensusState, &ConsensusState{},
		)
	}

	return consensusState, nil
}

// IterateProcessedTime iterates through the prefix store and applies the callback.
// If the cb returns true, then iterator will close and stop.
func IterateProcessedTime(store storetypes.KVStore, cb func(key, val []byte) bool) {
	iterator := storetypes.KVStorePrefixIterator(store, []byte(host.KeyConsensusStatePrefix))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		keySplit := strings.Split(string(key), "/")
		// processed time key in prefix store has format: "consensusState/<height>/processedTime"
		if len(keySplit) != 3 || keySplit[2] != "processedTime" {
			// ignore all consensus state keys
			continue
		}

		if cb(iterator.Key(), iterator.Value()) {
			break
		}
	}
}

// GetHeightFromIterationKey takes an iteration key and returns the height that it references
func GetHeightFromIterationKey(iterKey []byte) exported.Height {
	bigEndianBytes := iterKey[len([]byte(KeyIterateConsensusStatePrefix)):]
	revisionBytes := bigEndianBytes[0:8]
	heightBytes := bigEndianBytes[8:]
	revision := binary.BigEndian.Uint64(revisionBytes)
	height := binary.BigEndian.Uint64(heightBytes)
	return clienttypes.NewHeight(revision, height)
}

// IterateConsensusStateAscending iterates through the consensus states in ascending order. It calls the provided
// callback on each height, until stop=true is returned.
func IterateConsensusStateAscending(clientStore storetypes.KVStore,
	cb func(height exported.Height) (stop bool)) {
	iterator := storetypes.KVStorePrefixIterator(clientStore, []byte(KeyIterateConsensusStatePrefix))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		iterKey := iterator.Key()
		height := GetHeightFromIterationKey(iterKey)
		if cb(height) {
			return
		}
	}
}

// ProcessedTime Store code

// ProcessedTimeKey returns the key under which the processed time will be stored in the client store.
func ProcessedTimeKey(height exported.Height) []byte {
	return append(host.ConsensusStateKey(height), KeyProcessedTime...)
}

// SetProcessedTime stores the time at which a header was processed and the corresponding consensus state was created.
// This is useful when validating whether a packet has reached the specified delay period in the tendermint client's
// verification functions
func SetProcessedTime(clientStore storetypes.KVStore, height exported.Height, timeNs uint64) {
	key := ProcessedTimeKey(height)
	val := sdk.Uint64ToBigEndian(timeNs)
	clientStore.Set(key, val)
}

// GetProcessedTime gets the time (in nanoseconds) at which this chain received and processed a tendermint header.
// This is used to validate that a received packet has passed the delay period.
func GetProcessedTime(clientStore storetypes.KVStore, height exported.Height) (uint64, bool) {
	key := ProcessedTimeKey(height)
	bz := clientStore.Get(key)
	if bz == nil {
		return 0, false
	}
	return sdk.BigEndianToUint64(bz), true
}

// IterationKey returns the key under which the consensus state key will be stored.
// The iteration key is a BigEndian representation of the consensus state key to support efficient iteration.
func IterationKey(height exported.Height) []byte {
	heightBytes := bigEndianHeightBytes(height)
	return append([]byte(KeyIterateConsensusStatePrefix), heightBytes...)
}

// SetIterationKey stores the consensus state key under a key that is more efficient for ordered iteration
func SetIterationKey(clientStore storetypes.KVStore, height exported.Height) {
	key := IterationKey(height)
	val := host.ConsensusStateKey(height)
	clientStore.Set(key, val)
}

// GetIterationKey returns the consensus state key stored under the efficient iteration key.
// NOTE: This function is currently only used for testing purposes
func GetIterationKey(clientStore storetypes.KVStore, height exported.Height) []byte {
	key := IterationKey(height)
	return clientStore.Get(key)
}

// setConsensusMetadata sets context time as processed time and set context height as processed height
// as this is internal tendermint light client logic.
// client state and consensus state will be set by client keeper
// set iteration key to provide ability for efficient ordered iteration of consensus states.
func setConsensusMetadata(ctx sdk.Context, clientStore storetypes.KVStore, height exported.Height) {
	setConsensusMetadataWithValues(clientStore, height, clienttypes.GetSelfHeight(ctx), uint64(ctx.BlockTime().UnixNano()))
}

// deleteConsensusMetadata deletes the metadata stored for a particular consensus state.
func deleteConsensusMetadata(clientStore storetypes.KVStore, height exported.Height) {
	deleteProcessedTime(clientStore, height)
	deleteIterationKey(clientStore, height)
}

// setConsensusMetadataWithValues sets the consensus metadata with the provided values
func setConsensusMetadataWithValues(
	clientStore storetypes.KVStore, height,
	processedHeight exported.Height,
	processedTime uint64,
) {
	SetProcessedTime(clientStore, height, processedTime)
	SetIterationKey(clientStore, height)
}

// deleteConsensusState deletes the consensus state at the given height
func deleteConsensusState(clientStore storetypes.KVStore, height exported.Height) {
	key := host.ConsensusStateKey(height)
	clientStore.Delete(key)
}

// deleteProcessedTime deletes the processedTime for a given height
func deleteProcessedTime(clientStore storetypes.KVStore, height exported.Height) {
	key := ProcessedTimeKey(height)
	clientStore.Delete(key)
}

// deleteIterationKey deletes the iteration key for a given height
func deleteIterationKey(clientStore storetypes.KVStore, height exported.Height) {
	key := IterationKey(height)
	clientStore.Delete(key)
}

func bigEndianHeightBytes(height exported.Height) []byte {
	heightBytes := make([]byte, 16)
	binary.BigEndian.PutUint64(heightBytes, height.GetRevisionNumber())
	binary.BigEndian.PutUint64(heightBytes[8:], height.GetRevisionHeight())
	return heightBytes
}
