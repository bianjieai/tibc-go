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
		path *ibctesting.Path
		//err    error
	)

	testCases := []struct {
		msg            string
		malleate       func()
		awayFromSource bool
		expPass        bool
	}{
		{"successful transfer from source chain",
			func() {
				suite.coordinator.SetupClients(path)
				// issue denom
				issueDenomMsg := nfttypes.NewMsgIssueDenom("dog", "dog-name", "",
					suite.chainA.SenderAccount.GetAddress().String(), "", false, false)
				_, _ = suite.chainA.SendMsgs(issueDenomMsg)

				// mint nft
				mintNftMsg := nfttypes.NewMsgMintNFT("taidy", "dog", "",
					"", "", suite.chainA.SenderAccount.GetAddress().String(),
					suite.chainA.SenderAccount.GetAddress().String())
				_, _ = suite.chainA.SendMsgs(mintNftMsg)

			}, true, true},

		{"successful transfer from sink chain",
			func() {
				suite.coordinator.SetupClients(path)
				// issue denom
				issueDenomMsg := nfttypes.NewMsgIssueDenom("dog", "dog-name", "",
					suite.chainA.SenderAccount.GetAddress().String(), "", false, false)
				_, _ = suite.chainA.SendMsgs(issueDenomMsg)

				// mint nft
				mintNftMsg := nfttypes.NewMsgMintNFT("taidy", "dog", "",
					"", "", suite.chainA.SenderAccount.GetAddress().String(),
					suite.chainA.SenderAccount.GetAddress().String())
				_, _ = suite.chainA.SendMsgs(mintNftMsg)

			}, false, true},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset
			path = NewTransferPath(suite.chainA, suite.chainB)

			tc.malleate()
			if !tc.awayFromSource {

				// send nft from chainA to ChainB
				_ = suite.chainA.App.NftTransferKeeper.SendNftTransfer(
					suite.chainA.GetContext(), "dog", "taidy",
					suite.chainA.SenderAccount.GetAddress(), suite.chainB.SenderAccount.GetAddress().String(),
					suite.chainB.ChainID, "")

				// get nft from chainB
				dd, _ := suite.chainB.App.NftKeeper.GetDenom(suite.chainB.GetContext(), "dog")

				// send nft from chainB to chainA
				err := suite.chainB.App.NftTransferKeeper.SendNftTransfer(
					suite.chainB.GetContext(), dd.Id, "taidy",
					suite.chainB.SenderAccount.GetAddress(), suite.chainA.SenderAccount.GetAddress().String(),
					suite.chainA.ChainID, "")

				suite.Require().NoError(err) // message committed

			}
			err := suite.chainA.App.NftTransferKeeper.SendNftTransfer(
				suite.chainA.GetContext(), "dog", "taidy",
				suite.chainA.SenderAccount.GetAddress(), suite.chainB.SenderAccount.GetAddress().String(),
				suite.chainB.ChainID, "")

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}

		})
	}

}

func (suite *KeeperTestSuite) TestOnRecvPacket() {

	testCases := []struct {
		msg            string
		malleate       func()
		awayfromSource bool
		expPass        bool
	}{
		{"success receive on awayfromSource chain", func() {}, true, true},
		{"success receive on nearToSource chain", func() {}, false, true},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset
			path := NewTransferPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)

			seq := uint64(1)
			if tc.awayfromSource {
				data := types.NewNonFungibleTokenPacketData(CLASS, NFTID, "", suite.chainA.SenderAccount.GetAddress().String(),
					suite.chainB.SenderAccount.GetAddress().String(), true)

				packet := packettypes.NewPacket(data.GetBytes(), seq, suite.chainA.ChainID, suite.chainB.ChainID, "", string(routingtypes.NFT))
				err := suite.chainB.App.NftTransferKeeper.OnRecvPacket(suite.chainB.GetContext(), packet, data)
				if tc.expPass {
					suite.Require().NoError(err)
				} else {
					suite.Require().Error(err)
				}
			} else {
				newClass := PREFIX + "/" + suite.chainB.ChainID + "/" + CLASS
				data := types.NewNonFungibleTokenPacketData(newClass, NFTID, "", suite.chainB.SenderAccount.GetAddress().String(),
					suite.chainA.SenderAccount.GetAddress().String(), false)

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
