package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/ssm"
	yaml "gopkg.in/yaml.v3"
	"influencemobile.com/password_extension/extension"

	. "influencemobile.com/logging"
)

type (
	ConfigFile struct {
		KeyVals []string `yaml:"keys"`
	}

	KeyVal struct {
		Key   string                 `json:"key"`
		Value map[string]interface{} `json:"value"`
	}

	KeyValFile struct {
		KeyVals []KeyVal `json:"keys"`
	}
)

var (
	// Always retrieve environment variables first.
	env = NewEnv()

	extensionName   = filepath.Base(os.Args[0]) // extension name has to match the filename
	extensionClient = extension.NewClient(os.Getenv("AWS_LAMBDA_RUNTIME_API"))
	printPrefix     = fmt.Sprintf("[%s]", extensionName)

	secretsFilename    = fmt.Sprintf("%s%s", os.Getenv("LAMBDA_TASK_ROOT"), env.ConfigSecretsFile)
	parametersFilename = fmt.Sprintf("%s%s", os.Getenv("LAMBDA_TASK_ROOT"), env.ConfigParametersFile)

	// Create a new logger
	logger = NewLoggers(os.Stderr, env.LogLevel, 0).
		SetFields(Fields{
			EnvVars[ENV_APP_NAME]:    extensionName,
			EnvVars[ENV_ENVIRONMENT]: env.Environment})
)

func main() {
	logger.
		NewEntry(log_msg_processing, LogMsgTxt).
		Debug("In extension:main()")

	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		// TODO: What does this do? We don't actaully need in this this
		// extension.
		s := <-sigs
		cancel()
		logger.
			NewEntry(log_msg_processing, LogMsgTxt).
			Debug(fmt.Sprintf("%s Received: %s", printPrefix, s))
	}()

	// An extension must register with AWS to receive the SHUTDOWN event. Ok.
	// But this particular extension actually doesn't need to register because
	// all it does is write to a file and exit. So this event handling logic is
	// superfluous.
	_, err := extensionClient.Register(ctx, extensionName)
	if err != nil {
		logger.
			NewEntry(log_msg_extension_error, LogMsgTxt).
			Error(err.Error())
		os.Exit(1)
	}

	// Load the secrets file. If it does not exist, then log this and assume
	// there is no secret file to load.
	if env.ConfigSecretsFile != "" {
		secrets, err := LoadConfig(secretsFilename)
		if err != nil {
			logger.
				NewEntry(log_msg_config_secrets_error, LogMsgTxt).
				Info(err.Error())
		}

		if secrets != nil {
			if err := WriteSecrets(secrets); err != nil {
				logger.
					NewEntry(log_msg_secrets_file_error, LogMsgTxt).
					Error(err.Error())
				// os.Exit(1)
			}
		}
	}

	// Load the parameters file. If it does not exist, then log this and assume
	// there is no parameters file to load.
	if env.ConfigParametersFile != "" {
		parameters, err := LoadConfig(parametersFilename)
		if err != nil {
			logger.
				NewEntry(log_msg_config_parameters_error, LogMsgTxt).
				Info(err.Error())
		}

		if parameters != nil {
			if err := WriteParameters(parameters); err != nil {
				logger.
					NewEntry(log_msg_parameters_file_error, LogMsgTxt).
					Error(err.Error())
				// os.Exit(1)
			}
		}
	}

	// Block until shutdown event is received or cancelled via the context.
	processEvents(ctx)
}

func processEvents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			res, err := extensionClient.NextEvent(ctx)
			if err != nil {
				return
			}
			// Exit if we receive a SHUTDOWN event
			if res.EventType == extension.Shutdown {
				return
			}
		}
	}
}

func GetSecretValue(svc *secretsmanager.SecretsManager, key string) (map[string]interface{}, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(key),
	}

	val, err := svc.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				return nil, fmt.Errorf("%s, %s", secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeInvalidParameterException:
				return nil, fmt.Errorf("%s, %s", secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				return nil, fmt.Errorf("%s, %s", secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeDecryptionFailure:
				return nil, fmt.Errorf("%s, %s", secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				return nil, fmt.Errorf("%s, %s", secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			default:
				return nil, fmt.Errorf(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return nil, fmt.Errorf(err.Error())
		}
	}

	var result map[string]interface{} = make(map[string]interface{})
	err = json.Unmarshal([]byte(*val.SecretString), &result)
	if err != nil {
		return nil, fmt.Errorf(LogMsgTxt[log_msg_secretsclient_json_error])
	}

	return result, nil
}

func GetParameter(svc *ssm.SSM, key string) (map[string]interface{}, error) {
	input := &ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(true),
	}

	val, err := svc.GetParameter(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ssm.ErrCodeParameterNotFound:
				return nil, fmt.Errorf("%s, %s", ssm.ErrCodeParameterNotFound, aerr.Error())
			case ssm.ErrCodeParameterVersionNotFound:
				return nil, fmt.Errorf("%s, %s", ssm.ErrCodeParameterVersionNotFound, aerr.Error())
			case ssm.ErrCodeInvalidKeyId:
				return nil, fmt.Errorf("%s, %s", ssm.ErrCodeInvalidKeyId, aerr.Error())
			case ssm.ErrCodeInternalServerError:
				return nil, fmt.Errorf("%s, %s", ssm.ErrCodeInternalServerError, aerr.Error())
			default:
				return nil, fmt.Errorf(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return nil, fmt.Errorf(err.Error())
		}
	}

	var result map[string]interface{} = make(map[string]interface{})
	res := struct {
		Name  string `json:"Name"`
		Value string `json:"Value"`
	}{
		Name:  *val.Parameter.Name,
		Value: *val.Parameter.Value,
	}

	temp, _ := json.Marshal(res)
	err = json.Unmarshal(temp, &result)
	if err != nil {
		return nil, fmt.Errorf("%s, %s, %s", LogMsgTxt[log_msg_ssmclient_json_error], err.Error(), string(temp))
	}

	// {
	//     "Parameter": {
	//         "ARN": "arn:aws:ssm:us-east-2:111122223333:parameter/MyGitHubPassword",
	//         "DataType": "text",
	//         "LastModifiedDate": 1582657288.8,
	//         "Name": "MyGitHubPassword",
	//         "Type": "SecureString",
	//         "Value": "AYA39c3b3042cd2aEXAMPLE/AKIAIOSFODNN7EXAMPLE/fh983hg9awEXAMPLE==",
	//         "Version": 3
	//     }
	// }

	return result, nil
}

func LoadConfig(fileName string) (*ConfigFile, error) {
	// os.Exit(1) if a config file is not specified.
	if fileName == "" {
		return nil, fmt.Errorf(LogMsgTxt[log_msg_config_secrets_unspecified])
	}

	// os.Exit(1) if the config file is not found.
	config := ConfigFile{}
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("%s: \"%s\"", LogMsgTxt[log_msg_config_secrets_io_error], fileName)
	}

	// os.Exit(1) if the config file is not yaml.
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return nil, fmt.Errorf(LogMsgTxt[log_msg_config_secrets_unmarshal_error])
	}
	return &config, nil
}

func WriteSecrets(config *ConfigFile) error {
	// Open a connection to Secrets Manager. Exit with error if fail.
	sess, err := session.NewSession()
	if err != nil {
		return fmt.Errorf("%s, %s", LogMsgTxt[log_msg_secretsclient_error], err.Error())
	}
	ssm := secretsmanager.New(sess)

	results := KeyValFile{}
	for _, secret := range config.KeyVals {
		// Retrieve the secret value. A map of key / value pairs.
		result, err := GetSecretValue(ssm, secret)
		if err != nil {
			return fmt.Errorf("%s, %s", LogMsgTxt[log_msg_secretsclient_get_secret], err.Error())
		}
		results.KeyVals = append(results.KeyVals, KeyVal{Key: secret, Value: result})
	}

	// Write the secret(s) to a file in "Ephemeral" storage.
	f, err := os.Create(env.SecretsStore)
	if err != nil {
		return fmt.Errorf("%s, %s", LogMsgTxt[log_msg_secrets_file_io_error], err.Error())
	}

	bytes, _ := json.Marshal(results)
	_, _ = f.Write(bytes) // OK to ignore error, as main lambda will complain if it cannot read file.
	f.Close()

	return nil
}

func WriteParameters(config *ConfigFile) error {
	// Open a connection to Secrets Manager. Exit with error if fail.
	sess, err := session.NewSession()
	if err != nil {
		return fmt.Errorf("%s, %s", LogMsgTxt[log_msg_ssmclient_error], err.Error())
	}
	ssm := ssm.New(sess)

	results := KeyValFile{}
	for _, secret := range config.KeyVals {
		// Retrieve the secret value. A map of key / value pairs.
		result, err := GetParameter(ssm, secret)
		if err != nil {
			return fmt.Errorf("%s, %s", LogMsgTxt[log_msg_ssmclient_get_parameter], err.Error())
		}
		results.KeyVals = append(results.KeyVals, KeyVal{Key: secret, Value: result})
	}

	// Write the secret(s) to a file in "Ephemeral" storage.
	f, err := os.Create(env.ParametersStore)
	if err != nil {
		return fmt.Errorf("%s, %s", LogMsgTxt[log_msg_parameters_file_io_error], err.Error())
	}

	bytes, _ := json.Marshal(results)
	_, _ = f.Write(bytes) // OK to ignore error, as main lambda will complain if it cannot read file.
	f.Close()

	return nil
}
