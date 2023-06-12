package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	im_dynamo "influencemobile.com/libs/dynamo"
)

// This function retrieves the latest unique agreements for a given jurisdiction from a DynamoDB table
func GetUniqueAgreementsForJurisdiction(dynamoDBClient dynamodbiface.DynamoDBAPI, jurisdiction string) ([]DynamoJurisdictionMapping, error) {
	var lastEvaluatedKey map[string]*dynamodb.AttributeValue = nil
	uniqueAgreements := make(map[string]bool)
	var jurisdictionMappings []DynamoJurisdictionMapping

	for {
		queryResult, err := im_dynamo.QueryByIndexS(dynamoDBClient, "gdpr", "gdpr-index", "jurisdiction", jurisdiction, lastEvaluatedKey, 300, im_dynamo.ORDER_DESCENDING)
		if err != nil {
			return nil, err
		}

		for _, item := range queryResult.Items {
			agreement := item["agreement"].S
			if _, ok := uniqueAgreements[*agreement]; !ok {
				// If agreement is not in the map, add it
				uniqueAgreements[*agreement] = true
				var jurisdictionMapping DynamoJurisdictionMapping
				if err := dynamodbattribute.UnmarshalMap(item, &jurisdictionMapping); err != nil {
					return nil, err
				}
				jurisdictionMappings = append(jurisdictionMappings, jurisdictionMapping)
			}
		}
		if len(queryResult.LastEvaluatedKey) == 0 {
			break
		}
		lastEvaluatedKey = queryResult.LastEvaluatedKey
	}

	return jurisdictionMappings, nil
}

// This function retrieves the latest unique agreements for a given player id from a DynamoDB table
func GetUniqueAgreementsForPlayerId(dynamoDBClient dynamodbiface.DynamoDBAPI, player_id string, jurisdiction string) ([]DynamoPlayerIdMapping, error) {
	var lastEvaluatedKey map[string]*dynamodb.AttributeValue = nil
	uniqueAgreements := make(map[string]bool)
	var playeridMappings []DynamoPlayerIdMapping

	for {
		queryResult, err := im_dynamo.QueryByIndexN2(dynamoDBClient, "Consent", "consent-index", "player_id", "jurisdiction", player_id, jurisdiction, lastEvaluatedKey, 300, im_dynamo.ORDER_DESCENDING)
		if err != nil {
			return nil, err
		}
		for _, item := range queryResult.Items {
			agreement := item["agreement"].S
			if _, ok := uniqueAgreements[*agreement]; !ok {
				// If agreement is not in the map, add it
				uniqueAgreements[*agreement] = true
				var playerIdMapping DynamoPlayerIdMapping
				if err := dynamodbattribute.UnmarshalMap(item, &playerIdMapping); err != nil {
					return nil, err
				}
				playeridMappings = append(playeridMappings, playerIdMapping)
			}
		}
		if len(queryResult.LastEvaluatedKey) == 0 {
			break
		}
		lastEvaluatedKey = queryResult.LastEvaluatedKey
	}

	return playeridMappings, nil
}

func QueryByIndexN2(dynamoDBClient dynamodbiface.DynamoDBAPI, s1, s2 string, attributes map[string]string, lastEvaluatedKey map[string]*dynamodb.AttributeValue, i int, orderBy im_dynamo.OrderBy) {
	panic("unimplemented")
}

func PutItem(dynamoSvc dynamodbiface.DynamoDBAPI, table string, agreeRecord AgreeRecord) error {

	dbValues, err := dynamodbattribute.MarshalMap(&agreeRecord)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      dbValues,
		TableName: aws.String(table),
	}

	_, err = dynamoSvc.PutItem(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return fmt.Errorf("%s, %s", dynamodb.ErrCodeConditionalCheckFailedException, aerr.Error())
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				return fmt.Errorf("%s, %s", dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				return fmt.Errorf("%s, %s", dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
				return fmt.Errorf("%s, %s", dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
			case dynamodb.ErrCodeTransactionConflictException:
				return fmt.Errorf("%s, %s", dynamodb.ErrCodeTransactionConflictException, aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				return fmt.Errorf("%s, %s", dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				return fmt.Errorf("%s, %s", dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				return fmt.Errorf("%s", aerr.Error())
			}
		}
		return err
	}

	return nil
}
