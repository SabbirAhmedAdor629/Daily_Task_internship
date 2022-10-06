package main

import (
	"fmt"
)

func main (){
	s := sum(1,2,3,4,5)
	fmt.Println(*s)
}

func sum(value...int) *int{
	result := 0
	for _, V := range value{ 
		result += V
	}
	return &result
}



// func main(){
// 	greetings := "hello"
// 	name := "Sabbir"
// 	sayGreeting(&greetings, &name)
// 	fmt.Println(name)
// }

// func sayGreeting(greetings, name *string){
// 	*name = "Ahmed"
// 	fmt.Println(*name)
	
// }