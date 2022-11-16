package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
)


type Data struct {
	SkipInbox				bool		`json:"skip_inbox"`
	GUID                    string      `json:"guid"`
	PlayerID                json.Number `json:"player_id"`
	MessageTemplateID       json.Number `json:"message_template_id"`
	Badge                   bool        `json:"badge"`
	BgImage                 string      `json:"bg_image"`
	CallToAction            string      `json:"call_to_action"`
	CallToActionEvent       string      `json:"call_to_action_event"`
	CallToActionPoints      json.Number `json:"call_to_action_points"`
	CompletedAt             string      `json:"completed_at"`
	CompleteOnCallToAction  string      `json:"complete_on_call_to_action"`
	Component               string      `json:"component"`
	ComponentParams         string      `json:"component_params"`
	CreatedAt               string      `json:"created_at"`
	UpdatedAt    			string      `json:"updated_at"`
	Goto                    string      `json:"goto"`
	Image                   string      `json:"image"`
	MemberID                string      `json:"member_id"`
	NextMessage             string      `json:"next_message"`
	PushMessage             string      `json:"push_message"`
	S3ProgramImageExtension string      `json:"s3_program_image_extension"`
	S3ProgramImageName      string      `json:"s3_program_image_name"`
	Title                   string      `json:"title"`
	ViewPoints              json.Number `json:"view_points"`
	Push_message_body		string		`json:"push_message_body"`
	Template_body			string		`json:"template_body"`
}




type Push struct{
	Status            	string `json:"status"`
	ProviderMessageId 	string `json:"provider_message_id"`
	PushProvider        string `json:"push_provider"`
	CreatedAt         	string `json:"created_at"`
	UpdatedAt         	string `json:"updated_at"`
	MessageGuid         string `json:"message_guid"`
	Guid				string `json:"guid"`
	DeviceId  			string `json:"device_id"`
	AwsMessageId  		string `json:"aws_message_id"`
	AwsArnStatus  		string `json:"aws_arn_status"`
}



type SQSPushInboxLambdaPayload struct{

	Origin string			`json:"origin"`
	Operation string		`json:"operation"`
	TransactionID string	`json:"transaction_id"`
	SubmittedTs string		`json:"submitted_ts"`
	
	Data Data					`json:"data"`
	PushMessages []Push			`json:"push_messages"`
	
	BonusGuid		string      `json:"bonus_guid"`
	AwardedPoints	json.Number `json:"awarded_points"`
	CeatedAt		string     	`json:"created_at"`
	UpdatedAt		string		`json:"updated_at"`
	TemplateId		json.Number `json:"template_id"`
	ExpiresAt		string      `json:"expires_at"`
	ExpiresIn		json.Number `json:"expires_in"`
	AwardedAt		string      `json:"awarded_at"`
	RewardTitle		string      `json:"reward_title"`
	CompletedAt		string      `json:"completed_at"`
	EventName		string      `json:"event_name"`
	EventCount		json.Number `json:"event_count"`
	EventCounter	json.Number `json:"event_counter"`
	BonusType		string      `json:"bonus_type"`
}



// type LambdaEvent struct {
// 	Name string `json:"name"`
// 	Age int `json:"age"`
// }


// type LambdaResponse struct {
// 	Message string `json:"message"`
// }


func LambdaHandler(event SQSPushInboxLambdaPayload) (string, error) {
	return fmt.Sprintf(
		"Origin : %s SkipInbox : %s BonusGuid : %s",event.Origin, event.Data.GUID,event.BonusGuid,
		), nil
	
}



func main ()  {
	lambda.Start(LambdaHandler)
	
}

				// ACCESSING VALUES OF EVERY OBJECTS IN THE ARRAY OF STRUCT
	//var event SQSPushInboxLambdaPayload
	// for _, i := range event.PushMessages{
	// 	i.DeviceId = "kdjfsdf"
	// }

	// for _, i := range event.PushMessages{
	// 	fmt.Println(i.DeviceId)
	// }