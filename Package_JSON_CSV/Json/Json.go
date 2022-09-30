package MyJson

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"log"
)

type course struct {
	Name     string 
	Price    int
	Platform string
	Password string   
	Tags     []string 
}

// Encoding
func EncodeJson(network bytes.Buffer) {

	// Decode (receive) the value.
	// var network bytes.Buffer
	dec := gob.NewDecoder(&network)
	var list_of_subjects []course
	err := dec.Decode(&list_of_subjects)
	if err != nil {
		log.Fatal("decode error:", err)
	}
	
	// package this data as JSON data
	finaljson, err := json.MarshalIndent(list_of_subjects, "", "\t")
	if err != nil {
		panic(err)
	}
	_ = ioutil.WriteFile("test.json", finaljson, 0777)
}