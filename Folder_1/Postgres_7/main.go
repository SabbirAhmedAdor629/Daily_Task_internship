package main

import "fmt"

// Calculator is an interface that defines methods for addition and subtraction.
type Calculator interface {
	Add(a, b int) int
	Subtract(a, b int) int
}

type Calculator_4 struct{}

func (c *Calculator_4) Add(a, b int) int {
	return a + b
}

func (c *Calculator_4) Subtract(a, b int) int {
	return a - b
}

type cal_4_2 struct{}

func (c *cal_4_2) Add(a, b int) int {}

func (c *cal_4_2) Subtract(a, b int) int {}

func main() {
	var c Calculator = &Calculator_4{}
	result := c.Add(1, 2)
	fmt.Println(result) // Output: 3
	result = c.Subtract(1, 2)
	fmt.Println(result) // Output: -1

	var d Calculator = &cal_4_2{}
	res := d.Add(2, 2)
	fmt.Println(res) //4
	res = d.Subtract(2, 2)
	fmt.Println(res) //0

}
