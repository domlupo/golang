package main

type kind int

const (
	noKind kind = iota
	circle
	square
	cross
)

type shapeState int

const (
	unknown shapeState = iota
	free
	locked
)

type shape struct {
	kind       kind
	shapeState shapeState
}

func newShape(kind kind, state shapeState) *shape {
	s := shape{
		kind:       kind,
		shapeState: free,
	}
	return &s
}
