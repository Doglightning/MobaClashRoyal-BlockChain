package system

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/iterators"
	"pkg.world.dev/world-engine/cardinal/search/filter"

	comp "MobaClashRoyal/component"
	"MobaClashRoyal/msg"
)

// Spawns player units
// called from create_unit.go msg
func UnitSpawnerSystem(world cardinal.WorldContext) error {
	return cardinal.EachMessage(world,
		func(create cardinal.TxData[msg.CreateUnitMsg]) (msg.CreateUnitResult, error) {
			//create filter for matching ID's
			matchFilter := cardinal.ComponentFilter[comp.MatchId](func(m comp.MatchId) bool {
				return m.MatchId == create.Msg.MatchID
			})
			//get game state
			gameState, err := cardinal.NewSearch().Entity(
				filter.Exact(GameStateFilters())).
				Where(matchFilter).First(world)
			if err != nil {
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("error searching for match (unit_spawner.go): %w", err)
			}
			if gameState == iterators.BadID { // Assuming cardinal.NoEntity represents no result found
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("no match found with ID or missing components (unit_spawner.go): %s", create.Msg.MatchID)
			}

			//check if unit being spawned exsists in the unit registry
			unitType, ok := UnitRegistry[create.Msg.UnitType]
			if !ok {
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("unit type %s not found in registry (unit_spawner.go)", create.Msg.UnitType)
			}

			//get new UID
			UID, err := getNextUID(world, create.Msg.MatchID)
			if err != nil {
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("(unit_spawner.go) - %w", err)
			}

			//calculate distance from enemy spawn
			var tempDistance float32

			//check map exsists in registy
			mapData, exists := MapDataRegistry[create.Msg.MapName]
			if !exists {
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("error key for MapDataRegistry does not exsist (unit_spawner.go)")
			}
			if create.Msg.Team == "Blue" {
				tempDistance = distanceBetweenTwoPoints(float32(mapData.Bases[1][0]), float32(mapData.Bases[1][1]), create.Msg.PositionX, create.Msg.PositionY)
			} else {
				tempDistance = distanceBetweenTwoPoints(float32(mapData.Bases[0][0]), float32(mapData.Bases[0][1]), create.Msg.PositionX, create.Msg.PositionY)
			}

			//get collision Hash component from game state
			SpatialHash, err := cardinal.GetComponent[comp.SpatialHash](world, gameState)
			if err != nil {
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("error getting SpatialHash component (unit_spawner.go): %w", err)
			}
			//check if spawning on a taken spot in collision hash
			if CheckCollisionSpatialHash(SpatialHash, create.Msg.PositionX, create.Msg.PositionY, unitType.Radius) {
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("collision with unit (unit_spawner.go)")
			}

			//create unit
			entityID, err := cardinal.Create(world,
				comp.MatchId{MatchId: create.Msg.MatchID},
				comp.UID{UID: UID},
				comp.UnitName{UnitName: create.Msg.UnitType},
				comp.Team{Team: create.Msg.Team},
				comp.Health{CurrentHP: unitType.Health, MaxHP: unitType.Health},
				comp.Movespeed{CurrentMS: unitType.Speed},
				comp.Position{PositionVectorX: create.Msg.PositionX, PositionVectorY: create.Msg.PositionY, PositionVectorZ: create.Msg.PositionZ, RotationVectorX: create.Msg.RotationX, RotationVectorY: create.Msg.RotationY, RotationVectorZ: create.Msg.RotationZ},
				comp.MapName{MapName: create.Msg.MapName},
				comp.Distance{Distance: tempDistance},
				comp.UnitRadius{UnitRadius: unitType.Radius},
				comp.Attack{Combat: false, Damage: unitType.Damage, Rate: unitType.AttackRate, Frame: 0, DamageFrame: unitType.DamageFrame, Class: unitType.Class, AttackRadius: unitType.AttackRadius, AggroRadius: unitType.AggroRadius},
				comp.Sp{DmgSp: unitType.DmgSp, SpRate: unitType.SpRate, CurrentSp: unitType.CurrentSP, MaxSp: unitType.MaxSP, Animation: "default"},
			)
			if err != nil {
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("error creating unit (unit_spawner.go): %w", err)
			}

			//add unit to collision hash collision map
			AddObjectSpatialHash(SpatialHash, entityID, create.Msg.PositionX, create.Msg.PositionY, unitType.Radius, create.Msg.Team)

			return msg.CreateUnitResult{Success: true}, nil
		})
}
