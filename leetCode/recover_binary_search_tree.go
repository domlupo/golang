// https://leetcode.com/problems/recover-binary-search-tree/
// N is the number of nodes
// time complexity is O(N)
// space complexity is O(N)
// Morris traversal can be used to reduce space complexity to O(1)

package main

func recoverTree(root *TreeNode) {
	nodes := []*TreeNode{}

	// perform in order traversal to get a sorted array with two swapped values
	curr := root
	q := []*TreeNode{}
	for curr != nil || len(q) != 0 {
		for curr != nil {
			q = append(q, curr)
			curr = curr.Left
		}
		curr = q[len(q)-1]
		q = q[:len(q)-1]
		nodes = append(nodes, curr)
		curr = curr.Right
	}

	// determine which nodes need swapped
	var first, second *TreeNode
	for i := 0; i < len(nodes)-1; i++ {
		if nodes[i+1].Val < nodes[i].Val && first == nil {
			first = nodes[i]
			second = nodes[i+1]
		} else if nodes[i+1].Val < nodes[i].Val {
			second = nodes[i+1]
		}
	}

	// swap the values for the nodes
	first.Val, second.Val = second.Val, first.Val
}
