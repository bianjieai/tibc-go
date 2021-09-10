package keeper_test

import (
	"fmt"

	"github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
	ibctesting "github.com/bianjieai/tibc-go/modules/tibc/testing"
	ibcmock "github.com/bianjieai/tibc-go/modules/tibc/testing/mock"
)

type testCase = struct {
	msg      string
	malleate func()
	expPass  bool
}

var (
	validPacketData = []byte("VALID PACKET DATA")

	relayChain = ""
)

// TestSendPacket tests SendPacket from chainA to chainB
func (suite *KeeperTestSuite) TestSendPacket() {
	var (
		packet exported.PacketI
	)

	testCases := []testCase{
		{"success: UNORDERED channel", func() {
			path := ibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)

			packet = types.NewPacket(validPacketData, 1, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, ibctesting.MockPort)
		}, true},
		{"sending packet out of order on UNORDERED channel", func() {
			// setup creates an unordered channel
			path := ibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)
			packet = types.NewPacket(validPacketData, 5, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, ibctesting.MockPort)
		}, false},
		// {"packet basic validation failed, empty packet data", func() {
		// 	path := ibctesting.NewPath(suite.chainA, suite.chainB)
		// 	suite.coordinator.SetupClients(path)
		// 	packet = types.NewPacket([]byte{}, 1, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, ibctesting.MockPort)
		// }, false},
		// {"port not found", func() {
		// 	// use wrong port
		// 	path := ibctesting.NewPath(suite.chainA, suite.chainB)
		// 	suite.coordinator.SetupClients(path)
		// 	packet = types.NewPacket(validPacketData, 1, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, ibctesting.InvalidID)
		// }, false},
		{"client state not found", func() {
			path := ibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)

			packet = types.NewPacket(validPacketData, 1, path.EndpointA.ChainName, ibctesting.InvalidID, relayChain, ibctesting.MockPort)
		}, false},
		{"next sequence wrong", func() {
			path := ibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)
			packet = types.NewPacket(validPacketData, 1, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, ibctesting.MockPort)
			suite.chainA.App.TIBCKeeper.PacketKeeper.SetNextSequenceSend(suite.chainA.GetContext(), path.EndpointA.ChainName, path.EndpointB.ChainName, 5)
		}, false},
	}

	for i, tc := range testCases {
		tc := tc
		suite.Run(fmt.Sprintf("Case %s, %d/%d tests", tc.msg, i, len(testCases)), func() {
			suite.SetupTest() // reset

			tc.malleate()

			err := suite.chainA.App.TIBCKeeper.PacketKeeper.SendPacket(suite.chainA.GetContext(), packet)

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}

}

// TestRecvPacket test RecvPacket on chainB. Since packet commitment verification will always
// occur last (resource instensive), only tests expected to succeed and packet commitment
// verification tests need to simulate sending a packet from chainA to chainB.
func (suite *KeeperTestSuite) TestRecvPacket() {
	var (
		packet exported.PacketI
	)

	testCases := []testCase{
		{"success: ORDERED channel", func() {
			path := ibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)

			packet = types.NewPacket(validPacketData, 1, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, ibctesting.MockPort)
			err := path.EndpointA.SendPacket(packet)
			suite.Require().NoError(err)
		}, true},
		{"success with out of order packet: UNORDERED channel", func() {
			// setup uses an UNORDERED channel
			path := ibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)
			packet = types.NewPacket(validPacketData, 1, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, ibctesting.MockPort)

			// send 2 packets
			err := path.EndpointA.SendPacket(packet)
			suite.Require().NoError(err)
			// set sequence to 2
			packet = types.NewPacket(validPacketData, 2, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, ibctesting.MockPort)
			err = path.EndpointA.SendPacket(packet)
			suite.Require().NoError(err)
		}, true},
		{"port not found", func() {
			// use wrong port
			path := ibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)
			packet = types.NewPacket(validPacketData, 1, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, ibctesting.InvalidID)
		}, false},
		{"receipt already stored", func() {
			path := ibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)
			packet = types.NewPacket(validPacketData, 1, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, ibctesting.MockPort)
			err := path.EndpointA.SendPacket(packet)
			suite.Require().NoError(err)
			suite.chainB.App.TIBCKeeper.PacketKeeper.SetPacketReceipt(suite.chainB.GetContext(), path.EndpointA.ChainName, path.EndpointB.ChainName, 1)
		}, false},
		{"validation failed", func() {
			// packet commitment not set resulting in invalid proof
			path := ibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)
			packet = types.NewPacket(validPacketData, 1, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, ibctesting.InvalidID)
		}, false},
	}

	for i, tc := range testCases {
		tc := tc
		suite.Run(fmt.Sprintf("Case %s, %d/%d tests", tc.msg, i, len(testCases)), func() {
			suite.SetupTest() // reset
			tc.malleate()

			// get proof of packet commitment from chainA
			packetKey := host.PacketCommitmentKey(packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence())
			proof, proofHeight := suite.chainA.QueryProof(packetKey)

			err := suite.chainB.App.TIBCKeeper.PacketKeeper.RecvPacket(suite.chainB.GetContext(), packet, proof, proofHeight)

			if tc.expPass {
				suite.Require().NoError(err)

				receipt, receiptStored := suite.chainB.App.TIBCKeeper.PacketKeeper.GetPacketReceipt(suite.chainB.GetContext(), packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence())

				suite.Require().True(receiptStored, "packet receipt not stored after RecvPacket in UNORDERED channel")
				suite.Require().Equal(string([]byte{byte(1)}), receipt, "packet receipt is not empty string")

			} else {
				suite.Require().Error(err)
			}
		})
	}

}

func (suite *KeeperTestSuite) TestWriteAcknowledgement() {
	var (
		ack    []byte
		packet exported.PacketI
	)

	testCases := []testCase{
		{
			"success",
			func() {
				path := ibctesting.NewPath(suite.chainA, suite.chainB)
				suite.coordinator.SetupClients(path)
				packet = types.NewPacket(validPacketData, 1, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, ibctesting.MockPort)
				ack = ibctesting.TestHash
			},
			true,
		},
		// {"port not found", func() {
		// 	// use wrong port naming
		// 	path := ibctesting.NewPath(suite.chainA, suite.chainB)
		// 	suite.coordinator.SetupClients(path)
		// 	packet = types.NewPacket(validPacketData, 1, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, ibctesting.InvalidID)
		// 	ack = ibctesting.TestHash
		// }, false},
		{
			"no-op, already acked",
			func() {
				path := ibctesting.NewPath(suite.chainA, suite.chainB)
				suite.coordinator.SetupClients(path)
				packet = types.NewPacket(validPacketData, 1, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, ibctesting.MockPort)
				ack = ibctesting.TestHash
				suite.chainB.App.TIBCKeeper.PacketKeeper.SetPacketAcknowledgement(suite.chainB.GetContext(), packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence(), ack)
			},
			false,
		},
		{
			"empty acknowledgement",
			func() {
				path := ibctesting.NewPath(suite.chainA, suite.chainB)
				suite.coordinator.SetupClients(path)
				packet = types.NewPacket(validPacketData, 1, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, ibctesting.MockPort)
				ack = nil
			},
			false,
		},
	}
	for i, tc := range testCases {
		tc := tc
		suite.Run(fmt.Sprintf("Case %s, %d/%d tests", tc.msg, i, len(testCases)), func() {
			suite.SetupTest() // reset

			tc.malleate()

			err := suite.chainB.App.TIBCKeeper.PacketKeeper.WriteAcknowledgement(suite.chainB.GetContext(), packet, ack)

			if tc.expPass {
				suite.NoError(err, "Invalid Case %d passed: %s", i, tc.msg)
			} else {
				suite.Error(err, "Case %d failed: %s", i, tc.msg)
			}
		})
	}
}

// TestAcknowledgePacket tests the call AcknowledgePacket on chainA.
func (suite *KeeperTestSuite) TestAcknowledgePacket() {
	var (
		packet types.Packet
		ack    = ibcmock.MockAcknowledgement
	)

	testCases := []testCase{
		{"success", func() {
			// setup
			path := ibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)
			packet = types.NewPacket(validPacketData, 1, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, ibctesting.MockPort)

			// create packet commitment
			err := path.EndpointA.SendPacket(packet)
			suite.Require().NoError(err)

			// create packet receipt and acknowledgement
			err = path.EndpointB.RecvPacket(packet)
			suite.Require().NoError(err)
		}, true},
		{"port not found", func() {
			// use wrong port naming
			path := ibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)
			packet = types.NewPacket(validPacketData, 1, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, ibctesting.InvalidID)
		}, false},
		{"packet hasn't been sent", func() {
			// packet commitment never written
			path := ibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)
			packet = types.NewPacket(validPacketData, 1, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, ibctesting.MockPort)
		}, false},
	}

	for i, tc := range testCases {
		tc := tc
		suite.Run(fmt.Sprintf("Case %s, %d/%d tests", tc.msg, i, len(testCases)), func() {
			suite.SetupTest() // reset
			tc.malleate()

			packetKey := host.PacketAcknowledgementKey(packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence())
			proof, proofHeight := suite.chainB.QueryProof(packetKey)

			err := suite.chainA.App.TIBCKeeper.PacketKeeper.AcknowledgePacket(suite.chainA.GetContext(), packet, ack, proof, proofHeight)
			pc := suite.chainA.App.TIBCKeeper.PacketKeeper.GetPacketCommitment(suite.chainA.GetContext(), packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence())

			if tc.expPass {
				suite.NoError(err, "Case %d failed: %s", i, tc.msg)
				suite.Nil(pc)
			} else {
				suite.Error(err, "Invalid Case %d passed: %s", i, tc.msg)
			}
		})
	}
}
