package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"

	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/types"
	tibctesting "github.com/bianjieai/tibc-go/modules/tibc/testing"
)

type KeeperTestSuite struct {
	suite.Suite
	ctx sdk.Context

	coordinator *tibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA      *tibctesting.TestChain
	chainB      *tibctesting.TestChain
	chainC      *tibctesting.TestChain
	queryClient types.QueryClient
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.coordinator = tibctesting.NewCoordinator(suite.T(), 3)
	suite.chainA = suite.coordinator.GetChain(tibctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(tibctesting.GetChainID(1))
	suite.chainC = suite.coordinator.GetChain(tibctesting.GetChainID(2))

	queryHelper := baseapp.NewQueryServerTestHelper(suite.chainA.GetContext(), suite.chainA.App.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.chainA.App.NftTransferKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)
}

func (suite *KeeperTestSuite) TestGetTransferMoudleAddr() {
	expectedMaccAddr := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))

	macc := suite.chainA.App.NftTransferKeeper.GetNftTransferModuleAddr(types.ModuleName)

	suite.Require().NotNil(macc)
	suite.Require().Equal(expectedMaccAddr, macc)
}

func NewTransferPath(scChain, destChain *tibctesting.TestChain) *tibctesting.Path {
	path := tibctesting.NewPath(scChain, destChain)
	return path
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
