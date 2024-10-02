package system

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/iterators"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"

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
				filter.Contains(filter.Component[comp.GameStateTag]())).
				Where(matchFilter).First(world)
			if err != nil {
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("error searching for match (unit_spawner.go): %w", err)
			}
			if gameState == iterators.BadID { // Assuming cardinal.NoEntity represents no result found
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("no match found with ID or missing components (unit_spawner.go): %s", create.Msg.MatchID)
			}

			//get unit data
			unitType, spType, err := getUnitData(create.Msg.UnitType)
			if err != nil {
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("(unit_spawner.go): %w", err)
			}

			//check if mapName exsists and if direction vector exsists at (x, y) location
			if !moveDirectionExsist(create.Msg.PositionX, create.Msg.PositionY, create.Msg.MapName) {
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("map name or direction vector does not exsist for location")
			}

			mapData, exists := MapDataRegistry[create.Msg.MapName]
			if !exists {
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("error key for MapDataRegistry does not exsist (unit_spawner.go)")
			}

			//calculate distance from enemy spawn
			var tempDistance float32
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
			if CheckCollisionSpatialHash(SpatialHash, create.Msg.PositionX, create.Msg.PositionY, unitType.Radius, unitType.Class, true) {
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("collision with unit (unit_spawner.go)")
			}

			err = handLogic(world, gameState, create.Msg.UnitType, create.Msg.Team, unitType.Cost, create.Msg.UID)
			if err != nil {
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("(unit_spawner.go) - %w", err)
			}

			//get new UID
			UID, err := getNextUID(world, create.Msg.MatchID)
			if err != nil {
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("(unit_spawner.go) - %w", err)
			}

			zOffSet := create.Msg.PositionZ
			if unitType.Class == "air" {
				zOffSet += 450
			}

			//create unit
			entityID, err := cardinal.Create(world,
				comp.MatchId{MatchId: create.Msg.MatchID},
				comp.UID{UID: UID},
				comp.UnitName{UnitName: create.Msg.UnitType},
				comp.Team{Team: create.Msg.Team},
				comp.Health{CurrentHP: unitType.Health, MaxHP: unitType.Health},
				comp.Movespeed{CurrentMS: unitType.Speed},
				comp.Position{PositionVectorX: create.Msg.PositionX, PositionVectorY: create.Msg.PositionY, PositionVectorZ: zOffSet, RotationVectorX: create.Msg.RotationX, RotationVectorY: create.Msg.RotationY, RotationVectorZ: create.Msg.RotationZ},
				comp.MapName{MapName: create.Msg.MapName},
				comp.Distance{Distance: tempDistance},
				comp.Class{Class: unitType.Class},
				//comp.Destroyed{Destroyed: false},
				comp.UnitRadius{UnitRadius: unitType.Radius},
				comp.Attack{
					Combat:       false,
					Damage:       unitType.Damage,
					Rate:         unitType.AttackRate,
					Frame:        0,
					DamageFrame:  unitType.DamageFrame,
					AttackRadius: unitType.AttackRadius,
					AggroRadius:  unitType.AggroRadius,
					State:        "Default",
				},
				comp.Sp{
					DmgSp:               unitType.DmgSp,
					SpRate:              unitType.SpRate,
					CurrentSp:           unitType.CurrentSP,
					MaxSp:               unitType.MaxSP,
					Charged:             false,
					Rate:                spType.AttackRate,
					DamageFrame:         spType.DamageFrame,
					DamageEndFrame:      spType.DamageEndFrame,
					StructureTargetable: spType.StructureTargetable,
				},
				comp.CenterOffset{CenterOffset: unitType.CenterOffset},
				comp.CC{Stun: 0},
				comp.EffectsList{EffectsList: make(map[string]int)},
				comp.UnitTag{},
			)
			if err != nil {
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("error creating unit (unit_spawner.go): %w", err)
			}

			//add unit to collision hash collision map
			AddObjectSpatialHash(SpatialHash, entityID, create.Msg.PositionX, create.Msg.PositionY, unitType.Radius, create.Msg.Team, unitType.Class)

			err = cardinal.SetComponent(world, gameState, SpatialHash)
			if err != nil {
				return msg.CreateUnitResult{Success: false}, fmt.Errorf("error setting hash component (unit_spawner.go): %w", err)
			}

			return msg.CreateUnitResult{Success: true}, nil
		})
}

// Deals with the logic of playing a card from hand and drawing from deck to replace
func handLogic(world cardinal.WorldContext, gameState types.EntityID, name, team string, cost, UID int) error {

	var found bool = false

	if team == "Blue" {
		//get player1 component from game state
		player1, err := cardinal.GetComponent[comp.Player1](world, gameState)
		if err != nil {
			return fmt.Errorf("error getting player1 component (unit_spawner.go): %w", err)
		}
		//check if enough gold to spawn unit
		if player1.Gold < float32(cost) {
			return fmt.Errorf("not enough gold to spawn %s (unit_spawner.go): ", name)
		}
		//check unit spawned is in hand
		for _, v := range player1.Hand { //search hand
			if v == name { // if card == unit spawned
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("card not in hand (unit_spawner.go) ")
		}

		//reduce Gold
		player1.Gold -= float32(cost)
		//-1 key means to tell player to remove unit they were holding to transition it to this spawned unit
		player1.RemovalList[UID] = true

		//hand sorting
		tempCard := player1.Deck[0]                     //get top deck card
		player1.Deck = removeFirstElement(player1.Deck) //remove top deck card
		for i, v := range player1.Hand {                //search hand
			if v == name { // if card == unit spawned
				player1.Hand[i] = tempCard             //insert top card to hand
				player1.Deck = append(player1.Deck, v) // put spawned unit to back of deck
			}
		}

		err = cardinal.SetComponent(world, gameState, player1)
		if err != nil {
			return fmt.Errorf("error setting player1 component (unit_spawner.go): %w", err)
		}
	} else {
		// get player2 component from game state
		player2, err := cardinal.GetComponent[comp.Player2](world, gameState)
		if err != nil {
			return fmt.Errorf("error getting player2 component (unit_spawner.go): %w", err)
		}
		//check if enough gold to spawn unit
		if player2.Gold < float32(cost) {
			return fmt.Errorf("not enough gold to spawn %s (unit_spawner.go): ", name)
		}
		//check unit spawned is in hand
		for _, v := range player2.Hand { //search hand
			if v == name { // if card == unit spawned
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("card not in hand (unit_spawner.go) ")
		}

		//reduce Gold
		player2.Gold -= float32(cost)
		//-1 key means to tell player to remove unit they were holding to transition it to this spawned unit
		player2.RemovalList[UID] = true

		//hand sorting
		tempCard := player2.Deck[0]                     //get top deck card
		player2.Deck = removeFirstElement(player2.Deck) //remove top deck card
		for i, v := range player2.Hand {                //search hand
			if v == name { // if card == unit spawned
				player2.Hand[i] = tempCard             //insert top card to hand
				player2.Deck = append(player2.Deck, v) // put spawned unit to back of deck
			}
		}

		err = cardinal.SetComponent(world, gameState, player2)
		if err != nil {
			return fmt.Errorf("error setting player2 component (unit_spawner.go): %w", err)
		}
	}

	return nil
}

func removeFirstElement(slice []string) []string {
	if len(slice) > 0 {
		return slice[1:] // Slice from the second element to the end
	}
	return slice // Return the original slice if it's empty
}
