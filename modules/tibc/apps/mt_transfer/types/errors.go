package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidDenom            = sdkerrors.Register(ModuleName, 2, "invalid denom")
	ErrUnknownNFT              = sdkerrors.Register(ModuleName, 3, "unknown mt")
	ErrScChainEqualToDestChain = sdkerrors.Register(ModuleName, 4, "source chain equals to destination chain")
	ErrTraceNotFound           = sdkerrors.Register(ModuleName, 5, "class trace not found")
	ErrInvalidAmount           = sdkerrors.Register(ModuleName, 6, "the amount must be greater than 0")
)
