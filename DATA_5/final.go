package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {

	// Reading the JSON file
	jsonFile, _ := os.Open("mapping.json") // Openning our json file

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	// Declared an empty interface
	var result_map map[string]interface{}

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal([]byte(byteValue), &result_map)

	fmt.Println(result_map)

	//Simple JSON object which we will parse
	empJson := `{
		"message": {
			"guid" : "mm-b23cf415-5a5d-4e09-8def-a1c9c39ec246",
			"player_id" : 10802498,
			"message_template_id": 2628,
			"data": "Push message notification",
			"badge": true,
			"bg_image" : "",
			"body" : "Just like pups like water 💦, we heard you like BONUSES!  Here ya go!! 🥳 500 points to try a NEW game in the next 2 hours!  🎉",
			"call_to_action" : "Try a NEW game now!! 🎲🍭🎰",
			"call_to_action_event" : "",
			"call_to_action_points" : 0,
			"category_id" : 0,
			"completed_at" : "2022-07-28T16:28:48-07:00",
		   
			"complete_on_call_to_action": "FALSE",
			"component" : "",
			"component_params" : "",
			"created_at" : "2022-07-28T16:28:48-07:00",
			"created_date": "2022-07-28",
			"goto" : "home",
			"image" : "http://cdn.influencemobile.com/message_templates/images/000/002/628/header/header.jpeg?1657225678",
			"member_id" : "0a9b82cb-31f8-6452-a6bf-dc44a09342e7",
		   
			"next_message" :"",
			"push_message" : "Dogs like water 💦 You like BONUS time!! 🎉",
			"s3_program_image_extension" : "",
			"s3_program_image_name" : "",
			"title": "Bonus time RIGHT now!!! ⏱",
			"view_points" : 0
		},

		"push_messages":[
			{

				"status":"delivered",
				"provider_message_id":"12312",
				"push_provider": "aws::sn",
				"created_at":"01-01-2022",
				"updated_at":"01-01-2022",
				"guid":"mb-b23cf415-5a5d-4e90-8def-a1c9c39ec246"
			},
			{

				"status":"delivered",
				"provider_message_id":"12312",
				"push_provider": "aws::sn",
				"created_at":"01-01-2022",
				"updated_at":"01-01-2022",
				"guid":"mb-b23cf415-5a5d-4e90-8def-a1c9c39ec246"
			}
    	],

		"bonus_message":{

			"bonus_guid" :"mb-b23cf415-5a5d-4e90-8def-a1c9c39ec246",
			"awarded_points" : 500,
			"bonus_created_at" : "2022-07-28T16:28:48-07:00",
			"bonus_updated_at" : "2022-07-28T16:28:48-07:00",
			"bonus_template_id" : 210,
			"bonus_expires_at" : "2022-07-28T18:28:48-07:00",
			"bonus_expires_in" : 2,
			"bonus_awarded_at" : "2022-07-28T18:28:48-07:00",
			"bonus_reward_title" : "Bonus time RIGHT now!!! ⏱",
			"bonus_completed_at" : "",
			"bonus_event_name" : "engage_install_server",
			"bonus_event_count" : 0,
			"bonus_event_counter" : 0,
			"bonus_type" : "prize"
		}
	}`

	// Declared an empty interface
	var result map[string]interface{}

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal([]byte(empJson), &result)

	message := result["message"].(map[string]interface{})
	bonus_message := result["bonus_message"].(map[string]interface{})
	push_message := result["push_messages"].([]interface{})

	// Use make() to create the slice for better performance
	translationKeys := make([]string, len(message))
	translationKeys_bonus_message := make([]string, len(bonus_message))

	// We only need the keys
	// keys of message
	for key := range message {
		translationKeys = append(translationKeys, key)
	}
	// keys of bonus_message
	for key := range bonus_message {
		translationKeys_bonus_message = append(translationKeys_bonus_message, key)
	}

	// Printing keys
	for k := range translationKeys {
		fmt.Println(translationKeys[k])
	}
	for k2 := range translationKeys_bonus_message {
		fmt.Println(translationKeys_bonus_message[k2])
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

}
