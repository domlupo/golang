// https://leetcode.com/problems/permutations-ii/
// N is the length of nums
// time complexity is O(N * N!)
// space complexity is O(N)

package main

func permuteUnique(nums []int) [][]int {
	var ans [][]int
	counter := make(map[int]int)
	for _, n := range nums {
		counter[n]++
	}
	backtrack47(&ans, []int{}, nums, counter)
	return ans
}

func backtrack47(ans *[][]int, comb, nums []int, counter map[int]int) {
	if len(comb) == len(nums) {
		deepCopy := make([]int, len(comb))
		copy(deepCopy, comb)
		*ans = append(*ans, deepCopy)
		return
	}

	for n, count := range counter {
		if count == 0 {
			continue
		}
		comb = append(comb, n)
		counter[n]--
		backtrack47(ans, comb, nums, counter)
		counter[n]++
		comb = comb[:len(comb)-1]
	}
}
