package types

const (
	// SubModuleName defines the IBC client name
	SubModuleName string = "client"

	// RouterKey is the message route for IBC client
	RouterKey string = SubModuleName

	// QuerierRoute is the querier route for IBC client
	QuerierRoute string = SubModuleName

	// KeyClientName is the key used to store the chain name in
	// the keeper.
	KeyClientName = "chainName"

	// KeyRelayers is the key used to store the relayers address in
	// the keeper.
	KeyRelayers = "relayers"
)
