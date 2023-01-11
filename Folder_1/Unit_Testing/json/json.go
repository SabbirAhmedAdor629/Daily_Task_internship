package custom_json

import (
    "encoding/json"
    //"fmt"
)

type Name struct{
	Firstname string
	Lastname  string
}

func Decoding(jsondata string)(Name, error){
	var name1 Name
				  
	err := json.Unmarshal([]byte(jsondata), &name1)
	
	if err!=nil {
		return nil, err
	}
	return name1, nil
	

}