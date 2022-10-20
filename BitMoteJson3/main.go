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
	var event SNS_Event

	jsondata := `{
		"skip_inbox" : true,
		"push_message" :[
			{
				"member_guid" : "4a1147ac-494b-11ed-b878-0242ac120002",
				"device_id" : "3478-34534",
				"push_registration_id" : "354433",
				"aws_message_id" : "1231",
				"aws_arn_status" : "Failed"
			},
			{
				"member_guid" : "4a1147ac-494b-11ed-b878-0242ac120002",
				"device_id" : "238923-3432",
				"push_registration_id" : "8989234",
				"aws_message_id" : "178232",
				"aws_arn_status" : "Success"
			}
		],
		"message": {
			"player_id" : 12001,
			"message-template_id" : 24,
			"data" : "Push message notification",
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
			fmt.Println("\n")
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


