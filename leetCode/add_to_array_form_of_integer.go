// https://leetcode.com/problems/add-to-array-form-of-integer/

package main

func addToArrayForm(num []int, k int) []int {
	var carry int
	i := len(num) - 1
	var rev []int

	for i >= 0 || k > 0 {
		v1, v2 := 0, k%10
		k /= 10
		if i >= 0 {
			v1 = num[i]
		}
		i--

		total := carry + v1 + v2
		carry = total / 10
		rev = append(rev, total%10)
	}

	var ans []int
	for i := len(rev) - 1; i >= 0; i-- {
		ans = append(ans, rev[i])
	}

	if carry == 1 {
		return append([]int{1}, ans...)
	} else {
		return ans
	}
}
