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
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/honeybadger-io/honeybadger-go"
	im_badger "influencemobile.com/badger"
	im_targeting "influencemobile.com/libs/im_targeting"
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
	ctx                   context.Context
	dynamoSvc             dynamodbiface.DynamoDBAPI
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

type TargetingPlayerHash struct {
	PlayerId       int               `json:"player_id"`
	TimeZoneOffset *int              `json:"time_zone_offset"`
	PrimeTime      map[string]string `json:"prime_time"`
}

// very simple function 
func contains(elements []interface{}, v string) bool {
	for _, s := range elements {
		if v == s {
			return true
		}
	}
	return false
}

// getPrimeTime , Locad location are static
func timeToDeliver(campaignData *CampaignTable) bool {
	scheduleTime := int(campaignData.ScheduleTime.Int32)

	timeZone := getScheduleTimeZone(scheduleTime)

	if timeZone == "Prime" {
		return true
	} else if timeZone == "User" {
		return userPresentsInScheduleTime(scheduleTime)
	}getPrimeTime
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


// not clear
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


// locations are not dynamic no need to test
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

// incInvalidPlayerIds is missing
func getValidPlayerId(player_id string) (*int, error) {
	if player_id == "" {
		opCounters.incInvalidPlayerIds(1)
		return nil, fmt.Errorf("player_id can not be null")
	}

	playerId, err := strconv.Atoi(player_id)
	if err != nil {
		opCounters.incInvalidPlayerIds(1)
		return nil, fmt.Errorf("player_id must be of type int")
	}
	return &playerId, nil
}

func getPlayerHashData() ([]byte, error){
	hash := map[string]string {
		"time_zone_offset" : "-6",
	}
	jsonPlayerHash, err:= json.Marshal(playerHash)
	return jsonPlayerHash,err
}

func getPlayerHashValue(playerId int) *TargetingPlayerHash {
	var err error = nil
	hash, err := getPlayerHashData()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	playerHashMap := TargetingPlayerHash{}
	_ = json.Unmarshal([]byte(hash), &playerHashMap)
	logger.Debug(playerHashMap)
	return &playerHashMap
}

// incInvalidPlayerIds is missing in getValidPlayerId function
func userTimeToDeliver(player_id string, scheduleDays ScheduleDays) (bool, error) {
	if len(scheduleDays) == 0 {
		return true, nil
	}

	playerId, err := getValidPlayerId(player_id)
	if err != nil {
		return false, err
	}

	playerHash := getPlayerHashValue(*playerId) // placeholder :  Im::Targeting::Target.player(player_id)

	if playerHash == nil || playerHash.TimeZoneOffset == nil {
		return false, nil
	}

	userDay := time.Now().UTC().Add(time.Hour * (time.Duration(*playerHash.TimeZoneOffset))).Weekday()
	logger.Debug("User Current Day [3] : " + userDay.String())
	logger.Debug("User Schedule Day [3] : ", scheduleDays)
	if contains(scheduleDays, userDay.String()) {
		return true, nil
	}
	return false, nil
}

// incInvalidPlayerIds is missing in getValidPlayerId function
func primeTimeToDeliver(player_id string, campaignTableData *CampaignTable) (bool, error) {

	scheduleTime := int(campaignTableData.ScheduleTime.Int32)

	if len(campaignTableData.ScheduleDays) == 0 {
		return true, nil
	}

	playerId, err := getValidPlayerId(player_id)
	if err != nil {
		return false, err
	}

	playerHash := getPlayerHashValue(*playerId) // player_hash = Im::Targeting::Target.player(player_id)

	if playerHash == nil || playerHash.TimeZoneOffset == nil || playerHash.PrimeTime == nil {
		return false, nil
	}

	if len(campaignTableData.ScheduleDays) > 0 {
		primeDay := time.Now().UTC().Add(time.Hour * (time.Duration(*playerHash.TimeZoneOffset))).Weekday()
		logger.Debug("Prime Day current : " + primeDay.String())
		logger.Debug("Campaign Schedule Days : ", campaignTableData.ScheduleDays)
		if !contains(campaignTableData.ScheduleDays, primeDay.String()) {
			return false, nil
		}
	}

	primeTime := getPrimeTime(scheduleTime)
	logger.Debug("primeTime: ", primeTime)
	primeValue, ok := playerHash.PrimeTime[primeTime]
	if !ok || primeValue == "" {
		return false, fmt.Errorf(fmt.Sprintf("player hash does not contain prime_time: %s", primeTime))
	}

	primeTimeValue, err := strconv.Atoi(primeValue)
	if err != nil {
		return false, fmt.Errorf("prime time must be of type int")
	}

	localHour := userLocalHour(primeTimeValue, *playerHash.TimeZoneOffset)
	currentHour := userCurrentHour(*playerHash.TimeZoneOffset)
	logger.Debug("User Local Prime Hour: ", localHour)
	logger.Debug("User Current Hour: ", currentHour)
	if currentHour == localHour {
		return true, nil
	}

	return false, nil
}


// AVAILABLE KEY IN MAP (TEST 9)
func availableKeyInMap(mapData map[string]interface{}, keyList ...string) bool {
	for _, key := range keyList {
		if val, ok := mapData[key]; !ok || val == nil {
			return false
		}
	}
	return true
}

// CampaignTable and MessageTemplate struct is missing 
func schedulePlayerId(campaignTable *CampaignTable, messageTemplate *MessageTemplate, sqsData SQSPayload, player_id string) (bool, error) {
	scheduleTime := int(campaignTable.ScheduleTime.Int32)
	scheduleTimeZone := getScheduleTimeZone(scheduleTime)
	fmt.Println("Campaign Schedule Time Zone [2] : ", scheduleTimeZone)
	if scheduleTimeZone == "User" {
		inUserTime, err := userTimeToDeliver(player_id, campaignTable.ScheduleDays)
		if err != nil {
			return false, err
		}
		if !inUserTime {
			opCounters.incOutOfUserTime(1)	// increaments the counter
			return false, nil
		}
	} else if scheduleTimeZone == "Prime" {
		inPrimeTime, err := primeTimeToDeliver(player_id, campaignTable)
		if err != nil {
			return false, err
		}
		if !inPrimeTime {
			opCounters.incOutOfPrimeTime(1)
			return false, nil
		}
		logger.Debug("inPrimeTime : " + fmt.Sprint(inPrimeTime))
	}
	return true, nil
}


// openRedisStore function missing
func canDeliverTo(playerId string, messageTemplate *MessageTemplate) (bool, error) {
	redisCLient, err := openRedisStore(env.RedisHost, env.RedisPort, "", 0)			
	if err != nil {
		logger.Debug("redis connection failed")
		return true, nil // it will be return false, err
	}
	// set sample cap data into redis
	redisCLient.setSampleCapsValue(messageTemplate.Id, playerId)
	dailyCap := redisCLient.getDailyCap(messageTemplate.Id, playerId)
	logger.Debug("Redis Daily Cap: ", dailyCap)
	if messageTemplate == nil {
		return true, nil
	}
	return true, nil
}


// getCapKeys function missing in the file
func isCapped(playerId string, messageTemplate *MessageTemplate) (bool, error) {
	capKeys := getCapKeys(messageTemplate.Id, playerId)		
	fmt.Print(capKeys)
	if messageTemplate.DailyCap.Valid {
		return true, nil
	}
	if messageTemplate.MonthlyCap.Valid {
		return true, nil
	}
	if messageTemplate.TotalCap.Valid {
		return true, nil
	}
	// throttle
	return false, nil
}


// Getting the timezone (TESTING 1)
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

// returns user, pt or gmt of the specific index (TESTING 2)
func getScheduleTimeZone(schedule_time_index int) string {
	scheduleTimeList := scheduleTimes()
	if schedule_time_index < 0 || schedule_time_index >= len(scheduleTimeList) {
		return ""
	}
	scheduleTime := scheduleTimeList[schedule_time_index]
	return strings.Split(scheduleTime, " ")[0]
}

// returns time of day for valid index (last half) (TESTING 3)
func getScheduleTimeOfDay(schedule_time_index int) string {
	scheduleTimeList := scheduleTimes()
	if schedule_time_index < 0 || schedule_time_index >= len(scheduleTimeList) {
		return ""
	}
	scheduleTime := scheduleTimeList[schedule_time_index]
	scheduleTimeArray := strings.Split(scheduleTime, " ")
	return scheduleTimeArray[len(scheduleTimeArray)-1]
}


// returns time of day for valid index (last half) (TESTING 4)
func getPrimeTime(schedule_time_index int) string {
	if schedule_time_index < 0 {
		return ""
	}
	return strings.ToLower(getScheduleTimeOfDay(schedule_time_index))
}

// return user, pt and gmt for valid index (TESTING 5)
func getUserTime(schedule_time_index int) string {
	if schedule_time_index < 0 {
		return ""
	}
	return getScheduleTimeZone(schedule_time_index)
}

// returns the exact local hour	(TESTING 6)
func userLocalHour(hr int, userTimeZoneOffset int) int {
	localHr := hr + userTimeZoneOffset
	if localHr < 0 {
		return localHr + 24				// if the localHr+24 > 12 (then return localHr+24-12)
	} else {
		return localHr
	}
}

// this function will return hour in the specific time in the timezone   (TESTING 7)
func userCurrentHour(userTimeZoneOffset int) int {
	currentTime := time.Now().UTC().Add(time.Hour * (time.Duration(userTimeZoneOffset)))
	return currentTime.Hour()
}

// user current time without any minitue or seconds (TESTING 8)
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
			logger.
				NewEntry(log_msg_campaign_select_query_error, LogMsgTxt).
				Error(err, sqsPayload.Data.CampaignId)
			itemIdentifiers = append(itemIdentifiers, map[string]string{"itemIdentifier": message.MessageId})
			continue
		}

		// read from message_template table
		messageTemplateTable, err := pgDb.ReadMessageTemplateData(int(campaignTable.MessageTemplateId.Int32))
		logger.WithStruct(messageTemplateTable).Debug()
		if err != nil {
			logger.
				NewEntry(log_msg_message_template_select_query_error, LogMsgTxt).
				Error(err, campaignTable.MessageTemplateId.Int32)
			itemIdentifiers = append(itemIdentifiers, map[string]string{"itemIdentifier": message.MessageId})
			continue
		}

		fileScanner := bufio.NewScanner(result.Body)
		for fileScanner.Scan() {
			player_id := fileScanner.Text()
			if player_id == "" {
				continue
			}
			opCounters.incInitialTarget(1)
			schedule, err := schedulePlayerId(campaignTable, messageTemplateTable, sqsPayload, player_id)
			if err != nil {
				logger. // campaign schedule error
					NewEntry(log_msg_campaign_schedule_error, LogMsgTxt).
					Error(err.Error())
				opCounters.incScheduledFailedPlayerIds(1)
				continue
			}
			if !schedule {
				// player schedule failed
				opCounters.incScheduledFailedPlayerIds(1)
				logger.Debug("[ FAILED ] Campaign is scheduled failed for player_id - " + string(player_id))
				continue
			}
			opCounters.incScheduledSuccessPlayerIds(1)
			logger.Debug("[ SUCCESS ] Campaign is scheduled success for player_id - " + string(player_id))
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

	dynamoSvc = dynamodb.New(session.New())
	session, _ := session.NewSession()
	svc = s3Client{s3.New(session)}
	ctx = context.Background()
}

func main() {
	defer badger.Monitor()
	lambda.Start(handler)
}
