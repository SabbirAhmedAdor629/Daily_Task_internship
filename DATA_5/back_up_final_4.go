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
				// Dynamically decoding every objects of mapping.json
	mapping := map[string]map[string]interface{}{}
	for k,v := range result_map{
		switch c := v.(type){
		case map[string]interface{}:
			mapping[k] = v.(map[string]interface{})
		default:
			fmt.Printf("Not sure what type item %q is, but I think it might be %T\n", k, c)
		}
	}

				//Reading the data.json file and decoding
	empJson, _ := os.Open("data.json")
	defer empJson.Close()
	byteValue_2, _ := ioutil.ReadAll(empJson)
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue_2), &result)
				
	// Dynamically decoding every object of data.json
	data := map[string]map[string]interface{}{}

	
	// data2 := map[string]map[int]interface{}{}

	// data2["ad"][0] = "delivered"
	//  fmt.Println(data2["ad"]["push"])

	for k, v := range result {
		switch c := v.(type) {
		case map[string]interface{}:
			data[k] = v.(map[string]interface{})
		case []interface{}: 
			for _, item := range v.([]interface{}) {
				data[k] = item.(map[string]interface{})
				// fmt.Printf("%v : %v",k, data[k])
				// fmt.Println("")
			}
		default:
			fmt.Printf("Not sure what type item %q is, but I think it might be %T\n", k, c)
		}
	}

	//keys and values in data
	// for item := range data{
	// 	// for key := range data[item]{
	// 		fmt.Println(data[item])
	// 		fmt.Println()
	// 	// }
	// }


				// Taking another map to store keys and values
	CopiedMap:= make(map[string]interface{})
	for v := range result_map{
		for item := range mapping[v] {

			tF := strings.Contains(fmt.Sprint(mapping[v][item]), "[")

			if tF == true {
				value := (strings.Split(fmt.Sprint((mapping[v][item])),"."))
				v0 :=  (strings.Split(value[0], "["))
				for items := range data{
					if(items == v0[0]){
						for key := range data[v0[0]]{
							if(key == value[1]){
								CopiedMap[item] = data[items][key]
								fmt.Printf("%v : %v\n", item, CopiedMap[item])
							}
						}
					}
				}
			}else{
				value := (strings.Split(fmt.Sprint((mapping[v][item])),"."))
				CopiedMap[item] = data[value[0]][value[1]]
				fmt.Printf("%v : %v\n", item, CopiedMap[item])
				
			}
		}
	}

	// fmt.Println(data["push_messages[]"]["]status"])
	

}
































//https://stackoverflow.com/questions/23066758/how-can-i-write-an-array-of-maps-golang

// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// 	"strings"
// 	//"strings"
// )

// func main() {
// 				// Reading the mapping.json file and decoding
// 	jsonFile, _ := os.Open("mapping.json")
// 	defer jsonFile.Close()
// 	byteValue, _ := ioutil.ReadAll(jsonFile)
// 	var result_map map[string]interface{}
// 	json.Unmarshal([]byte(byteValue), &result_map)
// 				// Dynamically decoding every objects of mapping.json file
// 	mapping := map[string]map[string]interface{}{}
// 	for k,v := range result_map{
// 		switch c := v.(type){
// 		case map[string]interface{}:
// 			mapping[k] = v.(map[string]interface{})
// 		default:
// 			fmt.Printf("Not sure what type item %q is, but I think it might be %T\n", k, c)
// 		}
// 	}

// 				//Reading the data.json file and decoding
// 	empJson, _ := os.Open("data.json")
// 	defer empJson.Close()
// 	byteValue_2, _ := ioutil.ReadAll(empJson)
// 	var result map[string]interface{}
// 	json.Unmarshal([]byte(byteValue_2), &result)
				
// 	// Dynamically decoding every object of data.json file
// 	data := map[string]map[string]interface{}{}

// 	data2 := map[string][]map[string]interface{}{}
// 	fmt.Println(data2)
	
// 	for k, v := range result {
// 		switch c := v.(type) {
// 		case map[string]interface{}:
// 			data[k] = v.(map[string]interface{})
// 		case []interface{}: 
// 			for index, item := range v.([]interface{}) {
// 				data2[k][index] = (item.(map[string]interface{}))
// 				fmt.Printf("%v : %v",k, data2[k][index])
// 				fmt.Println("")
// 			}
// 		default:
// 			fmt.Printf("Not sure what type item %q is, but I think it might be %T\n", k, c)
// 		}
// 	}

// 	//keys and values in data
// 	// for item := range data{
// 	// 	// for key := range data[item]{
// 	// 		fmt.Println(data[item])
// 	// 		fmt.Println()
// 	// 	// }
// 	// }


// 				// Taking another map to store keys and values
// 	CopiedMap:= make(map[string]interface{})
// 	for v := range result_map{
// 		for item := range mapping[v] {

// 			tF := strings.Contains(fmt.Sprint(mapping[v][item]), "[")

// 			if tF == true {
// 				value := (strings.Split(fmt.Sprint((mapping[v][item])),"."))
// 				v0 :=  (strings.Split(value[0], "["))
// 				for items := range data{
// 					if(items == v0[0]){
// 						for key := range data[v0[0]]{
// 							if(key == value[1]){
// 								CopiedMap[item] = data[items][key]
// 								//fmt.Printf("%v : %v\n", item, CopiedMap[item])
// 							}
// 						}
// 					}
// 				}
// 			}else{
// 				value := (strings.Split(fmt.Sprint((mapping[v][item])),"."))
// 				CopiedMap[item] = data[value[0]][value[1]]
// 				//fmt.Printf("%v : %v\n", item, CopiedMap[item])
				
// 			}
// 		}
// 	}

// 	// fmt.Println(data["push_messages[]"]["]status"])
	

// }
