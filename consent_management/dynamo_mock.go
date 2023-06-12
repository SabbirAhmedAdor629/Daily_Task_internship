package main

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type MockDynamoClient struct{ dynamodbiface.DynamoDBAPI }

func (svc *MockDynamoClient) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	tableName := *input.TableName
	if v, ok := input.ExpressionAttributeNames["#GSI"]; ok {
		if tableName == "jurisdiction_mapping_test" && *v == "jurisdiction" {
			if input.ExpressionAttributeValues != nil {
				if _, ok := input.ExpressionAttributeValues[":GSI"]; ok {
					switch {
					case *input.ExpressionAttributeValues[":GSI"].S == "gdpr":
						return &dynamodb.QueryOutput{
							ConsumedCapacity: &dynamodb.ConsumedCapacity{},
							Count:            new(int64),
							Items: []map[string]*dynamodb.AttributeValue{
								// Descending order after querying the result based on the sort key [version]
								{
									"guid": &dynamodb.AttributeValue{
										S: aws.String("juris-e1f30f7d-5643"),
									},
									"jurisdiction": &dynamodb.AttributeValue{
										S: aws.String("gdpr"),
									},
									"agreement": &dynamodb.AttributeValue{
										S: aws.String("tou"),
									},
									"version": &dynamodb.AttributeValue{
										N: aws.String("1681118745"),
									},
									"agreement_locale": &dynamodb.AttributeValue{
										S: aws.String("en-uk"),
									},
								},
								{
									"guid": &dynamodb.AttributeValue{
										S: aws.String("juris-e1f30f7d-5643"),
									},
									"jurisdiction": &dynamodb.AttributeValue{
										S: aws.String("gdpr"),
									},
									"agreement": &dynamodb.AttributeValue{
										S: aws.String("tou"),
									},
									"version": &dynamodb.AttributeValue{
										N: aws.String("1681118712"),
									},
									"agreement_locale": &dynamodb.AttributeValue{
										S: aws.String("en-uk"),
									},
								},
								{
									"guid": &dynamodb.AttributeValue{
										S: aws.String("juris-e1f30f7d-5643"),
									},
									"jurisdiction": &dynamodb.AttributeValue{
										S: aws.String("gdpr"),
									},
									"agreement": &dynamodb.AttributeValue{
										S: aws.String("tou"),
									},
									"version": &dynamodb.AttributeValue{
										N: aws.String("1681118689"),
									},
									"agreement_locale": &dynamodb.AttributeValue{
										S: aws.String("en-uk"),
									},
								},
								{
									"guid": &dynamodb.AttributeValue{
										S: aws.String("juris-e1f30f7d-38f2"),
									},
									"jurisdiction": &dynamodb.AttributeValue{
										S: aws.String("gdpr"),
									},
									"agreement": &dynamodb.AttributeValue{
										S: aws.String("tou"),
									},
									"version": &dynamodb.AttributeValue{
										N: aws.String("1681118686"),
									},
									"agreement_locale": &dynamodb.AttributeValue{
										S: aws.String("en-uk"),
									},
								},
								{
									"guid": &dynamodb.AttributeValue{
										S: aws.String("juris-e1f30f7d-5465"),
									},
									"jurisdiction": &dynamodb.AttributeValue{
										S: aws.String("gdpr"),
									},
									"agreement": &dynamodb.AttributeValue{
										S: aws.String("pp"),
									},
									"version": &dynamodb.AttributeValue{
										N: aws.String("1681118791"),
									},
									"agreement_locale": &dynamodb.AttributeValue{
										S: aws.String("en-uk"),
									},
								},
								{
									"guid": &dynamodb.AttributeValue{
										S: aws.String("juris-e1f30f7d-5465"),
									},
									"jurisdiction": &dynamodb.AttributeValue{
										S: aws.String("gdpr"),
									},
									"agreement": &dynamodb.AttributeValue{
										S: aws.String("pp"),
									},
									"version": &dynamodb.AttributeValue{
										N: aws.String("1681118781"),
									},
									"agreement_locale": &dynamodb.AttributeValue{
										S: aws.String("en-uk"),
									},
								},
								{
									"guid": &dynamodb.AttributeValue{
										S: aws.String("juris-e1f30f7d-5465"),
									},
									"jurisdiction": &dynamodb.AttributeValue{
										S: aws.String("gdpr"),
									},
									"agreement": &dynamodb.AttributeValue{
										S: aws.String("pp"),
									},
									"version": &dynamodb.AttributeValue{
										N: aws.String("1681118765"),
									},
									"agreement_locale": &dynamodb.AttributeValue{
										S: aws.String("en-uk"),
									},
								},
							},
							LastEvaluatedKey: nil,
							ScannedCount:     new(int64),
						}, nil
					case *input.ExpressionAttributeValues[":GSI"].S == "eu":
						// Descending order after querying the result based on the sort key [version]
						return &dynamodb.QueryOutput{
							ConsumedCapacity: &dynamodb.ConsumedCapacity{},
							Count:            new(int64),
							Items: []map[string]*dynamodb.AttributeValue{
								{
									"guid": &dynamodb.AttributeValue{
										S: aws.String("juris-e1f30f7d-6655"),
									},
									"jurisdiction": &dynamodb.AttributeValue{
										S: aws.String("eu"),
									},
									"agreement": &dynamodb.AttributeValue{
										S: aws.String("pp"),
									},
									"version": &dynamodb.AttributeValue{
										N: aws.String("1680845499"),
									},
									"agreement_locale": &dynamodb.AttributeValue{
										S: aws.String("en-uk"),
									},
								},
								{
									"guid": &dynamodb.AttributeValue{
										S: aws.String("juris-e1f30f7d-7865"),
									},
									"jurisdiction": &dynamodb.AttributeValue{
										S: aws.String("eu"),
									},
									"agreement": &dynamodb.AttributeValue{
										S: aws.String("pp"),
									},
									"version": &dynamodb.AttributeValue{
										N: aws.String("1680845478"),
									},
									"agreement_locale": &dynamodb.AttributeValue{
										S: aws.String("en-uk"),
									},
								},
								{
									"guid": &dynamodb.AttributeValue{
										S: aws.String("juris-e1f30f7d-6334"),
									},
									"jurisdiction": &dynamodb.AttributeValue{
										S: aws.String("eu"),
									},
									"agreement": &dynamodb.AttributeValue{
										S: aws.String("tou"),
									},
									"version": &dynamodb.AttributeValue{
										N: aws.String("1680846146"),
									},
									"agreement_locale": &dynamodb.AttributeValue{
										S: aws.String("en-uk"),
									},
								},
								{
									"guid": &dynamodb.AttributeValue{
										S: aws.String("juris-e1f30f7d-7865"),
									},
									"jurisdiction": &dynamodb.AttributeValue{
										S: aws.String("eu"),
									},
									"agreement": &dynamodb.AttributeValue{
										S: aws.String("tou"),
									},
									"version": &dynamodb.AttributeValue{
										N: aws.String("1680845408"),
									},
									"agreement_locale": &dynamodb.AttributeValue{
										S: aws.String("en-uk"),
									},
								},
							},
							LastEvaluatedKey: nil,
							ScannedCount:     new(int64),
						}, nil
					case *input.ExpressionAttributeValues[":GSI"].S == "pipl":
						return &dynamodb.QueryOutput{
							ConsumedCapacity: &dynamodb.ConsumedCapacity{},
							Count:            new(int64),
							Items:            []map[string]*dynamodb.AttributeValue{},
						}, nil
					}

				}

			}
		}
	}

	return &dynamodb.QueryOutput{}, fmt.Errorf("something went wrong")
}

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
		}, errors.New("internal DynamoDB error")
	default:
		return &dynamodb.GetItemOutput{
			ConsumedCapacity: &dynamodb.ConsumedCapacity{},
			Item:             map[string]*dynamodb.AttributeValue{},
		}, nil
	}
}

func (svc *MockDynamoClient) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	result := dynamodb.PutItemOutput{}

	if input.Item == nil {
		return &result, errors.New("missing required field PutItemInput.Item")
	}

	if input.TableName == nil || *input.TableName == "" {
		return &result, errors.New("missing required field CreateTableInput.TableName")
	}

	return &result, nil
}
