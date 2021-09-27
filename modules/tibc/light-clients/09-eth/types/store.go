package types

import (
	"encoding/binary"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"

	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

var (
	KeyIterateConsensusStatePrefix = "iterateConsensusStates"
	KeyIndexConsensusStatePrefix   = "consensusStatesIndex"
	KeyProcessedTime               = []byte("/processedTime")
)

func ConsensusStateIndexPath(hash common.Hash) string {
	return fmt.Sprintf("%s/%s", KeyIndexConsensusStatePrefix, hash)
}

func ConsensusStateIndexKey(hash common.Hash) []byte {
	return []byte(ConsensusStateIndexPath(hash))
}

// setConsensusMetadata sets context time as processed time and set context height as processed height
// as this is internal tendermint light client logic.
// client state and consensus state will be set by client keeper
// set iteration key to provide ability for efficient ordered iteration of consensus states.
func setConsensusMetadata(ctx sdk.Context, clientStore sdk.KVStore, height exported.Height) {
	setConsensusMetadataWithValues(clientStore, height, clienttypes.GetSelfHeight(ctx), uint64(ctx.BlockTime().UnixNano()))
}

// setConsensusMetadataWithValues sets the consensus metadata with the provided values
func setConsensusMetadataWithValues(
	clientStore sdk.KVStore, height,
	processedHeight exported.Height,
	processedTime uint64,
) {
	SetProcessedTime(clientStore, height, processedTime)
	SetIterationKey(clientStore, height)
}

// SetProcessedTime stores the time at which a header was processed and the corresponding consensus state was created.
// This is useful when validating whether a packet has reached the specified delay period in the tendermint client's
// verification functions
func SetProcessedTime(clientStore sdk.KVStore, height exported.Height, timeNs uint64) {
	key := ProcessedTimeKey(height)
	val := sdk.Uint64ToBigEndian(timeNs)
	clientStore.Set(key, val)
}

// GetProcessedTime gets the time (in nanoseconds) at which this chain received and processed a tendermint header.
// This is used to validate that a received packet has passed the delay period.
func GetProcessedTime(clientStore sdk.KVStore, height exported.Height) (uint64, bool) {
	key := ProcessedTimeKey(height)
	bz := clientStore.Get(key)
	if bz == nil {
		return 0, false
	}
	return sdk.BigEndianToUint64(bz), true
}

// SetIterationKey stores the consensus state key under a key that is more efficient for ordered iteration
func SetIterationKey(clientStore sdk.KVStore, height exported.Height) {
	key := IterationKey(height)
	val := host.ConsensusStateKey(height)
	clientStore.Set(key, val)
}

// GetIterationKey returns the consensus state key stored under the efficient iteration key.
// NOTE: This function is currently only used for testing purposes
func GetIterationKey(clientStore sdk.KVStore, height exported.Height) []byte {
	key := IterationKey(height)
	return clientStore.Get(key)
}

// GetConsensusState retrieves the consensus state from the client prefixed
// store. An error is returned if the consensus state does not exist.
func GetConsensusState(store sdk.KVStore,
	cdc codec.BinaryMarshaler, height exported.Height) (*ConsensusState, error) {
	bz := store.Get(host.ConsensusStateKey(height))
	if bz == nil {
		return nil, sdkerrors.Wrapf(
			clienttypes.ErrConsensusStateNotFound,
			"consensus state does not exist for height %s", height,
		)
	}

	consensusStateI, err := clienttypes.UnmarshalConsensusState(cdc, bz)
	if err != nil {
		return nil, sdkerrors.Wrapf(clienttypes.ErrInvalidConsensus, "unmarshal error: %v", err)
	}

	consensusState, ok := consensusStateI.(*ConsensusState)
	if !ok {
		return nil, sdkerrors.Wrapf(
			clienttypes.ErrInvalidConsensus,
			"invalid consensus type %T, expected %T", consensusState, &ConsensusState{},
		)
	}

	return consensusState, nil
}

// IterateConsensusStateAscending iterates through the consensus states in ascending order. It calls the provided
// callback on each height, until stop=true is returned.
func IterateConsensusStateAscending(clientStore sdk.KVStore,
	cb func(height exported.Height) (stop bool)) {
	iterator := sdk.KVStorePrefixIterator(clientStore, []byte(KeyIterateConsensusStatePrefix))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		iterKey := iterator.Key()
		height := GetHeightFromIterationKey(iterKey)
		if cb(height) {
			return
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

// deleteConsensusMetadata deletes the metadata stored for a particular consensus state.
func deleteConsensusMetadata(clientStore sdk.KVStore, height exported.Height) {
	deleteProcessedTime(clientStore, height)
	deleteIterationKey(clientStore, height)
}

// deleteProcessedTime deletes the processedTime for a given height
func deleteProcessedTime(clientStore sdk.KVStore, height exported.Height) {
	key := ProcessedTimeKey(height)
	clientStore.Delete(key)
}

// ProcessedTimeKey returns the key under which the processed time will be stored in the client store.
func ProcessedTimeKey(height exported.Height) []byte {
	return append(host.ConsensusStateKey(height), KeyProcessedTime...)
}

// deleteConsensusState deletes the consensus state at the given height
func deleteConsensusState(cdc codec.BinaryMarshaler, clientStore sdk.KVStore, height exported.Height) error {
	key := host.ConsensusStateKey(height)
	consensusState := clientStore.Get(key)
	if consensusState == nil {
		return sdkerrors.Wrapf(
			clienttypes.ErrConsensusStateNotFound,
			"consensus state does not exist for height %s", height,
		)
	}
	var currenttmp exported.ConsensusState
	if err := cdc.UnmarshalInterface(consensusState, &currenttmp); err != nil {
		return ErrUnmarshalInterface
	}
	tmpConsensus, ok := currenttmp.(*ConsensusState)
	if !ok {
		return ErrUnmarshalInterface
	}
	header := tmpConsensus.Header
	clientStore.Delete(ConsensusStateIndexKey(header.Hash()))
	clientStore.Delete(key)
	return nil
}

// deleteIterationKey deletes the iteration key for a given height
func deleteIterationKey(clientStore sdk.KVStore, height exported.Height) {
	key := IterationKey(height)
	clientStore.Delete(key)
}

// IterationKey returns the key under which the consensus state key will be stored.
// The iteration key is a BigEndian representation of the consensus state key to support efficient iteration.
func IterationKey(height exported.Height) []byte {
	heightBytes := bigEndianHeightBytes(height)
	return append([]byte(KeyIterateConsensusStatePrefix), heightBytes...)
}

func bigEndianHeightBytes(height exported.Height) []byte {
	heightBytes := make([]byte, 16)
	binary.BigEndian.PutUint64(heightBytes, height.GetRevisionNumber())
	binary.BigEndian.PutUint64(heightBytes[8:], height.GetRevisionHeight())
	return heightBytes
}
