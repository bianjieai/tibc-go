package types



// TIBC transfer events
const (
	EventTypePacket       	 = "non_fungible_token_packet"
	EventTypeNftTransfer     = "tibc_nft_transfer"

	AttributeKeyClass        = "class"
	AttributeKeyId           = "id"
	AttributeKeyUri          = "uri"
	AttributeKeyReceiver     = "receiver"
	AttributeKeyAckSuccess   = "success"
	AttributeKeyAckError     = "error"
)
