package custom_math
import (
	"testing"
)

type AddData struct{
	x,y int
	output int
}

func TestAdd_Sub(t *testing.T){

	//TESTING ADD FUNCTION ONLY
	// result := Add(1,3)
	// if result != 4{
	// 	t.Errorf("Add(1,3) failed, we expected %d and we get %d\n",4,result)
	// }else{
	// 	t.Logf("Add(1,3) passed, we expected %d and we get %d\n",4,result)
	// }
	
	testData_add := []AddData{
		{3,5,8},
		{7,-4,3},
		{2,5,7},
	}

	for _, Datum := range(testData_add){
		result := Add(Datum.x,Datum.y)
		if(result != Datum.output){
			t.Errorf("Add failed, we expected %d and we get %d\n",Datum.output,result)
		}else{
			t.Logf("Add passed, we expected %d and we get %d\n",Datum.output,result)
		}
	}

	testData_sub := []AddData{
		{3,5,2},
		{7,8,1},
		{5,12,7},
	}

	for _, Datum := range(testData_sub){
		result := Subtract(Datum.x,Datum.y)
		if(result != Datum.output){
			t.Errorf("Subtract failed, we expected %d and we get %d\n",Datum.output,result)
		}else{
			t.Logf("Subtract passed, we expected %d and we get %d\n",Datum.output,result)
		}
	}
}

// func TestDivide(t *testing.T){
// 	result := divide(5,0)
// 	if result != 0.0{
// 		t.Errorf("Divide(5,0) failed, we expected %f and we got %f\n",0.0 , result)
// 	}else{
// 		t.Logf("Divide(5,0) passed, we expected %f and we got %f\n",0.0, result)
// 	}
// }