package main

import "fmt"

func main() {
	translations := map[string]string{
		"hi_msg":      "Hi!",
		"my_name_msg": "My name is",
		"welcome_msg": "Welcome to my website :)",
	}

	// Use make() to create the slice for better performance
	translationKeys := make([]string, 0, len(translations))

	// We only need the keys
	for key := range translations {
		translationKeys = append(translationKeys, key)
	}

	for k := range translationKeys {
		fmt.Println(translationKeys[k])
	}
	
}