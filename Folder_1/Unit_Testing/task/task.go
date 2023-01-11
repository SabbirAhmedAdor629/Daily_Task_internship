package task

import "encoding/json"


type SQSPushLambdaPayload struct {
	PushMessage  PushMessage `json:"push_message"`
	BonusMessage Bonus       `json:"bonus_message"`
	PushMessages []Push      `json:"push_messages"`
}

type PushMessage struct {
	SkipInbox               bool        `json:"skip_inbox"`
	GUID                    string      `json:"guid"`
	PlayerID                int64       `json:"player_id"`
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
	Body                    string      `json:"body"`
	TemplateBody            string      `json:"template_body"`
}

type Bonus struct {
	Guid          string      `json:"guid"`
	AwardedPoints json.Number `json:"awarded_points"`
	CreatedAt     string      `json:"created_at"`
	UpdatedAt     string      `json:"updated_at"`
	TemplateId    json.Number `json:"template_id"`
	ExpiresAt     string      `json:"expires_at"`
	ExpiresIn     json.Number `json:"expires_in"`
	AwardedAt     string      `json:"awarded_at"`
	RewardTitle   string      `json:"reward_title"`
	CompletedAt   string      `json:"completed_at"`
	EventName     string      `json:"event_name"`
	EventCount    json.Number `json:"event_count"`
	EventCounter  json.Number `json:"event_counter"`
	Type          string      `json:"type"`
}

type Push struct {
	Status             string `json:"status"`
	ProviderMessageId  string `json:"provider_message_id"`
	PushProvider       string `json:"push_provider"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
	MessageGuid        string `json:"message_guid"`
	Guid               string `json:"guid"`
	DeviceId           string `json:"device_id"`
	PushRegistrationId string `json:"push_registration_id"`
	AwsMessageId       string `json:"aws_message_id"`
	AwsArnStatus       string `json:"aws_arn_status"`
}


func main(){
	var parameter SQSPushLambdaPayload
	CreateInboxQueueMessage(parameter)
}




var Data = make(map[string]interface{})
func CreateInboxQueueMessage(sqsPushLambdaPayload SQSPushLambdaPayload) map[string]interface{} {
	sqsInboxPayload := make(map[string]interface{})

	sqsInboxPayload["origin"] = "engage"
	sqsInboxPayload["operation"] = "reward_point"
	sqsInboxPayload["transaction_id"] = "a80edda3-607d-4c85-bf24-fa5aea8a80bc"
	sqsInboxPayload["submitted_ts"] = "2022-04-20T20:46:31Z"

	// converting SQSPushLambdaPayload to map[string]interface{}
	var payloadDataMap = make(map[string]interface{})
	payload, _ := json.Marshal(sqsPushLambdaPayload)
	json.Unmarshal(payload, &payloadDataMap)

	// storing SQSPushLambdaPayload data to inbox queue data format
	
	for key, value := range payloadDataMap["push_message"].(map[string]interface{}) {
		Data[key] = value
	}

	Data["push_messages"] = payloadDataMap["push_messages"]

	prefix := "bonus_"
	for key, value := range payloadDataMap["bonus_message"].(map[string]interface{}) {
		if key == "awarded_points" {
			Data[key] = value
			continue
		}
		Data[prefix+key] = value
	}

	sqsInboxPayload["data"] = Data

	return sqsInboxPayload
}




// func Add(A int, B int) int{
// 	return A+B;
// }