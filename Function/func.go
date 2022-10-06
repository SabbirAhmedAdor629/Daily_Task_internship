package main

import (
	"fmt"
)


type Human struct{
	name string
	height float64
}

func main() {
	ador := Human{
		name : "Ador",
		height: 6.50,
	}

	ador.ChangeName()
	fmt.Println(ador.name, ador.height)
}

func (ador *Human) ChangeName(){
	ador.name = "Sabbir"
	fmt.Println(ador.name, ador.height)
}










							//PRACTICE
// func main(){
// 	var divide func(float64, float64) (float64, error)
// 	divide = func(a,b float64)(float64, error){
// 		if (b == 0){
// 			return 0.0, fmt.Errorf("Can not divide by zero")
// 		}else{
// 			return a/b, nil
// 		}
// 	}

// 	d, err := divide(5.0, 3.0)

// 	if (err != nil){
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println(d)
// }






	 	// While using function as a variable it must be initalize befor calling 
// func main () {
// 	var divide func(float64, float64) (float64, error)
// 	divide = func(a,b float64) (float64, error){
// 		if (b == 0){
// 			return 0.0, fmt.Errorf("Can not divide by Zero")
// 		}else{
// 			return a / b, nil
// 		}
// 	}

// 	d, err := divide(5.0, 3.0)
// 	if err != nil{
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println(d)
// }



// func main(){
// 	var f func() = func() {
// 		fmt.Println("hello world")
// 	}
// 	f()
// }



// func main () {
// 	for i := 0; i<5; i++ {
// 		func (i int){
// 			fmt.Println("Hello world")
// 		}(i)
// 	}
	
// }



// func main () {
// 	d, err := Divide (5.0, 3.0)
// 	if err != nil{
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println(d)
// }

// func Divide(a,b float64) (float64, error){
// 	if b == 0 {
// 		return 0.0, fmt.Errorf("Can not divide by zero")
// 	} 

// 	return a / b, nil
// }


				//RETURNING NAMED VALUE
// func main(){
// 	s := sum(1,2,3,4,5)
// 	fmt.Println(s)
// }

// func sum(value... int) (result int){
// 	for _,i := range(value){
// 		result += i
// 	}
// 	return 
// }


				// RETURNING ADDRESS
// func main (){
// 	s := sum(1,2,3,4,5)
// 	fmt.Println(*s)
// }

// func sum(value...int) *int{
// 	result := 0
// 	for _, i := range value{ 
// 		result += i
// 	}
// 	return &result
// }



// func main(){
// 	greetings := "hello"
// 	name := "Sabbir"
// 	sayGreeting(&greetings, &name)
// 	fmt.Println(name)
// }

// func sayGreeting(greetings, name *string){
// 	*name = "Ahmed"
// 	fmt.Println(*name)
	
// }