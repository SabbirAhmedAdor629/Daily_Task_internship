package main

import (
	"fmt"
	"time"
)

func main() {
	// Parse the input date string
	input := "2021-02-18T21:54:42.123Z"
	t, err := time.Parse(time.RFC3339, input)
	if err != nil {
		// Handle the error
		fmt.Println(err)
		return
	}

	// Format the time using the desired layout
	output := t.Format("2006-01-02 15:04:05.000 0700")
	fmt.Println(output)
}
