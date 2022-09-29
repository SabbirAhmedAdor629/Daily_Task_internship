package gfg_go

import (
	"testing"
	"fmt"
	"strings"
)

// function to test if the original
// go program is working
func TestFunction(test *testing.T){

	test_string1 := "go_modules"

	// calling the function from
	// the previous go file
	res := strings.Split(initialiser(), ":")

	// removing spaces and line-ending
	// punctuation before comparing
	test_string2 := strings.Trim(res[1], " .")

	if test_string1 == test_string2 {
			fmt.Printf("Successful!\n")

	} else {
				// this prints error message if
				// strings do not match
			test.Errorf("Error!\n")
	}

}

func initialiser() {
	panic("unimplemented")
}
