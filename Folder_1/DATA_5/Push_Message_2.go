package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	//Simple JSON object which we will parse
	empArray := `
    [
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
    ]`

	// Declared an empty interface of type Array
	var results []map[string]interface{}

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal([]byte(empArray), &results)
    //push_messages := results["push_messages"].([]map[string]interface{})

	for key, result := range results{
		
		fmt.Println("Reading Value for Key :", key)
        fmt.Println("\n")
		//Reading each value by its key
		fmt.Println("\nstatus :", result["status"],
			"\nName :", result["name"],
			"\nprovider_message_id :", result["provider_message_id"],
			"\npush_provider :", result["push_provider"])
		
	}
}
