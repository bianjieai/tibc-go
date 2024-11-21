package types

import (
	errorsmod "cosmossdk.io/errors"
	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
)

const moduleName = host.ModuleName + "-" + SubModuleName

// TIBC client sentinel errors
var (
	ErrClientExists                           = errorsmod.Register(moduleName, 2, "light client already exists")
	ErrInvalidClient                          = errorsmod.Register(moduleName, 3, "light client is invalid")
	ErrClientNotFound                         = errorsmod.Register(moduleName, 4, "light client not found")
	ErrClientFrozen                           = errorsmod.Register(moduleName, 5, "light client is frozen due to misbehaviour")
	ErrInvalidClientMetadata                  = errorsmod.Register(moduleName, 6, "invalid client metadata")
	ErrConsensusStateNotFound                 = errorsmod.Register(moduleName, 7, "consensus state not found")
	ErrInvalidConsensus                       = errorsmod.Register(moduleName, 8, "invalid consensus state")
	ErrClientTypeNotFound                     = errorsmod.Register(moduleName, 9, "client type not found")
	ErrInvalidClientType                      = errorsmod.Register(moduleName, 10, "invalid client type")
	ErrRootNotFound                           = errorsmod.Register(moduleName, 11, "commitment root not found")
	ErrInvalidHeader                          = errorsmod.Register(moduleName, 12, "invalid client header")
	ErrInvalidMisbehaviour                    = errorsmod.Register(moduleName, 13, "invalid light client misbehaviour")
	ErrFailedClientStateVerification          = errorsmod.Register(moduleName, 14, "client state verification failed")
	ErrFailedClientConsensusStateVerification = errorsmod.Register(moduleName, 15, "client consensus state verification failed")
	ErrFailedConnectionStateVerification      = errorsmod.Register(moduleName, 16, "connection state verification failed")
	ErrFailedChannelStateVerification         = errorsmod.Register(moduleName, 17, "channel state verification failed")
	ErrFailedPacketCommitmentVerification     = errorsmod.Register(moduleName, 18, "packet commitment verification failed")
	ErrFailedPacketAckVerification            = errorsmod.Register(moduleName, 19, "packet acknowledgement verification failed")
	ErrFailedPacketReceiptVerification        = errorsmod.Register(moduleName, 20, "packet receipt verification failed")
	ErrFailedNextSeqRecvVerification          = errorsmod.Register(moduleName, 21, "next sequence receive verification failed")
	ErrSelfConsensusStateNotFound             = errorsmod.Register(moduleName, 22, "self consensus state not found")
	ErrUpdateClientFailed                     = errorsmod.Register(moduleName, 23, "unable to update light client")
	ErrInvalidUpdateClientProposal            = errorsmod.Register(moduleName, 24, "invalid update client proposal")
	ErrInvalidUpgradeClient                   = errorsmod.Register(moduleName, 25, "invalid client upgrade")
	ErrRelayerExists                          = errorsmod.Register(moduleName, 26, "relayer already exists")
	ErrClientNotActive                        = errorsmod.Register(moduleName, 27, "client is not active")
)
