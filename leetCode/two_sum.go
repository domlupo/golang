// https://leetcode.com/problems/two-sum/

package main

func twoSum(nums []int, target int) []int {
	complements := make(map[int]int)

	for ni, n := range nums {
		if ci, ok := complements[n]; ok {
			return []int{ci, ni}
		}
		complements[target-n] = ni
	}

	panic("twoSum: no solution found")
}
