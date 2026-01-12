// Provides spatial partitions for different unit and analysis combinations.

package main

import (
	"fmt"
	"log/slog"
	"slices"
)

type partitions []*hexagon

// The type of analysis which will process the partitions.
type partitionAnalysis int

// A look up table (LUT) to find an optimized number of partitions. It avoids runtime collision in
// two ways. First, collision checks for all 1380 Pixels in the demotif hexagon are done at startup
// for the 132 adjacent Pixels. If a collision occurs the hexSide which caused it is stored. Second,
// optimized number of partitions means no extra partitions are considered for unit <-> unit
// collisions since spatial partitioning checks all units within the partition. The LUT is utilized
// via: Unit Pixel -> demotif -> hex sides from LUT -> check unit tile for collision -> check hex
// sides from LUT for collision.
type partitionLUT map[demotif][]hexSide

const (
	paNone partitionAnalysis = iota
	paCollision
	paVision
	paSteer
)

func (pa partitionAnalysis) string() string {
	switch pa {
	case paNone:
		return "paNone"
	case paCollision:
		return "paCollision"
	case paVision:
		return "paVision"
	case paSteer:
		return "paSteer"
	default:
		panic(fmt.Sprintf("partitionAnalysis.string unknown partitionAnalysis: %d", pa))
	}
}

// Typical log method.
func (pa partitionAnalysis) log() slog.Attr {
	return slog.String("partitionAnalysis", pa.string())
}

type partitioner struct {
	grid                 grid
	hexagons             hexagons
	unitConfigProvider   unitConfigProvider
	infantryCollisionLUT partitionLUT
	infantryVisionLUT    partitionLUT
}

func newPartitioner(g grid, h hexagons, ucProvider unitConfigProvider) partitioner {
	p := partitioner{
		grid:               g,
		hexagons:           h,
		unitConfigProvider: ucProvider,
	}
	p.infantryCollisionLUT = p.setupPartitionLUT(infantry, infantry, paCollision)
	p.infantryVisionLUT = p.setupPartitionLUT(infantry, infantry, paVision)
	return p
}

func (part partitioner) partitions(u unit, pa partitionAnalysis, buffer int) partitions {
	var tiles []Tile
	radius := buffer
	switch {
	case u.name == infantry && pa != paSteer:
		return part.lookUpPartitions(u, pa)
	case pa == paCollision:
		radius += u.radius
	case pa == paVision:
		radius += u.vision
	case pa == paSteer:
		radius += 0
	default:
		panic(fmt.Sprintf("partitioner.partitions: unknown partition analysis %s", pa.string()))
	}

	if pa == paSteer && u.name != infantry {
		panic(fmt.Sprintf(
			"partitioner.partitions: cannot get partitions for unit %s, analysis %s and buffer %d",
			u.name.string(),
			pa.string(),
			buffer))
	}

	tiles = part.defaultPartitions(u.Pos, radius)
	var parts partitions
	for _, t := range tiles {
		parts = append(parts, part.hexagons[t])
	}
	return parts
}

func (part partitioner) lookUpPartitions(u unit, pa partitionAnalysis) partitions {
	var lut partitionLUT
	switch pa {
	case paCollision:
		lut = part.infantryCollisionLUT
	case paVision:
		lut = part.infantryVisionLUT
	default:
		panic(fmt.Sprintf("partitioner.lookUpPartitions: unknown partition analysis %s", pa.string()))
	}

	demotif, demotifOk := part.grid.demotif(u.pixel())
	if !demotifOk {
		panic(fmt.Sprintf("partitioner.lookUpPartitions: could not demotif pixel %s", u.pixel().log()))
	}
	hexSides := lut[demotif]
	tile, tileOk := part.grid.getTile(u.pixel())
	tiles := []Tile{tile}
	if !tileOk {
		panic(fmt.Sprintf(
			"partitioner.lookUpPartitions: could not get Tile for Pixel %s",
			u.pixel().log()))
	}
	for _, hs := range hexSides {
		if t, ok := part.grid.adjacent(tile, hs); ok {
			tiles = append(tiles, t)
		}
	}
	var parts partitions
	for _, t := range tiles {
		parts = append(parts, part.hexagons[t])
	}
	return parts
}

// Returns conservative partitions which collide with the passed Pos and distance. The caller is
// responsible for verifying the analysis performed with the partitions is valid. For example,
// analyzing partitions for collision detection between two circles would be invalid if the circle
// with the Pos and distance not passed has a radius which exceeds the conservative hex radius. Not
// optimized because the hexagon collision uses conservative circles.
func (part partitioner) defaultPartitions(p Pos, d int) []Tile {
	t, ok := part.grid.getTile(p.pixel())
	if !ok {
		panic(fmt.Sprintf("partitioner.defaultPartitions: cannot find tile for %s", p.log()))
	}
	visited := make(map[Tile]struct{})
	visited[t] = struct{}{}
	q := []Tile{t}
	nextQ := []Tile{}
	for len(q) > 0 {
		current := q[0]
		q = q[1:]

		for _, aT := range part.grid.allAdjacent(current) {
			if _, ok := visited[aT]; !ok {
				visited[aT] = struct{}{}
				// TODO: Fix this. This isn't correct but will probably work fine for awhile.
				// It doesn't handle border partition collisions 100% correct but this should be rare.
				// Also, I suspect this partitions method will most likely be refactored before issue
				// blocks development.
				pa, okA := part.grid.getCenterWithFallback(aT)
				if !okA {
					panic(fmt.Sprintf(
						"partitioner.defaultPartitions: cannot find center for tile %s",
						aT.log()))
				}
				_, collides := collides(p.noLayer(), pa.pos(), d, part.grid.conservativeHexPixelRadius)
				if collides {
					nextQ = append(nextQ, aT)
				}
			}
		}

		if len(q) == 0 && len(nextQ) == 0 {
			break
		} else if len(q) == 0 {
			q = nextQ
			nextQ = []Tile{}
		}
	}

	partitions := make([]Tile, len(visited))
	i := 0
	for p := range visited {
		partitions[i] = p
		i++
	}
	return partitions
}

func (part partitioner) setupPartitionLUT(
	internalUnitName name,
	externalUnitName name,
	pa partitionAnalysis) partitionLUT {
	var internalRadius int
	var externalRadius int
	internalUC := part.unitConfigProvider.provide(internalUnitName)
	externalUC := part.unitConfigProvider.provide(externalUnitName)
	switch pa {
	case paCollision:
		internalRadius = internalUC.radius
		externalRadius = externalUC.radius
	case paVision:
		internalRadius = internalUC.vision
		externalRadius = externalUC.radius
	default:
		panic(fmt.Sprintf(
			"partitioner.setupPartitionsLUT: unsupported partition analysis %s",
			pa.string()))
	}

	lut := make(partitionLUT)
	collisionSources := part.grid.demotifHexPixelsAdjacent()
	pixels := part.grid.demotifHexPixels()
	for externalP, hs := range collisionSources {
		for _, internalP := range pixels {
			if _, ok := collides(internalP.pos(), externalP.pos(), internalRadius, externalRadius); ok {
				hexSides := lut[demotif{internalP}]
				if !slices.Contains(hexSides, hs) {
					hexSides = append(hexSides, hs)
					lut[demotif{internalP}] = hexSides
				}
			}
		}
	}
	return lut
}
