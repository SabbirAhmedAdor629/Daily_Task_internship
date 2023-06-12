package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

func createResponseBody(errorCode int, agreements []Agreement) APIResponseBody {
	response := ResponseObject{
		Error: errorCode,
		Data:  agreements,
	}
	body, _ := json.Marshal(response)
	return APIResponseBody(body)
}

func GenerateApiResponse(errorCode int, agreements []Agreement, statusCode int, responseHeaders map[string]string) *events.APIGatewayProxyResponse {
	body := createResponseBody(errorCode, agreements)
	// remove the api key from response headers
	delete(responseHeaders, "api_key")
	return &events.APIGatewayProxyResponse{
		StatusCode:      statusCode,
		Headers:         responseHeaders,
		Body:            string(body),
		IsBase64Encoded: false,
	}
}
