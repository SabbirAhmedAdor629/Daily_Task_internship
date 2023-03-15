package main

import (
    "context"
    "fmt"
    "time"

    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
)

func MyHandler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
    location, err := time.LoadLocation("Europe/London")
    if err != nil {
        return nil, err
    }

    timeZoneMap := map[string]string{
        "France":    "Europe/Paris",
        "Germany":   "Europe/Berlin",
        "Italy":     "Europe/Rome",
        "Netherlands": "Europe/Amsterdam",
        "Spain":     "Europe/Madrid",
        "Switzerland": "Europe/Zurich",
        "United Kingdom": "Europe/London",
    }

    for country, tz := range timeZoneMap {
        loc, err := time.LoadLocation(tz)
        if err != nil {
            return nil, err
        }
        if loc == location {
            fmt.Println(country, "is in Europe (timezone:", loc, ")")
        }
    }

    return &events.APIGatewayProxyResponse{
        StatusCode: 200,
        Body:       "Hello from Lambda!",
    }, nil
}

func main() {
    lambda.Start(MyHandler)
}
