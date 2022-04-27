package mt_transfer_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	mttypes "github.com/irisnet/irismod/modules/mt/types"

	"github.com/bianjieai/tibc-go/modules/tibc/apps/mt_transfer/types"
	packettypes "github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
	routingtypes "github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
	tibctesting "github.com/bianjieai/tibc-go/modules/tibc/testing"
)

const (
	ClassID = "c02a799c8fee067a7f9b944554d8431ee539847234441833e45a3a2d3123fd99" //sha256("mt-denom-1")
	MtID    = "ff6e57b41cb52ae7d58d854b2123da2c5657fd15d525821a13fe7da1b9cebd80" //sha256("mt-1")
	Amount  = 2
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
mt
A->B B->C
*/
func (suite *TransferTestSuite) TestHandleMsgTransfer() {
	// setup between chainA and chainB

	path := tibctesting.NewPath(suite.chainA, suite.chainB)

	suite.coordinator.SetupClients(path)

	// issue denom
	issueDenomMsg := mttypes.NewMsgIssueDenom(
		"mobile-name", "",
		suite.chainA.SenderAccount.GetAddress().String(),
	)
	_, _ = suite.chainA.SendMsgs(issueDenomMsg)

	// mint mt
	mintMtMsg := mttypes.NewMsgMintMT(
		"", ClassID, Amount, "",
		suite.chainA.SenderAccount.GetAddress().String(),
		suite.chainA.SenderAccount.GetAddress().String(),
	)
	_, _ = suite.chainA.SendMsgs(mintMtMsg)

	dd, has := suite.chainA.App.MtKeeper.GetDenom(suite.chainA.GetContext(), ClassID)
	suite.Require().Truef(has, "denom %s not found", ClassID)

	// send mt from A To B
	msg := types.NewMsgMtTransfer(
		dd.Id, MtID,
		suite.chainA.SenderAccount.GetAddress().String(),
		suite.chainB.SenderAccount.GetAddress().String(),
		suite.chainB.ChainID, "", "0xabcsda", 1,
	)

	_, err := suite.chainA.SendMsgs(msg)
	suite.Require().NoError(err) // message committed
	//// relay send
	multiTokenPacketData := types.NewMultiTokenPacketData(
		dd.Id, MtID,
		suite.chainA.SenderAccount.GetAddress().String(),
		suite.chainB.SenderAccount.GetAddress().String(),
		true,
		"0xabcsda",
		1,
		[]byte(""),
	)
	packet := packettypes.NewPacket(
		multiTokenPacketData.GetBytes(), 1,
		path.EndpointA.ChainName, path.EndpointB.ChainName,
		"", string(routingtypes.MT),
	)

	ack := packettypes.NewResultAcknowledgement([]byte{byte(1)})
	err = path.RelayPacket(packet, ack.GetBytes())
	suite.Require().NoError(err) // relay committed

	// check that voucher exists on chain B
	// denomID :tibc-hash(mt/chainA.chainID/chainB.chainId/mobile)
	classInChainB := "tibc-1AA90CD7273981C2E2CFFC5415D887C17FB03A7FCB49EFE4A7B6878E384F1670"
	mt, err := suite.chainB.App.MtKeeper.GetMT(suite.chainB.GetContext(), classInChainB, MtID)
	suite.Require().NoError(err) // message committed
	suite.Require().Equal(MtID, mt.GetID())

	balanceInChainA := suite.chainA.App.MtKeeper.GetBalance(suite.chainA.GetContext(), ClassID, MtID, suite.chainA.SenderAccount.GetAddress())
	suite.Require().Equal(uint64(1), balanceInChainA)

	balanceInChainB := suite.chainB.App.MtKeeper.GetBalance(suite.chainB.GetContext(), classInChainB, MtID, suite.chainB.SenderAccount.GetAddress())
	suite.Require().Equal(uint64(1), balanceInChainB)

	// setup between chainB to chainC
	pathBtoC := tibctesting.NewPath(suite.chainB, suite.chainC)
	suite.coordinator.SetupClients(pathBtoC)

	// send ft from chainB to chainC
	msgfromBToC := types.NewMsgMtTransfer(
		classInChainB, MtID,
		suite.chainB.SenderAccount.GetAddress().String(),
		suite.chainC.SenderAccount.GetAddress().String(),
		suite.chainC.ChainID, "",
		"0xabcsda",
		1,
	)

	_, err1 := suite.chainB.SendMsgs(msgfromBToC)
	suite.Require().NoError(err1) // message committed

	fullClassPathFromBToC := "mt" + "/" + suite.chainA.ChainID + "/" + suite.chainB.ChainID + "/" + ClassID
	// relay send
	mtPacketFromBToC := types.NewMultiTokenPacketData(
		fullClassPathFromBToC, MtID,
		suite.chainB.SenderAccount.GetAddress().String(),
		suite.chainC.SenderAccount.GetAddress().String(),
		true,
		"0xabcsda",
		1,
		[]byte(""),
	)
	packetFromBToC := packettypes.NewPacket(
		mtPacketFromBToC.GetBytes(), 1,
		pathBtoC.EndpointA.ChainName,
		pathBtoC.EndpointB.ChainName,
		"", string(routingtypes.MT),
	)

	ack1 := packettypes.NewResultAcknowledgement([]byte{byte(1)})
	err = pathBtoC.RelayPacket(packetFromBToC, ack1.GetBytes())
	suite.Require().NoError(err) // relay committed

	// check that voucher exists on chain C
	// denomID : tibc-{hash(mt/chainA.chainID/chainB.chainID/chainC.chainID/ClassID)}
	classInchainC := "tibc-14B78BCCA85B511A58AF995DFFE2621CC6BD00FAF43FC33B2270508709DE2C94"
	mtInC, err := suite.chainC.App.MtKeeper.GetMT(suite.chainC.GetContext(), classInchainC, MtID)
	suite.Require().NoError(err)
	suite.Require().Equal(MtID, mtInC.GetID())

	balanceInChainB = suite.chainB.App.MtKeeper.GetBalance(suite.chainB.GetContext(), classInChainB, MtID, suite.chainB.SenderAccount.GetAddress())
	suite.Require().Equal(uint64(0), balanceInChainB)

	balanceInChainC := suite.chainC.App.MtKeeper.GetBalance(suite.chainC.GetContext(), classInchainC, MtID, suite.chainC.SenderAccount.GetAddress())
	suite.Require().Equal(uint64(1), balanceInChainC)

	// send mt from chainC back to chainB
	msgfromCToB := types.NewMsgMtTransfer(
		classInchainC, MtID,
		suite.chainC.SenderAccount.GetAddress().String(),
		suite.chainB.SenderAccount.GetAddress().String(),
		suite.chainB.ChainID, "",
		"0xabcsda",
		1,
	)

	_, err2 := suite.chainC.SendMsgs(msgfromCToB)
	suite.Require().NoError(err2) // message committed

	fullClassPathFromCToB := "mt" + "/" + suite.chainA.ChainID + "/" + suite.chainB.ChainID + "/" + suite.chainC.ChainID + "/" + ClassID
	// relay send
	mtPacket := types.NewMultiTokenPacketData(
		fullClassPathFromCToB, MtID,
		suite.chainC.SenderAccount.GetAddress().String(),
		suite.chainB.SenderAccount.GetAddress().String(),
		false,
		"0xabcsda",
		1,
		[]byte(""),
	)
	packetFromCToB := packettypes.NewPacket(
		mtPacket.GetBytes(), 1,
		pathBtoC.EndpointB.ChainName,
		pathBtoC.EndpointA.ChainName,
		"", string(routingtypes.MT),
	)

	ack2 := packettypes.NewResultAcknowledgement([]byte{byte(1)})
	err = pathBtoC.RelayPacket(packetFromCToB, ack2.GetBytes())
	suite.Require().NoError(err) // relay committed
	balanceInChainB = suite.chainB.App.MtKeeper.GetBalance(suite.chainB.GetContext(), classInChainB, MtID, suite.chainB.SenderAccount.GetAddress())
	suite.Require().Equal(uint64(1), balanceInChainB)

	balanceInChainC = suite.chainC.App.MtKeeper.GetBalance(suite.chainC.GetContext(), classInchainC, MtID, suite.chainC.SenderAccount.GetAddress())
	suite.Require().Equal(uint64(0), balanceInChainC)

	// send nft  from chainB back to chainA
	msgFromBToA := types.NewMsgMtTransfer(
		classInChainB, MtID,
		suite.chainB.SenderAccount.GetAddress().String(),
		suite.chainA.SenderAccount.GetAddress().String(),
		suite.chainA.ChainID, "", "0xabcsda", 1,
	)

	_, err = suite.chainB.SendMsgs(msgFromBToA)
	suite.Require().NoError(err) // message committed

	fullClassPathFromBToA := "mt" + "/" + suite.chainA.ChainID + "/" + suite.chainB.ChainID + "/" + ClassID
	// relay send
	multiTokenPacketData = types.NewMultiTokenPacketData(
		fullClassPathFromBToA, MtID,
		suite.chainB.SenderAccount.GetAddress().String(),
		suite.chainA.SenderAccount.GetAddress().String(),
		false,
		"0xabcsda",
		1,
		[]byte(""),
	)
	packet = packettypes.NewPacket(
		multiTokenPacketData.GetBytes(), 1,
		path.EndpointB.ChainName, path.EndpointA.ChainName,
		"", string(routingtypes.MT),
	)

	ack = packettypes.NewResultAcknowledgement([]byte{byte(1)})
	err = path.RelayPacket(packet, ack.GetBytes())
	suite.Require().NoError(err) // relay committed

	// Query whether there are corresponding mt for C, B, and A respectively
	/* just do  A->B->C
	denom found in A: ClassID
	mt found in A: MtID

	denom found in B: tibc/mt/testchain0/ClassID
	nft found in B: MtID

	denom found in C: tibc/mt/testchain0/testchain1/ClassID
	nft found in C: MtID
	*/

	/* do A->B->C  then do C->B->A
	denom found in A: ClassID
	mt found in A: MtID

	denom found in B: tibc/mt/testchain0/ClassID
	mt not found in B

	denom found in C: tibc/mt/testchain0/testchain1/ClassID
	mt not found in C

	*/
	// query A
	denomInA, found := suite.chainA.App.MtKeeper.GetDenom(suite.chainA.GetContext(), ClassID)
	if found {
		fmt.Println("denom found in A:", denomInA.Id)
	} else {
		fmt.Println("denom not found in A:", denomInA.Id)
	}
	mtInA, errA := suite.chainA.App.MtKeeper.GetMT(suite.chainA.GetContext(), ClassID, MtID)

	if errA != nil {
		fmt.Println("mt not found in A")
	} else {
		fmt.Println("mt found in A:", mtInA.GetID())
	}

	// query B
	denomInB, found := suite.chainB.App.MtKeeper.GetDenom(suite.chainB.GetContext(), classInChainB)
	if found {
		fmt.Println("denom found in B:", denomInB.Id)
	} else {
		fmt.Println("denom not found in B:", denomInB.Id)
	}
	mtInB, errB := suite.chainB.App.MtKeeper.GetMT(suite.chainB.GetContext(), classInChainB, MtID)
	if errB != nil {
		fmt.Println("mt not found in B")
	} else {
		fmt.Println("mt found in B:", mtInB.GetID())
	}

	// query C
	denomInC, found := suite.chainC.App.MtKeeper.GetDenom(suite.chainC.GetContext(), classInchainC)
	if found {
		fmt.Println("denom found in C:", denomInC.Id)
	} else {
		fmt.Println("denom not found in C:", denomInC.Id)
	}
	mtInC, errC := suite.chainC.App.MtKeeper.GetMT(suite.chainC.GetContext(), classInchainC, MtID)
	if errC != nil {
		fmt.Println("mt not found in C")
	} else {
		fmt.Println("mt found in C:", mtInC.GetID())
	}
}

func TestTransferTestSuite(t *testing.T) {
	suite.Run(t, new(TransferTestSuite))
}
