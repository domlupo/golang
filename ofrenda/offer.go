package main

type offerType int

const (
	noOfferType offerType = iota
	bottle
	candy
	flower
	skull
)

type offer struct {
	offerType offerType
}

func newOffer(ot offerType) *offer {
	o := offer{
		offerType: ot,
	}
	return &o
}
