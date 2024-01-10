package cosmos_ai_queries

import "cosmossdk.io/collections"

const ModuleName = "cosmos_ai_queries"
const MaxIndexLength = 256

var (
	ParamsKey        = collections.NewPrefix("Params")
	StoredQueriesKey = collections.NewPrefix("StoredQueries/value/")
)
