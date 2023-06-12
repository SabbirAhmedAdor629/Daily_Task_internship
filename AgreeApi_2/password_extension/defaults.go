package main

import "influencemobile.com/logging"

const (
	DEFAULT_LOG_LEVEL          = logging.Log_level_info
	DEFAULT_APP_NAME           = "password_extension"
	DEFAULT_SECRETS_STORE      = "/tmp/secrets.txt"
	DEFAULT_PARAMETERS_STORE   = "/tmp/parameters.txt"
	DEFAULT_TARGET_ENVIRONMENT = "unknown"
	// DEFAULT_ERROR_URL = "https://github.com/influencemobile/go-lambdas/tree/master/password_extension"
	// DEFAULT_CONFIG_SECRETS_FILE    = "var/task/config/secrets.yaml"
	// DEFAULT_CONFIG_PARAMETERS_FILE = "var/task/config/parameters.yaml"
)

const (
	log_msg_end = iota + 100
	log_msg_begin
	log_msg_envars
	log_msg_processing
	log_msg_secretsclient_error
	log_msg_secretsclient_get_secret
	log_msg_secretsclient_json_error
	log_msg_ssmclient_error
	log_msg_ssmclient_get_parameter
	log_msg_ssmclient_json_error
	log_msg_secrets_file_error
	log_msg_parameters_file_error
	log_msg_secrets_file_io_error
	log_msg_parameters_file_io_error
	log_msg_config_secrets_error
	log_msg_config_parameters_error
	log_msg_config_secrets_unspecified
	log_msg_config_secrets_io_error
	log_msg_config_secrets_unmarshal_error
	log_msg_extension_error
)

var (
	LogMsgTxt = []string{
		log_msg_begin:                          "begin",
		log_msg_end:                            "end",
		log_msg_envars:                         "envars",
		log_msg_processing:                     "processing",
		log_msg_secretsclient_get_secret:       "Unable to retrieve secret",
		log_msg_secretsclient_json_error:       "Unable to parse the Secrets Client secret",
		log_msg_secretsclient_error:            "Unable to create a Secrets Client Session",
		log_msg_ssmclient_get_parameter:        "Unable to retrieve parameter",
		log_msg_ssmclient_json_error:           "Unable to parse the parameter ",
		log_msg_ssmclient_error:                "Unable to create a SSM session",
		log_msg_secrets_file_error:             "Unable to create the secrets file",
		log_msg_parameters_file_error:          "Unable to create the parameters file",
		log_msg_secrets_file_io_error:          "Unable to write to the secrets file",
		log_msg_parameters_file_io_error:       "Unable to write to the parameters file",
		log_msg_config_secrets_error:           "Secrets config file error",
		log_msg_config_parameters_error:        "Parameters config file error",
		log_msg_config_secrets_unspecified:     "Location to secrets file not specified",
		log_msg_config_secrets_io_error:        "Unable to read secrets file",
		log_msg_config_secrets_unmarshal_error: "Unable to parse secrets file",
		log_msg_extension_error:                "Lambda extension error",
	}
)
