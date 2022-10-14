package main

import (
	"encoding/json"
	"fmt"
	"os"
)


type Product struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

func main() {
	reader, _ := os.Open("dat.json")
	decoder := json.NewDecoder(reader)

	product := &Product{}
	decoder.Decode(product)

	fmt.Println(product)
}