module password_extension

go 1.18

replace influencemobile.com/password_extension/extension => ../password_extension/extension

replace influencemobile.com/logging => ../logging

require (
	github.com/aws/aws-sdk-go v1.43.35
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	influencemobile.com/logging v0.0.0-00010101000000-000000000000
	influencemobile.com/password_extension/extension v0.0.0-00010101000000-000000000000
)

require github.com/jmespath/go-jmespath v0.4.0 // indirect
