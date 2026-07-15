// M (3) is the number of digits, N (4) is the knumber of integers
// time complexity is O(M^N * N) = O(1)
// space complexity is O(M * N) = O(1)

package main

import "strconv"

func restoreIpAddresses(s string) []string {
	var ans []string
	backtrack93(&ans, "", s, 0, 0)
	return ans
}

func backtrack93(ans *[]string, curr, s string, dots, i int) {
	// There are other ways to early exit but they are not needed.
	// For example, len(s) > 12, or comparing dots and index.
	if dots == 4 && i == len(s) {
		*ans = append(*ans, curr[:len(curr)-1])
		return
	} else if dots == 4 {
		return
	}

	var one, two, three string
	if i+1 <= len(s) {
		one = s[i : i+1]
	}
	if i+2 <= len(s) {
		two = s[i : i+2]
	}
	if i+3 <= len(s) {
		three = s[i : i+3]
	}

	if validInteger(one) {
		backtrack93(ans, curr+one+".", s, dots+1, i+1)
	}
	if validInteger(two) {
		backtrack93(ans, curr+two+".", s, dots+1, i+2)
	}
	if validInteger(three) {
		backtrack93(ans, curr+three+".", s, dots+1, i+3)
	}
}

func validInteger(s string) bool {
	n, _ := strconv.Atoi(s)

	if len(s) == 1 {
		return true
	}
	if len(s) == 2 && n >= 10 && n <= 100 {
		return true
	}
	if len(s) == 3 && n >= 100 && n <= 255 {
		return true
	}

	return false
}
