package main

import "fmt"

//using for loop
func sum_to_n_a(n int) int {
	sum := 0
	for i := 1; i <= n; i++ {
			sum += i
	}
	return sum
}

	//using arithmetic series
func sum_to_n_b(n int) int {
	return n * (n + 1) / 2
}

	//using recursion
func sum_to_n_c(n int) int {
	if n == 0 {
			return 0
	}
	return n + sum_to_n_a(n-1)
}

func main() {
	//using for loop
	//space complexity is O(n)
	//time complexity is O(1)
	result := sum_to_n_a(5)
	fmt.Printf("%d\n", result)

	//using arithmetic series
	//space complexity is O(1)
	//time complexity is O(1)
	result2 := sum_to_n_b(5)
	fmt.Printf("%d\n", result2)

	//using recursion
	//Time Complexity: O(n)
	//Space Complexity: O(n)
	result3 := sum_to_n_c(5)
	fmt.Printf("%d\n", result3)
}