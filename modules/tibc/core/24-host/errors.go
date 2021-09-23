package host

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// SubModuleName defines the TICS 24 host
const SubModuleName = "host"

const moduleName = ModuleName + "-" + SubModuleName

// TIBC client sentinel errors
var (
	ErrInvalidID     = sdkerrors.Register(moduleName, 2, "invalid identifier")
	ErrInvalidPath   = sdkerrors.Register(moduleName, 3, "invalid path")
	ErrInvalidPacket = sdkerrors.Register(moduleName, 4, "invalid packet")
	ErrInvalidRule   = sdkerrors.Register(moduleName, 5, "invalid routing rule")
)
