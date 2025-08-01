package main

type moverDirection int

const (
	mNone moverDirection = iota
	mDown
	mUp
)

type moverLength int

const (
	none moverLength = iota
	short
	long
)

type mover struct {
	length    moverLength
	direction moverDirection
}
