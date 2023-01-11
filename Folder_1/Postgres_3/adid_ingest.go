package main

import (
	"bufio"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/honeybadger-io/honeybadger-go"

	. "influencemobile.com/logging"
	config "influencemobile.com/parameter_store"
)

const (
	s3Event_type = "ObjectCreated:Put"
)

type (
	hashBuffer [32]byte
)

var (
	opCounters *OperationalCounters
	env        env_vars
	logger     *Loggers
	svc        s3iface.S3API
	pgDb       PgClientIface
	badger     badgerIface

	parameterStore config.ParameterStoreIface
	secretsStore   config.ParameterStoreIface

	connectionEstablished bool
)

type PgClientIface interface {
	AddHashedAppId(hash hashBuffer, appId string) error
	Ping() error
}

type DBCredentialsType struct {
	DbName   string `json:"dbname"`
	Host     string `json:"host"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	Username string `json:"username"`
}

type ParameterStoreType struct {
	Name  string
	Value string
}

type s3Client struct{ *s3.S3 }

func CopyFile(svc s3iface.S3API, transactionId, appId string, opCounters *OperationalCounters, slackTimestamp, bucket, destinationPath, sourceKey string, isCancelled bool,
	source, suppressionWindow, operation, mmp string) error {
	destinationKey := filepath.Join(destinationPath, filepath.Base(sourceKey))
	_, err := svc.CopyObject(&s3.CopyObjectInput{
		CopySource:        aws.String(url.QueryEscape(filepath.Join(bucket, sourceKey))),
		Bucket:            aws.String(bucket),
		Key:               aws.String(destinationKey),
		StorageClass:      &s3.StorageClass_Values()[1], //  StorageClassReducedRedundancy
		ContentType:       aws.String(DEFAULT_CONTENT_TYPE),
		MetadataDirective: aws.String(s3.MetadataDirective_Values()[1]),
		Metadata: map[string]*string{
			DEFAULT_APP_ID:               aws.String(appId),
			DEFAULT_TRANSACTION_ID:       aws.String(transactionId),
			DEFAULT_TOTAL_HASH_OUT:       aws.String(strconv.Itoa(opCounters.UuidCount)),
			DEFAULT_ADID_TOTAL_NEW:       aws.String(strconv.Itoa(opCounters.NewAdIdCount)),
			DEFAULT_APPID_TOTAL_NEW:      aws.String>(strconv.Itoa(opCounters.SelectError)),
			DEFAULT_ERROR_UPDATE:         aws.String(strconv.Itoa(opCounters.UpdateError)),
			DEFAULT_ERROR_LOG_INSERT:     aws.String(strconv.Itoa(opCounters.SuppressionLogInsertError)),
			DEFAULT_ERROR_TRANSACTION:    aws.String(strconv.Itoa(opCounters.TransactionError)),
			DEFAULT_SLACK_TIMESTAMP:      aws.String(slackTimestamp),
			DEFAULT_INGEST_CANCELLED:     aws.String(strconv.FormatBool(isCancelled)),
			DEFAULT_SOURCE:               aws.String(source),
			DEFAULT_SUPPRESSION_WINDOW:   aws.String(suppressionWindow),
			DEFAULT_OPERATION:            aws.String(operation),
			DEFAULT_MMP:                  aws.String(mmp)}})
	if err != nil {
		return fmt.Errorf("Unable to copy item from: %q, to %q, err: %v",
			filepath.Join(bucket, sourceKey),
			filepath.Join(bucket, destinationKey),
			err)
	}
	err = svc.WaitUntilObjectExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(destinationKey)})
	if err != nil {
		return fmt.Errorf("Error occurred while waiting for item %q to be copied to %q, %v",
			filepath.Join(bucket, sourceKey),
			filepath.Join(bucket, destinationKey),
			err)
	}

	return nil
}

func DeleteFile(bucket, key string) error {
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key)})
	if err != nil {
		return fmt.Errorf("Could not delete %q, %v", filepath.Join(bucket, key), err)
	}
	return nil
}

func MoveFiles(svc s3iface.S3API, objects []events.S3EventRecord, transactionId, appId string, opCounters *OperationalCounters, slackTimestamp, bucket, destination string, isCancelled bool,
	source, suppressionWindow, operation, mmp string) {
	for _, obj := range objects {
		key, err := url.QueryUnescape(obj.S3.Object.Key)
		if err != nil {
			logger. // File not found
				NewEntry(log_msg_file_read_error, LogMsgTxt).
				Error(err.Error(), "Unreadable filename")
			_, _ = badger.Notify(err.Error(),
				honeybadger.Context{
					"ERROR_NUMBER": log_msg_file_read_error,
					"ERROR_TEXT":   LogMsgTxt[log_msg_file_read_error],
					"ERROR_URL":    "https://github.com/influencemobile/go-lambdas/tree/master/chunk#types-of-log-events",
				},
				honeybadger.Tags{
					env.ApplicationName,
					env.Environment,
					appId,
					LogMsgTxt[log_msg_file_read_error]})
			badger.Flush()
			continue
		}
		if err := CopyFile(svc, transactionId, appId, opCounters, slackTimestamp, bucket, destination, key, isCancelled,
			source, suppressionWindow, operation, mmp); err != nil {
			logger. // File not copied
				NewEntry(log_msg_file_copy_error, LogMsgTxt).
				Error(err.Error(), "File copy error.")
			_, _ = badger.Notify(err.Error(),
				honeybadger.Context{
					"ERROR_NUMBER": log_msg_file_copy_error,
					"ERROR_TEXT":   LogMsgTxt[log_msg_file_copy_error],
					"ERROR_URL":    "https://github.com/influencemobile/go-lambdas/tree/master/chunk#types-of-log-events",
				},
				honeybadger.Tags{
					env.ApplicationName,
					env.Environment,
					appId,
					LogMsgTxt[log_msg_file_copy_error]})
			badger.Flush()
			continue
		}
		if err := DeleteFile(bucket, key); err != nil {
			logger. // File not deleted
				NewEntry(log_msg_file_delete_error, LogMsgTxt).
				Error(err.Error(), "File delete error.")
			_, _ = badger.Notify(err.Error(),
				honeybadger.Context{
					"ERROR_NUMBER": log_msg_file_delete_error,
					"ERROR_TEXT":   LogMsgTxt[log_msg_file_delete_error],
					"ERROR_URL":    "https://github.com/influencemobile/go-lambdas/tree/master/chunk#types-of-log-events",
				},
				honeybadger.Tags{
					env.ApplicationName,
					env.Environment,
					appId,
					LogMsgTxt[log_msg_file_delete_error]})
			badger.Flush()
			continue
		}
	}
}

func WriteRecords(threadNum int, wg *sync.WaitGroup, appId string, chHash chan hashBuffer) {
	defer wg.Done()

	log_every := env.Records
	var goodCount, badCount int
	for h := range chHash {
		if err := pgDb.AddHashedAppId(h, appId); err != nil {
			badCount++
			if badCount%log_every == 0 {
				logger.
					NewEntry(log_msg_processing, LogMsgTxt).
					WithFields(Fields{"HASHES_SKIPPED": log_every}).
					Debug()
			}
			continue
		}
		goodCount++
		if goodCount%log_every == 0 {
			logger.
				NewEntry(log_msg_processing, LogMsgTxt).
				WithFields(Fields{"HASHES_PROCESSED": log_every}).
				Debug()
		}
	}
}

func ValidateHashedAdId(s string) ([]byte, error) {
	hex, err := hex.DecodeString(s) // Ensure hash is actually in hex.
	if err != nil {
		return nil, err
	}
	if len(hex) != 32 { // Ensure hash is correct length.
		return nil, fmt.Errorf("Hash incorrect length")
	}
	return hex, nil
}

func ProcessHashedAdIds(cout chan hashBuffer, r io.Reader) {
	var hashCount int
	var malformedHashCount int
	var buf [32]byte
	log_every := env.Records

	s := bufio.NewScanner(r)
	for s.Scan() {
		hash, err := ValidateHashedAdId(s.Text())
		if err != nil {
			malformedHashCount++
			continue
		}
		copy(buf[:], hash)
		cout <- buf
		hashCount++
		if hashCount%log_every == 0 && hashCount > 0 {
			logger.
				NewEntry(log_msg_processing, LogMsgTxt).
				WithFields(Fields{"HASHES_READ": hashCount}).
				Debug()
		}
	}
	opCounters.incUuid(hashCount)
	opCounters.incHashError(malformedHashCount)

	if s.Err() != nil {
		logger. // File IO error
			NewEntry(log_msg_io, LogMsgTxt).
			Error(s.Err().Error())
	}
}

func ProcessFiles(svc s3iface.S3API, chUuid chan hashBuffer, objects []events.S3EventRecord, appId string) {
	defer close(chUuid)

	for _, obj := range objects {
		key, err := url.QueryUnescape(obj.S3.Object.Key)
		if err != nil {
			logger. // File not found
				NewEntry(log_msg_file_read_error, LogMsgTxt).
				Error(err.Error(), "Unreadable filename")
			_, _ = badger.Notify(err.Error(),
				honeybadger.Context{
					"ERROR_NUMBER": log_msg_file_read_error,
					"ERROR_TEXT":   LogMsgTxt[log_msg_file_read_error],
					"ERROR_URL":    "https://github.com/influencemobile/go-lambdas/tree/master/chunk#types-of-log-events",
				},
				honeybadger.Tags{
					env.ApplicationName,
					env.Environment,
					appId,
					LogMsgTxt[log_msg_file_read_error]})
			badger.Flush()
			continue
		}

		result, err := svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(obj.S3.Bucket.Name),
			Key:    aws.String(key)})
		if err != nil {
			logger. // File not found
				NewEntry(log_msg_file_read_error, LogMsgTxt).
				Error(err.Error(), key)
			_, _ = badger.Notify(err.Error(),
				honeybadger.Context{
					"ERROR_NUMBER": log_msg_file_read_error,
					"ERROR_TEXT":   LogMsgTxt[log_msg_file_read_error],
					"ERROR_URL":    DEFAULT_ERROR_URL,
				},
				honeybadger.Tags{
					env.ApplicationName,
					env.Environment,
					appId,
					LogMsgTxt[log_msg_file_read_error]})
			badger.Flush()
			continue
		}
		defer result.Body.Close()

		// Make sure that the App_id hasn't been changed out from under us.
		if result.Metadata == nil {
			logger.
				WithFields(Fields{"APPID": "No Metadata found"}).
				Error()
			_, _ = badger.Notify(&AppIdError{
				Type: "APPID",
				Err:  fmt.Errorf("No metadata found")},
				honeybadger.Context{
					"ERROR_NUMBER": log_msg_object_metadata_missing,
					"ERROR_TEXT":   LogMsgTxt[log_msg_object_metadata_missing],
					"ERROR_URL":    DEFAULT_ERROR_URL,
				},
				honeybadger.Tags{
					env.ApplicationName,
					env.Environment,
					appId,
					"APPID"})
			badger.Flush()
			return
		}

		id, exists := result.Metadata[DEFAULT_APP_ID]
		if !exists {
			logger.
				WithFields(Fields{"APPID": fmt.Sprintf("%s not found in metadata. Verify that the first letter is uppercase?", DEFAULT_APP_ID)}).
				Error()
			_, _ = badger.Notify(&AppIdError{
				Type: "APPID",
				Err:  fmt.Errorf("%s not found in metadata.", appId)},
				honeybadger.Context{
					"ERROR_NUMBER": log_msg_app_id_missing,
					"ERROR_TEXT":   LogMsgTxt[log_msg_app_id_missing],
					"ERROR_URL":    DEFAULT_ERROR_URL,
				},
				honeybadger.Tags{
					env.ApplicationName,
					env.Environment,
					appId,
					"APPID"})
			badger.Flush()
			return
		}
		if appId != *id {
			logger.
				WithFields(Fields{"APPID": fmt.Sprintf("App_id mismatch. Expected %s, got %s", appId, *id)}).
				Error()
			_, _ = badger.Notify(&AppIdError{
				Type: "APPID",
				Err:  fmt.Errorf("Mistamth. Expected %s, got %s", appId, *id)},
				honeybadger.Context{
					"ERROR_NUMBER": log_msg_app_id_mismatch,
					"ERROR_TEXT":   LogMsgTxt[log_msg_app_id_mismatch],
					"ERROR_URL":    DEFAULT_ERROR_URL,
				},
				honeybadger.Tags{
					env.ApplicationName,
					env.Environment,
					appId,
					"APPID"})
			badger.Flush()
			return
		}

		ProcessHashedAdIds(chUuid, result.Body)
	}
}

func ValidateMetadata(svc s3iface.S3API, objects []events.S3EventRecord, metadata string) (string, error) {
	var val string
	for _, file := range objects {
		key, err := url.QueryUnescape(file.S3.Object.Key)
		if err != nil {
			return "", &FileError{
				Type: "FILE",
				Err:  fmt.Errorf("File not found: %s", err.Error()),
			}
		}
		id, err := RetrieveMetadata(svc, file.S3.Bucket.Name, key, 0, metadata)
		if err != nil {
			return "", err
		}
		if val == "" {
			val = id
		}
		if val != id {
			return "", &MetadataError{
				Type: metadata,
				Err:  fmt.Errorf("Mismatch. Expected %s, got %s", val, id)}
		}
	}
	return val, nil
}

func RetrieveMetadata(svc s3iface.S3API, bucket, key string, retry int, metadata string) (string, error) {
	var count int
	var err error
	var val *string = new(string)
	// Loop "retry" times. Return immediately when metadata value is retrieved, or the
	// last error when retries are exhausted.
	for {
		if count > retry {
			break
		}
		var result *s3.HeadObjectOutput
		result, err = svc.HeadObject(&s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key)})
		if err != nil { // Error retrieving object. Don't bother retrying.
			return "", err
		}

		if result.Metadata == nil { // Retry in this case.
			err = &MetadataError{
				Type: metadata,
				Err:  fmt.Errorf("No metadata found")}
			if retry > 0 {
				opCounters.setRetryAttempts(count) // Set this so we can report to Slack
				logger.WithFields(Fields{
					"METADATA_RETRY": opCounters.RetryAttempts,
					metadata:         metadata}).
					Debug()
				time.Sleep(DEFAULT_RETRY_DURATION * time.Second)
			}
			count++
			continue
		}
		var exists bool
		val, exists = result.Metadata[metadata]
		if !exists { // Retry in this case
			val = new(string)
			err = &MetadataError{
				Type: metadata,
				Err:  fmt.Errorf("%s not found in metadata.", metadata)}
			if retry > 0 {
				opCounters.setRetryAttempts(count) // Set this so we can report to Slack
				logger.WithFields(Fields{
					"METADA_RETRY": opCounters.RetryAttempts,
					metadata:       metadata}).
					Debug()
				time.Sleep(DEFAULT_RETRY_DURATION * time.Second)
			}
			count++
			continue
		}
		// If we are here, it means we have retrieved the metadata field. Any previous
		// error should be discarded. Exit the loop immediately.
		err = nil
		break
	}
	return *val, err
}

func ValidateSQSEvents(event *events.SQSEvent) error {
	// We need a valid s3Event containing at least one record.
	if event == nil {
		return &SQSEventError{
			Type: "SQSEvent",
			Err:  fmt.Errorf("No SQSEvent"),
		}
	}
	if event.Records == nil {
		return &SQSEventError{
			Type: "SQSEvent",
			Err:  fmt.Errorf("No SQSEvent records"),
		}
	}
	return nil
}

func ValidateS3Events(e string, s3Event *events.S3Event) error {
	// We need a valid s3Event containing at least one record.
	if s3Event == nil {
		return &S3EventError{
			Type: "s3Event",
			Err:  fmt.Errorf("No s3Event"),
		}
	}
	if s3Event.Records == nil {
		return &S3EventError{
			Type: "s3Event",
			Err:  fmt.Errorf("No s3Event records"),
		}
	}

	for _, record := range s3Event.Records {
		if e != record.EventName {
			return &S3EventError{
				Type: "s3Event",
				Err:  fmt.Errorf("Expected %s, got %s", e, record.EventName),
			}
		}
	}
	return nil
}

func TestBucket(bucket string) error {
	_, err := svc.HeadBucket(&s3.HeadBucketInput{Bucket: aws.String(bucket)})
	if err != nil {
		return fmt.Errorf("Cannot access bucket: %q, %v", bucket, err)
	}
	return nil
}

func TestPath(bucket, path string) error {
	_, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path)})
	if err != nil {
		return fmt.Errorf("Cannot access path: %q, %v", filepath.Join(bucket, path), err)
	}
	return nil
}

func HandleRequest(ctx context.Context, sqsEvent *events.SQSEvent) {
	// Start logging
	logger.
		NewEntry(log_msg_begin, LogMsgTxt).
		Info()
	// Log environment configuration parameters
	logger.NewEntry(log_msg_envars, LogMsgTxt).
		WithStruct(env).
		Debug()
	// Log lambda parameters
	logger.
		NewEntry(log_msg_processing, LogMsgTxt).
		Debug(sqsEvent)

	opCounters = &OperationalCounters{RWMutex: sync.RWMutex{}}

	// Attempt to load the parameters file written by the lambda extension.
	if err := parameterStore.ReadConfig(); err != nil {
		logger.
			NewEntry(log_msg_parameters_store_error, LogMsgTxt).
			Error(err.Error())
		os.Exit(1)
	}

	// configure honeybadger
	badgerParams, err := parameterStore.ParameterToBytes(DEFAULT_HONEYBADGER_KEY)
	if err != nil {
		logger.
			NewEntry(log_msg_parameters_store, LogMsgTxt).
			Error(err.Error())
		os.Exit(1)
		// params := fmt.Sprintf("{\"name\":\"/honeybadger/engage/api_key\",\"value\":\"%s\"}", env.honeybadgerApiKey)
		// badgerParams = []byte(params)
	}
	badgerApiKey := ParameterStoreType{}
	_ = json.Unmarshal(badgerParams, &badgerApiKey)
	badger.Configure(honeybadger.Configuration{APIKey: badgerApiKey.Value, Env: env.Environment})

	// Attempt to load the secrets file written by the lambda extension.
	if err := secretsStore.ReadConfig(); err != nil {
		logger.
			NewEntry(log_msg_secrets_store_error, LogMsgTxt).
			Error(err.Error())
		os.Exit(1)
	}

	dbSecrets, err := secretsStore.ParameterToBytes(fmt.Sprintf("dbuser/%s/engage/lambda_dsl", strings.ToLower(env.Environment)))
	if err != nil {
		logger.
			NewEntry(log_msg_secrets_store_error, LogMsgTxt).
			Error(err.Error())
		os.Exit(1)
		// params := fmt.Sprintf("{\"name\":\"dbuser/%s/engage/lambda_dsl\",\"value\":{\"dbname\":\"%s\",\"host\":\"%s\",\"password\":,\"%s\",\"port\":,\"%d\",\"username\":,\"%s\"}",
		// 	env.Environment, env.DbName, env.DbHost, env.dbPassword, env.DbPort, env.DbUser)
		// badgerParams = []byte(params)
	}
	dbCreds := DBCredentialsType{}
	_ = json.Unmarshal(dbSecrets, &dbCreds)
	badger.Configure(honeybadger.Configuration{APIKey: badgerApiKey.Value, Env: env.Environment})

	if !connectionEstablished {
		connectionEstablished = true
		logger.Debug("DATABASE", "connecting")
		// Connect to the database
		var err error
		pgDb, err = OpenDatabase(dbCreds.Host, dbCreds.Port, dbCreds.Username, dbCreds.Password, dbCreds.DbName, env.DbSsl)
		if err != nil {
			logger.
				WithFields(Fields{"DATABASE": err.Error()}).
				Error(LogMsgTxt[log_msg_pg_database_connect])
			_, _ = badger.Notify(err.Error(),
				honeybadger.Tags{
					env.ApplicationName,
					env.Environment,
					"DATABASE"})
			// TODO: Add context
			badger.Flush()
			os.Exit(1)
		}
		logger.Debug("DATABASE", "connected")
	}

	if err := pgDb.Ping(); err != nil {
		logger.
			WithFields(Fields{"DATABASE": err.Error()}).
			Error(LogMsgTxt[log_msg_pg_database_connect])
		_, _ = badger.Notify(err.Error(),
			honeybadger.Tags{
				env.ApplicationName,
				"Database",
				"PostgreSQL"})
		badger.Flush()
		os.Exit(1)
	}

	// // Bomb out early if we cannot access the "S3_DESTINATION_BUCKET".
	// if err := TestBucket(env.DestinationBucket); err != nil {
	// 	logger.
	// 		NewEntry(log_msg_processing, LogMsgTxt).
	// 		WithFields(Fields{"S3_DESTINATION_BUCKET": err.Error()}).
	// 		Error()
	// 	_, _ = badger.Notify(err.Error(),
	// 		honeybadger.Context{
	// 			"ERROR_NUMBER": log_msg_bucket_access_error,
	// 			"ERROR_TEXT":   LogMsgTxt[log_msg_bucket_access_error],
	// 			"ERROR_URL":    "https://github.com/influencemobile/go-lambdas/tree/master/chunk#types-of-log-events",
	// 		},
	// 		honeybadger.Tags{
	// 			env.ApplicationName,
	// 			env.Environment,
	// 			"S3_DESTINATION_BUCKET"})
	// 	badger.Flush()
	// 	return
	// }

	// // Bomb out early if we cannot access the "S3_COMPLETED_PATH"
	// if err := TestPath(env.DestinationBucket, env.CompletedPath); err != nil {
	// 	logger.
	// 		NewEntry(log_msg_processing, LogMsgTxt).
	// 		WithFields(Fields{"S3_COMPLETED_PATH": err.Error()}).
	// 		Error()
	// 	_, _ = badger.Notify(err.Error(),
	// 		honeybadger.Context{
	// 			"ERROR_NUMBER": log_msg_filepath_access_error,
	// 			"ERROR_TEXT":   LogMsgTxt[log_msg_filepath_access_error],
	// 			"ERROR_URL":    "https://github.com/influencemobile/go-lambdas/tree/master/chunk#types-of-log-events",
	// 		},
	// 		honeybadger.Tags{
	// 			env.ApplicationName,
	// 			env.Environment,
	// 			"S3_COMPLETED_PATH"})
	// 	badger.Flush()
	// 	return
	// }

	var isCancelled bool = false

	// Validate that the sqsEvent structure contains at least one event.
	if err := ValidateSQSEvents(sqsEvent); err != nil {
		logger.
			NewEntry(log_msg_processing, LogMsgTxt).
			WithFields(Fields{"SQSEVENT": err.Error()}).
			Error()
		_, _ = badger.Notify(err.Error(),
			honeybadger.Context{
				"ERROR_NUMBER": log_msg_sqsevent,
				"ERROR_TEXT":   LogMsgTxt[log_msg_sqsevent],
				"ERROR_URL":    DEFAULT_ERROR_URL,
			},
			honeybadger.Tags{
				env.ApplicationName,
				env.Environment,
				"SQSEVENT"})
		badger.Flush()
		return
	}

	var transactionId, appId, slackTimestamp, source, suppressionWindow, operation, mmp string

	// Each SQS message should contain one or more S3Event records.
	for _, record := range sqsEvent.Records {
		var s3Event events.S3Event
		// An error means the S3 event is not available
		if err := json.Unmarshal([]byte(record.Body), &s3Event); err != nil {
			logger.
				NewEntry(log_msg_processing, LogMsgTxt).
				WithFields(Fields{"S3EVENT": err.Error()}).
				Error()
			_, _ = badger.Notify(err.Error(),
				honeybadger.Context{
					"ERROR_NUMBER": log_msg_s3event,
					"ERROR_TEXT":   LogMsgTxt[log_msg_s3event],
					"ERROR_URL":    DEFAULT_ERROR_URL,
				},
				honeybadger.Tags{
					env.ApplicationName,
					env.Environment,
					"S3EVENT"})
			badger.Flush()
			return
		}
		// Validate that the s3Event structure before continuing. If there are
		// multiple events, make sure all have the same event type, etc.
		if err := ValidateS3Events(s3Event_type, &s3Event); err != nil {
			logger.
				NewEntry(log_msg_processing, LogMsgTxt).
				WithFields(Fields{"S3EVENT": err.Error()}).
				Error()
			_, _ = badger.Notify(err.Error(),
				honeybadger.Context{
					"ERROR_NUMBER": log_msg_s3event,
					"ERROR_TEXT":   LogMsgTxt[log_msg_s3event],
					"ERROR_URL":    DEFAULT_ERROR_URL,
				},
				honeybadger.Tags{
					env.ApplicationName,
					env.Environment,
					"S3EVENT"})
			badger.Flush()
			return
		}

		// If there are multiple files, ensure all have the same App_id.
		var err error
		appId, err = ValidateMetadata(svc, s3Event.Records, DEFAULT_APP_ID)
		if err != nil {
			logger.
				NewEntry(log_msg_processing, LogMsgTxt).
				WithFields(Fields{"APPID": err.Error()}).
				Error()
			_, _ = badger.Notify(err.Error(),
				honeybadger.Context{
					"ERROR_NUMBER": log_msg_app_id_mismatch,
					"ERROR_TEXT":   LogMsgTxt[log_msg_app_id_mismatch],
					"ERROR_URL":    DEFAULT_ERROR_URL,
				},
				honeybadger.Tags{
					env.ApplicationName,
					env.Environment,
					"APPID"})
			badger.Flush()
			return
		}

		// If there are multiple files, ensure all have the same Transaction_id
		transactionId, err = ValidateMetadata(svc, s3Event.Records, DEFAULT_TRANSACTION_ID)
		if err != nil {
			logger.
				NewEntry(log_msg_processing, LogMsgTxt).
				WithFields(Fields{"TRANSACTIONID": err.Error()}).
				Error()
			_, _ = badger.Notify(err.Error(),
				honeybadger.Context{
					"ERROR_NUMBER": log_msg_transaction_id_mismatch,
					"ERROR_TEXT":   LogMsgTxt[log_msg_transaction_id_mismatch],
					"ERROR_URL":    DEFAULT_ERROR_URL,
				},
				honeybadger.Tags{
					env.ApplicationName,
					env.Environment,
					"APPID"})
			badger.Flush()
			return
		}

		// If there are multiple files, ensure all have the same Slack_timestamp
		slackTimestamp, err = ValidateMetadata(svc, s3Event.Records, DEFAULT_SLACK_TIMESTAMP)
		if err != nil {
			logger.
				NewEntry(log_msg_processing, LogMsgTxt).
				WithFields(Fields{DEFAULT_SLACK_TIMESTAMP: err.Error()}).
				Error()
			_, _ = badger.Notify(err.Error(),
				honeybadger.Context{
					"ERROR_NUMBER": log_msg_transaction_id_mismatch,
					"ERROR_TEXT":   LogMsgTxt[log_msg_transaction_id_mismatch],
					"ERROR_URL":    DEFAULT_ERROR_URL,
				},
				honeybadger.Tags{
					env.ApplicationName,
					env.Environment,
					DEFAULT_SLACK_TIMESTAMP})
			badger.Flush()
		}

		logger. // Log the app ID
			NewEntry(log_msg_processing, LogMsgTxt).
			WithFields(Fields{DEFAULT_APP_ID: appId}).
			Debug()

		cancelledStr, _ := ValidateMetadata(svc, s3Event.Records, DEFAULT_INGEST_CANCELLED)
		if cancelledStr == "true" {
			isCancelled = true
		}

		suppressionWindow, err = ValidateMetadata(svc, s3Event.Records, DEFAULT_SUPPRESSION_WINDOW)
		if err != nil {
			logger.
				NewEntry(log_msg_ingest_suppression_window_error, LogMsgTxt).
				WithFields(Fields{DEFAULT_SUPPRESSION_WINDOW: err.Error()}).
				Error()
		}
		if suppressionWindow == "" {
			suppressionWindow = "0"
		}

		operation, err = ValidateMetadata(svc, s3Event.Records, DEFAULT_OPERATION)
		if err != nil {
			logger.
				NewEntry(log_msg_ingest_operation_error, LogMsgTxt).
				WithFields(Fields{DEFAULT_OPERATION: err.Error()}).
				Error()
		}

		if operation == "" {
			operation = "Add"
		}

		mmp, err = ValidateMetadata(svc, s3Event.Records, DEFAULT_MMP)
		if err != nil {
			logger.
				NewEntry(log_msg_ingest_mmp_error, LogMsgTxt).
				WithFields(Fields{DEFAULT_MMP: err.Error()}).
				Error()
		}

		source, err = ValidateMetadata(svc, s3Event.Records, DEFAULT_SOURCE)
		if err != nil {
			logger.
				NewEntry(log_msg_ingest_source_error, LogMsgTxt).
				WithFields(Fields{DEFAULT_SOURCE: err.Error()}).
				Error()
		}

		if source == "" {
			source = "Manual_bulk"
		}

	}

	if !isCancelled {
		// Create a channel to route hashes to worker threads.
		hashChannel := make(chan hashBuffer, DEFAULT_CHANNEL_BUFFERS)

		// Spin up the chunking threads. Each thread reads/blocks on the channel.
		var wg sync.WaitGroup
		var count int = 0
		for i := 0; i < env.Threads; i++ {
			wg.Add(1)
			go WriteRecords(count, &wg, appId, hashChannel)
			count++
		}

		for _, record := range sqsEvent.Records {
			var s3Event events.S3Event
			_ = json.Unmarshal([]byte(record.Body), &s3Event)
			// Read the AdIds from the input file(s)
			ProcessFiles(svc, hashChannel, s3Event.Records, appId)
		}

		// Wait until the chunking threads are done chunking
		wg.Wait()
	}

	for _, record := range sqsEvent.Records {
		var s3Event events.S3Event
		_ = json.Unmarshal([]byte(record.Body), &s3Event)

		// The metadata may have been updated while we were processing the records,
		// so check again to retrieve the latest value.
		if !isCancelled {
			cancelledStr, _ := ValidateMetadata(svc, s3Event.Records, DEFAULT_INGEST_CANCELLED)
			if cancelledStr == "true" {
				isCancelled = true
			}
		}
		// Move the input files to the completion directory.
		MoveFiles(svc, s3Event.Records, transactionId, appId, opCounters, slackTimestamp, env.DestinationBucket, env.CompletedPath, isCancelled,
			source, suppressionWindow, operation, mmp)
	}

	if opCounters.InsertError > 0 ||
		opCounters.SelectError > 0 ||
		opCounters.UpdateError > 0 ||
		opCounters.Utf8Error > 0 ||
		opCounters.HashError > 0 {
		_, _ = badger.
			Notify(&DatabaseError{
				Type: "DATABASE",
				Err: fmt.Errorf(
					"INSERT_ERROR: %d, SELECT_ERROR: %d, UPDATE_ERROR: %d, UTF8_ERROR: %d, MALFORMED_HASH_ERROR: %d",
					opCounters.InsertError, opCounters.SelectError,
					opCounters.UpdateError, opCounters.Utf8Error, opCounters.HashError)},
				honeybadger.Tags{
					env.ApplicationName,
					env.Environment,
					"DATABASE"},
				honeybadger.Context{
					"ERROR_NUMBER": log_msg_database,
					"ERROR_TEXT":   LogMsgTxt[log_msg_database],
					"ERROR_URL":    DEFAULT_ERROR_URL,
				})
		badger.Flush()
	}

	logger. // Log the operational metrics
		NewEntry(log_msg_end, LogMsgTxt).
		WithStruct(opCounters).
		Info()
}

func init() {
	// Retrieve environment variables from AWS. Log level, number of threads,
	// etc.
	env = NewEnv()

	// Create a new logger.
	logger = NewLoggers(os.Stderr, env.LogLevel, 0).SetFields(
		Fields{
			EnvVars[ENV_APPLICATION_NAME]: env.ApplicationName,
			EnvVars[ENV_JOB_ID]:           env.JobId,
		})
	badger = &badgerClient{}

	parameterStore = &config.KeyClient{Filename: env.ParametersStore}
	secretsStore = &config.KeyClient{Filename: env.SecretsStore}

	session, _ := session.NewSession()
	svc = s3Client{s3.New(session)}
}

func main() {
	defer badger.Monitor() // Report all unhandled panics into HoneyBadger
	lambda.Start(HandleRequest)
}
