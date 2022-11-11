package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
        Name string `json:"name"`
}

func HandleRequest(ctx context.Context, name MyEvent) {
		name.Name = "Sabbir"
        log.Println(fmt.Sprintf("Hello %s!", name.Name))
	}

func main() {
        lambda.Start(HandleRequest)
}