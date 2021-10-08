package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	nftexported "github.com/irisnet/irismod/modules/nft/exported"
	"github.com/irisnet/irismod/modules/nft/types"

	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

// NftKeeper defines the expected nft keeper
type NftKeeper interface {
	MintNFT(ctx sdk.Context, denomID, tokenID, tokenNm, tokenURI, tokenData string, owner sdk.AccAddress) error
	BurnNFT(ctx sdk.Context, denomID, tokenID string, owner sdk.AccAddress) error
	GetNFT(ctx sdk.Context, denomID, tokenID string) (nft nftexported.NFT, err error)
	TransferOwner(ctx sdk.Context, denomID, tokenID, tokenNm, tokenURI, tokenData string, srcOwner, dstOwner sdk.AccAddress) error
	GetDenom(ctx sdk.Context, id string) (denom types.Denom, found bool)
	IssueDenom(ctx sdk.Context, id, name, schema, symbol string, creator sdk.AccAddress, mintRestricted, updateRestricted bool) error
}

// AccountKeeper defines the expected account keeper
type AccountKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
}

// PacketKeeper defines the expected packet keeper
type PacketKeeper interface {
	GetNextSequenceSend(ctx sdk.Context, sourceChain, destChain string) uint64
	SendPacket(ctx sdk.Context, packet exported.PacketI) error
}

// ClientKeeper defines the expected client keeper
type ClientKeeper interface {
	GetChainName(ctx sdk.Context) string
}
