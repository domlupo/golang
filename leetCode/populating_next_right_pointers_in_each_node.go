// https://leetcode.com/problems/populating-next-right-pointers-in-each-node/
// N is the number of tree nodes
// time complexity is O(N)
// space complexity is O(1)

package main

type Node116 struct {
	Val   int
	Left  *Node116
	Right *Node116
	Next  *Node116
}

func connect(root *Node116) *Node116 {
	if root == nil {
		return root
	}

	leftmost := root
	for leftmost.Left != nil {
		n := leftmost
		for n != nil {
			n.Left.Next = n.Right
			if n.Next != nil {
				n.Right.Next = n.Next.Left
			}
			n = n.Next
		}
		leftmost = leftmost.Left
	}

	return root
}

/*

O(N) space solution

func connect(root *Node1) *Node1 {
	queue := []*Node1{root}
	depths := []int{1}

	for len(queue) != 0 {
		n := queue[0]
		d := depths[0]
		queue = queue[1:]
		depths = depths[1:]
		if n == nil {
			continue
		}
		queue = append(queue, n.Left, n.Right)
		depths = append(depths, d+1, d+1)
		if d == depths[0] {
			n.Next = queue[0]
		}
	}

	return root
}

*/
