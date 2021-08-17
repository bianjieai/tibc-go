package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
)

const moduleName = host.ModuleName + "-" + SubModuleName

// TIBC routing sentinel errors
var (
	ErrInvalidRoute     = sdkerrors.Register(moduleName, 2, "route not found")
	ErrInvalidRule      = sdkerrors.Register(moduleName, 3, "invalid rule")
	ErrFailMarshalRules = sdkerrors.Register(moduleName, 4, "ailed to marshal rules")
	ErrRoutingRulesNotFound      = sdkerrors.Register(moduleName, 22, "routing rules not found")
)
