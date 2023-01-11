package main

import "fmt"

type Shape struct {	// circle
	a int
	b int
	c int
}

func (c Cal) Add() { 
	fmt.Println(c.a+c.b)	
}

type Cal_2 struct { // Rectangle
	Cal
}
// reac
func main() {
	e := Cal_2{
		Cal{
			a: 5,
			b: 60,
		},
	}

	e.Add()
	// Call the Greet method of the Person field
	
}
