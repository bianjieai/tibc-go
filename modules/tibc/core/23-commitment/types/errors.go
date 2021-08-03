package types

import (
	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// SubModuleName is the error codespace
const SubModuleName string = "commitment"

const moduleName = host.ModuleName + "-" + SubModuleName

// IBC connection sentinel errors
var (
	ErrInvalidProof       = sdkerrors.Register(moduleName, 2, "invalid proof")
	ErrInvalidPrefix      = sdkerrors.Register(moduleName, 3, "invalid prefix")
	ErrInvalidMerkleProof = sdkerrors.Register(moduleName, 4, "invalid merkle proof")
)
