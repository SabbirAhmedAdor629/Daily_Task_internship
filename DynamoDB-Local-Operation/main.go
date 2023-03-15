package main

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)


var dynamo *dynamodb.DynamoDB



type Person struct {
	Id      int
	Name    string
	Website string
}

const TABLE_NAME = "people"

func init() {
	dynamo = connectDynamo()
}

// connectDynamo returns a dynamoDB client
func connectDynamo() (db *dynamodb.DynamoDB) {
	return dynamodb.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})))
}

// CreateTable creates a table
func CreateTable() {
	_, err := dynamo.CreateTable(&dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Id"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Id"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String(TABLE_NAME),

		// BillingMode: aws.String(dynamodb.BillingModePayPerRequest),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr)
		}
	}
}

// PutItem inserts the struct Person
func PutItem(person Person) {
	_, err := dynamo.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Id": {
				N: aws.String(strconv.Itoa(person.Id)),
			},
			"Name": {
				S: aws.String(person.Name),
			},
			"Website": {
				S: aws.String(person.Website),
			},
		},
		TableName: aws.String(TABLE_NAME),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr)
		}
	}
}

// UpdateItem updates the Person based on the Person.Id
func UpdateItem(person Person) {
	_, err := dynamo.UpdateItem(&dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#N": aws.String("Name"),
			"#W": aws.String("Website"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":Name": {
				S: aws.String(person.Name),
			},
			":Website": {
				S: aws.String(person.Website),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				N: aws.String(strconv.Itoa(person.Id)),
			},
		},
		TableName:        aws.String(TABLE_NAME),
		UpdateExpression: aws.String("SET #N = :Name, #W = :Website"),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr)
		}
	}
}

// DeleteItem deletes a Person based on the Person.Id
func DeleteItem(id int)  {
	_, err := dynamo.DeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				N: aws.String(strconv.Itoa(id)),
			},
		},
		TableName: aws.String(TABLE_NAME),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr)
		}
	}
}

// GetItem gets a Person based on the Id, returns a person
// func GetItem(id int) (All_Person Person, err error) {
// 	result, err := dynamo.GetItem(&dynamodb.GetItemInput{
// 		Key: map[string]*dynamodb.AttributeValue{
// 			"Id": {
// 				N: aws.String(strconv.Itoa(id)),
// 			},
// 		},
// 		TableName: aws.String(TABLE_NAME),
// 	})
// 	if err != nil {
// 		if aerr, ok := err.(awserr.Error); ok {
// 			fmt.Println(aerr)
// 		}
// 	}
// 	err = dynamodbattribute.UnmarshalMap(result.Item, &All_Person)
// 	return All_Person, err
// }


func GetItem(svc dynamodbiface.DynamoDBAPI, id int) (All_Person Person, err error) {
    result, err := svc.GetItem(&dynamodb.GetItemInput{
    	AttributesToGet:          []*string{},
    	ConsistentRead:           new(bool),
    	ExpressionAttributeNames: map[string]*string{},
    	Key:                      map[string]*dynamodb.AttributeValue{"Id": {N: aws.String(strconv.Itoa(id))}},
    	ProjectionExpression:     new(string),
    	ReturnConsumedCapacity:   new(string),
    	TableName:                aws.String(TABLE_NAME),
    })
    if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
            fmt.Println(aerr)
        }
        return All_Person, err
    }
    err = dynamodbattribute.UnmarshalMap(result.Item, &All_Person)
    return All_Person, err
}









// func GetItem(id int, svc *dynamodb.DynamoDB) (All_Person Person, err error) {
// 	result, err := svc.GetItem(&dynamodb.GetItemInput{
// 		Key: map[string]*dynamodb.AttributeValue{
// 			"Id": {
// 				N: aws.String(strconv.Itoa(id)),
// 			},
// 		},
// 		TableName: aws.String(TABLE_NAME),
// 	})

// 	if err != nil {
// 		if aerr, ok := err.(awserr.Error); ok {
// 			fmt.Println(aerr)
// 		}
// 	}

// 	err = dynamodbattribute.UnmarshalMap(result.Item, &All_Person)

// 	return All_Person, err
// }




func main() {
	// create a new DynamoDB client using an AWS session
	

	CreateTable()

	// INSERT
	var person_1 Person = Person{
		Id:      1,
		Name:    "Sabbir Ahmed",
		Website: "https://sabbirahmedador629.github.io/Personal_Portfolio/",
	}
	
	dynamoSvc := &MockDynamoClient{
		DynamoDBAPI: dynamo,
	}


	PutItem(person_1)
	fmt.Println(GetItem(dynamoSvc,1))


	// // UPDATE
	// person_1.Name = "Sabbir Ahmed Ador"
	// UpdateItem(person_1)
	

	// fmt.Println(GetItem(1))
	// // DELETE
	// DeleteItem(1)

	// fmt.Println(GetItem(1))
}






