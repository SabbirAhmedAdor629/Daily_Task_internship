module campaign_message

go 1.19

replace influencemobile.com/parameter_store => ../parameter_store

replace influencemobile.com/logging => ../logging

replace influencemobile.com/badger => ../badger

replace influencemobile.com/libs/im_targeting => ../libs/targeting

replace influencemobile.com/libs/dynamo => ../libs/dynamo

require (
	github.com/aws/aws-lambda-go v1.34.1
	github.com/aws/aws-sdk-go v1.44.165
	github.com/go-redis/redis/v9 v9.0.0-rc.2
	github.com/honeybadger-io/honeybadger-go v0.5.0
	github.com/lib/pq v1.10.7
	influencemobile.com/badger v0.0.0-00010101000000-000000000000
	influencemobile.com/libs/im_targeting v0.0.0-00010101000000-000000000000
	influencemobile.com/logging v0.0.0-00010101000000-000000000000
	influencemobile.com/parameter_store v0.0.0-00010101000000-000000000000
)

require (
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/google/uuid v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/pborman/uuid v1.2.0 // indirect
	github.com/shirou/gopsutil v2.18.12+incompatible // indirect
	golang.org/x/sys v0.2.0 // indirect
	influencemobile.com/libs/dynamo v0.0.0-00010101000000-000000000000 // indirect
)
