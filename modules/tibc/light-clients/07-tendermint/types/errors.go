package types

import (
	errorsmod "cosmossdk.io/errors"
	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
)

const (
	SubModuleName = "tendermint-client"
	moduleName    = host.ModuleName + "-" + SubModuleName
)

// TIBC tendermint client sentinel errors
var (
	ErrInvalidChainID         = errorsmod.Register(moduleName, 2, "invalid chain-id")
	ErrInvalidTrustingPeriod  = errorsmod.Register(moduleName, 3, "invalid trusting period")
	ErrInvalidUnbondingPeriod = errorsmod.Register(moduleName, 4, "invalid unbonding period")
	ErrInvalidHeaderHeight    = errorsmod.Register(moduleName, 5, "invalid header height")
	ErrInvalidHeader          = errorsmod.Register(moduleName, 6, "invalid header")
	ErrInvalidMaxClockDrift   = errorsmod.Register(moduleName, 7, "invalid max clock drift")
	ErrProcessedTimeNotFound  = errorsmod.Register(moduleName, 8, "processed time not found")
	ErrDelayPeriodNotPassed   = errorsmod.Register(moduleName, 9, "packet-specified delay period has not been reached")
	ErrTrustingPeriodExpired  = errorsmod.Register(moduleName, 10, "time since latest trusted state has passed the trusting period")
	ErrUnbondingPeriodExpired = errorsmod.Register(moduleName, 11, "time since latest trusted state has passed the unbonding period")
	ErrInvalidProofSpecs      = errorsmod.Register(moduleName, 12, "invalid proof specs")
	ErrInvalidValidatorSet    = errorsmod.Register(moduleName, 13, "invalid validator set")
)
