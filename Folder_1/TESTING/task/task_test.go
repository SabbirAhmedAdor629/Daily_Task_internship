package task

import (
	"reflect"
	"testing"
)
func Test_createSNSMessage(t *testing.T) {
	type args struct {
		sqsPushPayload    SQSPushLambdaPayload
		pushMessagesIndex int
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "#1: Invalid push messages index",
			args: args{
				sqsPushPayload:    SQSPushLambdaPayload{},
				pushMessagesIndex: -1,
			},

			want: nil,
		},
		{
			name: "#2: Empty push messages array",
			args: args{
				sqsPushPayload:    SQSPushLambdaPayload{},
				pushMessagesIndex: 0,
			},

			want: nil,
		},
		{
			name: "#3: Empty push messages array & invalid index",
			args: args{
				sqsPushPayload:    SQSPushLambdaPayload{},
				pushMessagesIndex: -1,
			},

			want: nil,
		},
		{
			name: "4#: valid sqs payload ",
			args: args{
				sqsPushPayload: SQSPushLambdaPayload{
					Origin:        "engage",
					Operation:     "reward_point",
					TransactionId: "123123-123123-sadfasfd-asdfsadf",
					SubmittedTs:   "2022-04-20T20:46:31Z",
					Data: Message{
						SkipInbox:                      true,
						GUID:                           "mm-933cf415-5a5d-4e09-8def-a1c9c39ec246",
						PlayerID:                       20802498,
						MessageTemplateID:              "2628",
						Badge:                          true,
						BgImage:                        "",
						CallToAction:                   "Try a NEW game now!! 🎲🍭🎰",
						CallToActionEvent:              "",
						CallToActionPoints:             "0",
						CompletedAt:                    "2022-07-28T16:28:48-07:00",
						CompleteOnCallToAction:         "FALSE",
						ComposqsPushPayloadonentParams: "",
						CreatedAt:                      "2022-07-28T16:28:48-07:00",
						UpdatedAt:                      "",
						ViewedAt:                       "",
						CampaignId:                     "",
						CreatedDate:                    "2022-07-28",
						CreatedHour:                    "",
						PushProvider:                   "",
						InAppPush:                      false,
						CausedUninstall:                false,
						RawPush:                        false,
						OutboundNumber:                 "",
						PushNotification:               false,
						PushResponse:                   "",
						PushDelivered:                  false,
						Goto:                           "home",
						Image:                          "http://cdn.influencemobile.com/message_templates/images/000/002/628/header/header.jpeg?1657225678",
						MemberID:                       "0a9b82cb-31f8-6452-a6bf-dc44a09342e7",
						NextMessage:                    "",
						PushMessage:                    "Dogs like water 💦 You like BONUS time!! 🎉",
						S3ProgramImageExtension:        "",
						S3ProgramImageName:             "",
						Title:                          "Bonus time RIGHT now!!! ⏱",
						ViewPoints:                     "0",
						Body:                           "",
						TemplateBody:                   "",
						BonusGUID:                      "",
						AwardedPoints:                  "500",
						BonusCreatedAt:                 "2022-07-28T16:28:48-07:00",
						BonusUpdatedAt:                 "2022-07-28T16:28:48-07:00",
						BonusTemplateID:                "",
						BonusExpiresAt:                 "2022-07-28T18:28:48-07:00",
						BonusExpiresIn:                 "43200",
						BonusAwardedAt:                 "2022-07-28T18:28:48-07:00",
						BonusRewardTitle:               "Bonus time RIGHT now!!! ⏱",
						BonusCompletedAt:               "",
						BonusEventName:                 "engage_install_server",
						BonusEventCount:                "0",
						BonusEventCounter:              "0",
						BonusType:                      "prize",
						PushMessages:                   []PushMessages{{Status: "delivered", ProviderMessageID: "12312", PushProvider: "aws: :sn", CreatedAt: "01-01-2022", UpdatedAt: "01-01-2022", MessageUUID: "ff9e7777-d4c6-4dcd-9422-40b4f6257370", GUID: "mb-b23cf415-5a5d-4e90-8def-a1c9c39ec246", MessagePK: "mm-933cf415-5a5d-4e09-8def-a1c9c39ec246", DeviceId: "123213213", PushRegistrationId: "13123232", AwsMessageId: "21321321", AwsArnStatus: "Success"}, {Status: "delivered", ProviderMessageID: "4321", PushProvider: "aws::sn:32", CreatedAt: "02-01-2022", UpdatedAt: "02-01-2022", MessageUUID: "ef9e7777-d4c6-4dcd-9422-40b4f6257360", GUID: "mb-423cf415-5a5d-4e90-8def-a1c9c39ec248", MessagePK: "mm-933cf415-5a5d-4e09-8def-a1c9c39ec246", DeviceId: "43432132", PushRegistrationId: "123213213", AwsMessageId: "21123213", AwsArnStatus: "Success"}},
					},
				},
				pushMessagesIndex: 0,
			},

			want: map[string]interface{}{
				"GCM": map[string]interface{}{
					"notification": map[string]interface{}{
						"title": "Bonus time RIGHT now!!! ⏱",
						// others field will be there
					},
					// others field will be there
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createSNSMessage(tt.args.sqsPushPayload, tt.args.pushMessagesIndex); !reflect.DeepEqual(got, tt.want) {
				if !(got != nil && tt.want != nil) {
					t.Errorf("createSNSMessage() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}