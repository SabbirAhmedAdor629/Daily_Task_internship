package main

import (
	"encoding/json"
	"fmt"
)

type DynamoInboxStoreItem struct {
	GUID                    string      `json:"guid"`
	PlayerID                json.Number `json:"player_id"`
	MessageTemplateID       json.Number `json:"message_template_id"`
	Data                    string      `json:"data"`
	Badge                   bool        `json:"badge"`
	BgImage                 string      `json:"bg_image"`
	Body                    string      `json:"body"`
	CallToAction            string      `json:"call_to_action"`
	CallToActionEvent       string      `json:"call_to_action_event"`
	CallToActionPoints      json.Number `json:"call_to_action_points"`
	CategoryID              json.Number `json:"category_id"`
	CompletedAt             string      `json:"completed_at"`
	CompleteOnCallToAction  string      `json:"complete_on_call_to_action"`
	Component               string      `json:"component"`
	ComponentParams         string      `json:"component_params"`
	CreatedAt               string      `json:"created_at"`
	CreatedDate             string      `json:"created_date"`
	Goto                    string      `json:"goto"`
	Image                   string      `json:"image"`
	MemberID                string      `json:"member_id"`
	NextMessage             string      `json:"next_message"`
	PushMessage             string      `json:"push_message"`
	S3ProgramImageExtension string      `json:"s3_program_image_extension"`
	S3ProgramImageName      string      `json:"s3_program_image_name"`
	Title                   string      `json:"title"`
	ViewPoints              json.Number `json:"view_points"`

	PushMessageStatus            string `json:"push_message_status"`
	PushMessageProviderMessageID string `json:"push_message_provider_message_id"`
	PushMessageProvider          string `json:"push_message_provider"`
	PushMessageCreatedAt         string `json:"push_message_created_at"`
	PushMessageUpdatedAt         string `json:"push_message_updated_at"`
	PushMessageGUID              string `json:"push_message_guid"`

	BonusGUID         string      `json:"bonus_guid"`
	AwardedPoints     json.Number `json:"awarded_points"`
	BonusCreatedAt    string      `json:"bonus_created_at"`
	BonusUpdatedAt    string      `json:"bonus_updated_at"`
	BonusTemplateID   json.Number `json:"bonus_template_id"`
	BonusExpiresAt    string      `json:"bonus_expires_at"`
	BonusExpiresIn    json.Number `json:"bonus_expires_in"`
	BonusAwardedAt    string      `json:"bonus_awarded_at"`
	BonusRewardTitle  string      `json:"bonus_reward_title"`
	BonusCompletedAt  string      `json:"bonus_completed_at"`
	BonusEventName    string      `json:"bonus_event_name"`
	BonusEventCount   json.Number `json:"bonus_event_count"`
	BonusEventCounter json.Number `json:"bonus_event_counter"`
	BonusType         string      `json:"bonus_type"`
}

func main() {
	//Simple Employee JSON object which we will parse
	empJson := `{
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
		"view_points" : 0, 
		"push_messages":{
			"status":"delivered",
			"provider_message_id":"12312",
			"push_provider": "aws::sn",
			"created_at":"01-01-2022",
			"updated_at":"01-01-2022",
			"guid":"mb-b23cf415-5a5d-4e90-8def-a1c9c39ec246"
		},
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
}`

	// Declared an empty interface
	var result map[string]interface{}

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal([]byte(empJson), &result)

	push_messages := result["push_messages"].(map[string]interface{})



	





	//Reading each value by its key

	fmt.Println(
		"\nPushMessageStatus :", push_messages["status"], 
		"\nPushMessageProviderMessageID :", push_messages["provider_message_id"], 
		"\nPushMessageProvider :", push_messages["push_provider"])
		
	

	
}