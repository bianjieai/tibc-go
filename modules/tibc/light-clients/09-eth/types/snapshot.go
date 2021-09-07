package types

import (
	"bytes"

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
	cdc   codec.BinaryMarshaler
	store sdk.KVStore
	// Block number where the snapshot was created
	Number uint64
	// Set of authorized signers at this moment
	Signers map[common.Address]struct{} `json:"signers"`
	// Set of authorized validators at this moment
	Validators map[common.Address]struct{}
	Recents    map[uint64]common.Address
}

func (m ClientState) snapshot(
	cdc codec.BinaryMarshaler,
	store sdk.KVStore,
) (*snapshot, error) {

	snap := &snapshot{
		cdc:    cdc,
		store:  store,
		Number: m.Header.Height.RevisionHeight,
	}

	return snap, nil
}
