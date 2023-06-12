package sqs

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/google/uuid"
)

type (
	MockDynamoClient struct{ dynamodbiface.DynamoDBAPI }
	MockSqsClient    struct{ sqsiface.SQSAPI }
)

// func TestCreateSqsMessage(t *testing.T) {
// 	u := uuid.New()
// 	ts := time.Now()

// 	type args struct {
// 		bytes     []byte
// 		prefix    string
// 		origin    OriginType
// 		operation OperationType
// 		priority  PriorityType
// 		ts        time.Time
// 		uuid      uuid.UUID
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want SqsMessageType
// 	}{
// 		{
// 			name: "#1: Remove the api_key",
// 			args: args{
// 				bytes:     []byte(`{"api_key":"should_be_nulled","member_guid":"guid", "device_id":"12345", "transaction_id":"present"}`),
// 				prefix:    "mpr",
// 				uuid:      u,
// 				origin:    "messaging",
// 				operation: "register",
// 				priority:  HighPriority,
// 				ts:        ts,
// 			},
// 			want: SqsMessageType{
// 				"guid":           fmt.Sprintf("mpr-%s", u.String()),
// 				"transaction_id": "present",
// 				"origin":         OriginType("messaging"),
// 				"operation":      OperationType("register"),
// 				"priority":       PriorityType(HighPriority),
// 				"ts":             ts,
// 				"data": map[string]interface{}{
// 					"device_id":      "12345",
// 					"member_guid":    "guid",
// 					"transaction_id": "present"},
// 			},
// 		},
// 		{
// 			name: "#2: transaction_id present",
// 			args: args{
// 				bytes:     []byte(`{"api_key":"should_be_nulled","member_guid":"guid", "device_id":"12345", "transaction_id":"present"}`),
// 				prefix:    "mpr",
// 				uuid:      u,
// 				origin:    OriginType("messaging"),
// 				operation: "register",
// 				priority:  HighPriority,
// 				ts:        ts,
// 			},
// 			want: map[string]interface{}{
// 				"guid":           fmt.Sprintf("mpr-%s", u.String()),
// 				"transaction_id": "present",
// 				"origin":         OriginType("messaging"),
// 				"operation":      OperationType("register"),
// 				"priority":       PriorityType(HighPriority),
// 				"ts":             ts,
// 				"data": map[string]interface{}{
// 					"device_id":      "12345",
// 					"member_guid":    "guid",
// 					"transaction_id": "present"},
// 			},
// 		},
// 		{
// 			name: "#3: transaction_id not present",
// 			args: args{
// 				bytes:     []byte(`{"api_key":"should_be_nulled","member_guid":"guid","device_id":"12345"}`),
// 				prefix:    "mpr",
// 				uuid:      u,
// 				origin:    OriginType("messaging"),
// 				operation: "register",
// 				priority:  HighPriority,
// 				ts:        ts,
// 			},
// 			want: map[string]interface{}{
// 				"guid":           fmt.Sprintf("mpr-%s", u.String()),
// 				"transaction_id": u,
// 				"origin":         OriginType("messaging"),
// 				"operation":      OperationType("register"),
// 				"priority":       PriorityType(HighPriority),
// 				"ts":             ts,
// 				"data": map[string]interface{}{
// 					"device_id":      "12345",
// 					"member_guid":    "guid",
// 					"transaction_id": u},
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := CreateSqsMessage(tt.args.bytes, tt.args.prefix, tt.args.origin, tt.args.operation, tt.args.priority, tt.args.ts, tt.args.uuid)
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("CreateSqsMessage() got = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func TestIngCreateSqsMessageTest(t *testing.T) {

	tid := "617dd2c4-4f37-11ed-bdc3-0242ac120002"
	u, _ := uuid.Parse(tid)
	u2 := uuid.New()

	ts := time.Now()

	test1 := map[string]interface{}{
		"guid":           fmt.Sprintf("mpr-%s", u.String()),
		"transaction_id": u,
		"origin":         OriginType("messaging"),
		"operation":      OperationType("register"),
		"priority":       PriorityType(HighPriority),
		"ts":             ts,
		"data": map[string]interface{}{
			"device_id":      "12345",
			"member_guid":    "guid",
			"transaction_id": u}}
	test1Bytes, _ := json.Marshal(test1)
	test1String := string(test1Bytes)

	test2 := map[string]interface{}{
		"guid":           fmt.Sprintf("mpr-%s", u2.String()),
		"transaction_id": u,
		"origin":         OriginType("messaging"),
		"operation":      OperationType("register"),
		"priority":       PriorityType(HighPriority),
		"ts":             ts,
		"data": map[string]interface{}{
			"device_id":      "12345",
			"member_guid":    "guid",
			"transaction_id": u}}
	test2Bytes, _ := json.Marshal(test2)
	test2String := string(test2Bytes)

	test3 := map[string]interface{}{
		"guid":           fmt.Sprintf("mpr-%s", u.String()),
		"transaction_id": u,
		"origin":         OriginType("messaging"),
		"operation":      OperationType("register"),
		"priority":       PriorityType(HighPriority),
		"ts":             ts,
		"data": map[string]interface{}{
			"device_id":      "12345",
			"member_guid":    "guid",
			"transaction_id": u}}
	test3Bytes, _ := json.Marshal(test3)
	test3String := string(test3Bytes)

	type args struct {
		bytes     []byte
		prefix    string
		origin    OriginType
		operation OperationType
		priority  PriorityType
		ts        time.Time
		uuid      uuid.UUID
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "#1: =======> Remove the api_key, transaction_id is present but not a UUID",
			args: args{
				bytes:     []byte(`{"api_key":"should_be_nulled","member_guid":"guid", "device_id":"12345", "transaction_id":"present"}`),
				prefix:    "mpr",
				uuid:      u,
				origin:    "messaging",
				operation: "register",
				priority:  HighPriority,
				ts:        ts,
			},
			want: test1String,
		},
		{
			name: "#2: =======> transaction_id present",
			args: args{
				bytes:     []byte(`{"api_key":"should_be_nulled","member_guid":"guid", "device_id":"12345", "transaction_id":"617dd2c4-4f37-11ed-bdc3-0242ac120002"}`),
				prefix:    "mpr",
				uuid:      u2,
				origin:    OriginType("messaging"),
				operation: "register",
				priority:  HighPriority,
				ts:        ts,
			},
			want: test2String,
		},
		{
			name: "#3: =======> transaction_id not present",
			args: args{
				bytes:     []byte(`{"api_key":"should_be_nulled","member_guid":"guid","device_id":"12345"}`),
				prefix:    "mpr",
				uuid:      u,
				origin:    OriginType("messaging"),
				operation: "register",
				priority:  HighPriority,
				ts:        ts,
			},
			want: test3String,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateSqsMessageTest(tt.args.bytes, tt.args.prefix, tt.args.origin, tt.args.operation, tt.args.priority, tt.args.ts, tt.args.uuid); got != tt.want {
				t.Errorf("CreateSqsMessageTest() = %v, want %v", got, tt.want)
			}
		})
	}
}
