package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

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
