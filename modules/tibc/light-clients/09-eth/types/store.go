package types

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/store/types"

	"github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

var (
	KeyIndexEthHeaderPrefix = "ethHeaderIndex"
	KeyMainRootPrefix       = "ethRootMain"
)

func EthHeaderIndexPath(hash common.Hash, height uint64) string {
	return fmt.Sprintf("%s/%s%d", KeyIndexEthHeaderPrefix, hash, height)
}

func EthHeaderIndexKey(hash common.Hash, height uint64) []byte {
	return []byte(EthHeaderIndexPath(hash, height))
}
func EthRootMainPath(root common.Hash, height uint64) string {
	return fmt.Sprintf("%s/%s%d", KeyMainRootPrefix, root, height)
}
func EthRootMainKey(root common.Hash, height uint64) []byte {
	return []byte(EthRootMainPath(root, height))
}

func SetEthHeaderIndex(
	clientStore sdk.KVStore,
	header Header,
	headerBytes []byte,
) {
	clientStore.Set(EthHeaderIndexKey(header.Hash(), header.Height.RevisionHeight), headerBytes)
}
func GetIterator(store sdk.KVStore, keyType string) types.Iterator {
	iterator := sdk.KVStorePrefixIterator(store, []byte(keyType))
	return iterator
}
func IteratorEthMetaDataByPrefix(store sdk.KVStore, keyType string, cb func(key, val []byte) bool) {
	iterator := GetIterator(store, keyType)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		if cb(iterator.Key(), iterator.Value()) {
			break
		}
	}
}

func GetParentHeaderFromIndex(
	clientStore sdk.KVStore,
	header Header,
) []byte {
	get := clientStore.Get(EthHeaderIndexKey(header.ToEthHeader().ParentHash, header.Height.RevisionHeight-1))
	return get
}

// SetEthConsensusRoot sets the consensus metadata with the provided values
func SetEthConsensusRoot(
	clientStore sdk.KVStore,
	height uint64,
	root, headerHash common.Hash,
) {
	clientStore.Set(EthRootMainKey(root, height), EthHeaderIndexKey(headerHash, height))
}
func GetHeaderIndexKeyByEthConsensusRoot(clientStore sdk.KVStore, root common.Hash, height uint64) []byte {
	return clientStore.Get(EthRootMainKey(root, height))
}

// GetConsensusState retrieves the consensus state from the client prefixed
// store. An error is returned if the consensus state does not exist.
func GetConsensusState(store sdk.KVStore, cdc codec.BinaryCodec, height exported.Height) (*ConsensusState, error) {
	bz := store.Get(host.ConsensusStateKey(height))
	if bz == nil {
		return nil, sdkerrors.Wrapf(
			clienttypes.ErrConsensusStateNotFound,
			"consensus state does not exist for height %s",
			height,
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
			"invalid consensus type %T, expected %T",
			consensusState, &ConsensusState{},
		)
	}

	return consensusState, nil
}

// IterateConsensusStateAscending iterates through the consensus states in ascending order. It calls the provided
// callback on each height, until stop=true is returned.
func IterateConsensusStateAscending(clientStore sdk.KVStore,
	cb func(height exported.Height) (stop bool)) {
	iterator := sdk.KVStorePrefixIterator(clientStore, []byte(host.KeyConsensusStatePrefix))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		keySplit := strings.Split(string(key), "/")
		// processed time key in prefix store has format: "consensusStates/<height>"
		if len(keySplit) != 2 {
			// ignore all not consensus state keys
			continue
		}
		height := GetHeightFromIterationKey(key)
		if cb(height) {
			return
		}
	}
}

// GetHeightFromIterationKey takes an iteration key and returns the height that it references
func GetHeightFromIterationKey(iterKey []byte) exported.Height {
	bigEndianBytes := iterKey[len([]byte(host.KeyConsensusStatePrefix+"/")):]
	revisionBytes := bigEndianBytes[0:8]
	heightBytes := bigEndianBytes[8:]
	revision := sdk.BigEndianToUint64(revisionBytes)
	height := sdk.BigEndianToUint64(heightBytes)
	return clienttypes.NewHeight(revision, height)
}
func deleteRootMain(clientStore sdk.KVStore, root common.Hash, height exported.Height) {
	clientStore.Delete(EthRootMainKey(root, height.GetRevisionHeight()))
}

// deleteConsensusStateAndIndexHeader deletes the consensus state at the given height
func deleteConsensusStateAndIndexHeader(cdc codec.BinaryCodec, clientStore sdk.KVStore, height exported.Height) error {
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
	root := tmpConsensus.Root
	headerIndexKey := GetHeaderIndexKeyByEthConsensusRoot(clientStore, common.BytesToHash(root), height.GetRevisionHeight())
	if headerIndexKey == nil {
		return sdkerrors.Wrapf(ErrHeader,
			"Header index not found in height %s", height,
		)
	}
	// delete header index
	clientStore.Delete(headerIndexKey)
	// delete root main
	deleteRootMain(clientStore, common.BytesToHash(root), height)
	// delete consensus state
	clientStore.Delete(key)
	return nil
}

func bigEndianHeightBytes(height exported.Height) []byte {
	heightBytes := make([]byte, 16)
	binary.BigEndian.PutUint64(heightBytes, height.GetRevisionNumber())
	binary.BigEndian.PutUint64(heightBytes[8:], height.GetRevisionHeight())
	return heightBytes
}
