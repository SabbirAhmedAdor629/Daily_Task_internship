package main

import (
	"fmt"
	greeting "greetings"
	//"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
)

func main() {
	var message string
	message = greeting.M()
	
	fmt.Println(greeting.M2("Sabbir"))
	fmt.Println(message)
}
