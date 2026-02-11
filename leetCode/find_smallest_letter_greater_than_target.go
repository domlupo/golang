// https://leetcode.com/problems/find-smallest-letter-greater-than-target/

package main

import "math"

func nextGreatestLetter(letters []byte, target byte) byte {
	m := math.MaxInt
	o := letters[0]

	var l int
	r := len(letters) - 1
	for l <= r {
		mid := l + (r-l)/2
		if letters[mid] <= target {
			l = mid + 1
		} else {
			r = mid - 1
			m = min(m, int(letters[mid]-target))
		}

		if m == int(letters[mid]-target) {
			o = letters[mid]
		}
	}

	return o
}
