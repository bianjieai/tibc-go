package keeper_test

import (
	"fmt"
	"testing"

	mttypes "github.com/irisnet/irismod/modules/mt/types"
	"github.com/stretchr/testify/suite"

	"github.com/bianjieai/tibc-go/modules/tibc/apps/mt_transfer/types"
	packettypes "github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
	routingtypes "github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
	ibctesting "github.com/bianjieai/tibc-go/modules/tibc/testing"
)

const (
	CLASS  = "c02a799c8fee067a7f9b944554d8431ee539847234441833e45a3a2d3123fd99" //sha256("mt-denom-1")
	PREFIX = "tibc/mt"
	MTID   = "ff6e57b41cb52ae7d58d854b2123da2c5657fd15d525821a13fe7da1b9cebd80" //sha256("mt-1")
)

func TestRelayTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestSendTransfer() {

	var (
		path     *ibctesting.Path
		newClass string
	)

	testCases := []struct {
		msg            string
		malleate       func()
		awayFromSource bool
		expPass        bool
	}{{
		"successful transfer from source chain",
		func() {
			// issue denom
			issueDenomMsg := mttypes.NewMsgIssueDenom(
				"dog-name", "",
				suite.chainA.SenderAccount.GetAddress().String(),
			)
			_, _ = suite.chainA.SendMsgs(issueDenomMsg)

			// mint nft
			mintNftMsg := mttypes.NewMsgMintMT(
				"", CLASS, 2, "",
				suite.chainA.SenderAccount.GetAddress().String(),
				suite.chainA.SenderAccount.GetAddress().String(),
			)
			_, _ = suite.chainA.SendMsgs(mintNftMsg)
		},
		true, true,
	}, {
		"successful transfer from sink chain",
		func() {
			// issue denom
			issueDenomMsg := mttypes.NewMsgIssueDenom(
				"dog-name", "",
				suite.chainA.SenderAccount.GetAddress().String(),
			)
			_, _ = suite.chainA.SendMsgs(issueDenomMsg)

			// mint nft
			mintNftMsg := mttypes.NewMsgMintMT(
				"", CLASS, 2, "none",
				suite.chainA.SenderAccount.GetAddress().String(),
				suite.chainA.SenderAccount.GetAddress().String(),
			)
			_, _ = suite.chainA.SendMsgs(mintNftMsg)

			data := types.NewMultiTokenPacketData(
				CLASS, MTID,
				suite.chainA.SenderAccount.GetAddress().String(),
				suite.chainB.SenderAccount.GetAddress().String(),
				true,
				"0xabcsda",
				1,
				[]byte(""),
			)

			packet := packettypes.NewPacket(data.GetBytes(), uint64(1), suite.chainA.ChainID, suite.chainB.ChainID, "", string(routingtypes.MT))
			_ = suite.chainB.App.MtTransferKeeper.OnRecvPacket(suite.chainB.GetContext(), packet, data)
		},
		false, true,
	}}

	for _, tc := range testCases {
		tc := tc
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest()
			path = NewTransferPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)

			tc.malleate()
			if !tc.awayFromSource {
				newClass = "tibc-1AA90CD7273981C2E2CFFC5415D887C17FB03A7FCB49EFE4A7B6878E384F1670"
				//newClass = PREFIX + "/" + suite.chainA.ChainID + "/" + CLASS
				// send nft from chainB to chainA
				err := suite.chainB.App.MtTransferKeeper.SendMtTransfer(
					suite.chainB.GetContext(),
					newClass,
					MTID,
					suite.chainB.SenderAccount.GetAddress(),
					suite.chainA.SenderAccount.GetAddress().String(),
					suite.chainA.ChainID, "",
					"0xabcsda", 1)

				suite.Require().NoError(err) // message committed

			}
			err := suite.chainA.App.MtTransferKeeper.SendMtTransfer(
				suite.chainA.GetContext(),
				CLASS,
				MTID,
				suite.chainA.SenderAccount.GetAddress(),
				suite.chainB.SenderAccount.GetAddress().String(),
				suite.chainB.ChainID, "",
				"0xabcsda", 1)
			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}

		})
	}
}

func (suite *KeeperTestSuite) TestOnRecvPacket() {
	var newClass string
	testCases := []struct {
		msg            string
		malleate       func()
		awayfromSource bool
		expPass        bool
	}{{
		"success receive on awayfromSource chain", func() {}, true, true,
	}, {
		"failed receive on nearToSource chain",
		func() {
			// issue denom
			issueDenomMsg := mttypes.NewMsgIssueDenom(
				"dog-name", "",
				suite.chainA.SenderAccount.GetAddress().String(),
			)
			_, _ = suite.chainA.SendMsgs(issueDenomMsg)

			// mint nft
			mintNftMsg := mttypes.NewMsgMintMT(
				"", CLASS, 2, "",
				suite.chainA.SenderAccount.GetAddress().String(),
				suite.chainA.SenderAccount.GetAddress().String(),
			)
			_, _ = suite.chainA.SendMsgs(mintNftMsg)
		},
		false, false,
	}, {
		"failed receive on nearToSource chain without prefix ",
		func() { newClass = CLASS },
		false, false,
	}}

	for _, tc := range testCases {
		tc := tc
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset
			path := NewTransferPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)

			seq := uint64(1)
			if tc.awayfromSource {
				// // send nft from A to B
				data := types.NewMultiTokenPacketData(
					CLASS, MTID,
					suite.chainA.SenderAccount.GetAddress().String(),
					suite.chainB.SenderAccount.GetAddress().String(),
					true,
					"0xabcsda", 1, []byte(""),
				)

				packet := packettypes.NewPacket(data.GetBytes(), seq, suite.chainA.ChainID, suite.chainB.ChainID, "", string(routingtypes.NFT))
				err := suite.chainB.App.MtTransferKeeper.OnRecvPacket(suite.chainB.GetContext(), packet, data)
				if tc.expPass {
					suite.Require().NoError(err)
				} else {
					suite.Require().Error(err)
				}
			} else {
				// send nft from B to A
				tc.malleate()

				data := types.NewMultiTokenPacketData(
					newClass, MTID,
					suite.chainB.SenderAccount.GetAddress().String(),
					suite.chainA.SenderAccount.GetAddress().String(),
					false,
					"0xabcsda", 1, []byte(""),
				)
				packet := packettypes.NewPacket(data.GetBytes(), seq, suite.chainB.ChainID, suite.chainA.ChainID, "", string(routingtypes.NFT))
				err := suite.chainA.App.MtTransferKeeper.OnRecvPacket(suite.chainA.GetContext(), packet, data)
				if tc.expPass {
					suite.Require().NoError(err)
				} else {
					suite.Require().Error(err)
				}
			}
		})
	}
}

// func (suite *KeeperTestSuite) TestOnAcknowledgementPacket() {
// 	var (
// 		failedAck     = packettypes.NewErrorAcknowledgement("failed packet transfer")
// 		path          *ibctesting.Path
// 		newClass      string
// 		fullClassPath string
// 	)

// 	testCases := []struct {
// 		msg            string
// 		ack            packettypes.Acknowledgement
// 		malleate       func()
// 		success        bool // success of ack
// 		awayFromOrigin bool
// 		expPass        bool
// 	}{{
// 		"successful refund from source chain ", failedAck,
// 		func() {
// 			// issue denom
// 			issueDenomMsg := nfttypes.NewMsgIssueDenom(
// 				"dog", "dog-name", "",
// 				suite.chainA.SenderAccount.GetAddress().String(),
// 				"", false, false,
// 			)
// 			_, _ = suite.chainA.SendMsgs(issueDenomMsg)

// 			// mint nft
// 			mintNftMsg := nfttypes.NewMsgMintNFT(
// 				"taidy", "dog", "",
// 				"", "", suite.chainA.SenderAccount.GetAddress().String(),
// 				suite.chainA.SenderAccount.GetAddress().String(),
// 			)
// 			_, _ = suite.chainA.SendMsgs(mintNftMsg)

// 			// sendTransfer
// 			if err := suite.chainA.App.NftTransferKeeper.SendNftTransfer(
// 				suite.chainA.GetContext(), "dog", "taidy",
// 				suite.chainA.SenderAccount.GetAddress(),
// 				suite.chainB.SenderAccount.GetAddress().String(),
// 				suite.chainB.ChainID, "", "0xabcsda",
// 			); err != nil {
// 				fmt.Println("occur err :", err.Error())
// 			}
// 		},
// 		false, true, true,
// 	}, {
// 		"successful refund from sink chain ", failedAck,
// 		func() {
// 			data := types.NewNonFungibleTokenPacketData(
// 				"dog", "taidy", "www.test.com",
// 				suite.chainA.SenderAccount.GetAddress().String(),
// 				suite.chainB.SenderAccount.GetAddress().String(),
// 				true,
// 				"0xabcsda",
// 			)
// 			packet := packettypes.NewPacket(
// 				data.GetBytes(),
// 				uint64(1),
// 				suite.chainA.ChainID,
// 				suite.chainB.ChainID,
// 				"",
// 				string(routingtypes.NFT))
// 			_ = suite.chainB.App.NftTransferKeeper.OnRecvPacket(suite.chainB.GetContext(), packet, data)

// 		}, true, false, true},
// 	}

// 	for _, tc := range testCases {
// 		tc := tc
// 		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
// 			suite.SetupTest() // reset
// 			path = NewTransferPath(suite.chainA, suite.chainB)
// 			suite.coordinator.SetupClients(path)

// 			tc.malleate()
// 			if tc.awayFromOrigin {
// 				data := types.NewNonFungibleTokenPacketData(
// 					"dog",
// 					"taidy",
// 					"",
// 					suite.chainA.SenderAccount.GetAddress().String(),
// 					suite.chainB.SenderAccount.GetAddress().String(),
// 					true,
// 					"0xabcsda",
// 				)

// 				err := suite.chainA.App.NftTransferKeeper.OnAcknowledgementPacket(suite.chainA.GetContext(), data, tc.ack)
// 				if tc.expPass {
// 					suite.Require().NoError(err)
// 					if tc.success {
// 						nft, err := suite.chainA.App.NftKeeper.GetNFT(suite.chainA.GetContext(), "dog", "taidy")
// 						if err == nil {
// 							// The nft owner before sending and the nft owner after ACK must be the same
// 							suite.Require().Equal(suite.chainA.SenderAccount.GetAddress().String(), nft.GetOwner().String())
// 						} else {
// 							fmt.Println("not found nft")
// 						}
// 					}
// 				}
// 			} else {
// 				fullClassPath = "nft" + "/" + suite.chainA.ChainID + "/" + suite.chainB.ChainID + "/" + CLASS
// 				newClass = "tibc-5F88F7B2F39E49BB64D9682E6D7F8E10F8AA7DD10F6438FBAF1D4C659025F691"
// 				// send nft from chainB to chainA
// 				_ = suite.chainB.App.NftTransferKeeper.SendNftTransfer(
// 					suite.chainB.GetContext(), newClass,
// 					"taidy",
// 					suite.chainB.SenderAccount.GetAddress(),
// 					suite.chainA.SenderAccount.GetAddress().String(),
// 					suite.chainA.ChainID,
// 					"",
// 					"0xabcsda")

// 				data := types.NewNonFungibleTokenPacketData(
// 					fullClassPath,
// 					"taidy",
// 					"",
// 					suite.chainB.SenderAccount.GetAddress().String(),
// 					suite.chainA.SenderAccount.GetAddress().String(),
// 					false,
// 					"0xabcsda")

// 				err := suite.chainB.App.NftTransferKeeper.OnAcknowledgementPacket(suite.chainB.GetContext(), data, tc.ack)

// 				if tc.expPass {
// 					suite.Require().NoError(err)
// 					if tc.success {
// 						nft, err := suite.chainA.App.NftKeeper.GetNFT(suite.chainA.GetContext(), newClass, "taidy")
// 						if err == nil {
// 							fmt.Println("found nft", nft.GetOwner())
// 						} else {
// 							fmt.Println("not found nft")
// 						}
// 					}
// 				}
// 			}
// 		})
// 	}
// }
