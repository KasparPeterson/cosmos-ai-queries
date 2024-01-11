package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kasparpeterson/cosmos_ai_queries"
)

var _ cosmos_ai_queries.QueryServer = queryServer{}

// NewQueryServerImpl returns an implementation of the module QueryServer.
func NewQueryServerImpl(k Keeper) cosmos_ai_queries.QueryServer {
	return queryServer{k}
}

type queryServer struct {
	k Keeper
}

// GetQuery defines the handler for the Query/GetQuery RPC method.
func (qs queryServer) GetQuery(ctx context.Context, req *cosmos_ai_queries.QueryGetQueryRequest) (*cosmos_ai_queries.QueryGetQueryResponse, error) {
	query, err := qs.k.StoredQueries.Get(ctx, req.Index)
	if err == nil {
		return &cosmos_ai_queries.QueryGetQueryResponse{Query: &query}, nil
	}
	if errors.Is(err, collections.ErrNotFound) {
		return &cosmos_ai_queries.QueryGetQueryResponse{Query: nil}, nil
	}

	return nil, status.Error(codes.Internal, err.Error())
}
