package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	opCounters *OperationalCounters
)

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
	scheduleTime := scheduleTimeList[schedule_time_index]
	scheduleTimeArray := strings.Split(scheduleTime, " ")
	return scheduleTimeArray[len(scheduleTimeArray)-1]
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

func main() {

	playerID := "12345"
	scheduleDays := []string{"Tuesday"}
	y := make([]interface{}, len(scheduleDays))
	for i, v := range scheduleDays {
		y[i] = v
	}
	isInUserTime, err := UserTimeToDeliver(playerID, y)
	if err != nil {
		fmt.Println(err)
		return
	}
	if isInUserTime {
		fmt.Println("within the user's schedule time.")
	} else {
		fmt.Println("not within the user's schedule time.")
	}
}

type TargetingPlayerHash struct {
	PlayerId       int               `json:"player_id"`
	TimeZoneOffset *int              `json:"time_zone_offset"`
	PrimeTime      map[string]string `json:"prime_time"`
}

func getPlayerHashData() ([]byte, error) {
	hash := map[string]string{
		"time_zone_offset": "-6",
	}
	jsonPlayerHash, err := json.Marshal(hash)
	return jsonPlayerHash, err
}

func getPlayerHashValue(playerId int) *TargetingPlayerHash {
	hash, err := getPlayerHashData() 	// recieve a json data here
	if err != nil {
		fmt.Println(err)
		return nil
	}
	playerHashMap := TargetingPlayerHash{}
	_ = json.Unmarshal([]byte(hash), &playerHashMap) 	// decoding the json data, store into struct
	return &playerHashMap
}

func getValidPlayerId(player_id string) (*int, error) {
	if player_id == "" {
		opCounters.incInvalidPlayerIds(1)				
		return nil, fmt.Errorf("player_id can not be null")
	}

	playerId, err := strconv.Atoi(player_id)
	if err != nil {
		opCounters.incInvalidPlayerIds(1)			// counter increase
		return nil, fmt.Errorf("player_id must be of type int")
	}
	return &playerId, nil
}

func contains(elements []interface{}, v string) bool {
	for _, s := range elements {
		if v == s {					// comparing days
			return true
		}
	}
	return false
}

func UserTimeToDeliver(player_id string, scheduleDays ScheduleDays) (bool, error) {
	if len(scheduleDays) == 0 {
		return true, nil
	}

	playerId, err := getValidPlayerId(player_id)  		//validates the player id
	if err != nil {
		return false, err
	}

	playerHash := getPlayerHashValue(*playerId) 		// placeholder :  Im::Targeting::Target.player(player_id)

	if playerHash == nil || playerHash.TimeZoneOffset == nil {
		return false, nil
	}

	userDay := time.Now().UTC().Add(time.Hour * (time.Duration(*playerHash.TimeZoneOffset))).Weekday()

	if contains(scheduleDays, userDay.String()) {
		return true, nil
	}
	return false, nil
}
