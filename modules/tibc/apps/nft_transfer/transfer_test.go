package nft_transfer_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	ibctesting "github.com/bianjieai/tibc-go/modules/tibc/testing"
)

const (
	isCheckTx = false
)

type TransferTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain
	chainC *ibctesting.TestChain
}

func (suite *TransferTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 3)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.chainC = suite.coordinator.GetChain(ibctesting.GetChainID(2))
}

/*
nft
A->B B->C
*/
// func (suite *TransferTestSuite) TestHandleMsgTransfer() {
// 	// setup between chainA and chainB

// 	path := ibctesting.NewPath(suite.chainA, suite.chainB)

// 	suite.coordinator.SetupClients(path)

// 	// issue denom
// 	issueDenomMsg := nfttypes.NewMsgIssueDenom("mobile", "mobile-name", "",
// 		suite.chainA.SenderAccount.GetAddress().String(), "", false, false)
// 	_, _ = suite.chainA.SendMsgs(issueDenomMsg)

// 	// mint nft
// 	mintNftMsg := nfttypes.NewMsgMintNFT("xiaomi", "mobile", "",
// 		"", "", suite.chainA.SenderAccount.GetAddress().String(),
// 		suite.chainA.SenderAccount.GetAddress().String())
// 	_, _ = suite.chainA.SendMsgs(mintNftMsg)

// 	dd, _ := suite.chainA.App.NftKeeper.GetDenom(suite.chainA.GetContext(), "mobile")

// 	// send nft from A To B
// 	msg := types.NewMsgNftTransfer(dd.Id, "xiaomi",
// 		suite.chainA.SenderAccount.GetAddress().String(),
// 		suite.chainB.SenderAccount.GetAddress().String(),
// 		suite.chainB.ChainID, "")

// 	_, err := suite.chainA.SendMsgs(msg)
// 	suite.Require().NoError(err) // message committed
// 	//// relay send
// 	//NonfungibleTokenPacket := types.NewNonFungibleTokenPacketData("mobile", "xiaomi",
// 	//	"", suite.chainA.SenderAccount.GetAddress().String(),
// 	//	suite.chainB.SenderAccount.GetAddress().String(), true,
// 	//)
// 	//packet := packettypes.NewPacket(NonfungibleTokenPacket.GetBytes(), 1,
// 	//	path.EndpointA.ChainName, path.EndpointB.ChainName, "", string(routingtypes.NFT))
// 	//
// 	//ack := packettypes.NewResultAcknowledgement([]byte{byte(1)})
// 	//err = path.RelayPacket(packet, ack.GetBytes())
// 	//suite.Require().NoError(err) // relay committed

// }

func TestTransferTestSuite(t *testing.T) {
	suite.Run(t, new(TransferTestSuite))
}
