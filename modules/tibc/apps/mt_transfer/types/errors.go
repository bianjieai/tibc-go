package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrInvalidDenom            = errorsmod.Register(ModuleName, 2, "invalid denom")
	ErrUnknownNFT              = errorsmod.Register(ModuleName, 3, "unknown mt")
	ErrScChainEqualToDestChain = errorsmod.Register(ModuleName, 4, "source chain equals to destination chain")
	ErrTraceNotFound           = errorsmod.Register(ModuleName, 5, "class trace not found")
	ErrInvalidAmount           = errorsmod.Register(ModuleName, 6, "the amount must be greater than 0")
)
