// https://leetcode.com/problems/largest-rectangle-in-histogram/
// N is the length of heights
// time complexity is O(N)
// space complexity is O(N)

package main

func largestRectangleArea(heights []int) int {
	stack := []int{-1}
	var maxArea int
	for i := range heights {
		for len(stack) > 1 && heights[stack[len(stack)-1]] >= heights[i] {
			h := heights[stack[len(stack)-1]]
			stack = stack[:len(stack)-1]
			w := i - stack[len(stack)-1] - 1
			maxArea = max(maxArea, h*w)
		}
		stack = append(stack, i)
	}

	for len(stack) > 1 {
		h := heights[stack[len(stack)-1]]
		stack = stack[:len(stack)-1]
		w := len(heights) - stack[len(stack)-1] - 1
		maxArea = max(maxArea, h*w)
	}

	return maxArea
}
