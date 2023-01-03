package main

import (
	"fmt"
	"time"
)

func main() {
	// Parse the input date string in the ISO 8601 format
	date, err := time.Parse("2006-01-02", "2022-12-21")
	if err != nil {
		// Handle error
	}

	// Format the date using the desired layout
	formattedDate := date.Format("Mon, 02 Jan 2006")

	// Print the result
	fmt.Println(formattedDate)
}
