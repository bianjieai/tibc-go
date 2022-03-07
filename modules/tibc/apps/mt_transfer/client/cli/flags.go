package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagRelayChain   = "relay-chain"
	FlagDestContract = "dest-contract"
)

var (
	FsMtTransfer = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsMtTransfer.String(FlagRelayChain, "", "relay chain used by cross-chain mt")
	FsMtTransfer.String(FlagDestContract, "", "the destination contract address to receive the mt")
}
