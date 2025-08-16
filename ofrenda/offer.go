package main

type kind int

const (
	noKind kind = iota
	skull
	candy
	flower
	bottle
)

type offer struct {
	kind kind
}

func newOffer(kind kind) *offer {
	s := offer{
		kind: kind,
	}
	return &s
}
