package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type MySelf struct {
	Name         string
	Parents_Info []parentsInfo
}

type parentsInfo struct {
	Type       string
	Name       string
	Occupation string
	Fav_colour []string
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
	Mother.Fav_colour = []string{"black", "green", "blue"}

	Parents_Info := []parentsInfo{Father, Mother}
	self := MySelf{"Sabbir", Parents_Info}

	fmt.Printf("Company is %v\\n", self)

	file, _ := json.MarshalIndent(self, "", " ")
	_ = ioutil.WriteFile("test.json", file, 0777)
}
