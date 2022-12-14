package main

import (
	//"encoding/json"
	"encoding/json"
	"fmt"
	"io/ioutil"
	//"io/ioutil"
)

type People struct {
	Myself []mySelf
}

type mySelf struct {
	Name         string
	Parents_Info []parentsInfo
}


type parentsInfo struct {
	Type       string
	Name       string
	Fav_colour []string
	Occupation string
}


func main() {
	Father := parentsInfo{
		Type:       "Father Information",
		Name:       "Ahmed",
		Occupation: "businessman",
		Fav_colour: []string{"red", "green", "blue"},
	}
	Mother := parentsInfo{
		Type:       "Mothers Information",
		Name:       "Begum",
		Occupation: "Teacher",
	}
	
	
	
	Parents_Info := []parentsInfo{Father, Mother}
	Self := mySelf{"Sabbir", Parents_Info}
	People := People{[]mySelf{Self}}

	fmt.Println("Peoples : %v\\n", People)

	file, _ := json.MarshalIndent(People, "", " ")
	_ = ioutil.WriteFile("test.json", file, 0777)

}
