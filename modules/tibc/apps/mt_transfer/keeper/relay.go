package keeper

import (
	"strconv"
	"strings"

	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bianjieai/tibc-go/modules/tibc/apps/mt_transfer/types"
	packetType "github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
	routingtypes "github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
	coretypes "github.com/bianjieai/tibc-go/modules/tibc/core/types"
)

const (
	CLASSPREFIX = "tibc-"

	CLASSPATHPREFIX = "mt"

	DELIMITER = "/"
)

func (k Keeper) SendMtTransfer(
	ctx sdk.Context,
	class, id string,
	sender sdk.AccAddress,
	receiver string,
	destChain string,
	relayChain string,
	destContract string,
	amount uint64,
) error {
	_, found := k.mk.GetDenom(ctx, class)
	if !found {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "class %s not existed ", class)
	}

	mt, err := k.mk.GetMT(ctx, class, id)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrUnknownNFT, "invalid mt %s from class %s", id, class)
	}

	sourceChain := k.ck.GetChainName(ctx)
	if sourceChain == destChain {
		return sdkerrors.Wrapf(types.ErrScChainEqualToDestChain, "invalid destChain %s equals to scChain %s", destChain, sourceChain)
	}

	fullClassPath := class

	// deconstruct the mt class into the class trace info to determine if the sender is the source chain
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

	// determine whether mt is sent from the source chain or sent back to the source chain from other chains
	awayFromOrigin := k.determineAwayFromOrigin(fullClassPath, destChain)

	// get the next sequence
	sequence := k.pk.GetNextSequenceSend(ctx, sourceChain, destChain)

	// get moudle address
	moudleAddr := k.GetMtTransferModuleAddr(types.ModuleName)

	labels = append(labels, telemetry.NewLabel(coretypes.LabelSource, strconv.FormatBool(awayFromOrigin)))
	if awayFromOrigin {
		// Two conversion scenarios
		// 1. mtClass -> tibc-hash(mt/A/B/mtClass)
		// 2. tibc-hash(mt/A/B/mtClass) -> tibc-hash(mt/A/B/C/mtClass)

		// Two things need to be done
		// 1. lock mt  |send to moduleAccount
		// 2. send packet
		// The mt attribute must be marked as unchanged (it cannot be changed in practice)
		// because the TransferOwner method will verify when UpdateRestricted is true
		if err := k.mk.TransferOwner(ctx, class, id, amount, sender, moudleAddr); err != nil {
			return err
		}
	} else {
		// burn mt
		if err := k.mk.BurnMT(ctx, class, id, amount, sender); err != nil {
			return err
		}
	}

	// constructs packet
	packetData := types.NewMultiTokenPacketData(
		fullClassPath,
		id,
		sender.String(),
		receiver,
		awayFromOrigin,
		destContract,
		amount,
		mt.GetData(),
	)

	packet := packetType.NewPacket(packetData.GetBytes(), sequence, sourceChain, destChain, relayChain, string(routingtypes.MT))

	defer func() {
		telemetry.SetGaugeWithLabels(
			[]string{"tx", "msg", "tibc", "mttransfer"},
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
	B receive packet from A : class -> mt/A/B/class
	c receive packet from B : mt/A/B/class -> mt/A/B/C/class
C->B->A  away_from_source == flase
	B receive packet from C : mt/A/B/C/class -> mt/A/B/class
	A receive packet from B : mt/A/B/class -> class
*/
func (k Keeper) OnRecvPacket(ctx sdk.Context, packet packetType.Packet, data types.MultiTokenPacketData) error {
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

	moudleAddr := k.GetMtTransferModuleAddr(types.ModuleName)
	var newClassPath string
	if data.AwayFromOrigin {
		labels = append(labels, telemetry.NewLabel(coretypes.LabelSource, "true"))
		newClassPath = k.getAwayNewClassPath(packet.SourceChain, packet.DestinationChain, data.Class)
		voucherClass := k.getIBCClassFromClassPath(ctx, newClassPath)

		_, found := k.mk.GetDenom(ctx, voucherClass)
		if !found {
			_ = k.mk.IssueDenom(ctx, voucherClass, "", moudleAddr, []byte(""))
		}

		if !k.mk.HasMT(ctx, voucherClass, data.Id) {
			if _, err = k.mk.IssueMT(ctx, voucherClass, data.Id, data.Amount, data.Data, moudleAddr); err != nil {
				return err
			}
		} else {
			if err := k.mk.MintMT(ctx, voucherClass, data.Id, data.Amount, moudleAddr); err != nil {
				return err
			}
		}

		if err := k.mk.TransferOwner(ctx, voucherClass, data.Id, data.Amount, moudleAddr, receiver); err != nil {
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
		if err := k.mk.TransferOwner(ctx, voucherClass, data.Id, data.Amount, moudleAddr, receiver); err != nil {
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

func (k Keeper) OnAcknowledgementPacket(ctx sdk.Context, data types.MultiTokenPacketData, ack packetType.Acknowledgement) error {
	switch ack.Response.(type) {
	case *packetType.Acknowledgement_Error:
		return k.refundPacketToken(ctx, data)
	default:
		// the acknowledgement succeeded on the receiving chain so nothing
		// needs to be executed and no error needs to be returned
		return nil
	}
}

func (k Keeper) refundPacketToken(ctx sdk.Context, data types.MultiTokenPacketData) error {
	sender, err := sdk.AccAddressFromBech32(data.Sender)
	if err != nil {
		return err
	}

	moudleAddr := k.GetMtTransferModuleAddr(types.ModuleName)
	classTrace := types.ParseClassTrace(data.Class)
	voucherClass := classTrace.IBCClass()

	if data.AwayFromOrigin {
		// unlock
		if err := k.mk.TransferOwner(ctx, voucherClass, data.Id, data.Amount,
			k.GetMtTransferModuleAddr(types.ModuleName), sender); err != nil {
			return err
		}
	} else {
		if err := k.mk.MintMT(ctx, voucherClass, data.Id, data.Amount, moudleAddr); err != nil {
			return err
		}

		if err := k.mk.TransferOwner(ctx, voucherClass, data.Id, data.Amount, moudleAddr, sender); err != nil {
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

			1. B -> C    class:mt/A/B/class 	| sourceChain:B  | destChain:C |awayFromOrigin = true
			2. C -> B    class:mt/A/B/C/class   | sourceChain:C  | destChain:B |awayFromOrigin = false
			3. B -> A    class:mt/A/B/class 	| sourceChain:B  | destChain:A |awayFromOrigin = false
	*/
	if !strings.HasPrefix(class, CLASSPATHPREFIX) ||
		(strings.HasPrefix(class, CLASSPATHPREFIX) && !strings.Contains(class, DELIMITER)) {
		return true
	}

	classSplit := strings.Split(class, DELIMITER)

	return classSplit[len(classSplit)-3] != destChain
}

//fullClassPath = "tibcmt" + "/" + packet.SourceChain + "/" + packet.DestinationChain + "/" + data.Class
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
		// mt/A/B/class -> mt/A/B/C/class
		// [mt][A][B][class] -> [mt][A][B][C][class]
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
		// mt/A/B/class -> class
		newClassPath = classSplit[len(classSplit)-1]
	} else {
		// mt/A/B/C/class -> mt/A/B/class
		classSplit = append(classSplit[:len(classSplit)-2], classSplit[len(classSplit)-1])
		newClassPath = strings.Join(classSplit, DELIMITER)
	}
	return
}

// example : mt/A/B/class -> tibc-hash(mt/A/B/class)
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
