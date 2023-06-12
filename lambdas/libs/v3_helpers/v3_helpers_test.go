package v3_helpers

import (
	"errors"
	"net/http"
	"reflect"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/google/uuid"
	im_keycache "influencemobile.com/libs/api_key"
)

type (
	MockDynamoClient struct{ dynamodbiface.DynamoDBAPI }
	MockSqsClient    struct{ sqsiface.SQSAPI }
)

func (svc *MockDynamoClient) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	switch {
	case *input.Key["api_key"].S == "test-2":
		return &dynamodb.GetItemOutput{
			ConsumedCapacity: &dynamodb.ConsumedCapacity{},
			Item:             nil,
		}, nil
	case *input.Key["api_key"].S == "test-3":
		return &dynamodb.GetItemOutput{
			ConsumedCapacity: &dynamodb.ConsumedCapacity{},
			Item: map[string]*dynamodb.AttributeValue{
				"api_key": {
					S: aws.String("test-3_api_key_token"),
				}},
		}, nil
	case *input.Key["api_key"].S == "test-4":
		return &dynamodb.GetItemOutput{
			ConsumedCapacity: &dynamodb.ConsumedCapacity{},
			Item:             map[string]*dynamodb.AttributeValue{},
		}, errors.New("Internal DynamoDB Error")
	default:
		return &dynamodb.GetItemOutput{
			ConsumedCapacity: &dynamodb.ConsumedCapacity{},
			Item:             map[string]*dynamodb.AttributeValue{},
		}, nil
	}
}

func TestValidateRequestMethod(t *testing.T) {
	type args struct {
		requestMethod string
		method        HttpMethodType
	}
	tests := []struct {
		name    string
		args    args
		want    ResponseCodeType
		wantErr bool
	}{
		{
			name: "#1: POST received",
			args: args{
				requestMethod: "POST",
				method:        "POST",
			},
			want:    ResponseSuccess,
			wantErr: false,
		},
		{
			name: "#2: GET Received",
			args: args{
				requestMethod: "GET",
				method:        "POST",
			},
			want:    ResponseInvalidRequestMethod,
			wantErr: true,
		},
		{
			name: "#2: Empty method Received",
			args: args{
				requestMethod: "",
				method:        "POST",
			},
			want:    ResponseInvalidRequestMethod,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateRequestMethod(tt.args.requestMethod, tt.args.method)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRequestMethod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateRequestMethod() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateRequestPath(t *testing.T) {
	type args struct {
		requestPath string
		path        *regexp.Regexp
	}
	tests := []struct {
		name    string
		args    args
		want    ResponseCodeType
		wantErr bool
	}{
		{
			name: "#1: Expected path received",
			args: args{
				requestPath: "/api/v2/push_registrations.json",
				path:        regexp.MustCompile("^/api/v2/push_registrations.json$"),
			},
			want:    ResponseSuccess,
			wantErr: false,
		},
		{
			name: "#2: Unexcpected path received",
			args: args{
				requestPath: "/api/v3/push_registrations.json",
				path:        regexp.MustCompile("^/api/v2/push_registrations.json$"),
			},
			want:    ResponseInvalidRequestPath,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateRequestPath(tt.args.requestPath, tt.args.path)
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

func TestValidateRequestHeaders(t *testing.T) {
	type args struct {
		request         map[string]string
		mandatoryParams []string
		responseCodes   map[string]int
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		want1   int
		wantErr bool
	}{
		{
			name: "#1: Good request",
			args: args{
				request: map[string]string{
					"api_key":          "abcde",
					"request_id":       "617dd2c4-4f37-11ed-bdc3-0242ac120002",
					"platform":         "android",
					"program_slug":     "rewarded_play",
					"platform_version": "100",
					"build_version":    "12.0.1",
				},
				mandatoryParams: MandatoryHeaders,
				responseCodes:   ResponseCodeToKey,
			},
			want: map[string]string{
				"api_key":          "abcde",
				"request_id":       "617dd2c4-4f37-11ed-bdc3-0242ac120002",
				"platform":         "android",
				"program_slug":     "rewarded_play",
				"platform_version": "100",
				"build_version":    "12.0.1",
			},
			want1:   0,
			wantErr: false,
		},
		{
			name: "#2: No api_key",
			args: args{
				request: map[string]string{
					"request_id":       "617dd2c4-4f37-11ed-bdc3-0242ac120002",
					"platform":         "android",
					"program_slug":     "rewarded_play",
					"platform_version": "100",
					"build_version":    "12.0.1",
				},
				mandatoryParams: MandatoryHeaders,
				responseCodes:   ResponseCodeToKey,
			},
			want: map[string]string{
				"request_id":       "617dd2c4-4f37-11ed-bdc3-0242ac120002",
				"platform":         "android",
				"program_slug":     "rewarded_play",
				"platform_version": "100",
				"build_version":    "12.0.1",
			},
			want1:   int(ResponseMissingApiKey),
			wantErr: true,
		},
		{
			name: "#3: Invalid request_id",
			args: args{
				request: map[string]string{
					"api_key":          "abcde",
					"request_id":       "617dd2c4-4f37-11ed-bdc3-0242ac1200022",
					"platform":         "android",
					"program_slug":     "rewarded_play",
					"platform_version": "100",
					"build_version":    "12.0.1",
				},
				mandatoryParams: MandatoryHeaders,
				responseCodes:   ResponseCodeToKey,
			},
			want: map[string]string{
				"api_key":          "abcde",
				"request_id":       "617dd2c4-4f37-11ed-bdc3-0242ac1200022",
				"platform":         "android",
				"program_slug":     "rewarded_play",
				"platform_version": "100",
				"build_version":    "12.0.1",
			},
			want1:   int(ResponseInvalidRequestIdHeader),
			wantErr: true,
		},
		{
			name: "#4: request_id not present",
			args: args{
				request: map[string]string{
					"api_key":          "abcde",
					"platform":         "android",
					"program_slug":     "rewarded_play",
					"platform_version": "100",
					"build_version":    "12.0.1",
				},
				mandatoryParams: MandatoryHeaders,
				responseCodes:   ResponseCodeToKey,
			},
			want: map[string]string{
				"api_key":          "abcde",
				"platform":         "android",
				"program_slug":     "rewarded_play",
				"platform_version": "100",
				"build_version":    "12.0.1",
			},
			want1:   int(ResponseMissingRequestIdHeader),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ValidateRequestHeaders(tt.args.request, tt.args.mandatoryParams, tt.args.responseCodes)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateBodyParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateBodyParams() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ValidateBodyParams() got = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestValidateBodyParameters(t *testing.T) {
	type args struct {
		requestBody     string
		isBase64Encoded bool
		mandatoryParams []string
		responseCodes   map[string]int
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		want1   int
		wantErr bool
	}{
		{
			name: "#1: Good request",
			args: args{
				requestBody:     "{ \"request_id\":\"617dd2c4-4f37-11ed-bdc3-0242ac120002\",\"api_key\": \"abcde\", \"member_guid\": \"1234\",\"device_id\": \"1234\",\"program_id\": \"1234\",\"device_info\": {},\"consumer\": {}}",
				isBase64Encoded: false,
				mandatoryParams: []string{"request_id"},
				responseCodes:   ResponseCodeToKey,
			},
			want: map[string]interface{}{
				"request_id":  "617dd2c4-4f37-11ed-bdc3-0242ac120002",
				"api_key":     "abcde",
				"member_guid": "1234",
				"device_id":   "1234",
				"program_id":  "1234",
				"device_info": map[string]interface{}{},
				"consumer":    map[string]interface{}{},
			},
			want1:   0,
			wantErr: false,
		},
		{
			name: "#2: No api_key",
			args: args{
				requestBody:     "{ \"request_id\":\"617dd2c4-4f37-11ed-bdc3-0242ac120002\",\"member_guid\": \"1234\" }",
				isBase64Encoded: false,
				mandatoryParams: []string{"request_id", "member_guid", "api_key"},
				responseCodes:   ResponseCodeToKey,
			},
			want:    nil,
			want1:   int(ResponseMissingApiKey),
			wantErr: true,
		},
		{
			name: "#2: No mandatory params",
			args: args{
				requestBody:     "{ \"request_id\":\"617dd2c4-4f37-11ed-bdc3-0242ac120002\",\"member_guid\": \"1234\" }",
				isBase64Encoded: false,
				mandatoryParams: nil,
				responseCodes:   ResponseCodeToKey,
			},
			want: map[string]interface{}{
				"request_id":  "617dd2c4-4f37-11ed-bdc3-0242ac120002",
				"member_guid": "1234",
			},
			want1:   0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ValidateBodyParameters(tt.args.requestBody, tt.args.isBase64Encoded, tt.args.mandatoryParams, tt.args.responseCodes)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateBodyParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateBodyParams() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ValidateBodyParams() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestValidateAuth(t *testing.T) {
	type args struct {
		svc       dynamodbiface.DynamoDBAPI
		apiCache  *im_keycache.ApiKeyType
		tableName string
		apiKey    string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		want1   int
		wantErr bool
	}{
		{
			name: "#1: Cached",
			args: args{
				svc:       &MockDynamoClient{},
				apiCache:  &im_keycache.ApiKeyType{"test-1": 1},
				tableName: "api_keys_production",
				apiKey:    "test-1",
			},
			want:    true,
			want1:   int(ResponseSuccess),
			wantErr: false,
		},
		{
			name: "#2: Retrieved from DynamoDB",
			args: args{
				svc:       &MockDynamoClient{},
				apiCache:  &im_keycache.ApiKeyType{"test-1": 1},
				tableName: "api_keys_production",
				apiKey:    "test-3",
			},
			want:    true,
			want1:   int(ResponseSuccess),
			wantErr: false,
		},
		{
			name: "#3: Not found",
			args: args{
				svc:       &MockDynamoClient{},
				apiCache:  &im_keycache.ApiKeyType{"test-1": 1},
				tableName: "api_keys_production",
				apiKey:    "test-2",
			},
			want:    false,
			want1:   int(ResponseInvalidApiKey),
			wantErr: false,
		},
		{
			name: "#4: Internal DynamoDB error",
			args: args{
				svc:       &MockDynamoClient{},
				apiCache:  &im_keycache.ApiKeyType{"test-1": 1},
				tableName: "api_keys_production",
				apiKey:    "test-4",
			},
			want:    false,
			want1:   int(ResponseInternalServerError),
			wantErr: true,
		},
		{
			name: "#5: api_key not specified",
			args: args{
				svc:       &MockDynamoClient{},
				apiCache:  &im_keycache.ApiKeyType{"test-1": 1},
				tableName: "api_keys_production",
			},
			want:    false,
			want1:   int(ResponseMissingApiKey),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ValidateAuth(tt.args.svc, tt.args.apiCache, tt.args.tableName, tt.args.apiKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAuth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateAuth() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ValidateAuth() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateRequestId(t *testing.T) {
	valid_tid := "617dd2c4-4f37-11ed-bdc3-0242ac120002"
	valid_u, _ := uuid.Parse(valid_tid)
	invalid_tid := "617dd2c4-4f37-11ed-bdc3-0242ac12000"

	type args struct {
		tid string
	}
	tests := []struct {
		name    string
		args    args
		want    *uuid.UUID
		wantErr bool
	}{
		{
			name: "#1: Valid transaction_id",
			args: args{
				tid: valid_tid,
			},
			want:    &valid_u,
			wantErr: false,
		},
		{
			name: "#2: Invalid transaction_id",
			args: args{
				tid: invalid_tid,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateRequestId(tt.args.tid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTransactionId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateTransactionId() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPtrS(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want *string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPtrS(tt.args.s); got != tt.want {
				t.Errorf("NewPtrS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPtrB(t *testing.T) {
	pVal := new(bool)
	*pVal = true
	type args struct {
		b bool
	}
	tests := []struct {
		name string
		args args
		want *bool
	}{
		{
			name: "Test #1",
			args: args{
				b: true,
			},
			want: pVal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPtrB(tt.args.b); *got != *tt.want {
				t.Errorf("NewPtrB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPtrI(t *testing.T) {
	pVal := new(int)
	*pVal = 100
	type args struct {
		i int
	}
	tests := []struct {
		name string
		args args
		want *int
	}{
		{
			name: "Test #1",
			args: args{
				i: 100,
			},
			want: pVal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPtrI(tt.args.i); *got != *tt.want {
				t.Errorf("NewPtrI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateParamI(t *testing.T) {
	type args struct {
		name  string
		value *int64
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		want1   string
		want2   int
		wantErr bool
	}{
		{
			name: "#1",
			args: args{
				name:  "member_id",
				value: nil,
			},
			want:    0,
			want1:   "[ Process Failed. ] A valid member_id is required",
			want2:   http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "#2",
			args: args{
				name:  "memberId",
				value: new(int64),
			},
			want:    0,
			want1:   "",
			want2:   http.StatusOK,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := ValidateParamI(tt.args.name, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateParamI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateParamI() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ValidateParamI() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("ValidateParamI() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestValidateParamS(t *testing.T) {
	type args struct {
		name  string
		value *string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		want2   int
		wantErr bool
	}{
		{
			name: "#1",
			args: args{
				name:  "member_id",
				value: nil,
			},
			want:    "",
			want1:   "[ Process Failed. ] A valid member_id is required",
			want2:   http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "#2",
			args: args{
				name:  "member_id",
				value: new(string),
			},
			want:    "",
			want1:   "[ Process Failed. ] A valid member_id is required",
			want2:   http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "#3",
			args: args{
				name:  "member_id",
				value: NewPtrS("1234"),
			},
			want:    "1234",
			want1:   "",
			want2:   http.StatusOK,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := ValidateParamS(tt.args.name, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateParamS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateParamS() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ValidateParamS() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("ValidateParamS() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestValidateAdId(t *testing.T) {
	type args struct {
		adId *string
	}
	tests := []struct {
		name    string
		args    args
		want    AdIdType
		want1   int
		wantErr bool
	}{
		{
			name:    "#1: ",
			args:    args{adId: NewPtrS("12345")},
			want:    "12345",
			want1:   int(ResponseSuccess),
			wantErr: false,
		},
		{
			name:    "#2: ",
			args:    args{adId: NewPtrS("")},
			want:    "",
			want1:   int(ResponseMissingAdvertisingId),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ValidateAdId(tt.args.adId)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAdId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateAdId() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ValidateMessageId() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestValidateAppId(t *testing.T) {
	type args struct {
		appId *string
	}
	tests := []struct {
		name    string
		args    args
		want    AppIdType
		want1   int
		wantErr bool
	}{
		{
			name:    "#1: ",
			args:    args{appId: NewPtrS("abcde")},
			want:    "abcde",
			want1:   int(ResponseSuccess),
			wantErr: false,
		},
		{
			name:    "#2: ",
			args:    args{appId: NewPtrS("")},
			want:    "",
			want1:   int(ResponseMissingAppId),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ValidateAppId(tt.args.appId)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAppId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateAppId() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ValidateMessageId() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestValidatePlayerId(t *testing.T) {
	type args struct {
		playerId *string
	}
	tests := []struct {
		name    string
		args    args
		want    PlayerIdType
		want1   int
		wantErr bool
	}{
		{
			name: "#1: Nil",
			args: args{
				playerId: nil,
			},
			want:    0,
			want1:   int(ResponseMissingMember),
			wantErr: true,
		},
		{
			name: "#2: Valid",
			args: args{
				playerId: NewPtrS("1234"),
			},
			want:    1234,
			want1:   int(ResponseSuccess),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ValidatePlayerId(tt.args.playerId)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePlayerId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidatePlayerId() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ValidatePlayerId() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestValidateMessageId(t *testing.T) {
	type args struct {
		messageId *string
	}
	tests := []struct {
		name    string
		args    args
		want    MessageIdType
		want1   int
		wantErr bool
	}{
		{
			name: "#1: Nil",
			args: args{
				messageId: nil,
			},
			want:    "",
			want1:   int(ResponseMissingMessageId),
			wantErr: true,
		},
		{
			name: "#2: Empty",
			args: args{
				messageId: NewPtrS(""),
			},
			want:    "",
			want1:   int(ResponseMissingMessageId),
			wantErr: true,
		},
		{
			name: "#3: Valid",
			args: args{
				messageId: NewPtrS("123"),
			},
			want:    "123",
			want1:   int(ResponseSuccess),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ValidateMessageId(tt.args.messageId)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMessageId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateMessageId() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ValidateMessageId() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestValidateMemberId(t *testing.T) {
	type args struct {
		memberId *string
	}
	tests := []struct {
		name    string
		args    args
		want    MemberIdType
		want1   int
		wantErr bool
	}{
		{
			name: "#1: Nil",
			args: args{
				memberId: nil,
			},
			want:    "",
			want1:   int(ResponseMissingMember),
			wantErr: true,
		},
		{
			name: "#2: Valid",
			args: args{
				memberId: NewPtrS("123"),
			},
			want:    "123",
			want1:   int(ResponseSuccess),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ValidateMemberId(tt.args.memberId)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMemberId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateMemberId() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ValidateMemberId() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
