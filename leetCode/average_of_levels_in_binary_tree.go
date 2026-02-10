// https://leetcode.com/problems/average-of-levels-in-binary-tree/

package main

func averageOfLevels(root *TreeNode) []float64 {
	var avgs []float64
	if root == nil {
		return avgs
	}

	q := []*TreeNode{root}
	for len(q) > 0 {
		var sum int
		n := len(q)
		for range n {
			curr := q[0]
			q = q[1:]
			sum += curr.Val

			if curr.Left != nil {
				q = append(q, curr.Left)
			}
			if curr.Right != nil {
				q = append(q, curr.Right)
			}
		}

		avgs = append(avgs, float64(sum)/float64(n))
	}

	return avgs
}

func averageOfLevelsRecursive(root *TreeNode) []float64 {
	sums := make(map[int]int)
	counts := make(map[int]int)
	h3(root, 0, sums, counts)

	var o []float64
	for i := range len(sums) {
		o = append(o, float64(sums[i])/float64(counts[i]))
	}
	return o
}

func h3(n *TreeNode, depth int, sums, counts map[int]int) {
	if n == nil {
		return
	}

	sums[depth] = sums[depth] + n.Val
	counts[depth] = counts[depth] + 1

	h3(n.Left, depth+1, sums, counts)
	h3(n.Right, depth+1, sums, counts)
}
