package types

import (
	storetypes "cosmossdk.io/store/types"
	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

// ExportMetadata exports all the processed times in the client store so they can be included in clients genesis
// and imported by a ClientKeeper
func (cs ClientState) ExportMetadata(store storetypes.KVStore) []exported.GenesisMetadata {
	gm := make([]exported.GenesisMetadata, 0)
	IterateProcessedTime(store, func(key, val []byte) bool {
		gm = append(gm, clienttypes.NewGenesisMetadata(key, val))
		return false
	})
	if len(gm) == 0 {
		return nil
	}
	return gm
}
