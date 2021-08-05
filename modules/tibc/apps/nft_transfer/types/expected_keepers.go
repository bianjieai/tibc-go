package types

import (
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irismod/modules/nft/types"
)

// NftKeeper defines the expected nft keeper
type NftKeeper interface {
	MintNFT(ctx sdk.Context, denomID, tokenID, tokenNm, tokenURI, tokenData string, owner sdk.AccAddress) error
	BurnNFT(ctx sdk.Context, denomID, tokenID string, owner sdk.AccAddress) error
	TransferOwner(ctx sdk.Context, denomID, tokenID, tokenNm, tokenURI, tokenData string, srcOwner, dstOwner sdk.AccAddress) error
	GetDenom(ctx sdk.Context, id string) (denom types.Denom, found bool)
}


// AccountKeeper defines the expected account keeper
type AccountKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
}


// PacketKeeper defines the expected packet keeper
type PacketKeeper interface {
	// todo getSequence

	SendPacket(ctx sdk.Context, packet exported.PacketI) error
}

// ClientKeeper defines the expected client keeper
type ClientKeeper interface {
	GetChainName(ctx sdk.Context) string
}
