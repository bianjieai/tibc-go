package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagSourceChain = "source-chain-name"
	FlagRelayChain  = "relay-chain-name"
)

var (
	FsSendCleanPacket = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsSendCleanPacket.String(FlagSourceChain, "", "Denom data structure definition")
	FsSendCleanPacket.String(FlagRelayChain, "", "The name of the denom")
}
