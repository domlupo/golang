// https://leetcode.com/problems/min-cost-climbing-stairs/
// N is the length of cost
// time complexity is O(N)
// space complexity is O(1)

package main

func minCostClimbingStairs(cost []int) int {
	var prev, curr int
	for i := 2; i < len(cost)+1; i++ {
		oneStepCost := curr + cost[i-1]
		twoStepCost := prev + cost[i-2]
		prev, curr = curr, min(oneStepCost, twoStepCost)
	}
	return curr
}

// 10 15 _
// min cost of len(2) would be min(oneStepBack, twoStepBack)
// min(10, 15)

// 10 15 20 _
// min cost of len(3), would be min(oneStepBack, twoStepBack)
// min (cost[0]+cost[2], cost[1])
// min (cost[i-2]+prev, cost[i-1]+curr
