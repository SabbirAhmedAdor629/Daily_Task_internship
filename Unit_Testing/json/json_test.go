package custom_json

import (
	//"bytes"
	//"fmt"
	"testing"
)

func TestDecoding(t *testing.T){
	jsondata := `{
		"firstname":"Sabbir",
		"lastname":"Ahmed"
	}`
	//jsonBytes := []byte(jsondata)
	person, err := Decoding(jsondata)

	
	if err != nil{
		if person.Firstname,"Sabbir"{
			t.Errorf("Decoding failed")
		}else{
			t.Logf("Decoding failed")
		}
	}
}
