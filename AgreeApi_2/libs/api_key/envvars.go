package api_key

import (
	"encoding/json"
	"os"
)

var ENV_REQUEST_API_KEYS string = "REQUEST_API_KEYS"

type RequestApiKeys struct {
	ApiKeys []string `json:"api_keys"`
}

func GetApiKeys() RequestApiKeys {
	apiKeys := RequestApiKeys{}
	_ = json.Unmarshal([]byte(os.Getenv(ENV_REQUEST_API_KEYS)), &apiKeys)

	var keys []string
	for _, v := range apiKeys.ApiKeys {
		if v == "" {
			continue
		}
		keys = append(keys, v)
	}
	apiKeys.ApiKeys = keys
	return apiKeys
}
