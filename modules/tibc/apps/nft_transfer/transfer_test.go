package nft_transfer_test

import (
	"github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/types"
	"github.com/bianjieai/tibc-go/modules/tibc/testing"
	nfttypes "github.com/irisnet/irismod/modules/nft/types"
	"github.com/stretchr/testify/suite"
	"testing"
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
func (suite *TransferTestSuite) TestHandleMsgTransfer() {
	// setup between chainA and chainB
	suite.coordinator.Setup(suite.chainA, suite.chainB)

	// issue denom
	issueDenomMsg :=  nfttypes.NewMsgIssueDenom("mobile", "mobile-name", "",
		suite.chainA.SenderAccount.GetAddress().String())
	_, _ = suite.chainA.SendMsgs(issueDenomMsg)

	// mint nft 
	mintNftMsg := nfttypes.NewMsgMintNFT("xiaomi", "mobile", "",
		"", "", suite.chainA.SenderAccount.GetAddress().String(),
		suite.chainA.SenderAccount.GetAddress().String())
	_, _ = suite.chainA.SendMsgs(mintNftMsg)

	dd, _:= suite.chainA.App.NftKeeper.GetDenom(suite.chainA.GetContext(), "mobile")


	// send nft from A To B
	msg := types.NewMsgNftTransfer(dd.Id, "xiaomi", "",
		suite.chainA.SenderAccount.GetAddress().String(),
		suite.chainB.SenderAccount.GetAddress().String(),
		true, suite.chainB.ChainID, "")

	_, err := suite.chainA.SendMsgs(msg)
	suite.Require().NoError(err) // message committed

}


func TestTransferTestSuite(t *testing.T) {
	suite.Run(t, new(TransferTestSuite))
}




