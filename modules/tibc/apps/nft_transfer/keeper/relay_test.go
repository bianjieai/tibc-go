package keeper_test

import (
	"fmt"

	nfttypes "github.com/irisnet/irismod/modules/nft/types"

	"github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/types"
	packettypes "github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
	routingtypes "github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
	ibctesting "github.com/bianjieai/tibc-go/modules/tibc/testing"
)

const (
	CLASS  = "dog"
	PREFIX = "tibc/nft"
	NFTID  = "taidy"
)

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
			issueDenomMsg := nfttypes.NewMsgIssueDenom(
				"dog", "dog-name", "",
				suite.chainA.SenderAccount.GetAddress().String(),
				"", false, false, "", "", "", "",
			)
			_, _ = suite.chainA.SendMsgs(issueDenomMsg)

			// mint nft
			mintNftMsg := nfttypes.NewMsgMintNFT(
				"taidy", "dog", "", "", "", "",
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
			issueDenomMsg := nfttypes.NewMsgIssueDenom(
				"dog", "dog-name", "",
				suite.chainA.SenderAccount.GetAddress().String(),
				"", false, false, "", "", "", "",
			)
			_, _ = suite.chainA.SendMsgs(issueDenomMsg)

			// mint nft
			mintNftMsg := nfttypes.NewMsgMintNFT(
				"taidy", "dog", "taidy", "", "www.test.com", "none",
				suite.chainA.SenderAccount.GetAddress().String(),
				suite.chainA.SenderAccount.GetAddress().String(),
			)
			_, _ = suite.chainA.SendMsgs(mintNftMsg)

			data := types.NewNonFungibleTokenPacketData(
				"dog", "taidy", "www.test.com",
				suite.chainA.SenderAccount.GetAddress().String(),
				suite.chainB.SenderAccount.GetAddress().String(),
				true,
				"0xabcsda",
			)

			packet := packettypes.NewPacket(data.GetBytes(), uint64(1), suite.chainA.ChainID, suite.chainB.ChainID, "", string(routingtypes.NFT))
			_ = suite.chainB.App.NftTransferKeeper.OnRecvPacket(suite.chainB.GetContext(), packet, data)
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
				newClass = "tibc-5F88F7B2F39E49BB64D9682E6D7F8E10F8AA7DD10F6438FBAF1D4C659025F691"
				//newClass = PREFIX + "/" + suite.chainA.ChainID + "/" + CLASS
				// send nft from chainB to chainA
				err := suite.chainB.App.NftTransferKeeper.SendNftTransfer(
					suite.chainB.GetContext(),
					newClass,
					"taidy",
					suite.chainB.SenderAccount.GetAddress(),
					suite.chainA.SenderAccount.GetAddress().String(),
					suite.chainA.ChainID, "",
					"0xabcsda")

				suite.Require().NoError(err) // message committed

			}
			err := suite.chainA.App.NftTransferKeeper.SendNftTransfer(
				suite.chainA.GetContext(),
				"dog",
				"taidy",
				suite.chainA.SenderAccount.GetAddress(),
				suite.chainB.SenderAccount.GetAddress().String(),
				suite.chainB.ChainID, "",
				"0xabcsda")

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
			issueDenomMsg := nfttypes.NewMsgIssueDenom(
				"dog", "dog-name", "",
				suite.chainA.SenderAccount.GetAddress().String(),
				"", false, false, "", "", "", "",
			)
			_, _ = suite.chainA.SendMsgs(issueDenomMsg)

			// mint nft
			mintNftMsg := nfttypes.NewMsgMintNFT(
				"taidy", "dog", "", "", "", "",
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
				data := types.NewNonFungibleTokenPacketData(
					"dog", "taidy", "",
					suite.chainA.SenderAccount.GetAddress().String(),
					suite.chainB.SenderAccount.GetAddress().String(),
					true,
					"0xabcsda",
				)

				packet := packettypes.NewPacket(data.GetBytes(), seq, suite.chainA.ChainID, suite.chainB.ChainID, "", string(routingtypes.NFT))
				err := suite.chainB.App.NftTransferKeeper.OnRecvPacket(suite.chainB.GetContext(), packet, data)
				if tc.expPass {
					suite.Require().NoError(err)
				} else {
					suite.Require().Error(err)
				}
			} else {
				// send nft from B to A
				tc.malleate()

				data := types.NewNonFungibleTokenPacketData(
					newClass, "taidy", "",
					suite.chainB.SenderAccount.GetAddress().String(),
					suite.chainA.SenderAccount.GetAddress().String(),
					false,
					"0xabcsda",
				)
				packet := packettypes.NewPacket(data.GetBytes(), seq, suite.chainB.ChainID, suite.chainA.ChainID, "", string(routingtypes.NFT))
				err := suite.chainA.App.NftTransferKeeper.OnRecvPacket(suite.chainA.GetContext(), packet, data)
				if tc.expPass {
					suite.Require().NoError(err)
				} else {
					suite.Require().Error(err)
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestOnAcknowledgementPacket() {
	var (
		failedAck     = packettypes.NewErrorAcknowledgement("failed packet transfer")
		path          *ibctesting.Path
		newClass      string
		fullClassPath string
	)

	testCases := []struct {
		msg            string
		ack            packettypes.Acknowledgement
		malleate       func()
		success        bool // success of ack
		awayFromOrigin bool
		expPass        bool
	}{{
		"successful refund from source chain ", failedAck,
		func() {
			// issue denom
			issueDenomMsg := nfttypes.NewMsgIssueDenom(
				"dog", "dog-name", "",
				suite.chainA.SenderAccount.GetAddress().String(),
				"", false, false, "", "", "", "",
			)
			_, _ = suite.chainA.SendMsgs(issueDenomMsg)

			// mint nft
			mintNftMsg := nfttypes.NewMsgMintNFT(
				"taidy", "dog", "", "",
				"", "", suite.chainA.SenderAccount.GetAddress().String(),
				suite.chainA.SenderAccount.GetAddress().String(),
			)
			_, _ = suite.chainA.SendMsgs(mintNftMsg)

			// sendTransfer
			if err := suite.chainA.App.NftTransferKeeper.SendNftTransfer(
				suite.chainA.GetContext(), "dog", "taidy",
				suite.chainA.SenderAccount.GetAddress(),
				suite.chainB.SenderAccount.GetAddress().String(),
				suite.chainB.ChainID, "", "0xabcsda",
			); err != nil {
				fmt.Println("occur err :", err.Error())
			}
		},
		false, true, true,
	}, {
		"successful refund from sink chain ", failedAck,
		func() {
			data := types.NewNonFungibleTokenPacketData(
				"dog", "taidy", "www.test.com",
				suite.chainA.SenderAccount.GetAddress().String(),
				suite.chainB.SenderAccount.GetAddress().String(),
				true,
				"0xabcsda",
			)
			packet := packettypes.NewPacket(
				data.GetBytes(),
				uint64(1),
				suite.chainA.ChainID,
				suite.chainB.ChainID,
				"",
				string(routingtypes.NFT))
			_ = suite.chainB.App.NftTransferKeeper.OnRecvPacket(suite.chainB.GetContext(), packet, data)

		}, true, false, true},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset
			path = NewTransferPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)

			tc.malleate()
			if tc.awayFromOrigin {
				data := types.NewNonFungibleTokenPacketData(
					"dog",
					"taidy",
					"",
					suite.chainA.SenderAccount.GetAddress().String(),
					suite.chainB.SenderAccount.GetAddress().String(),
					true,
					"0xabcsda",
				)

				err := suite.chainA.App.NftTransferKeeper.OnAcknowledgementPacket(suite.chainA.GetContext(), data, tc.ack)
				if tc.expPass {
					suite.Require().NoError(err)
					if tc.success {
						nft, err := suite.chainA.App.NftKeeper.GetNFT(suite.chainA.GetContext(), "dog", "taidy")
						if err == nil {
							// The nft owner before sending and the nft owner after ACK must be the same
							suite.Require().Equal(suite.chainA.SenderAccount.GetAddress().String(), nft.GetOwner().String())
						} else {
							fmt.Println("not found nft")
						}
					}
				}
			} else {
				fullClassPath = "nft" + "/" + suite.chainA.ChainID + "/" + suite.chainB.ChainID + "/" + CLASS
				newClass = "tibc-5F88F7B2F39E49BB64D9682E6D7F8E10F8AA7DD10F6438FBAF1D4C659025F691"
				// send nft from chainB to chainA
				_ = suite.chainB.App.NftTransferKeeper.SendNftTransfer(
					suite.chainB.GetContext(), newClass,
					"taidy",
					suite.chainB.SenderAccount.GetAddress(),
					suite.chainA.SenderAccount.GetAddress().String(),
					suite.chainA.ChainID,
					"",
					"0xabcsda")

				data := types.NewNonFungibleTokenPacketData(
					fullClassPath,
					"taidy",
					"",
					suite.chainB.SenderAccount.GetAddress().String(),
					suite.chainA.SenderAccount.GetAddress().String(),
					false,
					"0xabcsda")

				err := suite.chainB.App.NftTransferKeeper.OnAcknowledgementPacket(suite.chainB.GetContext(), data, tc.ack)

				if tc.expPass {
					suite.Require().NoError(err)
					if tc.success {
						nft, err := suite.chainA.App.NftKeeper.GetNFT(suite.chainA.GetContext(), newClass, "taidy")
						if err == nil {
							fmt.Println("found nft", nft.GetOwner())
						} else {
							fmt.Println("not found nft")
						}
					}
				}
			}
		})
	}
}
