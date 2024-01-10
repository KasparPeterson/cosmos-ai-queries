# cosmos-ai-queries

Logic:

1. Send tx with new query str
2. Able to query the query str
3. Send tx to add a result to query

minid tx cosmos_ai_queries create id7 hello --from alice --yes
minid query cosmos_ai_queries get-query id7
minid tx cosmos_ai_queries post id7 world --from alice --yes