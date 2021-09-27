package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
)

const (
	SubModuleName = "eth-client"
	moduleName    = host.ModuleName + "-" + SubModuleName
)

// TIBC bsc client sentinel errors
var (
	ErrInvalidGenesisBlock = sdkerrors.Register(moduleName, 2, "invalid genesis block")
	ErrFutureBlock         = sdkerrors.Register(moduleName, 3, "block in the future")

	// ErrInvalidMixDigest is returned if a block's mix digest is non-zero.
	ErrInvalidMixDigest = sdkerrors.Register(moduleName, 4, "non-zero mix digest")

	// ErrInvalidDifficulty is returned if the difficulty of a block is missing.
	ErrInvalidDifficulty = sdkerrors.Register(moduleName, 5, "invalid difficulty")
	ErrUnknownAncestor   = sdkerrors.Register(moduleName, 6, "unknown ancestor")

	// ErrWrongDifficulty is returned if the difficulty of a block doesn't match the
	// turn of the signer.
	ErrWrongDifficulty = sdkerrors.Register(moduleName, 7, "wrong difficulty")
	ErrInvalidProof    = sdkerrors.Register(moduleName, 8, "invalid proof")

	// ErrHeaderIsExist is returned if the header already in store
	ErrHeaderIsExist = sdkerrors.Register(moduleName, 9, "header already exist")

	// ErrUnmarshalInterface is returned if interface can not unmarshal
	ErrUnmarshalInterface = sdkerrors.Register(moduleName, 10, "unmarshal field")

	// ErrExtraLenth is returned if extra data is to long
	ErrExtraLenth = sdkerrors.Register(moduleName, 11, "extra-data to long")

	// ErrInvalidGas is returned if gas data is invalid
	ErrInvalidGas = sdkerrors.Register(moduleName, 12, "gas invalid")

	// ErrHeader is returned if the header invalid
	ErrHeader = sdkerrors.Register(moduleName, 13, "header invalid")

	// ErrMarshalInterface is returned if struct can not marshal
	ErrMarshalInterface = sdkerrors.Register(moduleName, 14, "marshal field")
)
