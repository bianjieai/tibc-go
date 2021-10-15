package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagRelayChain = "relay-chain-name"
)

var (
	FsSendCleanPacket = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsSendCleanPacket.String(FlagRelayChain, "", "The name of relay chain")
}
