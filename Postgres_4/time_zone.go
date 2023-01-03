package main

import (
	"fmt"
	"time"
)

func main() {
	// location1, err1 := time.LoadLocation("America/Adak")
	// location2, err2 := time.LoadLocation("America/New_York")

	// if err1 != nil || err2 != nil {
	//     // handle error
	// }

	// offset1 := time.Now().In(location1).UTC().Hour() / 3600
	// offset2 := time.Now().In(location2).UTC().Hour() / 3600

	// fmt.Println(offset2)
	// offset_ranges := []int{}
	// for i := offset1; i <= offset2; i++ {
	// 	offset_ranges = append(offset_ranges, i)
	// 	fmt.Println(i)
	// }

	adakLoc, _ := time.LoadLocation("America/Adak")
	_, adakOffset := time.Now().In(adakLoc).Zone()
	range1 := adakOffset / 3600
	loc, _ := time.LoadLocation("America/New_York") // use other time zones such as MST, IST
	_, offset := time.Now().In(loc).Zone()
	range2 := offset / 3600

	offset_ranges := []int{}
	for i := range1; i <= range2; i++ {
		offset_ranges = append(offset_ranges, i)
		//fmt.Println(i)
	}
	fmt.Println(offset_ranges)

}

// package main

// import (
// 	"time"
// 	"fmt"// adakLoc, _ := time.LoadLocation("America/Adak")
// _, adakOffset := time.Now().In(adakLoc).Zone()
// fmt.Println(adakOffset / 3600)
// loc, _ := time.LoadLocation("America/New_York") // use other time zones such as MST, IST
// _, offset := time.Now().In(loc).Zone()
// fmt.Println(offset / 3600)
// )

// func main() {

// adakLoc, _ := time.LoadLocation("America/Adak")
// _, adakOffset := time.Now().In(adakLoc).Zone()
// fmt.Println(adakOffset / 3600)
// loc, _ := time.LoadLocation("America/New_York") // use other time zones such as MST, IST
// _, offset := time.Now().In(loc).Zone()
// fmt.Println(offset / 3600)
// 	// Get the UTC offset in hours for the America/Adak time zone
// 	offsetAdak := time.Now().In(time.FixedZone("America/Adak", -36000)).UTC().Round(time.Second).Sub(time.Now().In(time.UTC)).Hours()

// 	// Get the UTC offset in hours for the America/New_York time zone
// 	offsetNewYork := time.Now().In(time.FixedZone("America/New_York", -18000)).UTC().Round(time.Second).Sub(time.Now().In(time.UTC)).Hours()

// 	// Create a range with the UTC offsets for America/Adak and America/New_York
// 	offsetRanges := (offsetAdak)..(offsetNewYork)
// 	fmt.Println(offsetRanges)
// }

// def users_present_in_schedule_time?

//     # as we only operate in US time zones, America/New_York to America/Adak for ‘User’ time campaigns
//     # offset_ranges = (-9..-4) when daylight saving is ON and (-10..-5) when it's OFF

//     offset_ranges = (((Time.now.in_time_zone('America/Adak').utc_offset/3600))..((Time.now.in_time_zone('America/New_York').utc_offset/3600)))

//    # offset_ranges_will_be_like = [-10,-9,-8,-7,-6,-5]

//     user_current_times = offset_ranges.map { |tz_offset| user_current_time(tz_offset) }
//     user_current_times.include?(schedule_time_of_day)

// end

// package main

// import (
// 	"time"
// )

// func usersPresentInScheduleTime() bool {
// 	// as we only operate in US time zones, America/New_York to America/Adak for ‘User’ time campaigns
// 	// offset_ranges = (-9..-4) when daylight saving is ON and (-10..-5) when it's OFF
// 	offsetRanges := []int{-10, -9, -8, -7, -6, -5}
// 	userCurrentTimes := make([]time.Time, len(offsetRanges))
// 	for i, tzOffset := range offsetRanges {
// 		userCurrentTimes[i] = userCurrentTime(tzOffset)
// 	}
// 	for _, t := range userCurrentTimes {
// 		if t == scheduleTimeOfDay {
// 			return true
// 		}
// 	}
// 	return false
// }

// func userCurrentTime(tzOffset int) time.Time {
// 	// Implement this function
// 	return time.Time{}
// }

// func scheduleTimeOfDay() time.Time {
// 	// Implement this function
// 	return time.Time{}
// }

// package main

// import (
//   "time"
// )

// func users_present_in_schedule_time(schedule_time_of_day time.Time) bool {
//   // as we only operate in US time zones, America/New_York to America/Adak for ‘User’ time campaigns
//   // offset_ranges = (-9..-4) when daylight saving is ON and (-10..-5) when it's OFF
//   offset_ranges := []int{}
//   for i := (time.Now().In(time.FixedZone("America/Adak", -9*3600)).UTC().Hour() / 3600); i <= (time.Now().In(time.FixedZone("America/New_York", -5*3600)).UTC().Hour() / 3600); i++ {
//     offset_ranges = append(offset_ranges, i)
//   }
//   // offset_ranges_will_be_like = [-10,-9,-8,-7,-6,-5]
//   user_current_times := []time.Time{}
//   for _, tz_offset := range offset_ranges {
//     user_current_times = append(user_current_times, user_current_time(tz_offset))
//   }
//   for _, user_current_time := range user_current_times {
//     if user_current_time == schedule_time_of_day {
//       return true
//     }
//   }
//   time.
//   return false
// }
