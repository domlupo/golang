// https://leetcode.com/problems/populating-next-right-pointers-in-each-node-ii/
// N is the number of nodes in the tree
// time complexity is O(N)
// space complexity is O(1)

package main

type Node117 struct {
	Val   int
	Left  *Node117
	Right *Node117
	Next  *Node117
}

func connect117(root *Node117) *Node117 {
	if root == nil {
		return root
	}

	leftmost := root
	for leftmost != nil {
		n := leftmost
		leftmost = nil
		var prev *Node117
		for n != nil {
			if leftmost == nil && n.Left != nil {
				leftmost = n.Left
				prev = leftmost
			} else if leftmost == nil && n.Right != nil {
				leftmost = n.Right
				prev = leftmost
			}

			if prev != nil && n.Left != nil && prev != n.Left {
				prev.Next = n.Left
				prev = prev.Next
			}
			if prev != nil && n.Right != nil && prev != n.Right {
				prev.Next = n.Right
				prev = prev.Next
			}

			n = n.Next
		}
	}

	return root
}
