
# Overview

The `password_extension` layer retrieves key/values from AWS Secrets Manager and
AWS Parameter Store and writes these to a location in shared ephemeral storage.

- The [`parameter_store`
library](https://github.com/influencemobile/go-lambdas/blob/master/lambdas/parameter_store/README.md#overview)
describes how to use this layer.
- The
  [`test_lambda`](https://github.com/influencemobile/go-lambdas/blob/master/lambdas/test_lambda/README.md)
  is an example of how to integrate this layer into a lambda.

# Environment Variables
 
- `TARGET_ENVIRONMENT`:: One of `production`, `staging`, `development` or
  `unit_testing`. 
- `LOG_LEVEL` :: Will write log entries at the specified log level and higher.
  Supports `DEBUG`, `INFO`, `ERROR`, `FATAL` and `NONE`.
- `CONFIG_PARAMETERS_FILE` :: Set to `/var/task/config/staging_parameters.yaml`.
  Used by the [`password_extension` layer](#Layers_/_Libraries).
- `CONFIG_SECRETS_FILE` :: `/var/task/config/staging_secrets.yaml`. Used by the
  [`password_extension` layer](#Layers_/_Libraries).
- `SECRETS_STORE` : `/tmp/secrets.txt`. Used by the [`password_extension`
  layer](#Layers_/_Libraries).
- `PARAMETERS_STORE` :: `/tmp/parameters.txt`. Used by the [`password_extension`
  layer](#Layers_/_Libraries).

# AWS Dependencies

- Permissions to output logs to Cloudwatch,
- Permissions to read from Parameter Store.
- Permissions to read from Secure Storage.

# Build and Packaging 

1. Build the layer as follows; 

``` shell
go build -o bin/extensions/password_extension
```

2. In the `bin/` subdirectory, package the binary as follows;

``` shell
zip extensions.zip extensions/password_extension
```
