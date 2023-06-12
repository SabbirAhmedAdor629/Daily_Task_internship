package v3_helpers

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/google/uuid"
	im_keycache "influencemobile.com/libs/api_key"
	im_uuid "influencemobile.com/libs/uuid"
)

type (
	AppIdType     string
	AdIdType      string
	PlayerIdType  int64
	MessageIdType string
	MemberIdType  string

	ContinuationTokenType map[string]*dynamodb.AttributeValue

	ResponseBody struct {
		Error int          `json:"error"`
		Data  *interface{} `json:"data"`
	}
)

func NewPtrS(s string) *string {
	ptr := new(string)
	*ptr = s
	return ptr
}
func NewPtrB(b bool) *bool {
	ptr := new(bool)
	*ptr = b
	return ptr
}
func NewPtrI(i int) *int {
	ptr := new(int)
	*ptr = i
	return ptr
}

// Merge header and body. Key/values in header will overwrite any duplicates in
// body
func MergeHeaderAndBody(headers map[string]string, body map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})

	bytes, _ := json.Marshal(body)
	json.Unmarshal(bytes, &merged)

	for k, v := range headers {
		merged[k] = v
	}

	return merged
}

// CreateResponseHeaders retuns a copy of headers, only including keys specified
// in keepList.
func CreateResponseHeaders(headers map[string]string, keepList []string) map[string]string {
	// Create a map of the keeps
	keep := make(map[string]int)
	for _, v := range keepList {
		keep[v]++
	}

	// Copy the headers
	h := make(map[string]string)
	for k, v := range headers {
		if _, ok := keep[k]; ok {
			h[k] = v
		}
	}
	return h
}

func ValidateRequestMethod(requestMethod string, method HttpMethodType) (ResponseCodeType, error) {
	if requestMethod != string(method) {
		return ResponseInvalidRequestMethod,
			fmt.Errorf("Invalid request method")
	}
	return ResponseSuccess, nil
}

func ValidateRequestPath(requestPath string, path *regexp.Regexp) (ResponseCodeType, error) {
	if ok := path.FindStringSubmatch(requestPath); ok == nil {
		return ResponseInvalidRequestPath,
			fmt.Errorf("Invalid request path")
	}
	return ResponseSuccess, nil
}

func ValidateRequestHeaders(request map[string]string, mandatory []string, responseCodes map[string]int) (map[string]string, int, error) {
	if mandatory != nil && responseCodes == nil {
		return request, int(ResponseInternalServerError), fmt.Errorf("No ResponseCodes provided")
	}
	// Ensure mandatory headers are present
	for _, mVal := range mandatory {
		if _, ok := request[mVal]; !ok {
			if _, ok := responseCodes[mVal]; !ok {
				return request, int(ResponseInternalServerError), fmt.Errorf("ResponseCode for %s not found", mVal)
			}
			return request, responseCodes[mVal], fmt.Errorf(mVal + " not present")
		}
	}

	// Validate request_id, if present
	if rid, present := request["request_id"]; present {
		if _, err := ValidateRequestId(fmt.Sprintf("%s", rid)); err != nil {
			return request, int(ResponseInvalidRequestIdHeader), err
		}
	}
	return request, int(ResponseSuccess), nil
}

func ValidateBodyParameters(body string, isBase64Encoded bool, mandatory []string, responseCodes map[string]int) (map[string]interface{}, int, error) {
	var byteBody []byte = []byte(body)

	if isBase64Encoded {
		var err error
		byteBody, err = hex.DecodeString(body)
		if err != nil {
			return nil, int(ResponseInvalidRequestBodyParameters), fmt.Errorf("Cannot retrieve parameters")
		}
	}

	if mandatory != nil && responseCodes == nil {
		return nil, int(ResponseInternalServerError), fmt.Errorf("No ResponseCodes provided")
	}

	// Ensure mandatory params are present
	mapParams := make(map[string]interface{})
	if err := json.Unmarshal(byteBody, &mapParams); err != nil {
		return nil, int(ResponseInvalidRequestBodyParameters), fmt.Errorf("Cannot retrieve parameters")
	}
	for _, mVal := range mandatory {
		if _, ok := mapParams[mVal]; !ok {
			if _, ok := responseCodes[mVal]; !ok {
				return nil, int(ResponseInternalServerError), fmt.Errorf("ResponseCode for %s not found", mVal)
			}
			return nil, responseCodes[mVal], fmt.Errorf(mVal + " not present")
		}
	}

	return mapParams, int(ResponseSuccess), nil
}

func ValidateContinuationId(cid string) (*uuid.UUID, error) { return ValidateRequestId(cid) }
func ValidateRequestId(rid string) (*uuid.UUID, error) {
	requestId, err := im_uuid.ValidateUUID(rid)
	if err != nil {
		return nil, err
	}

	return requestId, nil
}

func ValidateAdId(id *string) (AdIdType, int, error) {
	if id != nil && *id != *NewPtrS("") {
		return AdIdType(*id), int(ResponseSuccess), nil
	}
	return "", int(ResponseMissingAdvertisingId), fmt.Errorf("advertising_id invalid ")
}

func ValidateAppId(id *string) (AppIdType, int, error) {
	if id != nil && *id != *NewPtrS("") {
		return AppIdType(*id), int(ResponseSuccess), nil
	}
	return "", int(ResponseMissingAppId), fmt.Errorf("app_id invalid ")
}

func ValidatePlayerId(id *string) (PlayerIdType, int, error) {
	if id != nil && *id != *NewPtrS("") {
		i, err := strconv.Atoi(*id)
		if err != nil {
			return PlayerIdType(i), int(ResponseInvalidMember), fmt.Errorf("player_id invalid")
		}
		return PlayerIdType(i), int(ResponseSuccess), nil
	}
	return 0, int(ResponseMissingMember), fmt.Errorf("player_id missing")
}

func ValidateMessageId(id *string) (MessageIdType, int, error) {
	if id != nil && *id != *NewPtrS("") {
		return MessageIdType(*id), int(ResponseSuccess), nil
	}
	return "", int(ResponseMissingMessageId), fmt.Errorf("message_id invalid")
}

func ValidateMemberId(id *string) (MemberIdType, int, error) {
	if id != nil && *id != *NewPtrS("") {
		return MemberIdType(*id), int(ResponseSuccess), nil
	}
	return "", int(ResponseMissingMember), fmt.Errorf("member_id invalid")
}

func ValidateOfferId(id *string) (string, int, error) {
	if id != nil && *id != *NewPtrS("") {
		return *id, int(ResponseSuccess), nil
	}
	return "", int(ResponseMissingOfferId), fmt.Errorf("offer_id invalid")
}

func ValidateAuth(svc dynamodbiface.DynamoDBAPI, apiCache *im_keycache.ApiKeyType, tableName, apiKey string) (bool, int, error) {
	if apiKey == "" {
		return false, int(ResponseMissingApiKey), nil
	}

	found, err := apiCache.ValidateKey(svc, tableName, "api_key", apiKey)
	if err != nil {
		return false, int(ResponseInternalServerError), err
	}

	if !found {
		return false, int(ResponseInvalidApiKey), nil
	}

	return true, int(ResponseSuccess), nil
}

func ValidateParamI(name string, value *int64) (int64, string, int, error) {
	if value != nil {
		return *value, "", http.StatusOK, nil
	}
	return 0,
		"[ Process Failed. ] A valid " + name + " is required",
		http.StatusBadRequest,
		fmt.Errorf("%s required", name)
}

func ValidateParamS(name string, value *string) (string, string, int, error) {
	if value != nil && *value != *NewPtrS("") {
		return *value, "", http.StatusOK, nil
	}
	return "",
		"[ Process Failed. ] A valid " + name + " is required",
		http.StatusBadRequest,
		fmt.Errorf("%s required", name)
}

func CreateReallySimpleBody(error int) ResponseBody {
	return CreateSimpleBody(error, nil)
}
func CreateSimpleBody(error int, data map[string]interface{}) ResponseBody {
	return CreateBody(error, data, nil, nil)
}

func CreateBody(error int, data map[string]interface{}, numResults *int, continuationToken *ContinuationTokenType) ResponseBody {
	body := make(map[string]interface{})
	if numResults != nil {
		body["num_results"] = numResults
	}
	if continuationToken != nil {
		body["continuation_token"] = continuationToken
	}
	for k, v := range data {
		body[k] = v
	}
	var temp interface{} = data
	return ResponseBody{
		Error: error,
		Data:  &temp,
	}
}

func CreateHttpResponse(statusCode int, headers map[string]string, body ResponseBody) *events.ALBTargetGroupResponse {
	b, _ := json.Marshal(body)

	h := make(map[string]string)
	for k, v := range headers {
		h[k] = v
	}
	if _, ok := h["Set-cookie"]; !ok {
		h["Set-cookie"] = "cookies"
	}
	if _, ok := h["Content-Type"]; !ok {
		h["Content-Type"] = "application/json"
	}

	return &events.ALBTargetGroupResponse{
		StatusCode:        statusCode,
		StatusDescription: httpStatusCodes[statusCode],
		Headers:           h,
		Body:              string(b),
		IsBase64Encoded:   false,
	}
}
