// A unit config provider simplifies the complex process of creating unit configs. It has several
// benefits but the core benefit is handling unit config fields dependent on other unit config
// fields.
//
// Benefits:
// - Unit config fields can have respective verifications to ensure valid values (e.g. radius
//   cannot be negative).
// - Unit config fields which are dependent on other fields within the same unit are easy to set.
//   - The order for every set field is controlled. The dependent field is easily set after its
//     required field.
//   - The dependent field can have respective verifications to ensure valid values (e.g. vision >=
//     radius)
// - Unit config fields which are dependent on other unit configs are easy to set.
//   - All independent fields, for all units can be set before the fields dependent on other units
//     are set.
// - It simplifies file exports.
//   - Consumers only interact with the single provider and the relatively simple provides method.
//   - Does not export multiple global unit configs
//   - In the future, I may make this a separate package, but for now using "set" methods clarifies
//     exports.
// - Default values are more flexible.
//   - Implicit values can be set per field via a switch statement.
//   - Required explicit values can be set per field via a switch statement.
// - Better documentation and easier to tune configurations
//   - it is easy to understand a field when it is set for every unit within a single method.
// - It is easier to extend in the future since most likely upgrades will require dynamic unit
//   configs.
//   - Requests to change configs will all be processed by the config provider.
//   - Multiple approaches are viable for consuming dynamic configs:
//     - Consumers may use the unit config provider as the single source of truth
//     - The provider may be the single source to provide updates to all consumers.
//
// Dependent fields have two anticipated sources:
// - The first source is sometimes it is easier to understand a config when it is defined via field
//   X. But the config
//   must be processed via field Y. e.g. It is easier to understand unit A attacks unit B via an
//   attack field. Then this
//   attack field is used to determine what layers the unit should search for a unit to attack.
// - The second source is when a field has high CPU usage to calculate at runtime. It is better for
//   provider to
//   calculate and cache the value at startup, and retrieve it during steady state runtime.

package main

import (
	"fmt"
)

// Provides unit configs to various consumers
type unitConfigProvider struct {
	unitConfigs
}

type unitConfigs map[name]unitConfig

// A unit action which should occur when a condition is met.
type reaction struct {
	when  func(*unit) bool
	do    func(*unit)
	next  uint64
	delay uint64
}

func (ucProvider unitConfigProvider) provide(name name) unitConfig {
	return ucProvider.unitConfigs[name]
}

func newUnitConfigProvider() unitConfigProvider {
	configs := make(unitConfigs)
	for _, name := range names() {
		configs[name] = unitConfig{name: name}
	}
	for k, uc := range configs {
		uc.setUCFields()
		configs[k] = uc
	}
	ucProvider := unitConfigProvider{configs}
	ucProvider.setUCFieldsDependentOnOtherUCs()
	return ucProvider
}

func (ucProvider unitConfigProvider) setUCFieldsDependentOnOtherUCs() {
	ucProvider.setTracksLayers()
	ucProvider.setConsumables()
	ucProvider.setReactions()
	ucProvider.setCollisionLayerStatic()
}

func (ucProvider unitConfigProvider) setCollisionLayerStatic() {
	for _, uc := range ucProvider.unitConfigs {
		maxConservativeRadius := 0
		switch uc.name {
		case infantry, fighter, bomber, missileCruise, missileBallistic, missileAntiBallistic,
			missileHeatSeeking, missileRadarActive, flare, bombGravity, chaff:
			for _, other := range ucProvider.unitConfigs {
				if uc.moveLayer == other.moveLayer {
					maxConservativeRadius = max(maxConservativeRadius,
						other.staticBody.conservativeRadius())
				}
			}
		default:
			panic(fmt.Sprintf("unitConfigProvider.setCollisionLayerStatic: unknown unit %s", uc.name.string()))
		}
		uc.collisionLayerStatic = static{circleRadius: &maxConservativeRadius}
		ucProvider.unitConfigs[uc.name] = uc
	}
}

func (ucProvider unitConfigProvider) setTracksLayers() {
	for name, uc := range ucProvider.unitConfigs {
		layers := []layer{}
		if uc.tracks != nil {
			for _, other := range ucProvider.unitConfigs {
				if uc.tracks(other) {
					layers = append(layers, other.moveLayer)
				}
			}
		}
		uc.trackLayers = layers
		ucProvider.unitConfigs[name] = uc
	}
}

func (ucProvider unitConfigProvider) setConsumables() {
	cProvider := newConsumableProvider(ucProvider.unitConfigs)
	for _, uc := range ucProvider.unitConfigs {
		counts := make(map[name]int)
		switch uc.name {
		case fighter:
			counts[bombGravity] = 1
			counts[flare] = 100
			counts[chaff] = 100
			counts[missileHeatSeeking] = 100
			counts[missileRadarActive] = 100
		case bomber:
			counts[bombGravity] = 1
			counts[flare] = 1
			counts[chaff] = 1
		case infantry, missileCruise, missileBallistic, missileAntiBallistic, missileHeatSeeking,
			missileRadarActive, flare, bombGravity, chaff:
		default:
			panic(fmt.Sprintf("unitConfigProvider.setConsumables: unknown unit %s", uc.name.string()))
		}
		uc.consumables = make(consumables)
		for n, c := range counts {
			uc.consumables[n] = cProvider.provide(n, c)
		}
		ucProvider.unitConfigs[uc.name] = uc
	}
}

func (ucProvider unitConfigProvider) setReactions() {
	for name := range ucProvider.unitConfigs {
		var reactions []reaction
		switch name {
		case infantry, bombGravity, chaff, missileRadarActive, missileHeatSeeking, missileCruise,
			missileBallistic, missileAntiBallistic, flare:
		case bomber:
			bomberUC := ucProvider.unitConfigs[bomber]
			reactions = []reaction{
				{
					when: func(bomber *unit) bool {
						// TODO: Bresenham is not guaranteed to get to final position. I have not seen
						// aircraft not reach their final destination but missiles have. Need to ensure
						// this is not error prone or make it more conservative.
						return bomber.Pos.equals(bomber.sourcePosition) && bomber.isPayloadNotFull()
					},
					do: func(bomber *unit) {
						bomber.refillPayload()
					},
					delay: bomberUC.reloadPayloadDelay,
				},
				ucProvider.reactFlare(bomberUC.consumables[flare].fireDelay),
				ucProvider.reactChaff(bomberUC.consumables[chaff].fireDelay),
			}
		case fighter:
			fighterUC := ucProvider.unitConfigs[bomber]
			missileRadarActiveUC := ucProvider.unitConfigs[missileRadarActive]
			missileHeatSeekingUC := ucProvider.unitConfigs[missileHeatSeeking]
			reactions = []reaction{
				{
					when: func(fighter *unit) bool {
						// TODO: Bresenham is not guaranteed to get to final position. I have not seen
						// aircraft not reach their final destination but missiles have. Need to ensure
						// this is not error prone or make it more conservative.
						return fighter.Pos.equals(fighter.sourcePosition) && fighter.isPayloadNotFull()
					},
					do: func(fighter *unit) {
						fighter.refillPayload()
						// TODO: when a missile terminates, should it remove the target from skip targets
						// rather than do this during restock. Not sure if this is a realistic upgrade.
						// Also, multiple ordnance might make this complex. Just because a single ordnance skips
						// a target does not mean others should or can this be avoided?
						fighter.skipTargetIds = []id{}
					},
					delay: fighterUC.reloadPayloadDelay,
				},
				{
					when: func(fighter *unit) bool { return fighter.nearest != nil },
					do: func(fighter *unit) {
						var uc unitConfig
						// TODO: determine fighter missile preference? Configurable by player?
						// Arbitrarily set radar to have higher priority than heat seeking.
						switch {
						case fighter.consumables[missileRadarActive].count > 0 &&
							!missileRadarActiveUC.ignoreTarget(fighter.nearest):
							uc = missileRadarActiveUC
						case fighter.consumables[missileHeatSeeking].count > 0 &&
							!missileHeatSeekingUC.ignoreTarget(fighter.nearest):
							uc = missileHeatSeekingUC
						default:
							return
						}
						missiles := fighter.consumables[uc.name]
						missiles.count -= 1
						fighter.consumables[uc.name] = missiles
						spawn := newUnit(
							uc,
							fighter.layer(uc.moveLayer),
							fighter.player,
							[]Pixel{fighter.nearest.pixel()},
							[]action{},
							Tile{},
							fighter.selected)
						spawn.sourceUnit = fighter.id
						fighter.requests = append(fighter.requests, request{spawn: spawn})
					},
					// TODO: fix this. delay needs to be conditional on some unit config.
					delay: 50,
				},
				ucProvider.reactChaff(fighterUC.consumables[chaff].fireDelay),
				ucProvider.reactFlare(fighterUC.consumables[flare].fireDelay),
			}
		default:
			panic(fmt.Sprintf("unitConfigProvider.setReactions: unknown unit %s", name.string()))
		}
		uc := ucProvider.unitConfigs[name]
		uc.reactions = reactions
		ucProvider.unitConfigs[name] = uc
	}
}

func (ucProvider unitConfigProvider) reactChaff(fireOrdnanceDelay uint64) reaction {
	chaffUC := ucProvider.unitConfigs[chaff]
	return reaction{
		when: func(u *unit) bool {
			return u.nearest != nil &&
				u.nearest.name == missileRadarActive && u.consumables[chaff].count > 0
		},
		do: func(u *unit) {
			newChafff := newUnit(
				chaffUC,
				u.Pos.layer(chaffUC.moveLayer),
				u.player,
				[]Pixel{u.pixel()},
				[]action{},
				Tile{},
				u.selected)
			u.requests = append(u.requests, request{spawn: newChafff})
			consumable := u.consumables[chaff]
			consumable.count -= 1
			u.consumables[chaff] = consumable
		},
		delay: fireOrdnanceDelay,
	}
}

func (ucProvider unitConfigProvider) reactFlare(fireOrdnanceDelay uint64) reaction {
	flareUC := ucProvider.unitConfigs[flare]
	return reaction{
		when: func(u *unit) bool {
			return u.nearest != nil &&
				u.nearest.name == missileHeatSeeking && u.consumables[flare].count > 0
		},
		do: func(u *unit) {
			newFlare := newUnit(
				flareUC,
				u.layer(flareUC.moveLayer),
				u.player,
				[]Pixel{u.nearest.pixel()},
				[]action{},
				Tile{},
				u.selected)
			u.requests = append(u.requests, request{spawn: newFlare})
			consumable := u.consumables[flare]
			consumable.count -= 1
			u.consumables[flare] = consumable
		},
		delay: fireOrdnanceDelay,
	}
}
