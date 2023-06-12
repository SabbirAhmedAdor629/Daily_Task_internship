package v3_helpers

import "net/http"

type HttpMethodType string

var (
	httpStatusCodes = []string{
		http.StatusOK:                  "200 OK",
		http.StatusBadRequest:          "400 Bad Request",
		http.StatusUnauthorized:        "401 Unauthorized",
		http.StatusNotFound:            "404 Not Found",
		http.StatusInternalServerError: "500 Internal Server Error",
	}

	HttpGet  HttpMethodType = "GET"
	HttpPost HttpMethodType = "POST"
)

type ResponseCodeType int

// These are fixed.
// Do not change order. Do not insert. Only append.
const (
	ResponseSuccess ResponseCodeType = iota // Start from zero, increment by one
	ResponseNotAuthenticated
	ResponseNotAuthorized
	ResponseMissingPlatformHeader
	ResponseInvalidPlatformnHeader
	ResponseMissingPlatformVersionHeader
	ResponseInvalidPlatformVersionHeader
	ResponseMissingBuildVersionHeader
	ResponseInvalidBuildVersionHeader
	ResponseMissingProgramSlugHeader
	ResponseInvalidProgramSlugHeader
	ResponseMissingRequestIdHeader
	ResponseInvalidRequestIdHeader
	ResponseMissingMember
	ResponseInvalidMember
	ResponseMissingCorrelationId
	ResponseInvalidCorrelationId
)

// These are fixed.
// Do not change order. Do not insert. Only append.
const (
	ResponseInternalServerError  ResponseCodeType = iota + 500
	ResponseInvalidRequestMethod                  // POST vs GET
	ResponseInvalidRequestPath                    // /api/v3/blee/blah vs /api/v3/blah/blee
	ResponseInvalidRequestBodyParameters
	ResponseMissingApiKey
	ResponseInvalidApiKey
	ResponseMissingAdvertisingId
	ResponseInvalidAdvertisingId
	ResponseMissingAppId
	ResponseInvalidAppId
	ResponseMissingContinuationToken
	ResponseInvalidContinuationToken
	ResponseMissingMessageId
	ResponseInvalidMessageId
	ResponseMissingOfferId
	ResponseInvalidOfferId
)

var ResponseCodeToKey = map[string]int{
	"api_key":            int(ResponseMissingApiKey),
	"request_id":         int(ResponseMissingRequestIdHeader),
	"platform":           int(ResponseMissingPlatformHeader),
	"platform_version":   int(ResponseMissingPlatformVersionHeader),
	"build_version":      int(ResponseMissingBuildVersionHeader),
	"program_slug":       int(ResponseMissingProgramSlugHeader),
	"player_id":          int(ResponseMissingMember),
	"member_id":          int(ResponseMissingMember),
	"device_id":          int(ResponseMissingMember),
	"advertising_id":     int(ResponseMissingAdvertisingId),
	"app_id":             int(ResponseMissingAppId),
	"continuation_token": int(ResponseMissingContinuationToken),
	"message_id":         int(ResponseMissingMessageId),
	"offer_id":           int(ResponseMissingOfferId),
	"correlation_id":     int(ResponseMissingCorrelationId),
}

var MandatoryHeaders = []string{
	"api_key",
	"request_id",
	"platform",
	"platform_version",
	"build_version",
	"program_slug",
}
