package main

import (
	"fmt"
	"math"
)


type Shape interface{
	Area() float64
	Cicumf() float64
}


type Square struct{
	Length float32
}

type Circle struct{
	Radius float64
}

// Methods of Square
func (S Square) Area() float64{
	return float64(S.Length * S.Length)
}

func (S Square) Cicumf() float64 {
	return float64(math.Pi * S.Length)
}

// Mehtods of Circle
func (C Circle) Area() float64{
	return C.Radius*C.Radius
}

func (C Circle) Cicumf() float64{
	return math.Pi * C.Radius
}

// func PrintShapeInfo(S Shape){
// 	fmt.Printf("The area of %T is : %0.2f\n", S, S.Area())
// 	fmt.Printf("The Cicumf of %T is : %0.2f\n", S, S.Cicumf())
// }

func main(){
	shapes := []Shape{
		Square{Length: 2.75},
		Circle{Radius: 3.85},
	}

	for index, val := range shapes{
		// PrintShapeInfo(val)
		// fmt.Println("------")
		fmt.Printf("The area of %T is : %0.2f\n", index, val.Area())
		fmt.Printf("The Cicumf of %T is : %0.2f\n", index, val.Cicumf())
	}
}