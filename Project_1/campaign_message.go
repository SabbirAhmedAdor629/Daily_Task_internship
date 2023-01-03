package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
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
	im_badger "influencemobile.com/badger"
	. "influencemobile.com/logging"
	config "influencemobile.com/parameter_store"
)

var (
	logger                *Loggers
	env                   envVars
	secretsStore          config.ParameterStoreIface
	parameterStore        config.ParameterStoreIface
	badger                im_badger.BadgerIface
	pgDb                  PgClientIface
	connectionEstablished bool
	svc                   s3iface.S3API
	opCounters            *OperationalCounters
)

type (
	hashBuffer [32]byte
)

type HoneyBadgerType struct {
	Name  string
	Value string
}

type PgClientIface interface {
	Ping() error
	UpdateWorkLog(campaign_id int, status string) error
	ReadCampaignData(campaign_id int) (*CampaignTable, error)
	ReadMessageTemplateData(message_template_id int) (*MessageTemplate, error)
}

type DBCredentialsType struct {
	DbName   string `json:"dbname"`
	Host     string `json:"host"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	Username string `json:"username"`
}

type s3Client struct{ *s3.S3 }

type SQSPayload struct {
	Origin      string `json:"origin"`
	operation   string `json:"operation"`
	SubmittedTs string `json:"submitted_ts"`
	Data        Data   `json:"data"`
}

type Data struct {
	S3File     string `json:"s3_file"`
	CampaignId int    `json:"campaign_id"`
	LogChunkId int    `json:"log_chunk_id"`
}

func contains(elements []interface{}, v string) bool {
	for _, s := range elements {
		if v == s {
			return true
		}
	}
	return false
}

func timeToDeliver(campaignData *CampaignTable) bool {

	scheduleTime := int(campaignData.ScheduleTime.Int32)

	timeZone := getScheduleTimeZone(scheduleTime)

	if timeZone == "Prime" {
		return true
	} else if timeZone == "User" {
		return userPresentsInScheduleTime(scheduleTime)
	}

	changedTimeZone := "PST8PDT" // "Pacific Time (US & Canada)"
	if timeZone == "GMT" {
		changedTimeZone = "UTC"
	}
	loc, _ := time.LoadLocation(changedTimeZone)
	currentDay := time.Now().In(loc).Weekday().String()
	logger.Debug("Current Day", currentDay)

	if len(campaignData.ScheduleDays) > 0 || contains(campaignData.ScheduleDays, currentDay) {
		locPT, _ := time.LoadLocation("PST8PDT") // "Pacific Time (US & Canada)
		pt := time.Now().In(locPT).Format("03:04am")
		utc := time.Now().UTC().Format("03:04am")

		if campaignData.ScheduleTime.Valid {
			return true
		} else if timeZone == "PT" && getScheduleTimeOfDay(scheduleTime) == pt {
			return true
		} else if timeZone == "GMT" && getScheduleTimeOfDay(scheduleTime) == utc {
			return true
		}
	}
	return false
}

func userPresentsInScheduleTime(scheduleTimeIndex int) bool {
	offsetRange := getOffsetRange()
	userTimes := []interface{}{}
	for _, tzOffset := range offsetRange {
		userTimes = append(userTimes, userCurrentTime(tzOffset))
	}
	logger.Debug("User Times List [1] : ", userTimes)
	logger.Debug("User Time [1] : ", getScheduleTimeOfDay(scheduleTimeIndex))
	if contains(userTimes, getScheduleTimeOfDay(scheduleTimeIndex)) {
		return true
	}
	return false
}

func getOffsetRange() []int {
	locFirst, _ := time.LoadLocation("America/Adak")
	locLast, _ := time.LoadLocation("America/New_York")
	_, offsetFirst := time.Now().In(locFirst).Zone()
	_, offsetLast := time.Now().In(locLast).Zone()
	offsetFirst, offsetLast = offsetFirst/3600, offsetLast/3600
	offsetRange := []int{}
	for i := offsetFirst; i <= offsetLast; i++ {
		offsetRange = append(offsetRange, i)
	}
	return offsetRange
}

func userTimeToDeliver(player_id string, scheduleDays ScheduleDays) (bool, error) {
	if len(scheduleDays) == 0 {
		return true, nil
	}

	if player_id == "" {
		return false, fmt.Errorf("player_id can not be null")
	}

	playerId, err := strconv.Atoi(player_id)
	if err != nil {
		return false, fmt.Errorf("player_id must be of type int")
	}

	// sample playerHash data
	samplePlayerHash := map[string]interface{}{
		"player_id":        playerId,
		"time_zone_offset": -10,
	}

	playerHash := samplePlayerHash // placeholder :  Im::Targeting::Target.player(player_id)

	if len(playerHash) == 0 || !availableKeyInMap(playerHash, "time_zone_offset") {
		return false, nil
	}

	userDay := time.Now().UTC().Add(time.Hour * (time.Duration(playerHash["time_zone_offset"].(int)))).Weekday()
	logger.Debug("User Current Day [3] : " + userDay.String())
	logger.Debug("User Schedule Day [3] : ", scheduleDays)
	if contains(scheduleDays, userDay.String()) {
		return true, nil
	}
	return false, nil
}

func primeTimeToDeliver(player_id string, campaignTableData *CampaignTable) (bool, error) {

	scheduleTime := int(campaignTableData.ScheduleTime.Int32)

	if len(campaignTableData.ScheduleDays) == 0 {
		return true, nil
	}

	if player_id == "" {
		return false, fmt.Errorf("player_id can not be null")
	}

	playerId, err := strconv.Atoi(player_id)
	if err != nil {
		return false, fmt.Errorf("player_id must be of type int")
	}

	// sample playerHash data
	// player_hash = Im::Targeting::Target.player(player_id)
	playerHash := map[string]interface{}{
		"player_id": playerId,
		"prime_time": map[string]string{
			"best": fmt.Sprint(time.Now().UTC().Hour()),
		},
		"time_zone_offset": -10,
	}

	if len(playerHash) == 0 || !availableKeyInMap(playerHash, "prime_time", "time_zone_offset") {
		return false, nil
	}

	if len(campaignTableData.ScheduleDays) != 0 {
		primeDay := time.Now().UTC().Add(time.Hour * (time.Duration(playerHash["time_zone_offset"].(int)))).Weekday()
		logger.Debug("Prime Day current : " + primeDay.String())
		logger.Debug("Campaign Schedule Days : ", campaignTableData.ScheduleDays)
		if contains(campaignTableData.ScheduleDays, primeDay.String()) {
			return true, nil
		}
	}

	primeTime := getPrimeTime(scheduleTime)
	primeValue, ok := playerHash["prime_time"].(map[string]string)[primeTime]
	if !ok {
		return false, nil
	}

	primeTimeValue, err := strconv.Atoi(primeValue)
	if err != nil {
		return false, fmt.Errorf("player_id must be of type int")
	}

	localHour := userLocalHour(primeTimeValue, scheduleTime)
	currentHour := userCurrentHour(scheduleTime)
	logger.Debug("User Local Prime Hour: ", localHour)
	logger.Debug("User Current Hour: ", currentHour)
	if currentHour == localHour {
		return true, nil
	}

	return false, nil
}

func availableKeyInMap(mapData map[string]interface{}, keyList ...string) bool {
	for _, key := range keyList {
		if _, ok := mapData[key]; !ok {
			return false
		}
	}
	return true
}

func schedulePlayerId(campaignTable *CampaignTable, sqsData SQSPayload, player_id string) (bool, error) {
	scheduleTime := int(campaignTable.ScheduleTime.Int32)
	scheduleTimeZone := getScheduleTimeZone(scheduleTime)
	fmt.Println("Campaign Schedule Time Zone [2] : ", scheduleTimeZone)
	if scheduleTimeZone == "User" {
		inUserTime, err := userTimeToDeliver(player_id, campaignTable.ScheduleDays)
		if err != nil {
			return false, err
		}
		if !inUserTime {
			opCounters.incOutOfUserTime(1)
			return false, nil
		}
	} else if scheduleTimeZone == "Prime" {
		inPrimeTime, err := primeTimeToDeliver(player_id, campaignTable)
		if err != nil {
			return false, err
		}
		logger.Debug("inPrimeTime : " + fmt.Sprint(inPrimeTime))
	}
	return true, nil
}

func scheduleTimes() []string {
	var times []string

	times = append(times, []string{
		"Prime Time Best",
		"Prime Time Morning",
		"Prime Time Afternoon",
		"Prime Time Evening",
		"Prime Time Swing",
	}...)

	for _, tz := range []string{"User", "PT", "GMT"} {
		for hour := 0; hour < 24; hour++ {
			tIndicator := "am"
			if hour > 11 {
				tIndicator = "pm"
			}
			t, _ := time.Parse("3:04 pm", fmt.Sprintf("%d:00 %s", hour%12, tIndicator))
			times = append(times, fmt.Sprintf("%s %s", tz, t.Format("03pm")))
		}
	}

	return times
}

func getScheduleTimeZone(schedule_time_index int) string {
	scheduleTimeList := scheduleTimes()
	if schedule_time_index < 0 || schedule_time_index >= len(scheduleTimeList) {
		return ""
	}
	scheduleTime := scheduleTimeList[schedule_time_index]
	return strings.Split(scheduleTime, " ")[0]
}

func getScheduleTimeOfDay(schedule_time_index int) string {
	scheduleTimeList := scheduleTimes()
	if schedule_time_index < 0 || schedule_time_index >= len(scheduleTimeList) {
		return ""
	}
	scheduleTime := scheduleTimeList[schedule_time_index]
	scheduleTimeArray := strings.Split(scheduleTime, " ")
	return scheduleTimeArray[len(scheduleTimeArray)-1]
}

func getPrimeTime(schedule_time_index int) string {
	if schedule_time_index < 0 {
		return ""
	}
	return strings.ToLower(getScheduleTimeOfDay(schedule_time_index))
}

func getUserTime(schedule_time_index int) string {
	if schedule_time_index < 0 {
		return ""
	}
	return getScheduleTimeZone(schedule_time_index)
}

func userLocalHour(hr int, userTimeZoneOffset int) int {
	localHr := hr + userTimeZoneOffset
	if localHr < 0 {
		return localHr + 24
	} else {
		return localHr
	}
}

func userCurrentHour(userTimeZoneOffset int) int {
	currentTime := time.Now().UTC().Add(time.Hour * (time.Duration(userTimeZoneOffset)))
	return currentTime.Hour()
}

func userCurrentTime(userTimeZoneOffset int) string {
	currentTime := time.Now().UTC().Add(time.Hour * (time.Duration(userTimeZoneOffset)))
	return currentTime.Format("03pm")
}

func handler(ctx context.Context, sqsEvent *events.SQSEvent) (string, error) {

	// start logging
	logger.
		NewEntry(log_msg_begin, LogMsgTxt).
		Debug()
	// Log environment configuration parameters
	logger.
		NewEntry(log_msg_envars, LogMsgTxt).
		WithStruct(env).
		Debug()
	// Log lambda parameters
	logger.
		NewEntry(log_msg_request, LogMsgTxt).
		WithStruct(sqsEvent).
		Debug()

	opCounters = &OperationalCounters{RWMutex: sync.RWMutex{}}

	// Attempt to load the parameters file written by the lambda extension.
	if err := parameterStore.ReadConfig(); err != nil {
		logger.
			NewEntry(log_msg_parameters_store_parse_error, LogMsgTxt).
			Error(err.Error())
		os.Exit(1)
	}

	parameters, err := parameterStore.ParameterToBytes(DEFAULT_HONEYBADGER_KEY)
	if err != nil {
		logger.
			NewEntry(log_msg_parameters_store, LogMsgTxt).
			Error(err.Error())
		os.Exit(1)
	}

	badgerApiKey := HoneyBadgerType{}

	_ = json.Unmarshal(parameters, &badgerApiKey)

	badger.Configure(honeybadger.Configuration{APIKey: badgerApiKey.Value, Env: env.Environment})

	if err := secretsStore.ReadConfig(); err != nil {
		logger.
			NewEntry(log_msg_secrets_store_parse_error, LogMsgTxt).
			Error(err.Error())
		os.Exit(1)
	}

	secrets, err := secretsStore.ParameterToBytes(env.PostgresKey)
	if err != nil {
		logger.
			NewEntry(log_msg_parameters_store, LogMsgTxt).
			Error(err.Error())
		os.Exit(1)
	}

	dbCreds := DBCredentialsType{}
	_ = json.Unmarshal(secrets, &dbCreds)

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
			// _, _ = badger.Notify(err.Error(),
			// 	honeybadger.Tags{
			// 		env.ApplicationName,
			// 		env.Environment,
			// 		"DATABASE"})
			// // TODO: Add context
			// badger.Flush()
			os.Exit(1)
		}
		logger.Debug("DATABASE", "connected")
	}

	if err := pgDb.Ping(); err != nil {
		logger.
			WithFields(Fields{"DATABASE": err.Error()}).
			Error(LogMsgTxt[log_msg_pg_database_connect])
		// _, _ = badger.Notify(err.Error(),
		// 	honeybadger.Tags{
		// 		env.ApplicationName,
		// 		"Database",
		// 		"PostgreSQL"})
		// badger.Flush()
		os.Exit(1)
	}

	var itemIdentifiers []map[string]string // this array will store all sqs failed message id

	for _, message := range sqsEvent.Records {
		opCounters.incSQSMessageBatch(1)
		sqsPayload := SQSPayload{}
		err := json.Unmarshal([]byte(message.Body), &sqsPayload)

		logger.WithStruct(sqsPayload).Debug()

		if err != nil {
			logger.
				NewEntry(log_msg_json_unmarshal_error, LogMsgTxt).
				WithFields(Fields{"SQSEVENT": err.Error()}).
				Error()
			itemIdentifiers = append(itemIdentifiers, map[string]string{"itemIdentifier": message.MessageId})
			continue
		}

		if err := pgDb.UpdateWorkLog(sqsPayload.Data.LogChunkId, "processing"); err != nil {
			opCounters.incLogChunkStatusUpdateError(1)
			logger.
				NewEntry(log_msg_execution_log_chunk_status_update_error, LogMsgTxt).
				WithFields(Fields{"DATABASE": err.Error()}).
				Error()
			itemIdentifiers = append(itemIdentifiers, map[string]string{"itemIdentifier": message.MessageId})
			continue
		}

		result, err := svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(env.S3SourceBucket),
			Key:    aws.String(sqsPayload.Data.S3File)})
		if err != nil {
			logger. // File not found
				NewEntry(log_msg_file_read_error, LogMsgTxt).
				Error(err.Error(), sqsPayload.Data.S3File)
			// _, _ = badger.Notify(err.Error(),
			// 	honeybadger.Context{
			// 		"ERROR_NUMBER": log_msg_file_read_error,
			// 		"ERROR_TEXT":   LogMsgTxt[log_msg_file_read_error],
			// 		"ERROR_URL":    DEFAULT_ERROR_URL,
			// 	},
			// 	honeybadger.Tags{
			// 		env.ApplicationName,
			// 		env.Environment,
			// 		appId,
			// 		LogMsgTxt[log_msg_file_read_error]})
			// badger.Flush()
			itemIdentifiers = append(itemIdentifiers, map[string]string{"itemIdentifier": message.MessageId})
			continue
		}

		defer result.Body.Close()

		// read from campaign table
		campaignTable, err := pgDb.ReadCampaignData(sqsPayload.Data.CampaignId)
		logger.WithStruct(campaignTable).Debug()
		if err != nil {
			logger. // player schedule failed
				NewEntry(log_msg_campaign_select_query_error, LogMsgTxt).
				Error(err, sqsPayload.Data.CampaignId)
			itemIdentifiers = append(itemIdentifiers, map[string]string{"itemIdentifier": message.MessageId})
			continue
		}

		// check campaign is active
		if !campaignTable.Active {
			logger. // player schedule failed
				NewEntry(log_msg_campaign_active_error, LogMsgTxt).
				Error(fmt.Errorf("Campaign not active"), sqsPayload.Data.CampaignId)
			itemIdentifiers = append(itemIdentifiers, map[string]string{"itemIdentifier": message.MessageId})
			continue
		}

		// read from message_template table
		messageTemplateTable, err := pgDb.ReadMessageTemplateData(int(campaignTable.MessageTemplateId.Int32))
		logger.WithStruct(messageTemplateTable).Debug()
		if err != nil {
			logger. // player schedule failed
				NewEntry(log_msg_message_template_select_query_error, LogMsgTxt).
				Error(err, campaignTable.MessageTemplateId)
			itemIdentifiers = append(itemIdentifiers, map[string]string{"itemIdentifier": message.MessageId})
			continue
		}

		// check campaign message_template is active
		if !messageTemplateTable.Active {
			logger. // player schedule failed
				NewEntry(log_msg_campaign_message_template_active_error, LogMsgTxt).
				Error(fmt.Errorf("Campaign Message Template not active"), campaignTable.MessageTemplateId)
			itemIdentifiers = append(itemIdentifiers, map[string]string{"itemIdentifier": message.MessageId})
			continue
		}

		logger.Debug(campaignTable)
		logger.Debug(messageTemplateTable)

		if !timeToDeliver(campaignTable) {
			logger. // player schedule failed
				NewEntry(log_msg_campaign_not_time_to_deliver_error, LogMsgTxt).
				Error(fmt.Errorf("Not time to deliver Campaign"), campaignTable.MessageTemplateId)
			itemIdentifiers = append(itemIdentifiers, map[string]string{"itemIdentifier": message.MessageId})
			continue
		}

		fileScanner := bufio.NewScanner(result.Body)
		for fileScanner.Scan() {
			player_id := fileScanner.Text()
			if player_id == "" {
				continue
			}
			opCounters.incTotalPlayer(1)
			schedule, err := schedulePlayerId(campaignTable, sqsPayload, player_id)
			if err != nil {
				logger. // invalid player id
					NewEntry(log_msg_player_id_invalid_error, LogMsgTxt).
					Error(err, player_id)
				opCounters.incInvalidPlayerIds(1)
				continue
			}
			if !schedule {
				// player schedule failed
				opCounters.incScheduledFailedPlayerIds(1)
				logger.Debug("Campaign is scheduled failed for player_id - " + string(player_id))
				continue
			}
			opCounters.incScheduledSuccessPlayerIds(1)
			logger.Debug("Campaign is scheduled success for player_id - " + string(player_id))
		}

		if fileScanner.Err() != nil {
			logger. // file scan error
				NewEntry(log_msg_file_read_error, LogMsgTxt).
				Error(fileScanner.Err(), sqsPayload.Data.S3File)
			itemIdentifiers = append(itemIdentifiers, map[string]string{"itemIdentifier": message.MessageId})
			continue
		}

	}

	logger.NewEntry(log_msg_end, LogMsgTxt).WithStruct(opCounters).Info()

	return "", nil
}

func init() {
	env = NewEnv()

	logger = NewLoggers(os.Stderr, env.LogLevel, 0).
		SetFields(Fields{
			EnvVars[ENV_APPLICATION_NAME]: env.ApplicationName,
		})
	parameterStore = &config.KeyClient{Filename: env.ParametersStore}
	secretsStore = &config.KeyClient{Filename: env.SecretsStore}
	badger = &im_badger.BadgerClient{}

	session, _ := session.NewSession()
	svc = s3Client{s3.New(session)}
}

func main() {
	defer badger.Monitor()
	lambda.Start(handler)
}
