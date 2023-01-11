package main


import (
	"encoding/json"
)

type SQSPushLambdaPayload struct {
	Origin        string  `json:"origin"`
	Operation     string  `json:"operation"`
	TransactionId string  `json:"transaction_id"`
	SubmittedTs   string  `json:"submitted_ts"`
	Data          Message `json:"data"`
}

type Message struct {
	SkipInbox               bool           `json:"skip_inbox"`
	GUID                    string         `json:"guid"`
	PlayerID                int64          `json:"player_id"`
	MessageTemplateID       json.Number    `json:"message_template_id"`
	Badge                   bool           `json:"badge"`
	BgImage                 string         `json:"bg_image"`
	CallToAction            string         `json:"call_to_action"`
	CallToActionEvent       string         `json:"call_to_action_event"`
	CallToActionPoints      json.Number    `json:"call_to_action_points"`
	CompletedAt             string         `json:"completed_at"`
	CompleteOnCallToAction  string         `json:"complete_on_call_to_action"`
	Component               string         `json:"component"`
	ComponentParams         string         `json:"component_params"`
	CreatedAt               string         `json:"created_at"`
	UpdatedAt               string         `json:"updated_at"`
	ViewedAt                string         `json:"viewed_at"`
	CampaignId              json.Number    `json:"campaign_id"`
	CreatedDate             string         `json:"created_date"`
	CreatedHour             json.Number    `json:"created_hour"`
	PushProvider            string         `json:"push_provider"`
	InAppPush               bool           `json:"in_app_push"`
	CausedUninstall         bool           `json:"caused_uninstall"` // default false
	RawPush                 bool           `json:"raw_push"`
	OutboundNumber          string         `json:"outbound_number"`
	PushNotification        bool           `json:"push_notification"`
	PushResponse            string         `json:"push_response"`
	PushDelivered           bool           `json:"push_delivered"`
	Goto                    string         `json:"goto"`
	Image                   string         `json:"image"`
	MemberID                string         `json:"member_id"`
	NextMessage             string         `json:"next_message"`
	PushMessage             string         `json:"push_message"`
	S3ProgramImageExtension string         `json:"s3_program_image_extension"`
	S3ProgramImageName      string         `json:"s3_program_image_name"`
	Title                   string         `json:"title"`
	ViewPoints              json.Number    `json:"view_points"`
	Body                    string         `json:"body"`
	TemplateBody            string         `json:"template_body"`
	BonusGUID               string         `json:"bonus_guid"`
	AwardedPoints           json.Number    `json:"awarded_points"`
	BonusCreatedAt          string         `json:"bonus_created_at"`
	BonusUpdatedAt          string         `json:"bonus_updated_at"`
	BonusTemplateID         json.Number    `json:"bonus_template_id"`
	BonusExpiresAt          string         `json:"bonus_expires_at"`
	BonusExpiresIn          json.Number    `json:"bonus_expires_in"`
	BonusAwardedAt          string         `json:"bonus_awarded_at"`
	BonusRewardTitle        string         `json:"bonus_reward_title"`
	BonusCompletedAt        string         `json:"bonus_completed_at"`
	BonusEventName          string         `json:"bonus_event_name"`
	BonusEventCount         json.Number    `json:"bonus_event_count"`
	BonusEventCounter       json.Number    `json:"bonus_event_counter"`
	BonusType               string         `json:"bonus_type"`
	PushMessages            []PushMessages `json:"push_messages"`
}

type PushMessages struct {
	Status                   string `json:"status"`
	ProviderMessageID        string `json:"provider_message_id"`
	PushProvider             string `json:"push_provider"`
	CreatedAt                string `json:"created_at"`
	UpdatedAt                string `json:"updated_at"`
	MessageUUID              string `json:"message_uuid"`
	GUID                     string `json:"guid"`
	MessagePK                string `json:"message_pk"`
	DeviceId                 string `json:"device_id"`
	PushRegistrationId       string `json:"push_registration_id"`
	AwsMessageId             string `json:"aws_message_id"`
	AwsArnStatus             string `json:"aws_arn_status"`
	AppReceivedAt            string `json:"app_received_at"`
	PushNotificationsEnabled string `json:"push_notifications_enabled"`
	AppProcessingStatus      string `json:"app_processing_status"`
}


func createSNSMessage(sqsPushPayload SQSPushLambdaPayload, pushMessagesIndex int) map[string]interface{} {
	message_payload := make(map[string]interface{})

	if pushMessagesIndex < 0 || len(sqsPushPayload.Data.PushMessages) <= pushMessagesIndex {
		return nil
	}

	message_payload["GCM"] = map[string]interface{}{
		"notification": map[string]interface{}{
			"title":        sqsPushPayload.Data.Title,
			"body":         sqsPushPayload.Data.Body,
			"goto":         sqsPushPayload.Data.Goto,
			"image":        sqsPushPayload.Data.Image,
			"message_uuid": sqsPushPayload.Data.PushMessages[pushMessagesIndex].MessageUUID,
		},
		"data": map[string]interface{}{
			"body":         sqsPushPayload.Data.Body,
			"title":        sqsPushPayload.Data.Title,
			"icon":         "ic_notification",
			"largeIcon":    "ic_launcher",
			"goto":         sqsPushPayload.Data.Goto,
			"image":        sqsPushPayload.Data.Image,
			"color":        "#06192E",
			"sound":        "generic2",
			"channel":      "default",
			"message_uuid": sqsPushPayload.Data.PushMessages[pushMessagesIndex].MessageUUID,
		},
		"android": map[string]interface{}{
			"android": map[string]interface{}{
				"importance": "default",
			},
		},
	}

	// SNS Message JSON schema
	push_params := map[string]interface{}{
		"target_arn":        "device.aws_endpoint_arn",
		"message":           message_payload,
		"message_structure": "json",
	}

	return push_params


	
}