package main

import (
	//"context"
	"log"

	//"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	//"go.uber.org/zap"
)




func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error){

	log.Println("hello world")

	response := events.APIGatewayProxyResponse{
		StatusCode:        0,
		Headers:           map[string]string{},
		MultiValueHeaders: map[string][]string{},
		Body:              "",
		IsBase64Encoded:   false,
	}

	return response, nil
	
}

func main() {
	lambda.Start(handler)
}