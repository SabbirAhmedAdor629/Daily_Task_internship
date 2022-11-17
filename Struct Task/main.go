package main

import (
	"encoding/json"
)

type Data struct {
	SkipInbox               bool        `json:"skip_inbox"`
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
	Push_message_body       string      `json:"push_message_body"`
	Template_body           string      `json:"template_body"`
}

type Push struct {
	Status             string `json:"status"`
	ProviderMessageId  string `json:"provider_message_id"`
	PushProvider       string `json:"push_provider"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
	MessageGuid        string `json:"message_guid"`
	Guid               string `json:"guid"`
	MessagePk          string `json:"message_pk"`
	DeviceId           string `json:"device_id"`
	PushRegistrationId string `json:"push_registration_id"`
	AwsMessageId       string `json:"aws_message_id"`
	AwsArnStatus       string `json:"aws_arn_status"`
}

type SQSPushInboxLambdaPayload struct {
	Origin            string      `json:"origin"`
	Operation         string      `json:"operation"`
	TransactionID     string      `json:"transaction_id"`
	Priority          string      `json:"priority"`
	SubmittedTs       string      `json:"submitted_ts"`
	Data              Data        `json:"data"`
	PushMessages      []Push      `json:"push_messages"`
	BonusGuid         string      `json:"bonus_guid"`
	AwardedPoints     json.Number `json:"awarded_points"`
	BonusCeatedAt     string      `json:"bonus_created_at"`
	BonusUpdatedAt    string      `json:"bonus_updated_at"`
	BonusTemplateId   json.Number `json:"bonus_template_id"`
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
	obj := SQSPushInboxLambdaPayload{
		Origin:        "engage",
		Operation:     "reward_point",
		TransactionID: "123123-123123-sadfasfd-asdfsadf",
		Priority:      "high",
		SubmittedTs:   "2022-04-20T20:46:31Z",
		Data: Data{
			SkipInbox:          true,
			GUID:               "mm-933cf415-5a5d-4e09-8def-a1c9c39ec246",
			PlayerID:           20802498,
			MessageTemplateID:  2628,
			Badge:              true,
			BgImage:            "",
			CallToAction:       "Try a NEW game now!! 🎲🍭🎰",
			CallToActionEvent:  "",
			CallToActionPoints: 0,

			CompletedAt:            "2022-07-28T16:28:48-07:00",
			CompleteOnCallToAction: "FALSE",
			Component:              "",

			ComponentParams:         "",
			CreatedAt:               "2022-07-28T16:28:48-07:00",
			CreatedDate:             "2022-07-28",
			Goto:                    "home",
			Image:                   "http://cdn.influencemobile.com/message_templates/images/000/002/628/header/header.jpeg?1657225678",
			MemberID:                "0a9b82cb-31f8-6452-a6bf-dc44a09342e7",
			NextMessage:             "",
			PushMessage:             "Dogs like water 💦 You like BONUS time!! 🎉",
			S3ProgramImageExtension: "",
			S3ProgramImageName:      "",
			Title:                   "Bonus time RIGHT now!!! ⏱",
			ViewPoints:              0,
			Push_message_body:       "Just like pups like water 💦, we heard you like BONUSES!  Here ya go!! 🥳 500 points to try a NEW game in the next 2 hours!  🎉",
			Template_body:           "sample template body message",
		},
		PushMessages: []Push{
			{
				Status:             "delivered",
				ProviderMessageId:  "12312",
				PushProvider:       "aws: :sn",
				CreatedAt:          "01-01-2022",
				UpdatedAt:          "01-01-2022",
				MessageGuid:        "ff9e7777-d4c6-4dcd-9422-40b4f6257370",
				Guid:               "mb-b23cf415-5a5d-4e90-8def-a1c9c39ec246",
				MessagePk:          "mm-933cf415-5a5d-4e09-8def-a1c9c39ec246",
				DeviceId:           "123213213",
				PushRegistrationId: "13123232",
				AwsMessageId:       "21321321",
				AwsArnStatus:       "Success",
			},
			{
				Status:             "delivered",
				ProviderMessageId:  "4321",
				PushProvider:       "aws::sn:32",
				CreatedAt:          "02-01-2022",
				UpdatedAt:          "02-01-2022",
				MessageGuid:        "ef9e7777-d4c6-4dcd-9422-40b4f6257360",
				Guid:               "mb-423cf415-5a5d-4e90-8def-a1c9c39ec248",
				MessagePk:          "mm-933cf415-5a5d-4e09-8def-a1c9c39ec246",
				DeviceId:           "43432132",
				PushRegistrationId: "123213213",
				AwsMessageId:       "21123213",
				AwsArnStatus:       "Success",
			},
		},
		BonusGuid:         "mb-b23cf415-5a5d-4e90-8def-a1c9c39ec246",
		AwardedPoints:     500,
		BonusCeatedAt:     "2022-07-28T18:28:48-07:00",
		BonusUpdatedAt:    "2022-07-28T16:28:48-07:00",
		BonusTemplateId:   210,
		BonusExpiresAt:    "2022-07-28T18:28:48-07:00",
		BonusExpiresIn:    43200,
		BonusAwardedAt:    "2022-07-28T18:28:48-07:00",
		BonusRewardTitle:  "Bonus time RIGHT now!!! ⏱",
		BonusCompletedAt:  "",
		BonusEventName:    "engage_install_server",
		BonusEventCount:   0,
		BonusEventCounter: 0,
		BonusType:         "prize",
	}
}
