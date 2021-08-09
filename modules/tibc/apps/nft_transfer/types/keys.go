package types


const (
	// ModuleName defines the TIBC nft_transfer name
	ModuleName = "nftTransfer"

	// RouterKey is the message route for the nft-transfer module
	RouterKey = ModuleName

	// PortID is the default port id that nft-transfer module binds to
	PortID = "nftTransfer"

	// StoreKey is the store key string for TIBC nft-transfer
	StoreKey = 	ModuleName
	// QuerierRoute is the querier route for TIBC nft-transfer
	QuerierRoute = ModuleName
)
