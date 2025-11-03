package main

import "fmt"

func main() {
	sum := sumNumbers([]int{1, 2, 3})
	fmt.Println(sum)
}

func sumNumbers(arr []int) int {
	var sum int = 0
	for _, num := range arr {
		sum += num
	}
	return sum
}
