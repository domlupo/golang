// https://leetcode.com/problems/valid-mountain-array/

package main

func validMountainArray(arr []int) bool {
	var decreasing, increasing bool
	prev := arr[0]

	for i := 1; i < len(arr); i++ {
		if arr[i] > prev && decreasing {
			return false
		} else if arr[i] > prev {
			increasing = true
		} else if arr[i] < prev {
			decreasing = true
		} else {
			return false
		}
		prev = arr[i]
	}
	return increasing && decreasing
}
