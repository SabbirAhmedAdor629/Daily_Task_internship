package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestGetQueryStringParams(t *testing.T) {
	type args struct {
		request events.APIGatewayProxyRequest
	}
	type test struct {
		name         string
		args         args
		expectedResp *QueryStringParams
		expectedErr  error
	}
	tests := []test{
		{
			name: "#1. Valid QueryStringParams",
			args: args{
				request: events.APIGatewayProxyRequest{
					QueryStringParameters: map[string]string{
						"timezone": "Europe/London",
					},
				},
			},
			expectedResp: &QueryStringParams{
				Timezone: "Europe/London",
			},
			expectedErr: nil,
		},
		{
			name: "#2. Empty QueryStringParams",
			args: args{
				request: events.APIGatewayProxyRequest{
					QueryStringParameters: map[string]string{},
				},
			},
			expectedResp: &QueryStringParams{},
			expectedErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getQueryStringParams(tt.args.request)
			if !reflect.DeepEqual(got, tt.expectedResp) {
				t.Errorf("getQueryStringParams() got = %v, expected = %v", got, tt.expectedResp)
			}
			if err != tt.expectedErr {
				t.Errorf("getQueryStringParams() error = %v, expected error = %v", err, tt.expectedErr)
			}
		})
	}
}

func TestSetConsentUrl(t *testing.T) {
	type args struct {
		hostUrl      string
		jurisdiction string
		language     string
		documentType string
		version      int
	}
	type test struct {
		name        string
		args        args
		expectedUrl string
	}
	tests := []test{
		{
			name: "#1. Valid input",
			args: args{
				hostUrl:      "https://example.com",
				jurisdiction: "gdpr",
				language:     "en-uk",
				documentType: "pp",
				version:      1680166802,
			},
			expectedUrl: "https://example.com/agreement/gdpr/en-uk/pp/1680166802/pp.html",
		},
		{
			name: "#2. Valid input",
			args: args{
				hostUrl:      "https://example.com",
				jurisdiction: "gdpr",
				language:     "en-ca",
				documentType: "pp",
				version:      1680166789,
			},
			expectedUrl: "https://example.com/agreement/gdpr/en-ca/pp/1680166789/pp.html",
		},
		{
			name: "#3. Valid input",
			args: args{
				hostUrl:      "https://example.com",
				jurisdiction: "gdpr",
				language:     "en-uk",
				documentType: "tou",
				version:      1680166788,
			},
			expectedUrl: "https://example.com/agreement/gdpr/en-uk/tou/1680166788/tou.html",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SetConsentUrl(tt.args.hostUrl, tt.args.jurisdiction, tt.args.language, tt.args.documentType, tt.args.version)
			if got != tt.expectedUrl {
				t.Errorf("SetConsentUrl() got = %v, expected = %v", got, tt.expectedUrl)
			}
		})
	}
}

// NEW FUNCTIONS
func TestGetPlayerID(t *testing.T) {
	type args struct {
		request events.APIGatewayProxyRequest
	}
	type test struct {
		name     string
		args     args
		expected string
	}
	tests := []test{
		{
			name:     "#1. Valid player request",
			args:     args{request: events.APIGatewayProxyRequest{Resource: "", Path: "https://pgnzdf1s4m.execute-api.us-east-1.amazonaws.com/beta/12001", HTTPMethod: "", Headers: map[string]string{}, MultiValueHeaders: map[string][]string{}, QueryStringParameters: map[string]string{"locale": "en-uk", "timezone": "Europe/London"}, MultiValueQueryStringParameters: map[string][]string{}, PathParameters: map[string]string{}, StageVariables: map[string]string{}, RequestContext: events.APIGatewayProxyRequestContext{}, Body: "", IsBase64Encoded: false}},
			expected: "12001",
		},
		{
			name:     "#2. Valid player request",
			args:     args{request: events.APIGatewayProxyRequest{Resource: "", Path: "https://pgnzdf1s4m.execute-api.us-east-1.amazonaws.com/beta/12002", HTTPMethod: "", Headers: map[string]string{}, MultiValueHeaders: map[string][]string{}, QueryStringParameters: map[string]string{"locale": "en-uk", "timezone": "Europe/London"}, MultiValueQueryStringParameters: map[string][]string{}, PathParameters: map[string]string{}, StageVariables: map[string]string{}, RequestContext: events.APIGatewayProxyRequestContext{}, Body: "", IsBase64Encoded: false}},
			expected: "12002",
		},
		{
			name:     "#3. Valid player request",
			args:     args{request: events.APIGatewayProxyRequest{Resource: "", Path: "https://pgnzdf1s4m.execute-api.us-east-1.amazonaws.com/beta/12", HTTPMethod: "", Headers: map[string]string{}, MultiValueHeaders: map[string][]string{}, QueryStringParameters: map[string]string{"locale": "en-uk", "timezone": "Europe/London"}, MultiValueQueryStringParameters: map[string][]string{}, PathParameters: map[string]string{}, StageVariables: map[string]string{}, RequestContext: events.APIGatewayProxyRequestContext{}, Body: "", IsBase64Encoded: false}},
			expected: "12",
		},
		{
			name:     "#3. InValid player request",
			args:     args{request: events.APIGatewayProxyRequest{Resource: "", Path: "https://pgnzdf1s4m.execute-api.us-east-1.amazonaws.com/beta/", HTTPMethod: "", Headers: map[string]string{}, MultiValueHeaders: map[string][]string{}, QueryStringParameters: map[string]string{"locale": "en-uk", "timezone": "Europe/London"}, MultiValueQueryStringParameters: map[string][]string{}, PathParameters: map[string]string{}, StageVariables: map[string]string{}, RequestContext: events.APIGatewayProxyRequestContext{}, Body: "", IsBase64Encoded: false}},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := getPlayerID(tt.args.request)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("getPlayerID() = %v, want %v", got, tt.expected)
			}
		})
	}

}

func TestValidTimezonePrefix(t *testing.T) {
	testCases := []struct {
		name     string
		timezone string
		prefix   string
		expected bool
	}{
		{
			name:     "#1. Matching prefix",
			timezone: "Europe/London",
			prefix:   "europe/",
			expected: true,
		},
		{
			name:     "#2. Non-matching prefix",
			timezone: "Asia/Tokyo",
			prefix:   "europe/",
			expected: false,
		},
		{
			name:     "#3. Case-insensitive matching prefix",
			timezone: "Europe/London",
			prefix:   "EUROPE/",
			expected: true,
		},
		{
			name:     "#4. Empty timezone",
			timezone: "",
			prefix:   "europe/",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := validTimezonePrefix(tc.timezone, tc.prefix)
			if result != tc.expected {
				t.Errorf("Expected %t, but got %t", tc.expected, result)
			}
		})
	}
}

func TestGetIpAddress(t *testing.T) {
	testCases := []struct {
		name          string
		xForwardedFor string
		expected      string
		expectedError error
	}{
		{
			name:          "#1: Valid IP address",
			xForwardedFor: "192.168.0.1",
			expected:      "192.168.0.1",
			expectedError: nil,
		},
		{
			name:          "#2: Invalid IP address",
			xForwardedFor: "invalid",
			expected:      "",
			expectedError: fmt.Errorf("invalid IP address"),
		},
		{
			name:          "#3: Valid X-Forwarded-For header with multiple IP addresses",
			xForwardedFor: "192.168.0.1, 10.0.0.1",
			expected:      "192.168.0.1",
			expectedError: nil,
		},
		{
			name:          "#4: Valid X-Forwarded-For header with multiple IP addresses and spaces",
			xForwardedFor: "192.168.0.1 , 10.0.0.1",
			expected:      "192.168.0.1",
			expectedError: nil,
		},
		{
			name:          "#5: Valid X-Forwarded-For header with IPv6 address",
			xForwardedFor: "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			expected:      "2001:db8:85a3::8a2e:370:7334",
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ip, err := GetIpAddress(tc.xForwardedFor)
			if err != nil && tc.expectedError == nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if err == nil && tc.expectedError != nil {
				t.Error("Expected error, but got nil")
			}
			if ip != tc.expected {
				t.Errorf("Unexpected IP address: %s", ip)
			}
		})
	}
}
