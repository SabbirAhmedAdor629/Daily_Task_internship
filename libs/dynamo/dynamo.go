package dynamo

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type (
	SortBy           int
	ContinuationType map[string]*dynamodb.AttributeValue
	OrderBy          int
	QueryParamType   int
)

const (
	SORT_CREATED_AT SortBy = iota

	ORDER_ASCENDING OrderBy = iota
	ORDER_DESCENDING
)

func QueryByIndexN(svc dynamodbiface.DynamoDBAPI, table, index, fieldName, fieldValue string,
	page ContinuationType, pageCount int, orderBy OrderBy) (*dynamodb.QueryOutput, error) {
	if table == "" || index == "" || fieldName == "" || fieldValue == "" {
		return nil, nil
	}

	var scanForward bool = false // Descending sort order, newest to oldest.
	if orderBy == ORDER_ASCENDING {
		scanForward = true // Ascending sort order, oldest to newest.
	}

	// https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/dynamodbattribute/
	// https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/#example_DynamoDB_Query_shared00
	// https://stackoverflow.com/questions/56569183/query-a-secondary-index-on-dynamodb-in-golang
	// https://github.com/influencemobile/im2dynamodb/blob/5131a4f9f18d2bdddf738578c0385cfec83f3e8a/lib/im2/dynamodb.rb#L169
	result, err := svc.Query(&dynamodb.QueryInput{
		// ProjectionExpression: *string, // To request just a few fields
		ExpressionAttributeNames: map[string]*string{
			"#GSI": aws.String(fieldName),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":GSI": {
				N: aws.String(fieldValue),
			},
		},
		KeyConditionExpression: aws.String("#GSI = :GSI"),
		TableName:              aws.String(table),
		IndexName:              aws.String(index),
		ExclusiveStartKey:      page,
		Limit:                  aws.Int64(int64(pageCount)),
		ScanIndexForward:       aws.Bool(scanForward),
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
	return result, nil
}

func QueryByIndexS(svc dynamodbiface.DynamoDBAPI, table, index, fieldName, fieldValue string,
	page ContinuationType, pageCount int, orderBy OrderBy) (*dynamodb.QueryOutput, error) {
	if table == "" || index == "" || fieldName == "" || fieldValue == "" {
		return nil, nil
	}

	var scanForward bool = false // Descending sort order, newest to oldest.
	if orderBy == ORDER_ASCENDING {
		scanForward = true // Ascending sort order, oldest to newest.
	}

	// https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/dynamodbattribute/
	// https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/#example_DynamoDB_Query_shared00
	// https://stackoverflow.com/questions/56569183/query-a-secondary-index-on-dynamodb-in-golang
	// https://github.com/influencemobile/im2dynamodb/blob/5131a4f9f18d2bdddf738578c0385cfec83f3e8a/lib/im2/dynamodb.rb#L169
	result, err := svc.Query(&dynamodb.QueryInput{
		// ProjectionExpression: *string, // To request just a few fields
		ExpressionAttributeNames: map[string]*string{
			"#GSI": aws.String(fieldName),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":GSI": {
				S: aws.String(fieldValue),
			},
		},
		KeyConditionExpression: aws.String("#GSI = :GSI"),
		TableName:              aws.String(table),
		IndexName:              aws.String(index),
		ExclusiveStartKey:      page,
		Limit:                  aws.Int64(int64(pageCount)),
		ScanIndexForward:       aws.Bool(scanForward),
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
	return result, nil
}

func QueryByAppId(svc dynamodbiface.DynamoDBAPI, table, index string, appId string, adId string) (*dynamodb.QueryOutput, error) {
	if table == "" || index == "" || appId == "" || adId == "" {
		return nil, nil
	}
	// https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/dynamodbattribute/
	// https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/#example_DynamoDB_Query_shared00
	// https://stackoverflow.com/questions/56569183/query-a-secondary-index-on-dynamodb-in-golang
	// https://github.com/influencemobile/im2dynamodb/blob/5131a4f9f18d2bdddf738578c0385cfec83f3e8a/lib/im2/dynamodb.rb#L169
	result, err := svc.Query(&dynamodb.QueryInput{
		// ProjectionExpression: *string, // To request just a few fields
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":GSI": {
				S: aws.String(adId),
			},
			":ASI": {
				N: aws.String(appId),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#GSI": aws.String("advertising_id"),
			"#ASI": aws.String("app_id"),
		},
		KeyConditionExpression: aws.String("#GSI = :GSI"),
		FilterExpression:       aws.String("#ASI = :ASI"),
		TableName:              aws.String(table),
		IndexName:              aws.String(index),
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
	return result, nil
}

func QetAllItems(svc dynamodbiface.DynamoDBAPI, table, fieldName string, fieldValue int) (*[]dynamodb.QueryOutput, error) {
	if table == "" || fieldName == "" {
		return nil, nil
	}

	results := []dynamodb.QueryOutput{}
	var lastKey map[string]*dynamodb.AttributeValue = nil
	for {
		result, err := svc.Query(&dynamodb.QueryInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":q": {
					N: aws.String(strconv.Itoa(fieldValue)),
				},
			},
			ExpressionAttributeNames: map[string]*string{
				"#Q": aws.String(fieldName),
			},
			KeyConditionExpression: aws.String("#Q = :q"),
			TableName:              aws.String(table),
			ExclusiveStartKey:      lastKey,
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
		if result.Items != nil {
			results = append(results, *result)
			// Exit loop if there are no more results
			if result.LastEvaluatedKey == nil {
				break
			}
		}
		lastKey = result.LastEvaluatedKey
	}
	return &results, nil
}

func queryByMessageId(svc dynamodbiface.DynamoDBAPI, table, index string, playerIdKey, playerIdVal, memberIdKey, memberIdVal, messageIdKey, messageIdValue string) (*dynamodb.QueryOutput, error) {
	if table == "" || index == "" {
		return nil, nil
	}
	// https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/dynamodbattribute/
	// https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/#example_DynamoDB_Query_shared00
	// https://stackoverflow.com/questions/56569183/query-a-secondary-index-on-dynamodb-in-golang
	// https://github.com/influencemobile/im2dynamodb/blob/5131a4f9f18d2bdddf738578c0385cfec83f3e8a/lib/im2/dynamodb.rb#L169
	result, err := svc.Query(&dynamodb.QueryInput{
		// ProjectionExpression: *string, // To request just a few fields
		ExpressionAttributeNames: map[string]*string{
			"#GSI": aws.String(playerIdKey),
			"#ASI": aws.String(memberIdKey),
			"#BSI": aws.String(messageIdKey),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":GSI": {
				N: aws.String(playerIdVal),
			},
			":ASI": {
				N: aws.String(memberIdVal),
			},
			":BSI": {
				S: aws.String(messageIdValue),
			},
		},
		KeyConditionExpression: aws.String("#GSI = :GSI"),
		FilterExpression:       aws.String("#ASI = :ASI"),
		TableName:              aws.String(table),
		IndexName:              aws.String(index),
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

	return result, nil
}

func QueryByPlayerId(svc dynamodbiface.DynamoDBAPI, table, fieldName string, fieldValue int, indexName string) (*[]dynamodb.QueryOutput, error) {
	if table == "" || fieldName == "" {
		return nil, nil
	}

	results := []dynamodb.QueryOutput{}
	var lastKey map[string]*dynamodb.AttributeValue = nil
	for {
		result, err := svc.Query(&dynamodb.QueryInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":q": {
					N: aws.String(strconv.Itoa(fieldValue)),
				},
			},
			ExpressionAttributeNames: map[string]*string{
				"#Q": aws.String(fieldName),
			},
			KeyConditionExpression: aws.String("#Q = :q"),
			TableName:              aws.String(table),
			IndexName:              aws.String(indexName),
			ExclusiveStartKey:      lastKey,
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
		if result.Items != nil {
			results = append(results, *result)
			// Exit loop if there are no more results
			if result.LastEvaluatedKey == nil {
				break
			}
		}
		lastKey = result.LastEvaluatedKey
	}
	return &results, nil
}