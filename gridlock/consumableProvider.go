// See unitConfigProvider.go file for documentation. This is a separate file for readability.

package main

import (
	"fmt"
	"log/slog"
)

type consumableProvider struct {
	consumables
}

type consumables map[name]consumable

type consumable struct {
	count         int
	countMax      int
	fireDelay     step
	firesAt       func(uc unitConfig) bool
	firesAtLayers layers
}

func consumableNames() []name {
	return []name{bombGravity, missileHeatSeeking, missileRadarActive, missileBallistic,
		missileCruise, missileAntiBallistic, flare, chaff}
}

func newConsumableProvider(unitConfigs unitConfigs) consumableProvider {
	consumables := make(consumables)
	for _, name := range consumableNames() {
		consumables[name] = consumable{}
	}
	for name, uc := range consumables {
		uc.setFireDelay(name)
		uc.setFiresAt(name)
		consumables[name] = uc
	}
	cProvider := consumableProvider{consumables}
	cProvider.setFiresAtLayers(unitConfigs)
	return cProvider
}

func (cProvider consumableProvider) provide(name name, countMax int) consumable {
	c := cProvider.consumables[name]
	c.countMax = countMax
	return c
}

func (c *consumable) setFireDelay(name name) {
	var fireDelay step
	switch name {
	case missileHeatSeeking:
		fireDelay = 50
	case missileRadarActive:
		fireDelay = 50
	case flare:
		fireDelay = 50
	case chaff:
		fireDelay = 50
	case bombGravity, missileCruise, missileBallistic, missileAntiBallistic:
	default:
		panic(fmt.Sprintf("consumable.setFireDelay: unknown consumable %s", name.string()))
	}
	c.fireDelay = fireDelay
}

func (c *consumable) setFiresAt(name name) {
	firesAt := func(uc unitConfig) bool { return false }
	switch name {
	case missileHeatSeeking:
		firesAt = func(uc unitConfig) bool { return uc.isEmitsHeat && uc.name != flare }
	case missileRadarActive:
		firesAt = func(uc unitConfig) bool { return uc.isEmitsRadar && uc.name != chaff }
	case bombGravity, flare, chaff, missileBallistic, missileAntiBallistic, missileCruise:
	default:
		panic(fmt.Sprintf("consumable.setFiresAt: unknown consumable %s", name.string()))
	}
	c.firesAt = firesAt
}

func (cProvider consumableProvider) setFiresAtLayers(unitConfigs unitConfigs) {
	layers := []layer{}
	for name, c := range cProvider.consumables {
		for _, uc := range unitConfigs {
			if c.firesAt(uc) {
				layers = append(layers, uc.produceLayer)
			}
		}
		c.firesAtLayers = layers
		cProvider.consumables[name] = c
	}
}

// Logs consumables in a deterministic order.
func (c consumables) log() []slog.Attr {
	a := []slog.Attr{}
	if len(c) == 0 {
		a = append(a, slog.String("consumables", "none"))
	}
	for _, name := range consumableNames() {
		if c, ok := c[name]; ok {
			a = append(a, nest(name.string(), c.log()))
		}
	}
	return a
}

// Typical log method.
func (c consumable) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, slog.Int("count", c.count))
	a = append(a, slog.Int("countMax", c.countMax))
	a = append(a, slog.Uint64("fireDelay", c.fireDelay))
	a = append(a, slog.String("firesAt", c.firesAtLayers.string()))
	return a
}
