package keeper_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
	tibctesting "github.com/bianjieai/tibc-go/modules/tibc/testing"
)

//used by TestSetRoutingRules
type testCase1 struct {
	msg     string
	rules   []string
	expPass bool
}

// used by TestAuthenticate
type testCase2 struct {
	msg                string
	rules              []string
	source, dest, port string
	expPass            bool
}

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
	suite.coordinator.CommitNBlocks(suite.chain, 2)
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

func (suite KeeperTestSuite) TestSetRoutingRules() {
	testCases := []testCase1{
		{
			"1 success with no *",
			[]string{"source,dest,port"},
			true,
		},
		{
			"2 success with *",
			[]string{"*,*,*"},
			true,
		},
		{
			"3 success no rule",
			[]string{},
			true,
		},
		{
			"4 fail due to invalid character",
			[]string{"xxx,*dd,??dd"},
			false,
		},
		{
			"5 fail due to number of commas -1",
			[]string{"a.b,c"},
			false,
		},
		{
			"6 fail due to number of commas -2",
			[]string{"a,b,c,d"},
			false,
		},
		{
			"7 fail due to no content",
			[]string{",dd,"},
			false,
		},
		{
			"8 fail no content",
			[]string{""},
			false,
		},
	}

	for i, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s, %d/%d tests", tc.msg, i+1, len(testCases)), func() {
			suite.SetupTest() // reset the context
			rules := tc.rules
			err := suite.chain.App.TIBCKeeper.RoutingKeeper.SetRoutingRules(suite.chain.GetContext(), rules)
			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
func (suite KeeperTestSuite) TestGetRoutingRules() {
	rules := []string{"xxx,xxx,xxx", "yyy,yyy,yyy"}
	ctx := suite.chain.GetContext()
	err := suite.chain.App.TIBCKeeper.RoutingKeeper.SetRoutingRules(ctx, rules)
	suite.Require().NoError(err)
	rules2, ok := suite.chain.App.TIBCKeeper.RoutingKeeper.GetRoutingRules(ctx)
	suite.Require().True(ok)
	suite.Require().Equal(rules, rules2)
}

func (suite KeeperTestSuite) TestAuthenticate() {
	testCases := []testCase2{
		{
			"1 success with no *",
			[]string{"source,dest,port"},
			"source",
			"dest",
			"port",
			true,
		},
		{
			"2 success with *",
			[]string{"*,*,*"},
			"aabb",
			"cc",
			"dd",
			true,
		},
		{
			"3 fail,null rules",
			[]string{},
			"aabb",
			"cc",
			"dd",
			false,
		},
		{
			"4 fail,unauthenticated",
			[]string{"source,dest,port"},
			"source",
			"dest",
			"dgsbl",
			false,
		},
	}

	for i, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s, %d/%d tests", tc.msg, i+1, len(testCases)), func() {
			suite.SetupTest() // reset the context
			rules := tc.rules
			err := suite.chain.App.TIBCKeeper.RoutingKeeper.SetRoutingRules(suite.chain.GetContext(), rules)
			suite.Require().NoError(err)
			ok := suite.chain.App.TIBCKeeper.RoutingKeeper.Authenticate(suite.chain.GetContext(), tc.source, tc.dest, tc.port)
			if tc.expPass {
				suite.Require().True(ok)
			} else {
				suite.Require().False(ok)
			}
		})
	}
}
