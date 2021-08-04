package keeper

import (
	"github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	packetType "github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
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

