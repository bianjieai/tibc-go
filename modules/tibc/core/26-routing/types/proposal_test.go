package types_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
	ibctesting "github.com/bianjieai/tibc-go/modules/tibc/testing"
)

type TypesTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	chain *ibctesting.TestChain
}

func (suite *TypesTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 1)
	suite.chain = suite.coordinator.GetChain(ibctesting.GetChainID(0))
}

func TestTypesTestSuite(t *testing.T) {
	suite.Run(t, new(TypesTestSuite))
}

func (suite *TypesTestSuite) TestNewSetRoutingRulesProposal() {
	p, err := types.NewSetRoutingRulesProposal(ibctesting.Title, ibctesting.Description, []string{"source.dest.dgsbl"})
	suite.Require().NoError(err)
	suite.Require().NotNil(p)
}
