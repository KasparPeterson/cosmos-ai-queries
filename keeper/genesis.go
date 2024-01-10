package keeper

import (
	"context"

	"github.com/kasparpeterson/cosmos_ai_queries"
)

// InitGenesis initializes the module state from a genesis state.
func (k *Keeper) InitGenesis(ctx context.Context, data *cosmos_ai_queries.GenesisState) error {
	if err := k.Params.Set(ctx, data.Params); err != nil {
		return err
	}

	for _, indexedStoredQuery := range data.IndexedStoredQueryList {
		if err := k.StoredQueries.Set(ctx, indexedStoredQuery.Index, indexedStoredQuery.StoredQuery); err != nil {
			return err
		}
	}

	return nil
}

// ExportGenesis exports the module state to a genesis state.
func (k *Keeper) ExportGenesis(ctx context.Context) (*cosmos_ai_queries.GenesisState, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	var indexedStoredQueries []cosmos_ai_queries.IndexedStoredQuery
	if err := k.StoredQueries.Walk(ctx, nil, func(index string, storedQuery cosmos_ai_queries.StoredQuery) (bool, error) {
		indexedStoredQueries = append(indexedStoredQueries, cosmos_ai_queries.IndexedStoredQuery{
			Index:       index,
			StoredQuery: storedQuery,
		})
		return false, nil
	}); err != nil {
		return nil, err
	}

	return &cosmos_ai_queries.GenesisState{
		Params:                 params,
		IndexedStoredQueryList: indexedStoredQueries,
	}, nil
}
