// https://leetcode.com/problems/array-partition/

package main

func arrayPairSum(nums []int) int {
	counts := make([]int, 20_001)
	for _, n := range nums {
		i := n + 10_000
		counts[i] = counts[i] + 1
	}

	var sum int
	skip := true
	add := !skip
	p := len(counts) - 1
	for p >= 0 {
		if counts[p] == 0 {
			p--
			continue
		}
		if add {
			sum += p - 10_000
		}
		counts[p] = counts[p] - 1
		add = !add
		skip = !skip
	}

	return sum
}
