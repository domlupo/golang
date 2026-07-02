// https://leetcode.com/problems/count-and-say/

package main

import "strconv"

func countAndSay(n int) string {
	ans := "1"

	for n > 1 {
		var chars string
		i := 0

		for i < len(ans) {
			count := 1
			for i+1 < len(ans) && ans[i] == ans[i+1] {
				count++
				i++
			}
			chars = chars + strconv.Itoa(count) + string(ans[i])
			i++
		}

		ans = chars
		n--
	}

	return ans
}

// 1 -> 1
// 2 -> 11
// 3 -> 21
// 4 -> 1211
// 5 -> 111221
// 6 -> 312211
// 7 -> 13112221
// 8 -> 111321321
