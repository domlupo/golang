package main

type producer struct {
	kind kind
}

func newProducer(kind kind) *producer {
	p := producer{
		kind: kind,
	}
	return &p
}
