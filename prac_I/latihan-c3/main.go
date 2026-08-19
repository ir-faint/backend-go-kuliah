package main

import "fmt"

func swapByValue(x int, y int) (int, int) {
	tmp := x
	x = y
	y = tmp
	return x, y
}

func swapByReference(x *int, y *int) {
	tmp := *x
	*x = *y
	*y = tmp
}

func addSliceByValue(slice []string, value string) []string {
	slice = append(slice, value)
	return slice
}

func addSliceByReference(slice *[]string, value string) {
	*slice = append(*slice, value)
}

func main() {

	x, y := 10, 20
	fmt.Println("Before Swap By Value :", x, y)

	x, y = swapByValue(x, y)
	fmt.Println("After Swap By Value :", x, y)
	swapByReference(&x, &y)
	fmt.Println("After Swap By Reference :", x, y)

	fmt.Println("------------------------")

	var slice []string = []string{"Raihan"}
	fmt.Println("Before Add Slice By Value :", slice)

	slice = addSliceByValue(slice, "Irfandi")
	fmt.Println("After Add Slice By Value :", slice)
	addSliceByReference(&slice, "Khaliq")
	fmt.Println("After Add Slice By Reference :", slice)
}
