package types

const (
	// ModuleName defines the TIBC mt_transfer name
	ModuleName = "MT"

	// RouterKey is the message route for the mt-transfer module
	RouterKey = ModuleName

	// PortID is the default port id that mt-transfer module binds to
	PortID = ModuleName

	// StoreKey is the store key string for TIBC mt-transfer
	StoreKey = ModuleName

	// QuerierRoute is the querier route for TIBC mt-transfer
	QuerierRoute = ModuleName

	// ClassPrefix is the prefix used for mt class.
	ClassPrefix = "tibc"
)

var (
	// ClassTraceKey defines the key to store the class trace info in store
	ClassTraceKey = []byte{0x01}
)
