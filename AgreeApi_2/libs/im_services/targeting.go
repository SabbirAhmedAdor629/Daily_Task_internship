package im_services

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	. "influencemobile.com/libs/dynamo"
)

type PlayerIdType int

func Player(svc dynamodbiface.DynamoDBAPI, dynamoTableName string, playerId PlayerIdType) ([]byte, error) {

	result, err := GetAllItems(svc, dynamoTableName, "player_id", int(playerId))
	if err != nil {
		return nil, err
	}
	players := []map[string]interface{}{}

	for _, r := range *result {
		for _, v := range r.Items {
			player := make(map[string]interface{})
			err = dynamodbattribute.UnmarshalMap(v, &player)
			if err != nil {
				return nil, err
			}
			players = append(players, player)
		}
	}
	if len(players) == 1 {
		playerJson, err := extractDynamoHashValues(players[0], getSchemaHash())
		if err != nil {
			return nil, err
		}
		playerHash := map[string]interface{}{}
		json.Unmarshal([]byte(playerJson), &playerHash)

		countDays, err := dateDiffer(players[0]["account_created_at"], time.Now())
		if err != nil {
			playerHash["install_age"] = nil
		} else {
			playerHash["install_age"] = countDays
		}

		countDays, err = dateDiffer(players[0]["session_last"], time.Now())
		if err != nil {
			playerHash["inactive_days"] = nil
		} else {
			playerHash["inactive_days"] = countDays
		}

		playerHash = profilesValidate(playerHash)
		jsonPlayerHash, err := json.Marshal(playerHash)
		if err != nil {
			return nil, err
		}
		return jsonPlayerHash, nil
	} else {
		return nil, nil
	}
}

func profilesValidate(playerHash map[string]interface{}) map[string]interface{} {
	if _, ok := playerHash["profiles"]; ok {
		profiles := playerHash["profiles"]
		playerHash["profiles"] = unique(profiles)
	} else {
		playerHash["profiles"] = []string{"All Users"}
	}
	return playerHash
}

func unique(arr interface{}) []string {
	occurred := map[string]bool{}
	result := []string{}
	for _, value := range arr.([]interface{}) {
		if occurred[value.(string)] != true {
			occurred[value.(string)] = true
			result = append(result, value.(string))
		}
	}
	return result
}

func getSchemaHash() map[string]string {
	schemaHash := map[string]string{
		"age_max":                      "Fixnum",
		"age_min":                      "Fixnum",
		"session_avg_points":           "Fixnum",
		"engage_installs":              "Fixnum",
		"daily_ongoing_usage_installs": "Fixnum",
		"time_zone_offset":             "Fixnum",
		"redemption_count":             "Fixnum",
		"redemption_points":            "Fixnum",
		"rewarded_offers_available":    "Fixnum",
		"cpi":                          "Float",
		"latitude":                     "Float",
		"longitude":                    "Float",
		"interstitial_revenue":         "Float",
		"interstitial_ctr":             "Float",
		"total_revenue":                "Float",
		"cpi_blended":                  "Float",
		"engage_revenue":               "Float",
	}

	stringFields := []string{"country", "city", "state", "zip", "gender", "session_first", "session_last", "account_created_at", "time_zone"}
	integerFields := []string{"player_id", "xp", "active_age", "install_age", "usage_days", "inactive_days", "session_count", "session_avg_daily", "session_avg_duration"}

	for _, v := range stringFields {
		schemaHash[v] = "string"
	}
	for _, v := range integerFields {
		schemaHash[v] = "integer"
	}
	return schemaHash
}

func dateDiffer(dateA interface{}, dateB time.Time) (int, error) {
	if dateA != nil && isDate(dateA.(string)) {
		t, _ := time.Parse("2006-01-02", dateA.(string))
		today := dateB
		diff := today.Sub(t)
		return int(diff.Hours() / 24), nil
	} else if dateA != nil && isTime(dateA.(string)) {
		t, _ := time.Parse(time.RFC3339, dateA.(string))
		today := dateB
		diff := today.Sub(t)
		return int(diff.Hours() / 24), nil
	} else {
		return -1, fmt.Errorf("data is not valid")
	}
}

func extractDynamoHashValues(playerData map[string]interface{}, schemaHash map[string]string) ([]byte, error) {
	playerHash := make(map[string]interface{})
	for key, value := range playerData {
		_, ok := schemaHash[key]
		if value == "" {
			playerHash[key] = nil
		} else if ok && (schemaHash[key] == "Fixnum" || schemaHash[key] == "integer") {
			switch value.(type) {
			case string:
				intValue, _ := strconv.Atoi(value.(string))
				playerHash[key] = intValue
			case float64:
				playerHash[key] = int(value.(float64))
			default:
				playerHash[key] = value
			}
		} else if ok && schemaHash[key] == "Float" {
			switch value.(type) {
			case float64:
				playerHash[key] = strconv.FormatFloat(value.(float64), 'E', -1, 64)
			default:
				playerHash[key] = nil
			}
		} else if ok && schemaHash[key] == "Date" {
			date, _ := time.Parse("2006-01-02", value.(string))
			playerHash[key] = date.Format("Mon, 02 Jan 2006")
		} else if value == "<empty_string>" {
			playerHash[key] = ""
		} else if key[len(key)-3:] == "_at" {
			if isTime(value.(string)) {
				t1, _ := time.Parse(time.RFC3339, value.(string))
				playerHash[key] = t1
			} else if isDate(value.(string)) {
				d1, _ := time.Parse("2006-01-02", value.(string))
				playerHash[key] = d1.Format("Mon, 02 Jan 2006")
			} else {
				playerHash[key] = value
			}
		} else {
			playerHash[key] = value
		}
	}

	jsonPlayerHash, err := json.Marshal(playerHash)
	if err != nil {
		return nil, err
	}
	return jsonPlayerHash, nil
}

func isTime(value string) bool {
	_, e := time.Parse(time.RFC3339, value)
	return e == nil
}

func isDate(value string) bool {
	_, e := time.Parse("2006-01-02", value)
	return e == nil
}
