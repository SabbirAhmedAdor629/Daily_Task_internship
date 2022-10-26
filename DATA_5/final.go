package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {

	// Reading the mapping.json file and decoding
	jsonFile, _ := os.Open("mapping.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result_map map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result_map)
	// Have to decode mapping struct seperately
	mapping := result_map["mapping"].(map[string]interface{})
	
	//Accessing keys of mappping struct from mappping.json file
	for key := range mapping {
		fmt.Println(key)
	}

	
	//Reading the data.json file and decoding
	empJson, _ := os.Open("data.json")
	defer empJson.Close()
	byteValue_2, _ := ioutil.ReadAll(empJson)
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue_2), &result)
	// Have to decode three struct seperately
	message := result["message"].(map[string]interface{})
	push_message := result["push_messages"].([]interface{})
	bonus_message := result["bonus_message"].(map[string]interface{})

	// We need the keys
	// Use make() to create the slice for better performance
	Keys_of_message := make([]string, len(message))
	Keys_of_bonus_message  := make([]string, len(bonus_message))	
	// Storing keys of message in to Keys_of_message slice
	for key := range message {
		Keys_of_message = append(Keys_of_message, key)
	}
	// Storing keys of bonus_message into Keys_of_bonus_message slice
	for key := range bonus_message {
		Keys_of_bonus_message  = append(Keys_of_bonus_message , key)
	}

	// Printing keys
	for key := range Keys_of_message {
		fmt.Println(Keys_of_message[key])
	}
			// push-message
	for _,item := range push_message {
		for key := range item.(map[string]interface{}){
			fmt.Println(key)
		}
		fmt.Println(" ")

	}
	for k2 := range Keys_of_bonus_message  {
		fmt.Println(Keys_of_bonus_message [k2])
	}


	// Printing values
	fmt.Println(
		"\nGuid :", message["guid"],
		"\nbonus_guid :", bonus_message["bonus_guid"],

	)

	for _,item:=range push_message {
		fmt.Printf("created_at : %v\n", item.(map[string]interface{})["created_at"])
		fmt.Printf("status : %v\n", item.(map[string]interface{})["status"])
	}

	// Comparing the objects and storing data 
	fmt.Println(" ")
	for key := range mapping{
		for key_2 := range bonus_message{
			if (key == key_2){
				mapping[key] = bonus_message[key_2]
				fmt.Printf("%v : %v\n", key, mapping[key])
			}
		}
	}



	finaljson, err := json.MarshalIndent(result_map, "", "\t")
	if err != nil {
		panic(err)
	}
	_ = ioutil.WriteFile("newMapping.json", finaljson, 0777)
	


}
