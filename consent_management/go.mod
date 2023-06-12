module consent_management

go 1.19

replace influencemobile.com/logging => ../libs/logging

replace influencemobile.com/libs/dynamo => ../libs/dynamo

replace influencemobile.com/libs/api_key => ../libs/api_key

replace influencemobile.com/libs/uuid => ../libs/uuid

replace influencemobile.com/libs/v3_helpers => ../libs/v3_helpers

replace influencemobile.com/libs/badger => ../libs/badger

replace influencemobile.com/parameter_store => ../parameter_store

require (
	github.com/aws/aws-lambda-go v1.39.1
	github.com/aws/aws-sdk-go v1.44.238
	github.com/google/uuid v1.0.0
	github.com/honeybadger-io/honeybadger-go v0.5.0
	influencemobile.com/libs/api_key v0.0.0-00010101000000-000000000000
	influencemobile.com/libs/badger v0.0.0-00010101000000-000000000000
	influencemobile.com/libs/dynamo v0.0.0-00010101000000-000000000000
	influencemobile.com/libs/v3_helpers v0.0.0-00010101000000-000000000000
	influencemobile.com/logging v0.0.0-00010101000000-000000000000
	influencemobile.com/parameter_store v0.0.0-00010101000000-000000000000
)

require (
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/pborman/uuid v1.2.0 // indirect
	github.com/shirou/gopsutil v2.18.12+incompatible // indirect
	golang.org/x/sys v0.1.0 // indirect
	influencemobile.com/libs/uuid v0.0.0-00010101000000-000000000000 // indirect
)
