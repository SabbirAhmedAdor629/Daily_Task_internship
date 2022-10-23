package main

import "encoding/json"



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
	
	
	PushMessage    []push_message  

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

type push_message struct {
	status	string 				`json:"status"`
	provider_message_id	string	`json:"provider_message_id"`
	push_provider	string		`json:"push_provider"`
	created_at					`json:"created_at"`
	updated_at					`json:"updated_at"`
	guid	string				`json:"guid"`
}