package main

import (
	"encoding/json"
	"fmt"
	//"os"
)

type Push_Messages struct {
	status   string  
	provider_message_id      string
	push_provider string     
	created_at    string    
	updated_at       string     
	guid        string 
}

type Messages struct {
	player_id           int
	message_template_id int
	data                string
	badge               bool
}

type Bonus_Message struct {
	bonus_guid          string
	awarded_points      int
	bonus_created_at    string
	bonus_updated_at    string
	bonus_template_id   int
	bonus_expires_at    string
	bonus_expires_in    int
	bonus_awarded_at    string
	bonus_reward_title  string
	bonus_completed_at  string
	bonus_event_name    string
	bonus_event_count   int
	bonus_event_counter int
	bonus_type          string
}

type SNS_Event struct {
	guid                       string
	bg_image                   string
	body                       string
	call_to_action             string
	call_to_action_event       string
	call_to_action_points      string
	category_id                int
	completed_at               string
	complete_on_call_to_action string
	component                  string
	component_params           string
	created_at                 string
	created_date               string
	go_to                      string
	image                      string
	member_id                  string
	next_message               string
	push_message               string
	s3_program_image_extension string
	s3_program_image_name      string
	title                      string
	view_points                int

	Push_Message Push_Messages
	Message      Messages
	Bonus        Bonus_Message
}

func main() {
	var event SNS_Event

	jsondata := `{
		"guid" : "mm-b23cf415-5a5d-4e09-8def-a1c9c39ec246",
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
		"go_to" : "home",
		"image" : "http://cdn.influencemobile.com/message_templates/images/000/002/628/header/header.jpeg?1657225678",
		"member_id" : "0a9b82cb-31f8-6452-a6bf-dc44a09342e7",
		"next_message" :"",
		"push_message" : "Dogs like water 💦 You like BONUS time!! 🎉",
		"s3_program_image_extension" : "",
		"s3_program_image_name" : "",
		"title": "Bonus time RIGHT now!!! ⏱",
		"view_points" : 0, 
		"push_messages":{
			"status":"delivered",
			"provider_message_id":"12312",
			"push_provider": "aws::sn",
			"created_at":"01-01-2022",
			"updated_at":"01-01-2022",
			"guid":"mb-b23cf415-5a5d-4e90-8def-a1c9c39ec246"
		},
		"message": {
			"player_id" : 10802498,
			"message_template_id": 2628,
			"data": "Push message notification",
			"badge": true
		},
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

	err := json.Unmarshal([]byte(jsondata), &event)

	if err == nil {
		fmt.Println(event)
	} else {
		fmt.Println(err)
	}

	// if err == nil{
	// 	fmt.Println("Skip_Inbox: ",event.Skip_Inbox)
	// 	fmt.Println(" ")

	// 	for _,obj := range event.Push_Message{
	// 		fmt.Println("Member_Guid : ", obj.Member_Guid)
	// 		fmt.Println("Device_Id : ",obj.Device_Id)
	// 		fmt.Println("Push_Registration_Id : ",obj.Push_Registration_Id)
	// 		fmt.Println("Aws_Message_Id : ",obj.Aws_Message_Id)
	// 		fmt.Println("Aws_Arn_Status : ",obj.Aws_Arn_Status)
	// 		fmt.Println("\n")
	// 	}
	// 	fmt.Println(" ")

	// 	fmt.Println("Player_Id : ",event.Message.Player_Id)
	// 	fmt.Println("Message_Template_Id : ",event.Message.Message_Template_Id)
	// 	fmt.Println("Data : ",event.Message.Data)
	// 	fmt.Println("Badge : ",event.Message.Badge)

	// }else{
	// 	fmt.Println(err)
	// }
}
