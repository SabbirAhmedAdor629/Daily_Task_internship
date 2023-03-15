package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var dynamo *dynamodb.DynamoDB

type Row struct {
	Type         string `json:"type"`
	Version      string `json:"Version"`
	FileType     string `json:"file_type"`
	Jurisdiction string `json:"Jurisdiction"`
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

func MyHandler(ct context.Context, event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	var res *events.APIGatewayProxyResponse

	locale := event.QueryStringParameters["locale"]
	timezone := event.QueryStringParameters["timezone"]
	PartitionKey := event.QueryStringParameters["PartitionKey"]
	SortKey := event.QueryStringParameters["SortKey"]
	fmt.Println("locale=" + locale)
	fmt.Println("timezone=" + timezone)
	fmt.Println("PartitionKey=" + PartitionKey)
	fmt.Println("SortKey=" + SortKey)

	fmt.Println(GetItem(PartitionKey, SortKey))

	if event.Path == "/hello" {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 0,
			Body:       "Hello World !",
		}
	} else {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "defult ed",
		}
	}

	return res, nil
}

func main() {
	lambda.Start(MyHandler)
}

// 	// Get the value of the "name" query parameter
// 	name := r.URL.Query().Get("name")
// 	// Get the value of the "age" query parameter
// 	age := r.URL.Query().Get("age")
// 	// Print the values of the query parameters
// 	fmt.Fprintf(w, "Name: %s\n", name)
// 	fmt.Fprintf(w, "Age: %s\n", age)
// }

// queryParams := event.QueryStringParameters
// 	log.Println("queryParams",queryParams)
// 	locale := queryParams["locale"]
// 	timezone := queryParams["timezone"]
// 	log.Println("locale = ", locale)
// 	log.Println("timezone =", timezone)

// locale := event.QueryStringParameters["locale"]
// timezone := event.QueryStringParameters["timezone"]
// log.Println("locale = ", locale)
// log.Println("timezone =", timezone)

// parsedURL, _ := url.ParseRequestURI(event.Path)
// query := parsedURL.Query()
// locale := query.Get("locale")
// timezone := query.Get("timezone")
// log.Println("locale = ", locale)
// log.Println("timezone =", timezone)

// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 	locale := r.URL.Query().Get("locale")
// 	log.Println("locale = ", locale)
// 	timezone := r.URL.Query().Get("timezone")
// 	log.Println("timezone =", timezone)
// })

//logger.Info("recieved event", zap.Any("method", event.HTTPMethod), zap.Any("path", event.Path), zap.Any("body", event.Body))
// var logger *zap.Logger
// func init() {
// 	l, _ := zap.NewProduction()
// 	logger = l
// }
