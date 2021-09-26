package types

const (
	// ModuleName defines the TIBC nft_transfer name
	ModuleName = "NFT"

	// RouterKey is the message route for the nft-transfer module
	RouterKey = ModuleName

	// PortID is the default port id that nft-transfer module binds to
	PortID = ModuleName

	// StoreKey is the store key string for TIBC nft-transfer
	StoreKey = ModuleName

	// QuerierRoute is the querier route for TIBC nft-transfer
	QuerierRoute = ModuleName

	// ClassPrefix is the prefix used for nft class.
	ClassPrefix = "tibc"
)

var (
	// ClassTraceKey defines the key to store the class trace info in store
	ClassTraceKey = []byte{0x01}
)
