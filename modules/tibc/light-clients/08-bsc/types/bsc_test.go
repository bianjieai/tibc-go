package types_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bianjieai/tibc-go/simapp"
)

type BSCTestSuite struct {
	suite.Suite

	ctx sdk.Context
	app *simapp.SimApp
}

func (suite *BSCTestSuite) SetupTest() {
	app := simapp.Setup(suite.T())

	suite.ctx = app.BaseApp.NewContextLegacy(false, tmproto.Header{})
	suite.app = app
}

func TestBSCTestSuite(t *testing.T) {
	suite.Run(t, new(BSCTestSuite))
}
