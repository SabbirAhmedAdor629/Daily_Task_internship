package main

import (
	"encoding/json"
	"fmt"
	//"os"
)


type Push_Messages struct {
	Member_Guid	string
	Device_Id   string
	Push_Registration_Id  string
	Aws_Message_Id  string
	Aws_Arn_Status	string
}

type Messages struct{
	Player_Id  int
	Message_Template_Id  int
	Data  string
	Badge  bool
}

type SNS_Event struct{
	Skip_Inbox bool
	Push_Message []Push_Messages
	Message Messages
}

func main(){

	map := map[string]interface{
		completed_at : 2022-07-28T16:28:48-07:00
		completed_at_2 : 2022-07-28T16:28:48-07:00

		push_message : Dogs like water 💦 You like BONUS time!! 🎉
		push_message_2 : Dogs like water 💦 You like BONUS time!! 🎉

		call_to_action : 0
		call_to_action_2 : 0
		
		push_messages : [map[created_at:01-01-2022 guid:mb-b23cf415-5a5d-4e90-8def-a1c9c39ec246 provider_message_id:12312 push_provider:aws::sn status:not-delivered updated_at:01-01-2022] map[created_at:01-01-2022 guid:mb-b23cf415-5a5d-4e90-8def-a1c9c39ec246 provider_message_id:12312 push_provider:aws::sn status:delivered updated_at:01-01-2022]]
		
		member_id : 0a9b82cb-31f8-6452-a6bf-dc44a09342e7
		data : Push message notification
		bonus_guid : mb-b23cf415-5a5d-4e90-8def-a1c9c39ec246
	}


	var event SNS_Event

	jsondata := `{
		"skip_inbox" : true,
		"push_message" :[
			{
				"member_guid" : "4a1147ac-494b-11ed-b878-0242ac120002",
				"device_id" : "sabbir2",
				"push_registration_id" : "sabbir3",
				"aws_message_id" : "sabbir4",
				"aws_arn_status" : "sabbir5"
			}
		],
		"message": {
			"player_id" : 12,
			"message-template_id" : 24,
			"data" : "sabbir9",
			"badge" : true
		}
	   
	}`

	err := json.Unmarshal([]byte(jsondata), &event)

	if err == nil{
		fmt.Println("Skip_Inbox: ",event.Skip_Inbox)
		fmt.Println(" ")
		
		for _,obj := range event.Push_Message{
			fmt.Println("Member_Guid : ", obj.Member_Guid)
			fmt.Println("Device_Id : ",obj.Device_Id)
			fmt.Println("Push_Registration_Id : ",obj.Push_Registration_Id)
			fmt.Println("Aws_Message_Id : ",obj.Aws_Message_Id)
			fmt.Println("Aws_Arn_Status : ",obj.Aws_Arn_Status)
		}
		fmt.Println(" ")
		
		fmt.Println("Player_Id : ",event.Message.Player_Id)
		fmt.Println("Message_Template_Id : ",event.Message.Message_Template_Id)
		fmt.Println("Data : ",event.Message.Data)
		fmt.Println("Badge : ",event.Message.Badge)

	}else{
		fmt.Println(err)
	}

	
}



// func main() {
// 	reader, _ := os.Open("data.json")
// 	decoder := json.NewDecoder(reader)

// 	event := &SNS_Event{}
// 	decoder.Decode(event)

// 	fmt.Println(event)
// }