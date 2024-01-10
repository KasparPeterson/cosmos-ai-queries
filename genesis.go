package cosmos_ai_queries

// NewGenesisState creates a new genesis state with default values.
func NewGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
func (gs *GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	unique := make(map[string]bool)
	for _, indexedStoredQuery := range gs.IndexedStoredQueryList {
		if length := len([]byte(indexedStoredQuery.Index)); MaxIndexLength < length || length < 1 {
			return ErrIndexTooLong
		}
		if _, ok := unique[indexedStoredQuery.Index]; ok {
			return ErrDuplicateAddress
		}
		if err := indexedStoredQuery.StoredQuery.Validate(); err != nil {
			return err
		}
		unique[indexedStoredQuery.Index] = true
	}

	return nil
}
