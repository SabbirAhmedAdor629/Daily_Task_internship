package api_key

import (
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

func (svc *MockDynamoClient) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	switch {
	case *input.Key["api_key"].S == "test-2":
		// Not found
		return &dynamodb.GetItemOutput{
			ConsumedCapacity: &dynamodb.ConsumedCapacity{},
			Item:             nil,
		}, nil
	case *input.Key["api_key"].S == "test-3":
		// Found
		return &dynamodb.GetItemOutput{
			ConsumedCapacity: &dynamodb.ConsumedCapacity{},
			Item: map[string]*dynamodb.AttributeValue{
				"api_key": {
					S: aws.String("test-3_api_key_token"),
				}},
		}, nil
	case *input.Key["api_key"].S == "test-4":
		// Error
		return &dynamodb.GetItemOutput{
			ConsumedCapacity: &dynamodb.ConsumedCapacity{},
			Item:             map[string]*dynamodb.AttributeValue{},
		}, errors.New("")
	default:
		// Not found
		return &dynamodb.GetItemOutput{
			ConsumedCapacity: &dynamodb.ConsumedCapacity{},
			Item:             nil,
		}, nil
	}
}

func TestLoadCache(t *testing.T) {
	type args struct {
		keys []string
	}
	tests := []struct {
		name string
		args args
		want ApiKeyType
	}{
		{
			name: "#1",
			args: args{
				keys: nil,
			},
			want: map[string]int{},
		},
		{
			name: "#2",
			args: args{
				keys: []string{"abc", "def"},
			},
			want: map[string]int{"abc": 1, "def": 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LoadCache(tt.args.keys); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadApiCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApiKeyType_ValidateKey(t *testing.T) {
	type args struct {
		svc       dynamodbiface.DynamoDBAPI
		tableName string
		keyName   string
		apiKey    string
	}
	tests := []struct {
		name    string
		a       *ApiKeyType
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "#1.1: No Table name",
			a:    &ApiKeyType{},
			args: args{
				svc:       &MockDynamoClient{},
				tableName: "",
				keyName:   "api_key",
				apiKey:    "test-0",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "#1.2: No Key name",
			a:    &ApiKeyType{},
			args: args{
				svc:       &MockDynamoClient{},
				tableName: "api_keys_production",
				keyName:   "",
				apiKey:    "test-0",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "#1.3: No Api key",
			a:    &ApiKeyType{},
			args: args{
				svc:       &MockDynamoClient{},
				tableName: "api_keys_production",
				keyName:   "api_key",
				apiKey:    "",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "#2: Result already cached",
			a:    &ApiKeyType{"test-1": 1},
			args: args{
				svc:       &MockDynamoClient{},
				tableName: "api_keys_production",
				keyName:   "api_key",
				apiKey:    "test-1",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "#3: Result not cached, and not found in DynamoDB",
			a:    &ApiKeyType{"test-1": 1},
			args: args{
				svc:       &MockDynamoClient{},
				tableName: "api_keys_production",
				keyName:   "api_key",
				apiKey:    "test-2",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "#4: Result not cached, found in DynamoDB, then cached",
			a:    &ApiKeyType{"test-1": 1},
			args: args{
				svc:       &MockDynamoClient{},
				tableName: "api_keys_production",
				keyName:   "api_key",
				apiKey:    "test-3",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "#5: Cannot query DynamoDB",
			a:    &ApiKeyType{"test-1": 1},
			args: args{
				svc:       &MockDynamoClient{},
				tableName: "api_keys_production",
				keyName:   "api_key",
				apiKey:    "test-4",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.ValidateKey(tt.args.svc, tt.args.tableName, tt.args.keyName, tt.args.apiKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("ApiKeyType.ValidateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ApiKeyType.ValidateKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
