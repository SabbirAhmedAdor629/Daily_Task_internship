package main

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	v3 "influencemobile.com/libs/v3_helpers"
	im_config "influencemobile.com/parameter_store"
)

func setEnvVars(es map[TypeEnvVar]string) {
	for k, v := range es {
		os.Setenv(EnvVars[k], v)
	}
}

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

func TestValidateQueryStringParams(t *testing.T) {
	type args struct {
		request events.APIGatewayProxyRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *QueryStringParams
		want1   int
		wantErr bool
	}{
		{
			name: "#1: Valid query string params timezone [Europe/London]",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/required_agreements",
					Path:       "/api/v1/required_agreements",
					HTTPMethod: "GET",
					Headers: map[string]string{
						"Content-Type": "application/json",
						"api_key":      "test-3",
						"request_id":   "C713F659-1ABF-4EDF-8E86-655295CF7CD7",
					},
					QueryStringParameters: map[string]string{
						"timezone": "Europe/London",
					},
					PathParameters: map[string]string{
						"proxy": "example",
					},
					Body:            "",
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
			},
			want: &QueryStringParams{
				Timezone: "Europe/London",
			},
			want1:   int(v3.ResponseSuccess),
			wantErr: false,
		},
		{
			name: "#2: Invalid query string params [timezone is empty]",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/required_agreements",
					Path:       "/api/v1/required_agreements/old",
					HTTPMethod: "GET",
					Headers: map[string]string{
						"Content-Type": "application/json",
						"api_key":      "test-3",
						"request_id":   "C713F659-1ABF-4EDF-8E86-655295CF7CD7",
					},
					QueryStringParameters: map[string]string{
						"timezone": "",
					},
					PathParameters: map[string]string{
						"proxy": "example",
					},
					Body:            "",
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
			},
			want:    &QueryStringParams{},
			want1:   int(ResponseInvalidTimezone),
			wantErr: true,
		},
		{
			name: "#3: No query string params",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/required_agreements",
					Path:       "/api/v1/required_agreements/old",
					HTTPMethod: "GET",
					Headers: map[string]string{
						"Content-Type": "application/json",
						"api_key":      "test-3",
						"request_id":   "C713F659-1ABF-4EDF-8E86-655295CF7CD7",
					},
					PathParameters: map[string]string{
						"proxy": "example",
					},
					Body:            "",
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
			},
			want:    nil,
			want1:   int(v3.ResponseSuccess),
			wantErr: false,
		},
		{
			name: "#4: Valid Query string params [non european]",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/required_agreements",
					Path:       "/api/v1/required_agreements/old",
					HTTPMethod: "GET",
					Headers: map[string]string{
						"Content-Type": "application/json",
						"api_key":      "test-3",
						"request_id":   "C713F659-1ABF-4EDF-8E86-655295CF7CD7",
					},
					QueryStringParameters: map[string]string{
						"timezone": "America/Los_Angels",
					},
					PathParameters: map[string]string{
						"proxy": "example",
					},
					Body:            "",
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
			},
			want:    &QueryStringParams{Timezone: "America/Los_Angels"},
			want1:   int(v3.ResponseSuccess),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ValidateQueryStringParams(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateQueryStringParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateQueryStringParams() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ValidateQueryStringParams() got1 = %v, want1 %v", got1, tt.want1)
			}
		})
	}
}

func TestDetermineJurisdiction(t *testing.T) {
	type args struct {
		params  *QueryStringParams
		api_key string
		headers map[string]string
	}
	type test struct {
		name     string
		args     args
		expected string
		wantCode int
	}
	tests := []test{
		{
			name: "#1. European timezone",
			args: args{
				params: &QueryStringParams{
					Timezone: "Europe/London",
				},
				api_key: "",
				headers: map[string]string{
					"X-Forwarded-For": "",
				},
			},
			expected: "gdpr",
			wantCode: int(v3.ResponseSuccess),
		},
		{
			name: "#2. Empty Query params and Source IP is in Europe",
			args: args{
				params:  nil,
				api_key: "test-api-key",
				headers: map[string]string{
					"X-Forwarded-For": "127.0.0.1",
				},
			},
			expected: "gdpr",
			wantCode: int(v3.ResponseSuccess),
		},
		{
			name: "#3. Empty Query params and Ip not in europe",
			args: args{
				params:  nil,
				api_key: "test-api-key",
				headers: map[string]string{
					"X-Forwarded-For": "127.1.1.1",
				},
			},
			expected: "",
			wantCode: int(v3.ResponseInternalServerError),
		},
		{
			name: "#4. Non-European timezone and Invalid API key",
			args: args{
				params:  nil,
				api_key: "invalid_api_key",
				headers: map[string]string{
					"X-Forwarded-For": "0.0.0.0",
				},
			},
			expected: "",
			wantCode: int(v3.ResponseInternalServerError),
		},
		{
			name: "#5. Now European timezone",
			args: args{
				params: &QueryStringParams{
					Timezone: "America/Los_Angels",
				},
				api_key: "test-api-key",
				headers: map[string]string{
					"X-Forwarded-For": "512.0.0.40",
				},
			},
			expected: "",
			wantCode: int(v3.ResponseSuccess),
		},
	}

	httpClient = &ClientMock{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, gotCode, _ := DetermineJurisdiction(tt.args.params, tt.args.api_key, tt.args.headers)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("DetermineJurisdiction() returned %v, but expected %v", result, tt.expected)
			}
			if gotCode != tt.wantCode {
				t.Errorf("DetermineJurisdiction() got1 = %v, want1 %v", gotCode, tt.wantCode)
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

func TestCreateAgreement(t *testing.T) {
	type args struct {
		jurisdiction string
		agreement    string
		version      int
		locale       string
	}
	type test struct {
		name     string
		args     args
		expected Agreement
	}
	tests := []test{
		{
			name: "#1. CreateAgreement() with valid input",
			args: args{
				jurisdiction: "gdpr",
				agreement:    "tou",
				version:      1680166789,
				locale:       "en-uk",
			},
			expected: Agreement{
				Jurisdiction: "gdpr",
				Agreement:    "tou",
				Version:      1680166789,
				Locale:       "en-uk",
				URL:          SetConsentUrl(env.CloudfrontHostUrl, "gdpr", "en-uk", "tou", 1680166789),
			},
		},
		{
			name: "#2. CreateAgreement() with valid input",
			args: args{
				jurisdiction: "gdpr",
				agreement:    "pp",
				version:      1680166802,
				locale:       "en-uk",
			},
			expected: Agreement{
				Jurisdiction: "gdpr",
				Agreement:    "pp",
				Version:      1680166802,
				Locale:       "en-uk",
				URL:          SetConsentUrl(env.CloudfrontHostUrl, "gdpr", "en-uk", "pp", 1680166802),
			},
		},
		{
			name: "#3. CreateAgreement() with random input",
			args: args{
				jurisdiction: "eu",
				agreement:    "pp",
				version:      1680166902,
				locale:       "en-us",
			},
			expected: Agreement{
				Jurisdiction: "eu",
				Agreement:    "pp",
				Version:      1680166902,
				Locale:       "en-us",
				URL:          SetConsentUrl(env.CloudfrontHostUrl, "eu", "en-us", "pp", 1680166902),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateAgreement(tt.args.jurisdiction, tt.args.agreement, tt.args.locale, tt.args.version)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("CreateAgreement() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestValidateRequestPath(t *testing.T) {
	type args struct {
		request events.APIGatewayProxyRequest
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// API: /api/v1/required_agreements
		{
			name: "#1: When request path is valid for [ /api/v1/required_agreements ]",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/required_agreements",
					Path:       "/api/v1/required_agreements",
					HTTPMethod: "GET",
					Headers: map[string]string{
						"Content-Type": "application/json",
						"api_key":      "test-3",
						"request_id":   "C713F659-1ABF-4EDF-8E86-655295CF7CD7",
					},
					QueryStringParameters: map[string]string{
						"timezone": "Europe/London",
					},
					PathParameters: map[string]string{
						"proxy": "example",
					},
					Body:            "",
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
			},
			want:    int(v3.ResponseSuccess),
			wantErr: false,
		},
		{
			name: "#2: Invalid Request path",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/required",
					Path:       "/api/v1/required",
					HTTPMethod: "GET",
					Headers: map[string]string{
						"Content-Type": "application/json",
						"api_key":      "test-3",
						"request_id":   "C713F659-1ABF-4EDF-8E86-655295CF7CD7",
					},
					QueryStringParameters: map[string]string{
						"timezone": "Europe/London",
					},
					PathParameters: map[string]string{
						"proxy": "example",
					},
					Body:            "",
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
			},
			want:    int(v3.ResponseInvalidRequestPath),
			wantErr: true,
		},
		// API: /api/v1/required_agreements/{player_id}
		{
			name: "#3: When request path is valid for [ /api/v1/required_agreements/12001 ]",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/required_agreements/{player_id}",
					Path:       "/api/v1/required_agreements/12001",
					HTTPMethod: "GET",
					Headers: map[string]string{
						"Content-Type": "application/json",
						"api_key":      "test-3",
						"request_id":   "C713F659-1ABF-4EDF-8E86-655295CF7CD7",
					},
					QueryStringParameters: map[string]string{
						"timezone": "Europe/London",
					},
					PathParameters: map[string]string{
						"proxy": "example",
					},
					Body:            "",
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
			},
			want:    int(v3.ResponseSuccess),
			wantErr: false,
		},
		{
			name: "#4: When request path is invalid for [ /api/v1/required_agreements/12001 ]",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/required_agreements/12001/12321",
					Path:       "/api/v1/required_agreements/12002/123213",
					HTTPMethod: "GET",
					Headers: map[string]string{
						"Content-Type": "application/json",
						"api_key":      "test-3",
						"request_id":   "C713F659-1ABF-4EDF-8E86-655295CF7CD7",
					},
					QueryStringParameters: map[string]string{
						"timezone": "Europe/London",
					},
					PathParameters: map[string]string{
						"proxy": "example",
					},
					Body:            "",
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
			},
			want:    int(v3.ResponseInvalidRequestPath),
			wantErr: true,
		},

		// API: /api/v1/agree/{player_id}
		{
			name: "#5: When request path is valid for [ /api/v1/agree/12001 ]",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/agree/{player_id}",
					Path:       "/api/v1/agree/12001",
					HTTPMethod: "GET",
					Headers: map[string]string{
						"Content-Type": "application/json",
						"api_key":      "test-3",
						"request_id":   "C713F659-1ABF-4EDF-8E86-655295CF7CD7",
					},
					QueryStringParameters: map[string]string{
						"timezone": "Europe/London",
					},
					PathParameters: map[string]string{
						"proxy": "example",
					},
					Body:            "",
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
			},
			want:    int(v3.ResponseSuccess),
			wantErr: false,
		},
		{
			name: "#6: Invalid request path",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/agree/12001/21313",
					Path:       "/api/v1/agree/12002/1232",
					HTTPMethod: "GET",
					Headers: map[string]string{
						"Content-Type": "application/json",
						"api_key":      "test-3",
						"request_id":   "C713F659-1ABF-4EDF-8E86-655295CF7CD7",
					},
					QueryStringParameters: map[string]string{
						"timezone": "Europe/London",
					},
					PathParameters: map[string]string{
						"proxy": "example",
					},
					Body:            "",
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
			},
			want:    int(v3.ResponseInvalidRequestPath),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateRequestPath(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRequestPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateRequestPath() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateJurisdiction(t *testing.T) {
	tests := []struct {
		name    string
		args    AgreeRecord
		want    bool
		want1   int
		wantErr bool
	}{
		{
			name: "#1: Valid Jurisdiction [gdpr]",
			args: AgreeRecord{
				Jurisdiction: "gdpr",
			},
			want:    true,
			want1:   int(v3.ResponseSuccess),
			wantErr: false,
		},
		{
			name: "#2: Invalid Jurisdiction [abcs]",
			args: AgreeRecord{
				Jurisdiction: "abc",
			},
			want:    false,
			want1:   int(ResponseInvalidJurisdiction),
			wantErr: true,
		},
		{
			name: "#3: Empty Jurisdiction",
			args: AgreeRecord{
				Jurisdiction: "",
			},
			want:    false,
			want1:   int(ResponseInvalidJurisdiction),
			wantErr: true,
		},
		{
			name:    "#4: Jurisdiction key is missing",
			args:    AgreeRecord{},
			want:    false,
			want1:   int(ResponseInvalidJurisdiction),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ValidateJurisdiction(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJurisdiction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateJurisdiction() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ValidateJurisdiction() got1 = %v, want1 %v", got1, tt.want1)
			}
		})
	}
}

func TestValidateAgreement(t *testing.T) {
	tests := []struct {
		name    string
		args    AgreeRecord
		want    bool
		want1   int
		wantErr bool
	}{
		{
			name: "#1: When agreement is pp [pp]",
			args: AgreeRecord{
				Agreement: "pp",
			},
			want:    true,
			want1:   int(v3.ResponseSuccess),
			wantErr: false,
		},
		{
			name: "#2: When agreement is tou [tou]",
			args: AgreeRecord{
				Agreement: "tou",
			},
			want:    true,
			want1:   int(v3.ResponseSuccess),
			wantErr: false,
		},
		{
			name: "#3: When agreement is not-allowed list [abcd]",
			args: AgreeRecord{
				Agreement: "abcd",
			},
			want:    false,
			want1:   int(ResponseInvalidAgreement),
			wantErr: true,
		},
		{
			name: "#4: When agreement is empty",
			args: AgreeRecord{
				Agreement: "",
			},
			want:    false,
			want1:   int(ResponseInvalidAgreement),
			wantErr: true,
		},
		{
			name:    "#5: When agreement key is missing",
			args:    AgreeRecord{},
			want:    false,
			want1:   int(ResponseInvalidAgreement),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ValidateAgreement(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAgreement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateAgreement() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ValidateAgreement() got1 = %v, want1 %v", got1, tt.want1)
			}
		})
	}
}

func TestValidateVersion(t *testing.T) {
	tests := []struct {
		name    string
		args    AgreeRecord
		want    bool
		want1   int
		wantErr bool
	}{
		{
			name: "#1: When version is valid",
			args: AgreeRecord{
				Version: 1620031200, // May 3, 2021, 12:00:00 AM UTC
			},
			want:    true,
			want1:   int(v3.ResponseSuccess),
			wantErr: false,
		},
		{
			name: "#2: When version is zero",
			args: AgreeRecord{
				Version: 0,
			},
			want:    false,
			want1:   int(ResponseInvalidVersion),
			wantErr: true,
		},
		{
			name: "#3: When version is negative",
			args: AgreeRecord{
				Version: -1620031200, // November 9, 1919, 12:00:00 AM UTC
			},
			want:    false,
			want1:   int(ResponseInvalidVersion),
			wantErr: true,
		},
		{
			name:    "#4: When version is not provided",
			args:    AgreeRecord{},
			want:    false,
			want1:   int(ResponseInvalidVersion),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ValidateVersion(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ValidateVersion() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ValidateVersion() got1 = %v, want1 %v", got1, tt.want1)
			}
		})
	}
}

func TestValidateAgreementLocale(t *testing.T) {
	tests := []struct {
		name    string
		args    AgreeRecord
		want    bool
		want1   int
		wantErr bool
	}{
		{
			name: "#1: Valid agreement locale [en-uk]",
			args: AgreeRecord{
				Locale: "en-uk",
			},
			want:    true,
			want1:   int(v3.ResponseSuccess),
			wantErr: false,
		},
		{
			name: "#2: When locale is not in allowed list [en-ca]",
			args: AgreeRecord{
				Locale: "en-ca",
			},
			want:    false,
			want1:   int(ResponseInvalidLocale),
			wantErr: true,
		},

		{
			name: "#3: When locale is empty",
			args: AgreeRecord{
				Locale: "",
			},
			want:    false,
			want1:   int(ResponseInvalidLocale),
			wantErr: true,
		},
		{
			name:    "#4: When locale key is missing",
			args:    AgreeRecord{},
			want:    false,
			want1:   int(ResponseInvalidLocale),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ValidateAgreementLocale(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAgreementLocale() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateAgreementLocale() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ValidateAgreementLocale() got1 = %v, want1 %v", got1, tt.want1)
			}
		})
	}
}

func TestValidateTimestamp(t *testing.T) {
	type args struct {
		agreeRecord AgreeRecord
		currentTime time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		want1   int
		wantErr bool
	}{
		{
			name: "#1: Valid timestamp [2023-04-20T12:34:01Z]",
			args: args{
				agreeRecord: AgreeRecord{
					Timestamp: "2023-04-20T12:34:01Z",
				},
				currentTime: time.Date(2023, time.April, 20, 12, 34, 20, 0, time.UTC),
			},
			want:    true,
			want1:   int(v3.ResponseSuccess),
			wantErr: false,
		},
		{
			name: "#2: When timestamp is not valid",
			args: args{
				agreeRecord: AgreeRecord{
					Timestamp: "2023-04-21T12:07:48+0000",
				},
				currentTime: time.Date(2023, time.April, 20, 12, 34, 56, 0, time.UTC),
			},
			want:    false,
			want1:   int(ResponseInvalidTimestamp),
			wantErr: true,
		},

		{
			name: "#3: When timestamp is empty",
			args: args{
				agreeRecord: AgreeRecord{
					Timestamp: "",
				},
				currentTime: time.Date(2023, time.April, 20, 05, 50, 56, 0, time.UTC),
			},
			want:    false,
			want1:   int(ResponseInvalidTimestamp),
			wantErr: true,
		},
		{
			name: "#4: When timestamp key is missing",
			args: args{
				agreeRecord: AgreeRecord{},
				currentTime: time.Date(2023, time.April, 20, 12, 34, 56, 0, time.UTC),
			},
			want:    false,
			want1:   int(ResponseInvalidTimestamp),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ValidateTimestamp(tt.args.agreeRecord, tt.args.currentTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTimestamp() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ValidateTimestamp() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ValidateTimestamp() got1 = %v, want1 %v", got1, tt.want1)
			}
		})
	}
}

func TestIsTimeWithinThreshold(t *testing.T) {

	type args struct {
		currentTime    time.Time
		givenTimestamp string
		timeRange      int
	}
	tests := []struct {
		name      string
		args      args
		want      bool
		wantError bool
	}{
		{
			name: "#1: When timestamp is valid [2023-04-20T12:34:01Z]",
			args: args{
				currentTime:    time.Date(2023, time.April, 20, 12, 34, 20, 0, time.UTC),
				givenTimestamp: "2023-04-20T12:34:01Z",
				timeRange:      30,
			},
			want:      true,
			wantError: false,
		},
		{
			name: "#2: When timestamp is not valid [2023-04-20T12:00:01Z]",
			args: args{
				currentTime:    time.Date(2023, time.April, 20, 12, 34, 20, 0, time.UTC),
				givenTimestamp: "2023-04-20T12:00:01Z",
				timeRange:      30,
			},
			want:      false,
			wantError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsTimeWithinThreshold(tt.args.currentTime, tt.args.givenTimestamp, tt.args.timeRange)
			if (err != nil) != tt.wantError {
				t.Errorf("IsTimeWithinThreshold() gotError = %v, wantError %v", err, tt.wantError)
			}
			if got != tt.want {
				t.Errorf("IsTimeWithinThreshold() got = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestAgreementTypeIsAvailable(t *testing.T) {
	type args struct {
		ValidAgreementList []string
		Agreement          string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "#1: Agreement is available in the list",
			args: args{
				ValidAgreementList: []string{"pp", "tou"},
				Agreement:          "pp",
			},
			want: true,
		},
		{
			name: "#2: Agreement is available in the list",
			args: args{
				ValidAgreementList: []string{"pp", "tou"},
				Agreement:          "tou",
			},
			want: true,
		},
		{
			name: "#3: Agreement is not available in the list",
			args: args{
				ValidAgreementList: []string{"pp", "tou"},
				Agreement:          "abc",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AgreementTypeIsAvailable(tt.args.Agreement, tt.args.ValidAgreementList)
			if got != tt.want {
				t.Errorf("AgreementTypeIsAvailable() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateBodyParams(t *testing.T) {

	type args struct {
		requestBody         string
		mandatoryParameters []string
		currentTime         time.Time
	}

	tests := []struct {
		name     string
		args     args
		want     []AgreeRecord
		wantCode int
		wantErr  bool
	}{
		{
			name: "#1: Valid Request body all parameters",
			args: args{
				requestBody: `[
					{
						"jurisdiction": "gdpr",
						"agreement": "pp",
						"version": 3242343243,
						"locale": "en-uk",
						"timestamp": "2023-04-20T12:34:01Z"
					},
					{
						"jurisdiction": "gdpr",
						"agreement": "tou",
						"version": 3242343244,
						"locale": "en-uk",
						"timestamp": "2023-04-20T12:34:01Z"
					}
				]`,
				mandatoryParameters: []string{
					"jurisdiction",
					"agreement",
					"version",
					"locale",
					"timestamp",
				},
				currentTime: time.Date(2023, time.April, 20, 12, 34, 20, 0, time.UTC),
			},
			want: []AgreeRecord{
				AgreeRecord{
					Agreement:    "pp",
					Jurisdiction: "gdpr",
					Locale:       "en-uk",
					Timestamp:    "2023-04-20T12:34:01Z",
					Version:      3242343243,
				},
				AgreeRecord{
					Agreement:    "tou",
					Jurisdiction: "gdpr",
					Locale:       "en-uk",
					Timestamp:    "2023-04-20T12:34:01Z",
					Version:      3242343244,
				},
			},
			wantCode: int(v3.ResponseSuccess),
			wantErr:  false,
		},
		{
			name: "#2: Invalid jurisdiction [abcd]",
			args: args{
				requestBody: `[
					{
						"jurisdiction": "abcd",
						"agreement": "pp",
						"version": 3242343243,
						"locale": "en-uk",
						"timestamp": "2023-04-20T12:34:01Z"
					}
				]`,
				mandatoryParameters: []string{
					"jurisdiction",
					"agreement",
					"version",
					"locale",
					"timestamp",
				},
				currentTime: time.Date(2023, time.April, 20, 12, 34, 20, 0, time.UTC),
			},
			want:     nil,
			wantCode: int(ResponseInvalidJurisdiction),
			wantErr:  true,
		},
		{
			name: "#3: Empty jurisdiction",
			args: args{
				requestBody: `[
					{
						"jurisdiction": "",
						"agreement": "pp",
						"version": 3242343243,
						"locale": "en-uk",
						"timestamp": "2023-04-20T12:34:01Z"
					}
				]`,
				mandatoryParameters: []string{
					"jurisdiction",
					"agreement",
					"version",
					"locale",
					"timestamp",
				},
				currentTime: time.Date(2023, time.April, 20, 12, 34, 20, 0, time.UTC),
			},
			want:     nil,
			wantCode: int(ResponseInvalidJurisdiction),
			wantErr:  true,
		},
		{
			name: "#4: Missing jurisdiction",
			args: args{
				requestBody: `[
					{
						"agreement": "pp",
						"version": 3242343243,
						"locale": "en-uk",
						"timestamp": "2023-04-20T12:34:01Z"
					}
				]`,
				mandatoryParameters: []string{
					"jurisdiction",
					"agreement",
					"version",
					"locale",
					"timestamp",
				},
				currentTime: time.Date(2023, time.April, 20, 12, 34, 20, 0, time.UTC),
			},
			want:     nil,
			wantCode: int(ResponseInvalidJurisdiction),
			wantErr:  true,
		},
		{
			name: "#5: Invalid agreement type [abcd]",
			args: args{
				requestBody: `[
					{
						"jurisdiction": "gdpr",
						"agreement": "abcd",
						"version": 3242343243,
						"locale": "en-uk",
						"timestamp": "2023-04-20T12:34:01Z"
					}
				]`,
				mandatoryParameters: []string{
					"jurisdiction",
					"agreement",
					"version",
					"locale",
					"timestamp",
				},
				currentTime: time.Date(2023, time.April, 20, 12, 34, 20, 0, time.UTC),
			},
			want:     nil,
			wantCode: int(ResponseInvalidAgreement),
			wantErr:  true,
		},

		{
			name: "#6: Empty agreement type",
			args: args{
				requestBody: `[
					{
						"jurisdiction": "gdpr",
						"agreement": "",
						"version": 3242343243,
						"locale": "en-uk",
						"timestamp": "2023-04-20T12:34:01Z"
					}
				]`,
				mandatoryParameters: []string{
					"jurisdiction",
					"agreement",
					"version",
					"locale",
					"timestamp",
				},
				currentTime: time.Date(2023, time.April, 20, 12, 34, 20, 0, time.UTC),
			},
			want:     nil,
			wantCode: int(ResponseInvalidAgreement),
			wantErr:  true,
		},

		{
			name: "#7: Agreement type is missing",
			args: args{
				requestBody: `[
					{
						"jurisdiction": "gdpr",
						"version": 3242343243,
						"locale": "en-uk",
						"timestamp": "2023-04-20T12:34:01Z"
					}
				]`,
				mandatoryParameters: []string{
					"jurisdiction",
					"agreement",
					"version",
					"locale",
					"timestamp",
				},
				currentTime: time.Date(2023, time.April, 20, 12, 34, 20, 0, time.UTC),
			},
			want:     nil,
			wantCode: int(ResponseInvalidAgreement),
			wantErr:  true,
		},
		{
			name: "#8: Request body params is invalid when version is zero",
			args: args{
				requestBody: `[
					{
						"jurisdiction": "gdpr",
						"agreement": "pp",
						"version": 0,
						"locale": "en-uk",
						"timestamp": "2023-04-20T12:34:01Z"
					}
				]`,
				mandatoryParameters: []string{
					"jurisdiction",
					"agreement",
					"version",
					"locale",
					"timestamp",
				},
				currentTime: time.Date(2023, time.April, 20, 12, 34, 20, 0, time.UTC),
			},
			want:     nil,
			wantCode: int(ResponseInvalidVersion),
			wantErr:  true,
		},
		{
			name: "#9: Request body params is invalid when version is missing",
			args: args{
				requestBody: `[
					{
						"jurisdiction": "gdpr",
						"agreement": "pp",
						"locale": "en-uk",
						"timestamp": "2023-04-20T12:34:01Z"
					}
				]`,
				mandatoryParameters: []string{
					"jurisdiction",
					"agreement",
					"version",
					"locale",
					"timestamp",
				},
				currentTime: time.Date(2023, time.April, 20, 12, 34, 20, 0, time.UTC),
			},
			want:     nil,
			wantCode: int(ResponseInvalidVersion),
			wantErr:  true,
		},
		{
			name: "#10: Request body params is invalid for locale is en-ca [en-ca]",
			args: args{
				requestBody: `[
					{
						"jurisdiction": "gdpr",
						"agreement": "pp",
						"version": 3242343243,
						"locale": "en-ca",
						"timestamp": "2023-04-20T12:34:01Z"
					}
				]`,
				mandatoryParameters: []string{
					"jurisdiction",
					"agreement",
					"version",
					"locale",
					"timestamp",
				},
				currentTime: time.Date(2023, time.April, 20, 12, 34, 20, 0, time.UTC),
			},
			want:     nil,
			wantCode: int(ResponseInvalidLocale),
			wantErr:  true,
		},

		{
			name: "#11: Request body params is invalid for locale is empty",
			args: args{
				requestBody: `[
					{
						"jurisdiction": "gdpr",
						"agreement": "pp",
						"version": 3242343243,
						"locale": "",
						"timestamp": "2023-04-20T12:34:01Z"
					}
				]`,
				mandatoryParameters: []string{
					"jurisdiction",
					"agreement",
					"version",
					"locale",
					"timestamp",
				},
				currentTime: time.Date(2023, time.April, 20, 12, 34, 20, 0, time.UTC),
			},
			want:     nil,
			wantCode: int(ResponseInvalidLocale),
			wantErr:  true,
		},
		{
			name: "#12: Request body params is invalid for locale is missing",
			args: args{
				requestBody: `[
					{
						"jurisdiction": "gdpr",
						"agreement": "tou",
						"version": 3242343243,
						"timestamp": "2023-04-20T12:34:01Z"
					}
				]`,
				mandatoryParameters: []string{
					"jurisdiction",
					"agreement",
					"version",
					"locale",
					"timestamp",
				},
				currentTime: time.Date(2023, time.April, 20, 12, 34, 20, 0, time.UTC),
			},
			want:     nil,
			wantCode: int(ResponseInvalidLocale),
			wantErr:  true,
		},
		{
			name: "#13: Invalid timestamp",
			args: args{
				requestBody: `[
					{
						"jurisdiction": "gdpr",
						"agreement": "pp",
						"version": 3242343243,
						"locale": "en-uk",
						"timestamp": "Fri, 21 Apr 2023 13:37:24 +0000"
					}
				]`,
				mandatoryParameters: []string{
					"jurisdiction",
					"agreement",
					"version",
					"locale",
					"timestamp",
				},
				currentTime: time.Date(2023, time.April, 20, 12, 35, 20, 0, time.UTC),
			},
			want:     nil,
			wantCode: int(ResponseInvalidTimestamp),
			wantErr:  true,
		},
		{
			name: "#14: Request body is invalid for timestamp is empty",
			args: args{
				requestBody: `[
					{
						"jurisdiction": "gdpr",
						"agreement": "pp",
						"version": 3242343243,
						"locale": "en-uk",
						"timestamp": ""
					}
				]`,
				mandatoryParameters: []string{
					"jurisdiction",
					"agreement",
					"version",
					"locale",
					"timestamp",
				},
				currentTime: time.Date(2023, time.April, 20, 12, 34, 20, 0, time.UTC),
			},
			want:     nil,
			wantCode: int(ResponseInvalidTimestamp),
			wantErr:  true,
		},
		{
			name: "#15: Request body params is invalid for timestamp is missing",
			args: args{
				requestBody: `[
					{
						"jurisdiction": "gdpr",
						"agreement": "pp",
						"version": 3242343243,
						"locale": "en-uk"
					}
				]`,
				mandatoryParameters: []string{
					"jurisdiction",
					"agreement",
					"version",
					"locale",
					"timestamp",
				},
				currentTime: time.Date(2023, time.April, 20, 12, 34, 20, 0, time.UTC),
			},
			want:     nil,
			wantCode: int(ResponseInvalidTimestamp),
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotCode, err := ValidateBodyParams(tt.args.requestBody, tt.args.mandatoryParameters, tt.args.currentTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateBodyParams() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateBodyParams() got = %v, want %v", got, tt.want)
			}
			if gotCode != tt.wantCode {
				t.Errorf("ValidateBodyParams() gotCode = %v, wantCode %v", gotCode, tt.wantCode)
			}
		})
	}
}

func TestGetUniqueAgreementsForJurisdiction(t *testing.T) {
	type args struct {
		jurisdiction string
	}
	tests := []struct {
		name    string
		args    args
		want    []DynamoJurisdictionMapping
		wantErr bool
	}{
		{
			name: "#1: Data available in Dynamo Table [jurisdiction is gdpr]",
			args: args{
				jurisdiction: "gdpr",
			},
			want: []DynamoJurisdictionMapping{
				{
					Jurisdiction:    "gdpr",
					Agreement:       "tou",
					AgreementLocale: "en-uk",
					Version:         1681118745,
				},
				{
					Jurisdiction:    "gdpr",
					Agreement:       "pp",
					AgreementLocale: "en-uk",
					Version:         1681118791,
				},
			},
			wantErr: false,
		},
		{
			name: "#2: Data available in Dynamo Table [jurisdiction is eu]",
			args: args{
				jurisdiction: "eu",
			},
			want: []DynamoJurisdictionMapping{
				{
					Jurisdiction:    "eu",
					Agreement:       "pp",
					AgreementLocale: "en-uk",
					Version:         1680845499,
				},
				{
					Jurisdiction:    "eu",
					Agreement:       "tou",
					AgreementLocale: "en-uk",
					Version:         1680846146,
				},
			},
			wantErr: false,
		},
		{
			name: "#3: Data not found in Dynamo table [Jurisdiction is pipl]",
			args: args{
				jurisdiction: "pipl",
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "#4: Generate an error in dynamo query",
			args: args{
				jurisdiction: "test",
			},
			want:    nil,
			wantErr: true,
		},
	}

	setEnvVars(map[TypeEnvVar]string{
		ENV_DYNAMO_MAPPING_TABLE_NAME: "jurisdiction_mapping_test",
	})

	env = NewEnv()

	dynamoSvc = &MockDynamoClient{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUniqueAgreementsForJurisdiction(dynamoSvc, tt.args.jurisdiction)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUniqueAgreementsForJurisdiction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUniqueAgreementsForJurisdiction() got= %v, want %v", got, tt.want)
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

func currentTimeISO8601() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05Z")
}

func TestHandleRequest(t *testing.T) {
	parameterStore = &im_config.MockKeyClient{
		FileBytes: []byte(`{"keys":[{"key":"/honeybadger/admin/api_key","value":{"Name":"/honeybadger/admin/api_key","Value":"testing"}}, {"key":"/ip_quality_score/private_key","value":{"Name":"/ip_quality_score/private_key","Value":"test-api-key"}}]}`),
	}

	type args struct {
		ctx     context.Context
		request events.APIGatewayProxyRequest
	}

	setEnvVars(map[TypeEnvVar]string{
		ENV_APPLICATION_NAME:           "consent_management",
		ENV_LOG_LEVEL:                  "debug",
		ENV_AWS_REGION:                 "us-east-1",
		ENV_PARAMETERS_STORE:           "/tmp/parameters.txt",
		ENV_DYNAMO_CONSENT_TABLE_NAME:  "consent_management_test",
		ENV_DYNAMO_API_KEYS_TABLE_NAME: "api_keys_test",
		ENV_DYNAMO_MAPPING_TABLE_NAME:  "jurisdiction_mapping_test",
		ENV_CLOUDFRONT_HOST_URL:        "https://test.cloudfront.net",
		ENV_ALLOWED_AGREEMENT_LIST:     "pp,tou",
		ENV_ALLOWED_JURISDICTION_LIST:  "gdpr",
		ENV_ALLOWED_LOCALE_LIST:        "en-uk",
	})

	tests := []struct {
		name    string
		args    args
		want    events.APIGatewayProxyResponse
		wantErr bool
	}{
		// API: /api/v1/required_agreements
		{
			name: "#1: Success API request",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/required_agreements",
					Path:       "/api/v1/required_agreements",
					HTTPMethod: "GET",
					Headers: map[string]string{
						"Content-Type":   "application/json",
						"api_key":        "test-api-key",
						"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
						"correlation_id": "test-correlation-id",
					},
					QueryStringParameters: map[string]string{
						"timezone": "Europe/London",
					},
					PathParameters: map[string]string{
						"proxy": "example",
					},
					Body:            "",
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
				ctx: context.TODO(),
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"correlation_id": "test-correlation-id",
					"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
				},
				Body:            `{"error":0,"data":[{"jurisdiction":"gdpr","agreement":"tou","version":1681118745,"locale":"en-uk","url":"https://test.cloudfront.net/agreement/gdpr/en-uk/tou/1681118745/tou.html"},{"jurisdiction":"gdpr","agreement":"pp","version":1681118791,"locale":"en-uk","url":"https://test.cloudfront.net/agreement/gdpr/en-uk/pp/1681118791/pp.html"}]}`,
				IsBase64Encoded: false,
			},
			wantErr: false,
		},
		{
			name: "#2: Missing API Key from headers",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/required_agreements",
					Path:       "/api/v1/required_agreements",
					HTTPMethod: "GET",
					Headers: map[string]string{
						"Content-Type":   "application/json",
						"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
						"correlation_id": "correlation_id-1",
					},
					QueryStringParameters: map[string]string{
						"timezone": "Europe/London",
					},
					PathParameters: map[string]string{
						"proxy": "example",
					},
					Body:            "",
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
				ctx: context.TODO(),
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"correlation_id": "correlation_id-1",
					"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
				},
				Body:            fmt.Sprintf(`{"error":%d,"data":null}`, int(v3.ResponseNotAuthenticated)),
				IsBase64Encoded: false,
			},
			wantErr: false,
		},
		{
			name: "#3: Empty api_key",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/required_agreements",
					Path:       "/api/v1/required_agreements",
					HTTPMethod: "GET",
					Headers: map[string]string{
						"Content-Type":   "application/json",
						"api_key":        "",
						"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
						"correlation_id": "correlation_id-1",
					},
					QueryStringParameters: map[string]string{
						"timezone": "Europe/London",
					},
					PathParameters: map[string]string{
						"proxy": "example",
					},
					Body:            "",
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
				ctx: context.TODO(),
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"correlation_id": "correlation_id-1",
					"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
				},
				Body:            fmt.Sprintf(`{"error":%d,"data":null}`, int(v3.ResponseNotAuthenticated)),
				IsBase64Encoded: false,
			},
			wantErr: false,
		},
		{
			name: "#4: Invalid request_id in headers",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/required_agreements",
					Path:       "/api/v1/required_agreements",
					HTTPMethod: "GET",
					Headers: map[string]string{
						"Content-Type":   "application/json",
						"api_key":        "test-3",
						"request_id":     "",
						"correlation_id": "correlation_id-1",
					},
					QueryStringParameters: map[string]string{
						"timezone": "Europe/London",
					},
					PathParameters: map[string]string{
						"proxy": "example",
					},
					Body:            "",
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
				ctx: context.TODO(),
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"correlation_id": "correlation_id-1",
					"request_id":     "",
				},
				Body:            fmt.Sprintf(`{"error":%d,"data":null}`, int(v3.ResponseInvalidRequestIdHeader)),
				IsBase64Encoded: false,
			},
			wantErr: false,
		},
		{
			name: "#5: Missing correlation_id",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/required_agreements",
					Path:       "/api/v1/required_agreements",
					HTTPMethod: "GET",
					Headers: map[string]string{
						"Content-Type": "application/json",
						"api_key":      "test-api-key-1",
						"request_id":   "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
						// "correlation_id": "correlation_id-1",
					},
					QueryStringParameters: map[string]string{
						"timezone": "Europe/London",
					},
					PathParameters: map[string]string{
						"proxy": "example",
					},
					Body:            "",
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
				ctx: context.TODO(),
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"request_id": "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
				},
				Body:            fmt.Sprintf(`{"error":%d,"data":null}`, int(v3.ResponseMissingCorrelationId)),
				IsBase64Encoded: false,
			},
			wantErr: false,
		},
		{
			name: "#6: Success API request when the timezone is within Europe [timezone: Europe/Berlin]",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/required_agreements",
					Path:       "/api/v1/required_agreements",
					HTTPMethod: "GET",
					Headers: map[string]string{
						"Content-Type":   "application/json",
						"api_key":        "test-3",
						"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
						"correlation_id": "correlation_id-1",
					},
					QueryStringParameters: map[string]string{
						"timezone": "Europe/London",
					},
					PathParameters: map[string]string{
						"proxy": "example",
					},
					Body:            "",
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
				ctx: context.TODO(),
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"correlation_id": "correlation_id-1",
					"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
				},
				Body:            `{"error":0,"data":[{"jurisdiction":"gdpr","agreement":"tou","version":1681118745,"locale":"en-uk","url":"https://test.cloudfront.net/agreement/gdpr/en-uk/tou/1681118745/tou.html"},{"jurisdiction":"gdpr","agreement":"pp","version":1681118791,"locale":"en-uk","url":"https://test.cloudfront.net/agreement/gdpr/en-uk/pp/1681118791/pp.html"}]}`,
				IsBase64Encoded: false,
			},
			wantErr: false,
		},
		{
			name: "#7: Empty agreement list while the timezone is outside of Europe [timezone: Asia/Kabul]",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/required_agreements",
					Path:       "/api/v1/required_agreements",
					HTTPMethod: "GET",
					Headers: map[string]string{
						"Content-Type":   "application/json",
						"api_key":        "test-3",
						"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
						"correlation_id": "correlation_id-1",
					},
					QueryStringParameters: map[string]string{
						"timezone": "Asia/Kabul",
					},
					PathParameters: map[string]string{
						"proxy": "example",
					},
					Body:            "",
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
				ctx: context.TODO(),
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"correlation_id": "correlation_id-1",
					"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
				},
				Body:            `{"error":0,"data":[]}`,
				IsBase64Encoded: false,
			},
			wantErr: false,
		},
		{
			name: "#8: Invalid timezone",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/required_agreements",
					Path:       "/api/v1/required_agreements",
					HTTPMethod: "GET",
					Headers: map[string]string{
						"Content-Type":   "application/json",
						"api_key":        "test-3",
						"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
						"correlation_id": "correlation_id-1",
					},
					QueryStringParameters: map[string]string{
						"timezone": "",
					},
					PathParameters: map[string]string{
						"proxy": "example",
					},
					Body:            "",
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
				ctx: context.TODO(),
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"correlation_id": "correlation_id-1",
					"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
				},
				Body:            fmt.Sprintf(`{"error":%d,"data":null}`, int(ResponseInvalidTimezone)),
				IsBase64Encoded: false,
			},
			wantErr: false,
		},
		{
			name: "#9: Empty Query Params and Ip address is in Europe",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/required_agreements",
					Path:       "/api/v1/required_agreements",
					HTTPMethod: "GET",
					Headers: map[string]string{
						"Content-Type":    "application/json",
						"api_key":         "test-3",
						"request_id":      "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
						"correlation_id":  "correlation_id-1",
						"X-Forwarded-For": "127.0.0.1, 37.111.194.51, 130.176.187.12",
					},
					QueryStringParameters: nil,
					PathParameters: map[string]string{
						"proxy": "example",
					},
					Body:            "",
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
				ctx: context.TODO(),
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"correlation_id": "correlation_id-1",
					"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
				},
				Body:            fmt.Sprintf(`{"error":%d,"data":[{"jurisdiction":"gdpr","agreement":"tou","version":1681118745,"locale":"en-uk","url":"https://test.cloudfront.net/agreement/gdpr/en-uk/tou/1681118745/tou.html"},{"jurisdiction":"gdpr","agreement":"pp","version":1681118791,"locale":"en-uk","url":"https://test.cloudfront.net/agreement/gdpr/en-uk/pp/1681118791/pp.html"}]}`, int(v3.ResponseSuccess)),
				IsBase64Encoded: false,
			},
			wantErr: false,
		},
		{
			name: "#10: Invalid Request Method",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/required_agreements",
					Path:       "/api/v1/required_agreements",
					HTTPMethod: "PUT",
					Headers: map[string]string{
						"Content-Type":   "application/json",
						"api_key":        "test-3",
						"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
						"correlation_id": "correlation_id-1",
					},
					QueryStringParameters: map[string]string{
						"timezone": "Europe/London",
					},
					PathParameters: map[string]string{
						"proxy": "example",
					},
					Body:            "",
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
				ctx: context.TODO(),
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"correlation_id": "correlation_id-1",
					"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
				},
				Body:            fmt.Sprintf(`{"error":%d,"data":null}`, int(v3.ResponseInvalidRequestMethod)),
				IsBase64Encoded: false,
			},
			wantErr: false,
		},
		{
			name: "#11: Invalid Request Path",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/required",
					Path:       "/api/v1/required",
					HTTPMethod: "GET",
					Headers: map[string]string{
						"Content-Type":   "application/json",
						"api_key":        "test-3",
						"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
						"correlation_id": "correlation_id-1",
					},
					QueryStringParameters: map[string]string{
						"timezone": "Europe/London",
					},
					PathParameters: map[string]string{
						"proxy": "example",
					},
					Body:            "",
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
				ctx: context.TODO(),
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
					"correlation_id": "correlation_id-1",
				},
				Body:            fmt.Sprintf(`{"error":%d,"data":null}`, int(v3.ResponseInvalidRequestPath)),
				IsBase64Encoded: false,
			},
			wantErr: false,
		},
		{
			name: "#12: GEO IP address Timezone within Europe [127.0.0.10]",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/required_agreements",
					Path:       "/api/v1/required_agreements",
					HTTPMethod: "GET",
					Headers: map[string]string{
						"Content-Type":    "application/json",
						"api_key":         "test-api-key",
						"request_id":      "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
						"correlation_id":  "test-correlation-id",
						"X-Forwarded-For": "127.0.0.1, 37.111.194.51, 130.176.187.12",
					},
					QueryStringParameters: nil,
					PathParameters: map[string]string{
						"proxy": "example",
					},
					Body:            "",
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
				ctx: context.TODO(),
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"correlation_id": "test-correlation-id",
					"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
				},
				Body:            `{"error":0,"data":[{"jurisdiction":"gdpr","agreement":"tou","version":1681118745,"locale":"en-uk","url":"https://test.cloudfront.net/agreement/gdpr/en-uk/tou/1681118745/tou.html"},{"jurisdiction":"gdpr","agreement":"pp","version":1681118791,"locale":"en-uk","url":"https://test.cloudfront.net/agreement/gdpr/en-uk/pp/1681118791/pp.html"}]}`,
				IsBase64Encoded: false,
			},
			wantErr: false,
		},
		{
			name: "#13: GEO IP address Timezone not within Europe [127.0.0.50]",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/required_agreements",
					Path:       "/api/v1/required_agreements",
					HTTPMethod: "GET",
					Headers: map[string]string{
						"Content-Type":    "application/json",
						"api_key":         "test-api-key",
						"request_id":      "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
						"correlation_id":  "test-correlation-id",
						"X-Forwarded-For": "127.0.0.50",
					},
					QueryStringParameters: nil,
					PathParameters: map[string]string{
						"proxy": "example",
					},
					Body:            "",
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
				ctx: context.TODO(),
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"correlation_id": "test-correlation-id",
					"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
				},
				Body:            fmt.Sprintf(`{"error":%d,"data":null}`, int(v3.ResponseInternalServerError)),
				IsBase64Encoded: false,
			},
			wantErr: false,
		},

		// API: /api/v1/agree/{player_id}
		{
			name: "#14: Agree API request success",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/agree/{player_id}",
					Path:       "/api/v1/agree/12001",
					HTTPMethod: "POST",
					Headers: map[string]string{
						"Content-Type": "application/json",
						"api_key":      "test-api-key",
						"request_id":   "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
					},
					QueryStringParameters: map[string]string{},
					PathParameters: map[string]string{
						"player_id": "12001",
					},
					Body: fmt.Sprintf(`[
								{
									"jurisdiction": "gdpr",
									"agreement": "pp",
									"version": 3242343243,
									"locale": "en-uk",
									"timestamp": "%s"
								},
								{
									"jurisdiction": "gdpr",
									"agreement": "tou",
									"version": 214214214,
									"locale": "en-uk",
									"timestamp": "%s"
								}
							]`, currentTimeISO8601(), currentTimeISO8601()),
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
				ctx: context.TODO(),
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"request_id": "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
				},
				Body:            fmt.Sprintf(`{"error":%d,"data":null}`, int(v3.ResponseSuccess)),
				IsBase64Encoded: false,
			},
			wantErr: false,
		},
		{
			name: "#15: Provide and response back correlation_id",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/agree/{player_id}",
					Path:       "/api/v1/agree/12001",
					HTTPMethod: "POST",
					Headers: map[string]string{
						"Content-Type":   "application/json",
						"api_key":        "test-api-key",
						"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
						"correlation_id": "323232323232323232323",
					},
					QueryStringParameters: map[string]string{},
					PathParameters: map[string]string{
						"player_id": "12001",
					},
					Body: fmt.Sprintf(`[
								{
									"jurisdiction": "gdpr",
									"agreement": "pp",
									"version": 3242343243,
									"locale": "en-uk",
									"timestamp": "%s"
								},
								{
									"jurisdiction": "gdpr",
									"agreement": "tou",
									"version": 214214214,
									"locale": "en-uk",
									"timestamp": "%s"
								}
							]`, currentTimeISO8601(), currentTimeISO8601()),
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
				ctx: context.TODO(),
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
					"correlation_id": "323232323232323232323",
				},
				Body:            fmt.Sprintf(`{"error":%d,"data":null}`, int(v3.ResponseSuccess)),
				IsBase64Encoded: false,
			},
			wantErr: false,
		},
		{
			name: "#16: Invalid Player Id",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/agree/{player_id}",
					Path:       "/api/v1/agree/FGK12001",
					HTTPMethod: "POST",
					Headers: map[string]string{
						"Content-Type":   "application/json",
						"api_key":        "test-api-key",
						"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
						"correlation_id": "323232323232323232323",
					},
					QueryStringParameters: map[string]string{},
					PathParameters: map[string]string{
						"player_id": "FGK12001",
					},
					Body: fmt.Sprintf(`[
								{
									"jurisdiction": "gdpr",
									"agreement": "pp",
									"version": 3242343243,
									"locale": "en-uk",
									"timestamp": "%s"
								},
								{
									"jurisdiction": "gdpr",
									"agreement": "tou",
									"version": 214214214,
									"locale": "en-uk",
									"timestamp": "%s"
								}
							]`, currentTimeISO8601(), currentTimeISO8601()),
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
				ctx: context.TODO(),
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
					"correlation_id": "323232323232323232323",
				},
				Body:            fmt.Sprintf(`{"error":%d,"data":null}`, int(v3.ResponseInvalidMember)),
				IsBase64Encoded: false,
			},
			wantErr: false,
		},
		{
			name: "#17: Invalid Jurisdiction",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/agree/{player_id}",
					Path:       "/api/v1/agree/12001",
					HTTPMethod: "POST",
					Headers: map[string]string{
						"Content-Type":   "application/json",
						"api_key":        "test-api-key",
						"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
						"correlation_id": "323232323232323232323",
					},
					QueryStringParameters: map[string]string{},
					PathParameters: map[string]string{
						"player_id": "12001",
					},
					Body: fmt.Sprintf(`[
								{
									"jurisdiction": "non-gdpr",
									"agreement": "pp",
									"version": 3242343243,
									"locale": "en-uk",
									"timestamp": "%s"
								},
								{
									"jurisdiction": "non-gdpr",
									"agreement": "tou",
									"version": 214214214,
									"locale": "en-uk",
									"timestamp": "%s"
								}
							]`, currentTimeISO8601(), currentTimeISO8601()),
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
				ctx: context.TODO(),
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
					"correlation_id": "323232323232323232323",
				},
				Body:            fmt.Sprintf(`{"error":%d,"data":null}`, int(ResponseInvalidJurisdiction)),
				IsBase64Encoded: false,
			},
			wantErr: false,
		},
		{
			name: "#18: Missing Jurisdiction",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/agree/{player_id}",
					Path:       "/api/v1/agree/12001",
					HTTPMethod: "POST",
					Headers: map[string]string{
						"Content-Type":   "application/json",
						"api_key":        "test-api-key",
						"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
						"correlation_id": "323232323232323232323",
					},
					QueryStringParameters: map[string]string{},
					PathParameters: map[string]string{
						"player_id": "12001",
					},
					Body: fmt.Sprintf(`[
								{
									
									"agreement": "pp",
									"version": 3242343243,
									"locale": "en-uk",
									"timestamp": "%s"
								},
								{
									"jurisdiction": "gdpr",
									"agreement": "tou",
									"version": 214214214,
									"locale": "en-uk",
									"timestamp": "%s"
								}
							]`, currentTimeISO8601(), currentTimeISO8601()),
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
				ctx: context.TODO(),
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
					"correlation_id": "323232323232323232323",
				},
				Body:            fmt.Sprintf(`{"error":%d,"data":null}`, int(ResponseInvalidJurisdiction)),
				IsBase64Encoded: false,
			},
			wantErr: false,
		},
		{
			name: "#19: Invalid Agreement",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/agree/{player_id}",
					Path:       "/api/v1/agree/12001",
					HTTPMethod: "POST",
					Headers: map[string]string{
						"Content-Type":   "application/json",
						"api_key":        "test-api-key",
						"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
						"correlation_id": "323232323232323232323",
					},
					QueryStringParameters: map[string]string{},
					PathParameters: map[string]string{
						"player_id": "12001",
					},
					Body: fmt.Sprintf(`[
								{
									"jurisdiction": "gdpr",
									"agreement": "invalid_ppr",
									"version": 3242343243,
									"locale": "en-uk",
									"timestamp": "%s"
								},
								{
									"jurisdiction": "gdpr",
									"agreement": "tou",
									"version": 214214214,
									"locale": "en-uk",
									"timestamp": "%s"
								}
							]`, currentTimeISO8601(), currentTimeISO8601()),
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
				ctx: context.TODO(),
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
					"correlation_id": "323232323232323232323",
				},
				Body:            fmt.Sprintf(`{"error":%d,"data":null}`, int(ResponseInvalidAgreement)),
				IsBase64Encoded: false,
			},
			wantErr: false,
		},
		{
			name: "#20: Missing Agreement",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/agree/{player_id}",
					Path:       "/api/v1/agree/12001",
					HTTPMethod: "POST",
					Headers: map[string]string{
						"Content-Type":   "application/json",
						"api_key":        "test-api-key",
						"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
						"correlation_id": "323232323232323232323",
					},
					QueryStringParameters: map[string]string{},
					PathParameters: map[string]string{
						"player_id": "12001",
					},
					Body: fmt.Sprintf(`[
								{
									"jurisdiction": "gdpr",
									"version": 3242343243,
									"locale": "en-uk",
									"timestamp": "%s"
								},
								{
									"jurisdiction": "gdpr",
									"agreement": "tou",
									"version": 214214214,
									"locale": "en-uk",
									"timestamp": "%s"
								}
							]`, currentTimeISO8601(), currentTimeISO8601()),
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
				ctx: context.TODO(),
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
					"correlation_id": "323232323232323232323",
				},
				Body:            fmt.Sprintf(`{"error":%d,"data":null}`, int(ResponseInvalidAgreement)),
				IsBase64Encoded: false,
			},
			wantErr: false,
		},
		{
			name: "#21: Invalid Locale",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/agree/{player_id}",
					Path:       "/api/v1/agree/12001",
					HTTPMethod: "POST",
					Headers: map[string]string{
						"Content-Type":   "application/json",
						"api_key":        "test-api-key",
						"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
						"correlation_id": "323232323232323232323",
					},
					QueryStringParameters: map[string]string{},
					PathParameters: map[string]string{
						"player_id": "12001",
					},
					Body: fmt.Sprintf(`[
								{
									"jurisdiction": "gdpr",
									"agreement": "pp",
									"version": 3242343243,
									"locale": "en-uk",
									"timestamp": "%s"
								},
								{
									"jurisdiction": "gdpr",
									"agreement": "tou",
									"version": 214214214,
									"locale": "en-ca",
									"timestamp": "%s"
								}
							]`, currentTimeISO8601(), currentTimeISO8601()),
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
				ctx: context.TODO(),
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
					"correlation_id": "323232323232323232323",
				},
				Body:            fmt.Sprintf(`{"error":%d,"data":null}`, int(ResponseInvalidLocale)),
				IsBase64Encoded: false,
			},
			wantErr: false,
		},
		{
			name: "#22: Missing Locale",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/agree/{player_id}",
					Path:       "/api/v1/agree/12001",
					HTTPMethod: "POST",
					Headers: map[string]string{
						"Content-Type":   "application/json",
						"api_key":        "test-api-key",
						"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
						"correlation_id": "323232323232323232323",
					},
					QueryStringParameters: map[string]string{},
					PathParameters: map[string]string{
						"player_id": "12001",
					},
					Body: fmt.Sprintf(`[
								{
									"jurisdiction": "gdpr",
									"agreement": "pp",
									"version": 3242343243,
									"locale": "en-uk",
									"timestamp": "%s"
								},
								{
									"jurisdiction": "gdpr",
									"agreement": "tou",
									"version": 214214214,
									"timestamp": "%s"
								}
							]`, currentTimeISO8601(), currentTimeISO8601()),
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
				ctx: context.TODO(),
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
					"correlation_id": "323232323232323232323",
				},
				Body:            fmt.Sprintf(`{"error":%d,"data":null}`, int(ResponseInvalidLocale)),
				IsBase64Encoded: false,
			},
			wantErr: false,
		},
		{
			name: "#23: Invalid Timestamp",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/agree/{player_id}",
					Path:       "/api/v1/agree/12001",
					HTTPMethod: "POST",
					Headers: map[string]string{
						"Content-Type":   "application/json",
						"api_key":        "test-api-key",
						"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
						"correlation_id": "323232323232323232323",
					},
					QueryStringParameters: map[string]string{},
					PathParameters: map[string]string{
						"player_id": "12001",
					},
					Body: fmt.Sprintf(`[
								{
									"jurisdiction": "gdpr",
									"agreement": "pp",
									"version": 3242343243,
									"locale": "en-uk",
									"timestamp": "%s"
								},
								{
									"jurisdiction": "gdpr",
									"agreement": "tou",
									"version": 214214214,
									"locale": "en-uk",
									"timestamp": "%s"
								}
							]`, "2023-05-02T10:29:40+00:00", currentTimeISO8601()),
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
				ctx: context.TODO(),
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
					"correlation_id": "323232323232323232323",
				},
				Body:            fmt.Sprintf(`{"error":%d,"data":null}`, int(ResponseInvalidTimestamp)),
				IsBase64Encoded: false,
			},
			wantErr: true,
		},
		{
			name: "#24: Missing Timestamp",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/agree/{player_id}",
					Path:       "/api/v1/agree/12001",
					HTTPMethod: "POST",
					Headers: map[string]string{
						"Content-Type":   "application/json",
						"api_key":        "test-api-key",
						"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
						"correlation_id": "323232323232323232323",
					},
					QueryStringParameters: map[string]string{},
					PathParameters: map[string]string{
						"player_id": "12001",
					},
					Body: fmt.Sprintf(`[
								{
									"jurisdiction": "gdpr",
									"agreement": "pp",
									"version": 3242343243,
									"locale": "en-uk",
									"timestamp": "%s"
								},
								{
									"jurisdiction": "gdpr",
									"agreement": "tou",
									"version": 214214214,
									"locale": "en-uk"									
								}
							]`, currentTimeISO8601()),
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
				ctx: context.TODO(),
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
					"correlation_id": "323232323232323232323",
				},
				Body:            fmt.Sprintf(`{"error":%d,"data":null}`, int(ResponseInvalidTimestamp)),
				IsBase64Encoded: false,
			},
			wantErr: false,
		},
		{
			name: "#25: Invalid Body Parameters",
			args: args{
				request: events.APIGatewayProxyRequest{
					Resource:   "/api/v1/agree/{player_id}",
					Path:       "/api/v1/agree/12001",
					HTTPMethod: "POST",
					Headers: map[string]string{
						"Content-Type":   "application/json",
						"api_key":        "test-api-key",
						"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
						"correlation_id": "323232323232323232323",
					},
					QueryStringParameters: map[string]string{},
					PathParameters: map[string]string{
						"player_id": "12001",
					},
					//comma not provide after jurisdiction
					Body: fmt.Sprintf(`[
								{
									"jurisdiction": "gdpr"

									"agreement": "pp",
									"version": 3242343243,
									"locale": "en-uk",
									"timestamp": "%s"
								},
								{
									"jurisdiction": "gdpr",
									"agreement": "tou",
									"version": 214214214,
									"locale": "en-uk"
									"timestamp": "%s"									
								}
							]`, currentTimeISO8601(), currentTimeISO8601()),
					IsBase64Encoded: false,
					RequestContext: events.APIGatewayProxyRequestContext{
						AccountID:  "123456789012",
						ResourceID: "resource-id",
						Stage:      "prod",
						RequestID:  "request-id",
						Identity: events.APIGatewayRequestIdentity{
							CognitoIdentityPoolID: "",
							AccountID:             "",
							CognitoIdentityID:     "",
							Caller:                "",
							APIKey:                "",
							SourceIP:              "127.0.0.1",
							AccessKey:             "",
						},
					},
				},
				ctx: context.TODO(),
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"request_id":     "8713F659-1ABF-4EDF-8E86-655295CF7CD9",
					"correlation_id": "323232323232323232323",
				},
				Body:            fmt.Sprintf(`{"error":%d,"data":null}`, int(v3.ResponseInvalidRequestBodyParameters)),
				IsBase64Encoded: false,
			},
			wantErr: false,
		},
	}

	env = NewEnv()

	dynamoSvc = &MockDynamoClient{}
	httpClient = &ClientMock{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, _ := HandleRequest(tt.args.ctx, tt.args.request)
			if !reflect.DeepEqual(response, tt.want) {
				t.Errorf("HandleRequest() response = %v, wantResponse %v", response, tt.want)
			}
		})
	}
}
