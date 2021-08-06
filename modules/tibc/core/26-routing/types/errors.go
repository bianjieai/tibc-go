package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
)

const moduleName = host.ModuleName + "-" + SubModuleName

// TIBC routing sentinel errors
var (
	ErrInvalidRoute = sdkerrors.Register(moduleName, 2, "route not found")
)
