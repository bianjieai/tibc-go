package keeper

import (
	"github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	packetType "github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
	"strings"
)

const (
	Prefix = "tibc/nft"
)

func (k Keeper) SendNftTransfer(
	ctx sdk.Context,
	class, id, uri string,
	sender sdk.AccAddress,
	receiver string,
	awayFromOrigin bool,
	destChain, relayChain string,
) error {
	// get the next sequence
	// todo
	var sequence = uint64(0)

	if awayFromOrigin{
		// lock nft
		if err := k.nftKeeper.TransferOwner(ctx, class, id, "", uri, "",
			sender, k.GetNftTransferModuleAddr(types.ModuleName)); err != nil{
			return err
		}
	} else {
		// burn nft
		if err := k.nftKeeper.BurnNFT(ctx, class, id,
			k.GetNftTransferModuleAddr(types.ModuleName)); err != nil{
			return err
		}
	}

	// packetdata
	packetdata  := types.NewNonFungibleTokenPacketData(class, id, uri,
		sender.String(), receiver, awayFromOrigin)

	// constructs packet
	// todo
	packet := packetType.NewPacket([]byte(""), sequence, "sourceChain", destChain, "port")

	// send packet
	// todo


	return nil
}


func (k Keeper)OnRecvPacket(ctx sdk.Context, packet packetType.Packet, data types.NonFungibleTokenPacketData) error{
	// validate packet data upon receiving
	if err := data.ValidateBasic(); err != nil {
		return err
	}

	// decode the sender address
	sender, err := sdk.AccAddressFromBech32(data.Sender)
	if err != nil {
		return err
	}

	// decode the receiver address
	receiver, err := sdk.AccAddressFromBech32(data.Receiver)
	if err != nil {
		return err
	}

	if data.AwayFromOrigin{
		if strings.HasPrefix(data.Class, Prefix){
			// has prefix  tibc/nft/a/b/class
			classSplit := strings.Split(data.Class, "/")
			classSplit = append(classSplit[:len(classSplit) - 2], packet.SourceChain)
			newClass := strings.Join(classSplit, "/")
			if err := k.nftKeeper.MintNFT(ctx, newClass, data.Id, "", data.Uri, "", sender); err != nil{
				return err
			}
		} else {
			// not has prefix  tibc/nft/a/class
			newClass := Prefix + "/" + packet.SourceChain + data.Class
			if err := k.nftKeeper.MintNFT(ctx, newClass, data.Id, "", data.Uri, "", sender); err != nil{
				return err
			}
			// lock todo
			// send packet  need judge relay chain empty todo
		}
	} else {
		if strings.HasPrefix(data.Class, Prefix){
			classSplit := strings.Split(data.Class, "/")
			destChain := classSplit[len(classSplit) - 2]
			if destChain != packet.DestinationChain{
				// return err  must equal
			}
			var newClass string
			if len(classSplit) == 4{
				// tibc/nft/A/class -> class
				newClass = classSplit[len(classSplit) - 1]
			} else {
				// tibc/nft/A/B/class -> tibc/nft/A/class
				classSplit = append(classSplit[:len(classSplit) - 3], classSplit[len(classSplit) - 1])
				newClass = strings.Join(classSplit, "/")
			}
			// unlock : from moduleAddr to receiver
			if err := k.nftKeeper.TransferOwner(ctx, newClass, data.Id, "", data.Uri, "",
				k.GetNftTransferModuleAddr(types.ModuleName), receiver); err != nil{
				return err
			}

			// if two skip
			// need create packet &&sendpacket todo

		} else {
			//  return err must has prefix if awayfromchain todo
		}
	}
	return nil
}


func (k Keeper) OnAcknowledgementPacket(ctx sdk.Context, packet packetType.Packet, data types.NonFungibleTokenPacketData, ack channeltypes.Acknowledgement) error {
	switch ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:
		return k.refundPacketToken(ctx, packet, data)
	default:
		// the acknowledgement succeeded on the receiving chain so nothing
		// needs to be executed and no error needs to be returned
		return nil
	}
}

func (k Keeper) refundPacketToken(ctx sdk.Context, packet packetType.Packet, data types.NonFungibleTokenPacketData) error {
	// decode the sender address
	sender, err := sdk.AccAddressFromBech32(data.Sender)
	if err != nil {
		return err
	}

	// decode the recevier address
	receiver, err := sdk.AccAddressFromBech32(data.Receiver)
	if err != nil {
		return err
	}

	if data.AwayFromOrigin{
		// unlock
		if err := k.nftKeeper.TransferOwner(ctx, data.Class, data.Id, "", data.Uri, "",
			k.GetNftTransferModuleAddr(types.ModuleName), receiver); err != nil{
			return err
		}

	} else {
		// mintNFT
		if err := k.nftKeeper.MintNFT(ctx, data.Class, data.Id, "", data.Uri, "", sender); err != nil{
			return err
		}
	}
	return nil
}