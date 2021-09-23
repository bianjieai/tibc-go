package types

import (
	"bytes"
	"sort"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

// validatorsAscending implements the sort interface to allow sorting a list of addresses
type validatorsAscending []common.Address

func (s validatorsAscending) Len() int           { return len(s) }
func (s validatorsAscending) Less(i, j int) bool { return bytes.Compare(s[i][:], s[j][:]) < 0 }
func (s validatorsAscending) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type snapshot struct {
	cdc   codec.BinaryCodec
	store sdk.KVStore
	// Block number where the snapshot was created
	Number uint64
	// Set of authorized validators at this moment
	Validators map[common.Address]struct{}
	Recents    map[uint64]common.Address
}

func (m ClientState) snapshot(cdc codec.BinaryCodec, store sdk.KVStore) (*snapshot, error) {

	recentSingers, err := GetRecentSigners(store)
	if err != nil {
		return nil, err
	}

	snap := &snapshot{
		cdc:        cdc,
		store:      store,
		Number:     m.Header.Height.RevisionHeight,
		Validators: make(map[common.Address]struct{}, len(m.Validators)),
		Recents:    make(map[uint64]common.Address, len(recentSingers)),
	}

	for _, validator := range m.Validators {
		snap.Validators[common.BytesToAddress(validator)] = struct{}{}
	}

	for _, singer := range recentSingers {
		snap.Recents[singer.Height.RevisionHeight] = common.BytesToAddress(singer.Validator)
	}
	return snap, nil
}

// validators retrieves the list of validators in ascending order.
func (s *snapshot) validators() []common.Address {
	validators := make([]common.Address, 0, len(s.Validators))
	for v := range s.Validators {
		validators = append(validators, v)
	}
	sort.Sort(validatorsAscending(validators))
	return validators
}

// inturn returns if a validator at a given block height is in-turn or not.
func (s *snapshot) inturn(validator common.Address) bool {
	validators := s.validators()
	offset := (s.Number + 1) % uint64(len(validators))
	return validators[offset] == validator
}
