package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	. "influencemobile.com/libs/dynamo"
)

func sendMessage(params map[string]interface{}) (bool, error) {
	messageTemplate, err := pgDb.ReadMessageTemplateData(params["message_template_id"].(int)) // read message_template data form rds
	if err != nil || messageTemplate == nil {
		return false, err
	}

	bonusTemplate, err := pgDb.ReadBonusTemplateData(messageTemplate.BonusTemplateId) // read bonus_template data for rds
	if err != nil {
		return false, err
	}

	if bonusTemplate != nil {
		isBonusAvailable, err := pgDb.BonusAvailable(params["player_id"].(int), time.Now().UTC().Format(time.RFC3339))
		if err != nil {
			return false, err
		}
		if isBonusAvailable {
			return false, err
		}
	}
	payload, err := sendTo(params, bonusTemplate, messageTemplate)
	if err != nil {
		return false, err
	}
	// logger.WithStruct(payload).Debug("SQS Message Payload")
	return false, nil
}

func sendTo(params map[string]interface{}, bonusTemplate *BonusTemplate, messageTemplate *MessageTemplate) (map[string]interface{}, error) {
	bonus, err := createBonus(params["player_id"].(int), bonusTemplate, messageTemplate)
	if err != nil {
		fmt.Println("Send_to method error", err)
	}
	logger.WithStruct(bonus).Debug("Bonus Message")

	messagePayload := make(map[string]interface{})

	messagePayload["guid"] = assignGuid("Message")
	messagePayload["created_at"] = time.Now().UTC().Format(time.RFC3339)
	messagePayload["created_date"] = time.Now().UTC().Format("2006-01-02")
	messagePayload["created_hour"] = time.Now().Hour()
	messagePayload["caused_uninstall"] = false //need to check from rails app
	messagePayload["raw_push"] = false         //need to check from rails app
	messagePayload["player_id"] = params["player_id"]
	messagePayload["campaign id"] = params["campaign_id"]
	messagePayload["message_template_id"] = messageTemplate.Id
	messagePayload["badge"] = true // by default true

	messageTemplate = macro(params["player_id"].(int), messageTemplate)

	if ok, _ := regexp.MatchString("missing", messageTemplate.BgImage); !ok {
		messagePayload["bg_image"] = messageTemplate.BgImage
	}
	if ok, _ := regexp.MatchString("missing", messageTemplate.Image); !ok {
		messagePayload["image"] = imageUrlFormat(messageTemplate.Image)
	}

	messagePayload["body"] = messageTemplate.Body
	messagePayload["call_to_action"] = messageTemplate.CallToAction

	if messageTemplate.CallToActionEvent.Valid {
		messagePayload["call_to_action_event"] = messageTemplate.CallToActionEvent.String
	} else {
		messagePayload["call_to_action_event"] = ""
	}
	messagePayload["call_to_action_points"] = messageTemplate.CallToActionPoints
	messagePayload["completed_at"] = nil // define nil according to excel sheets
	messagePayload["complete_on_call_to_action"] = messageTemplate.CompleteOnCallToAction
	messagePayload["component"] = messageTemplate.Component
	messagePayload["component_params"] = messageTemplate.ComponentParams.String

	if messageTemplate.ComponentParams.Valid {
		messagePayload["component_params"] = messageTemplate.ComponentParams.String
	} else {
		messagePayload["component_params"] = ""
	}

	messagePayload["goto"] = messageTemplate.Goto
	messagePayload["image"] = messageTemplate.Image
	// messagePayload["member_id"] = messageTemplate.member_id  missing form member_id in message-template
	nexMessageTemplate, err := pgDb.ReadMessageTemplateData(messageTemplate.NextMessageTemplateId)
	if err != nil {
		messagePayload["next_message"] = ""
	} else {
		messagePayload["next_message"] = nexMessageTemplate.Name
	}

	if messageTemplate.Name.Valid {
		messagePayload["push_message"] = messageTemplate.Name.String
	} else {
		messagePayload["push_message"] = messageTemplate.PushMessage.String
	}
	messagePayload["s3_program_image_extension"] = messageTemplate.S3ProgramImageExtension
	messagePayload["s3_program_image_name"] = messageTemplate.S3ProgramImageName
	messagePayload["title"] = messageTemplate.Title
	messagePayload["view_points"] = messageTemplate.ViewPoints
	messagePayload["updated_at"] = time.Now().UTC().Format(time.RFC3339)
	messagePayload["viewed_at "] = nil

	// messagePayload["campaign id"] = "campaig" //should come from upstream
	// messagePayload["push response"] = "push response" //should come from upstream

	messagePayload["push_provider"] = "aws::sns"

	if messageTemplate.SkipInbox.Valid {
		messagePayload["skip_inbox"] = messageTemplate.SkipInbox.Bool
	} else {
		messagePayload["skip_inbox"] = false
	}

	if messageTemplate.InAppPush.Valid {
		messagePayload["in_app_push"] = messageTemplate.InAppPush.Bool
	} else {
		messagePayload["in_app_push"] = true
	}

	if messageTemplate.PushNotification.Valid {
		messagePayload["push_notification"] = messageTemplate.PushNotification.Bool
	} else {
		messagePayload["push_notification"] = true
	}

	pushMessagePayloadList, memberIds := push_all(messageTemplate, playerId)
	if len(pushMessagePayloadList) > 0 {
		messagePayload["push_delivered"] = true
	} else {
		messagePayload["push_delivered"] = false
	}

	if len(memberIds) > 1 {
		messagePayload["member_id"] = nil
	}

	payload := map[string]interface{}{}
	payload["data"] = messagePayload
	payload["push_message"] = pushMessagePayloadList
	payload["bonus"] = bonus

	return nil, nil
}

func macro(playerId int, messageTemplate *MessageTemplate) *MessageTemplate {
	// fields := []string{"Title", "Body", "PushMessage", "Goto", "CallToAction"}
	fields := []string{"title", "body", "push_message", "goto", "call_to_action"}

	vars := structToArrayConvert(messageTemplate)
	vars["attgetter"] = attentionGetter()

	// re := regexp.MustCompile("first_name")

	// for _, field := range fields {
	// 	if re.MatchString(vars[field].(string)) {
	// 		vars["first_name"] = firstName("test_id_123")
	// 		break
	// 	}
	// }

	for _, field := range fields {
		switch field {
		case "title":
			messageTemplate.Title = replaceText(messageTemplate.Title, vars)
		case "body":
			messageTemplate.Body = replaceText(messageTemplate.Body, vars)
		case "push_message":
			if messageTemplate.PushMessage.Valid {
				messageTemplate.PushMessage.String = replaceText(messageTemplate.PushMessage.String, vars)
			}
		case "goto":
			messageTemplate.Goto = replaceText(messageTemplate.Goto, vars)
		case "call_to_action":
			messageTemplate.CallToAction = replaceText(messageTemplate.CallToAction, vars)
		}
	}
	return messageTemplate
}

func structToArrayConvert(messageTemplate *MessageTemplate) map[string]interface{} {

	vars := make(map[string]interface{})

	vars["name"] = messageTemplate.Name.String
	vars["title"] = messageTemplate.Title
	vars["body"] = messageTemplate.Body
	vars["call_to_action"] = messageTemplate.CallToAction
	vars["goto"] = messageTemplate.Goto
	vars["image_file_name"] = messageTemplate.Image
	vars["image_content_type"] = messageTemplate.ImageContentType
	vars["image_file_size"] = messageTemplate.ImageFileSize.Int32 //int
	vars["image_updated_at"] = messageTemplate.ImageUpdatedAt.String
	vars["push_message"] = messageTemplate.PushMessage.String
	vars["bonus_template_id"] = messageTemplate.BonusTemplateId //int
	vars["event_name"] = messageTemplate.EventName.String       //need
	vars["event_count"] = messageTemplate.EventCount.Int32      //int
	vars["daily_cap"] = messageTemplate.DailyCap.Int32          //int
	vars["monthly_cap"] = messageTemplate.MonthlyCap.Int32      //int
	vars["total_cap"] = messageTemplate.TotalCap.Int32          //int
	vars["category_id"] = messageTemplate.CategoryId            //int
	vars["layout"] = messageTemplate.Layout
	vars["component"] = messageTemplate.Component
	vars["call_to_action_event"] = messageTemplate.CallToActionEvent.String
	vars["bg_image_file_name"] = messageTemplate.BgImage
	vars["bg_image_content_type"] = messageTemplate.BgImageContentType
	vars["bg_image_file_size"] = messageTemplate.BgImageFileSize.Int32    //int
	vars["bg_image_updated_at"] = messageTemplate.BgImageUpdatedAt.String //dateTime
	vars["next_message"] = messageTemplate.NextMessage.String
	vars["call_to_action_points"] = messageTemplate.CallToActionPoints //int
	vars["component_params"] = messageTemplate.ComponentParams.String
	vars["view_points"] = messageTemplate.ViewPoints //int
	vars["s3_program_image_name"] = messageTemplate.S3ProgramImageName.String
	vars["s3_program_image_extension"] = messageTemplate.S3ProgramImageExtension.String
	vars["on_complete_next_message"] = messageTemplate.OnCompleteNextImage.String                  //string
	vars["delay_next_message"] = messageTemplate.DelayNextMessage.Int32                            //int
	vars["delay_next_message_units"] = messageTemplate.DelayNextMessageUnits.String                //string
	vars["delay_next_message_time"] = messageTemplate.DelayNextMessageTime.String                  //string
	vars["on_complete_event_name"] = messageTemplate.OncompleteEventName.String                    //string
	vars["next_message_template_id"] = messageTemplate.NextMessageTemplateId                       //int
	vars["on_complete_next_message_template_id"] = messageTemplate.OnCompleteNextMessageTemplateId //int
	vars["complete_on_call_to_action"] = messageTemplate.CompleteOnCallToAction                    //bool
	vars["active"] = messageTemplate.Active                                                        //string
	vars["in_app_push"] = messageTemplate.InAppPush.Bool
	vars["skip_inbox"] = messageTemplate.SkipInbox.Bool               //string
	vars["push_notification"] = messageTemplate.PushNotification.Bool //int
	vars["onboarding_state"] = messageTemplate.OnboardingState.String //string

	return vars
}

func replaceText(text string, vars map[string]interface{}) string {
	pattern := regexp.MustCompile(`\[[^\]]*\]`)
	replacedText := pattern.ReplaceAllStringFunc(text, func(match string) string {
		// Remove the brackets from the match
		key := strings.Trim(match, "[]")
		// Look up the value in vars, or return an empty string if it doesn't exist
		if vars[key] == nil {
			return ""
		}
		return vars[key].(string)
	})
	// fmt.Println("--replace text: ", replacedText)
	return replacedText
}

func firstName(member_id string) string {
	lr_member := lrMember(member_id)
	if lr_member["first_name"] != nil {
		return lr_member["first_name"].(string)
	} else {
		return "Sports Fan"
	}
}

func lrMember(member_id string) map[string]interface{} {

	if member_id == "" {
		return nil
	}
	redisKey := fmt.Sprint("lr/member/v3/", member_id)
	val, err := redisDB.Get(ctx, redisKey).Result()
	member := make(map[string]interface{})

	if err != nil { // not set yet
		member["first_name"] = "hafizul"
		member["last_name"] = "islam"
		// QetAllItems(dynamoSvc, "members_"+env.Environment, "player_id",playerId)
		memberJsonData, err := json.Marshal(member)
		if err != nil {
			return nil
		}
		err = redisDB.Set(ctx, redisKey, memberJsonData, 0).Err()
		if err != nil {
			return nil
		}
		if _, err := redisDB.Expire(ctx, redisKey, time.Hour*time.Duration(24)).Result(); err != nil {
			return nil
		}
	}
	json.Unmarshal([]byte(val), &member)
	return member
}

func imageUrlFormat(image_url string) string {
	if strings.HasPrefix(image_url, "http://cdn.influencemobile.com") {
		return image_url
	}

	new_url := image_url

	if strings.Contains(image_url, "https") {
		new_url = strings.Replace(new_url, "https", "http", -1)
	} else if (image_url[0] == '/') && (image_url[1] == '/') {
		new_url = "http:" + new_url
	}

	if strings.Contains(image_url, "s3.amazonaws.com/affinityis") {
		new_url = strings.Replace(new_url, "s3.amazonaws.com/affinityis", "cdn.influencemobile.com", -1)
	} else if strings.Contains(image_url, "affinityis.s3.amazonaws.com") {
		new_url = strings.Replace(new_url, "affinityis.s3.amazonaws.com", "cdn.influencemobile.com", -1)

	}
	return new_url
}

func createBonus(playerId int, bonusTemplate *BonusTemplate, messageTemplate *MessageTemplate) (map[string]interface{}, error) {
	isBonusAvailable, err := pgDb.BonusAvailable(playerId, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	if isBonusAvailable {
		return nil, fmt.Errorf("bonus already available for player id - ", playerId)
	}

	if !bonusTemplate.ExpiresIn.Valid {
		bonusTemplate.ExpiresIn.Int32 = 1
	}

	if ok := strings.ToUpper(bonusTemplate.ExpiresInUnits.String); ok == "NULL" || ok == "" {
		bonusTemplate.ExpiresInUnits.String = "minutes"
	}

	bonusPayload := make(map[string]interface{})
	currentTime := time.Now().UTC()
	bonusPayload["bonus_guid"] = assignGuid("Bonus")
	bonusPayload["awarded_points"] = 0
	bonusPayload["bonus_created_at"] = currentTime.Format(time.RFC3339)
	bonusPayload["bonus_updated_at"] = currentTime.Format(time.RFC3339)
	bonusPayload["bonus_template_id"] = bonusTemplate.Id
	bonusPayload["bonus_expires_at"] = timeFormatWithUnit(bonusTemplate.ExpiresIn.Int32, bonusTemplate.ExpiresInUnits.String)
	bonusPayload["bonus_expires_in"] = bonusTemplate.ExpiresIn.Int32
	bonusPayload["bonus_awarded_at"] = nil //NULL
	bonusPayload["bonus_reward_title"] = messageTemplate.Title
	bonusPayload["bonus_completed_at"] = nil //NULL
	bonusPayload["bonus_event_name"] = bonusTemplate.EventName.String
	bonusPayload["bonus_event_count"] = bonusTemplate.EventCount.Int32
	bonusPayload["bonus_event_counter"] = nil // not exists in bonus_template table
	bonusPayload["bonus_type"] = bonusTemplate.BonusType.Int32

	// if bonus_template prorated is true then calculate the bonus_points
	if bonusTemplate.Prorated.Valid && bonusTemplate.Prorated.Bool {
		if bonusTemplate.Points.Valid && bonusTemplate.EventCount.Valid && bonusTemplate.EventCount.Int32 != 0 {
			bonusPayload["bonus_points"] = (bonusTemplate.Points.Int32 / bonusTemplate.EventCount.Int32)
		} else {
			bonusPayload["bonus_points"] = 0
		}
	} else {
		bonusPayload["bonus_points"] = bonusTemplate.Points.Int32
	}
	return bonusPayload, nil
}




func assignGuid(guidType string) string {
	uuidNew := uuid.New()
	switch guidType {
	case "PushMessage":
		return fmt.Sprintf("mpm-%s", uuidNew)
	case "Bonus":
		return fmt.Sprintf("mb-%s", uuidNew)
	case "Message":
		return fmt.Sprintf("mm-%s", uuidNew)
	default:
		return fmt.Sprintf("%s", uuidNew)
	}
}




func timeFormatWithUnit(expires_in int32, unit string) time.Time {
	currentTime := time.Now().UTC()
	switch strings.ToLower(unit) {
	case "minutes":
		return currentTime.Add(time.Minute * time.Duration(expires_in))
	case "hours":
		return currentTime.Add(time.Hour * time.Duration(expires_in))
	case "days":
		return currentTime.Add(time.Hour * time.Duration(expires_in*24))
	default:
		return currentTime
	}
}

func attentionGetter() []string {
	attentions := []string{
		"Limited 😀",
		"Jaw-dropping 😮",
		"Shocking 😱",
		"We must love you 😘",
		"We must ❤️ you!",
		"Amazing 🆕",
		"Exclusive 🎯",
		"Don’t Wait ⏰",
		"Look 👀",
		"Sweet 🍭",
		//   "Amazing \u{1F923}",
		"❤️  Incredible",
		"😃  Bonus",
		"⏳ Ends Soon",
		"Hot  🔥",
		"🌶  Hot",
		"Huge 🆒",
		"🆓 While it lasts",
	}

	return attentions
}

func onboarding(playerId string) (bool, error) {
	validPlayerId, err := getValidPlayerId(playerId)
	if err != nil {
		return false, err
	}

	result, err := QetAllItems(dynamoSvc, "players_"+env.Environment, "id", *validPlayerId)
	if err != nil {
		return false, err
	}

	players := []map[string]interface{}{}
	for _, r := range *result {
		for _, v := range r.Items {
			player := make(map[string]interface{})
			err = dynamodbattribute.UnmarshalMap(v, &player)
			if err != nil {
				return false, err
			}
			players = append(players, player)
		}
	}
	if len(players) > 0 {
		if players[0]["onboarding_completed_at"] != nil {
			return true, nil
		}
	}
	return false, nil
}


func getMemberId(playerId int) int {
	result, err := QetAllItems(dynamoSvc, "members_"+env.Environment, "player_id", playerId)
	if err != nil {
		return false, err
	}

	members := []map[string]interface{}{}
	for _, r := range *result {
		for _, v := range r.Items {
			member := make(map[string]interface{})
			err = dynamodbattribute.UnmarshalMap(v, &member)
			if err != nil {
				return false, err
			}
			members = append(members, member)
		}
	}

	if len(members) > 0 {
		return members[len(members)-1]["member_id"].(int)
	}

	return -1
}
