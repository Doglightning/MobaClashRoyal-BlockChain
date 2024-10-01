package system

// import (
// 	comp "MobaClashRoyal/component"
// 	"fmt"

// 	"pkg.world.dev/world-engine/cardinal"
// 	"pkg.world.dev/world-engine/cardinal/search/filter"
// 	"pkg.world.dev/world-engine/cardinal/types"
// )

// // destroy units with no health
// func UnitDestroyerSystem(world cardinal.WorldContext) error {
// 	// Filter for no HP
// 	healthFilter := cardinal.ComponentFilter(func(m comp.Health) bool {
// 		return m.CurrentHP <= 0
// 	})
// 	//for each unit with no hp's ids
// 	err := cardinal.NewSearch().Entity(
// 		filter.Contains(filter.Component[comp.UnitTag]())).
// 		Where(healthFilter).Each(world, func(id types.EntityID) bool {

// 		//get needed compoenents
// 		MatchID, uid, UnitPosition, UnitRadius, err := GetComponents4[comp.MatchId, comp.UID, comp.Position, comp.UnitRadius](world, id)
// 		if err != nil {
// 			fmt.Printf("4 (unit_destroyer.go): %v \n", err)
// 			return false
// 		}

// 		//get game state
// 		gameState, err := getGameStateGSS(world, MatchID)
// 		if err != nil {
// 			fmt.Printf("(unit_destroyer.go) %v \n", err)
// 			return false
// 		}
// 		//get player components
// 		p1, p2, err := getPlayerComponentsGSS(world, gameState)
// 		if err != nil {
// 			fmt.Printf("(unit_destroyer.go) %v \n", err)
// 			return false
// 		}

// 		//filter for units targeting self
// 		targetFilter := cardinal.ComponentFilter(func(m comp.Attack) bool {
// 			return m.Target == id
// 		})

// 		//for units targetting self, reset combat
// 		err = resetUnitsTargetingSelfUD(world, targetFilter)
// 		if err != nil {
// 			fmt.Printf("(unit_destroyer.go) %v \n", err)
// 			return false
// 		}

// 		//for Structures targetting self, reset combat
// 		err = resetStructuresTargetingSelfUD(world, targetFilter)
// 		if err != nil {
// 			fmt.Printf("(unit_destroyer.go) %v \n", err)
// 			return false
// 		}

// 		//for projectiles targetting self destroy
// 		err = destroyProjectilesTargetingSelfUD(world, targetFilter, p1, p2)
// 		if err != nil {
// 			fmt.Printf("(unit_destroyer.go) %v \n", err)
// 			return false
// 		}

// 		//filter for units targeting self
// 		targetFilter = cardinal.ComponentFilter(func(m comp.Target) bool {
// 			return m.Target == id
// 		})
// 		//for app special powers targettting self
// 		err = destroySPTargetingSelfUD(world, targetFilter)
// 		if err != nil {
// 			fmt.Printf("(unit_destroyer.go) %v \n", err)
// 			return false
// 		}

// 		//get collision Hash
// 		CollisionSpartialHash, err := cardinal.GetComponent[comp.SpatialHash](world, gameState)
// 		if err != nil {
// 			fmt.Printf("error retrieving SpartialHash component on tempSpartialHash (unit_destroyer.go): %s \n", err)
// 			return false
// 		}
// 		RemoveObjectFromSpatialHash(CollisionSpartialHash, id, UnitPosition.PositionVectorX, UnitPosition.PositionVectorY, UnitRadius.UnitRadius)

// 		//remove entity
// 		if err := cardinal.Remove(world, id); err != nil {
// 			fmt.Println("Error removing entity (unit_destroyer.go): \n", err) // Log error if any
// 			return false                                                      // Stop iteration on error
// 		}

// 		p1.RemovalList[uid.UID] = true //add removed units to players removal list
// 		p2.RemovalList[uid.UID] = true

// 		//add removed unit to player1 removal list component
// 		if err := cardinal.SetComponent(world, gameState, p1); err != nil {
// 			fmt.Printf("error updating player1 component (unit_destroyer.go): %s \n", err)
// 			return false
// 		}

// 		//add removed unit to player2 removal list component
// 		if err := cardinal.SetComponent(world, gameState, p2); err != nil {
// 			fmt.Printf("error updating player2 component (unit_destroyer.go): %s \n", err)
// 			return false
// 		}
// 		return true
// 	})
// 	if err != nil {
// 		return fmt.Errorf("error retrieving unit entities (unit_destroyer.go): %w", err)
// 	}
// 	return nil
// }
