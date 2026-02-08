// https://leetcode.com/problems/reshape-the-matrix/

package main

func matrixReshape(mat [][]int, r int, c int) [][]int {
	if len(mat)*len(mat[0]) != r*c {
		return mat
	}

	o := make([][]int, r)
	for i := 0; i < r; i++ {
		o[i] = make([]int, c)
	}

	var ri, ci int
	for i := range len(mat) {
		for j := range len(mat[i]) {
			if ci >= c {
				ri += 1
				ci = 0
			}
			o[ri][ci] = mat[i][j]
			ci++
		}
	}

	return o
}
