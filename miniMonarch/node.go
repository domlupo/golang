package main

type player int

const (
	noPlayer = iota
	playerOne
	playerTwo
)

type node struct {
	player player
	target *gridPosition
}

func newNode(p player) node {
	n := node{
		player: p,
		target: nil,
	}
	return n
}
