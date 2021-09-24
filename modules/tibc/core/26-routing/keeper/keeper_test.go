package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
	tibctesting "github.com/bianjieai/tibc-go/modules/tibc/testing"
)

var (
	validPort   = "validportid"
	invalidPort = "(invalidPortID)"
)

type KeeperTestSuite struct {
	suite.Suite
	coordinator *tibctesting.Coordinator
	chain       *tibctesting.TestChain
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.coordinator = tibctesting.NewCoordinator(suite.T(), 1)
	suite.chain = suite.coordinator.GetChain(tibctesting.GetChainID(0))
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestSetRouter() {
	// There is a sealed router in the RoutingKeeper
	suite.Require().True(suite.chain.App.TIBCKeeper.RoutingKeeper.Router.Sealed())
	suite.chain.App.TIBCKeeper.RoutingKeeper.Router = nil
	router := types.NewRouter()
	router.AddRoute("1", nil)
	router.Sealed()
	suite.chain.App.TIBCKeeper.RoutingKeeper.SetRouter(router)
	suite.Require().Equal(router, suite.chain.App.TIBCKeeper.RoutingKeeper.Router)
}

func (suite KeeperTestSuite) TestGetRoutingRules() {
	rules := []string{"xxx.xxx.xxx", "yyy.yyy.yyy"}
	ctx := suite.chain.GetContext()
	err := suite.chain.App.TIBCKeeper.RoutingKeeper.SetRoutingRules(ctx, rules)
	suite.Require().NoError(err)
	rules2, ok := suite.chain.App.TIBCKeeper.RoutingKeeper.GetRoutingRules(ctx)
	suite.Require().True(ok)
	suite.Require().Equal(rules, rules2)
}

func (suite KeeperTestSuite) TestAuthenticate() {
	ctx := suite.chain.GetContext()
	rules := []string{"a*b*.c*d*.tt", "b.c.d"}
	err := suite.chain.App.TIBCKeeper.RoutingKeeper.SetRoutingRules(ctx, rules)
	suite.Require().NoError(err)
	ok := suite.chain.App.TIBCKeeper.RoutingKeeper.Authenticate(ctx, "aabb", "ccdd", "tt")
	suite.Require().True(ok)
	ok = suite.chain.App.TIBCKeeper.RoutingKeeper.Authenticate(ctx, "aadb", "ccbe", "tt")
	suite.Require().False(ok)
}
