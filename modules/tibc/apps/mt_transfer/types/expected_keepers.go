package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	mtexported "github.com/irisnet/irismod/modules/mt/exported"
	"github.com/irisnet/irismod/modules/mt/types"

	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

// MtKeeper defines the expected mt keeper
type MtKeeper interface {
	IssueDenom(ctx sdk.Context,
		id, name string, sender sdk.AccAddress, data []byte,
	) types.Denom

	IssueMT(ctx sdk.Context,
		denomID string, mtID string,
		amount uint64,
		data []byte,
		recipient sdk.AccAddress,
	) (types.MT, error)

	MintMT(ctx sdk.Context,
		denomID, mtID string,
		amount uint64,
		recipient sdk.AccAddress,
	) error

	TransferOwner(ctx sdk.Context,
		denomID, mtID string,
		amount uint64,
		srcOwner, dstOwner sdk.AccAddress,
	) error

	BurnMT(ctx sdk.Context,
		denomID, mtID string,
		amount uint64,
		owner sdk.AccAddress) error

	HasMT(ctx sdk.Context, denomID, mtID string) bool
	GetMT(ctx sdk.Context, denomID, mtID string) (mt mtexported.MT, err error)
	GetDenom(ctx sdk.Context, id string) (denom types.Denom, found bool)
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
