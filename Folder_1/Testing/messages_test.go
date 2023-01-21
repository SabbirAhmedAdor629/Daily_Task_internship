package main

import (
	"testing"
	//	"time"
	"fmt"
	"strings"

	"github.com/google/uuid"
)


// TESTING 3 (build failed because of other fuction)
type TestCase struct {
    name string
    text string
    vars map[string]interface{}
    expectedOutput string
}

var testCases = []TestCase{
    {name: "valid replacement", text: "Hello [name]", vars: map[string]interface{}{"name": "John"}, expectedOutput: "Hello John"},
    {name: "invalid replacement", text: "Hello [age]", vars: map[string]interface{}{"name": "John"}, expectedOutput: "Hello "},
    {name: "no replacement", text: "Hello", vars: map[string]interface{}{"name": "John"}, expectedOutput: "Hello"},
}

func TestReplaceText(t *testing.T) {
    for _, tc := range testCases {
        output := replaceText(tc.text, tc.vars)
        if output != tc.expectedOutput {
            t.Errorf("Test case: %s, expected: %s, got: %s", tc.name, tc.expectedOutput, output)
        }
    }
}

















// TESTING 2 (Facing issues while building the code on uuid)
// type TestCase struct {
//     name string
//     guidType string
//     expectedOutput string
// }
// func TestAssignGuid(t *testing.T) {
// 	var testCases = []TestCase{
// 		{name: "PushMessage", guidType: "PushMessage", expectedOutput: "mpm-" + uuid.New()},
// 		{name: "Bonus", guidType: "Bonus", expectedOutput: "mb-" + uuid.New()},
// 		{name: "Message", guidType: "Message", expectedOutput: "mm-" + uuid.New()},
// 		{name: "Default", guidType: "", expectedOutput: uuid.New()},
// 	}

//     for _, tc := range testCases {
//         output := assignGuid(tc.guidType)
//         if output != tc.expectedOutput {
//             t.Errorf("Test case: %s, expected: %s, got: %s", tc.name, tc.expectedOutput, output)
//         }
//     }
// }



// TESTING 1 (Fraction result are little bit different)
// type testCase struct {
// 	expiresIn int32
// 	unit      string
// 	expected  time.Time
// }
// func TestTimeFormatWithUnit(t *testing.T) {
// 	testCases := []testCase{
// 		{
// 			expiresIn: 10,
// 			unit:      "minutes",
// 			expected:  time.Now().UTC().Add(time.Minute * 10),
// 		},
// 		{
// 			expiresIn: 2,
// 			unit:      "hours",
// 			expected:  time.Now().UTC().Add(time.Hour * 2),
// 		},
// 		{
// 			expiresIn: 3,
// 			unit:      "days",
// 			expected:  time.Now().UTC().Add(time.Hour * 72),
// 		},
// 		{
// 			expiresIn: 3,
// 			unit:      "Invalid",
// 			expected:  time.Now().UTC(),
// 		},
// 	}
// 	for _, tc := range testCases {
// 		result := timeFormatWithUnit(tc.expiresIn, tc.unit)
// 		if result != tc.expected {
// 			t.Errorf("For expiresIn = %d and unit = %s, expected %s but got %s", tc.expiresIn, tc.unit, tc.expected, result)
// 		}
// 	}
// }
