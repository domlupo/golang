// https://leetcode.com/problems/two-sum/
// N is the length of nums
// time complexity is O(N)
// space complexity is O(N)

package main

func twoSum1(nums []int, target int) []int {
	complements := make(map[int]int)

	for ni, n := range nums {
		if ci, ok := complements[n]; ok {
			return []int{ci, ni}
		}
		complements[target-n] = ni
	}

	panic("twoSum: no solution found")
}
