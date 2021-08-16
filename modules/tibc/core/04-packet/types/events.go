package types

import (
	"fmt"

	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
)

// IBC channel events
const (
	AttributeKeyConnectionID       = "connection_id"
	AttributeKeyPortID             = "port_id"
	AttributeKeyChannelID          = "channel_id"
	AttributeCounterpartyPortID    = "counterparty_port_id"
	AttributeCounterpartyChannelID = "counterparty_channel_id"

	EventTypeSendPacket        = "send_packet"
	EventTypeRecvPacket        = "recv_packet"
	EventTypeWriteAck          = "write_acknowledgement"
	EventTypeAcknowledgePacket = "acknowledge_packet"
	EventTypeTimeoutPacket     = "timeout_packet"
	EventTypeSendCleanPacket   = "send_clean_packet"
	EventTypeRecvCleanPacket   = "recv_clean_packet"

	AttributeKeyData             = "packet_data"
	AttributeKeyAck              = "packet_ack"
	AttributeKeyTimeoutHeight    = "packet_timeout_height"
	AttributeKeyTimeoutTimestamp = "packet_timeout_timestamp"
	AttributeKeySequence         = "packet_sequence"
	AttributeKeyPort             = "packet_port"
	AttributeKeySrcChain         = "packet_src_chain"
	AttributeKeyDstChain         = "packet_dst_port"
	AttributeKeyRelayChain       = "packet_relay_channel"
	AttributeKeyConnection       = "packet_connection"
)

// IBC channel events vars
var (
	AttributeValueCategory = fmt.Sprintf("%s_%s", host.ModuleName, SubModuleName)
)
