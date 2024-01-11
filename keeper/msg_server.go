package keeper

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

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

func (ms msgServer) CreateSyncQuery(ctx context.Context, msg *cosmos_ai_queries.MsgCreateSyncQuery) (*cosmos_ai_queries.MsgCreateSyncQueryResponse, error) {
	if length := len([]byte(msg.Index)); cosmos_ai_queries.MaxIndexLength < length || length < 1 {
		return nil, cosmos_ai_queries.ErrIndexTooLong
	}
	if _, err := ms.k.StoredQueries.Get(ctx, msg.Index); err == nil || errors.Is(err, collections.ErrEncoding) {
		return nil, fmt.Errorf("query already exists at index: %s", msg.Index)
	}

	var answer = getLlmAnswer(msg.Query)
	storedQuery := cosmos_ai_queries.StoredQuery{
		Query:  msg.Query,
		User:   msg.Creator,
		Answer: answer,
	}
	if err := storedQuery.Validate(); err != nil {
		return nil, err
	}
	if err := ms.k.StoredQueries.Set(ctx, msg.Index, storedQuery); err != nil {
		return nil, err
	}

	return &cosmos_ai_queries.MsgCreateSyncQueryResponse{
		Answer: answer,
	}, nil
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

type MessageRequest struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type RequestBody struct {
	Model    string           `json:"model"`
	Messages []MessageRequest `json:"messages"`
}

type Message struct {
	Role       string      `json:"role"`
	Content    string      `json:"content"`
	ToolCalls  interface{} `json:"tool_calls"`
	ToolCallID interface{} `json:"tool_call_id"`
}

type Choice struct {
	Message      Message `json:"message"`
	Index        int     `json:"index"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Post struct {
	Id      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

func getLlmAnswer(query string) (answer string) {
	posturl := "https://api.llmos.dev/v1/chat/completions"

	requestBody := RequestBody{
		Model: "gpt-3.5-turbo",
		Messages: []MessageRequest{
			{
				Content: "You are a helpful assistant.",
				Role:    "system",
			},
			{
				Content: query, // Replace "Hello!" with your custom content
				Role:    "user",
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Authorization", "AddYOUR KEY")

	client := &http.Client{}
	res, err := client.Do(r)

	fmt.Println("go some response")

	if err != nil {
		panic(err)
	}

	fmt.Println("closing")
	defer res.Body.Close()

	fmt.Println("parsing json")

	post := &Post{}
	derr := json.NewDecoder(res.Body).Decode(post)
	fmt.Println("derr %s:", derr)
	if derr != nil {
		panic(derr)
	}

	fmt.Println("Id:", post.Id)
	fmt.Println("Object:", post.Object)
	fmt.Println("Model:", post.Model)
	fmt.Println("Choices:", post.Choices)
	fmt.Println("Choice0.Message.Role:", post.Choices[0].Message.Role)
	fmt.Println("Choice0.Message.Content:", post.Choices[0].Message.Content)

	return post.Choices[0].Message.Content
}
