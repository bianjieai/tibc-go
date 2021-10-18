package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagRelayChain   = "relay-chain"
	FlagDestContract = "dest-contract"
)

var (
	FsNftTransfer = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsNftTransfer.String(FlagRelayChain, "", "relay chain used by cross-chain NFT")
	FsNftTransfer.String(FlagDestContract, "", "the destination contract address to receive the nft")
}
