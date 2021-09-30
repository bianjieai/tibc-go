package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagSourceChain = "source-chain"
	FlagRelayChain  = "relay-chain"
)

var (
	FsSendCleanPacket = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsSendCleanPacket.String(FlagRelayChain, "", "The name of relay chain")
}
