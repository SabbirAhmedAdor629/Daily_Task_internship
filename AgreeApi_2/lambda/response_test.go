package main

import (
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestGenerateApiResponse(t *testing.T) {
	type args struct {
		errorCode  int
		agreements []Agreement
		statusCode int
		headers    map[string]string
	}
	type test struct {
		name     string
		args     args
		expected *events.APIGatewayProxyResponse
	}
	tests := []test{
		{
			name: "#1. Returns expected response with 200 status code [Multiple Agreements]",
			args: args{
				errorCode: 0,
				agreements: []Agreement{{
					Jurisdiction: "gdpr",
					Agreement:    "tou",
					Version:      1680166789,
					Locale:       "en-uk",
					URL:          "https://example.com/agreement/gdpr/en-uk/tou/1680166789/tou.html",
				}, {
					Jurisdiction: "gdpr",
					Agreement:    "pp",
					Version:      1680166799,
					Locale:       "en-us",
					URL:          "https://example.com/agreement/gdpr/en-uk/pp/1680166799/pp.html",
				}},
				statusCode: 200,
				headers: map[string]string{
					"api_key":    "test-3",
					"request_id": "test-request-id",
				},
			},
			expected: &events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"request_id": "test-request-id",
				},
				Body:            `{"error":0,"data":[{"jurisdiction":"gdpr","agreement":"tou","version":1680166789,"locale":"en-uk","url":"https://example.com/agreement/gdpr/en-uk/tou/1680166789/tou.html"},{"jurisdiction":"gdpr","agreement":"pp","version":1680166799,"locale":"en-us","url":"https://example.com/agreement/gdpr/en-uk/pp/1680166799/pp.html"}]}`,
				IsBase64Encoded: false,
			},
		},
		{
			name: "#2. Returns expected response with 200 status code [Single Agreements]",
			args: args{
				errorCode: 0,
				agreements: []Agreement{{
					Jurisdiction: "gdpr",
					Agreement:    "pp",
					Version:      1680166779,
					Locale:       "en-uk",
					URL:          "https://example.com/agreement/gdpr/en-uk/pp/1680166779/tou.html",
				}},
				statusCode: 200,
				headers: map[string]string{
					"request_id": "test-request-id",
				},
			},
			expected: &events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"request_id": "test-request-id",
				},
				Body:            `{"error":0,"data":[{"jurisdiction":"gdpr","agreement":"pp","version":1680166779,"locale":"en-uk","url":"https://example.com/agreement/gdpr/en-uk/pp/1680166779/tou.html"}]}`,
				IsBase64Encoded: false,
			},
		},
		{
			name: "#3. Returns expected response with 200 status code",
			args: args{
				errorCode:  12,
				agreements: nil,
				statusCode: 200,
				headers: map[string]string{
					"request_id": "test-request-id",
				},
			},
			expected: &events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"request_id": "test-request-id",
				},
				Body:            `{"error":12,"data":null}`,
				IsBase64Encoded: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateApiResponse(tt.args.errorCode, tt.args.agreements, tt.args.statusCode, tt.args.headers)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("GenerateApiResponse() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
