package keeper_test

import (
	"fmt"
	ibctesting "github.com/bianjieai/tibc-go/modules/tibc/testing"
	nfttypes "github.com/irisnet/irismod/modules/nft/types"
)

//
func (suite *KeeperTestSuite) TestSendTransfer() {

	var (
		path   *ibctesting.Path
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
				issueDenomMsg :=  nfttypes.NewMsgIssueDenom("dog", "dog-name", "",
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
				issueDenomMsg :=  nfttypes.NewMsgIssueDenom("dog", "dog-name", "",
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
					suite.chainA.GetContext(), "dog", "taidy", "",
					suite.chainA.SenderAccount.GetAddress(), suite.chainB.SenderAccount.GetAddress().String(),
					true, suite.chainB.ChainID, "")

				// get nft from chainB
				dd, _:= suite.chainB.App.NftKeeper.GetDenom(suite.chainB.GetContext(), "dog")


				// send nft from chainB to chainA
				err := suite.chainB.App.NftTransferKeeper.SendNftTransfer(
					suite.chainB.GetContext(), dd.Id, "taidy", "",
					suite.chainB.SenderAccount.GetAddress(), suite.chainA.SenderAccount.GetAddress().String(),
					true, suite.chainA.ChainID, "")

				suite.Require().NoError(err) // message committed

			}
			err := suite.chainA.App.NftTransferKeeper.SendNftTransfer(
				suite.chainA.GetContext(), "mobile", "xiaomi", "",
				suite.chainA.SenderAccount.GetAddress(), suite.chainB.SenderAccount.GetAddress().String(),
				true, suite.chainB.ChainID, "")


			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}

		})
	}

}