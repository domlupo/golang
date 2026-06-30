# LeetCode

## Need to practice

- how to implement min and max heaps after import "container/heap"

## Need to learn

- how hash sets and hash maps are implemented.

## Revisit

- 22
  - Generate Parenthesis
  - time and space complexity aka Catalan numbers.
  - This seems useful https://leetcode.com/problems/generate-parentheses/solutions/3473534/generate-parentheses-by-leetcode-sf26/comments/2217787/
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
  - SQL question. I think this is for data engineers?
- 193, 195
  - bash

## Tips

- Can convert a string "s" into bytes via "bytes := []byte(s)".
- Can do s = "1" + s to avoid string mutation
- To get integer do, i, _ := strconv.Atoi(string(number1[p1]))
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

## Min and Max Heaps

- examples are LeetCode 703 and 1086
- complex to implement in Golang
- remember to call Heap.Push(*IntHeap, int) and Heap.Pop(*IntHeap).(int)

```
// IntHeap is a min-heap of ints.
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x any) {
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
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
