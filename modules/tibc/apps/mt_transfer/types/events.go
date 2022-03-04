package types

// TIBC transfer events
const (
	EventTypePacket      = "multi_token_packet"
	EventTypeNftTransfer = "tibc_mt_transfer"
	EventTypeClassTrace  = "class_trace"

	AttributeKeyClass      = "class"
	AttributeKeyId         = "id"
	AttributeKeyData       = "data"
	AttributeKeyReceiver   = "receiver"
	AttributeKeyAck        = "ack"
	AttributeKeyAckSuccess = "success"
	AttributeKeyAckError   = "error"
	AttributeKeyTraceHash  = "trace_hash"
)
