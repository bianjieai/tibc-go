package types

import (
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/store/types"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"

	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

var (
	PrefixKeyRecentSingers  = "recentSingers"
	PrefixPendingValidators = "pendingValidators"
)

func GetIterator(store storetypes.KVStore, keyType string) types.Iterator {
	iterator := storetypes.KVStorePrefixIterator(store, []byte(keyType))
	return iterator
}

func IteratorTraversal(store storetypes.KVStore, keyType string, cb func(key, val []byte) bool) {
	iterator := GetIterator(store, keyType)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		if cb(iterator.Key(), iterator.Value()) {
			break
		}
	}
}

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

// SetRecentSigners sets the recent singer list in the client prefixed store
func SetRecentSigners(store storetypes.KVStore, recentSingers []Signer) {
	for _, singer := range recentSingers {
		store.Set(keyRecentSinger(singer), singer.Validator)
	}
}

func SetSigner(store storetypes.KVStore, signer Signer) {
	store.Set(keyRecentSinger(signer), signer.Validator)
}

func DeleteSigner(store storetypes.KVStore, height clienttypes.Height) {
	keyBz := []byte(fmt.Sprintf("%s/%s", PrefixKeyRecentSingers, height))
	store.Delete(keyBz)
}

// GetRecentSigners retrieves the recent singer list from the client prefixed
func GetRecentSigners(store storetypes.KVStore) (recentSingers []Signer, err error) {
	iterator := storetypes.KVStorePrefixIterator(store, []byte(PrefixKeyRecentSingers))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		iterKey := iterator.Key()
		keys := strings.Split(string(iterKey), "/")
		height, err := clienttypes.ParseHeight(keys[1])
		if err != nil {
			return nil, err
		}
		recentSingers = append(recentSingers, Signer{
			Height:    height,
			Validator: iterator.Value(),
		})
	}
	return
}

// SetPendingValidators sets the validators to be updated in the client prefixed store
func SetPendingValidators(store storetypes.KVStore, cdc codec.BinaryCodec, validators [][]byte,
) {
	validatorSet := ValidatorSet{
		Validators: validators,
	}
	bz := cdc.MustMarshal(&validatorSet)
	store.Set([]byte(PrefixPendingValidators), bz)
}

// GetPendingValidators retrieves the validators to be updated from the client prefixed store
func GetPendingValidators(cdc codec.BinaryCodec, store storetypes.KVStore) ValidatorSet {
	bz := store.Get([]byte(PrefixPendingValidators))

	var validatorSet ValidatorSet
	cdc.MustUnmarshal(bz, &validatorSet)
	return validatorSet
}

func keyRecentSinger(singer Signer) []byte {
	return []byte(fmt.Sprintf("%s/%s", PrefixKeyRecentSingers, singer.Height))
}
