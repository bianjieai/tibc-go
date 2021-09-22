package types

import (
	"fmt"

	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
)

// TIBC client events
const (
	AttributeKeyChainName       = "chain_name"
	AttributeKeyClientType      = "client_type"
	AttributeKeyConsensusHeight = "consensus_height"
	AttributeKeyHeader          = "header"
)

// TIBC client events vars
var (
	EventTypeCreateClientProposal  = "create_client_proposal"
	EventTypeUpdateClient          = "update_client"
	EventTypeUpgradeClientProposal = "upgrade_client_proposal"
	EventTypeUpdateClientProposal  = "update_client_proposal"

	AttributeValueCategory = fmt.Sprintf("%s_%s", host.ModuleName, SubModuleName)
)
