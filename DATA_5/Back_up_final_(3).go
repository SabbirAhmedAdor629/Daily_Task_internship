package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	//"strings"
)

func main() {

	// Reading the mapping.json file and decoding
	jsonFile, _ := os.Open("mapping.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result_map map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result_map)
	// Dynamically decoding every objects of mapping.json
	mapping := map[string]map[string]interface{}{}
	for k, v := range result_map {
		switch c := v.(type) {
		case map[string]interface{}:
			mapping[k] = v.(map[string]interface{})
		default:
			fmt.Printf("Not sure what type item %q is, but I think it might be %T\n", k, c)
		}
	}

	//Reading the data.json file and decoding
	empJson, _ := os.Open("data.json")
	defer empJson.Close()
	byteValue_2, _ := ioutil.ReadAll(empJson)
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue_2), &result)
	// Dynamically decoding every object of data.json
	data := map[string]map[string]interface{}{}
	for k, v := range result {
		switch c := v.(type) {
		case map[string]interface{}:
			data[k] = v.(map[string]interface{})
		case []interface{}:
			for _, item := range v.([]interface{}) {
				data[k] = item.(map[string]interface{})
			}
		default:
			fmt.Printf("Not sure what type item %q is, but I think it might be %T\n", k, c)
		}
	}

	//Accessing values from data.json file
	// fmt.Println("GUID ", data["message"]["guid"])
	// fmt.Println("status ", data["push_messages"]["provider_message_id"])
	// fmt.Println("bonus_guid",data["bonus_message"]["bonus_guid"])
	// for v := range result{
	// 	for item := range data[v] {
	// 		fmt.Println("keys of data : ",item)
	// 		fmt.Println("Values : ",data[v][item])
	// 	}
	// }

	// Accessing keys and values from mapping.json file
	// fmt.Println(mapping["mapping"]["guid"])
	// for v := range result_map{
	// 	for item := range mapping[v] {
	// 		fmt.Println("keys : ",item)
	// 		fmt.Println("Values : ",mapping[v][item])
	// 	}
	// }

	// Taking another map to store keys and values
	CopiedMap := make(map[string]interface{})
	for v := range result_map {
		for item := range mapping[v] {

			value := (strings.Split(fmt.Sprint((mapping[v][item])), "."))
			CopiedMap[item] = data[value[0]][value[1]]
			fmt.Printf("%v : %v\n", item, CopiedMap[item])
			fmt.Println(" ")

		}
	}

	// Spliting
	// b := (strings.Split(fmt.Sprint((mapping["mapping"]["guid"])),"."))
	// fmt.Println(b[0])
	// fmt.Println(b[1])

	// // We need the keys
	// // Use make() to create the slice for better performance
	// Keys_of_message := make([]string, len(message))
	// Keys_of_bonus_message  := make([]string, len(bonus_message))
	// // Storing keys of message in to Keys_of_message slice
	// for key := range message {
	// 	Keys_of_message = append(Keys_of_message, key)
	// }
	// // Storing keys of bonus_message into Keys_of_bonus_message slice
	// for key := range bonus_message {
	// 	Keys_of_bonus_message  = append(Keys_of_bonus_message , key)
	// }

	// // Printing keys
	// for key := range result {
	// 	fmt.Println(result[key])
	// }
	// 		// push-message
	// for _,item := range push_message {
	// 	for key := range item.(map[string]interface{}){
	// 		fmt.Println(key)
	// 	}
	// 	fmt.Println(" ")

	// }
	// for k2 := range Keys_of_bonus_message  {
	// 	fmt.Println(Keys_of_bonus_message [k2])
	// }

	// // Printing values
	// fmt.Println(
	// 	"\nGuid :", message["guid"],
	// 	"\nbonus_guid :", bonus_message["bonus_guid"],
	// )
	// for _,item:=range push_message {
	// 	fmt.Printf("created_at : %v\n", item.(map[string]interface{})["created_at"])
	// 	fmt.Printf("status : %v\n", item.(map[string]interface{})["status"])
	// }

	// // Comparing the objects and storing data into new map
	// CopiedMap:= make(map[string]interface{})

	// fmt.Println(" ")
	// for key := range mapping{
	// 	k := 0
	// 	for key_2 := range bonus_message{
	// 		if (key == key_2){
	// 			CopiedMap[key] = bonus_message[key_2]
	// 			k++
	// 			fmt.Printf("%v : %v\n", key, CopiedMap[key])
	// 		}
	// 	}
	// 	if (k!=0){
	// 		continue
	// 	}
	// 	for key_2 := range message{
	// 		if (key == key_2){
	// 			CopiedMap[key] = message[key_2]
	// 			k++
	// 			fmt.Printf("%v : %v\n", key, CopiedMap[key])
	// 		}
	// 	}
	// }

	// println("          ")
	// println(2222222222222)

	// for _,item := range push_message {
	// 	for key_4 := range item.(map[string]interface{}){
	// 		for key := range mapping{
	// 			if (key == "push_messages[]."+key_4){
	// 				CopiedMap[key] = item.(map[string]interface{})[key_4]
	// 				fmt.Printf("%v : %v\n", key, CopiedMap[key])
	// 			}
	// 		}

	// 	}

	// }

}
