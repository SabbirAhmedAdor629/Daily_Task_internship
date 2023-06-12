package main

type OperationalCounters struct {
	DynamoDbError            int `json:"dynamodb_error"`
	InvalidQueryParameters   int `json:"invalid_query_parameters"`
	InvalidBodyParameters    int `json:"invalid_body_parameters"`
	InvalidAuth              int `json:"invalid_auth"`
	InvalidPath              int `json:"invalid_path"`
	InvalidRequestHeaders    int `json:"invalid_headers"`
	UnexpectedMarshalError   int `json:"unexpected_marshal_error,omitempty"`
	InvalidPlayerId          int `json:"invalid_player_id,omitempty"`
	UnexpectedUnmarshalError int `json:"unexpected_unmarshal_error,omitempty"`
}

func (c *OperationalCounters) incDynamoDbError(num int) {
	c.DynamoDbError = c.DynamoDbError + num
}

func (c *OperationalCounters) incUnexpectedMarshalError(num int) {
	c.UnexpectedMarshalError = c.UnexpectedMarshalError + num
}

func (c *OperationalCounters) incUnexpectedUnmarshalError(num int) {
	c.UnexpectedUnmarshalError = c.UnexpectedUnmarshalError + num
}

func (c *OperationalCounters) incInvalidBodyParameters(num int) {
	c.InvalidBodyParameters = c.InvalidBodyParameters + num
}

func (c *OperationalCounters) incInvalidAuth(num int) {
	c.InvalidAuth = c.InvalidAuth + num
}

func (c *OperationalCounters) incInvalidQueryParameters(num int) {
	c.InvalidQueryParameters = c.InvalidQueryParameters + num
}

func (c *OperationalCounters) incInvalidRequestHeaders(num int) {
	c.InvalidRequestHeaders = c.InvalidRequestHeaders + num
}

func (c *OperationalCounters) incInvalidPath(num int) {
	c.InvalidPath = c.InvalidPath + num
}

func (c *OperationalCounters) incInvalidPlayerId(num int) {
	c.InvalidPlayerId = c.InvalidPlayerId + num
}
