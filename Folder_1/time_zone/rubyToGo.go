package main

import (
	"fmt"
	"time"
)

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
			tFormate := "am"
			if hour > 11 {
				tFormate = "pm"
			}
			t, _ := time.Parse("3:04 pm", fmt.Sprintf("%d:00 %s", hour%12, tFormate))
			fmt.Println(t)
			times = append(times, fmt.Sprintf("%s %s", tz, t.Format("03pm")))
		}
	}

	return times
}

func main() {

	a := scheduleTimes()
	for _, v := range a {
		fmt.Println(v)
	}
	// fmt.Println(scheduleTimes())
}


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
// 			tFormate := "am"
// 			if hour > 11 {
// 				tFormate = "pm"
// 			}
// 			t, _ := time.Parse("3:04 pm", fmt.Sprintf("%d:00 %s", hour%12, tFormate))
// 			fmt.Println(t)
// 			times = append(times, fmt.Sprintf("%s %s", tz, t.Format("03pm")))
// 		}
// 	}

// 	return times
// }