// https://leetcode.com/problems/zigzag-conversion/
// N is the length of s
// time complexity is O(N)
// space complexity is O(1)

package main

func convert(s string, numRows int) string {
	if numRows == 1 {
		return s
	}

	var ans []byte
	var start, secondStep int
	firstStep := (numRows - 1) * 2
	for len(ans) < len(s) {
		i := start
		incFirst := true
		for i < len(s) {
			ans = append(ans, s[i])
			if firstStep == 0 {
				i += secondStep
			} else if secondStep == 0 {
				i += firstStep
			} else if incFirst {
				i += firstStep
				incFirst = false
			} else {
				i += secondStep
				incFirst = true
			}
		}
		start++
		firstStep -= 2
		secondStep += 2
	}

	return string(ans)
}

// 2 rows
// 0 2 4 6
// 1 3 5
// first and last rows steps are 2
// middle steps are 0

// 3 rows
// 0   4
// 1 3 5 7
// 2   6
// first and last rows steps are 4
// middle steps
// 2 -> 2 -> 2

// 4 rows
// 0     6     12
// 1   5 7  11 13
// 2 4   8 10  14
// 3     9     15
// first and last rows steps are 6
// middle steps
// 4 -> 2 -> 4 -> 2
// 2 -> 4 -> 2 -> 4

// 5 row
// 0      8     16
// 1    7 9    15
// 2  6   10  14
// 3 5    11 13
// 4      12
// first and last rows steps are 8
// middle steps
// 6 -> 2 -> 6
// 4 -> 4 -> 4
// 2 -> 6 ->
