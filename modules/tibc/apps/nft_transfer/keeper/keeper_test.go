package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/crypto"

	"github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/types"
	ibctesting "github.com/bianjieai/tibc-go/modules/tibc/testing"
)

type KeeperTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain
	chainC *ibctesting.TestChain
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 3)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.chainC = suite.coordinator.GetChain(ibctesting.GetChainID(2))
}

func (suite *KeeperTestSuite) TestGetTransferMoudleAddr() {
	expectedMaccAddr := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))

	macc := suite.chainA.App.NftTransferKeeper.GetNftTransferModuleAddr(types.ModuleName)

	suite.Require().NotNil(macc)
	suite.Require().Equal(expectedMaccAddr, macc)
}

func NewTransferPath(scChain, destChain *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(scChain, destChain)
	// setport
	return path
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
