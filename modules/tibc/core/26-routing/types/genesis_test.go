package types_test

import (
	"fmt"

	"github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
)

func (suite *TypesTestSuite) TestValidateGenesis() {

	testCases := []struct {
		name     string
		genState types.GenesisState
		expPass  bool
	}{
		{
			name:     "empty",
			genState: types.NewGenesisState(nil),
			expPass:  true,
		},
		{
			name:     "invalid genesis1",
			genState: types.NewGenesisState([]string{"1a"}),
			expPass:  false,
		},
		{
			name:     "invalid genesis2",
			genState: types.NewGenesisState([]string{"1a,2b"}),
			expPass:  false,
		},
		{
			name:     "valid genesis",
			genState: types.NewGenesisState([]string{fmt.Sprintf("source,dest,port")}),
			expPass:  true,
		},
	}
	for _, tc := range testCases {
		tc := tc
		err := tc.genState.Validate()
		if tc.expPass {
			suite.Require().NoError(err, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
	}
}
