// A shape is a collection of fully connected lines. This file is hacky since it isnt foundational, performance sensitive, etc.

package main

type shape struct {
	boundary []line
}

// Performs basic validation but assumes caller passes a valid boundary. Specifically, collinear can create invalid shapes.
// Creating a line method which continually merges lines to perform validation seemed very complex once I started.
// Also, I didn't seen an obvious need for the merge method elsewhere.
func newShape(boundary []line) shape {
	uniq := make(map[line]struct{})
	for _, l := range boundary {
		uniq[l] = struct{}{}
	}
	if len(uniq) != len(boundary) {
		panic("newShape: boundary lines must be unique")
	}
	if len(boundary) < 3 {
		panic("newShape: cannot create a shape with less than 3 boundary lines")
	}

	remove := len(boundary) - 1
	line := boundary[remove]
	pixel := line.p1
	endPixel := line.p2
	for len(boundary) > 1 {
		boundary = append(boundary[:remove], boundary[remove+1:]...)
		remove = -1
		for i, l := range boundary {
			switch {
			case remove != -1:
				continue
			case pixel.equals(l.p1):
				remove = i
				pixel = l.p2
			case pixel.equals(l.p2):
				remove = i
				pixel = l.p1
			}
		}
		if remove == -1 {
			panic("newShape: did not find a boundary line to connect")
		}
	}
	if pixel != endPixel {
		panic("newShape: the last boundary line pixel connected is not equivalent to the end pixel")
	}
	return shape{boundary}
}
