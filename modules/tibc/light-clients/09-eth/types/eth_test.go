package types_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bianjieai/tibc-go/simapp"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

type ETHTestSuite struct {
	suite.Suite

	ctx sdk.Context
	app *simapp.SimApp
}

func (suite *ETHTestSuite) SetupTest() {
	app := simapp.Setup(false)

	suite.ctx = app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now()})
	suite.app = app
}

func TestETHTestSuite(t *testing.T) {
	suite.Run(t, new(ETHTestSuite))
}
