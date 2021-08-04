package types




// NewNonFungibleTokenPacketData contructs a new NonFungibleTokenPacketData instance
func NewNonFungibleTokenPacketData(
	class, id, uri, sender, receiver string,
	awayFromOrigin bool,
	) NonFungibleTokenPacketData {
	return NonFungibleTokenPacketData{
		Class:  	class,
		Id:   		id,
		Uri:		uri,
		Sender:     sender,
		Receiver: receiver,
		AwayFromOrigin: awayFromOrigin,
	}
}

