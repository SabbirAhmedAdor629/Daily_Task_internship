package main

import (
	"os"
	"strings"

	"influencemobile.com/logging"
)

type TypeEnvVar int

const (
	ENV_ENVIRONMENT TypeEnvVar = iota
	ENV_APP_NAME
	ENV_LOG_LEVEL
	ENV_SECRETS_STORE
	ENV_PARAMETERS_STORE
	ENV_CONFIG_SECRETS_FILE
	ENV_CONFIG_PARAMETERS_FILE
)

var (
	EnvVars = []string{
		ENV_ENVIRONMENT:            "TARGET_ENVIRONMENT",
		ENV_APP_NAME:               "APP_NAME",
		ENV_LOG_LEVEL:              "LOG_LEVEL",
		ENV_SECRETS_STORE:          "SECRETS_STORE",
		ENV_PARAMETERS_STORE:       "PARAMETERS_STORE",
		ENV_CONFIG_SECRETS_FILE:    "CONFIG_SECRETS_FILE",
		ENV_CONFIG_PARAMETERS_FILE: "CONFIG_PARAMETERS_FILE",
	}
)

type envVars struct {
	LogLevel             int    `json:"log_level"`
	Environment          string `json:"target_environment"`
	ApplicationName      string `json:"app_name"`
	SecretsStore         string `json:"secrets_store"`
	ParametersStore      string `json:"parameters_store"`
	ConfigSecretsFile    string `json:"config_secrets_file"`
	ConfigParametersFile string `json:"config_parameters_file"`
}

func NewEnv() envVars {
	return envVars{
		LogLevel:             getLogLevel(),
		Environment:          getEnvironment(),
		ApplicationName:      getApplicationName(),
		SecretsStore:         getSecretsStore(),
		ParametersStore:      getParametersStore(),
		ConfigParametersFile: getConfigParametersFile(),
		ConfigSecretsFile:    getConfigSecretsFile(),
	}
}

func getLogLevel() int {
	level := DEFAULT_LOG_LEVEL

	temp, exist := os.LookupEnv(EnvVars[ENV_LOG_LEVEL])
	if !exist || temp == "" {
		return level
	}
	txtLevel := strings.ToUpper(temp)
	for i, l := range logging.LogLevelText {
		if txtLevel == l {
			level = i
		}
	}
	return level
}

func getApplicationName() string {
	appname := os.Getenv(EnvVars[ENV_APP_NAME])
	if appname == "" {
		return DEFAULT_APP_NAME
	}
	return appname
}

func getEnvironment() string {
	temp, exist := os.LookupEnv(EnvVars[ENV_ENVIRONMENT])
	if !exist || temp == "" {
		return DEFAULT_TARGET_ENVIRONMENT
	}
	return temp
}

func getSecretsStore() string {
	temp, exist := os.LookupEnv(EnvVars[ENV_SECRETS_STORE])
	if !exist || temp == "" {
		return DEFAULT_SECRETS_STORE
	}
	return temp
}

func getParametersStore() string {
	temp, exist := os.LookupEnv(EnvVars[ENV_PARAMETERS_STORE])
	if !exist || temp == "" {
		return DEFAULT_PARAMETERS_STORE
	}
	return temp
}

func getConfigSecretsFile() string {
	temp, exist := os.LookupEnv(EnvVars[ENV_CONFIG_SECRETS_FILE])
	if !exist {
		return ""
	}
	// if !exist || temp == "" {
	// 	return DEFAULT_CONFIG_SECRETS_FILE
	// }
	return temp
}

func getConfigParametersFile() string {
	temp, exist := os.LookupEnv(EnvVars[ENV_CONFIG_PARAMETERS_FILE])
	if !exist {
		return ""
	}
	// if !exist || temp == "" {
	// 	return DEFAULT_CONFIG_PARAMETERS_FILE
	// }
	return temp
}
