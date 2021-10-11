package types

const (
	// SubModuleName defines the TIBC port name
	SubModuleName = "routing"

	// StoreKey is the store key string for TIBC ports
	StoreKey = SubModuleName

	// RouterKey is the message route for TIBC ports
	RouterKey = SubModuleName

	// QuerierRoute is the querier route for TIBC ports
	QuerierRoute = SubModuleName
)

const (
	// RulePattern format "source,dest,port"
	RulePattern = "^(([a-zA-Z0-9\\.\\_\\+\\-\\#\\[\\]\\<\\>]{1,64}|[*]),){2}([a-zA-Z0-9\\.\\_\\+\\-\\#\\[\\]\\<\\>]{1,64}|[*])$"
)
