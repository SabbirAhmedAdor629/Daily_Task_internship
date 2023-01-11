package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)
func main(value str[]) {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	
	err := enc.Encode(value)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	fmt.Println(network.Bytes())	
}