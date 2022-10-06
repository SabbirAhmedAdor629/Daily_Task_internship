package custom_math

func Add(x, y int) int {
	return x + y
}

func Subtract(x, y int) int {
	return y-x
}

func divide(x, y int) float64 {
	if (y == 0){
		return float64(0.0)
	}
	return float64(x / y)
}

func multiply(x, y int) int {
	return x * y
}
