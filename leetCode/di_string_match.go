// https://leetcode.com/problems/di-string-match/

package main

func diStringMatch(s string) []int {
	ans := make([]int, len(s)+1)
	low, high := 0, len(s)

	for i, c := range s {
		if c == 'D' {
			ans[i] = high
			high--
		} else if c == 'I' {
			ans[i] = low
			low++
		}
	}

	ans[len(s)] = high
	return ans
}
