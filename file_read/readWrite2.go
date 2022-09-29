package main

import (
	//"encoding/json"
	"encoding/json"
	"fmt"
	"io/ioutil"
	//"io/ioutil"
	//"os"
	// "strconv"
)

type Mother struct {
	M_name       string
	M_fav_colour []string
}

type Father struct {
	F_name       string
	F_fav_colour []string
}

type Details struct {
	Name  string
	Id    int32
	EMAIL string
	Father
	Mother
}

func main() {
	detail1 := Details{}
	detail1.Name = "Sabbir"
	detail1.Id = 31112
	detail1.EMAIL = "Sabbir11@gmail.com"
	detail1.F_name = "Ahmed"
	detail1.Father.F_fav_colour = []string{"red", "green", "black"}
	detail1.M_name = "Begum"
	detail1.M_fav_colour = []string{"red", "green", "black"}

	fmt.Println(detail1)


	// Writing data in json
	file, _ := json.MarshalIndent(detail1, "", " ")
	_ = ioutil.WriteFile("test2.json", file, 0777)
	

}
