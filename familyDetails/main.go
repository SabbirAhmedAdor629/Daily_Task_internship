package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type MySelf struct {
	Name         string
	Parents_Info []Parents_Info
}

type Parents_Info struct {
	Type       string
	Name       string
	Occupation string
	Fav_colour []string
}

func main(){
	Father := Parents_Info{
		Type:       "Father Information",
		Name:       "Ahmed",
		Occupation: "businessman",
		Fav_colour: []string{"red", "green", "blue"},
	}
	Mother := Parents_Info{
		Type:       "Mothers Information",
		Name:       "Begum",
		Occupation: "Teacher",
		Fav_colour:  []string{"black", "green", "blue"},
	}

	Parents_Information := []Parents_Info{Father, Mother}
	self := MySelf{"Sabbir", Parents_Information}

	fmt.Printf("Details is : %v\\n", self)

	file, _ := json.MarshalIndent(self, "", " ")
	_ = ioutil.WriteFile("test.json", file, 0777)
}
