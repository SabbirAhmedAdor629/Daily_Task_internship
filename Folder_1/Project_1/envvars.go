package main

import (
	"os"
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
	ENV_SECRET_STORE
	ENV_POSTGRES_KEY
	ENV_POSTGRES_DB
	ENV_POSTGRES_DB_SSL
	ENV_S3_SOURCE_BUCKET
	ENV_DB_CAMPAIGN_TABLE
	ENV_DB_LOG_CHUNK_TABLE
	ENV_DB_MESSAGE_TEMPLATE_TABLE
	ENV_DB_CATEGORY_TABLE
)

var (
	EnvVars = []string{
		ENV_ENVIRONMENT:               "TARGET_ENVIRONMENT",
		ENV_APPLICATION_NAME:          "APP_NAME",
		ENV_LOG_LEVEL:                 "LOG_LEVEL",
		ENV_AWS_REGION:                "REGION",
		ENV_PARAMETERS_STORE:          "PARAMETERS_STORE",
		ENV_SECRET_STORE:              "SECRETS_STORE",
		ENV_POSTGRES_KEY:              "POSTGRES_KEY",
		ENV_POSTGRES_DB:               "POSTGRES_DB",
		ENV_POSTGRES_DB_SSL:           "POSTGRES_DB_SSL",
		ENV_S3_SOURCE_BUCKET:          "S3_SOURCE_BUCKET",
		ENV_DB_CAMPAIGN_TABLE:         "DB_CAMPAIGN_TABLE",
		ENV_DB_LOG_CHUNK_TABLE:        "DB_LOG_CHUNK_TABLE",
		ENV_DB_MESSAGE_TEMPLATE_TABLE: "DB_MESSAGE_TEMPLATE_TABLE",
		ENV_DB_CATEGORY_TABLE:         "DB_CATEGORY_TABLE",
	}
)

// Retrive region, sqs name, kinesis firehose name from lambda environment values
type envVars struct {
	LogLevel               int     `json:"LOG_LEVEL"`
	Environment            string  `json:"TARGET_ENVIRONMENT"`
	ApplicationName        string  `json:"APP_NAME"`
	region                 string  `json:"REGION"`
	ParametersStore        string  `json:"PARAMETERS_STORE"`
	SecretsStore           string  `json:"SECRETS_STORE"`
	PostgresKey            string  `json:"POSTGRES_KEY"`
	PostgresDB             string  `json:"POSTGRES_DB"`
	DbSsl                  SslType `json:"POSTGRES_DB_SSL"`
	S3SourceBucket         string  `json:"S3_SOURCE_BUCKET"`
	DbCampaignTable        string  `json:"DB_CAMPAIGN_TABLE"`
	DbLogChunkTable        string  `json:"DB_LOG_CHUNK_TABLE"`
	DbMessageTemplateTable string  `json:"DB_MESSAGE_TEMPLATE_TABLE"`
	DbCategoryTable        string  `json:"DB_CATEGORY_TABLE"`
}

func NewEnv() envVars {
	return envVars{
		LogLevel:               getLogLevel(),
		Environment:            getEnvironment(),
		ApplicationName:        getApplicationName(),
		region:                 getRegion(),
		ParametersStore:        getParametersStore(),
		SecretsStore:           getSecretsStore(),
		PostgresKey:            getPostgresKey(),
		PostgresDB:             getPostgresDB(),
		DbSsl:                  getDbSsl(),
		S3SourceBucket:         getS3SourceBucket(),
		DbCampaignTable:        getDbCampaignTable(),
		DbLogChunkTable:        getDbLogChunkTable(),
		DbMessageTemplateTable: getDbMessageTemplateTable(),
		DbCategoryTable:        getDbCategoryTable(),
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
func getSecretsStore() string    { return os.Getenv(EnvVars[ENV_SECRET_STORE]) }
func getPostgresKey() string     { return os.Getenv(EnvVars[ENV_POSTGRES_KEY]) }
func getPostgresDB() string      { return os.Getenv(EnvVars[ENV_POSTGRES_DB]) }
func getS3SourceBucket() string  { return os.Getenv(EnvVars[ENV_S3_SOURCE_BUCKET]) }

func getDbSsl() SslType {
	sslStr := os.Getenv(EnvVars[ENV_POSTGRES_DB_SSL])
	for i, val := range SslConnection {
		if val == sslStr {
			return SslType(i)
		}
	}
	return AlwaysSsl
}

func getDbCampaignTable() string {
	val, ok := os.LookupEnv(EnvVars[ENV_DB_CAMPAIGN_TABLE])
	if !ok {
		return DEFAULT_DB_CAMPAIGN_TABLE
	}
	return val
}

func getDbLogChunkTable() string {
	val, ok := os.LookupEnv(EnvVars[ENV_DB_LOG_CHUNK_TABLE])
	if !ok {
		return DEFAULT_LOG_CHUNK_TABLE
	}
	return val
}

func getDbMessageTemplateTable() string {
	val, ok := os.LookupEnv(EnvVars[ENV_DB_MESSAGE_TEMPLATE_TABLE])
	if !ok {
		return DEFAULT_DB_MESSAGE_TEMPLATE_TABLE
	}
	return val
}

func getDbCategoryTable() string {
	val, ok := os.LookupEnv(EnvVars[ENV_DB_CATEGORY_TABLE])
	if !ok {
		return DEFAULT_DB_CATEGORY_TABLE
	}
	return val
}
