package types

import (
	errorsmod "cosmossdk.io/errors"
	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
)

// SubModuleName is the error codespace
const SubModuleName string = "commitment"

const moduleName = host.ModuleName + "-" + SubModuleName

// TIBC connection sentinel errors
var (
	ErrInvalidProof       = errorsmod.Register(moduleName, 2, "invalid proof")
	ErrInvalidPrefix      = errorsmod.Register(moduleName, 3, "invalid prefix")
	ErrInvalidMerkleProof = errorsmod.Register(moduleName, 4, "invalid merkle proof")
)
