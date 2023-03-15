package main

import (
	// "errors"
	// "strconv"
	// "reflect"
	"reflect"
	"testing"

	// "github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/dynamodb"
	// "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	// "github.com/aws/aws-sdk-go/service/dynamodb"
	// "github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/assert"
)

func TestGetItem(t *testing.T) {
	type args struct {
		id int
	}
	type testcase struct {
		name    string
		args    args
		dynamo  *MockDynamoClient
		want    Person
		wantErr bool
	}
	tests := []testcase{
		{
			name: "#1: Valid item",
			args: args{
				id: 12,
			},
			dynamo: &MockDynamoClient{},
			want: Person{
				Id:      12,
				Name:    "Sabbir Ahmed",
				Website: "https://sabbirahmedador629.github.io/Personal_Portfolio/",
			},
			wantErr: false,
		},
		{
			name: "#2: Invalid item",
			args: args{
				id: 2,
			},
			dynamo:  &MockDynamoClient{},
			want:    Person{},
			wantErr: true,
		},
	}

	dynamoSvc := &MockDynamoClient{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err:= GetItem(dynamoSvc, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetItem() = %v, want %v", got, tt.want)
			}
		})
	}
}
