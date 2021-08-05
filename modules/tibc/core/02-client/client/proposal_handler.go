package client

import (
	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/client/cli"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

var ProposalHandler = govclient.NewProposalHandler(cli.NewCmdCreateClientProposal, nil)
