module im_services

go 1.13

replace influencemobile.com/libs/dynamo => ../dynamo

require (
	github.com/aws/aws-sdk-go v1.44.165
	github.com/go-redis/redis/v9 v9.0.0-rc.2
	influencemobile.com/libs/dynamo v0.0.0-00010101000000-000000000000
)
