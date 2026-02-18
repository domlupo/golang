// https://leetcode.com/problems/flipping-an-image/

package main

func flipAndInvertImage(image [][]int) [][]int {
	for i := 0; i < len(image); i++ {
		l := 0
		r := len(image[i]) - 1
		for l <= r {
			image[i][l], image[i][r] = 1-image[i][r], 1-image[i][l]
			l++
			r--
		}
	}

	return image
}
