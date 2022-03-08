package types

import (
	"strings"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewMultiTokenPacketData contructs a new MultiTokenPacketData instance
func NewMultiTokenPacketData(
	class, id, sender, receiver string,
	awayFromOrigin bool, destContract string,
	amount uint64, data []byte,
) MultiTokenPacketData {
	return MultiTokenPacketData{
		Class:          class,
		Id:             id,
		Data:           data,
		Sender:         sender,
		Receiver:       receiver,
		AwayFromOrigin: awayFromOrigin,
		DestContract:   destContract,
		Amount:         amount,
	}
}

// ValidateBasic is used for validating the mt transfer.
// NOTE: The addresses formats are not validated as the sender and recipient can have different
// formats defined by their corresponding chains that are not known to TIBC.
func (mtpd MultiTokenPacketData) ValidateBasic() error {
	if strings.TrimSpace(mtpd.Sender) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "sender address cannot be blank")
	}
	if strings.TrimSpace(mtpd.Receiver) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "receiver address cannot be blank")
	}
	if mtpd.Amount <= 0 {
		return ErrInvalidAmount
	}
	return nil
}

// GetBytes is a helper for serialising
func (mtpd MultiTokenPacketData) GetBytes() []byte {
	return ModuleCdc.MustMarshal(&mtpd)
}
