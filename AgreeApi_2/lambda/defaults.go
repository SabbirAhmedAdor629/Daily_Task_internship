package main

import (
	"influencemobile.com/logging"
)

const (
	DEFAULT_LOG_LEVEL                 = logging.Log_level_info
	DEFAULT_APPLICATION_NAME          = "consent_management"
	DEFAULT_HONEYBADGER_KEY           = "/honeybadger/admin/api_key"
	DEFAULT_IPQS_KEY                  = "/ip_quality_score/private_key"
	DEFAULT_DYNAMO_CONSENT_TABLE_NAME = "consent_management"
	DEFAULT_DYNAMO_MAPPING_TABLE_NAME = "jurisdiction_mapping"
	DEFAULT_TIMESTAMP_LIMIT_SECONDS   = 30
	DEFAULT_TIME_FORMAT               = "2006-01-02T15:04:05Z" // ISO 8601
)

var (
	DEFAULT_ALLOWED_AGREEMENT_LIST    = []string{"pp", "tou", "idc"}
	DEFAULT_ALLOWED_JURISDICTION_LIST = []string{"gdpr"}
	DEFAULT_ALLOWED_LOCALE_LIST       = []string{"en-uk"}
)

const (
	headerForwardedFor = "X-Forwarded-For"
)

type QueryParamType int

const (
	queryStringTimezone QueryParamType = iota
)

var (
	QueryParams = []string{
		queryStringTimezone: "timezone",
	}
)

const (
	log_msg_end = iota + 100

	log_msg_begin
	log_msg_envars
	log_msg_processing
	log_msg_request
	log_msg_request_headers
	log_msg_request_body
	log_msg_request_parameters
	log_msg_invalid_path
	log_msg_invalid_query_params
	log_msg_invalid_body_parameters
	log_msg_invalid_headers
	log_msg_http_method
	log_msg_invalid_auth
	log_msg_json_marshal_error
	log_msg_json_unmarshal_error
	log_msg_data_field_missing
	log_msg_parameters_store_parse_error
	log_msg_parameters_store
	log_msg_dynamodb_error
	log_msg_dynamodb_write_error
	log_msg_dynamodb_table_name_error
	log_msg_invalid_player_id
	log_msg_jurisdiction_error
)

var (
	LogMsgTxt = []string{
		log_msg_begin:                        "begin",
		log_msg_envars:                       "envars",
		log_msg_processing:                   "processing",
		log_msg_request:                      "request",
		log_msg_request_headers:              "headers",
		log_msg_request_body:                 "body",
		log_msg_request_parameters:           "parameters",
		log_msg_end:                          "end",
		log_msg_invalid_path:                 "invalid request path",
		log_msg_invalid_query_params:         "invalid request query params",
		log_msg_invalid_body_parameters:      "invalid body parameters",
		log_msg_invalid_headers:              "invalid headers",
		log_msg_http_method:                  "invalid request method",
		log_msg_invalid_auth:                 "not authenticated",
		log_msg_json_marshal_error:           "JSON marshal error.",
		log_msg_json_unmarshal_error:         "JSON unmarshal error",
		log_msg_parameters_store_parse_error: "cannot parse parameters store",
		log_msg_parameters_store:             "parameters store",
		log_msg_dynamodb_error:               "DynamoDB error",
		log_msg_dynamodb_write_error:         "DynamoDB write error",
		log_msg_dynamodb_table_name_error:    "DynamoDB table name error",
		log_msg_invalid_player_id:            "Player Id is invalid",
		log_msg_jurisdiction_error:           "Jurisdiction set error",
	}
)

type ResponseCodeType int

// These are fixed. Do not change the order
const (
	ResponseInvalidPlayerId ResponseCodeType = iota + 2001 // Start from 2001, increment by one
	ResponseMissingJurisdiction
	ResponseInvalidJurisdiction
	ResponseMissingAgreement
	ResponseInvalidAgreement
	ResponseMissingLocale
	ResponseInvalidLocale
	ResponseMissingVersion
	ResponseInvalidVersion
	ResponseMissingTimestamp
	ResponseInvalidTimestamp
)

// These are fixed. Do not change the order
const (
	ResponseMissingTimezone ResponseCodeType = iota + 1002 // Start from 1001, increment by one
	ResponseInvalidTimezone
	ResponseInvalidIp
)
