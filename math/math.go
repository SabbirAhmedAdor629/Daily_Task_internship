package custom_math

func add(x, y int) int {
	return x + y
}

func subtract(x, y int) int {
	return x - y
}

func divide(x, y int) float64 {
	if (y==0){
		return	0;
	}
	return float64(x / y)
}

func multiply(x, y int) int {
	return x * y
}
