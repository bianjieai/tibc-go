package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagRelayChain = "relay-chain"
)

var (
	FsNftTransfer = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsNftTransfer.String(FlagRelayChain, "", "relay chain used by cross-chain NFT")
}
