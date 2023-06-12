## Overview

The parameter store library contains functions for reading key/values retrieved
from Secrets Manager and Parameter Store by the
[password_extension](../password_extension) layer, and written to a location in
shared [ephemeral
storage](https://docs.aws.amazon.com/lambda/latest/dg/API_EphemeralStorage.html).

1. Create the parameters and/or secrets configuration files in the
`config/`subdirectory, see [test_lambda](https://github.com/influencemobile/go-lambdas/tree/master/lambdas/test_lambda/config). 
Each file must list the keys to retrieve from Secrets Manager and Parameter Store. Ensure that the 
filenames match the filenames specified in (2) below.

``` yaml
keys:
  - "/honeybadger/lambdas/test_lambda/api_key"
```

2 Update the deployment yaml file, adding following environment variables 
to each environment (Production, Staging, and Development) per the directions 
in this section.

  - `CONFIG_PARAMETERS_FILE`, and
  - `PARAMETERS_STORE`, and / or
  - `CONFIG_SECRETS_FILE`, and
  - `SECRETS_STORE`

The [password_extension layer](../password_extension)  layer reads the 
files set in `CONFIG_PARAMETERS_FILE` and `CONFIG_SECRETS_FILE`, 
retrieves the key/values from Parameter Store and Secrets Manager, and 
writes these key/values to the files at the locations specified in
`PARAMETERS_STORE` and `SECRETS_STORE`.

It is possible to use different configurations for each environment 
(or share the same configurations across two or more environments) by 
naming the files in the `config/`subdirectory for the intended deployment 
environment,ie. `production_parameters.yaml`, `staging_parameters.yaml`, and
`development_parameters.yaml`. 

``` yaml
Production:
  EnvironmentVariables:
    CONFIG_SECRETS_FILE: '/var/task/config/production_secrets.yaml'
    SECRETS_STORE: '/tmp/secrets.txt'
    CONFIG_PARAMETERS_FILE: '/var/task/config/production_parameters.yaml'
    PARAMETERS_STORE: '/tmp/parameters.txt'
    
Production:
  EnvironmentVariables:
    CONFIG_SECRETS_FILE: '/var/task/config/staging_secrets.yaml'
    SECRETS_STORE: '/tmp/secrets.txt'
    CONFIG_PARAMETERS_FILE: '/var/task/config/staging_parameters.yaml'
    PARAMETERS_STORE: '/tmp/parameters.txt'
    
Production:
  EnvironmentVariables:
    CONFIG_SECRETS_FILE: '/var/task/config/development_secrets.yaml'
    SECRETS_STORE: '/tmp/secrets.txt'
    CONFIG_PARAMETERS_FILE: '/var/task/config/development_parameters.yaml'
    PARAMETERS_STORE: '/tmp/parameters.txt'
```

Note the `/var/task/` prefix that must be present as these files are placed at 
this location in the filesystem of the .zip package uploaded to AWS. The
`password_extension` will always write files into the `/tmp/` directory.

3. Add the `password_extension` layer to the deployment yaml file.

``` yaml
Production:
  Extensions: ['arn:aws:lambda:us-east-1:237634799245:layer:PasswordExtension:7']
```

4. Import the library at the top of your `.go` source file.
``` go
import (
	config "influencemobile.com/parameter_store"
)
```

5. Then create and assign the `KeyClient` object to a global variable.
Initialization is performed in `init()` so as only to be done only once when the
lambda starts.

``` go
var (
	parameterStore config.ParameterStoreIface
)

func init() {
	parameterStore = &config.KeyClient{Filename: os.Getenv("PARAMETERS_STORE")}
}
```

6. Add the following to the `go.mod` file.

``` shell
replace influencemobile.com/parameter_store => ../parameter_store
```

From the shell, have Go retrieve the library for use.

``` shell
% go get
```

7. Create a structure to hold the specific configuration parameters retrieved
   from Parameter Store or Secrets Manager. The example below will hold the
   Honeybadger API key retrieved from Parameter Store.

``` go
type ParameterStoreType struct {
	Name  string
	Value string
}
```

8. In the function called by `lambda.Start()` (in this example we use
   `lambda.Start(HandleRequest)`), read the configuration from the file and
   retrieve the key. The key in parameter store must match that specified in the
   configuration file in step #1.

``` go
func HandleRequest(ctx context.Context, kinesisEvent *events.KinesisEvent) {

	if err := parameterStore.ReadConfig(); err != nil {
		logger.
			NewEntry(log_msg_parameters_store, LogMsgTxt).
			Error(err.Error())
		os.Exit(1)
	}

	badgerParams, err := parameterStore.ParameterToBytes("/honeybadger/lambdas/test_lambda/api_key")
	if err != nil {
		logger.
			NewEntry(log_msg_parameters_store, LogMsgTxt).
			Error(err.Error())
		os.Exit(1)
	}
```

9. Then `Unmarshal` into a struct and read the API key.

``` go
	badgerApiKey := ParameterStoreType{}
	_ = json.Unmarshal(badgerParams, &badgerApiKey)

	badger.Configure(honeybadger.Configuration{
		APIKey: badgerApiKey.Value,
		Env:    env.Environment})
```

10. When unit testing, create an instance of `MockKeyClient` and assign dummy
    data that looks like what would be retrieved by the extension layer.

``` go
func TestHandleRequest(t *testing.T) {
	setEnvVars(map[TypeEnvVar]string{
		ENV_APPLICATION_NAME:     "adid_ingest_data_stream_test",
		ENV_LOG_LEVEL:            LogLevelText[Log_level_debug],
		ENV_ENVIRONMENT:          "unit_testing",
	})

	secretsStore = &config.MockKeyClient{
		FileBytes: []byte(`{"keys":[{"key":"dbuser/staging/admin/third_party_api_delete","value":{"dbInstanceIdentifier":"staging-databases-20210609","dbname":"adids","engine":"postgres","host":"localhost","password":"postgres","port":5432,"username":"postgres"}}]}`),
	}

	parameterStore = &config.MockKeyClient{
		FileBytes: []byte(`{"keys":[{"key":"/honeybadger/engage/api_key","value":{"Name":"/honeybadger/engage/api_key","Value":"api_key"}}]}`),
	}
```

11. It is easy to fallback to an AWS environment variables in the event that the
    Roles/Policies for the lambda are not yet authorized for Parameter Store, as
    follows;

```go
	// Retrieve Honeybadger API key
	badgerParams, err := parameterStore.ParameterToBytes(DEFAULT_HONEYBADGER_KEY)
	if err != nil {
		logger.
			NewEntry(log_msg_parameters_store, LogMsgTxt).
			Error(err.Error())
		params := fmt.Sprintf("{\"name\":\"/honeybadger/engage/api_key\",\"value\":\"%s\"}", os.Getenv(DEFAULT_HONEYBADGER_KEY))
		badgerParams = []byte(params)
	}
	badgerApiKey := ParameterStoreType{}
	_ = json.Unmarshal(badgerParams, &badgerApiKey)
	badger.Configure(honeybadger.Configuration{APIKey: badgerApiKey.Value, Env: env.Environment})
```




