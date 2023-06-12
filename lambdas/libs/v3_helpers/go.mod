module v3_helpers

go 1.18

require (
	github.com/aws/aws-lambda-go v1.34.1
	github.com/aws/aws-sdk-go v1.44.118
	influencemobile.com/libs/api_key v0.0.0-00010101000000-000000000000
	influencemobile.com/libs/uuid v0.0.0-00010101000000-000000000000
)

require (
	github.com/google/uuid v1.0.0
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)

replace influencemobile.com/libs/uuid => ../uuid

replace influencemobile.com/libs/api_key => ../api_key
