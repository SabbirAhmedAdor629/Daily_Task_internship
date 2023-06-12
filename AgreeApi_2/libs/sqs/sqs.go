package sqs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/google/uuid"
)

type (
	OriginType     string
	OperationType  string
	PriorityType   string
	SqsMessageType map[string]interface{}
	SQSDataType    map[string]interface{}
)

const (
	HighPriority   PriorityType = "high"
	MediumPriority PriorityType = "medium"
	LowPriority    PriorityType = "low"
)

func CreateSqsMessage(bytes []byte, prefix string, origin OriginType, operation OperationType, priority PriorityType, ts time.Time, guid uuid.UUID) SqsMessageType {
	var sqsMessage SqsMessageType = make(SqsMessageType)
	var data SQSDataType = make(SQSDataType)

	json.Unmarshal(bytes, &data)

	delete(data, "api_key")
	delete(data, "advertising_id")

	sqsMessage["origin"] = origin
	sqsMessage["operation"] = operation
	sqsMessage["priority"] = priority
	sqsMessage["ts"] = ts

	sqsMessage["guid"] = prefix + "-" + guid.String() //fmt.Sprintf("%s-%s", prefix, uuid)

	if tid, present := data["transaction_id"]; present {
		if u, err := uuid.Parse(fmt.Sprintf("%s", tid)); err != nil {
			data["transaction_id"] = guid
		} else {
			data["transaction_id"] = u
		}
	} else {
		data["transaction_id"] = guid
	}

	sqsMessage["transaction_id"] = data["transaction_id"]
	sqsMessage["data"] = data

	return sqsMessage
}

func SendSqsMessage(svc sqsiface.SQSAPI, body string, queueUrl *string, delay int64) (*sqs.SendMessageOutput, error) {
	return svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(delay),
		QueueUrl:     queueUrl,
		MessageBody:  aws.String(body),
	})
}

func CreateSqsMessageTest(bytes []byte, prefix string, origin OriginType, operation OperationType, priority PriorityType, ts time.Time, uuid uuid.UUID) string {
	msg := CreateSqsMessage(bytes, prefix, origin, operation, priority, ts, uuid)
	b, _ := json.Marshal(msg)
	return string(b)
}
