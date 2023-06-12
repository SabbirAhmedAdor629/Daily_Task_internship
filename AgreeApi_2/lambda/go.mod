module AgreeApi

go 1.17

replace influencemobile.com/logging => ../libs/logging

replace influencemobile.com/libs/dynamo => ../libs/dynamo

replace influencemobile.com/libs/api_key => ../libs/api_key

replace influencemobile.com/libs/uuid => ../libs/uuid

replace influencemobile.com/libs/v3_helpers => ../libs/v3_helpers

replace influencemobile.com/libs/badger => ../libs/badger

replace influencemobile.com/parameter_store => ../parameter_store

require (
	github.com/aws/aws-lambda-go v1.41.0
	github.com/aws/aws-sdk-go v1.44.263
	influencemobile.com/libs/dynamo v0.0.0-00010101000000-000000000000
	influencemobile.com/libs/v3_helpers v0.0.0-00010101000000-000000000000
	influencemobile.com/parameter_store v0.0.0-00010101000000-000000000000
)

require (
	github.com/google/uuid v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	influencemobile.com/libs/api_key v0.0.0-00010101000000-000000000000 // indirect
	influencemobile.com/libs/uuid v0.0.0-00010101000000-000000000000 // indirect
	influencemobile.com/logging v0.0.0-00010101000000-000000000000
)
