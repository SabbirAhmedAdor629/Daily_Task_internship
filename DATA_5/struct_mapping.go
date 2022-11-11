package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	//"strings"
	//"strings"
)



type MESSAGE struct {
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
}
type Push struct{
	Status            string `json:"push_message_status"`
   	ProviderMessageID string `json:"push_message_provider_message_id"`
   	Provider          string `json:"push_message_pr
	
	ovider"`
   	CreatedAt         string `json:"push_message_created_at"`
   	UpdatedAt         string `json:"push_message_updated_at"`
   	GUID              string `json:"push_message_guid"`

}

type BONUS_MESSAGE struct{
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

type employee struct{
	Bonus_message BONUS_MESSAGE
	Message MESSAGE
	Push_message []PUSH_MESSAGE
}

func main() {
				// Reading the mapping.json file and decoding
	jsonFile, _ := os.Open("mapping.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result_map map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result_map)
				// Dynamically decoding every objects of mapping.json
	mapping := map[string]map[string]interface{}{}
	for k,v := range result_map{
		switch c := v.(type){
		case map[string]interface{}:
			mapping[k] = v.(map[string]interface{})
		default:
			fmt.Printf("Not sure what type item %q is, but I think it might be %T\n", k, c)
		}
	}

	var employee1 employee = employee{}

	//Reading the data.json file and decoding
	empJson, _ := os.Open("data.json")
	defer empJson.Close()
	byteValue_2, _ := ioutil.ReadAll(empJson)
	json.Unmarshal([]byte(byteValue_2), &employee1)
	fmt.Println(employee1.Message.Body)


	//Taking another map to store keys and values
	
	// for v := range result_map{
	// 	for item := range mapping[v] {
	// 		value := (strings.Split(fmt.Sprint((mapping[v][item])) ,"."))
	// 		strings.ToUpper(value[0])
	// 		strings.ToUpper(value[1])
	// 		fmt.Printf("%v : %v\n", item, employee1.item)
	// 		fmt.Println(" ")
	// 	}
	// }
}