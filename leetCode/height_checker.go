// https://leetcode.com/problems/height-checker/

package main

func heightChecker(heights []int) int {
	sorted := make([]int, len(heights))
	copy(sorted, heights)

	for i := 0; i < len(sorted)-1; i++ {
		if i >= 0 && sorted[i] > sorted[i+1] {
			sorted[i], sorted[i+1] = sorted[i+1], sorted[i]
			i -= 2
		}
	}

	var unexpected int
	for i, h := range heights {
		if h != sorted[i] {
			unexpected++
		}
	}

	return unexpected
}
