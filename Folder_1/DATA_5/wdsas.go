// package main
 
// import (
//     "fmt"
//     "strings"
// )
 
// func main() {       
//     s := "This,is,a,delimited,string"
//     v := strings.Split(s, ",")  
//     fmt.Println(v)     // [This is a delimited string]
// }

// package main
 
// import (
//     "fmt"
//     "strings"
// )
 
// func main() {       
//     s := "String,with,delimiters."
//     v := strings.SplitAfter(s, ",") // It doesn't removes delimiters when splitting
//     fmt.Println(v)     // [String, with, delimiters.]
// } 

// Go program to illustrate how to
// trim right-hand side elements
// from the string
package main
  
import (
    "fmt"
    "strings"
)
  
// Main method
func main() {
  
    // Creating and initializing the
    // string using shorthand declaration
    str1 := "elx.tssdfo"
    str2 := "@@This is the tutorial of Golang$$"
  
    // Displaying strings
    fmt.Println("Strings before trimming:")
    fmt.Println("String 1: ", str1)
    fmt.Println("String 2:", str2)
  
    // Trimming the given strings
    // Using TrimRight() function
    res1 := strings.TrimRight(str1, ".t")
    res2 := strings.TrimRight(str2, "$")
  
    // Displaying the results
    fmt.Println("\nStrings after trimming:")
    fmt.Println("Result 1: ", res1)
    fmt.Println("Result 2:", res2)
}