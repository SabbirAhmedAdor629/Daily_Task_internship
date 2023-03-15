package main

import (
	"encoding/json"
	//"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Property struct {
	
}

func lambdaHandler(event map[string]interface{}) (interface{}, error) {
	
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	
	var p Property
	err := json.Unmarshal([]byte(event["body"].(string)), &p)
	if err != nil {
		return nil, err
	}

	
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			
		},
		TableName: aws.String("Property"),
	}

	
	_, err = svc.PutItem(input)
	if err != nil {
		return nil, err
	}

	
	result, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String("Property"),
	})
	if err != nil {
		return nil, err
	}

	
	response, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return string(response), nil
}

func main() {
	lambda.Start(lambdaHandler)
}
