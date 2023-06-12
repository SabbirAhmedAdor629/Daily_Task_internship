package main

import (
	"os"
	"strconv"
	"strings"

	"influencemobile.com/logging"
)

type TypeEnvVar int

const (
	ENV_ENVIRONMENT TypeEnvVar = iota
	ENV_APPLICATION_NAME
	ENV_LOG_LEVEL
	ENV_AWS_REGION
	ENV_PARAMETERS_STORE
	ENV_DYNAMO_CONSENT_TABLE_NAME
	ENV_DYNAMO_MAPPING_TABLE_NAME
	ENV_CLOUDFRONT_HOST_URL
	ENV_DYNAMO_API_KEYS_TABLE_NAME
	ENV_TIMESTAMP_LIMIT_SECONDS
	ENV_ALLOWED_AGREEMENT_LIST
	ENV_ALLOWED_JURISDICTION_LIST
	ENV_ALLOWED_LOCALE_LIST
)

var (
	EnvVars = []string{
		ENV_ENVIRONMENT:                "TARGET_ENVIRONMENT",
		ENV_APPLICATION_NAME:           "APP_NAME",
		ENV_LOG_LEVEL:                  "LOG_LEVEL",
		ENV_AWS_REGION:                 "REGION",
		ENV_PARAMETERS_STORE:           "PARAMETERS_STORE",
		ENV_DYNAMO_CONSENT_TABLE_NAME:  "DYNAMO_CONSENT_TABLE_NAME",
		ENV_DYNAMO_MAPPING_TABLE_NAME:  "DYNAMO_MAPPING_TABLE_NAME",
		ENV_CLOUDFRONT_HOST_URL:        "CLOUDFRONT_HOST_URL",
		ENV_DYNAMO_API_KEYS_TABLE_NAME: "API_KEYS_TABLE_NAME",
		ENV_TIMESTAMP_LIMIT_SECONDS:    "TIMESTAMP_LIMIT_SECONDS",
		ENV_ALLOWED_AGREEMENT_LIST:     "ALLOWED_AGREEMENT_LIST",
		ENV_ALLOWED_JURISDICTION_LIST:  "ALLOWED_JURISDICTION_LIST",
		ENV_ALLOWED_LOCALE_LIST:        "ALLOWED_LOCALE_LIST",
	}
)

// Retrieve region, sqs name, kinesis firehose name from lambda environment values
type EnvVarsStruct struct {
	LogLevel                int      `json:"LOG_LEVEL"`
	Environment             string   `json:"TARGET_ENVIRONMENT"`
	ApplicationName         string   `json:"APP_NAME"`
	Region                  string   `json:"REGION"`
	ParametersStore         string   `json:"PARAMETERS_STORE"`
	DynamoConsentTableName  string   `json:"CONSENT_DYNAMO_TABLE_NAME"`
	DynamoMappingTableName  string   `json:"MAPPING_TABLE_NAME"`
	DynamoApiKeysTableName  string   `json:"API_KEYS_TABLE_NAME"`
	CloudfrontHostUrl       string   `json:"CLOUDFRONT_HOST_URL"`
	TimestampLimitSeconds   int      `json:"TIMESTAMP_LIMIT_SECONDS"`
	AllowedAgreementList    []string `json:"ALLOWED_AGREEMENT_LIST"`
	AllowedJurisdictionList []string `json:"ALLOWED_JURISDICTION_LIST"`
	AllowedLocaleList       []string `json:"ALLOWED_LOCALE_LIST"`
}

func NewEnv() EnvVarsStruct {
	return EnvVarsStruct{
		LogLevel:                getLogLevel(),
		Environment:             getEnvironment(),
		ApplicationName:         getApplicationName(),
		Region:                  getRegion(),
		ParametersStore:         getParametersStore(),
		DynamoConsentTableName:  getDynamoConsentTableName(),
		DynamoMappingTableName:  getDynamoMappingTableName(),
		DynamoApiKeysTableName:  getDynamoApiKeyTableName(),
		CloudfrontHostUrl:       getCloudfrontHostUrl(),
		TimestampLimitSeconds:   getTimestampLimitSeconds(),
		AllowedAgreementList:    getAllowedAgreementList(),
		AllowedJurisdictionList: getAllowedJurisdictionList(),
		AllowedLocaleList:       getAllowedLocaleList(),
	}
}

func getLogLevel() int {
	txtLevel := strings.ToUpper(os.Getenv(EnvVars[ENV_LOG_LEVEL]))
	level := DEFAULT_LOG_LEVEL
	for i, l := range logging.LogLevelText {
		if txtLevel == l {
			level = i
		}
	}
	return level
}

func getApplicationName() string {
	val, ok := os.LookupEnv(EnvVars[ENV_APPLICATION_NAME])
	if !ok {
		return DEFAULT_APPLICATION_NAME
	}
	return val
}

func getEnvironment() string     { return os.Getenv(EnvVars[ENV_ENVIRONMENT]) }
func getRegion() string          { return os.Getenv(EnvVars[ENV_AWS_REGION]) }
func getParametersStore() string { return os.Getenv(EnvVars[ENV_PARAMETERS_STORE]) }

func getDynamoConsentTableName() string {
	val, ok := os.LookupEnv(EnvVars[ENV_DYNAMO_CONSENT_TABLE_NAME])
	if !ok {
		return DEFAULT_DYNAMO_CONSENT_TABLE_NAME
	}
	return val
}

func getDynamoMappingTableName() string {
	val, ok := os.LookupEnv(EnvVars[ENV_DYNAMO_MAPPING_TABLE_NAME])
	if !ok {
		return DEFAULT_DYNAMO_MAPPING_TABLE_NAME
	}
	return val
}

func getDynamoApiKeyTableName() string { return os.Getenv(EnvVars[ENV_DYNAMO_API_KEYS_TABLE_NAME]) }

func getCloudfrontHostUrl() string { return os.Getenv(EnvVars[ENV_CLOUDFRONT_HOST_URL]) }

func getTimestampLimitSeconds() int {
	val, ok := os.LookupEnv(EnvVars[ENV_TIMESTAMP_LIMIT_SECONDS])
	if !ok {
		return DEFAULT_TIMESTAMP_LIMIT_SECONDS
	}

	limit, err := strconv.Atoi(val)
	if err != nil {
		return DEFAULT_TIMESTAMP_LIMIT_SECONDS
	}

	return limit
}

func getAllowedAgreementList() []string {
	val, ok := os.LookupEnv(EnvVars[ENV_ALLOWED_AGREEMENT_LIST])
	if !ok {
		return DEFAULT_ALLOWED_AGREEMENT_LIST
	}
	return strings.Split(strings.ReplaceAll(strings.ToLower(val), " ", ""), ",")
}

func getAllowedJurisdictionList() []string {
	val, ok := os.LookupEnv(EnvVars[ENV_ALLOWED_JURISDICTION_LIST])
	if !ok {
		return DEFAULT_ALLOWED_JURISDICTION_LIST
	}
	return strings.Split(strings.ReplaceAll(strings.ToLower(val), " ", ""), ",")
}

func getAllowedLocaleList() []string {
	val, ok := os.LookupEnv(EnvVars[ENV_ALLOWED_LOCALE_LIST])
	if !ok {
		return DEFAULT_ALLOWED_LOCALE_LIST
	}
	return strings.Split(strings.ReplaceAll(strings.ToLower(val), " ", ""), ",")
}
