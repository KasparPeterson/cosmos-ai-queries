# cosmos-ai-queries

Logic:

1. Send tx with new query str
2. Able to query the query str
3. Send tx to add a result to query

```shell
minid tx cosmos_ai_queries create id8 "hello wollrd" --from alice --yes
minid query cosmos_ai_queries get-query id8
minid tx cosmos_ai_queries post id8 "hello to you too" --from alice --yes

minid tx cosmos_ai_queries createsync id14 "what is the capital of Latvia?" --from alice --yes
minid query cosmos_ai_queries get-query id14
```