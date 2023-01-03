package main

import (
    "time"
    "fmt"
)

func main(){
	dateString := "2021-02-18T21:54:42.123Z"
	date, err := time.Parse(time.RFC3339, dateString)
	if (err != nil){
		fmt.Println(err)
	}
	fmt.Println(date)
	
    // fmt.Println(time.Now().Format(time.RFC3339))

	//time.Parse(time.RFC3339, dateString.(string))

}

