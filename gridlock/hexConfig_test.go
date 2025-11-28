package main

import (
	"fmt"
	"testing"
)

func TestDemotifHexLine(t *testing.T) {
	t.Parallel()
	hc := newHexConfig(newWindowConfig(), 46, 40, 19, 10)

	nw := Pixel{0, 10}
	n := Pixel{23, 0}
	ne := Pixel{46, 10}

	sw := Pixel{0, 30}
	s := Pixel{23, 40}
	se := Pixel{46, 30}

	var tests = []struct {
		hs        hexSide
		p1        Pixel
		p2        Pixel
		wantPanic bool
	}{
		{hsNone, Pixel{}, Pixel{}, true},
		{hsWest, nw, sw, false},
		{hsNorthwest, n, nw, false},
		{hsNortheast, n, ne, false},
		{hsEast, ne, se, false},
		{hsSouthwest, s, sw, false},
		{hsSoutheast, s, se, false},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("test demotif hex line for hex side %s", tt.hs.string())
		t.Run(testname, func(t *testing.T) {
			defer func() {
				r := recover()
				switch {
				case r == nil && tt.wantPanic:
					t.Errorf("want panic but did not panic")
				case r != nil && !tt.wantPanic:
					t.Errorf("do not want panic but did panic")
				}
			}()
			line := hc.demotifHexLine(tt.hs)
			if !(line.p1.equals(tt.p1) || line.p1.equals(tt.p2)) {
				t.Errorf("demotifHexLine has point %s rather than expected points %s %s", line.p1.log(), tt.p1.log(), tt.p2.log())
			}
			if !(line.p2.equals(tt.p1) || line.p2.equals(tt.p2)) {
				t.Errorf("demotifHexLine has point %s rather than expected points %s %s", line.p2.log(), tt.p1.log(), tt.p2.log())
			}
		})
	}
}

func TestDemotifPixels(t *testing.T) {
	hc := newHexConfig(newWindowConfig(), 46, 40, 19, 10)
	demotifPixels := hc.demotifHexPixels()
	if len(demotifPixels) != 1380 {
		t.Error("demotifPixels sanity check failed")
	}
}
