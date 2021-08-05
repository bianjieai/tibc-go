package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/client/cli"
)

var ProposalHandler = govclient.NewProposalHandler(cli.NewCmdCreateClientProposal, nil)
