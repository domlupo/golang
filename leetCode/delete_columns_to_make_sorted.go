// https://leetcode.com/problems/delete-columns-to-make-sorted/

package main

func minDeletionSize(strs []string) int {
	var size int

	for i := 0; i < len(strs[0]); i++ {
		prior := byte('a')
		for j := 0; j < len(strs); j++ {
			curr := strs[j][i]
			if curr < prior {
				size++
				break
			}
			prior = curr
		}
	}

	return size
}
