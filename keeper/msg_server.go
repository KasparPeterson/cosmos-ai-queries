package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	"github.com/kasparpeterson/cosmos_ai_queries"
)

type msgServer struct {
	k Keeper
}

var _ cosmos_ai_queries.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the module MsgServer interface.
func NewMsgServerImpl(keeper Keeper) cosmos_ai_queries.MsgServer {
	return &msgServer{k: keeper}
}

// CreateQuery defines the handler for the MsgCreateQuery message.
func (ms msgServer) CreateQuery(ctx context.Context, msg *cosmos_ai_queries.MsgCreateQuery) (*cosmos_ai_queries.MsgCreateQueryResponse, error) {
	if length := len([]byte(msg.Index)); cosmos_ai_queries.MaxIndexLength < length || length < 1 {
		return nil, cosmos_ai_queries.ErrIndexTooLong
	}
	if _, err := ms.k.StoredQueries.Get(ctx, msg.Index); err == nil || errors.Is(err, collections.ErrEncoding) {
		return nil, fmt.Errorf("query already exists at index: %s", msg.Index)
	}

	storedQuery := cosmos_ai_queries.StoredQuery{
		Query:  msg.Query,
		User:   msg.Creator,
		Answer: "",
	}
	if err := storedQuery.Validate(); err != nil {
		return nil, err
	}
	if err := ms.k.StoredQueries.Set(ctx, msg.Index, storedQuery); err != nil {
		return nil, err
	}

	return &cosmos_ai_queries.MsgCreateQueryResponse{}, nil
}

func (ms msgServer) PostQueryAnswer(ctx context.Context, msg *cosmos_ai_queries.MsgPostQueryAnswer) (*cosmos_ai_queries.MsgPostQueryAnswerResponse, error) {
	query, err := ms.k.StoredQueries.Get(ctx, msg.Index)
	if err != nil {
		return nil, fmt.Errorf("query does not exist at index: %s", msg.Index)
	}
	query.Answer = msg.Answer
	if err := ms.k.StoredQueries.Set(ctx, msg.Index, query); err != nil {
		return nil, fmt.Errorf("failed to add answer at index: %s", msg.Index)
	}

	return &cosmos_ai_queries.MsgPostQueryAnswerResponse{}, nil
}
