package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
)

func (suite *KeeperTestSuite) TestQueryPacketCommitment() {
	var (
		req   *types.QueryRoutingRulesRequest
		rules []string
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{{
		"empty request",
		func() {
			req = nil
		},
		false,
	}, {
		"success",
		func() {
			rules = []string{"source,dest,dgsbl"}
			suite.chain.App.TIBCKeeper.RoutingKeeper.SetRoutingRules(suite.chain.GetContext(), rules)
			req = &types.QueryRoutingRulesRequest{}
		},
		true,
	}}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset
			tc.malleate()
			ctx := sdk.WrapSDKContext(suite.chain.GetContext())

			res, err := suite.chain.QueryServer.RoutingRules(ctx, req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(rules, res.Rules)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
