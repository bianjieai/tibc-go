package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
)

const (
	SubModuleName = "bsc-client"
	moduleName    = host.ModuleName + "-" + SubModuleName
)

// IBC bsc client sentinel errors
var (
	ErrInvalidGenesisBlock   = sdkerrors.Register(moduleName, 2, "invalid chain-id")
	ErrInvalidValidatorBytes = sdkerrors.Register(moduleName, 3, "invalid validators bytes length")
)
