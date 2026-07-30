// https://leetcode.com/problems/largest-number/
// N is the length of nums
// time complexity is O(N * log(N))
// space complexity is O(N)

package main

import (
	"sort"
	"strconv"
)

func largestNumber(nums []int) string {
	sort.Slice(nums, func(i, j int) bool {
		iStr := strconv.Itoa(nums[i])
		jStr := strconv.Itoa(nums[j])
		ijStr, _ := strconv.Atoi(iStr + jStr)
		jiStr, _ := strconv.Atoi(jStr + iStr)
		return ijStr > jiStr
	})

	var sum int
	var ans string
	for _, n := range nums {
		sum += n
		ans = ans + strconv.Itoa(n)
	}

	if sum == 0 {
		return "0"
	} else {
		return ans
	}
}
