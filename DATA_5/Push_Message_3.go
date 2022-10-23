package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	//Simple JSON object which we will parse
	empArray := `{
		"push_messages":[
        {

            "status":"delivered",
            "provider_message_id":"12312",
            "push_provider": "aws::sn",
            "created_at":"01-01-2022",
            "updated_at":"01-01-2022",
            "guid":"mb-b23cf415-5a5d-4e90-8def-a1c9c39ec246"
         },
         {

            "status":"delivered",
            "provider_message_id":"12312",
            "push_provider": "aws::sn",
            "created_at":"01-01-2022",
            "updated_at":"01-01-2022",
            "guid":"mb-b23cf415-5a5d-4e90-8def-a1c9c39ec246"
        }
    ]
}`

	// Declared an empty interface of type Array
	var results map[string]interface{}

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal([]byte(empArray), &results)
   

	push_message := results["push_messages"].([]interface{})

	
	
	
	for _,item:=range push_message {
		fmt.Printf("%v\n", item.(map[string]interface{})["created_at"])
		fmt.Printf("%v\n", item.(map[string]interface{})["status"])
	}
    


}
