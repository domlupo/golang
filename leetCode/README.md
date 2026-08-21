# LeetCode

## Need to practice

- multiple imports
- morris traversal
- heap implementation
- trie implementation
- LRU implementation
- divide and conquer (e.g. 23 merge k sorted lists)

## Need to learn

- how hash sets and hash maps are implemented.

## Revisit

- 71, 95, 96, 105, 106, 139, 148
  - skipped and need to return
- 22
  - Generate Parenthesis
  - time and space complexity aka Catalan numbers.
  - This seems useful https://leetcode.com/problems/generate-parentheses/solutions/3473534/generate-parentheses-by-leetcode-sf26/comments/2217787/
- 43
  - Very involved and specific. Work on this later.
- 157
  - Read N Characters Given Read4
  - This one's explanation is confusing
- 408, 1089
  - This is harder than an easy
- 170, 225, 232
  - Very few companies have used this recently and it is an implement class
- 338, 645
  - Bit manipulation
- 705, 706
  - I need to research how hash sets and hash maps are implemented

## Skipped

- 175, 181, 182, 183, 196, 197, 511, 512, 577, 584, 586, 595, 596, 597, 603, 607, 610, 613, 619,
  620, 627, 1050, 1068, 1069, 1075, 1076, 1082, 1083, 1084
  - Easy SQL question. I think this is for data engineers?
- 176, 177, 178, 180, 184
  - Med SQL question. I think this is for data engineers?
- 193, 195
  - Easy bash
- 192, 194
  - Med bash
- 156 only Linkedin and question is strange
- 164 Solution is strange and very few companies used recently
- 186 Easier than 151 and very few companies used recently
- 251 only Airbnb
- 254 only Linkedin
- 255 not many companies use

## Tips

- Can convert a string "s" into bytes via "bytes := []byte(s)".
- strings are immutable in Golang
- Can do s = "1" + s to avoid string mutation
- To get integer do, i, err := strconv.Atoi(string(number1[p1]))
- It is usually easiest to form strings by using s := rune[] or s := byte[] then string(s).
- string('A'+1) -> B, string('A'+2) -> C.
- strconv.Itoa(1)+"->"+strconv.Itoa(2) returns "1->2"
- s := "hello", s[5:] is "". Even though intuitively 5 is out of bounds
- odd position bit mask is 0x55555555 -> 0x55555555 & 0b11111111 -> 0b01010101
- even position bit mask:  0xAAAAAAAA -> 0xAAAAAAAA & 0b11111111 -> 0b10101010
- power of two check: (n & (n-1)) == 0. power of two means only 1 bit is set.
- When _, ok := contains[n]. ok is true if contains[n] was ever set. Even for contains[n] = false.
- fmt.Println(string('a'-32)) lower to upper
- fmt.Println(string('A'+32)) upper to lower

## Data structures

- heap / priority queue
- trie
- array
- monotonic stack
- linked list
- tries
- trees
- graphs

## Problem types

- arrays and hashing
- two pointers
- sliding window
- stack
- binary search
- linked list
- trees
- heap / priority queue
- backtracking
- tries
- graphs
- dynamic programming
- greedy
- intervals
- math
- bit manipulation

## Min and Max Heaps

- examples are LeetCode 23, 703 and 1086
- complex to implement in Golang
- try to think of ways to keep size of heap small e.g. push and pop 
  nodes of linked list as you go rather than push all node values first.

```
import "container/heap"

h := &Heap{}
heap.Init(h)
heap.Push(h, x)
n := heap.Pop(h).(*ListNode)

// minimum val *ListNode heap
type Heap []*ListNode

func (h Heap) Len() int           { return len(h) }
func (h Heap) Less(i, j int) bool { return h[i].Val < h[j].Val }
func (h Heap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *Heap) Push(x any) {
	*h = append(*h, x.(*ListNode))
}

func (h *Heap) Pop() any {
	old := *h
	n := len(old)-1
	x := old[n]
	*h = old[:n]
	return x
}
```

### Tree Traversals

For pre-order, in-order and post-order, left always occurs before right but when root is traversed
changes.

   1
  / \
 2   3

pre-order -> root left right -> 1 2 3
in-order -> left root right -> 2 1 3
post-order -> left right root -> 2 3 1

### Pascal's Triangle trick

next := prev * (rowIndex - col + 1) / col

where rowIndex starts at 0 and col starts at 0.

e.g. row 4 calculations

col 0 -> 1
col 1 -> 4 := 1 * (4 - 1 + 1 ) / 1
col 2 -> 6 := 4 * (4 - 2 + 1 ) / 2
col 3 -> 4 := 6 * (4 - 3 + 1 ) / 3
col 4 -> 1 := 4 * (4 - 4 + 1 ) / 4

### Binary Search Tree

- Clarify with interviewer or LeetCode question if the values are unique or not.
- Is a tree where all values left are lower, and all values right are higher unless there are
  duplicate values.
- The structure is typically:
  - call recursive function on left child of node
  - do work on node
  - call recursive function on right child of node

### Area of a triangle

// the point order should rotate in increasing or decreasing order.
// increasing
// 1 -> 2 -> 3
// 2 -> 3 -> 1
// 3 -> 1 -> 2
area := 0.5 * float64( p1x*(p2y-p3y) + p2x*(p3y-p1y) + p3x*(p1y-p2y) )
if area < 0 {
  area *= -1
}

// decreasing shown below
area := 0.5 * float64( p1x*(p3y-p2y) + p2x*(p1y-p3y) + p3x*(p2y-p1y) )
if area < 0 {
  area *= -1
}

### math.MaxInt bounds

Avoid math.MaxInt and math.MinInt when possible. They have weirdness near bounds for addition and subtraction.

i := math.MinInt
i -= 1
i *= -1 is still a negative number

math.MaxInt32 is 2147483647
math.MinInt32 is -2147483648
math.MaxInt is 9223372036854775807
math.MinInt is -9223372036854775808

### Brian Kernighan's algorithm to set count bits

func countSetBits(n int) int {
	var count int
	for n != 0 {
		n &= (n - 1) // Clears the rightmost set bit
		count++
	}
	return count
}

count = 0
111 &= 110 -> 110
count = 1
110 &= 101 -> 100
count = 2
100 &= 011 -> 000
count = 3

### List of common algorithms

- Rabin Karp algorithm uses hashing, rolling hash and pattern matching to find a short pattern
  inside a larger text. You can also use bitmask. See LeetCode 187

### Greatest Common Divisor found from Euclid's algorithm

func gcd(v1, v2 int) int{
  for v2 != 0 {
    v1, v2 = v2, v1%v2
  }
  return v1
}

### strings.Join(x, "")

can create a string from:

  x := []byte{'a', 'b', 'c'}
  string(x)

but use strings.Join(x, "") for below:

  x := []string{"a", "b", "c"}
  strings.Join(x, "")

### Complicated space complexities

Remember, permutation is where position matters and combination is collection only.

TODO: verify these space complexities are correct, I think they are wrong and differ by problem.

permutations: N!

               k      N!
combinations: C  = --------
               N   (N-k)!k!

subsets: 2^N

### Finding two swapped values in a sorted array

There are two cases to handle.

The first is when the swapped values are next to each other.
      _ _
1 2 3 5 4 6
The second is when the swapped values are not next to each other.
      _   _
1 2 3 6 5 4

Below is an example from LeetCode #99 Recover Binary Search Tree

var first, second *TreeNode
for i:=0; i<len(nodes)-1; i++ {
    if nodes[i+1].Val < nodes[i].Val && first == nil {
        first = nodes[i]
        second = nodes[i+1]
    } else if nodes[i+1].Val < nodes[i].Val {
        second = nodes[i+1]
    }
}

### Golang standard library sort

slices.Sort()
time complexity is O(n*log(n))
space complexity is O(log(n))

Don't use gnome sort aka stupid sort.
Gnome sort is O(n*n) time complexity and O(1) space complexity.

example sorting intervals [][]int

[[0,30],[5,10],[15,20],[0, 15]] -> [[0 30] [0 15] [5 10] [15 20]]

slices.SortFunc(intervals, func(a, b []int) int {
  if a[0] < b[0] {
    return -1
  } else {
    return 1
  }
})

### Morris Traversal

LeetCode 114 Flatten Binary Tree To Linked List has good solution to learn Morris Traversal

https://leetcode.com/problems/flatten-binary-tree-to-linked-list

### Questions to ask interviewer

- Do I need to explicitly write out imports?
- If input should be modified or not modified?
- Edge cases? e.g. 64 bit vs 32 bit, invalid cases like popping an empty array 
- What data structures can I import?
- Is the tree a binary search tree?

### Golang deleting a key from a maps

delete(mapName, keyName) is most correct.

Using delete ensures for val, ok := mapName[keyName], ok will be false instead of true. 
If you instead do mapName[keyName] = nil, ok will be true instead of false.

Another option is to check if val is equal to nil but that is a bit strange.


### Comparing slices and arrays

- compare arrays with a == b
- compare slices with slices.Equal(a, b)
