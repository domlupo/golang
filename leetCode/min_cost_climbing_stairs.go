// https://leetcode.com/problems/min-cost-climbing-stairs/

package main

func minCostClimbingStairs(cost []int) int {
	oneStepBackCost := 0
	twoStepBackCost := 0

	for i := 2; i < len(cost)+1; i++ {
		oneStepCost := oneStepBackCost + cost[i-1]
		twoStepCost := twoStepBackCost + cost[i-2]
		oneStepBackCost, twoStepBackCost = min(oneStepCost, twoStepCost), oneStepBackCost
	}

	return oneStepBackCost
}
