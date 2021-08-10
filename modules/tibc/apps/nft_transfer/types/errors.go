package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var(
	ErrInvalidDenom      		 = sdkerrors.Register(ModuleName, 1, "invalid denom")
)

