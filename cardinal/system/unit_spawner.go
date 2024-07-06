package system

import (
	"fmt"
	"math"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/iterators"
	"pkg.world.dev/world-engine/cardinal/search/filter"

	comp "MobaClashRoyal/component"
	"MobaClashRoyal/msg"
)

type UnitType struct {
	Name         string
	Health       float32
	Damage       int
	AttackRate   int //tick based 5 = 5 ticks (100ms tickrate = 500ms attack rate)
	Target       int
	Speed        float32
	Cost         int
	Radius       int
	AttackRadius int
	DamageFrame  int
}

var UnitRegistry = map[string]UnitType{
	"Vampire": {Health: 100, Damage: 13, AttackRate: 10, DamageFrame: 4, Speed: 50, Cost: 50, Radius: 120, AttackRadius: 1400},
}

// Spawns player units
func UnitSpawnerSystem(world cardinal.WorldContext) error {
	return cardinal.EachMessage[msg.CreateUnitMsg, msg.CreateUnitResult](
		world,
		func(create cardinal.TxData[msg.CreateUnitMsg]) (msg.CreateUnitResult, error) {
			//create filter for matching ID's
			matchFilter := cardinal.ComponentFilter[comp.MatchId](func(m comp.MatchId) bool {
				return m.MatchId == create.Msg.MatchID
			})
			teamSearch := cardinal.NewSearch().Entity(
				filter.Exact(TeamFilters())).
				Where(matchFilter)

			found, err := teamSearch.First(world)

			if err != nil {
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("error searching for match: %w", err)
			}

			if found == iterators.BadID { // Assuming cardinal.NoEntity represents no result found
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("no match found with ID or missing components: %s", create.Msg.MatchID)
			}

			//check if unit being spawned exsists in the registry above
			unitType, ok := UnitRegistry[create.Msg.UnitType]
			if !ok {
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("unit type %s not found in registry", create.Msg.UnitType)
			}

			//get UID
			tempUID, err := cardinal.GetComponent[comp.UID](world, found)
			if err != nil {
				return msg.CreateUnitResult{}, fmt.Errorf("error getting UID: %w", err)
			}

			//search filter to find map that unit was placed on
			matchFilter = cardinal.ComponentFilter[comp.MapName](func(m comp.MapName) bool {
				return m.MapName == create.Msg.MapName
			})
			mapSearch := cardinal.NewSearch().Entity(
				filter.Exact(filter.Component[comp.DirectionMap](), filter.Component[comp.GridUtils](), filter.Component[comp.MapName]())).
				Where(matchFilter)

			foundMap, err := mapSearch.First(world)

			if err != nil {
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("error searching for map: %w", err)
			}

			if foundMap == iterators.BadID { // Assuming cardinal.NoEntity represents no result found
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("no map found with MapName or missing components: %s", create.Msg.MatchID)
			}

			//get Map Team endpoint coords
			tempGridUtils, err := cardinal.GetComponent[comp.GridUtils](world, foundMap)
			if err != nil {
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("error getting GridUtils: %w", err)
			}

			//calculate distance from enemy spawn
			var tempDistance float64
			if create.Msg.Team == "Blue" {
				tempDistance = math.Sqrt(((float64(create.Msg.PositionX) - float64(tempGridUtils.RedX)) * (float64(create.Msg.PositionX) - float64(tempGridUtils.RedX))) + ((float64(create.Msg.PositionY) - float64(tempGridUtils.RedY)) * (float64(create.Msg.PositionY) - float64(tempGridUtils.RedY))))
			} else {
				tempDistance = math.Sqrt(((float64(create.Msg.PositionX) - float64(tempGridUtils.BlueX)) * (float64(create.Msg.PositionX) - float64(tempGridUtils.BlueX))) + ((float64(create.Msg.PositionY) - float64(tempGridUtils.BlueY)) * (float64(create.Msg.PositionY) - float64(tempGridUtils.BlueY))))
			}

			//get SpatialHash component from game state
			SpatialHash, err := cardinal.GetComponent[comp.SpatialHash](world, found)
			if err != nil {
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("error getting SpatialHash component (unit Spawner): %w", err)
			}
			if CheckCollisionSpatialHash(SpatialHash, create.Msg.PositionX, create.Msg.PositionY, unitType.Radius) {
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("collision with unit")
			}

			//create unit
			entityID, errr := cardinal.Create(world,
				comp.MatchId{MatchId: create.Msg.MatchID},
				comp.UID{UID: tempUID.UID},
				comp.UnitName{UnitName: create.Msg.UnitType},
				comp.Team{Team: create.Msg.Team},
				comp.UnitHealth{CurrentHP: unitType.Health, MaxHP: unitType.Health},
				comp.Movespeed{CurrentMS: unitType.Speed},
				comp.Position{PositionVectorX: create.Msg.PositionX, PositionVectorY: create.Msg.PositionY, PositionVectorZ: create.Msg.PositionZ, RotationVectorX: create.Msg.RotationX, RotationVectorY: create.Msg.RotationY, RotationVectorZ: create.Msg.RotationZ},
				comp.MapName{MapName: create.Msg.MapName},
				comp.Distance{Distance: tempDistance},
				comp.UnitRadius{UnitRadius: unitType.Radius},
				comp.AttackRadius{AttackRadius: unitType.AttackRadius},
				comp.Attack{Combat: false, Damage: unitType.Damage, Rate: unitType.AttackRate, Frame: 0, DamageFrame: unitType.DamageFrame},
			)
			if errr != nil {
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("error creating unit: %w", err)
			}

			AddObjectSpatialHash(SpatialHash, entityID, create.Msg.PositionX, create.Msg.PositionY, unitType.Radius, create.Msg.Team)

			//incriment UID
			tempUID.UID++
			if err := cardinal.SetComponent[comp.UID](world, found, tempUID); err != nil {
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("error updating UID: %w", err)
			}

			return msg.CreateUnitResult{Success: true}, nil
		})
}
