package keeper

import "github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/types"

// MustMarshalClassTrace attempts to encode an ClassTrace object and returns the
// raw encoded bytes. It panics on error.
func (k Keeper) MustMarshalClassTrace(classTrace types.ClassTrace) []byte {
	return k.cdc.MustMarshal(&classTrace)
}

// UnmarshalClassTrace attempts to decode and return an ClassTrace object from
// raw encoded bytes.
func (k Keeper) UnmarshalClassTrace(bz []byte) (types.ClassTrace, error) {
	var classTrace types.ClassTrace
	if err := k.cdc.Unmarshal(bz, &classTrace); err != nil {
		return types.ClassTrace{}, err
	}

	return classTrace, nil
}
