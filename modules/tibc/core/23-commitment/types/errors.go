package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
)

// SubModuleName is the error codespace
const SubModuleName string = "commitment"

const moduleName = host.ModuleName + "-" + SubModuleName

// TIBC connection sentinel errors
var (
	ErrInvalidProof       = sdkerrors.Register(moduleName, 2, "invalid proof")
	ErrInvalidPrefix      = sdkerrors.Register(moduleName, 3, "invalid prefix")
	ErrInvalidMerkleProof = sdkerrors.Register(moduleName, 4, "invalid merkle proof")
)
