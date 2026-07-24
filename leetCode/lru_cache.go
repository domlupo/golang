// https://leetcode.com/problems/lru-cache/
// Get and Put are O(1) due to hash map and double linked list

package main

type LRUCache struct {
	Head, Tail *Node146
	Mp         map[int]*Node146
	Capacity   int
}

type Node146 struct {
	Prev, Next *Node146
	Key, Value int
}

func Constructor146(capacity int) LRUCache {
	head := &Node146{}
	tail := &Node146{}
	head.Next = tail
	tail.Prev = head
	return LRUCache{
		Head:     head,
		Tail:     tail,
		Mp:       make(map[int]*Node146),
		Capacity: capacity,
	}
}

func (this *LRUCache) Get(key int) int {
	if n, ok := this.Mp[key]; ok {
		this.remove(n)
		this.insert(n)
		return n.Value
	}

	return -1
}

func (this *LRUCache) Put(key int, value int) {
	// remove key if it exists
	if _, ok := this.Mp[key]; ok {
		this.remove(this.Mp[key])
	}

	// remove tail if needed
	if len(this.Mp) == this.Capacity {
		this.remove(this.Tail.Prev)
	}

	// insert key
	this.insert(&Node146{Key: key, Value: value})
}

func (this *LRUCache) remove(node *Node146) {
	delete(this.Mp, node.Key)
	node.Prev.Next = node.Next
	node.Next.Prev = node.Prev
}

func (this *LRUCache) insert(node *Node146) {
	this.Mp[node.Key] = node
	next := this.Head.Next
	this.Head.Next = node
	node.Prev = this.Head
	next.Prev = node
	node.Next = next
}
