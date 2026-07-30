// https://leetcode.com/problems/binary-search-tree-iterator/
// H is the height of the tree
// time complexity is O(1)
// space complexity is O(H)

package main

type TreeNode173 struct {
	Val   int
	Left  *TreeNode173
	Right *TreeNode173
}

type BSTIterator173 struct {
	Stack []*TreeNode173
}

func Constructor173(root *TreeNode173) BSTIterator173 {
	var stack []*TreeNode173
	n := root
	for n != nil {
		stack = append(stack, n)
		n = n.Left
	}
	return BSTIterator173{
		Stack: stack,
	}
}

func (this *BSTIterator173) Next() int {
	n := this.Stack[len(this.Stack)-1]
	this.Stack = this.Stack[:len(this.Stack)-1]
	if n.Right != nil {
		curr := n.Right
		for curr != nil {
			this.Stack = append(this.Stack, curr)
			curr = curr.Left
		}
	}
	return n.Val
}

func (this *BSTIterator173) HasNext() bool {
	return len(this.Stack) > 0
}
