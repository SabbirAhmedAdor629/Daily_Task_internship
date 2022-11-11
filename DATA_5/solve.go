
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	//"strings"
)

func main() {
	// Reading the mapping.json file and decoding
	jsonFile, _ := os.Open("mapping.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result_map map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result_map)

	//Reading the data.json file and decoding
	empJson, _ := os.Open("data.json")
	defer empJson.Close()
	byteValue_2, _ := ioutil.ReadAll(empJson)
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue_2), &result)

	error, CopiedMap := getDyanmicMappingStruct(result_map, result)
	if error == nil{
		for key := range CopiedMap{
			fmt.Printf(" %v : %v\n", key, CopiedMap[key] )
		}
	}

	
				
}



func getDyanmicMappingStruct(mappingData map[string]interface{}, jsonData map[string]interface{}) (error, map[string]interface{}) {
	dynamoItemStruct := make(map[string]interface{})
	for k, v := range mappingData["mapping"].(map[string]interface{}) {

		keys := strings.Split(k, ".")
		values := strings.Split(v.(string), ".")

		if len(keys) > 1 && len(values) > 1 {

			if keys[0][len(keys[0])-2:] == "[]" {
				structKyeFirst := keys[0][:len(keys[0])-2]
				jsonKey := values[0][:(len(values[0]) - 2)]
				objArr, ok := jsonData[jsonKey].([]interface{})
				if !ok {
					return fmt.Errorf("Json array []interface conversion error"), nil
				}

				if _, ok := dynamoItemStruct[structKyeFirst]; !ok {
					dynamoItemStruct[structKyeFirst] = []map[string]interface{}{}
					// assigned empty map array
					for i := 0; i < len(objArr); i++ {
						dynamoItemStruct[structKyeFirst] = append(dynamoItemStruct[structKyeFirst].([]map[string]interface{}), map[string]interface{}{})
					}
				}
				for i := 0; i < len(objArr); i++ {
					dynamoItemStruct[structKyeFirst].([]map[string]interface{})[i][keys[1]] = objArr[i].(map[string]interface{})[values[1]]
				}
			}
		} else {
			dynamoItemStruct[keys[0]] = jsonData[values[0]].(map[string]interface{})[values[1]]
		}
	}
	return nil, dynamoItemStruct

}