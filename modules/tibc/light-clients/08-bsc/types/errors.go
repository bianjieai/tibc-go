package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
)

const (
	SubModuleName = "bsc-client"
	moduleName    = host.ModuleName + "-" + SubModuleName
)

// TIBC bsc client sentinel errors
var (
	ErrInvalidGenesisBlock   = sdkerrors.Register(moduleName, 2, "invalid genesis block")
	ErrInvalidValidatorBytes = sdkerrors.Register(moduleName, 3, "invalid validators bytes length")

	// ErrUnknownBlock is returned when the list of validators is requested for a block
	// that is not part of the local blockchain.
	ErrUnknownBlock = sdkerrors.Register(moduleName, 4, "unknown block")
	ErrFutureBlock  = sdkerrors.Register(moduleName, 5, "block in the future")

	// ErrMissingVanity is returned if a block's extra-data section is shorter than
	// 32 bytes, which is required to store the signer vanity.
	ErrMissingVanity = sdkerrors.Register(moduleName, 6, "extra-data 32 byte vanity prefix missing")

	// ErrMissingSignature is returned if a block's extra-data section doesn't seem
	// to contain a 65 byte secp256k1 signature.
	ErrMissingSignature = sdkerrors.Register(moduleName, 7, "extra-data 65 byte signature suffix missing")

	// ErrInvalidMixDigest is returned if a block's mix digest is non-zero.
	ErrInvalidMixDigest = sdkerrors.Register(moduleName, 8, "non-zero mix digest")

	// ErrInvalidUncleHash is returned if a block contains an non-empty uncle list.
	ErrInvalidUncleHash = sdkerrors.Register(moduleName, 9, "non empty uncle hash")

	// ErrInvalidDifficulty is returned if the difficulty of a block is missing.
	ErrInvalidDifficulty = sdkerrors.Register(moduleName, 10, "invalid difficulty")
	ErrUnknownAncestor   = sdkerrors.Register(moduleName, 11, "unknown ancestor")
	// ErrCoinBaseMisMatch is returned if a header's coinbase do not match with signature
	ErrCoinBaseMisMatch = sdkerrors.Register(moduleName, 12, "coinbase do not match with signature")
	// ErrUnauthorizedValidator is returned if a header is signed by a non-authorized entity.
	ErrUnauthorizedValidator = sdkerrors.Register(moduleName, 13, "unauthorized validator")
	// ErrRecentlySigned is returned if a header is signed by an authorized entity
	// that already signed a header recently, thus is temporarily not allowed to.
	ErrRecentlySigned = sdkerrors.Register(moduleName, 14, "recently signed")
	// ErrWrongDifficulty is returned if the difficulty of a block doesn't match the
	// turn of the signer.
	ErrWrongDifficulty = sdkerrors.Register(moduleName, 15, "wrong difficulty")
	// ErrExtraValidators is returned if non-sprint-end block contain validator data in
	// their extra-data fields.
	ErrExtraValidators = sdkerrors.Register(moduleName, 16, "non-sprint-end block contains extra validator list")
	// ErrInvalidSpanValidators is returned if a block contains an
	// invalid list of validators (i.e. non divisible by 20 bytes).
	ErrInvalidSpanValidators = sdkerrors.Register(moduleName, 17, "invalid validator list on sprint end block")

	ErrInvalidProof = sdkerrors.Register(moduleName, 18, "invalid proof")
)
