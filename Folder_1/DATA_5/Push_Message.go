package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	//Simple JSON object which we will parse
	empArray := `{
        "guid":"mb-b23cf415-5a5d-4e90-8def-a1c9c39ec246",
        "push_messages":{

            "status":"delivered",
            "provider_message_id":"12312",
            "push_provider": "aws::sn",
            "created_at":"01-01-2022",
            "updated_at":"01-01-2022",
            "guid":"mb-b23cf415-5a5d-4e90-8def-a1c9c39ec246"
         },
        "status":"delivered"
}`

	// Declared an empty interface of type Array
	var results map[string]interface{}

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal([]byte(empArray), &results)
	//push_messages := results["push_messages"].([]interface{})
	
	// Use make() to create the slice for better performance
	translationKeys := make([]string, len(results))
	for key := range results {
		translationKeys = append(translationKeys, key)
	}
	

	// Printing keys
	for k := range translationKeys {
		fmt.Println(translationKeys[k])
	}

    fmt.Println("\nstatus :", results["status"])
    
    push_messages := results["push_messages"].([]map[string]interface{})

    //fmt.Println(push_messages)

    for key, result := range push_messages{
		
		fmt.Println("Reading Value for Key :", key)
        fmt.Println("\n")
		//Reading each value by its key
		fmt.Println("\nstatus :", result["status"],
			"\nprovider_message_id :", result["provider_message_id"],
			"\npush_provider :", result["push_provider"])
		
	}


}
