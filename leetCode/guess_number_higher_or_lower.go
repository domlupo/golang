// https://leetcode.com/problems/guess-number-higher-or-lower/

package main

func guessNumber(n int) int {
	var a int
	l := 1
	r := n * 2
	for {
		n = l + (r-l)/2
		a = guess(n)
		if a == 0 {
			break
		}

		if a == 1 {
			l = n + 1
		} else {
			r = n - 1
		}
	}

	return n
}

func guess(num int) int {
	return 0
}
