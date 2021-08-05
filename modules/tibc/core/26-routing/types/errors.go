package types

import (
	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const moduleName = host.ModuleName + "-" + SubModuleName

// TIBC routing sentinel errors
var (
	ErrInvalidRoute = sdkerrors.Register(moduleName, 2, "route not found")
)
