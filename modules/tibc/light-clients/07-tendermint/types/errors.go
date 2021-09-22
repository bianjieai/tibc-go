package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
)

const (
	SubModuleName = "tendermint-client"
	moduleName    = host.ModuleName + "-" + SubModuleName
)

// TIBC tendermint client sentinel errors
var (
	ErrInvalidChainID         = sdkerrors.Register(moduleName, 2, "invalid chain-id")
	ErrInvalidTrustingPeriod  = sdkerrors.Register(moduleName, 3, "invalid trusting period")
	ErrInvalidUnbondingPeriod = sdkerrors.Register(moduleName, 4, "invalid unbonding period")
	ErrInvalidHeaderHeight    = sdkerrors.Register(moduleName, 5, "invalid header height")
	ErrInvalidHeader          = sdkerrors.Register(moduleName, 6, "invalid header")
	ErrInvalidMaxClockDrift   = sdkerrors.Register(moduleName, 7, "invalid max clock drift")
	ErrProcessedTimeNotFound  = sdkerrors.Register(moduleName, 8, "processed time not found")
	ErrDelayPeriodNotPassed   = sdkerrors.Register(moduleName, 9, "packet-specified delay period has not been reached")
	ErrTrustingPeriodExpired  = sdkerrors.Register(moduleName, 10, "time since latest trusted state has passed the trusting period")
	ErrUnbondingPeriodExpired = sdkerrors.Register(moduleName, 11, "time since latest trusted state has passed the unbonding period")
	ErrInvalidProofSpecs      = sdkerrors.Register(moduleName, 12, "invalid proof specs")
	ErrInvalidValidatorSet    = sdkerrors.Register(moduleName, 13, "invalid validator set")
)
