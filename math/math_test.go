package custom_math
import (
	"testing"
)

func TestAdd(t *testing.T){
	result := add(0,3)

	if result != 4{
		t.Errorf("Add(1,3) failed, we expected %d and we get %d\n",4,result)
	}else{
		t.Logf("Add(1,3) passed, we expected %d and we get %d\n",4,result)
	}
}