package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

type course struct {
	Name     string
	Price    int
	Platform string
	Password string
	Tags     []string
}

type subject struct {
	Name     string
	Price    *int32
	Platform string
	Password string
	Tags     []string
}

func main() {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	dec := gob.NewDecoder(&network)

	// Encode (send) the value.
	list_of_courses := []course{
		{"Reacjs Bootcamp", 299, "LearnCodeOnline.in", "abc123", []string{"web-dev", "js"}},
		{"MERN Bootcamp", 199, "LearnCodeOnline.in", "bcd123", []string{"full_stack", "js"}},
		{"Angular Bootcamp", 299, "LearnCodeOnline.in", "hit123", nil},
	}

	err := enc.Encode(list_of_courses)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	fmt.Println(network.Bytes())

	//arr1 := network.Bytes()

	// Decode (receive) the value.
	var list_of_subjects []course
	err = dec.Decode(&list_of_subjects)
	if err != nil {
		log.Fatal("decode error:", err)
	}

	fmt.Println(
		list_of_subjects,
	)


	//str1 := string(byteArray[:])
	//fmt.Println("String =",str1)
}
