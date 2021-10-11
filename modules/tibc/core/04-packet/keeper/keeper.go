package keeper

import (
	"strconv"
	"strings"

	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

// Keeper defines the TIBC packet keeper
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           codec.BinaryCodec
	clientKeeper  types.ClientKeeper
	routingKeeper types.RoutingKeeper
}

// NewKeeper creates a new tibc packet Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, key sdk.StoreKey,
	clientKeeper types.ClientKeeper,
	routingKeeper types.RoutingKeeper,
) Keeper {
	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		clientKeeper:  clientKeeper,
		routingKeeper: routingKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+host.ModuleName+"/"+types.SubModuleName)
}

// GetNextSequenceSend gets a channel's next send sequence from the store
func (k Keeper) GetNextSequenceSend(ctx sdk.Context, sourceChain, destChain string) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(host.NextSequenceSendKey(sourceChain, destChain))
	if bz == nil {
		return 1
	}

	return sdk.BigEndianToUint64(bz)
}

// SetNextSequenceSend sets a channel's next send sequence to the store
func (k Keeper) SetNextSequenceSend(ctx sdk.Context, sourceChain, destChain string, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(sequence)
	store.Set(host.NextSequenceSendKey(sourceChain, destChain), bz)
}

// GetPacketReceipt gets a packet receipt from the store
func (k Keeper) GetPacketReceipt(ctx sdk.Context, sourceChain, destChain string, sequence uint64) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(host.PacketReceiptKey(sourceChain, destChain, sequence))
	if bz == nil {
		return "", false
	}

	return string(bz), true
}

// SetPacketReceipt sets an empty packet receipt to the store
func (k Keeper) SetPacketReceipt(ctx sdk.Context, sourceChain, destChain string, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(host.PacketReceiptKey(sourceChain, destChain, sequence), []byte{byte(1)})
}

func (k Keeper) deletePacketReceipt(ctx sdk.Context, sourceChain, destChain string, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(host.PacketReceiptKey(sourceChain, destChain, sequence))
}

// HasPacketAcknowledgement check if the packet ack hash is already on the store
func (k Keeper) HasPacketReceipt(ctx sdk.Context, sourceChain, destChain string, sequence uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(host.PacketReceiptKey(sourceChain, destChain, sequence))
}

// GetPacketCommitment gets the packet commitment hash from the store
func (k Keeper) GetPacketCommitment(ctx sdk.Context, sourceChain, destChain string, sequence uint64) []byte {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(host.PacketCommitmentKey(sourceChain, destChain, sequence))
	return bz
}

// HasPacketCommitment returns true if the packet commitment exists
func (k Keeper) HasPacketCommitment(ctx sdk.Context, sourceChain, destChain string, sequence uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(host.PacketCommitmentKey(sourceChain, destChain, sequence))
}

// SetPacketCommitment sets the packet commitment hash to the store
func (k Keeper) SetPacketCommitment(ctx sdk.Context, sourceChain, destChain string, sequence uint64, commitmentHash []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(host.PacketCommitmentKey(sourceChain, destChain, sequence), commitmentHash)
}

func (k Keeper) deletePacketCommitment(ctx sdk.Context, sourceChain, destChain string, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(host.PacketCommitmentKey(sourceChain, destChain, sequence))
}

// GetCleanPacketCommitment gets the packet commitment hash from the store
func (k Keeper) GetCleanPacketCommitment(ctx sdk.Context, sourceChain, destChain string) []byte {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(host.CleanPacketCommitmentKey(sourceChain, destChain))
	return bz
}

// SetCleanPacketCommitment sets the packet commitment hash to the store
func (k Keeper) SetCleanPacketCommitment(ctx sdk.Context, sourceChain, destChain string, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(host.CleanPacketCommitmentKey(sourceChain, destChain), sdk.Uint64ToBigEndian(sequence))
}

// SetMaxAckSequence sets the max ack height to the store
func (k Keeper) SetMaxAckSequence(ctx sdk.Context, sourceChain, destChain string, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	currentMaxSeq := sdk.BigEndianToUint64(store.Get(host.MaxAckSeqKey(sourceChain, destChain)))
	if sequence > currentMaxSeq {
		currentMaxSeq = sequence
	}
	store.Set(host.MaxAckSeqKey(sourceChain, destChain), sdk.Uint64ToBigEndian(currentMaxSeq))
}

// GetMaxAckSequence gets the max ack height from the store
func (k Keeper) GetMaxAckSequence(ctx sdk.Context, sourceChain, destChain string) uint64 {
	store := ctx.KVStore(k.storeKey)
	return sdk.BigEndianToUint64(store.Get(host.MaxAckSeqKey(sourceChain, destChain)))
}

// SetPacketAcknowledgement sets the packet ack hash to the store
func (k Keeper) SetPacketAcknowledgement(ctx sdk.Context, sourceChain, destChain string, sequence uint64, ackHash []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(host.PacketAcknowledgementKey(sourceChain, destChain, sequence), ackHash)
}

func (k Keeper) deletePacketAcknowledgement(ctx sdk.Context, sourceChain, destChain string, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(host.PacketAcknowledgementKey(sourceChain, destChain, sequence))
}

// GetPacketAcknowledgement gets the packet ack hash from the store
func (k Keeper) GetPacketAcknowledgement(ctx sdk.Context, sourceChain, destChain string, sequence uint64) ([]byte, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(host.PacketAcknowledgementKey(sourceChain, destChain, sequence))
	if bz == nil {
		return nil, false
	}
	return bz, true
}

// HasPacketAcknowledgement check if the packet ack hash is already on the store
func (k Keeper) HasPacketAcknowledgement(ctx sdk.Context, sourceChain, destChain string, sequence uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(host.PacketAcknowledgementKey(sourceChain, destChain, sequence))
}

// IteratePacketSequence provides an iterator over all send, receive or ack sequences.
// For each sequence, cb will be called. If the cb returns true, the iterator
// will close and stop.
func (k Keeper) IteratePacketSequence(ctx sdk.Context, iterator db.Iterator, cb func(sourceChain, destChain string, sequence uint64) bool) {
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		sourceChain, destChain, err := host.ParseChannelPath(string(iterator.Key()))
		if err != nil {
			// return if the key is not a channel key
			return
		}

		sequence := sdk.BigEndianToUint64(iterator.Value())

		if cb(sourceChain, destChain, sequence) {
			break
		}
	}
}

// GetAllPacketSendSeqs returns all stored next send sequences.
func (k Keeper) GetAllPacketSendSeqs(ctx sdk.Context) (seqs []types.PacketSequence) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(host.KeyNextSeqSendPrefix))
	k.IteratePacketSequence(ctx, iterator, func(sourceChain, destChain string, nextSendSeq uint64) bool {
		ps := types.NewPacketSequence(sourceChain, destChain, nextSendSeq)
		seqs = append(seqs, ps)
		return false
	})
	return seqs
}

// GetAllPacketRecvSeqs returns all stored next recv sequences.
func (k Keeper) GetAllPacketRecvSeqs(ctx sdk.Context) (seqs []types.PacketSequence) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(host.KeyNextSeqRecvPrefix))
	k.IteratePacketSequence(ctx, iterator, func(sourceChain, destChain string, nextRecvSeq uint64) bool {
		ps := types.NewPacketSequence(sourceChain, destChain, nextRecvSeq)
		seqs = append(seqs, ps)
		return false
	})
	return seqs
}

// GetAllPacketAckSeqs returns all stored next acknowledgements sequences.
func (k Keeper) GetAllPacketAckSeqs(ctx sdk.Context) (seqs []types.PacketSequence) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(host.KeyNextSeqAckPrefix))
	k.IteratePacketSequence(ctx, iterator, func(sourceChain, destChain string, nextAckSeq uint64) bool {
		ps := types.NewPacketSequence(sourceChain, destChain, nextAckSeq)
		seqs = append(seqs, ps)
		return false
	})
	return seqs
}

// IteratePacketCommitment provides an iterator over all PacketCommitment objects. For each
// packet commitment, cb will be called. If the cb returns true, the iterator will close
// and stop.
func (k Keeper) IteratePacketCommitment(ctx sdk.Context, cb func(sourceChain, destChain string, sequence uint64, hash []byte) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(host.KeyPacketCommitmentPrefix))
	k.iterateHashes(ctx, iterator, cb)
}

// GetAllPacketCommitments returns all stored PacketCommitments objects.
func (k Keeper) GetAllPacketCommitments(ctx sdk.Context) (commitments []types.PacketState) {
	k.IteratePacketCommitment(ctx, func(sourceChain, destChain string, sequence uint64, hash []byte) bool {
		pc := types.NewPacketState(sourceChain, destChain, sequence, hash)
		commitments = append(commitments, pc)
		return false
	})
	return commitments
}

// IteratePacketCommitmentAtChannel provides an iterator over all PacketCommmitment objects
// at a specified channel. For each packet commitment, cb will be called. If the cb returns
// true, the iterator will close and stop.
func (k Keeper) IteratePacketCommitmentAtChannel(ctx sdk.Context, sourceChain, destChain string, cb func(_, _ string, sequence uint64, hash []byte) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(host.PacketCommitmentPrefixPath(sourceChain, destChain)))
	k.iterateHashes(ctx, iterator, cb)
}

// GetAllPacketCommitmentsAtChannel returns all stored PacketCommitments objects for a specified
// port ID and channel ID.
func (k Keeper) GetAllPacketCommitmentsAtChannel(ctx sdk.Context, sourceChain, destChain string) (commitments []types.PacketState) {
	k.IteratePacketCommitmentAtChannel(ctx, sourceChain, destChain, func(_, _ string, sequence uint64, hash []byte) bool {
		pc := types.NewPacketState(sourceChain, destChain, sequence, hash)
		commitments = append(commitments, pc)
		return false
	})
	return commitments
}

// IteratePacketReceipt provides an iterator over all PacketReceipt objects. For each
// receipt, cb will be called. If the cb returns true, the iterator will close
// and stop.
func (k Keeper) IteratePacketReceipt(ctx sdk.Context, cb func(sourceChain, destChain string, sequence uint64, receipt []byte) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(host.KeyPacketReceiptPrefix))
	k.iterateHashes(ctx, iterator, cb)
}

// GetAllPacketReceipts returns all stored PacketReceipt objects.
func (k Keeper) GetAllPacketReceipts(ctx sdk.Context) (receipts []types.PacketState) {
	k.IteratePacketReceipt(ctx, func(sourceChain, destChain string, sequence uint64, receipt []byte) bool {
		packetReceipt := types.NewPacketState(sourceChain, destChain, sequence, receipt)
		receipts = append(receipts, packetReceipt)
		return false
	})
	return receipts
}

// IteratePacketAcknowledgement provides an iterator over all PacketAcknowledgement objects. For each
// acknowledgement, cb will be called. If the cb returns true, the iterator will close
// and stop.
func (k Keeper) IteratePacketAcknowledgement(ctx sdk.Context, cb func(sourceChain, destChain string, sequence uint64, hash []byte) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(host.KeyPacketAckPrefix))
	k.iterateHashes(ctx, iterator, cb)
}

// GetAllPacketAcks returns all stored PacketAcknowledgements objects.
func (k Keeper) GetAllPacketAcks(ctx sdk.Context) (acks []types.PacketState) {
	k.IteratePacketAcknowledgement(ctx, func(sourceChain, destChain string, sequence uint64, ack []byte) bool {
		packetAck := types.NewPacketState(sourceChain, destChain, sequence, ack)
		acks = append(acks, packetAck)
		return false
	})
	return acks
}

// common functionality for IteratePacketCommitment and IteratePacketAcknowledgement
func (k Keeper) iterateHashes(_ sdk.Context, iterator db.Iterator, cb func(sourceChain, destChain string, sequence uint64, hash []byte) bool) {
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		keySplit := strings.Split(string(iterator.Key()), "/")
		sourceChain := keySplit[1]
		destChain := keySplit[2]

		sequence, err := strconv.ParseUint(keySplit[len(keySplit)-1], 10, 64)
		if err != nil {
			panic(err)
		}

		if cb(sourceChain, destChain, sequence, iterator.Value()) {
			break
		}
	}
}

// ValidatePacket validates packet sequence
func (k Keeper) ValidatePacket(ctx sdk.Context, packet exported.PacketI) error {
	if err := packet.ValidateBasic(); err != nil {
		return err
	}
	chainName := k.clientKeeper.GetChainName(ctx)
	if packet.GetRelayChain() != chainName && packet.GetDestChain() != chainName && packet.GetSourceChain() != chainName {
		return sdkerrors.Wrap(types.ErrInvalidPacket, "packet/ack illegal!")
	}
	currentCleanSeq := sdk.BigEndianToUint64(k.GetCleanPacketCommitment(ctx, packet.GetSourceChain(), packet.GetDestChain()))
	if packet.GetSequence() <= currentCleanSeq {
		return sdkerrors.Wrap(types.ErrInvalidPacket, "sequence illegal!")
	}
	return nil
}

// ValidateCleanPacket validates clean packet legal or not
func (k Keeper) ValidateCleanPacket(ctx sdk.Context, cleanPacket exported.CleanPacketI) error {
	packetSeq := cleanPacket.GetSequence()
	sourceChain := cleanPacket.GetSourceChain()
	destChain := cleanPacket.GetDestChain()
	currentCleanSeq := sdk.BigEndianToUint64(k.GetCleanPacketCommitment(ctx, sourceChain, destChain))
	currentMaxAckSeq := k.GetMaxAckSequence(ctx, sourceChain, destChain)
	if packetSeq <= currentCleanSeq || packetSeq > currentMaxAckSeq {
		return sdkerrors.Wrap(types.ErrInvalidCleanPacket, "sequence illegal!")
	}
	for seq := currentCleanSeq; seq <= packetSeq; seq++ {
		if k.HasPacketCommitment(ctx, sourceChain, destChain, seq) {
			return sdkerrors.Wrapf(types.ErrInvalidCleanPacket, "packet with sequence %d has not been ack", seq)
		}
	}
	return nil
}

// cleanAcknowledgementBySeq clean acknowledgement by sequence
func (k Keeper) cleanAcknowledgementBySeq(ctx sdk.Context, sourceChain, destChain string, sequence uint64) {
	currentCleanSeq := sdk.BigEndianToUint64(k.GetCleanPacketCommitment(ctx, sourceChain, destChain))
	for seq := currentCleanSeq + 1; seq <= sequence; seq++ {
		if k.HasPacketAcknowledgement(ctx, sourceChain, destChain, seq) {
			k.deletePacketAcknowledgement(ctx, sourceChain, destChain, seq)
		}
	}
}

// cleanReceiptBySeq clean receipt by sequence
func (k Keeper) cleanReceiptBySeq(ctx sdk.Context, sourceChain, destChain string, sequence uint64) {
	currentCleanSeq := sdk.BigEndianToUint64(k.GetCleanPacketCommitment(ctx, sourceChain, destChain))
	for seq := currentCleanSeq + 1; seq <= sequence; seq++ {
		if k.HasPacketReceipt(ctx, sourceChain, destChain, seq) {
			k.deletePacketReceipt(ctx, sourceChain, destChain, seq)
		}
	}
}
