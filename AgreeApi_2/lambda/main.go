package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	// "os"

	// "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	im_config "influencemobile.com/parameter_store"
)

var (
	dynamo         *dynamodb.DynamoDB
	tableName      = "Consent"
	tableName2     = "gdpr"
	httpClient     HttpClient
	parameterStore im_config.ParameterStoreIface
)

type ParameterStoreType struct {
	Name  string
	Value string
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type APIResponseBody string

type ResponseObject struct {
	Error int         `json:"error"`
	Data  []Agreement `json:"data"`
}

type QueryStringParams struct {
	Timezone string `json:"timezone"`
}

type AgreeRecord struct {
	Jurisdiction string `json:"jurisdiction"`
	Agreement    string `json:"agreement"`
	Version      int    `json:"version"`
	Locale       string `json:"locale"`
	Timestamp    string `json:"timestamp"`
	Guid         string `json:"guid,omitempty"`
	PlayerId     string `json:"player_id,omitempty"`
}

type DynamoPlayerIdMapping struct {
	Jurisdiction string `json:"jurisdiction"`
	Agreement    string `json:"agreement"`
	Version      int    `json:"version"`
	Locale       string `json:"locale"`
	Timestamp    string `json:"timestamp"`
	Guid         string `json:"guid"`
	PlayerId     string `json:"player_id"`
}

type DynamoJurisdictionMapping struct {
	Jurisdiction    string `json:"jurisdiction"`
	Agreement       string `json:"agreement"`
	Version         int    `json:"version"`
	AgreementLocale string `json:"agreement_locale"`
}

type ConsentData struct {
	Guid         string `json:"guid"`
	PlayerId     string `json:"player_id"`
	Agreement    string `json:"agreement"`
	Jurisdiction string `json:"jurisdiction"`
	Locale       string `json:"locale"`
	Timestamp    string `json:"timestamp"`
	Version      int    `json:"version"`
}

type JurisdictionData struct {
	GUID            string `json:"guid"`
	Jurisdiction    string `json:"jurisdiction"`
	Agreement       string `json:"agreement"`
	AgreementLocale string `json:"agreement_locale"`
	Version         int    `json:"version"`
}

type Agreement struct {
	Jurisdiction string `json:"jurisdiction"`
	Agreement    string `json:"agreement"`
	Version      int    `json:"version"`
	Locale       string `json:"locale"`
	URL          string `json:"url"`
}

func connectDynamo() *dynamodb.DynamoDB {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		panic(err)
	}
	return dynamodb.New(sess)
}

func init() {
	dynamo = connectDynamo()
}

func GetLattestVersionOfAgreementsFromJurisdictionTable(jurisdiction string) (latestVersionOfPP int, latestVersionOfTOU int, rows []map[string]interface{}, err error) {
	// Query the table
	input := &dynamodb.QueryInput{
		TableName: aws.String(tableName2),
		KeyConditions: map[string]*dynamodb.Condition{
			"jurisdiction": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(jurisdiction),
					},
				},
			},
		},
		IndexName: aws.String("gdpr-index"),
	}

	result, err := dynamo.Query(input)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("error querying jurisdiction table: %s", err)
	}

	var jurisdictionData []JurisdictionData
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &jurisdictionData)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("error unmarshalling jurisdiction table data: %s", err)
	}

	// Find the latest version of tou and pp agreements
	for _, jd := range jurisdictionData {

		if jd.Agreement == "tou" {
			if jd.Version > latestVersionOfTOU {
				latestVersionOfTOU = jd.Version
			}
		} else if jd.Agreement == "pp" {
			if jd.Version > latestVersionOfPP {
				latestVersionOfPP = jd.Version
			}
		}
	}

	// Get the rows with the latest version number of tou and pp agreements
	for _, jd := range jurisdictionData {
		if (jd.Agreement == "tou" && jd.Version == latestVersionOfTOU) || (jd.Agreement == "pp" && jd.Version == latestVersionOfPP) {
			row, err := dynamodbattribute.MarshalMap(jd)
			if err != nil {
				return 0, 0, nil, fmt.Errorf("error marshalling jurisdiction data: %s", err)
			}
			interfaceMap := make(map[string]interface{})
			for key, value := range row {
				interfaceMap[key] = value
			}
			rows = append(rows, interfaceMap)
		}
	}

	return latestVersionOfPP, latestVersionOfTOU, rows, nil
}

func GetLatestVersionOfAgreementsFromConsentTableByPlayerId(playerID string) (touLatestVersion int, ppLatestVersion int, err error) {
	// Query Consent table by player_id
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(tableName),
		KeyConditions: map[string]*dynamodb.Condition{
			"player_id": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						N: aws.String(playerID),
					},
				},
			},
		},
		IndexName: aws.String("consent-index"),
	}

	queryOutput, err := dynamo.Query(queryInput)
	if err != nil {
		return 0, 0, err
	}

	var items []ConsentData
	err = dynamodbattribute.UnmarshalListOfMaps(queryOutput.Items, &items)
	if err != nil {
		return 0, 0, err
	}

	if len(items) == 0 {
		return 0, 0, fmt.Errorf("Player id did not found")
	}

	// Find the latest version of tou and pp agreements
	for _, item := range items {
		if item.Agreement == "tou" {
			if item.Version > touLatestVersion {
				touLatestVersion = item.Version
			}
		} else if item.Agreement == "pp" {
			if item.Version > ppLatestVersion {
				ppLatestVersion = item.Version
			}
		}
	}

	return touLatestVersion, ppLatestVersion, nil
}

func SetConsentUrl(hostUrl string, jurisdiction string, language string, document_type string, version int) string {
	return fmt.Sprintf("%s/agreement/%s/%s/%s/%d/%s.html", hostUrl, jurisdiction, language, document_type, version, document_type)
}

func CreateAgreement(jurisdiction, agreement, locale string, version int) Agreement {
	return Agreement{
		Jurisdiction: jurisdiction,
		Agreement:    agreement,
		Version:      version,
		Locale:       locale,
		URL:          SetConsentUrl("fkldf", jurisdiction, locale, agreement, version),
	}
}

func getPlayerID(req events.APIGatewayProxyRequest) (string, error) {
	u, err := url.Parse(req.Path)
	if err != nil {
		return "", fmt.Errorf("unable to parse")
	}
	// Split the path by '/'
	pathParts := strings.Split(u.Path, "/")

	// Retrieve the player ID from the pathParts
	var playerID string
	if len(pathParts) > 0 {
		playerID = pathParts[len(pathParts)-1]
	}
	return playerID, nil
}

func validTimezonePrefix(timezone string, prefix string) bool {
	return strings.HasPrefix(strings.ToLower(timezone), strings.ToLower(prefix))
}

func GetIpAddress(xForwardedFor string) (string, error) {
	ips := strings.Split(xForwardedFor, ",")
	if ip := net.ParseIP(strings.TrimSpace(ips[0])); ip != nil {
		return ip.String(), nil
	}
	return "", fmt.Errorf("invalid IP address")
}

func DetermineJurisdiction(params *QueryStringParams) string {
	const (
		gdprJurisdiction     = "gdpr"
		europeTimezonePrefix = "europe/"
	)
	if validTimezonePrefix(params.Timezone, europeTimezonePrefix) {
		return gdprJurisdiction
	}
	return ""
}

// int(v3.ResponseSuccess) == 0
func ValidateQueryStringParams(request events.APIGatewayProxyRequest) (*QueryStringParams, int, error) {
	params, err := getQueryStringParams(request)
	if err != nil {
		fmt.Println("Error getQueryStringParams function: ", err)
	}
	if params != nil && params.Timezone == "" {
		fmt.Println("timezone required")
	}
	return params, 0, nil
}

func getQueryStringParams(request events.APIGatewayProxyRequest) (*QueryStringParams, error) {
	queryJsonData, err := json.Marshal(request.QueryStringParameters)
	if err != nil {
		fmt.Println("Error marshalling query string")
	}
	queryStringParams := &QueryStringParams{}
	if err := json.Unmarshal(queryJsonData, &queryStringParams); err != nil {
		fmt.Println("Error unmarshalling query string")
	}
	return queryStringParams, nil
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	playerID, _ := getPlayerID(req)
	fmt.Println("Player Id : ", playerID)

	params, _, _ := ValidateQueryStringParams(req)
	fmt.Println("params : ", params)

	jurisdiction := DetermineJurisdiction(params)
	fmt.Println("Jurisdiction : ", jurisdiction)

	if jurisdiction != "gdpr" || playerID == "" {
		return *GenerateApiResponse(0, []Agreement{}, 200, nil), nil
	}

	// getting lattest unique aggrements from jurisdictions table
	agreementsObjList, _ := GetUniqueAgreementsForJurisdiction(dynamo, jurisdiction)

	var agreements []Agreement
	for _, agr := range agreementsObjList {
		agreement := CreateAgreement(agr.Jurisdiction, agr.Agreement, agr.AgreementLocale, agr.Version)
		agreements = append(agreements, agreement) // lattest uniqe aggrements list are here now
	}

	// getting lattest uniqe aggrements of players
	agreementsObjListOfPlayer, _ := GetUniqueAgreementsForPlayerId(dynamo, playerID, jurisdiction)

	fmt.Println("agreementsObjListOfPlayer : ", agreementsObjListOfPlayer)
	if agreementsObjListOfPlayer == nil {
		return *GenerateApiResponse(0, []Agreement{}, 200, nil), nil
	}

	var responseAgreement []Agreement
	for _, agrplayer := range agreementsObjListOfPlayer {
		for _, agr := range agreements {
			if agr.Agreement == agrplayer.Agreement {
				if agr.Version > agrplayer.Version {
					agreement := CreateAgreement(agr.Jurisdiction, agr.Agreement, agr.Locale, agr.Version)
					responseAgreement = append(responseAgreement, agreement)
				}
			}
		}
	}
	return *GenerateApiResponse(0, responseAgreement, 200, nil), nil
}

func main() {
	// request := events.APIGatewayProxyRequest{
	// 	Resource:          "",
	// 	Path:              "https://pgnzdf1s4m.execute-api.us-east-1.amazonaws.com/beta/12001",
	// 	HTTPMethod:        "",
	// 	Headers:           map[string]string{},
	// 	MultiValueHeaders: map[string][]string{},
	// 	QueryStringParameters: map[string]string{
	// 		"locale":   "en-uk",
	// 		"timezone": "Europe/London",
	// 	},
	// 	MultiValueQueryStringParameters: map[string][]string{},
	// 	PathParameters:                  map[string]string{},
	// 	StageVariables:                  map[string]string{},
	// 	RequestContext:                  events.APIGatewayProxyRequestContext{},
	// 	Body:                            "",
	// 	IsBase64Encoded:                 false,
	// }

	// params, _, _ := ValidateQueryStringParams(request)
	// fmt.Println(params)
	// jurisdiction := DetermineJurisdiction(params)
	// fmt.Println(jurisdiction)
	// playerID := getPlayerID(request)
	// fmt.Println("Player ID:", playerID)
	// agreementsObjList, _ := GetUniqueAgreementsForJurisdiction(dynamo, jurisdiction)
	// fmt.Println(agreementsObjList)
	// agreementsObjListOfPlayer, _ := GetUniqueAgreementsForPlayerId(dynamo, playerID, jurisdiction)
	// fmt.Println(agreementsObjListOfPlayer)

	// var agreements []Agreement
	// for _, agr := range agreementsObjList {
	// 	agreement := CreateAgreement(agr.Jurisdiction, agr.Agreement, agr.AgreementLocale, agr.Version)
	// 	agreements = append(agreements, agreement) // lattest uniqe aggrements list are here now
	// }

	// var responseAgreement []Agreement
	// for _, agrplayer := range agreementsObjListOfPlayer {
	// 	for _, agr := range agreements {
	// 		if agr.Agreement == agrplayer.Agreement {
	// 			if agr.Version > agrplayer.Version {
	// 				agreement := CreateAgreement(agr.Jurisdiction, agr.Agreement, agr.Locale, agr.Version)
	// 				responseAgreement = append(responseAgreement, agreement)
	// 			}
	// 		}
	// 	}
	// }
	// fmt.Println(*GenerateApiResponse(0, responseAgreement, 200, nil), nil)

	lambda.Start(handler)
}

// zip main.zip main

// go get github.com/aws/aws-lambda-go/lambda

// GOOS=linux GOARCH=amd64 go build -o main main.go

// GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go
