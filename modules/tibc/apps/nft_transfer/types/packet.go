package types

import (
	"strings"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewNonFungibleTokenPacketData contructs a new NonFungibleTokenPacketData instance
func NewNonFungibleTokenPacketData(
	class, id, uri, sender, receiver string, awayFromOrigin bool, destContract string,
) NonFungibleTokenPacketData {
	return NonFungibleTokenPacketData{
		Class:          class,
		Id:             id,
		Uri:            uri,
		Sender:         sender,
		Receiver:       receiver,
		AwayFromOrigin: awayFromOrigin,
		DestContract:   destContract,
	}
}

// ValidateBasic is used for validating the nft transfer.
// NOTE: The addresses formats are not validated as the sender and recipient can have different
// formats defined by their corresponding chains that are not known to TIBC.
func (nftpd NonFungibleTokenPacketData) ValidateBasic() error {
	if strings.TrimSpace(nftpd.Sender) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "sender address cannot be blank")
	}
	if strings.TrimSpace(nftpd.Receiver) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "receiver address cannot be blank")
	}
	return nil
}

// GetBytes is a helper for serialising
func (nftpd NonFungibleTokenPacketData) GetBytes() []byte {
	return ModuleCdc.MustMarshal(&nftpd)
}
