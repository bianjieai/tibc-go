package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NftKeeper defines the expected nft keeper
type NftKeeper interface {
	MintNFT(ctx sdk.Context, denomID, tokenID, tokenNm, tokenURI, tokenData string, owner sdk.AccAddress) error
	BurnNFT(ctx sdk.Context, denomID, tokenID string, owner sdk.AccAddress) error
	TransferOwner(ctx sdk.Context, denomID, tokenID, tokenNm, tokenURI, tokenData string, srcOwner, dstOwner sdk.AccAddress) error
}


// AccountKeeper defines the expected account keeper
type AccountKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
}