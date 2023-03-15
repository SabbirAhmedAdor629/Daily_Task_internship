package main

import (
	"fmt"
	"strconv"

	//"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type (
	MockDynamoClient struct {
		dynamodbiface.DynamoDBAPI
	}
)

func (svc *MockDynamoClient) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	id, err := strconv.Atoi(*input.Key["Id"].N) // .N only working for first digit
	if err != nil {
		return nil, fmt.Errorf("invalid input key")
	}

	if id == 12 {
		person := Person{
			Id:      12,
			Name:    "Sabbir Ahmed",
			Website: "https://sabbirahmedador629.github.io/Personal_Portfolio/",
		}
		item, _ := dynamodbattribute.MarshalMap(person)
		return &dynamodb.GetItemOutput{
			Item: item,
		}, nil
	}

	return nil, awserr.New(dynamodb.ErrCodeResourceNotFoundException, "not found", nil)
}
