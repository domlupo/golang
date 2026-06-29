// https://leetcode.com/problems/zigzag-conversion/

package main

func convert(s string, numRows int) string {
	if numRows == 1 {
		return s
	}

	var ans []byte

	for row := range numRows {
		var step int
		firstStep := true
		for i := row; i < len(s); i += step {
			if row == 0 || row == numRows-1 {
				step = (numRows*2 - 2)
			} else if firstStep {
				step = (numRows*2 - 2) - (2 * row)
				firstStep = false
			} else {
				step = 2 * row
				firstStep = true
			}
			ans = append(ans, s[i])
		}
	}

	return string(ans)
}

// first row and last row step is (numRows * 2) - 2
// middle row has two partial steps. The first step
// (numRows * 2 - 2) - (2 * row) and the second step
// is 2 * row

// 2 rows
// 0 2 4 6
// 1 3 5

// 3 rows
// 0   4
// 1 3 5 7
// 2   6
// middle steps
// 2 -> 2 -> 2

// 4 rows
// 0     6     12
// 1   5 7  11 13
// 2 4   8 10  14
// 3     9     15
// middle steps
// 4 -> 2 -> 4 -> 2
// 2 -> 4 -> 2 -> 4

// 5 row
// 0      8     16
// 1    7 9    15
// 2  6   10  14
// 3 5    11 13
// 4      12
// middle steps
// 6 -> 2 -> 6
// 4 -> 4 -> 4
// 2 -> 6 ->
