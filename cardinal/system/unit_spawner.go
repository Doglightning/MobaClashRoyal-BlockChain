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

// Spawns player units
func UnitSpawnerSystem(world cardinal.WorldContext) error {
	return cardinal.EachMessage[msg.CreateUnitMsg, msg.CreateUnitResult](
		world,
		func(create cardinal.TxData[msg.CreateUnitMsg]) (msg.CreateUnitResult, error) {
			//create filter for matching ID's
			matchFilter := cardinal.ComponentFilter[comp.MatchId](func(m comp.MatchId) bool {
				return m.MatchId == create.Msg.MatchID
			})
			gameStateSearch := cardinal.NewSearch().Entity(
				filter.Exact(GameStateFilters())).
				Where(matchFilter)

			gameState, err := gameStateSearch.First(world)

			if err != nil {
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("error searching for match: %w", err)
			}

			if gameState == iterators.BadID { // Assuming cardinal.NoEntity represents no result found
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("no match found with ID or missing components: %s", create.Msg.MatchID)
			}

			//check if unit being spawned exsists in the registry above
			unitType, ok := UnitRegistry[create.Msg.UnitType]
			if !ok {
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("unit type %s not found in registry", create.Msg.UnitType)
			}

			//get UID
			tempUID, err := cardinal.GetComponent[comp.UID](world, gameState)
			if err != nil {
				return msg.CreateUnitResult{}, fmt.Errorf("error getting UID: %w", err)
			}

			//calculate distance from enemy spawn
			var tempDistance float32
			if create.Msg.Team == "Blue" {
				tempDistance = float32(math.Sqrt(((float64(create.Msg.PositionX) - float64(MapDataRegistry[create.Msg.MapName].Bases[1][0])) * (float64(create.Msg.PositionX) - float64(MapDataRegistry[create.Msg.MapName].Bases[1][0]))) + ((float64(create.Msg.PositionY) - float64(MapDataRegistry[create.Msg.MapName].Bases[1][1])) * (float64(create.Msg.PositionY) - float64(MapDataRegistry[create.Msg.MapName].Bases[1][1])))))
			} else {
				tempDistance = float32(math.Sqrt(((float64(create.Msg.PositionX) - float64(MapDataRegistry[create.Msg.MapName].Bases[0][0])) * (float64(create.Msg.PositionX) - float64(MapDataRegistry[create.Msg.MapName].Bases[0][0]))) + ((float64(create.Msg.PositionY) - float64(MapDataRegistry[create.Msg.MapName].Bases[0][1])) * (float64(create.Msg.PositionY) - float64(MapDataRegistry[create.Msg.MapName].Bases[0][1])))))
			}

			//get SpatialHash component from game state
			SpatialHash, err := cardinal.GetComponent[comp.SpatialHash](world, gameState)
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
				comp.Health{CurrentHP: unitType.Health, MaxHP: unitType.Health},
				comp.Movespeed{CurrentMS: unitType.Speed},
				comp.Position{PositionVectorX: create.Msg.PositionX, PositionVectorY: create.Msg.PositionY, PositionVectorZ: create.Msg.PositionZ, RotationVectorX: create.Msg.RotationX, RotationVectorY: create.Msg.RotationY, RotationVectorZ: create.Msg.RotationZ},
				comp.MapName{MapName: create.Msg.MapName},
				comp.Distance{Distance: tempDistance},
				comp.UnitRadius{UnitRadius: unitType.Radius},
				comp.Attack{Combat: false, Damage: unitType.Damage, Rate: unitType.AttackRate, Frame: 0, DamageFrame: unitType.DamageFrame, Class: unitType.Class, AttackRadius: unitType.AttackRadius, AggroRadius: unitType.AggroRadius},
			)
			if errr != nil {
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("error creating unit: %w", err)
			}
			//add unit to spatial hash collision map
			AddObjectSpatialHash(SpatialHash, entityID, create.Msg.PositionX, create.Msg.PositionY, unitType.Radius, create.Msg.Team)

			//incriment UID
			tempUID.UID++
			if err := cardinal.SetComponent[comp.UID](world, gameState, tempUID); err != nil {
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("error updating UID: %w", err)
			}

			return msg.CreateUnitResult{Success: true}, nil
		})
}
