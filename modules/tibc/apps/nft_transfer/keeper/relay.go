package keeper

import (
	"strconv"
	"strings"

	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/types"
	packetType "github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
	routingtypes "github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
	coretypes "github.com/bianjieai/tibc-go/modules/tibc/core/types"
)

const (
	CLASSPREFIX = "tibc-"

	CLASSPATHPREFIX = "nft"

	DELIMITER = "/"

	// DoNotModify used to indicate that some field should not be updated
	DoNotModify = "[do-not-modify]"
)

func (k Keeper) SendNftTransfer(
	ctx sdk.Context,
	class, id string,
	sender sdk.AccAddress,
	receiver string,
	destChain string,
	relayChain string,
	destContract string,
) error {
	_, found := k.nk.GetDenom(ctx, class)
	if !found {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "class %s not existed ", class)
	}

	nft, err := k.nk.GetNFT(ctx, class, id)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrUnknownNFT, "invalid NFT %s from collection %s", id, class)
	}

	sourceChain := k.ck.GetChainName(ctx)
	if sourceChain == destChain {
		return sdkerrors.Wrapf(types.ErrScChainEqualToDestChain, "invalid destChain %s equals to scChain %s", destChain, sourceChain)
	}

	fullClassPath := class

	// deconstruct the nft class into the class trace info to determine if the sender is the source chain
	if strings.HasPrefix(class, CLASSPREFIX) {
		fullClassPath, err = k.ClassPathFromHash(ctx, class)
		if err != nil {
			return err
		}
	}

	labels := []metrics.Label{
		telemetry.NewLabel(coretypes.LabelSourceChain, sourceChain),
		telemetry.NewLabel(coretypes.LabelDestinationChain, destChain),
	}

	// determine whether nft is sent from the source chain or sent back to the source chain from other chains
	awayFromOrigin := k.determineAwayFromOrigin(fullClassPath, destChain)

	// get the next sequence
	sequence := k.pk.GetNextSequenceSend(ctx, sourceChain, destChain)

	// get moudle address
	moudleAddr := k.GetNftTransferModuleAddr(types.ModuleName)

	labels = append(labels, telemetry.NewLabel(coretypes.LabelSource, strconv.FormatBool(awayFromOrigin)))
	if awayFromOrigin {
		// Two conversion scenarios
		// 1. nftClass -> tibc-hash(nft/A/B/nftClass)
		// 2. tibc-hash(nft/A/B/nftClass) -> tibc-hash(nft/A/B/C/nftClass)

		// Two things need to be done
		// 1. lock nft  |send to moduleAccount
		// 2. send packet
		// The nft attribute must be marked as unchanged (it cannot be changed in practice)
		// because the TransferOwner method will verify when UpdateRestricted is true
		if err := k.nk.TransferOwner(ctx, class, id, DoNotModify, DoNotModify, DoNotModify, sender, moudleAddr); err != nil {
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
		fullClassPath,
		id,
		nft.GetURI(),
		sender.String(),
		receiver,
		awayFromOrigin,
		destContract,
	)

	packet := packetType.NewPacket(packetData.GetBytes(), sequence, sourceChain, destChain, relayChain, string(routingtypes.NFT))

	defer func() {
		telemetry.SetGaugeWithLabels(
			[]string{"tx", "msg", "tibc", "nfttransfer"},
			1,
			[]metrics.Label{telemetry.NewLabel(coretypes.LabelDenom, fullClassPath)},
		)

		telemetry.IncrCounterWithLabels(
			[]string{"tibc", types.ModuleName, "send"},
			1,
			labels,
		)
	}()
	// send packet
	return k.pk.SendPacket(ctx, packet)
}

/*
OnRecvPacket
A->B->C  away_from_source == true
	B receive packet from A : class -> nft/A/B/class
	c receive packet from B : nft/A/B/class -> nft/A/B/C/class
C->B->A  away_from_source == flase
	B receive packet from C : nft/A/B/C/class -> nft/A/B/class
	A receive packet from B : nft/A/B/class -> class
*/
func (k Keeper) OnRecvPacket(ctx sdk.Context, packet packetType.Packet, data types.NonFungibleTokenPacketData) error {
	// validate packet data upon receiving
	if err := data.ValidateBasic(); err != nil {
		return err
	}

	receiver, err := sdk.AccAddressFromBech32(data.Receiver)
	if err != nil {
		return err
	}

	labels := []metrics.Label{
		telemetry.NewLabel(coretypes.LabelSourceChain, packet.SourceChain),
		telemetry.NewLabel(coretypes.LabelDestinationChain, packet.DestinationChain),
	}

	moudleAddr := k.GetNftTransferModuleAddr(types.ModuleName)
	var newClassPath string
	if data.AwayFromOrigin {
		labels = append(labels, telemetry.NewLabel(coretypes.LabelSource, "true"))

		newClassPath = k.getAwayNewClassPath(packet.SourceChain, packet.DestinationChain, data.Class)

		voucherClass := k.getIBCClassFromClassPath(ctx, newClassPath)

		_, found := k.nk.GetDenom(ctx, voucherClass)
		if !found {
			// The creator of cross-chain denom must be a module account,
			// and only the owner of Denom can issue NFT under this category,
			// and no one under this category can update NFT,
			// that is, updateRestricted is true and mintRestricted is true
			if err := k.nk.IssueDenom(ctx, voucherClass, "", "", "", moudleAddr, true, true); err != nil {
				return err
			}
		}

		// Only module accounts can mint nft, because mintRestricted is true,
		// you must first mint nft to the module account, and then transfer nft ownership to the receiver
		if err := k.nk.MintNFT(ctx, voucherClass, data.Id, "", data.Uri, "", moudleAddr); err != nil {
			return err
		}

		if err := k.nk.TransferOwner(ctx, voucherClass, data.Id, DoNotModify, DoNotModify, DoNotModify, moudleAddr, receiver); err != nil {
			return err
		}
	} else {
		labels = append(labels, telemetry.NewLabel(coretypes.LabelSource, "false"))

		if !strings.HasPrefix(data.Class, CLASSPATHPREFIX) {
			return sdkerrors.Wrapf(types.ErrInvalidDenom, "class has no prefix: %s", data.Class)
		}

		newClassPath = k.getBackNewClassPath(data.Class)
		classTrace := types.ParseClassTrace(newClassPath)
		voucherClass := classTrace.IBCClass()
		// unlock
		if err := k.nk.TransferOwner(ctx, voucherClass, data.Id, DoNotModify, DoNotModify, DoNotModify, moudleAddr, receiver); err != nil {
			return err
		}
	}

	defer func() {
		telemetry.SetGaugeWithLabels(
			[]string{"tx", "msg", "ibc", "transfer"},
			1,
			[]metrics.Label{telemetry.NewLabel(coretypes.LabelDenom, newClassPath)},
		)

		telemetry.IncrCounterWithLabels(
			[]string{"ibc", types.ModuleName, "send"},
			1,
			labels,
		)
	}()

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
	sender, err := sdk.AccAddressFromBech32(data.Sender)
	if err != nil {
		return err
	}

	moudleAddr := k.GetNftTransferModuleAddr(types.ModuleName)

	classTrace := types.ParseClassTrace(data.Class)
	voucherClass := classTrace.IBCClass()

	if data.AwayFromOrigin {
		// unlock
		if err := k.nk.TransferOwner(ctx, voucherClass, data.Id, DoNotModify, DoNotModify, DoNotModify,
			k.GetNftTransferModuleAddr(types.ModuleName), sender); err != nil {
			return err
		}
	} else {
		// mintNFT
		// Corresponding to burnNft, because the mintRestricted attribute of denom generated by any cross-chain nft is true,
		// so to re-mint nft, you must first mintnft to the module account, and then transfer the nft ownership to the sender account
		if err := k.nk.MintNFT(ctx, voucherClass, data.Id, "", data.Uri, "", moudleAddr); err != nil {
			return err
		}

		if err := k.nk.TransferOwner(ctx, voucherClass, data.Id, DoNotModify, DoNotModify, DoNotModify, moudleAddr, sender); err != nil {
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
			 First take the source chain from the path
		    (this path represents the path generated from the source chain in the target chain),
		    and then judge whether the source chain is equal to the target chain,
		    if it is equal, it means it is close to the source chain,
		    if it is not equal then Indicates that we continue to stay away from the source chain

			1. B -> C    class:nft/A/B/class 	| sourceChain:B  | destChain:C |awayFromOrigin = true
			2. C -> B    class:nft/A/B/C/class  | sourceChain:C  | destChain:B |awayFromOrigin = false
			3. B -> A    class:nft/A/B/class 	| sourceChain:B  | destChain:A |awayFromOrigin = false
	*/
	if !strings.HasPrefix(class, CLASSPATHPREFIX) ||
		(strings.HasPrefix(class, CLASSPATHPREFIX) && !strings.Contains(class, DELIMITER)) {
		return true
	}

	classSplit := strings.Split(class, DELIMITER)

	return classSplit[len(classSplit)-3] != destChain
}

//fullClassPath = "tibcnft" + "/" + packet.SourceChain + "/" + packet.DestinationChain + "/" + data.Class
func (k Keeper) concatClassPath(scChain, destChain, class string) string {
	var b strings.Builder
	b.WriteString(CLASSPATHPREFIX)
	b.WriteString(DELIMITER)
	b.WriteString(scChain)
	b.WriteString(DELIMITER)
	b.WriteString(destChain)
	b.WriteString(DELIMITER)
	b.WriteString(class)
	return b.String()
}

// getAwayNewClassPath
func (k Keeper) getAwayNewClassPath(scChain, destChain, class string) (newClassPath string) {
	if strings.HasPrefix(class, CLASSPATHPREFIX) && strings.Contains(class, DELIMITER) {
		// nft/A/B/class -> nft/A/B/C/class
		// [nft][A][B][class] -> [nft][A][B][C][class]
		classSplit := strings.Split(class, DELIMITER)
		classSplit = append(classSplit[:len(classSplit)-1], append([]string{destChain}, classSplit[len(classSplit)-1:]...)...)
		newClassPath = strings.Join(classSplit, DELIMITER)
	} else {
		// class -> nft/A/B/class
		newClassPath = k.concatClassPath(scChain, destChain, class)
	}
	return
}

// getBackNewClassPath
func (k Keeper) getBackNewClassPath(class string) (newClassPath string) {
	classSplit := strings.Split(class, DELIMITER)

	if len(classSplit) == 4 {
		// nft/A/B/class -> class
		newClassPath = classSplit[len(classSplit)-1]
	} else {
		// nft/A/B/C/class -> nft/A/B/class
		classSplit = append(classSplit[:len(classSplit)-2], classSplit[len(classSplit)-1])
		newClassPath = strings.Join(classSplit, DELIMITER)
	}
	return
}

// example : nft/A/B/class -> tibc-hash(nft/A/B/class)
func (k Keeper) getIBCClassFromClassPath(ctx sdk.Context, classPath string) string {
	// construct the class trace from the full raw class
	classTrace := types.ParseClassTrace(classPath)
	traceHash := classTrace.Hash()

	if !k.HasClassTrace(ctx, traceHash) {
		k.SetClassTrace(ctx, classTrace)
	}

	return classTrace.IBCClass()
}

// ClassPathFromHash returns the full class path prefix from an ibc class with a hash
// component.
func (k Keeper) ClassPathFromHash(ctx sdk.Context, class string) (string, error) {
	// trim the class prefix, by default "tibc-"
	hexHash := class[len(types.ClassPrefix+"-"):]

	hash, err := types.ParseHexHash(hexHash)
	if err != nil {
		return "", sdkerrors.Wrap(types.ErrInvalidDenom, err.Error())
	}

	denomTrace, found := k.GetClassTrace(ctx, hash)
	if !found {
		return "", sdkerrors.Wrap(types.ErrTraceNotFound, hexHash)
	}

	fullDenomPath := denomTrace.GetFullClassPath()
	return fullDenomPath, nil
}
