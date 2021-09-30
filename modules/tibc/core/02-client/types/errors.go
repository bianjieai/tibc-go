package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
)

const moduleName = host.ModuleName + "-" + SubModuleName

// TIBC client sentinel errors
var (
	ErrClientExists                           = sdkerrors.Register(moduleName, 2, "light client already exists")
	ErrInvalidClient                          = sdkerrors.Register(moduleName, 3, "light client is invalid")
	ErrClientNotFound                         = sdkerrors.Register(moduleName, 4, "light client not found")
	ErrClientFrozen                           = sdkerrors.Register(moduleName, 5, "light client is frozen due to misbehaviour")
	ErrInvalidClientMetadata                  = sdkerrors.Register(moduleName, 6, "invalid client metadata")
	ErrConsensusStateNotFound                 = sdkerrors.Register(moduleName, 7, "consensus state not found")
	ErrInvalidConsensus                       = sdkerrors.Register(moduleName, 8, "invalid consensus state")
	ErrClientTypeNotFound                     = sdkerrors.Register(moduleName, 9, "client type not found")
	ErrInvalidClientType                      = sdkerrors.Register(moduleName, 10, "invalid client type")
	ErrRootNotFound                           = sdkerrors.Register(moduleName, 11, "commitment root not found")
	ErrInvalidHeader                          = sdkerrors.Register(moduleName, 12, "invalid client header")
	ErrInvalidMisbehaviour                    = sdkerrors.Register(moduleName, 13, "invalid light client misbehaviour")
	ErrFailedClientStateVerification          = sdkerrors.Register(moduleName, 14, "client state verification failed")
	ErrFailedClientConsensusStateVerification = sdkerrors.Register(moduleName, 15, "client consensus state verification failed")
	ErrFailedConnectionStateVerification      = sdkerrors.Register(moduleName, 16, "connection state verification failed")
	ErrFailedChannelStateVerification         = sdkerrors.Register(moduleName, 17, "channel state verification failed")
	ErrFailedPacketCommitmentVerification     = sdkerrors.Register(moduleName, 18, "packet commitment verification failed")
	ErrFailedPacketAckVerification            = sdkerrors.Register(moduleName, 19, "packet acknowledgement verification failed")
	ErrFailedPacketReceiptVerification        = sdkerrors.Register(moduleName, 20, "packet receipt verification failed")
	ErrFailedNextSeqRecvVerification          = sdkerrors.Register(moduleName, 21, "next sequence receive verification failed")
	ErrSelfConsensusStateNotFound             = sdkerrors.Register(moduleName, 22, "self consensus state not found")
	ErrUpdateClientFailed                     = sdkerrors.Register(moduleName, 23, "unable to update light client")
	ErrInvalidUpdateClientProposal            = sdkerrors.Register(moduleName, 24, "invalid update client proposal")
	ErrInvalidUpgradeClient                   = sdkerrors.Register(moduleName, 25, "invalid client upgrade")
	ErrRelayerExists                          = sdkerrors.Register(moduleName, 26, "relayer already exists")
	ErrClientNotActive                        = sdkerrors.Register(moduleName, 27, "client is not active")
)
