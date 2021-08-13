package keeper

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/types"
	packetType "github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
	routingtypes "github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
)

const (
	PREFIX = "tibc/nft"
)

func (k Keeper) SendNftTransfer(
	ctx sdk.Context,
	class, id string,
	sender sdk.AccAddress,
	receiver, destChain, relayChain string,
) error {
	// class must be existed
	_, found := k.nk.GetDenom(ctx, class)
	if !found {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "class %s not existed ", class)
	}
	// get nft
	nft, err := k.nk.GetNFT(ctx, class, id)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrUnknownNFT, "invalid NFT %s from collection %s", id, class)
	}

	// decode the sender address
	sender, err = sdk.AccAddressFromBech32(sender.String())
	if err != nil {
		return err
	}

	moudleAddr := k.GetNftTransferModuleAddr(types.ModuleName)

	// sourceChain cannot be equal to destChain
	sourceChain := k.ck.GetChainName(ctx)
	if sourceChain == destChain {
		return sdkerrors.Wrapf(types.ErrScChainEqualToDestChain, "invalid destChain %s equals to scChain %s", destChain, sourceChain)
	}

	// determine whether nft is sent from the source chain or sent back to the source chain from other chains
	awayFromOrigin := k.determineAwayFromOrigin(class, destChain)

	// get the next sequence
	sequence, _ := k.pk.GetNextSequenceSend(ctx, sourceChain, destChain)

	if awayFromOrigin {
		// Two conversion scenarios
		// 1. class -> tibc/nft/A/class
		// 2. tibc/nft/A/class -> tibc/nft/A/B/class

		// Two things need to be done
		// 1. lock nft  |send to moduleAccount
		// 2. send packet
		if err := k.nk.TransferOwner(ctx, class, id, "", nft.GetURI(), "", sender, moudleAddr); err != nil {
			return err
		}
	} else {
		// burn nft
		if err := k.nk.BurnNFT(ctx, class, id, sender); err != nil {
			return err
		}
	}

	// constructs packet
	packetData := types.NewNonFungibleTokenPacketData(
		class,
		id,
		nft.GetURI(),
		sender.String(),
		receiver,
		awayFromOrigin,
	)

	packet := packetType.NewPacket(packetData.GetBytes(), sequence, sourceChain, destChain, relayChain, string(routingtypes.NFT))

	// send packet
	if err := k.pk.SendPacket(ctx, packet); err != nil {
		return err
	}

	return nil
}

/*
OnRecvPacket
A->B->C  away_from_source == true
	B receive packet from A : class -> tibc/nft/A/class
	c receive packet from B : tibc/nft/A/class -> tibc/nft/A/B/class
C->B->A  away_from_source == flase
	B receive packet from C : tibc/nft/A/B/class -> tibc/nft/A/class
	A receive packet from B : tibc/nft/A/class -> class
*/
func (k Keeper) OnRecvPacket(ctx sdk.Context, packet packetType.Packet, data types.NonFungibleTokenPacketData) error {
	// validate packet data upon receiving
	if err := data.ValidateBasic(); err != nil {
		return err
	}

	// decode the sender address
	receiver, err := sdk.AccAddressFromBech32(data.Receiver)
	if err != nil {
		return err
	}

	var newClass string
	if data.AwayFromOrigin {
		if strings.HasPrefix(data.Class, PREFIX) {
			// tibc/nft/A/class -> tibc/nft/A/B/class
			// [tibc][nft][A][class] -> [tibc][nft][A][B][class]
			classSplit := strings.Split(data.Class, "/")
			classSplit = append(classSplit[:len(classSplit)-2], packet.SourceChain)
			newClass = strings.Join(classSplit, "/")
		} else {
			// class -> tibc/nft/A/classs
			newClass = PREFIX + "/" + packet.SourceChain + "/" + data.Class
		}

		_, found := k.nk.GetDenom(ctx, newClass)
		if !found {
			if err := k.nk.IssueDenom(ctx, newClass, "", "", "", receiver, false, false); err != nil {
				return err
			}
		}

		if err := k.nk.MintNFT(ctx, newClass, data.Id, "", data.Uri, "", receiver); err != nil {
			return err
		}
	} else {
		if strings.HasPrefix(data.Class, PREFIX) {
			classSplit := strings.Split(data.Class, "/")

			if len(classSplit) == 4 {
				// tibc/nft/A/class -> class
				newClass = classSplit[len(classSplit)-1]
			} else {
				// tibc/nft/A/B/class -> tibc/nft/A/class
				classSplit = append(classSplit[:len(classSplit)-3], classSplit[len(classSplit)-1])
				newClass = strings.Join(classSplit, "/")
			}

			// burn nft
			if err := k.nk.BurnNFT(ctx, newClass, data.Id, receiver); err != nil {
				return err
			}
		} else {
			return sdkerrors.Wrapf(types.ErrInvalidDenom, "class has no prefix: %s", data.Class)
		}
	}
	return nil
}

func (k Keeper) OnAcknowledgementPacket(ctx sdk.Context, data types.NonFungibleTokenPacketData, ack packetType.Acknowledgement) error {
	switch ack.Response.(type) {
	case *packetType.Acknowledgement_Error:
		return k.refundPacketToken(ctx, data)
	default:
		// the acknowledgement succeeded on the receiving chain so nothing
		// needs to be executed and no error needs to be returned
		return nil
	}
}

func (k Keeper) refundPacketToken(ctx sdk.Context, data types.NonFungibleTokenPacketData) error {
	// decode the sender address
	sender, err := sdk.AccAddressFromBech32(data.Sender)
	if err != nil {
		return err
	}

	if data.AwayFromOrigin {
		// unlock
		if err := k.nk.TransferOwner(ctx, data.Class, data.Id, "", data.Uri, "",
			k.GetNftTransferModuleAddr(types.ModuleName), sender); err != nil {
			return err
		}

	} else {
		// mintNFT
		if err := k.nk.MintNFT(ctx, data.Class, data.Id, "", data.Uri, "", sender); err != nil {
			return err
		}
	}
	return nil
}

// determineAwayFromOrigin determine whether nft is sent from the source chain or sent back to the source chain from other chains
func (k Keeper) determineAwayFromOrigin(class, destChain string) (awayFromOrigin bool) {
	/*
		-- not has prefix
		1. A -> B  class:class | sourceChain:A  | destChain:B |awayFromOrigin = true
		-- has prefix
		1. B -> C    class:tibc/nft/A/class | sourceChain:B  | destChain:C |awayFromOrigin = true
		2. C -> B    class:tibc/nft/A/B/class | sourceChain:C  | destChain:B |awayFromOrigin = false
		3. B -> A    class:tibc/nft/A/class | sourceChain:B  | destChain:A |awayFromOrigin = false
	*/
	if !strings.HasPrefix(class, PREFIX) {
		return true
	}

	classSplit := strings.Split(class, "/")
	if classSplit[len(classSplit)-2] == destChain {
		return false
	}
	return true
}
