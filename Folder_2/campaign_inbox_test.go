package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	config "influencemobile.com/parameter_store"
)

type (
	mockDynamoClient struct{ dynamodbiface.DynamoDBAPI }
)

func setEnvVars(es map[TypeEnvVar]string) {
	for k, v := range es {
		os.Setenv(EnvVars[k], v)
	}
}

func (m mockDynamoClient) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	result := dynamodb.PutItemOutput{}

	if input.Item == nil {
		return &result, errors.New("Missing required field PutItemInput.Item")
	}

	if input.TableName == nil || *input.TableName == "" {
		return &result, errors.New("Missing required field CreateTableInput.TableName")
	}

	return &result, nil
}

func TestPutItemDynamoTable(t *testing.T) {

	type args struct {
		tableName string
		dynamoSvc dynamodbiface.DynamoDBAPI
	}
	tests := []struct {
		name          string
		args          args
		dynamoItemMap map[string]interface{}
		tableName     string
		wantErr       bool
	}{
		{
			name: "#1: Valid map data",
			args: args{
				dynamoSvc: &mockDynamoClient{},
				tableName: "CampaignInboxTest",
			},
			dynamoItemMap: map[string]interface{}{
				"origin":                     "engage",
				"operation":                  "reward_point",
				"transaction_id":             "123123-123123-sadfasfd-asdfsadf",
				"priority":                   "high",
				"submitted_ts":               "2022-04-20T20:46:31Z",
				"guid":                       "mm-933cf415-5a5d-4e09-8def-a1c9c39ec246",
				"player_id":                  20802498,
				"message_template_id":        2628,
				"badge":                      true,
				"bg_image":                   "",
				"call_to_action":             "Try a NEW game now!! 🎲🍭🎰",
				"call_to_action_event":       "",
				"call_to_action_points":      0,
				"completed_at":               "2022-07-28T16:28:48-07:00",
				"complete_on_call_to_action": "FALSE",
				"component":                  "",
				"component_params":           "",
				"created_at":                 "2022-07-28T16:28:48-07:00",
				"created_date":               "2022-07-28",
				"goto":                       "home",
				"image":                      "http://cdn.influencemobile.com/message_templates/images/000/002/628/header/header.jpeg?1657225678",
				"member_id":                  "0a9b82cb-31f8-6452-a6bf-dc44a09342e7",
				"next_message":               "",
				"push_message":               "Dogs like water 💦 You like BONUS time!! 🎉",
				"s3_program_image_extension": "",
				"s3_program_image_name":      "",
				"title":                      "Bonus time RIGHT now!!! ⏱",
				"view_points":                0,
				"push_message_body":          "Just like pups like water 💦, we heard you like BONUSES!  Here ya go!! 🥳 500 points to try a NEW game in the next 2 hours!  🎉",
				"template_body":              "sample template body message",
			},
			wantErr: false,
		},
		{
			name: "#2: Empty map data",
			args: args{
				dynamoSvc: &mockDynamoClient{},
				tableName: "CampaignInboxTest",
			},
			dynamoItemMap: map[string]interface{}{},
			wantErr:       true,
		},
		{
			name: "#3: Nil map data",
			args: args{
				dynamoSvc: &mockDynamoClient{},
				tableName: "CampaignInboxTest",
			},
			dynamoItemMap: nil,
			wantErr:       true,
		},
		{
			name: "#3: Empty table name",
			args: args{
				dynamoSvc: &mockDynamoClient{},
				tableName: "",
			},
			dynamoItemMap: map[string]interface{}{
				"origin":                     "engage",
				"operation":                  "reward_point",
				"transaction_id":             "123123-123123-sadfasfd-asdfsadf",
				"priority":                   "high",
				"submitted_ts":               "2022-04-20T20:46:31Z",
				"guid":                       "mm-933cf415-5a5d-4e09-8def-a1c9c39ec246",
				"player_id":                  20802498,
				"message_template_id":        2628,
				"badge":                      true,
				"bg_image":                   "",
				"call_to_action":             "Try a NEW game now!! 🎲🍭🎰",
				"call_to_action_event":       "",
				"call_to_action_points":      0,
				"completed_at":               "2022-07-28T16:28:48-07:00",
				"complete_on_call_to_action": "FALSE",
				"component":                  "",
				"component_params":           "",
				"created_at":                 "2022-07-28T16:28:48-07:00",
				"created_date":               "2022-07-28",
				"goto":                       "home",
				"image":                      "http://cdn.influencemobile.com/message_templates/images/000/002/628/header/header.jpeg?1657225678",
				"member_id":                  "0a9b82cb-31f8-6452-a6bf-dc44a09342e7",
				"next_message":               "",
				"push_message":               "Dogs like water 💦 You like BONUS time!! 🎉",
				"s3_program_image_extension": "",
				"s3_program_image_name":      "",
				"title":                      "Bonus time RIGHT now!!! ⏱",
				"view_points":                0,
				"push_message_body":          "Just like pups like water 💦, we heard you like BONUSES!  Here ya go!! 🥳 500 points to try a NEW game in the next 2 hours!  🎉",
				"template_body":              "sample template body message",
			},
			wantErr: true,
		},
	}

	for _, testData := range tests {
		t.Run("PutItem DynamoDB Table", func(t *testing.T) {
			if err := putItemDynamoTable(testData.args.dynamoSvc, testData.args.tableName, testData.dynamoItemMap); (err != nil) != testData.wantErr {
				t.Errorf("putItem error = %v, wantErr %v", err, testData.wantErr)
			}
		})
	}

}

func TestHandleRequest(t *testing.T) {

	parameterStore = &config.MockKeyClient{
		FileBytes: []byte(`{"keys":[{"key":"/honeybadger/messaging/api_key","value":{"Name":"/honeybadger/messaging/api_key","Value":"api_key"}}]}`),
	}

	type args struct {
		ctx      context.Context
		sqsEvent *events.SQSEvent
	}

	setEnvVars(map[TypeEnvVar]string{
		ENV_APPLICATION_NAME: "campaign_inbox",
		ENV_LOG_LEVEL:        "debug",
		ENV_AWS_REGION:       "us-east-1",
		ENV_PARAMETERS_STORE: "/tmp/parameters.txt",
		ENV_DB_TABLE_NAME:    "TestCampaignInbox",
	})

	tests := []struct {
		name                     string
		args                     args
		wantPartialBatchResponse map[string]interface{}
		wantErr                  bool
	}{
		{
			name: "#1: Valid SQS Json Message",
			args: args{
				ctx: context.TODO(),
				sqsEvent: &events.SQSEvent{
					Records: []events.SQSMessage{
						events.SQSMessage{
							Body: `{
								"origin": "engage",
								"operation": "reward_point",
								"transaction_id": "123123-123123-sadfasfd-asdfsadf",
								"priority": "high",
								"submitted_ts": "2022-04-20T20:46:31Z",
								"data": {
								  "skip_inbox": true,
								  "guid": "mm-933cf415-5a5d-4e09-8def-a1c9c39ec246",
								  "player_id": 20802498,
								  "message_template_id": 2628,
								  "badge": true,
								  "bg_image": "",
								  "call_to_action": "Try a NEW game now!! 🎲🍭🎰",
								  "call_to_action_event": "",
								  "call_to_action_points": 0,
								  "completed_at": "2022-07-28T16:28:48-07:00",
								  "complete_on_call_to_action": "FALSE",
								  "component": "",
								  "component_params": "",
								  "created_at": "2022-07-28T16:28:48-07:00",
								  "created_date": "2022-07-28",
								  "goto": "home",
								  "image": "http://cdn.influencemobile.com/message_templates/images/000/002/628/header/header.jpeg?1657225678",
								  "member_id": "0a9b82cb-31f8-6452-a6bf-dc44a09342e7",
								  "next_message": "",
								  "push_message": "Dogs like water 💦 You like BONUS time!! 🎉",
								  "s3_program_image_extension": "",
								  "s3_program_image_name": "",
								  "title": "Bonus time RIGHT now!!! ⏱",
								  "view_points": 0,
								  "push_message_body": "Just like pups like water 💦, we heard you like BONUSES!  Here ya go!! 🥳 500 points to try a NEW game in the next 2 hours!  🎉",
								  "template_body":"sample template body message",
								  "push_messages": [
									{
									  "status": "delivered",
									  "provider_message_id": "12312",
									  "push_provider": "aws::sn",
									  "created_at": "01-01-2022",
									  "updated_at": "01-01-2022",
									  "message_uuid": "ff9e7777-d4c6-4dcd-9422-40b4f6257370",
									  "guid": "mb-b23cf415-5a5d-4e90-8def-a1c9c39ec246",
									  "message_pk": "mm-933cf415-5a5d-4e09-8def-a1c9c39ec246",
									  "device_id": "123213213",
									  "push_registration_id": "13123232",
									  "aws_message_id": "21321321",
									  "aws_arn_status": "Success"
									},
									{
									  "status": "delivered",
									  "provider_message_id": "4321",
									  "push_provider": "aws::sn:32",
									  "created_at": "02-01-2022",
									  "updated_at": "02-01-2022",
									  "message_uuid": "ef9e7777-d4c6-4dcd-9422-40b4f6257360",
									  "guid": "mb-423cf415-5a5d-4e90-8def-a1c9c39ec248",
									  "message_pk": "mm-933cf415-5a5d-4e09-8def-a1c9c39ec246",
									  "device_id": "43432132",
									  "push_registration_id": "123213213",
									  "aws_message_id": "21123213",
									  "aws_arn_status": "Success"
									}
								  ],
								  "bonus_guid": "mb-b23cf415-5a5d-4e90-8def-a1c9c39ec246",
								  "awarded_points": 500,
								  "bonus_created_at": "2022-07-28T16:28:48-07:00",
								  "bonus_updated_at": "2022-07-28T16:28:48-07:00",
								  "bonus_template_id": 210,
								  "bonus_expires_at": "2022-07-28T18:28:48-07:00",
								  "bonus_expires_in": 43200,
								  "bonus_awarded_at": "2022-07-28T18:28:48-07:00",
								  "bonus_reward_title": "Bonus time RIGHT now!!! ⏱",
								  "bonus_completed_at": "",
								  "bonus_event_name": "engage_install_server",
								  "bonus_event_count": 0,
								  "bonus_event_counter": 0,
								  "bonus_type": "prize"
								}
							  }`,
							MessageId: "34f2d472-e9c3-49f8-914c-h52453452rr67",
						},
					},
				},
			},
			wantErr:                  false,
			wantPartialBatchResponse: nil,
		},
		{
			name: "#2: Invalid SQS Message [data field missing]",
			args: args{
				ctx: context.TODO(),
				sqsEvent: &events.SQSEvent{
					Records: []events.SQSMessage{
						events.SQSMessage{
							Body: `{
								  "skip_inbox": true,
								  "guid": "mm-933cf415-5a5d-4e09-8def-a1c9c39ec246",
								  "player_id": 20802498,
								  "message_template_id": 2628,
								  "badge": true,
								  "bg_image": "",
								  "call_to_action": "Try a NEW game now!! 🎲🍭🎰",
								  "call_to_action_event": "",
								  "call_to_action_points": 0,
								  "completed_at": "2022-07-28T16:28:48-07:00",
								  "complete_on_call_to_action": "FALSE",
								  "component": "",
								  "component_params": "",
								  "created_at": "2022-07-28T16:28:48-07:00",
								  "created_date": "2022-07-28",
								  "goto": "home",
								  "image": "http://cdn.influencemobile.com/message_templates/images/000/002/628/header/header.jpeg?1657225678",
								  "member_id": "0a9b82cb-31f8-6452-a6bf-dc44a09342e7",
								  "next_message": "",
								  "push_message": "Dogs like water 💦 You like BONUS time!! 🎉",
								  "s3_program_image_extension": "",
								  "s3_program_image_name": "",
								  "title": "Bonus time RIGHT now!!! ⏱",
								  "view_points": 0,
								  "push_message_body": "Just like pups like water 💦, we heard you like BONUSES!  Here ya go!! 🥳 500 points to try a NEW game in the next 2 hours!  🎉",
								  "template_body":"sample template body message",
								  "push_messages": [
									{
									  "status": "delivered",
									  "provider_message_id": "12312",
									  "push_provider": "aws::sn",
									  "created_at": "01-01-2022",
									  "updated_at": "01-01-2022",
									  "message_uuid": "ff9e7777-d4c6-4dcd-9422-40b4f6257370",
									  "guid": "mb-b23cf415-5a5d-4e90-8def-a1c9c39ec246",
									  "message_pk": "mm-933cf415-5a5d-4e09-8def-a1c9c39ec246",
									  "device_id": "123213213",
									  "push_registration_id": "13123232",
									  "aws_message_id": "21321321",
									  "aws_arn_status": "Success"
									},
									{
									  "status": "delivered",
									  "provider_message_id": "4321",
									  "push_provider": "aws::sn:32",
									  "created_at": "02-01-2022",
									  "updated_at": "02-01-2022",
									  "message_uuid": "ef9e7777-d4c6-4dcd-9422-40b4f6257360",
									  "guid": "mb-423cf415-5a5d-4e90-8def-a1c9c39ec248",
									  "message_pk": "mm-933cf415-5a5d-4e09-8def-a1c9c39ec246",
									  "device_id": "43432132",
									  "push_registration_id": "123213213",
									  "aws_message_id": "21123213",
									  "aws_arn_status": "Success"
									}
								  ],
								  "bonus_guid": "mb-b23cf415-5a5d-4e90-8def-a1c9c39ec246",
								  "awarded_points": 500,
								  "bonus_created_at": "2022-07-28T16:28:48-07:00",
								  "bonus_updated_at": "2022-07-28T16:28:48-07:00",
								  "bonus_template_id": 210,
								  "bonus_expires_at": "2022-07-28T18:28:48-07:00",
								  "bonus_expires_in": 43200,
								  "bonus_awarded_at": ""GMT 07pm", "",
								  "bonus_event_name": "engage_install_server",
								  "bonus_event_count": 0,
								  "bonus_event_counter": 0,
								  "bonus_type": "prize"
								}`,
							MessageId: "34f2d472-e9c3-49f8-914c-h52453452rr67",
						},
					},
				},
			},
			wantErr: true,
			wantPartialBatchResponse: map[string]interface{}{
				"batchItemFailures": []map[string]string{
					map[string]string{
						"itemIdentifier": "e4f2d472-e9c3-49f8-914c-h32453450278",
					},
				},
			},
		},
		{
			name: "#3: Invalid Json Message",
			args: args{
				ctx: context.TODO(),
				sqsEvent: &events.SQSEvent{
					Records: []events.SQSMessage{
						events.SQSMessage{
							Body: `{
								  "skip_inbox": true,
								  "guid": "mm-933cf415-5a5d-4e09-8def-a1c9c39ec246",
								  "player_id": 20802498,
								  "message_template_id": 2628
								`,
							MessageId: "34f2d472-e9c3-49f8-914c-h52453452rr67",
						},
					},
				},
			},
			wantErr: true,
			wantPartialBatchResponse: map[string]interface{}{
				"batchItemFailures": []map[string]string{
					map[string]string{
						"itemIdentifier": "e4f2d472-e9c3-49f8-914c-h32453450278",
					},
				},
			},
		},
		{
			name: "#4: Invalid SQS Message [player_id not available]",
			args: args{
				ctx: context.TODO(),
				sqsEvent: &events.SQSEvent{
					Records: []events.SQSMessage{
						events.SQSMessage{
							Body: `{
								"origin": "engage",
								"operation": "reward_point",
								"transaction_id": "123123-123123-sadfasfd-asdfsadf",
								"priority": "high",
								"submitted_ts": "2022-04-20T20:46:31Z",
								"data": {
								  "skip_inbox": true,
								  "guid": "mm-933cf415-5a5d-4e09-8def-a1c9c39ec246",
								  "message_template_id": 2628,
								  "badge": true,
								  "bg_image": "",
								  "call_to_action": "Try a NEW game now!! 🎲🍭🎰",
								  "call_to_action_event": "",
								  "call_to_action_points": 0,
								  "completed_at": "2022-07-28T16:28:48-07:00",
								  "complete_on_call_to_action": "FALSE",
								  "component": "",
								  "component_params": "",
								  "created_at": "2022-07-28T16:28:48-07:00",
								  "created_date": "2022-07-28",
								  "goto": "home",
								  "image": "http://cdn.influencemobile.com/message_templates/images/000/002/628/header/header.jpeg?1657225678",
								  "member_id": "0a9b82cb-31f8-6452-a6bf-dc44a09342e7",
								  "next_message": "",
								  "push_message": "Dogs like water 💦 You like BONUS time!! 🎉",
								  "s3_program_image_extension": "",
								  "s3_program_image_name": "",
								  "title": "Bonus time RIGHT now!!! ⏱",
								  "view_points": 0,
								  "push_message_body": "Just like pups like water 💦, we heard you like BONUSES!  Here ya go!! 🥳 500 points to try a NEW game in the next 2 hours!  🎉",
								  "template_body":"sample template body message",
								  "push_messages": [
									{
									  "status": "delivered",
									  "provider_message_id": "12312",
									  "push_provider": "aws::sn",
									  "created_at": "01-01-2022",
									  "updated_at": "01-01-2022",
									  "message_uuid": "ff9e7777-d4c6-4dcd-9422-40b4f6257370",
									  "guid": "mb-b23cf415-5a5d-4e90-8def-a1c9c39ec246",
									  "message_pk": "mm-933cf415-5a5d-4e09-8def-a1c9c39ec246",
									  "device_id": "123213213",
									  "push_registration_id": "13123232",
									  "aws_message_id": "21321321",
									  "aws_arn_status": "Success"
									},
									{
									  "status": "delivered",
									  "provider_message_id": "4321",
									  "push_provider": "aws::sn:32",
									  "created_at": "02-01-2022",
									  "updated_at": "02-01-2022",
									  "message_uuid": "ef9e7777-d4c6-4dcd-9422-40b4f6257360",
									  "guid": "mb-423cf415-5a5d-4e90-8def-a1c9c39ec248",
									  "message_pk": "mm-933cf415-5a5d-4e09-8def-a1c9c39ec246",
									  "device_id": "43432132",
									  "push_registration_id": "123213213",
									  "aws_message_id": "21123213",
									  "aws_arn_status": "Success"
									}
								  ],
								  "bonus_guid": "mb-b23cf415-5a5d-4e90-8def-a1c9c39ec246",
								  "awarded_points": 500,
								  "bonus_created_at": "2022-07-28T16:28:48-07:00",
								  "bonus_updated_at": "2022-07-28T16:28:48-07:00",
								  "bonus_template_id": 210,
								  "bonus_expires_at": "2022-07-28T18:28:48-07:00",
								  "bonus_expires_in": 43200,
								  "bonus_awarded_at": "2022-07-28T18:28:48-07:00",
								  "bonus_reward_title": "Bonus time RIGHT now!!! ⏱",
								  "bonus_completed_at": "",
								  "bonus_event_name": "engage_install_server",
								  "bonus_event_count": 0,
								  "bonus_event_counter": 0,
								  "bonus_type": "prize"
								}
							  }`,
							MessageId: "34f2d472-e9c3-49f8-914c-h52453452rr67",
						},
					},
				},
			},
			wantErr: true,
			wantPartialBatchResponse: map[string]interface{}{
				"batchItemFailures": []map[string]string{
					map[string]string{
						"itemIdentifier": "e4f2d472-e9c3-49f8-914c-h32453450278",
					},
				},
			},
		},
		{
			name: "#5: Invalid SQS Message [ guid not available]",
			args: args{
				ctx: context.TODO(),
				sqsEvent: &events.SQSEvent{
					Records: []events.SQSMessage{
						events.SQSMessage{
							Body: `{
								"origin": "engage",
								"operation": "reward_point",
								"transaction_id": "123123-123123-sadfasfd-asdfsadf",
								"priority": "high",
								"submitted_ts": "2022-04-20T20:46:31Z",
								"data": {
								  "skip_inbox": true,
								  "player_id": 20802498,
								  "message_template_id": 2628,
								  "badge": true,
								  "bg_image": "",
								  "call_to_action": "Try a NEW game now!! 🎲🍭🎰",
								  "call_to_action_event": "",
								  "call_to_action_points": 0,
								  "completed_at": "2022-07-28T16:28:48-07:00",
								  "complete_on_call_to_action": "FALSE",
								  "component": "",
								  "component_params": "",
								  "created_at": "2022-07-28T16:28:48-07:00",
								  "created_date": "2022-07-28",
								  "goto": "home",
								  "image": "http://cdn.influencemobile.com/message_templates/images/000/002/628/header/header.jpeg?1657225678",
								  "member_id": "0a9b82cb-31f8-6452-a6bf-dc44a09342e7",
								  "next_message": "",
								  "push_message": "Dogs like water 💦 You like BONUS time!! 🎉",
								  "s3_program_image_extension": "",
								  "s3_program_image_name": "",
								  "title": "Bonus time RIGHT now!!! ⏱",
								  "view_points": 0,
								  "push_message_body": "Just like pups like water 💦, we heard you like BONUSES!  Here ya go!! 🥳 500 points to try a NEW game in the next 2 hours!  🎉",
								  "template_body":"sample template body message",
								  "push_messages": [
									{
									  "status": "delivered",
									  "provider_message_id": "12312",
									  "push_provider": "aws::sn",
									  "created_at": "01-01-2022",
									  "updated_at": "01-01-2022",
									  "message_uuid": "ff9e7777-d4c6-4dcd-9422-40b4f6257370",
									  "guid": "mb-b23cf415-5a5d-4e90-8def-a1c9c39ec246",
									  "message_pk": "mb-b23cf415-5a5d-4e90-8def-a1c9c39ec655",
									  "device_id": "123213213",
									  "push_registration_id": "13123232",
									  "aws_message_id": "21321321",
									  "aws_arn_status": "Success"
									},
									{
									  "status": "delivered",
									  "provider_message_id": "4321",
									  "push_provider": "aws::sn:32",
									  "created_at": "02-01-2022",
									  "updated_at": "02-01-2022",
									  "message_uuid": "ef9e7777-d4c6-4dcd-9422-40b4f6257360",
									  "guid": "mb-423cf415-5a5d-4e90-8def-a1c9c39ec248",
									  "message_pk": "mm-933cf415-5a5d-4e09-8def-a1c9c39ec246",
									  "device_id": "43432132",
									  "push_registration_id": "123213213",
									  "aws_message_id": "21123213",
									  "aws_arn_status": "Success"
									}
								  ],
								  "bonus_guid": "mb-b23cf415-5a5d-4e90-8def-a1c9c39ec246",
								  "awarded_points": 500,
								  "bonus_created_at": "2022-07-28T16:28:48-07:00",
								  "bonus_updated_at": "2022-07-28T16:28:48-07:00",
								  "bonus_template_id": 210,
								  "bonus_expires_at": "2022-07-28T18:28:48-07:00",
								  "bonus_expires_in": 43200,
								  "bonus_awarded_at": "2022-07-28T18:28:48-07:00",
								  "bonus_reward_title": "Bonus time RIGHT now!!! ⏱",
								  "bonus_completed_at": "",
								  "bonus_event_name": "engage_install_server",
								  "bonus_event_count": 0,
								  "bonus_event_counter": 0,
								  "bonus_type": "prize"
								}
							  }`,
							MessageId: "34f2d472-e9c3-49f8-914c-h52453452rr67",
						},
					},
				},
			},
			wantErr: true,
			wantPartialBatchResponse: map[string]interface{}{
				"batchItemFailures": []map[string]string{
					map[string]string{
						"itemIdentifier": "e4f2d472-e9c3-49f8-914c-h32453450278",
					},
				},
			},
		},
	}

	env = NewEnv()

	dynamoSvc = mockDynamoClient{}
	// tableName := "CampaignInboxTest"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, _ := HandleRequest(tt.args.ctx, tt.args.sqsEvent)
			if got := fmt.Sprintf("%v", tt.wantPartialBatchResponse) != fmt.Sprintf("%v", response); got != tt.wantErr {
				t.Errorf("HandleRequest() responseBatch = %v, wantBatch %v", response, tt.wantPartialBatchResponse)
			}
		})
	}
}
