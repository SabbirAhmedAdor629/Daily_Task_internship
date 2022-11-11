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