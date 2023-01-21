package dynamo

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type (
	MockDynamoClient struct{ dynamodbiface.DynamoDBAPI }
)

func (svc *MockDynamoClient) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	if v, ok := input.ExpressionAttributeNames["#GSI"]; ok {
		if *v == "player_id" {
			if input.ExpressionAttributeValues != nil {
				if _, ok := input.ExpressionAttributeValues[":GSI"]; ok {
					switch {
					case *input.ExpressionAttributeValues[":GSI"].N == "123456":
						return &dynamodb.QueryOutput{
							ConsumedCapacity: &dynamodb.ConsumedCapacity{},
							Count:            new(int64),
							Items: []map[string]*dynamodb.AttributeValue{
								{
									"badge": {
										BOOL: aws.Bool(true),
									},
									"player_id": {
										N: input.ExpressionAttributeValues[":GSI"].N,
									},
								},
							},
							LastEvaluatedKey: map[string]*dynamodb.AttributeValue{},
							ScannedCount:     new(int64),
						}, nil
					}

				}

			}
		}
	}
	return &dynamodb.QueryOutput{}, fmt.Errorf("Something went wrong")
}

// func queryByIndexN(svc dynamodbiface.DynamoDBAPI, table, index, fieldName, fieldValue string) (*dynamodb.QueryOutput, error) {
// 	if table == "" || index == "" || fieldName == "" || fieldValue == "" {
// 		return nil, nil
// 	}
// 	// https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/dynamodbattribute/
// 	// https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/#example_DynamoDB_Query_shared00
// 	// https://stackoverflow.com/questions/56569183/query-a-secondary-index-on-dynamodb-in-golang
// 	// https://github.com/influencemobile/im2dynamodb/blob/5131a4f9f18d2bdddf738578c0385cfec83f3e8a/lib/im2/dynamodb.rb#L169
// 	result, err := svc.Query(&dynamodb.QueryInput{
// 		// ProjectionExpression: *string, // To request just a few fields
// 		ExpressionAttributeNames: map[string]*string{
// 			"#GSI": aws.String(fieldName),
// 		},
// 		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
// 			":GSI": {
// 				N: aws.String(fieldValue),
// 			},
// 		},
// 		KeyConditionExpression: aws.String("#GSI = :GSI"),
// 		TableName:              aws.String(table),
// 		IndexName:              aws.String(index),
// 	})
// 	if err != nil {
// 		if aerr, ok := err.(awserr.Error); ok {
// 			switch aerr.Code() {
// 			case dynamodb.ErrCodeProvisionedThroughputExceededException:
// 				return nil, fmt.Errorf("%s, %s", dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
// 			case dynamodb.ErrCodeResourceNotFoundException:
// 				return nil, fmt.Errorf("%s, %s", dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
// 			case dynamodb.ErrCodeRequestLimitExceeded:
// 				return nil, fmt.Errorf("%s, %s", dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
// 			case dynamodb.ErrCodeInternalServerError:
// 				return nil, fmt.Errorf("%s, %s", dynamodb.ErrCodeInternalServerError, aerr.Error())
// 			default:
// 				return nil, fmt.Errorf("%s", aerr.Error())
// 			}
// 		}
// 		return nil, err
// 	}
// 	return result, nil
// }
