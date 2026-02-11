// https://leetcode.com/problems/self-dividing-numbers/

package main

import "strconv"

func selfDividingNumbers(left int, right int) []int {
	var o []int

	for left <= right {
		selfDividing := true
		num := strconv.Itoa(left)

		for _, c := range num {
			n, _ := strconv.Atoi(string(c))
			if n == 0 || left%n != 0 {
				selfDividing = false
				break
			}
		}

		if selfDividing {
			o = append(o, left)
		}

		left++
	}

	return o
}
