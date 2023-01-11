package main

import (
    "encoding/json"
    "fmt"
)

type Name struct{
	Firstname string
	Lastname string
}

func main(){

	var name1 Name

    jsondata := `{
		"firstname":"Sabbir",
		"lastname":"Ahmed"
	}`
				  
	err := json.Unmarshal([]byte(jsondata), &name1)
	
	if err == nil {
	        fmt.Println(name1.Firstname)
	        fmt.Println(name1.Lastname)
	    } else {
	        fmt.Println(err)
	    }
}








// type People struct {
//     Firstname string
//     Lastname string
// }

// func main() {
//     var person People

//     jsonString := `{"firstname":"Sabbir",
//                     "lastname":"Ahmed"}`

//     err := json.Unmarshal([]byte(jsonString), &person)

//     if err == nil {
//         fmt.Println(person.Firstname)
//         fmt.Println(person.Lastname)
//     } else {
//         fmt.Println(err)
//     }
// }