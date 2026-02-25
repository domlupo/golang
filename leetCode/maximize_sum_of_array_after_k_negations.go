// https://leetcode.com/problems/maximize-sum-of-array-after-k-negations/

package main

import "math"

func largestSumAfterKNegations(nums []int, k int) int {
	for i := 0; i < len(nums)-1; i++ {
		if i >= 0 && nums[i] > nums[i+1] {
			nums[i], nums[i+1] = nums[i+1], nums[i]
			i -= 2
		}
	}

	smallestAbsVal := math.MaxInt
	var smallestAbsValI int
	for i, n := range nums {
		absValN := n
		if n < 0 {
			absValN *= -1
		}
		smallestAbsVal = min(smallestAbsVal, absValN)
		if smallestAbsVal == absValN {
			smallestAbsValI = i
		}
	}

	var ptr int
	for k > 0 {
		if ptr < len(nums) && nums[ptr] < 0 {
			nums[ptr] *= -1
			ptr++
		} else {
			nums[smallestAbsValI] *= -1
		}
		k--
	}

	var sum int
	for _, n := range nums {
		sum += n
	}
	return sum
}
