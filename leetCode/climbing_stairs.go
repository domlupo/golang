// https://leetcode.com/problems/climbing-stairs/

package main

// iterative (best space and time complexity)

func climbStairs(n int) int {
	prev := 1
	curr := 2

	for ; n > 1; n-- {
		prev, curr = curr, prev+curr
	}

	return prev
}

// dynamic programming

var ways map[int]int = map[int]int{1: 1, 2: 2}

func climbStairsDP(n int) int {
	if w, ok := ways[n]; ok {
		return w
	} else {
		ways[n] = climbStairsDP(n-1) + climbStairsDP(n-2)
		return ways[n]
	}
}
