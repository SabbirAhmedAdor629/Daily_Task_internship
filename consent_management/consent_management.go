package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/honeybadger-io/honeybadger-go"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	im_keycache "influencemobile.com/libs/api_key"
	im_badger "influencemobile.com/libs/badger"
	v3 "influencemobile.com/libs/v3_helpers"
	logging "influencemobile.com/logging"
	im_config "influencemobile.com/parameter_store"
)

var (
	logger         *logging.Loggers
	env            EnvVarsStruct
	dynamoSvc      dynamodbiface.DynamoDBAPI
	opCounters     *OperationalCounters
	apiKeys        im_keycache.ApiKeyType
	badger         im_badger.BadgerIface
	parameterStore im_config.ParameterStoreIface
	httpClient     HttpClient
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

const (
	REQUIRED_AGREEMENTS_API             string = "/api/v1/required_agreements"
	REQUIRED_AGREEMENTS_WITH_PLAYER_API string = "/api/v1/required_agreements/{player_id}"
	AGREE_API                           string = "/api/v1/agree/{player_id}"
)

var (
	validApiPaths = []*regexp.Regexp{
		regexp.MustCompile("^/api/v1/required_agreements$"),
		regexp.MustCompile("^/api/v1/required_agreements/(?P<player_id>[^/]+)$"),
		regexp.MustCompile("^/api/v1/agree/(?P<player_id>[^/]+)$"),
	}

	expectedAllRequestHeaders = []string{"api_key", "request_id", "correlation_id"}

	endpointsMandatoryHeaders = map[string][]string{
		REQUIRED_AGREEMENTS_API:             {"api_key", "request_id", "correlation_id"},
		REQUIRED_AGREEMENTS_WITH_PLAYER_API: {"api_key", "request_id", "correlation_id"},
		AGREE_API:                           {"api_key", "request_id"},
	}

	mandatoryBodyParameters = []string{
		"jurisdiction",
		"agreement",
		"version",
		"locale",
		"timestamp",
	}
)

type ParameterStoreType struct {
	Name  string
	Value string
}

type (
	dynamoClient      struct{ dynamodbiface.DynamoDBAPI }
	QueryStringParams struct {
		Timezone string `json:"timezone"`
	}
)

type APIResponseBody string

type ResponseObject struct {
	Error int         `json:"error"`
	Data  []Agreement `json:"data"`
}

type Agreement struct {
	Jurisdiction string `json:"jurisdiction"`
	Agreement    string `json:"agreement"`
	Version      int    `json:"version"`
	Locale       string `json:"locale"`
	URL          string `json:"url"`
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

// MY CODE STARTS HERE
type DynamoPlayerIdMapping struct {
	Jurisdiction string `json:"jurisdiction"`
	Agreement    string `json:"agreement"`
	Version      int    `json:"version"`
	Locale       string `json:"locale"`
	Timestamp    string `json:"timestamp"`
	Guid         string `json:"guid"`
	PlayerId     string `json:"player_id"`
}

// MY CODE ENDS HERE
type DynamoJurisdictionMapping struct {
	Jurisdiction    string `json:"jurisdiction"`
	Agreement       string `json:"agreement"`
	Version         int    `json:"version"`
	AgreementLocale string `json:"agreement_locale"`
}

func GetIpAddress(xForwardedFor string) (string, error) {
	ips := strings.Split(xForwardedFor, ",")
	if ip := net.ParseIP(strings.TrimSpace(ips[0])); ip != nil {
		return ip.String(), nil
	}
	return "", fmt.Errorf("invalid IP address")
}

func ValidateRequestPath(request events.APIGatewayProxyRequest) (int, error) {
	for _, validPath := range validApiPaths {
		if path := validPath.FindStringSubmatch(request.Path); path != nil {
			return int(v3.ResponseSuccess), nil
		}
	}
	return int(v3.ResponseInvalidRequestPath), fmt.Errorf("invalid path")
}

func ValidateQueryStringParams(request events.APIGatewayProxyRequest) (*QueryStringParams, int, error) {
	params, err := getQueryStringParams(request)
	if err != nil {
		return params, int(v3.ResponseInternalServerError), err
	}

	if params != nil && params.Timezone == "" {
		return params, int(ResponseInvalidTimezone), fmt.Errorf("timezone required")
	}

	return params, int(v3.ResponseSuccess), nil
}

func getQueryStringParams(request events.APIGatewayProxyRequest) (*QueryStringParams, error) {
	queryJsonData, err := json.Marshal(request.QueryStringParameters)
	if err != nil {
		opCounters.incUnexpectedMarshalError(1)
		logger.NewEntry(log_msg_json_marshal_error, LogMsgTxt).WithStruct(opCounters).Error(err.Error())
		return nil, fmt.Errorf("cannot retrieve parameters: " + err.Error())
	}

	queryStringParams := &QueryStringParams{}
	if err := json.Unmarshal(queryJsonData, &queryStringParams); err != nil {
		opCounters.incUnexpectedUnmarshalError(1)
		logger.NewEntry(log_msg_json_unmarshal_error, LogMsgTxt).Error(err.Error())
		return nil, fmt.Errorf("cannot retrieve parameters: " + err.Error())
	}
	return queryStringParams, nil
}

func validTimezonePrefix(timezone string, prefix string) bool {
	return strings.HasPrefix(strings.ToLower(timezone), strings.ToLower(prefix))
}

func DetermineJurisdiction(params *QueryStringParams, ipqsKey string, headers map[string]string) (string, int, error) {
	const (
		gdprJurisdiction     = "gdpr"
		europeTimezonePrefix = "europe/"
	)

	if params == nil {
		sourceIp, err := GetIpAddress(headers[headerForwardedFor])
		if err != nil {
			return "", int(ResponseInvalidIp), err
		}
		ipResp, code, err := GetTimezoneFromIPQS(ipqsKey, sourceIp)
		if err != nil {
			return "", code, err
		}
		if validTimezonePrefix(ipResp.Timezone, europeTimezonePrefix) {
			return gdprJurisdiction, int(v3.ResponseSuccess), nil
		}
	} else if validTimezonePrefix(params.Timezone, europeTimezonePrefix) {
		return gdprJurisdiction, int(v3.ResponseSuccess), nil
	}
	return "", int(v3.ResponseSuccess), nil
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
		URL:          SetConsentUrl(env.CloudfrontHostUrl, jurisdiction, locale, agreement, version),
	}
}

func ValidateBodyParams(requestBody string, mandatoryParams []string, currentTime time.Time) ([]AgreeRecord, int, error) {
	agreeRecords := []AgreeRecord{}
	if err := json.Unmarshal([]byte(requestBody), &agreeRecords); err != nil {
		opCounters.incUnexpectedUnmarshalError(1)
		logger.NewEntry(log_msg_json_unmarshal_error, LogMsgTxt).Error(err.Error())
		return nil, int(v3.ResponseInvalidRequestBodyParameters), fmt.Errorf("cannot retrieve body parameters: " + err.Error())
	}
	for _, agreeRecord := range agreeRecords {
		if valid, errorCode, err := ValidateJurisdiction(agreeRecord); !valid {
			return nil, errorCode, err
		}
		if valid, errorCode, err := ValidateAgreement(agreeRecord); !valid {
			return nil, errorCode, err
		}
		if valid, errorCode, err := ValidateVersion(agreeRecord); !valid {
			return nil, errorCode, err
		}
		if valid, errorCode, err := ValidateAgreementLocale(agreeRecord); !valid {
			return nil, errorCode, err
		}
		if valid, errorCode, err := ValidateTimestamp(agreeRecord, currentTime); !valid {
			return nil, errorCode, err
		}
	}

	return agreeRecords, int(v3.ResponseSuccess), nil
}

func ValidateJurisdiction(agreeRecord AgreeRecord) (bool, int, error) {
	for _, value := range env.AllowedJurisdictionList {
		if agreeRecord.Jurisdiction == value {
			return true, int(v3.ResponseSuccess), nil
		}
	}

	return false, int(ResponseInvalidJurisdiction), fmt.Errorf("invalid jurisdiction")
}

func ValidateAgreement(agreeRecord AgreeRecord) (bool, int, error) {
	if ok := AgreementTypeIsAvailable(agreeRecord.Agreement, env.AllowedAgreementList); !ok {
		return false, int(ResponseInvalidAgreement), fmt.Errorf("invalid agreement")
	}

	return true, int(v3.ResponseSuccess), nil
}

func AgreementTypeIsAvailable(agreement string, agreementsList []string) bool {
	for _, value := range agreementsList {
		if value == agreement {
			return true
		}
	}
	return false
}

func ValidateVersion(agreeRecord AgreeRecord) (bool, int, error) {
	if agreeRecord.Version <= 0 {
		return false, int(ResponseInvalidVersion), fmt.Errorf("invalid version")
	}

	versionTime := time.Unix(int64(agreeRecord.Version), 0)
	if versionTime.IsZero() {
		return false, int(ResponseInvalidVersion), fmt.Errorf("invalid version")
	}

	return true, int(v3.ResponseSuccess), nil
}

func ValidateAgreementLocale(agreeRecord AgreeRecord) (bool, int, error) {
	for _, locale := range env.AllowedLocaleList {
		if agreeRecord.Locale == locale {
			return true, int(v3.ResponseSuccess), nil
		}
	}

	return false, int(ResponseInvalidLocale), fmt.Errorf("invalid locale")
}

func ValidateTimestamp(agreeRecord AgreeRecord, currentTime time.Time) (bool, int, error) {
	valid, err := IsTimeWithinThreshold(currentTime, agreeRecord.Timestamp, env.TimestampLimitSeconds)
	if err != nil {
		return false, int(ResponseInvalidTimestamp), fmt.Errorf("invalid timestamp - " + err.Error())
	}

	if !valid {
		return false, int(ResponseInvalidTimestamp), fmt.Errorf("invalid timestamp")
	}

	return true, int(v3.ResponseSuccess), nil
}

func IsTimeWithinThreshold(currentTime time.Time, givenTimestamp string, timeRange int) (bool, error) {
	givenTime, err := time.Parse(DEFAULT_TIME_FORMAT, givenTimestamp)
	if err != nil {
		return false, err
	}

	diff := currentTime.Sub(givenTime)
	// Check if the time difference is within current time and time range
	return (givenTime.Before(currentTime) || givenTime.Equal(currentTime)) && diff <= (time.Second*time.Duration(timeRange)), nil
}

func ConsentGUID(guid uuid.UUID) string {
	return fmt.Sprintf("agree-%s", guid)
}

func StoreAgreements(playerId v3.PlayerIdType, params []AgreeRecord) (int, error) {
	for _, item := range params {
		item.Guid = ConsentGUID(uuid.New())
		item.PlayerId = fmt.Sprintf("%d", playerId)
		err := PutItem(dynamoSvc, env.DynamoConsentTableName, item)
		if err != nil {
			return int(v3.ResponseInternalServerError), err
		}
	}
	return int(v3.ResponseSuccess), nil
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger.
		NewEntry(log_msg_begin, LogMsgTxt).
		Debug()
	logger. // Log environment parameters received from AWS
		NewEntry(log_msg_envars, LogMsgTxt).
		WithStruct(env).
		Debug()

	logger.
		NewEntry(log_msg_request_headers, LogMsgTxt).
		Info(request.Headers)

	logger. // Log lambda parameters
		NewEntry(log_msg_request, LogMsgTxt).
		WithStruct(request).
		Debug()

	opCounters = &OperationalCounters{}

	// Attempt to load the parameters file written by the lambda extension.
	if err := parameterStore.ReadConfig(); err != nil {
		logger.
			NewEntry(log_msg_parameters_store_parse_error, LogMsgTxt).
			Error(err.Error())
		os.Exit(1)
	}

	badgerParams, err := parameterStore.ParameterToBytes(DEFAULT_HONEYBADGER_KEY)
	if err != nil {
		logger.
			NewEntry(log_msg_parameters_store, LogMsgTxt).
			Error(err.Error())
		os.Exit(1)
	}

	badgerApiKey := ParameterStoreType{}
	_ = json.Unmarshal(badgerParams, &badgerApiKey)

	ipqsParams, err := parameterStore.ParameterToBytes(DEFAULT_IPQS_KEY)
	if err != nil {
		logger.
			NewEntry(log_msg_parameters_store, LogMsgTxt).
			Error(err.Error())
		os.Exit(1)
	}

	ipqsApiKey := ParameterStoreType{}
	_ = json.Unmarshal(ipqsParams, &ipqsApiKey)

	badger.Configure(honeybadger.Configuration{APIKey: badgerApiKey.Value, Env: env.Environment})

	responseHeaders := v3.CreateResponseHeaders(request.Headers, expectedAllRequestHeaders)

	// Validate required headers
	_, code, err := v3.ValidateRequestHeaders(responseHeaders, endpointsMandatoryHeaders[request.Resource], v3.ResponseCodeToKey)
	if err != nil {
		opCounters.incInvalidRequestHeaders(1)
		logger.NewEntry(log_msg_invalid_headers, LogMsgTxt).WithStruct(opCounters).Error(err.Error())
		return *GenerateApiResponse(code, nil, http.StatusOK, responseHeaders), nil
	}

	// validate request path
	code, err = ValidateRequestPath(request)
	if err != nil {
		opCounters.incInvalidPath(1)
		logger.NewEntry(log_msg_invalid_path, LogMsgTxt).WithStruct(opCounters).Error(err.Error())
		return *GenerateApiResponse(code, nil, http.StatusOK, responseHeaders), nil
	}

	// Validate the api key.
	var authorized bool
	authorized, code, err = v3.ValidateAuth(dynamoSvc, &apiKeys, env.DynamoApiKeysTableName, responseHeaders["api_key"])
	if err != nil {
		opCounters.incDynamoDbError(1)
		logger.NewEntry(log_msg_dynamodb_error, LogMsgTxt).WithStruct(opCounters).Error(err.Error())
		badger.Notify(err.Error(),
			honeybadger.Context{
				"ERROR_NUMBER": log_msg_dynamodb_error,
				"ERROR_TEXT":   LogMsgTxt[log_msg_dynamodb_error],
			},
			honeybadger.Tags{
				env.ApplicationName,
				env.Environment,
				LogMsgTxt[log_msg_dynamodb_error]})
		badger.Flush()
		return *GenerateApiResponse(code, nil, http.StatusOK, responseHeaders), nil
	}
	if !authorized {
		opCounters.incInvalidAuth(1)
		logger.NewEntry(log_msg_invalid_auth, LogMsgTxt).WithStruct(opCounters).Info()
		return *GenerateApiResponse(code, nil, http.StatusOK, responseHeaders), nil
	}

	switch request.HTTPMethod {
	case "GET":
		params, code, err := ValidateQueryStringParams(request)
		if err != nil {
			opCounters.incInvalidQueryParameters(1)
			logger.NewEntry(log_msg_invalid_query_params, LogMsgTxt).WithStruct(opCounters).Error(err.Error())
			return *GenerateApiResponse(code, nil, http.StatusOK, responseHeaders), nil
		}

		switch request.Resource {
		case REQUIRED_AGREEMENTS_API:
			jurisdiction, code, err := DetermineJurisdiction(params, ipqsApiKey.Value, request.Headers)
			if err != nil {
				logger.NewEntry(log_msg_jurisdiction_error, LogMsgTxt).WithStruct(opCounters).Error(err.Error())
				return *GenerateApiResponse(code, nil, http.StatusOK, responseHeaders), nil
			}

			if jurisdiction != "" {
				agreementsObjList, err := GetUniqueAgreementsForJurisdiction(dynamoSvc, jurisdiction)
				if err != nil {
					opCounters.incDynamoDbError(1)
					logger.NewEntry(log_msg_dynamodb_error, LogMsgTxt).WithStruct(opCounters).Error(err.Error())
					badger.Notify(err.Error(),
						honeybadger.Context{
							"ERROR_NUMBER": log_msg_dynamodb_error,
							"ERROR_TEXT":   LogMsgTxt[log_msg_dynamodb_error],
						},
						honeybadger.Tags{
							env.ApplicationName,
							env.Environment,
							LogMsgTxt[log_msg_dynamodb_error]})
					badger.Flush()
					return *GenerateApiResponse(int(v3.ResponseInternalServerError), nil, http.StatusOK, responseHeaders), nil
				}

				var agreements []Agreement
				for _, agr := range agreementsObjList {
					agreement := CreateAgreement(agr.Jurisdiction, agr.Agreement, agr.AgreementLocale, agr.Version)
					agreements = append(agreements, agreement) // lattest uniqe aggrements list are here now
				}

				logger.NewEntry(log_msg_end, LogMsgTxt).WithStruct(opCounters).Info()
				return *GenerateApiResponse(int(v3.ResponseSuccess), agreements, http.StatusOK, responseHeaders), nil
			} else {
				return *GenerateApiResponse(int(v3.ResponseSuccess), []Agreement{}, http.StatusOK, responseHeaders), nil
			}
		case REQUIRED_AGREEMENTS_WITH_PLAYER_API:
			// My code starts here
			playerId := request.PathParameters["player_id"]

			jurisdiction, code, err := DetermineJurisdiction(params, ipqsApiKey.Value, request.Headers)
			if err != nil {
				logger.NewEntry(log_msg_jurisdiction_error, LogMsgTxt).WithStruct(opCounters).Error(err.Error())
				return *GenerateApiResponse(code, nil, http.StatusOK, responseHeaders), nil
			}
			if jurisdiction != "" && playerId != "" {
				agreementsObjList, err := GetUniqueAgreementsForJurisdiction(dynamoSvc, jurisdiction)
				if err != nil {
					opCounters.incDynamoDbError(1)
					logger.NewEntry(log_msg_dynamodb_error, LogMsgTxt).WithStruct(opCounters).Error(err.Error())
					badger.Notify(err.Error(),
						honeybadger.Context{
							"ERROR_NUMBER": log_msg_dynamodb_error,
							"ERROR_TEXT":   LogMsgTxt[log_msg_dynamodb_error],
						},
						honeybadger.Tags{
							env.ApplicationName,
							env.Environment,
							LogMsgTxt[log_msg_dynamodb_error]})
					badger.Flush()
					return *GenerateApiResponse(int(v3.ResponseInternalServerError), nil, http.StatusOK, responseHeaders), nil
				}
				var agreements []Agreement
				for _, agr := range agreementsObjList {
					agreement := CreateAgreement(agr.Jurisdiction, agr.Agreement, agr.AgreementLocale, agr.Version)
					agreements = append(agreements, agreement) // lattest uniqe aggrements list are here now
				}
				agreementsObjListOfPlayer, err := GetUniqueAgreementsForPlayerId(dynamoSvc, playerId)
				if err != nil {
					opCounters.incDynamoDbError(1)
					logger.NewEntry(log_msg_dynamodb_error, LogMsgTxt).WithStruct(opCounters).Error(err.Error())
					badger.Notify(err.Error(),
						honeybadger.Context{
							"ERROR_NUMBER": log_msg_dynamodb_error,
							"ERROR_TEXT":   LogMsgTxt[log_msg_dynamodb_error],
						},
						honeybadger.Tags{
							env.ApplicationName,
							env.Environment,
							LogMsgTxt[log_msg_dynamodb_error]})
					badger.Flush()
					return *GenerateApiResponse(int(v3.ResponseInternalServerError), nil, http.StatusOK, responseHeaders), nil
				}

				for _, agrplayer := range agreementsObjListOfPlayer {
					for _, agr := range agreements {
						if agr.Agreement == agrplayer.Agreement {
							if agr.Version == agrplayer.Version {
								return *GenerateApiResponse(int(v3.ResponseSuccess), []Agreement{}, http.StatusOK, responseHeaders), nil
							} else if agr.Version > agrplayer.Version {
								return *GenerateApiResponse(int(v3.ResponseSuccess), agreements, http.StatusOK, responseHeaders), nil
							}
						}
					}
				}
				logger.NewEntry(log_msg_end, LogMsgTxt).WithStruct(opCounters).Info()
			} else {
				return *GenerateApiResponse(code, nil, http.StatusOK, responseHeaders), nil
			}
			// My code ends here
		default:
			return *GenerateApiResponse(int(v3.ResponseInvalidRequestPath), nil, http.StatusOK, responseHeaders), nil
		}

	case "POST":
		switch request.Resource {
		case AGREE_API:
			requestPlayerId := request.PathParameters["player_id"]
			playerId, code, err := v3.ValidatePlayerId(&requestPlayerId)
			if err != nil {
				opCounters.incInvalidPlayerId(1)
				logger.NewEntry(log_msg_invalid_player_id, LogMsgTxt).WithStruct(opCounters).Error(err.Error())
				return *GenerateApiResponse(code, nil, http.StatusOK, responseHeaders), nil
			}

			params, code, err := ValidateBodyParams(request.Body, mandatoryBodyParameters, time.Now().UTC())
			if err != nil {
				opCounters.incInvalidBodyParameters(1)
				logger.NewEntry(log_msg_invalid_body_parameters, LogMsgTxt).WithStruct(opCounters).Error(err.Error())
				return *GenerateApiResponse(code, nil, http.StatusOK, responseHeaders), nil
			}

			code, err = StoreAgreements(playerId, params)
			if err != nil {
				opCounters.incDynamoDbError(1)
				logger.NewEntry(log_msg_dynamodb_error, LogMsgTxt).WithStruct(opCounters).Error(err.Error())
				return *GenerateApiResponse(code, nil, http.StatusOK, responseHeaders), nil
			}

			return *GenerateApiResponse(int(v3.ResponseSuccess), nil, http.StatusOK, responseHeaders), nil
		default:
			return *GenerateApiResponse(int(v3.ResponseInvalidRequestPath), nil, http.StatusOK, responseHeaders), nil
		}
	}

	return *GenerateApiResponse(int(v3.ResponseInvalidRequestMethod), nil, http.StatusOK, responseHeaders), nil
}

func init() {

	env = NewEnv()

	// Create a new logger.
	logger = logging.NewLoggers(os.Stderr, env.LogLevel, 0).
		SetFields(logging.Fields{
			EnvVars[ENV_APPLICATION_NAME]: env.ApplicationName,
		})

	apiKeys = im_keycache.LoadCache(im_keycache.GetApiKeys().ApiKeys)
	s, err := session.NewSession()
	if err != nil {
		fmt.Printf("error creating sessionL: %s\n", err.Error())
		os.Exit(1)

	}

	dynamoSvc = dynamoClient{dynamodb.New(s)}
	parameterStore = &im_config.KeyClient{Filename: env.ParametersStore}
	badger = &im_badger.BadgerClient{}
	httpClient = &http.Client{}
}

func main() {
	defer badger.Monitor()
	lambda.Start(HandleRequest)
}
