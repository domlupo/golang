// Misc code which is used in several other files but not large enough to need its own file.

package main

type player int

const (
	noPlayer player = iota
	bluePlayer
	redPlayer
	neutralPlayer
)
