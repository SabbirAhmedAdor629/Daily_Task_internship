package main

import (
	"net/http"

	// "influencemobile.com/logging"
)

const (
	// DEFAULT_LOG_LEVEL                 = logging.Log_level_info
	DEFAULT_CHANNEL_BUFFERS           = 16
	DEFAULT_WRITE_THREADS             = 2
	DEFAULT_APPLICATION_NAME          = "campaign_message"
	DEFAULT_HONEYBADGER_KEY           = "/honeybadger/messaging/api_key"
	DEFAULT_DB_CAMPAIGN_TABLE         = "campaign"
	DEFAULT_LOG_CHUNK_TABLE           = "campaign_execution_log_chunks"
	DEFAULT_DB_MESSAGE_TEMPLATE_TABLE = "message_templates"
	DEFAULT_DB_CATEGORY_TABLE         = "category"
	DEFAULT_REDIS_PORT                = 6379
)

const (
	headerXTrace = "x-amzn-trace-id"
)

var (
	httpStatusCodes = []string{
		http.StatusOK:                  "200 OK",
		http.StatusBadRequest:          "400 Bad Request",
		http.StatusNotFound:            "404 Not Found",
		http.StatusInternalServerError: "500 Internal Server Error",
	}
)

const (
	log_msg_end = iota + 100

	log_msg_begin
	log_msg_envars
	log_msg_processing
	log_msg_is_attributed_invalid_error
	log_msg_sqs_send_error
	log_msg_firehose_send_error
	log_msg_unknown_event
	log_msg_request
	log_msg_request_headers
	log_msg_request_body
	log_msg_invalid_body_parameters
	log_msg_invalid_path
	log_msg_http_method
	log_msg_json_marshal_error
	log_msg_json_unmarshal_error
	log_msg_processed_event
	log_msg_http_error
	log_msg_sqs_event_process_success
	log_msg_sqs_event_process_failed
	log_msg_http_internal_alb_get
	log_msg_http_internal_alb_post

	log_msg_secrets_store_error
	log_msg_secrets_store_parse_error
	log_msg_secrets_store

	log_msg_parameters_store_error
	log_msg_parameters_store_parse_error
	log_msg_parameters_store
	log_msg_pg_database_connect

	log_msg_file_read_error
	log_msg_execution_log_chunk_status_update_error
	log_msg_campaign_schedule_error
	log_msg_player_id_invalid_error
	log_msg_campaign_select_query_error
	log_msg_message_template_select_query_error
	log_msg_campaign_active_error
	log_msg_campaign_message_template_active_error
	log_msg_campaign_not_time_to_deliver_error
)

var (
	LogMsgTxt = []string{
		log_msg_begin:                       "begin",
		log_msg_envars:                      "envars",
		log_msg_processing:                  "processing",
		log_msg_end:                         "end",
		log_msg_is_attributed_invalid_error: "is_attributed value not found.",
		log_msg_sqs_send_error:              "data send to sqs failed.",
		log_msg_firehose_send_error:         "data send to firehose failed.",
		log_msg_unknown_event:               "Unknown Kinesis event received.",
		log_msg_request:                     "request",
		log_msg_request_headers:             "headers",
		log_msg_json_marshal_error:          "JSON marshal error.",
		log_msg_request_body:                "body",
		log_msg_json_unmarshal_error:        "JSON unmarshal error.",
		log_msg_invalid_body_parameters:     "invalid body parameters",
		log_msg_invalid_path:                "invalid request path",
		log_msg_http_method:                 "invalid request method",
		log_msg_processed_event:             "Processed Event",
		log_msg_http_error:                  "http request error",
		log_msg_sqs_event_process_success:   "sqs queue for the event processed correctly",
		log_msg_sqs_event_process_failed:    "sqs queue for the event not processed correctly",
		log_msg_http_internal_alb_get:       "GET request to internal ALB",
		log_msg_http_internal_alb_post:      "POST request to internal ALB",
		log_msg_secrets_store_error:         "cannot load secrets store file",
		log_msg_secrets_store_parse_error:   "cannot parse secrets store",
		log_msg_secrets_store:               "secrets store",

		log_msg_parameters_store_error:                  "cannot load parameters store file",
		log_msg_parameters_store_parse_error:            "cannot parse parameters store",
		log_msg_parameters_store:                        "parameters store",
		log_msg_pg_database_connect:                     "Cannot open connection to postgres database.",
		log_msg_file_read_error:                         "s3 file read error",
		log_msg_execution_log_chunk_status_update_error: "status update error into execution log chunk table",
		log_msg_campaign_schedule_error:                 "campaign schedule failed",
		log_msg_player_id_invalid_error:                 "player id invalid",
		log_msg_campaign_select_query_error:             "Get query failed from campaign table",
		log_msg_message_template_select_query_error:     "Get query failed from message_template table",
		log_msg_campaign_active_error:                   "Campaign not active",
		log_msg_campaign_message_template_active_error:  "Campaign Message Template not active",
		log_msg_campaign_not_time_to_deliver_error:      "Campaign not time to deliver",
	}
)
