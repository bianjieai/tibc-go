package keeper

import (
	"strconv"
	"strings"

	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
)

// Keeper defines the IBC channel keeper
type Keeper struct {
	// implements gRPC QueryServer interface
	types.QueryServer

	storeKey      sdk.StoreKey
	cdc           codec.BinaryMarshaler
	clientKeeper  types.ClientKeeper
	routingKeeper types.RoutingKeeper
}

// NewKeeper creates a new IBC channel Keeper instance
func NewKeeper(
	cdc codec.BinaryMarshaler, key sdk.StoreKey,
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

// GetNextChannelSequence gets the next channel sequence from the store.
func (k Keeper) GetNextChannelSequence(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(types.KeyNextChannelSequence))
	if bz == nil {
		panic("next channel sequence is nil")
	}

	return sdk.BigEndianToUint64(bz)
}

// SetNextChannelSequence sets the next channel sequence to the store.
func (k Keeper) SetNextChannelSequence(ctx sdk.Context, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(sequence)
	store.Set([]byte(types.KeyNextChannelSequence), bz)
}

// GetNextSequenceSend gets a channel's next send sequence from the store
func (k Keeper) GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(host.NextSequenceSendKey(portID, channelID))
	if bz == nil {
		return 0, false
	}

	return sdk.BigEndianToUint64(bz), true
}

// SetNextSequenceSend sets a channel's next send sequence to the store
func (k Keeper) SetNextSequenceSend(ctx sdk.Context, sourceChain, destChain string, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(sequence)
	store.Set(host.NextSequenceSendKey(sourceChain, destChain), bz)
}

// GetNextSequenceRecv gets a channel's next receive sequence from the store
func (k Keeper) GetNextSequenceRecv(ctx sdk.Context, portID, channelID string) (uint64, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(host.NextSequenceRecvKey(portID, channelID))
	if bz == nil {
		return 0, false
	}

	return sdk.BigEndianToUint64(bz), true
}

// SetNextSequenceRecv sets a channel's next receive sequence to the store
func (k Keeper) SetNextSequenceRecv(ctx sdk.Context, portID, channelID string, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(sequence)
	store.Set(host.NextSequenceRecvKey(portID, channelID), bz)
}

// GetNextSequenceAck gets a channel's next ack sequence from the store
func (k Keeper) GetNextSequenceAck(ctx sdk.Context, portID, channelID string) (uint64, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(host.NextSequenceAckKey(portID, channelID))
	if bz == nil {
		return 0, false
	}

	return sdk.BigEndianToUint64(bz), true
}

// SetNextSequenceAck sets a channel's next ack sequence to the store
func (k Keeper) SetNextSequenceAck(ctx sdk.Context, portID, channelID string, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(sequence)
	store.Set(host.NextSequenceAckKey(portID, channelID), bz)
}

// GetPacketReceipt gets a packet receipt from the store
func (k Keeper) GetPacketReceipt(ctx sdk.Context, portID, channelID string, sequence uint64) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(host.PacketReceiptKey(portID, channelID, sequence))
	if bz == nil {
		return "", false
	}

	return string(bz), true
}

// SetPacketReceipt sets an empty packet receipt to the store
func (k Keeper) SetPacketReceipt(ctx sdk.Context, portID, channelID string, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(host.PacketReceiptKey(portID, channelID, sequence), []byte{byte(1)})
}

// GetPacketCommitment gets the packet commitment hash from the store
func (k Keeper) GetPacketCommitment(ctx sdk.Context, portID, channelID string, sequence uint64) []byte {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(host.PacketCommitmentKey(portID, channelID, sequence))
	return bz
}

// HasPacketCommitment returns true if the packet commitment exists
func (k Keeper) HasPacketCommitment(ctx sdk.Context, portID, channelID string, sequence uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(host.PacketCommitmentKey(portID, channelID, sequence))
}

// SetPacketCommitment sets the packet commitment hash to the store
func (k Keeper) SetPacketCommitment(ctx sdk.Context, portID, channelID string, sequence uint64, commitmentHash []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(host.PacketCommitmentKey(portID, channelID, sequence), commitmentHash)
}

func (k Keeper) deletePacketCommitment(ctx sdk.Context, portID, channelID string, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(host.PacketCommitmentKey(portID, channelID, sequence))
}

// SetPacketAcknowledgement sets the packet ack hash to the store
func (k Keeper) SetPacketAcknowledgement(ctx sdk.Context, sourceChain, destinationChain string, sequence uint64, ackHash []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(host.PacketAcknowledgementKey(sourceChain, destinationChain, sequence), ackHash)
}

// GetPacketAcknowledgement gets the packet ack hash from the store
func (k Keeper) GetPacketAcknowledgement(ctx sdk.Context, portID, channelID string, sequence uint64) ([]byte, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(host.PacketAcknowledgementKey(portID, channelID, sequence))
	if bz == nil {
		return nil, false
	}
	return bz, true
}

// HasPacketAcknowledgement check if the packet ack hash is already on the store
func (k Keeper) HasPacketAcknowledgement(ctx sdk.Context, portID, channelID string, sequence uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(host.PacketAcknowledgementKey(portID, channelID, sequence))
}

// IteratePacketSequence provides an iterator over all send, receive or ack sequences.
// For each sequence, cb will be called. If the cb returns true, the iterator
// will close and stop.
func (k Keeper) IteratePacketSequence(ctx sdk.Context, iterator db.Iterator, cb func(portID, channelID string, sequence uint64) bool) {
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		portID, channelID, err := host.ParseChannelPath(string(iterator.Key()))
		if err != nil {
			// return if the key is not a channel key
			return
		}

		sequence := sdk.BigEndianToUint64(iterator.Value())

		if cb(portID, channelID, sequence) {
			break
		}
	}
}

// GetAllPacketSendSeqs returns all stored next send sequences.
func (k Keeper) GetAllPacketSendSeqs(ctx sdk.Context) (seqs []types.PacketSequence) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(host.KeyNextSeqSendPrefix))
	k.IteratePacketSequence(ctx, iterator, func(portID, channelID string, nextSendSeq uint64) bool {
		ps := types.NewPacketSequence(portID, channelID, nextSendSeq)
		seqs = append(seqs, ps)
		return false
	})
	return seqs
}

// GetAllPacketRecvSeqs returns all stored next recv sequences.
func (k Keeper) GetAllPacketRecvSeqs(ctx sdk.Context) (seqs []types.PacketSequence) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(host.KeyNextSeqRecvPrefix))
	k.IteratePacketSequence(ctx, iterator, func(portID, channelID string, nextRecvSeq uint64) bool {
		ps := types.NewPacketSequence(portID, channelID, nextRecvSeq)
		seqs = append(seqs, ps)
		return false
	})
	return seqs
}

// GetAllPacketAckSeqs returns all stored next acknowledgements sequences.
func (k Keeper) GetAllPacketAckSeqs(ctx sdk.Context) (seqs []types.PacketSequence) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(host.KeyNextSeqAckPrefix))
	k.IteratePacketSequence(ctx, iterator, func(portID, channelID string, nextAckSeq uint64) bool {
		ps := types.NewPacketSequence(portID, channelID, nextAckSeq)
		seqs = append(seqs, ps)
		return false
	})
	return seqs
}

// IteratePacketCommitment provides an iterator over all PacketCommitment objects. For each
// packet commitment, cb will be called. If the cb returns true, the iterator will close
// and stop.
func (k Keeper) IteratePacketCommitment(ctx sdk.Context, cb func(portID, channelID string, sequence uint64, hash []byte) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(host.KeyPacketCommitmentPrefix))
	k.iterateHashes(ctx, iterator, cb)
}

// GetAllPacketCommitments returns all stored PacketCommitments objects.
func (k Keeper) GetAllPacketCommitments(ctx sdk.Context) (commitments []types.PacketState) {
	k.IteratePacketCommitment(ctx, func(portID, channelID string, sequence uint64, hash []byte) bool {
		pc := types.NewPacketState(portID, channelID, sequence, hash)
		commitments = append(commitments, pc)
		return false
	})
	return commitments
}

// IteratePacketCommitmentAtChannel provides an iterator over all PacketCommmitment objects
// at a specified channel. For each packet commitment, cb will be called. If the cb returns
// true, the iterator will close and stop.
func (k Keeper) IteratePacketCommitmentAtChannel(ctx sdk.Context, portID, channelID string, cb func(_, _ string, sequence uint64, hash []byte) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(host.PacketCommitmentPrefixPath(portID, channelID)))
	k.iterateHashes(ctx, iterator, cb)
}

// GetAllPacketCommitmentsAtChannel returns all stored PacketCommitments objects for a specified
// port ID and channel ID.
func (k Keeper) GetAllPacketCommitmentsAtChannel(ctx sdk.Context, portID, channelID string) (commitments []types.PacketState) {
	k.IteratePacketCommitmentAtChannel(ctx, portID, channelID, func(_, _ string, sequence uint64, hash []byte) bool {
		pc := types.NewPacketState(portID, channelID, sequence, hash)
		commitments = append(commitments, pc)
		return false
	})
	return commitments
}

// IteratePacketReceipt provides an iterator over all PacketReceipt objects. For each
// receipt, cb will be called. If the cb returns true, the iterator will close
// and stop.
func (k Keeper) IteratePacketReceipt(ctx sdk.Context, cb func(portID, channelID string, sequence uint64, receipt []byte) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(host.KeyPacketReceiptPrefix))
	k.iterateHashes(ctx, iterator, cb)
}

// GetAllPacketReceipts returns all stored PacketReceipt objects.
func (k Keeper) GetAllPacketReceipts(ctx sdk.Context) (receipts []types.PacketState) {
	k.IteratePacketReceipt(ctx, func(portID, channelID string, sequence uint64, receipt []byte) bool {
		packetReceipt := types.NewPacketState(portID, channelID, sequence, receipt)
		receipts = append(receipts, packetReceipt)
		return false
	})
	return receipts
}

// IteratePacketAcknowledgement provides an iterator over all PacketAcknowledgement objects. For each
// aknowledgement, cb will be called. If the cb returns true, the iterator will close
// and stop.
func (k Keeper) IteratePacketAcknowledgement(ctx sdk.Context, cb func(portID, channelID string, sequence uint64, hash []byte) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(host.KeyPacketAckPrefix))
	k.iterateHashes(ctx, iterator, cb)
}

// GetAllPacketAcks returns all stored PacketAcknowledgements objects.
func (k Keeper) GetAllPacketAcks(ctx sdk.Context) (acks []types.PacketState) {
	k.IteratePacketAcknowledgement(ctx, func(portID, channelID string, sequence uint64, ack []byte) bool {
		packetAck := types.NewPacketState(portID, channelID, sequence, ack)
		acks = append(acks, packetAck)
		return false
	})
	return acks
}

// common functionality for IteratePacketCommitment and IteratePacketAcknowledgement
func (k Keeper) iterateHashes(_ sdk.Context, iterator db.Iterator, cb func(portID, channelID string, sequence uint64, hash []byte) bool) {
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		keySplit := strings.Split(string(iterator.Key()), "/")
		portID := keySplit[2]
		channelID := keySplit[4]

		sequence, err := strconv.ParseUint(keySplit[len(keySplit)-1], 10, 64)
		if err != nil {
			panic(err)
		}

		if cb(portID, channelID, sequence, iterator.Value()) {
			break
		}
	}
}
