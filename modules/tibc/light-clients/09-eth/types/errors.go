package types

import (
	errorsmod "cosmossdk.io/errors"
	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
)

const (
	SubModuleName = "eth-client"
	moduleName    = host.ModuleName + "-" + SubModuleName
)

// TIBC bsc client sentinel errors
var (
	ErrInvalidGenesisBlock = errorsmod.Register(moduleName, 2, "invalid genesis block")
	ErrFutureBlock         = errorsmod.Register(moduleName, 3, "block in the future")

	// ErrInvalidMixDigest is returned if a block's mix digest is non-zero.
	ErrInvalidMixDigest = errorsmod.Register(moduleName, 4, "non-zero mix digest")

	// ErrInvalidDifficulty is returned if the difficulty of a block is missing.
	ErrInvalidDifficulty = errorsmod.Register(moduleName, 5, "invalid difficulty")
	ErrUnknownAncestor   = errorsmod.Register(moduleName, 6, "unknown ancestor")

	// ErrWrongDifficulty is returned if the difficulty of a block doesn't match the
	// turn of the signer.
	ErrWrongDifficulty = errorsmod.Register(moduleName, 7, "wrong difficulty")
	ErrInvalidProof    = errorsmod.Register(moduleName, 8, "invalid proof")

	// ErrHeaderIsExist is returned if the header already in store
	ErrHeaderIsExist = errorsmod.Register(moduleName, 9, "header already exist")

	// ErrUnmarshalInterface is returned if interface can not unmarshal
	ErrUnmarshalInterface = errorsmod.Register(moduleName, 10, "unmarshal field")

	// ErrExtraLenth is returned if extra data is to long
	ErrExtraLenth = errorsmod.Register(moduleName, 11, "extra-data to long")

	// ErrInvalidGas is returned if gas data is invalid
	ErrInvalidGas = errorsmod.Register(moduleName, 12, "gas invalid")

	// ErrHeader is returned if the header invalid
	ErrHeader = errorsmod.Register(moduleName, 13, "header invalid")

	// ErrMarshalInterface is returned if struct can not marshal
	ErrMarshalInterface = errorsmod.Register(moduleName, 14, "marshal field")
)
