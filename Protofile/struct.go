package main

import "encoding/json"


type SQSPushInboxLambdaPayload struct{
	PushMessage PushMessage		`json:"push_message"`
	BonusMessage Bonus					`json:"bonus_message"`
	PushMessages []Push			`json:"push_messages"`
}



type PushMessage struct {
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
	PushMessageBody		string		`json:"push_message_body"`
	TemplateBody			string		`json:"template_body"`
}


type Bonus struct{
	Guuid         string      `json:"guid"`
	AwardedPoints json.Number `json:"awarded_points"`
	CeatedAt    string      `json:"created_at"`
	UpdatedAt    string      `json:"updated_at"`
	TemplateId   json.Number `json:"template_id"`
	ExpiresAt    string      `json:"expires_at"`
	ExpiresIn    json.Number `json:"expires_in"`
	AwardedAt    string      `json:"awarded_at"`
	RewardTitle  string      `json:"reward_title"`
	CompletedAt  string      `json:"completed_at"`
	EventName    string      `json:"event_name"`
	EventCount   json.Number `json:"event_count"`
	EventCounter json.Number `json:"event_counter"`
	Type         string      `json:"type"` 
}

type Push struct{
	Status            	string `json:"status"`
	ProviderMessageId 	string `json:"provider_message_id"`
	PushProvider        string `json:"push_provider"`
	CreatedAt         	string `json:"created_at"`
	UpdatedAt         	string `json:"updated_at"`
	MessageGuid         string `json:"message_guid"`
	Guid				string  `json:"guid"`
	Id 					string	`json:"id"`
	DeviceId  			string	`json:"device_id"`
	PushRegistrationId  string	`json:"push_registration_id"`
	AwsMessageId  		string	`json:"aws_message_id"`
	AwsArnStatus  		string	`json:"aws_arn_status"`
}
