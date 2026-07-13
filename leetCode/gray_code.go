// https://leetcode.com/problems/two-sum/
// n is the number of bits used
// time complexity is O(2^n)
// space complexity is O(1)

package main

func grayCode(n int) []int {
	ans := []int{0}
	for i := 1; i <= n; i++ {
		prevLength := len(ans)
		mask := 1 << (i - 1)
		// grayCodes have an iterative pattern.
		// To find []int for n, iterate (n-1) in reverse
		// while adding a 1 digit in front.
		for j := prevLength - 1; j >= 0; j-- {
			ans = append(ans, mask+ans[j])
		}
	}

	return ans
}

// n = 1
// 0
// 1

// n = 2
// 00
// 01
// 11
// 10

// n = 3
// 000
// 001
// 011
// 010
// 110
// 111
// 101
// 10
