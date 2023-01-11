package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"math/rand"
    "time"
)

type course struct {
	Name     string `json:"coursename"`
	Price    int
	Platform string
	Password string   `json:"-"`
	Tags     []string `json:"tags,omitempty"`
}
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomString(n int) string {
	rand.Seed(time.Now().UnixNano())
    sb := strings.Builder{}
    sb.Grow(n)
    for i := 0; i < n; i++ {
        sb.WriteByte(charset[rand.Intn(len(charset))])
    }
    return sb.String()
}

func main() {
	fmt.Println("Bit")
	EncodeJson()
	//DecodeJson()

}

// Encoding
func EncodeJson() {
	list_of_courses := []course{
		{randomString(20), 299, "LearnCodeOnline.in", "abc123", []string{"web-dev", "js"}},
		{randomString(20), 199, "LearnCodeOnline.in", "bcd123", []string{"full_stack", "js"}},
		{randomString(20), 299, "LearnCodeOnline.in", "hit123", nil},
	}

	// package this data as JSON data
	finaljson, err := json.MarshalIndent(list_of_courses, "", "\t")
	if err != nil {
		panic(err)
	}
	_ = ioutil.WriteFile("test.json", finaljson, 0777)
	fmt.Printf("%s\n",finaljson)

}



//	Decoding
func DecodeJson() {
	jsonDataFromWeb := []byte(`{
		"Name": "MERN bootcamp",
		"Price": 199,
		"Platform": "LearnCodeOnline.in",
		"tags": ["full_stack","js"]
	}`)

	var lcocourse course

	checkValid := json.Valid(jsonDataFromWeb)
	if checkValid {
		fmt.Println("json was valid")
		json.Unmarshal(jsonDataFromWeb, &lcocourse)
		fmt.Printf("%#v\n", lcocourse)
	} else {
		fmt.Println("Json was not valid")
	}
}
