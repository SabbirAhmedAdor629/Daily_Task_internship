module github.com/sabbir/go_practice/Testing

go 1.19

replace influencemobile.com/libs/dynamo => ../libs/dynamo

require (
	github.com/lib/pq v1.10.7
	influencemobile.com/libs/dynamo v0.0.0-00010101000000-000000000000
)

require (
	github.com/aws/aws-sdk-go v1.44.118 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)
