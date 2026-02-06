// https://leetcode.com/problems/fizz-buzz/

package main

import "strconv"

func fizzBuzz(n int) []string {
	var o []string

	for i := 1; i <= n; i++ {
		if i%3 == 0 && i%5 == 0 {
			o = append(o, "FizzBuzz")
		} else if i%3 == 0 {
			o = append(o, "Fizz")
		} else if i%5 == 0 {
			o = append(o, "Buzz")
		} else {
			o = append(o, strconv.Itoa(i))
		}
	}

	return o
}
