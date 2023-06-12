package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	v3 "influencemobile.com/libs/v3_helpers"
)

type IPQSResponse struct {
	Timezone string `json:"timezone"`
	Success  bool   `json:"success"`
	Message  string `json:"message"`
}

func GetTimezoneFromIPQS(apiKey string, ip string) (*IPQSResponse, int, error) {
	requestURL := fmt.Sprintf("https://ipqualityscore.com/api/json/ip/%s/%s", apiKey, ip)

	request, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, int(v3.ResponseInternalServerError), err
	}

	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, int(v3.ResponseInternalServerError), err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, int(v3.ResponseInternalServerError), err
	}

	var response IPQSResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, int(v3.ResponseInternalServerError), err
	}

	if !response.Success {
		return nil, int(v3.ResponseInternalServerError), fmt.Errorf("IPQS %s", response.Message)
	}

	return &response, int(v3.ResponseSuccess), nil
}
