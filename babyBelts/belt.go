package main

type Direction string

const (
	Up    Direction = "up"
	Left  Direction = "left"
	Down  Direction = "down"
	Right Direction = "right"
)

type belt struct {
	direction Direction
	x         int
	y         int
}
