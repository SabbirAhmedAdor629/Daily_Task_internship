* Overview

Ensure the following are set in your unit tests;

``` json
	os.Setenv("API_KEYS_TABLE_NAME", "api_keys_staging")
	env = NewEnv()

	os.Setenv("REQUEST_API_KEYS", "{ \"api_keys\": [ \"test-3\" ] }")
	apiKeys = im_keycache.LoadCache(im_keycache.GetApiKeys().ApiKeys)
```
