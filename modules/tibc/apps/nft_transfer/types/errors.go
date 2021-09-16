package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrInvalidDenom            = sdkerrors.Register(ModuleName, 2, "invalid denom")
	ErrUnknownNFT              = sdkerrors.Register(ModuleName, 3, "unknown nft")
	ErrScChainEqualToDestChain = sdkerrors.Register(ModuleName, 4, "scChain equals to destChain")
)
