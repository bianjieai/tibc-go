package types

import (
	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const moduleName = host.ModuleName + "-" + SubModuleName

// IBC port sentinel errors
var (
	ErrPortExists   = sdkerrors.Register(moduleName, 2, "port is already binded")
	ErrPortNotFound = sdkerrors.Register(moduleName, 3, "port not found")
	ErrInvalidPort  = sdkerrors.Register(moduleName, 4, "invalid port")
	ErrInvalidRoute = sdkerrors.Register(moduleName, 5, "route not found")
)
