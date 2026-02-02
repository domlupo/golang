// https://leetcode.com/problems/summary-ranges/

package main

import "strconv"

func summaryRanges(nums []int) []string {
	var o []string
	for i := 0; i < len(nums); i++ {
		a := nums[i]
		for i+1 < len(nums) && nums[i]+1 == nums[i+1] {
			i++
		}
		b := nums[i]

		if a == b {
			o = append(o, strconv.Itoa(a))
		} else {
			o = append(o, strconv.Itoa(a)+"->"+strconv.Itoa(b))
		}
	}

	return o
}
