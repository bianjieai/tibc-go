package types

// TIBC transfer events
const (
	EventTypePacket      = "non_fungible_token_packet"
	EventTypeNftTransfer = "tibc_nft_transfer"
	EventTypeClassTrace  = "class_trace"

	AttributeKeyClass      = "class"
	AttributeKeyId         = "id"
	AttributeKeyUri        = "uri"
	AttributeKeyReceiver   = "receiver"
	AttributeKeyAck        = "ack"
	AttributeKeyAckSuccess = "success"
	AttributeKeyAckError   = "error"
	AttributeKeyTraceHash  = "trace_hash"
)
