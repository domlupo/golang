// https://leetcode.com/problems/fibonacci-number/

package main

func fib(n int) int {
	if n == 0 {
		return 0
	}

	var prior int
	curr := 1
	for n > 2 {
		curr, prior = curr+prior, curr
		n--
	}

	return curr + prior
}
