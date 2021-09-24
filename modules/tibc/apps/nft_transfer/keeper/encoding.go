package keeper

import "github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/types"

// MustMarshalClassTrace attempts to encode an ClassTrace object and returns the
// raw encoded bytes. It panics on error.
func (k Keeper) MustMarshalClassTrace(classTrace types.ClassTrace) []byte {
	return k.cdc.MustMarshalBinaryBare(&classTrace)
}
