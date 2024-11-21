package types

import (
	errorsmod "cosmossdk.io/errors"
	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
)

const moduleName = host.ModuleName + "-" + SubModuleName

// TIBC routing sentinel errors
var (
	ErrInvalidRoute         = errorsmod.Register(moduleName, 2, "route not found")
	ErrInvalidRule          = errorsmod.Register(moduleName, 3, "invalid rule")
	ErrFailMarshalRules     = errorsmod.Register(moduleName, 4, "failed to marshal rules")
	ErrFailUnmarshalRules   = errorsmod.Register(moduleName, 5, "failed to unmarshal rules")
	ErrRoutingRulesNotFound = errorsmod.Register(moduleName, 6, "routing rules not found")
)
