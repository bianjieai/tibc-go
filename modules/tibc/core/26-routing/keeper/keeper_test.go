package keeper_test

import (
	"testing"

	ibctesting "github.com/bianjieai/tibc-go/modules/tibc/testing"
	"github.com/stretchr/testify/suite"
)

var (
	validPort   = "validportid"
	invalidPort = "(invalidPortID)"
)

type KeeperTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator
	chain *ibctesting.TestChain
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)
	suite.chain = suite.coordinator.GetChain(ibctesting.GetChainID(0))
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
