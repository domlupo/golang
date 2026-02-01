// https://leetcode.com/problems/pascals-triangle/

// There are several solutions, I think this is simplest.

package main

func generate(numRows int) [][]int {
	t := [][]int{[]int{1}}

	for numRows > 1 {
		numRows--

		p := t[len(t)-1]
		n := []int{1}
		for i := 1; i < len(p); i++ {
			n = append(n, p[i-1]+p[i])
		}
		n = append(n, 1)

		t = append(t, n)
	}

	return t
}
