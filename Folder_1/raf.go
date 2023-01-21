

vars := make(map[string]interface{})

vars["delay_next_message_units"] = messageTemplate.DelayNextMessageUnits.String 


func timeFormatWithUnit(expires_in int32, unit string) time.Time {
	currentTime := time.Now().UTC()
	switch strings.ToLower(unit) {
	case "minutes":
		return currentTime.Add(time.Minute * time.Duration(expires_in))
	case "hours":
		return currentTime.Add(time.Hour * time.Duration(expires_in))
	case "days":
		return currentTime.Add(time.Hour * time.Duration(expires_in*24))
	default:
		return currentTime
	}
}














// type CampaignTable struct {
// 	Id                string
// 	Active            bool
// 	ScheduleTime      sql.NullInt32
// 	EventName         string
// 	MessageTemplateId sql.NullInt32
// 	ScheduleDays      ScheduleDays
// }


// campaignTableData *CampaignTable

// scheduleTime := int(campaignTableData.ScheduleTime.Int32)



// primeTime := getPrimeTime(scheduleTime)



// func getPrimeTime(schedule_time_index int) string {
// 	if schedule_time_index < 0 {
// 		return ""
// 	}
// 	return strings.ToLower(getScheduleTimeOfDay(schedule_time_index))				// morning afternoon and evening
// }

// func getScheduleTimeOfDay(schedule_time_index int) string {
// 	scheduleTimeList := scheduleTimes()

// 	if schedule_time_index < 0 || schedule_time_index >= len(scheduleTimeList) {
// 		return ""
// 	}

// 	scheduleTime := scheduleTimeList[schedule_time_index] // schedule time string from the schedule time list
// 	scheduleTimeArray := strings.Split(scheduleTime, " ") // splits this schedule time string by the space character

// 	return scheduleTimeArray[len(scheduleTimeArray)-1] //last element of the resulting array, the time of day morning/evening etc
// }
// func scheduleTimes() []string {

// 	var times []string

// 	times = append(times, []string{
// 		"Prime Time Best",
// 		"Prime Time Morning",
// 		"Prime Time Afternoon",
// 		"Prime Time Evening",
// 		"Prime Time Swing",
// 	}...)

// 	for _, tz := range []string{"User", "PT", "GMT"} {
// 		for hour := 0; hour < 24; hour++ {
// 			tIndicator := "am"
// 			if hour > 11 {
// 				tIndicator = "pm"
// 			}
// 			t, _ := time.Parse("3:04 pm", fmt.Sprintf("%d:00 %s", hour%12, tIndicator))
// 			times = append(times, fmt.Sprintf("%s %s", tz, t.Format("03pm")))
// 		}
// 	}

// 	return times
// }

















// // type CampaignTable struct {
// // 	Id                string
// // 	Active            bool
// // 	ScheduleTime      sql.NullInt32
// // 	EventName         string
// // 	MessageTemplateId sql.NullInt32
// // 	ScheduleDays      ScheduleDays
// // }


// // if len(campaignTableData.ScheduleDays) > 0 {
// // 	primeDay := time.Now().UTC().Add(time.Hour * (time.Duration(*playerHash.TimeZoneOffset))).Weekday()
// // 	fmt.Println(primeDay)
// // 	if !contains(campaignTableData.ScheduleDays, primeDay.String()) {
// // 		//fmt.Println("hey")
// // 		return false, nil
// // 	}
// // }

// // func contains(elements []interface{}, v string) bool {
// // 	for _, s := range elements {
// // 		if v == s { // comparing days
// // 			return true
// // 		}
// // 	}
// // 	return false
// // }