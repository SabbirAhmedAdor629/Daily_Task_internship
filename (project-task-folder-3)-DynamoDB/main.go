package main

import (
	"fmt"
	// "strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var dynamo *dynamodb.DynamoDB

type Row struct {
	Type     string
	Version    string
	FileType    string 	`json:"file_type"`
	Jurisdiction string
}

const TABLE_NAME = "NewTable"

func init() {
	dynamo = connectDynamo()
}

// connectDynamo returns a dynamoDB client
func connectDynamo() (db *dynamodb.DynamoDB) {
	return dynamodb.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})))
}


// GetItem gets a Person based on the Id, returns a person
func GetItem(Type string, Version string) (rows Row, err error) {
    result, err := dynamo.GetItem(&dynamodb.GetItemInput{
        Key: map[string]*dynamodb.AttributeValue{
            "type": {
                S: aws.String(Type),
            },
            "Version": {
                S: aws.String(Version),
            },
        },
        TableName: aws.String(TABLE_NAME),
    })

    if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
            fmt.Println(aerr)
        }
    }
    err = dynamodbattribute.UnmarshalMap(result.Item, &rows)
    return rows, err
}


func main() {
    d := "en-uk-tou"
	fmt.Println(GetItem(d,"1.2.0"))
}