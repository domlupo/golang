// Provides spatial partitions for different unit and analysis combinations.

package main

import (
	"fmt"
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
	infantryUC := ucProvider.provide(infantry)
	p.infantryCollisionLUT = p.setupPartitionLUT(infantryUC.radius, infantryUC.radius)
	p.infantryVisionLUT = p.setupPartitionLUT(infantryUC.vision, infantryUC.radius)
	return p
}

func (part partitioner) partitions(d dynamic, s static) partitions {
	if tiles, ok := part.lookUpPartitionsFromTable(d, s); ok {
		return tiles
	}
	tiles := part.defaultPartitions(d)
	var parts partitions
	for _, t := range tiles {
		parts = append(parts, part.hexagons[t])
	}
	return parts
}

func (part partitioner) lookUpPartitionsFromTable(d dynamic, s static) (partitions, bool) {
	var parts partitions
	var lut partitionLUT
	infantryUC := part.unitConfigProvider.provide(infantry)
	switch {
	case d.circleRadius != nil && *d.circleRadius == infantryUC.radius && s.circleRadius != nil &&
		*s.circleRadius == infantryUC.radius && s.offset.equals(infantryUC.staticBody.offset):
		lut = part.infantryCollisionLUT
	case d.circleRadius != nil && *d.circleRadius == infantryUC.vision && s.circleRadius != nil &&
		*s.circleRadius == infantryUC.radius && s.offset.equals(infantryUC.staticBody.offset):
		lut = part.infantryVisionLUT
	default:
		return parts, false
	}

	demotif, demotifOk := part.grid.demotif(d.Pos.pixel())
	if !demotifOk {
		panic(fmt.Sprintf(
			"partitioner.lookUpPartitionsFromTable: could not demotif pixel %s",
			d.Pos.pixel().log()))
	}
	hexSides := lut[demotif]
	tile, tileOk := part.grid.getTile(d.Pos.pixel())
	tiles := []Tile{tile}
	if !tileOk {
		panic(fmt.Sprintf(
			"partitioner.lookUpPartitionsFromTable: could not get Tile for Pixel %s",
			d.Pos.pixel().log()))
	}
	for _, hs := range hexSides {
		if t, ok := part.grid.adjacent(tile, hs); ok {
			tiles = append(tiles, t)
		}
	}
	for _, t := range tiles {
		parts = append(parts, part.hexagons[t])
	}
	return parts, true
}

// Returns conservative partitions which collide with the passed Pos and distance. The caller is
// responsible for verifying the analysis performed with the partitions is valid. For example,
// analyzing partitions for collision detection between two circles would be invalid if the circle
// with the Pos and distance not passed has a radius which exceeds the conservative hex radius. Not
// optimized because the hexagon collision uses conservative circles.
func (part partitioner) defaultPartitions(d dynamic) []Tile {
	conservativeStatic := newStaticCircle(part.grid.conservativeHexPixelRadius, Pos{})
	p := d.Pos
	t, ok := part.grid.getTile(p.pixel())
	if !ok {
		panic(fmt.Sprintf("partitioner.defaultPartitions: cannot find tile for %s", p.pixel().log()))
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
				collides := d.collides(*conservativeStatic.newDynamic(pa.pos()))
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

// does not work for polygons but may be easy to extend. Also this method signature can probably be
// improved by using partitionerRequest or something else but it is low priority.
func (part partitioner) setupPartitionLUT(
	internalRadius int,
	externalRadius int) partitionLUT {

	lut := make(partitionLUT)
	collisionSources := part.grid.demotifHexPixelsAdjacent()
	pixels := part.grid.demotifHexPixels()
	for externalP, hs := range collisionSources {
		for _, internalP := range pixels {
			d1 := newStaticCircle(internalRadius, Pos{}).newDynamic(internalP.pos())
			d2 := newStaticCircle(externalRadius, Pos{}).newDynamic(externalP.pos())
			if d1.collides(*d2) {
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
