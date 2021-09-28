package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
	tibctesting "github.com/bianjieai/tibc-go/modules/tibc/testing"
)

// KeeperTestSuite is a testing suite to test keeper functions.
type KeeperTestSuite struct {
	suite.Suite

	coordinator *tibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA *tibctesting.TestChain
	chainB *tibctesting.TestChain
}

// TestKeeperTestSuite runs all the tests within this package.
func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// SetupTest creates a coordinator with 2 test chains.
func (suite *KeeperTestSuite) SetupTest() {
	suite.coordinator = tibctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(tibctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(tibctesting.GetChainID(1))
	// commit some blocks so that QueryProof returns valid proof (cannot return valid query if height <= 1)
	suite.coordinator.CommitNBlocks(suite.chainA, 2)
	suite.coordinator.CommitNBlocks(suite.chainB, 2)
}

// TestGetAllSequences sets all packet sequences for two different channels on chain A and
// tests their retrieval.
func (suite KeeperTestSuite) TestGetAllSequences() {
	path := tibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupClients(path)

	seq1 := types.NewPacketSequence(path.EndpointA.ChainName, path.EndpointB.ChainName, 1)
	seq2 := types.NewPacketSequence(path.EndpointA.ChainName, path.EndpointB.ChainName, 2)

	// seq1 should be overwritten by seq2
	expSeqs := []types.PacketSequence{seq2}

	ctxA := suite.chainA.GetContext()

	for _, seq := range []types.PacketSequence{seq1, seq2} {
		suite.chainA.App.TIBCKeeper.PacketKeeper.SetNextSequenceSend(ctxA, path.EndpointA.ChainName, path.EndpointB.ChainName, seq.Sequence)
	}

	sendSeqs := suite.chainA.App.TIBCKeeper.PacketKeeper.GetAllPacketSendSeqs(ctxA)
	suite.Len(sendSeqs, 1)

	suite.Equal(expSeqs, sendSeqs)
}

// TestGetAllPacketState creates a set of acks, packet commitments, and receipts on two different
// channels on chain A and tests their retrieval.
func (suite KeeperTestSuite) TestGetAllPacketState() {
	path := tibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupClients(path)

	// channel 0 acks
	ack1 := types.NewPacketState(path.EndpointA.ChainName, path.EndpointB.ChainName, 1, []byte("ack"))
	ack2 := types.NewPacketState(path.EndpointA.ChainName, path.EndpointB.ChainName, 2, []byte("ack"))

	// duplicate ack
	ack2dup := types.NewPacketState(path.EndpointA.ChainName, path.EndpointB.ChainName, 2, []byte("ack"))

	// create channel 0 receipts
	receipt := string([]byte{byte(1)})
	rec1 := types.NewPacketState(path.EndpointA.ChainName, path.EndpointB.ChainName, 1, []byte(receipt))
	rec2 := types.NewPacketState(path.EndpointA.ChainName, path.EndpointB.ChainName, 2, []byte(receipt))

	// channel 0 packet commitments
	comm1 := types.NewPacketState(path.EndpointA.ChainName, path.EndpointB.ChainName, 1, []byte("hash"))
	comm2 := types.NewPacketState(path.EndpointA.ChainName, path.EndpointB.ChainName, 2, []byte("hash"))

	expAcks := []types.PacketState{ack1, ack2}
	expReceipts := []types.PacketState{rec1, rec2}
	expCommitments := []types.PacketState{comm1, comm2}

	ctxA := suite.chainA.GetContext()

	// set acknowledgements
	for _, ack := range []types.PacketState{ack1, ack2, ack2dup} {
		suite.chainA.App.TIBCKeeper.PacketKeeper.SetPacketAcknowledgement(ctxA, path.EndpointA.ChainName, path.EndpointB.ChainName, ack.Sequence, ack.Data)
	}

	// set packet receipts
	for _, rec := range expReceipts {
		suite.chainA.App.TIBCKeeper.PacketKeeper.SetPacketReceipt(ctxA, path.EndpointA.ChainName, path.EndpointB.ChainName, rec.Sequence)
	}

	// set packet commitments
	for _, comm := range expCommitments {
		suite.chainA.App.TIBCKeeper.PacketKeeper.SetPacketCommitment(ctxA, path.EndpointA.ChainName, path.EndpointB.ChainName, comm.Sequence, comm.Data)
	}

	acks := suite.chainA.App.TIBCKeeper.PacketKeeper.GetAllPacketAcks(ctxA)
	receipts := suite.chainA.App.TIBCKeeper.PacketKeeper.GetAllPacketReceipts(ctxA)
	commitments := suite.chainA.App.TIBCKeeper.PacketKeeper.GetAllPacketCommitments(ctxA)

	suite.Require().Len(acks, len(expAcks))
	suite.Require().Len(commitments, len(expCommitments))
	suite.Require().Len(receipts, len(expReceipts))

	suite.Require().Equal(expAcks, acks)
	suite.Require().Equal(expReceipts, receipts)
	suite.Require().Equal(expCommitments, commitments)
}

// TestSetSequence verifies that the keeper correctly sets the sequence counters.
func (suite *KeeperTestSuite) TestSetSequence() {
	path := tibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupClients(path)

	ctxA := suite.chainA.GetContext()
	one := uint64(1)

	// initialized channel has next send seq of 1
	seq := suite.chainA.App.TIBCKeeper.PacketKeeper.GetNextSequenceSend(ctxA, path.EndpointA.ChainName, path.EndpointB.ChainName)
	suite.Equal(one, seq)

	nextSeqSend := uint64(10)
	suite.chainA.App.TIBCKeeper.PacketKeeper.SetNextSequenceSend(ctxA, path.EndpointA.ChainName, path.EndpointB.ChainName, nextSeqSend)

	storedNextSeqSend := suite.chainA.App.TIBCKeeper.PacketKeeper.GetNextSequenceSend(ctxA, path.EndpointA.ChainName, path.EndpointB.ChainName)
	suite.Equal(nextSeqSend, storedNextSeqSend)
}

// TestGetAllPacketCommitmentsAtChannel verifies that the keeper returns all stored packet
// commitments for a specific channel. The test will store consecutive commitments up to the
// value of "seq" and then add non-consecutive up to the value of "maxSeq". A final commitment
// with the value maxSeq + 1 is set on a different channel.
func (suite *KeeperTestSuite) TestGetAllPacketCommitmentsAtChannel() {
	path := tibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupClients(path)

	ctxA := suite.chainA.GetContext()
	expectedSeqs := make(map[uint64]bool)
	hash := []byte("commitment")

	seq := uint64(15)
	maxSeq := uint64(25)
	suite.Require().Greater(maxSeq, seq)

	// create consecutive commitments
	for i := uint64(1); i < seq; i++ {
		suite.chainA.App.TIBCKeeper.PacketKeeper.SetPacketCommitment(ctxA, path.EndpointA.ChainName, path.EndpointB.ChainName, i, hash)
		expectedSeqs[i] = true
	}

	// add non-consecutive commitments
	for i := seq; i < maxSeq; i += 2 {
		suite.chainA.App.TIBCKeeper.PacketKeeper.SetPacketCommitment(ctxA, path.EndpointA.ChainName, path.EndpointB.ChainName, i, hash)
		expectedSeqs[i] = true
	}

	// add sequence on different channel/port
	suite.chainA.App.TIBCKeeper.PacketKeeper.SetPacketCommitment(ctxA, path.EndpointA.ChainName, "EndpointBChainName", maxSeq+1, hash)

	commitments := suite.chainA.App.TIBCKeeper.PacketKeeper.GetAllPacketCommitmentsAtChannel(ctxA, path.EndpointA.ChainName, path.EndpointB.ChainName)

	suite.Equal(len(expectedSeqs), len(commitments))
	// ensure above for loops occurred
	suite.NotEqual(0, len(commitments))

	// verify that all the packet commitments were stored
	for _, packet := range commitments {
		suite.True(expectedSeqs[packet.Sequence])
		suite.Equal(path.EndpointA.ChainName, packet.SourceChain)
		suite.Equal(path.EndpointB.ChainName, packet.DestinationChain)
		suite.Equal(hash, packet.Data)

		// prevent duplicates from passing checks
		expectedSeqs[packet.Sequence] = false
	}
}

// TestSetPacketAcknowledgement verifies that packet acknowledgements are correctly
// set in the keeper.
func (suite *KeeperTestSuite) TestSetPacketAcknowledgement() {
	path := tibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupClients(path)

	ctxA := suite.chainA.GetContext()
	seq := uint64(10)

	storedAckHash, found := suite.chainA.App.TIBCKeeper.PacketKeeper.GetPacketAcknowledgement(ctxA, path.EndpointA.ChainName, path.EndpointB.ChainName, seq)
	suite.Require().False(found)
	suite.Require().Nil(storedAckHash)

	ackHash := []byte("ackhash")
	suite.chainA.App.TIBCKeeper.PacketKeeper.SetPacketAcknowledgement(ctxA, path.EndpointA.ChainName, path.EndpointB.ChainName, seq, ackHash)

	storedAckHash, found = suite.chainA.App.TIBCKeeper.PacketKeeper.GetPacketAcknowledgement(ctxA, path.EndpointA.ChainName, path.EndpointB.ChainName, seq)
	suite.Require().True(found)
	suite.Require().Equal(ackHash, storedAckHash)
	suite.Require().True(suite.chainA.App.TIBCKeeper.PacketKeeper.HasPacketAcknowledgement(ctxA, path.EndpointA.ChainName, path.EndpointB.ChainName, seq))
}
