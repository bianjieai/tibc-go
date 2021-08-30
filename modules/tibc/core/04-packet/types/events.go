package types

import (
	"fmt"

	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
)

// TIBC packet events
const (
	EventTypeSendPacket        = "send_packet"
	EventTypeRecvPacket        = "recv_packet"
	EventTypeWriteAck          = "write_acknowledgement"
	EventTypeAcknowledgePacket = "acknowledge_packet"
	EventTypeSendCleanPacket   = "send_clean_packet"
	EventTypeRecvCleanPacket   = "recv_clean_packet"

	AttributeKeyData       = "packet_data"
	AttributeKeyAck        = "packet_ack"
	AttributeKeySequence   = "packet_sequence"
	AttributeKeyPort       = "packet_port"
	AttributeKeySrcChain   = "packet_src_chain"
	AttributeKeyDstChain   = "packet_dst_port"
	AttributeKeyRelayChain = "packet_relay_channel"
)

// tibc packet events vars
var (
	AttributeValueCategory = fmt.Sprintf("%s_%s", host.ModuleName, SubModuleName)
)
