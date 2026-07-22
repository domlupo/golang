// https://leetcode.com/problems/sum-root-to-leaf-numbers/
// N is the number of nodes
// time complexity is O(N)
// space complexity is O(1) due to Morris traversal

package main

func sumNumbers(root *TreeNode) int {
	var sum, curr, steps int
	var pred *TreeNode
	n := root

	for n != nil {
		if n.Left != nil {
			pred = n.Left
			steps = 1
			for pred.Right != nil && pred.Right != n {
				pred = pred.Right
				steps++
			}
			if pred.Right == nil {
				curr = curr*10 + n.Val
				// create link
				pred.Right = n
				n = n.Left
			} else {
				// only add if pred is leaf node, not interior pred
				if pred.Left == nil {
					sum += curr
				}
				for i := 0; i < steps; i++ {
					curr /= 10
				}
				// break link
				pred.Right = nil
				n = n.Right
			}
		} else {
			// handles going right continually from root node
			curr = curr*10 + n.Val
			if n.Right == nil {
				sum += curr
			}
			n = n.Right
		}
	}

	return sum
}

/*

O(N) space complexity solution when not using Morris traversal

func sumNumbers(root *TreeNode) int {
	var total int
	sum(root, 0, &total)
	return total
}

func sum(n *TreeNode, curr int, total *int) {
	if n == nil {
		return
	}

	curr = curr*10 + n.Val

	if n.Left == nil && n.Right == nil {
		*total += curr
	} else {
		sum(n.Left, curr, total)
		sum(n.Right, curr, total)
	}
}

*/
