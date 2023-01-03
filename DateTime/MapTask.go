package main

import "fmt"

func main() {

	m := map[string]string{"a": "sabbir", "b": "Ahmed", "c": "ador", true: "sabbir", "e": "sabbir"}

	seen := make(map[string]bool)
	
	for k, v := range m {
		fmt.Println(v)
		if seen[v] {
			delete(m, k)
		} else {
			seen[v] = true
		}
	}

	
	fmt.Println(m) 
}
