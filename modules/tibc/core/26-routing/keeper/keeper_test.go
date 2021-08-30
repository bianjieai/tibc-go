package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	ibctesting "github.com/bianjieai/tibc-go/modules/tibc/testing"
)

var (
	validPort   = "validportid"
	invalidPort = "(invalidPortID)"
)

type KeeperTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator
	chain       *ibctesting.TestChain
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 1)
	suite.chain = suite.coordinator.GetChain(ibctesting.GetChainID(0))
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
