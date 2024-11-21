package host

import (
	errorsmod "cosmossdk.io/errors"
)

// SubModuleName defines the TICS 24 host
const SubModuleName = "host"

const moduleName = ModuleName + "-" + SubModuleName

// TIBC client sentinel errors
var (
	ErrInvalidID     = errorsmod.Register(moduleName, 2, "invalid identifier")
	ErrInvalidPath   = errorsmod.Register(moduleName, 3, "invalid path")
	ErrInvalidPacket = errorsmod.Register(moduleName, 4, "invalid packet")
	ErrInvalidRule   = errorsmod.Register(moduleName, 5, "invalid routing rule")
)
