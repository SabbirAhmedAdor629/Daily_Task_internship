package Custom_math2

import (
	"testing"
)


type Value struct{
	input1, input2 int
	output_sum int
	output_sub int
}
		// PRACTICING SINGLE FUNCTION TESTING
// func TestAdd(t *testing.T){
// 	adding := Add(2,3)
// 	if adding != 5 {
// 		t.Errorf("Adding failed, We expected %d and we got %d\n", 5,adding)
// 	}else{
// 		t.Logf("Adding sucedded, We expected %d and we got %d\n", 5, adding)
// 	}
// }


			// TESTING MULTIPLE FUNCTIONS TOGETHER
func TestAddSub(t *testing.T){
	testData := Value{
		input1: 1,
		input2: 2,
		output_sum: 3,
		output_sub: -1,
	}
	sum := Add(testData.input1, testData.input2)
	sub := Subtract(testData.input1, testData.input2)
	if (sum != testData.output_sum && sub != testData.output_sub){
		t.Errorf("Failed")
	}else{
		t.Logf("Passed")
	}
}
