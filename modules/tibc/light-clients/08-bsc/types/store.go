package types

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

var (
	PrefixKeyValidators    = "validators"
	PrefixKeyRecentSingers = "RecentSingers"
)

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

// SetRecentSingers sets the recent singer list in the client prefixed store
func SetRecentSingers(store sdk.KVStore, recentSingers []Singer) {
	for _, singer := range recentSingers {
		store.Set(keyRecentSinger(singer), singer.Validator)
	}
}

// GetRecentSingers retrieves the recent singer list from the client prefixed
func GetRecentSingers(store sdk.KVStore) (recentSingers []Singer, err error) {
	iterator := sdk.KVStorePrefixIterator(store, []byte(PrefixKeyRecentSingers))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		iterKey := iterator.Key()
		keys := strings.Split(string(iterKey), "/")
		height, err := clienttypes.ParseHeight(keys[1])
		if err != nil {
			return nil, err
		}
		recentSingers = append(recentSingers, Singer{
			Height:    height,
			Validator: iterator.Value(),
		})
	}
	return
}

// GetValidators retrieves the validators from the client prefixed store
func GetValidators(
	cdc codec.BinaryMarshaler,
	store sdk.KVStore,
	height exported.Height,
) ValidatorSet {
	bz := store.Get(keyValidators(height))

	var validatorSet ValidatorSet
	cdc.MustUnmarshalBinaryBare(bz, &validatorSet)
	return validatorSet
}

// SetValidators sets the validators in the client prefixed store
func SetValidators(store sdk.KVStore,
	cdc codec.BinaryMarshaler,
	height exported.Height,
	validators [][]byte,
) {
	validatorSet := ValidatorSet{
		Validators: validators,
	}
	bz := cdc.MustMarshalBinaryBare(&validatorSet)
	store.Set(keyValidators(height), bz)
}

func keyRecentSinger(singer Singer) []byte {
	return []byte(fmt.Sprintf("%s/%s", PrefixKeyValidators, singer.Height))
}

func keyValidators(height exported.Height) []byte {
	return []byte(fmt.Sprintf("%s/%s", PrefixKeyRecentSingers, height))
}
