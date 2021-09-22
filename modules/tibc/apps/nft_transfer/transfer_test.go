package nft_transfer_test

import (
	"fmt"
	"testing"

	nfttypes "github.com/irisnet/irismod/modules/nft/types"

	"github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/types"
	packettypes "github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
	routingtypes "github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"

	"github.com/stretchr/testify/suite"

	tibctesting "github.com/bianjieai/tibc-go/modules/tibc/testing"
)

const (
	isCheckTx = false
)

type TransferTestSuite struct {
	suite.Suite

	coordinator *tibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA *tibctesting.TestChain
	chainB *tibctesting.TestChain
	chainC *tibctesting.TestChain
}

func (suite *TransferTestSuite) SetupTest() {
	suite.coordinator = tibctesting.NewCoordinator(suite.T(), 3)
	suite.chainA = suite.coordinator.GetChain(tibctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(tibctesting.GetChainID(1))
	suite.chainC = suite.coordinator.GetChain(tibctesting.GetChainID(2))
}

/*
nft
A->B B->C
*/
func (suite *TransferTestSuite) TestHandleMsgTransfer() {
	// setup between chainA and chainB

	path := tibctesting.NewPath(suite.chainA, suite.chainB)

	suite.coordinator.SetupClients(path)

	// issue denom
	issueDenomMsg := nfttypes.NewMsgIssueDenom("mobile", "mobile-name", "",
		suite.chainA.SenderAccount.GetAddress().String(), "", false, false)
	_, _ = suite.chainA.SendMsgs(issueDenomMsg)

	// mint nft
	mintNftMsg := nfttypes.NewMsgMintNFT("xiaomi", "mobile", "",
		"", "", suite.chainA.SenderAccount.GetAddress().String(),
		suite.chainA.SenderAccount.GetAddress().String())
	_, _ = suite.chainA.SendMsgs(mintNftMsg)

	dd, _ := suite.chainA.App.NftKeeper.GetDenom(suite.chainA.GetContext(), "mobile")

	// send nft from A To B
	msg := types.NewMsgNftTransfer(dd.Id, "xiaomi",
		suite.chainA.SenderAccount.GetAddress().String(),
		suite.chainB.SenderAccount.GetAddress().String(),
		suite.chainB.ChainID, "")

	_, err := suite.chainA.SendMsgs(msg)
	suite.Require().NoError(err) // message committed
	//// relay send
	NonfungibleTokenPacket := types.NewNonFungibleTokenPacketData("mobile", "xiaomi",
		"", suite.chainA.SenderAccount.GetAddress().String(),
		suite.chainB.SenderAccount.GetAddress().String(), true,
	)
	packet := packettypes.NewPacket(NonfungibleTokenPacket.GetBytes(), 1,
		path.EndpointA.ChainName, path.EndpointB.ChainName, "", string(routingtypes.NFT))

	ack := packettypes.NewResultAcknowledgement([]byte{byte(1)})
	err = path.RelayPacket(packet, ack.GetBytes())
	suite.Require().NoError(err) // relay committed

	// check that voucher exists on chain B
	// denomID : tibc/nft/chainA.chainID/mobile
	classInchainB := "tibc/nft" + "/" + suite.chainA.ChainID + "/" + "mobile"
	nft, _ := suite.chainB.App.NftKeeper.GetNFT(suite.chainB.GetContext(), classInchainB, "xiaomi")
	suite.Require().Equal("xiaomi", nft.GetID())

	// setup between chainB to chainC
	pathBtoC := tibctesting.NewPath(suite.chainB, suite.chainC)
	suite.coordinator.SetupClients(pathBtoC)

	// send nft from chainB to chainC
	msgfromBToC := types.NewMsgNftTransfer(classInchainB, "xiaomi",
		suite.chainB.SenderAccount.GetAddress().String(),
		suite.chainC.SenderAccount.GetAddress().String(),
		suite.chainC.ChainID, "")

	_, err1 := suite.chainB.SendMsgs(msgfromBToC)
	suite.Require().NoError(err1) // message committed

	// relay send
	nftPacketFromBToC := types.NewNonFungibleTokenPacketData(classInchainB, "xiaomi",
		"", suite.chainB.SenderAccount.GetAddress().String(),
		suite.chainC.SenderAccount.GetAddress().String(), true,
	)
	packetFromBToC := packettypes.NewPacket(nftPacketFromBToC.GetBytes(), 1,
		pathBtoC.EndpointA.ChainName, pathBtoC.EndpointB.ChainName, "", string(routingtypes.NFT))

	ack1 := packettypes.NewResultAcknowledgement([]byte{byte(1)})
	err = pathBtoC.RelayPacket(packetFromBToC, ack1.GetBytes())
	suite.Require().NoError(err) // relay committed

	// check that voucher exists on chain C
	// denomID : tibc/nft/chainA.chainID/chainB.chainID/mobile
	classInchainC := "tibc/nft" + "/" + suite.chainA.ChainID + "/" + suite.chainB.ChainID + "/" + "mobile"
	nftInC, _ := suite.chainC.App.NftKeeper.GetNFT(suite.chainC.GetContext(), classInchainC, "xiaomi")
	suite.Require().Equal("xiaomi", nftInC.GetID())

	// send nft  from chainC back to chainB
	msgfromCToB := types.NewMsgNftTransfer(classInchainC, "xiaomi",
		suite.chainC.SenderAccount.GetAddress().String(),
		suite.chainB.SenderAccount.GetAddress().String(),
		suite.chainB.ChainID, "")

	_, err2 := suite.chainC.SendMsgs(msgfromCToB)
	suite.Require().NoError(err2) // message committed

	// relay send
	nftPacket := types.NewNonFungibleTokenPacketData(classInchainC, "xiaomi",
		"", suite.chainC.SenderAccount.GetAddress().String(),
		suite.chainB.SenderAccount.GetAddress().String(), false,
	)
	packetFromCToB := packettypes.NewPacket(nftPacket.GetBytes(), 1,
		pathBtoC.EndpointB.ChainName, pathBtoC.EndpointA.ChainName, "", string(routingtypes.NFT))

	ack2 := packettypes.NewResultAcknowledgement([]byte{byte(1)})
	err = pathBtoC.RelayPacket(packetFromCToB, ack2.GetBytes())
	suite.Require().NoError(err) // relay committed

	// send nft  from chainB back to chainA
	msgFromBToA := types.NewMsgNftTransfer(classInchainB, "xiaomi",
		suite.chainB.SenderAccount.GetAddress().String(),
		suite.chainA.SenderAccount.GetAddress().String(),
		suite.chainA.ChainID, "")

	_, err = suite.chainB.SendMsgs(msgFromBToA)
	suite.Require().NoError(err) // message committed

	// relay send
	NonfungibleTokenPacket = types.NewNonFungibleTokenPacketData(classInchainB, "xiaomi",
		"", suite.chainB.SenderAccount.GetAddress().String(),
		suite.chainA.SenderAccount.GetAddress().String(), false,
	)
	packet = packettypes.NewPacket(NonfungibleTokenPacket.GetBytes(), 1,
		path.EndpointB.ChainName, path.EndpointA.ChainName, "", string(routingtypes.NFT))

	ack = packettypes.NewResultAcknowledgement([]byte{byte(1)})
	err = path.RelayPacket(packet, ack.GetBytes())
	suite.Require().NoError(err) // relay committed

	// Query whether there are corresponding nfts for C, B, and A respectively
	/* just do  A->B->C
	denom found in A: mobile
	nft found in A: xiaomi

	denom found in B: tibc/nft/testchain0/mobile
	nft found in B: xiaomi

	denom found in C: tibc/nft/testchain0/testchain1/mobile
	nft found in C: xiaomi
	*/

	/* do A->B->C  then do C->B->A
	denom found in A: mobile
	nft found in A: xiaomi

	denom found in B: tibc/nft/testchain0/mobile
	nft not found in B

	denom found in C: tibc/nft/testchain0/testchain1/mobile
	nft not found in C

	*/
	// query A
	denomInA, found := suite.chainA.App.NftKeeper.GetDenom(suite.chainA.GetContext(), "mobile")
	if found {
		fmt.Println("denom found in A:", denomInA.Id)
	} else {
		fmt.Println("denom not found in A:", denomInA.Id)
	}
	nftInA, errA := suite.chainA.App.NftKeeper.GetNFT(suite.chainA.GetContext(), "mobile", "xiaomi")

	if errA != nil {
		fmt.Println("nft not found in A")
	} else {
		fmt.Println("nft found in A:", nftInA.GetID())
	}

	// query B
	denomInB, found := suite.chainB.App.NftKeeper.GetDenom(suite.chainB.GetContext(), classInchainB)
	if found {
		fmt.Println("denom found in B:", denomInB.Id)
	} else {
		fmt.Println("denom not found in B:", denomInB.Id)
	}
	nftInB, errB := suite.chainB.App.NftKeeper.GetNFT(suite.chainB.GetContext(), classInchainB, "xiaomi")
	if errB != nil {
		fmt.Println("nft not found in B")
	} else {
		fmt.Println("nft found in B:", nftInB.GetID())
	}

	// query C
	denomInC, found := suite.chainC.App.NftKeeper.GetDenom(suite.chainC.GetContext(), classInchainC)
	if found {
		fmt.Println("denom found in C:", denomInC.Id)
	} else {
		fmt.Println("denom not found in C:", denomInC.Id)
	}
	nftInC, errC := suite.chainC.App.NftKeeper.GetNFT(suite.chainC.GetContext(), classInchainC, "xiaomi")
	if errC != nil {
		fmt.Println("nft not found in C")
	} else {
		fmt.Println("nft found in C:", nftInC.GetID())
	}

}

func TestTransferTestSuite(t *testing.T) {
	suite.Run(t, new(TransferTestSuite))
}
