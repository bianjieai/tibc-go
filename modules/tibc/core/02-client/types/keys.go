package types

const (
	// SubModuleName defines the TIBC client name
	SubModuleName string = "client"

	// RouterKey is the message route for TIBC client
	RouterKey string = SubModuleName

	// QuerierRoute is the querier route for TIBC client
	QuerierRoute string = SubModuleName

	// KeyClientName is the key used to store the chain name in the keeper.
	KeyClientName = "chainName"

	// KeyRelayers is the key used to store the relayers address in the keeper.
	KeyRelayers = "relayers"
)
