package main

import (
    "fmt"
    "time"
)

const timeFormat = "1am"

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
            t, _ := time.Parse("15:04", fmt.Sprintf("%d:00", hour))
            times = append(times, fmt.Sprintf("%s %s", tz, t.Format(timeFormat)))
        }
    }

    return times
}

func main() {
    fmt.Println(scheduleTimes())
}
