package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagRelayChain = "relayChain"
	FlagAwayFromChain = "awayFromChain"
)

var (
	FsNftTransfer    = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsNftTransfer.String(FlagRelayChain, "", "relay chain used by cross-chain NFT")
	FsNftTransfer.String(FlagAwayFromChain, "", "indicates whether nft is far away from the source")
}

