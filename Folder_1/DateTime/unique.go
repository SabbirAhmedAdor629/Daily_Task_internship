package main

import (
	"fmt"

	"github.com/mpvl/unique"
	_ "github.com/mpvl/unique"
)

func main() {
	a := []string{"h1", "h2", "h3", "h1", "h2", "sabbir"}
	unique.Strings(&a)
	fmt.Println(a)

}

 
// import "fmt"
 
// func unique(arr []string) []string {
//     occurred := map[string]bool{}
//     result := []string{}
//     for e := range arr {
//         if occurred[arr[e]] != true {
//             occurred[arr[e]] = true
//             result = append(result, arr[e])
//         }
//     }
 
//     return result
// }
 
// func main() {
//     array1 := []string{"one", "one", "two", "2", "3", "3"}
//     fmt.Println(array1)
//     unique_items := unique(array1)
//     fmt.Println(unique_items)
// }