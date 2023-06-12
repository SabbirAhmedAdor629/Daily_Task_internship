package api_key

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type MockDynamoClient struct{ dynamodbiface.DynamoDBAPI }
type ApiKeyType map[string]int

func LoadCache(keys []string) ApiKeyType {
	apiKeys := make(ApiKeyType)
	if keys != nil {
		for _, v := range keys {
			apiKeys[v]++
		}
	}
	return apiKeys
}

func (a *ApiKeyType) ValidateKey(svc dynamodbiface.DynamoDBAPI, tableName, keyName string, apiKey string) (bool, error) {
	if tableName == "" || keyName == "" || apiKey == "" {
		return false, nil
	}

	// Check cache, return true if found.
	keyCache := *a
	if _, found := keyCache[apiKey]; found {
		return true, nil
	}

	// Check DynamoDB.
	authKey, err := getAuth(svc, tableName, keyName, apiKey)
	if err != nil { // Cannot connect/query database.
		return false, err
	}
	// Return false if not found.
	if authKey == nil {
		return false, nil
	}
	// fmt.Printf("%+v\n", *authKey)

	// Found. Add to cache.
	keyCache[authKey.ApiKey]++
	// fmt.Printf("%+v\n", keyCache)
	return true, nil
}

type Auth struct {
	ApiKey             string `json:"api_key"`
	Service            string `json:"-"` // `json:"service"`
	DynamoLastSyncedAt string `json:"-"` // `json:"dynamo_last_synced_at"`
}

func getAuth(svc dynamodbiface.DynamoDBAPI, keyTableName, keyName, apiKey string) (*Auth, error) {
	auth := Auth{}

	if keyTableName == "" || keyName == "" || apiKey == "" {
		return nil, nil
	}

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			keyName: {
				S: aws.String(apiKey),
			}},
		TableName: aws.String(keyTableName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				return nil, fmt.Errorf("%s, %s", dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				return nil, fmt.Errorf("%s, %s", dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				return nil, fmt.Errorf("%s, %s", dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				return nil, fmt.Errorf("%s, %s", dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				return nil, fmt.Errorf("%s", aerr.Error())
			}
		}
		return nil, err
	}

	// No record found
	if result.Item == nil {
		return nil, nil
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &auth)
	if err != nil {
		return nil, err
	}

	return &auth, nil
}
