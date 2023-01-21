package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	//"sync"
	"database/sql"
	"time"
)

// var (
// 	opCounters *OperationalCounters
// )

// TEST 1
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

// TEST 2
func getScheduleTimeZone(schedule_time_index int) string {
	scheduleTimeList := scheduleTimes()
	if schedule_time_index < 0 || schedule_time_index >= len(scheduleTimeList) {
		return ""
	}
	scheduleTime := scheduleTimeList[schedule_time_index]
	return strings.Split(scheduleTime, " ")[0]
}

// TEST 3
func getScheduleTimeOfDay(schedule_time_index int) string {
	scheduleTimeList := scheduleTimes()

	if schedule_time_index < 0 || schedule_time_index >= len(scheduleTimeList) {
		return ""
	}

	scheduleTime := scheduleTimeList[schedule_time_index] // schedule time string from the schedule time list
	scheduleTimeArray := strings.Split(scheduleTime, " ") // splits this schedule time string by the space character

	return scheduleTimeArray[len(scheduleTimeArray)-1] //last element of the resulting array, the time of day.
}

// TEST 4
func getPrimeTime(schedule_time_index int) string {
	if schedule_time_index < 0 {
		return ""
	}
	return strings.ToLower(getScheduleTimeOfDay(schedule_time_index))
}

// TEST 5
func getUserTime(schedule_time_index int) string {
	if schedule_time_index < 0 {
		return ""
	}
	return getScheduleTimeZone(schedule_time_index)
}

// TEST 8
func userCurrentTime(userTimeZoneOffset int) string {
	currentTime := time.Now().UTC().Add(time.Hour * (time.Duration(userTimeZoneOffset)))
	return currentTime.Format("03pm")
}

// TEST 9
func AvailableKeyInMap(mapData map[string]interface{}, keyList ...string) bool {
	for _, key := range keyList {
		if val, ok := mapData[key]; !ok || val == nil {
			return false
		}
	}
	return true
}

// TEST 10

//	func main() {
//		opCounters = &OperationalCounters{RWMutex: sync.RWMutex{}}
//		playerID := "343"
//		scheduleDays := []string{"Tuesday"}
//		y := make([]interface{}, len(scheduleDays))
//		for i, v := range scheduleDays {
//			y[i] = v
//		}
//		isInUserTime, err := UserTimeToDeliver(playerID, y)
//		if err != nil {
//			fmt.Println(err.Error())
//			return
//		}
//		if isInUserTime {
//			fmt.Println("within the user's schedule time.")
//		} else {
//			fmt.Println("not within the user's schedule time.")
//		}
//	}

type TargetingPlayerHash struct {
	PlayerId       int               `json:"player_id"`
	TimeZoneOffset *int              `json:"time_zone_offset"`
	PrimeTime      map[string]string `json:"prime_time"`
}

func getPlayerHashData() ([]byte, error) {
	hash := map[string]interface{}{
		//"time_zone_offset": -23,
	}
	jsonPlayerHash, err := json.Marshal(hash)
	return jsonPlayerHash, err
}

// func getPlayerHashValue(playerId int) *TargetingPlayerHash {
// 	hash, err := getPlayerHashData() // recieve a json data here
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil
// 	}
// 	//fmt.Println(hash)
// 	playerHashMap := TargetingPlayerHash{}
// 	_ = json.Unmarshal([]byte(hash), &playerHashMap) // decoding the json data, store into struct
// 	//fmt.Println(*playerHashMap.TimeZoneOffset)
// 	return &playerHashMap
// }

func UserTimeToDeliver(player_id string, scheduleDays ScheduleDays) (bool, error) {
	if len(scheduleDays) == 0 {
		return true, nil
	}

	playerId, err := getValidPlayerId(player_id) //validates the player id
	if err != nil {
		return false, err
	}

	playerHash := getPlayerHashValue(*playerId) // hash value of the targeted player id
	//fmt.Println(playerHash)
	if playerHash == nil || playerHash.TimeZoneOffset == nil {
		return false, nil
	}
	//fmt.Println(*playerHash.TimeZoneOffset)
	userDay := time.Now().UTC().Add(time.Hour * (time.Duration(*playerHash.TimeZoneOffset))).Weekday()
	//fmt.Println(userDay)
	if contains(scheduleDays, userDay.String()) {
		return true, nil
	}
	return false, nil
}

func getValidPlayerId(player_id string) (*int, error) {
	if player_id == "" {
		//opCounters.incInvalidPlayerIds(1)
		return nil, fmt.Errorf("player_id can not be null")
	}

	playerId, err := strconv.Atoi(player_id)
	if err != nil {
		//opCounters.incInvalidPlayerIds(1) // counter increase
		return nil, fmt.Errorf("player_id must be of type int")
	}
	return &playerId, nil
}

// HASH VALUE FOR PRIMETIME TABLE
func getPlayerHashValue(playerID int) *TargetingPlayerHash {
	// mock function to return a dummy PlayerHash struct
	tz := +6
	return &TargetingPlayerHash{
		TimeZoneOffset: &tz,
		PrimeTime: map[string]string{
			"morning":   "8",
			"afternoon": "16",
			"evening":   "24",
			"best":      "22",
			"swing":     "21",
		},
	}
}

// TEST 11

func main() {
	campaign_1 := CampaignTable{
		Id:                "1",
		Active:            true,
		ScheduleTime:      sql.NullInt32{Int32: 3, Valid: true},
		EventName:         "event1",
		MessageTemplateId: sql.NullInt32{Int32: 1, Valid: true},
		ScheduleDays:      ScheduleDays{"Tuesday", "Thursday", "Friday"},
	}
	value, err := PrimeTimeToDeliver("1", &campaign_1)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(value)
	}
}

func PrimeTimeToDeliver(player_id string, campaignTableData *CampaignTable) (bool, error) {
	scheduleTime := int(campaignTableData.ScheduleTime.Int32)
	if len(campaignTableData.ScheduleDays) == 0 {
		return true, nil
	}
	playerId, err := getValidPlayerId(player_id)
	if err != nil {
		fmt.Println("hey 1")
		return false, err
	}
	playerHash := getPlayerHashValue(*playerId) // player_hash = Im::Targeting::Target.player(player_id)
	if playerHash == nil || playerHash.TimeZoneOffset == nil || playerHash.PrimeTime == nil {
		fmt.Println("hey 2")
		return false, nil
	}
	if len(campaignTableData.ScheduleDays) > 0 {
		primeDay := time.Now().UTC().Add(time.Hour * (time.Duration(*playerHash.TimeZoneOffset))).Weekday()
		fmt.Println(primeDay)
		if !contains(campaignTableData.ScheduleDays, primeDay.String()) {
			fmt.Println("hey 3")
			return false, nil
		}
	}
	primeTime := getPrimeTime(scheduleTime) // morning(1) / afternoon(2) /evening(3)
	fmt.Println(primeTime)
	primeValue, ok := playerHash.PrimeTime[primeTime]
	if !ok || primeValue == "" {
		fmt.Print("122222222222222222222222")
		return false, fmt.Errorf(fmt.Sprintf("player hash does not contain prime_time: %s", primeTime))
	}
	primeTimeValue, err := strconv.Atoi(primeValue)
	fmt.Println(primeTimeValue)
	if err != nil {
		fmt.Println("hey 5")
		return false, fmt.Errorf("prime time must be of type int")
	}
	return false, nil
}

// TEST 6
func userLocalHour(hr int, userTimeZoneOffset int) int {
	localHr := hr + userTimeZoneOffset
	if localHr < 0 {
		return localHr + 24 // if the localHr+24 > 12 (then return localHr+24-12)
	} else {
		return localHr
	}
}

// TEST 7
func userCurrentHour(userTimeZoneOffset int) int {
	currentTime := time.Now().UTC().Add(time.Hour * (time.Duration(userTimeZoneOffset)))
	return currentTime.Hour()
}

func contains(elements []interface{}, v string) bool {
	for _, s := range elements {
		if v == s { // comparing days
			return true
		}
	}
	return false
}
