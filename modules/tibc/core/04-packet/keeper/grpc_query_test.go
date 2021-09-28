package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
	tibctesting "github.com/bianjieai/tibc-go/modules/tibc/testing"
)

func (suite *KeeperTestSuite) TestQueryPacketCommitment() {
	var (
		req           *types.QueryPacketCommitmentRequest
		expCommitment []byte
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{{
		"empty request",
		func() {
			req = nil
		},
		false,
	}, {
		"invalid source chain name",
		func() {
			req = &types.QueryPacketCommitmentRequest{
				SourceChain: "",
				DestChain:   "dest-chain",
				Sequence:    0,
			}
		},
		false,
	}, {
		"invalid destination chain name",
		func() {
			req = &types.QueryPacketCommitmentRequest{
				SourceChain: "source-chain",
				DestChain:   "",
				Sequence:    0,
			}
		},
		false,
	}, {
		"invalid sequence",
		func() {
			req = &types.QueryPacketCommitmentRequest{
				SourceChain: "source-chain",
				DestChain:   "dest-chain",
				Sequence:    0,
			}
		},
		false,
	}, {
		"dest chain not found",
		func() {
			req = &types.QueryPacketCommitmentRequest{
				SourceChain: "source-chain",
				DestChain:   "dest-chain",
				Sequence:    1,
			}
		},
		false,
	}, {
		"success",
		func() {
			path := tibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)

			expCommitment = []byte("hash")
			suite.chainA.App.TIBCKeeper.PacketKeeper.SetPacketCommitment(suite.chainA.GetContext(), path.EndpointA.ChainName, path.EndpointB.ChainName, 1, expCommitment)

			req = &types.QueryPacketCommitmentRequest{
				SourceChain: path.EndpointA.ChainName,
				DestChain:   path.EndpointB.ChainName,
				Sequence:    1,
			}
		},
		true,
	}}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset

			tc.malleate()
			ctx := sdk.WrapSDKContext(suite.chainA.GetContext())

			res, err := suite.chainA.QueryServer.PacketCommitment(ctx, req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(expCommitment, res.Commitment)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQueryPacketCommitments() {
	var (
		req            *types.QueryPacketCommitmentsRequest
		expCommitments = []*types.PacketState{}
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{{
		"empty request",
		func() {
			req = nil
		},
		false,
	}, {
		"invalid ID",
		func() {
			req = &types.QueryPacketCommitmentsRequest{
				SourceChain: "",
				DestChain:   "dest-chain",
			}
		},
		false,
	}, {
		"success, empty res",
		func() {
			expCommitments = []*types.PacketState{}

			req = &types.QueryPacketCommitmentsRequest{
				SourceChain: "source-chain",
				DestChain:   "dest-chain",
				Pagination: &query.PageRequest{
					Key:        nil,
					Limit:      2,
					CountTotal: true,
				},
			}
		},
		true,
	}, {
		"success",
		func() {
			path := tibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)

			expCommitments = make([]*types.PacketState, 9)

			for i := uint64(0); i < 9; i++ {
				commitment := types.NewPacketState(path.EndpointA.ChainName, path.EndpointB.ChainName, i, []byte(fmt.Sprintf("hash_%d", i)))
				suite.chainA.App.TIBCKeeper.PacketKeeper.SetPacketCommitment(suite.chainA.GetContext(), path.EndpointA.ChainName, path.EndpointB.ChainName, commitment.Sequence, commitment.Data)
				expCommitments[i] = &commitment
			}

			req = &types.QueryPacketCommitmentsRequest{
				SourceChain: path.EndpointA.ChainName,
				DestChain:   path.EndpointB.ChainName,
				Pagination: &query.PageRequest{
					Key:        nil,
					Limit:      11,
					CountTotal: true,
				},
			}
		},
		true,
	}}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset

			tc.malleate()
			ctx := sdk.WrapSDKContext(suite.chainA.GetContext())

			res, err := suite.chainA.QueryServer.PacketCommitments(ctx, req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(expCommitments, res.Commitments)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQueryPacketReceipt() {
	var (
		req         *types.QueryPacketReceiptRequest
		expReceived bool
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{{
		"empty request",
		func() {
			req = nil
		},
		false,
	}, {
		"invalid source chain name",
		func() {
			req = &types.QueryPacketReceiptRequest{
				SourceChain: "",
				DestChain:   "dest-chain",
				Sequence:    1,
			}
		},
		false,
	}, {
		"invalid destination chain name",
		func() {
			req = &types.QueryPacketReceiptRequest{
				SourceChain: "source-chain",
				DestChain:   "",
				Sequence:    1,
			}
		},
		false,
	}, {
		"invalid sequence",
		func() {
			req = &types.QueryPacketReceiptRequest{
				SourceChain: "source-chain",
				DestChain:   "dest-chain",
				Sequence:    0,
			}
		},
		false,
	}, {
		"success: receipt not found",
		func() {
			path := tibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)
			suite.chainA.App.TIBCKeeper.PacketKeeper.SetPacketReceipt(suite.chainA.GetContext(), path.EndpointA.ChainName, path.EndpointB.ChainName, 1)

			req = &types.QueryPacketReceiptRequest{
				SourceChain: path.EndpointA.ChainName,
				DestChain:   path.EndpointB.ChainName,
				Sequence:    3,
			}
			expReceived = false
		},
		true,
	}, {
		"success: receipt found",
		func() {
			path := tibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)

			suite.chainA.App.TIBCKeeper.PacketKeeper.SetPacketReceipt(suite.chainA.GetContext(), path.EndpointA.ChainName, path.EndpointB.ChainName, 1)

			req = &types.QueryPacketReceiptRequest{
				SourceChain: path.EndpointA.ChainName,
				DestChain:   path.EndpointB.ChainName,
				Sequence:    1,
			}
			expReceived = true
		},
		true,
	}}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset

			tc.malleate()
			ctx := sdk.WrapSDKContext(suite.chainA.GetContext())

			res, err := suite.chainA.QueryServer.PacketReceipt(ctx, req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(expReceived, res.Received)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQueryPacketAcknowledgement() {
	var (
		req    *types.QueryPacketAcknowledgementRequest
		expAck []byte
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{{
		"empty request",
		func() {
			req = nil
		},
		false,
	}, {
		"invalid source chain name",
		func() {
			req = &types.QueryPacketAcknowledgementRequest{
				SourceChain: "",
				DestChain:   "dest-chain",
				Sequence:    0,
			}
		},
		false,
	}, {
		"invalid destination chain name",
		func() {
			req = &types.QueryPacketAcknowledgementRequest{
				SourceChain: "source-chain",
				DestChain:   "",
				Sequence:    0,
			}
		},
		false,
	}, {
		"invalid sequence",
		func() {
			req = &types.QueryPacketAcknowledgementRequest{
				SourceChain: "source-chain",
				DestChain:   "dest-chain",
				Sequence:    0,
			}
		},
		false,
	}, {
		"dest chain not found",
		func() {
			req = &types.QueryPacketAcknowledgementRequest{
				SourceChain: "source-chain",
				DestChain:   "dest-chain",
				Sequence:    1,
			}
		},
		false,
	}, {
		"success",
		func() {
			path := tibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)

			expAck = []byte("hash")
			suite.chainA.App.TIBCKeeper.PacketKeeper.SetPacketAcknowledgement(suite.chainA.GetContext(), path.EndpointA.ChainName, path.EndpointB.ChainName, 1, expAck)

			req = &types.QueryPacketAcknowledgementRequest{
				SourceChain: path.EndpointA.ChainName,
				DestChain:   path.EndpointB.ChainName,
				Sequence:    1,
			}
		},
		true,
	}}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset

			tc.malleate()
			ctx := sdk.WrapSDKContext(suite.chainA.GetContext())

			res, err := suite.chainA.QueryServer.PacketAcknowledgement(ctx, req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(expAck, res.Acknowledgement)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQueryPacketAcknowledgements() {
	var (
		req                 *types.QueryPacketAcknowledgementsRequest
		expAcknowledgements = []*types.PacketState{}
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{{
		"empty request",
		func() {
			req = nil
		},
		false,
	}, {
		"invalid ID",
		func() {
			req = &types.QueryPacketAcknowledgementsRequest{
				SourceChain: "",
				DestChain:   "dest-chain",
			}
		},
		false,
	}, {
		"success, empty res",
		func() {
			expAcknowledgements = []*types.PacketState{}

			req = &types.QueryPacketAcknowledgementsRequest{
				SourceChain: "source-chain",
				DestChain:   "dest-chain",
				Pagination: &query.PageRequest{
					Key:        nil,
					Limit:      2,
					CountTotal: true,
				},
			}
		},
		true,
	}, {
		"success",
		func() {
			path := tibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)

			expAcknowledgements = make([]*types.PacketState, 9)

			for i := uint64(0); i < 9; i++ {
				ack := types.NewPacketState(path.EndpointA.ChainName, path.EndpointB.ChainName, i, []byte(fmt.Sprintf("hash_%d", i)))
				suite.chainA.App.TIBCKeeper.PacketKeeper.SetPacketAcknowledgement(suite.chainA.GetContext(), path.EndpointA.ChainName, path.EndpointB.ChainName, ack.Sequence, ack.Data)
				expAcknowledgements[i] = &ack
			}

			req = &types.QueryPacketAcknowledgementsRequest{
				SourceChain: path.EndpointA.ChainName,
				DestChain:   path.EndpointB.ChainName,
				Pagination: &query.PageRequest{
					Key:        nil,
					Limit:      11,
					CountTotal: true,
				},
			}
		},
		true,
	}}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset

			tc.malleate()
			ctx := sdk.WrapSDKContext(suite.chainA.GetContext())

			res, err := suite.chainA.QueryServer.PacketAcknowledgements(ctx, req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(expAcknowledgements, res.Acknowledgements)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQueryUnreceivedPackets() {
	var (
		req    *types.QueryUnreceivedPacketsRequest
		expSeq = []uint64{}
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{{
		"empty request",
		func() {
			req = nil
		},
		false,
	}, {
		"invalid source chain name",
		func() {
			req = &types.QueryUnreceivedPacketsRequest{
				SourceChain: "",
				DestChain:   "dest-chain",
			}
		},
		false,
	}, {
		"invalid destination chain name",
		func() {
			req = &types.QueryUnreceivedPacketsRequest{
				SourceChain: "source-chain",
				DestChain:   "",
			}
		},
		false,
	}, {
		"invalid seq",
		func() {
			req = &types.QueryUnreceivedPacketsRequest{
				SourceChain:               "source-chain",
				DestChain:                 "dest-chain",
				PacketCommitmentSequences: []uint64{0},
			}
		},
		false,
	}, {
		"basic success unreceived packet commitments",
		func() {
			path := tibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)

			// no ack exists

			expSeq = []uint64{1}
			req = &types.QueryUnreceivedPacketsRequest{
				SourceChain:               path.EndpointA.ChainName,
				DestChain:                 path.EndpointB.ChainName,
				PacketCommitmentSequences: []uint64{1},
			}
		},
		true,
	}, {
		"basic success unreceived packet commitments, nothing to relay",
		func() {
			path := tibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)

			suite.chainA.App.TIBCKeeper.PacketKeeper.SetPacketReceipt(suite.chainA.GetContext(), path.EndpointA.ChainName, path.EndpointB.ChainName, 1)

			expSeq = []uint64{}
			req = &types.QueryUnreceivedPacketsRequest{
				SourceChain:               path.EndpointA.ChainName,
				DestChain:                 path.EndpointB.ChainName,
				PacketCommitmentSequences: []uint64{1},
			}
		},
		true,
	}, {
		"success multiple unreceived packet commitments",
		func() {
			path := tibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)

			expSeq = []uint64{} // reset
			packetCommitments := []uint64{}

			// set packet receipt for every other sequence
			for seq := uint64(1); seq < 10; seq++ {
				packetCommitments = append(packetCommitments, seq)

				if seq%2 == 0 {
					suite.chainA.App.TIBCKeeper.PacketKeeper.SetPacketReceipt(suite.chainA.GetContext(), path.EndpointA.ChainName, path.EndpointB.ChainName, seq)
				} else {
					expSeq = append(expSeq, seq)
				}
			}

			req = &types.QueryUnreceivedPacketsRequest{
				SourceChain:               path.EndpointA.ChainName,
				DestChain:                 path.EndpointB.ChainName,
				PacketCommitmentSequences: packetCommitments,
			}
		},
		true,
	}}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset

			tc.malleate()
			ctx := sdk.WrapSDKContext(suite.chainA.GetContext())

			res, err := suite.chainA.QueryServer.UnreceivedPackets(ctx, req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(expSeq, res.Sequences)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQueryUnreceivedAcks() {
	var (
		req    *types.QueryUnreceivedAcksRequest
		expSeq = []uint64{}
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{{
		"empty request",
		func() {
			req = nil
		},
		false,
	}, {
		"invalid source chain name",
		func() {
			req = &types.QueryUnreceivedAcksRequest{
				SourceChain: "",
				DestChain:   "dest-chain",
			}
		},
		false,
	}, {
		"invalid destination chain name",
		func() {
			req = &types.QueryUnreceivedAcksRequest{
				SourceChain: "source-chain",
				DestChain:   "",
			}
		},
		false,
	}, {
		"invalid seq",
		func() {
			req = &types.QueryUnreceivedAcksRequest{
				SourceChain:        "source-chain",
				DestChain:          "dest-chain",
				PacketAckSequences: []uint64{0},
			}
		},
		false,
	}, {
		"basic success unreceived packet acks",
		func() {
			path := tibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)

			suite.chainA.App.TIBCKeeper.PacketKeeper.SetPacketCommitment(suite.chainA.GetContext(), path.EndpointA.ChainName, path.EndpointB.ChainName, 1, []byte("commitment"))

			expSeq = []uint64{1}
			req = &types.QueryUnreceivedAcksRequest{
				SourceChain:        path.EndpointA.ChainName,
				DestChain:          path.EndpointB.ChainName,
				PacketAckSequences: []uint64{1},
			}
		},
		true,
	}, {
		"basic success unreceived packet acknowledgements, nothing to relay",
		func() {
			path := tibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)

			expSeq = []uint64{}
			req = &types.QueryUnreceivedAcksRequest{
				SourceChain:        path.EndpointA.ChainName,
				DestChain:          path.EndpointB.ChainName,
				PacketAckSequences: []uint64{1},
			}
		},
		true,
	}, {
		"success multiple unreceived packet acknowledgements",
		func() {
			path := tibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)

			expSeq = []uint64{} // reset
			packetAcks := []uint64{}

			// set packet commitment for every other sequence
			for seq := uint64(1); seq < 10; seq++ {
				packetAcks = append(packetAcks, seq)

				if seq%2 == 0 {
					suite.chainA.App.TIBCKeeper.PacketKeeper.SetPacketCommitment(suite.chainA.GetContext(), path.EndpointA.ChainName, path.EndpointB.ChainName, seq, []byte("commitement"))
					expSeq = append(expSeq, seq)
				}
			}

			req = &types.QueryUnreceivedAcksRequest{
				SourceChain:        path.EndpointA.ChainName,
				DestChain:          path.EndpointB.ChainName,
				PacketAckSequences: packetAcks,
			}
		},
		true,
	}}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset

			tc.malleate()
			ctx := sdk.WrapSDKContext(suite.chainA.GetContext())

			res, err := suite.chainA.QueryServer.UnreceivedAcks(ctx, req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(expSeq, res.Sequences)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQueryCleanPacketCommitment() {
	var req *types.QueryCleanPacketCommitmentRequest

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{{
		"empty request",
		func() {
			req = nil
		},
		false,
	}, {
		"invalid source chain name",
		func() {
			req = &types.QueryCleanPacketCommitmentRequest{
				SourceChain: "",
				DestChain:   "dest-chain",
			}
		},
		false,
	}, {
		"invalid destination chain name",
		func() {
			req = &types.QueryCleanPacketCommitmentRequest{
				SourceChain: "source-chain",
				DestChain:   "",
			}
		},
		false,
	}, {
		"dest chain not found",
		func() {
			req = &types.QueryCleanPacketCommitmentRequest{
				SourceChain: "source-chain",
				DestChain:   "dest-chain",
			}
		},
		false,
	}, {
		"success",
		func() {
			path := tibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)
			suite.chainA.App.TIBCKeeper.PacketKeeper.SetCleanPacketCommitment(suite.chainA.GetContext(), path.EndpointA.ChainName, path.EndpointB.ChainName, 1)

			req = &types.QueryCleanPacketCommitmentRequest{
				SourceChain: path.EndpointA.ChainName,
				DestChain:   path.EndpointB.ChainName,
			}
		},
		true,
	}}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset

			tc.malleate()
			ctx := sdk.WrapSDKContext(suite.chainA.GetContext())

			res, err := suite.chainA.QueryServer.CleanPacketCommitment(ctx, req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(sdk.Uint64ToBigEndian(1), res.Commitment)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
