package main

import "fmt"

func main() {
	// Create a map with some duplicate values
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 2, "e": 1}

	// Create a new map to hold the unique values
	uniqueMap := make(map[string]int)

	// Iterate through the original map
	for k, v := range m {
		// Check if the value has already been added to the new map
		if _, ok := uniqueMap[v]; !ok {
			// If it hasn't been added, add it to the new map
			uniqueMap[k] = v
		}
	}

	// Print the new map
	fmt.Println(uniqueMap) // Output: map[a:1 b:2 c:3]
}
