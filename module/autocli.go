package module

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	cosmosaiqueriesv1 "github.com/kasparpeterson/cosmos_ai_queries/api/v1"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: cosmosaiqueriesv1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "GetQuery",
					Use:       "get-query index",
					Short:     "Get the current answer of the query at index",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "index"},
					},
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: cosmosaiqueriesv1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "CreateQuery",
					Use:       "create index query",
					Short:     "Creates a new AI query at the index",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "index"},
						{ProtoField: "query"},
					},
				},
			},
		},
	}
}
