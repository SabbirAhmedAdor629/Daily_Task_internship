//	(Date.today - Date.parse("2022-09-26")).to_i
//  => 87

package main

import (
	"fmt"
	"time"
)

func main() {
	// Parse the input date string
	t, err := time.Parse(time.RFC3339, "2021-02-18")
	if err != nil {
		// Handle the error
		fmt.Println(err)
		return
	}

	today := time.Now()
	diff := today.Sub(t)
	fmt.Println(int(diff.Hours() / 24))

}
