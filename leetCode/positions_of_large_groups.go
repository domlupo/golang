// https://leetcode.com/problems/positions-of-large-groups/

package main

func largeGroupPositions(s string) [][]int {
	var l int
	var char rune

	var o [][]int
	for i, ch := range s {
		if char == ch {
			continue
		}

		if i-l >= 3 {
			o = append(o, []int{l, i - 1})
		}

		char = ch
		l = i
	}

	if len(s)-l >= 3 {
		o = append(o, []int{l, len(s) - 1})
	}

	return o
}
