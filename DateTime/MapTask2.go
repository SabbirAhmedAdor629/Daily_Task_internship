package main

import "fmt"

func removeDuplicatesUnordered(elements []string) []string {
    encountered := map[string]bool{}

    // Create a map of all unique elements.
    for v:= range elements {
        encountered[elements[v]] = true
    }

    // Place all keys from the map into a slice.
    result := []string{}
    for key, _ := range encountered {
        result = append(result, key)
    }
    return result
}

func main() {
    elements := []string{"cat", "dog", "cat", "bird"}

    // Remove string duplicates, ignoring order.
    result := removeDuplicatesUnordered(elements)
    fmt.Println(result)
}