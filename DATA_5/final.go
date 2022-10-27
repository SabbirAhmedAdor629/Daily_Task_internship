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

	// Comparing the objects and storing data into new map 
	CopiedMap:= make(map[string]interface{})

	fmt.Println(" ")
	for key := range mapping{
		k := 0
		for key_2 := range bonus_message{
			if (key == key_2){
				CopiedMap[key] = bonus_message[key_2]
				k++
				fmt.Printf("%v : %v\n", key, CopiedMap[key])
			}
		}
		if (k!=0){
			continue
		}
		for key_2 := range message{
			if (key == key_2){
				CopiedMap[key] = message[key_2]
				k++
				fmt.Printf("%v : %v\n", key, CopiedMap[key])
			}
		}
	}

	println("          ")
	println(2222222222222)
	
	for _,item := range push_message {
		for key_4 := range item.(map[string]interface{}){
			for key := range mapping{
				if (key == "push_messages[]."+key_4){
					CopiedMap[key] = item.(map[string]interface{})[key_4]
					fmt.Printf("%v : %v\n", key, CopiedMap[key])
				}
			}
			
		}
		
	}
	
}
