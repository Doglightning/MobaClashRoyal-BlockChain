package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

func DestroyerSystem(world cardinal.WorldContext) error {
	// Filter for no HP
	healthFilter := cardinal.ComponentFilter(func(m comp.Health) bool {
		return m.CurrentHP <= 0
	})

	// Filter for destoryed Entities
	destroyedFilter := cardinal.ComponentFilter(func(m comp.Destroyed) bool {
		return m.Destroyed
	})

	//for each unit with no hp's ids
	err := cardinal.NewSearch().Entity(filter.Contains()).
		Where(cardinal.OrFilter(healthFilter, destroyedFilter)).Each(world, func(id types.EntityID) bool {

		// get attack component
		atk, err := cardinal.GetComponent[comp.Attack](world, id)
		if err != nil {
			fmt.Printf("error retrieving unit Attack component (phase_Destroyer.go): %v \n", err)
			return false
		}

		// // projectile attack logic
		// if atk.Class == "projectile" {
		// 	err = ProjectileAttack(world, id, atk)
		// 	if err != nil {
		// 		fmt.Printf("%v \n", err)
		// 		return false
		// 	}

		// 	// basic melee/range attack logic
		// } else
		if atk.Class == "melee" || atk.Class == "range" || atk.Class == "air" {
			err = ClassDestroySystem(world, id)
			if err != nil {
				fmt.Printf("%v \n", err)
				return false
			}
		}
		// else if atk.Class == "structure" {
		// 	err = StructureAttack(world, id, atk)
		// 	if err != nil {
		// 		fmt.Printf("%v \n", err)
		// 		return false
		// 	}
		// }

		return true
	})

	return err
}

func unitDestroyerDefault(world cardinal.WorldContext, id types.EntityID) error {
	//get needed compoenents
	MatchID, uid, UnitPosition, UnitRadius, err := GetComponents4[comp.MatchId, comp.UID, comp.Position, comp.UnitRadius](world, id)
	if err != nil {
		return fmt.Errorf("4 (unit_destroyer): %v ", err)
	}

	//get game state
	gameState, err := getGameStateGSS(world, MatchID)
	if err != nil {
		return fmt.Errorf("(unit_destroyer) %v ", err)
	}
	//get player components
	p1, p2, err := getPlayerComponentsGSS(world, gameState)
	if err != nil {
		return fmt.Errorf("(unit_destroyer) %v ", err)
	}

	//filter for units targeting self
	targetFilter := cardinal.ComponentFilter(func(m comp.Attack) bool {
		return m.Target == id
	})

	//for units targetting self, reset combat
	err = resetUnitsTargetingSelf(world, targetFilter)
	if err != nil {
		return fmt.Errorf("(unit_destroyer) %v ", err)
	}

	//for Structures targetting self, reset combat
	err = resetStructuresTargetingSelfUD(world, targetFilter)
	if err != nil {
		return fmt.Errorf("(unit_destroyer) %v ", err)
	}

	//filter for units targeting self
	destroyedFilter := cardinal.ComponentFilter(func(m comp.Destroyed) bool {
		return !m.Destroyed
	})

	projectileFilter := cardinal.AndFilter(targetFilter, destroyedFilter)

	//for projectiles targetting self destroy
	err = destroyProjectilesTargetingSelfUD(world, projectileFilter, p1, p2)
	if err != nil {
		return fmt.Errorf("(unit_destroyer) %v ", err)
	}

	//filter for units targeting self
	targetFilter = cardinal.ComponentFilter(func(m comp.Target) bool {
		return m.Target == id
	})
	//for app special powers targettting self
	err = destroySPTargetingSelfUD(world, targetFilter)
	if err != nil {
		return fmt.Errorf("(unit_destroyer) %v ", err)
	}

	//get collision Hash
	CollisionSpartialHash, err := cardinal.GetComponent[comp.SpatialHash](world, gameState)
	if err != nil {
		return fmt.Errorf("error retrieving SpartialHash component on tempSpartialHash (unit_destroyer): %s ", err)
	}

	//remove entity
	if err := cardinal.Remove(world, id); err != nil {
		return fmt.Errorf("error removing entity (unit_destroyer): %v", err) // Log error if any
	}

	RemoveObjectFromSpatialHash(CollisionSpartialHash, id, UnitPosition.PositionVectorX, UnitPosition.PositionVectorY, UnitRadius.UnitRadius)

	p1.RemovalList[uid.UID] = true //add removed units to players removal list
	p2.RemovalList[uid.UID] = true

	//add removed unit to player1 removal list component
	if err := cardinal.SetComponent(world, gameState, p1); err != nil {
		return fmt.Errorf("error updating player1 component (unit_destroyer): %s ", err)
	}

	//add removed unit to player2 removal list component
	if err := cardinal.SetComponent(world, gameState, p2); err != nil {
		return fmt.Errorf("error updating player2 component (unit_destroyer): %s ", err)
	}
	return nil
}

// for each unit targeting targetFilter, reset combat to false and attack frame to 0
func resetUnitsTargetingSelf(world cardinal.WorldContext, targetFilter cardinal.FilterFn) error {
	//for each targetting unit
	err := cardinal.NewSearch().Entity(
		filter.Contains(filter.Component[comp.UnitTag]())).
		Where(targetFilter).Each(world, func(enemyID types.EntityID) bool {

		name, err := cardinal.GetComponent[comp.UnitName](world, enemyID)
		if err != nil {
			fmt.Printf("error getting unit name component (unit_destroyer.go) \n")
			return false
		}

		err = ClassResetCombat(world, enemyID, name.UnitName)
		if err != nil {
			fmt.Printf("error running ResetCombat (unit_destroyer.go) \n")
			return false
		}

		return true
	})
	if err != nil {
		return fmt.Errorf("error retrieving unit entities (resetUnitsTargetingSelfUD): %s", err)
	}
	return nil
}

// for each unit targeting targetFilter, reset combat to false and attack frame to 0
func resetStructuresTargetingSelfUD(world cardinal.WorldContext, targetFilter cardinal.FilterFn) error {
	//for each targetting unit
	err := cardinal.NewSearch().Entity(
		filter.Contains(filter.Component[comp.StructureTag]())).
		Where(targetFilter).Each(world, func(structID types.EntityID) bool {
		//reset attack component
		cardinal.UpdateComponent(world, structID, func(attack *comp.Attack) *comp.Attack {
			if attack == nil {
				fmt.Printf("error retrieving enemy attack component (unit_destroyer.go): \n")
				return nil
			}
			attack.Combat = false
			attack.Frame = 0
			return attack
		})
		return true
	})
	if err != nil {
		return fmt.Errorf("error retrieving unit entities (resetUnitsTargetingSelfUD): %s", err)
	}
	return nil
}

// for each projectile targeting targetFilter, destroy
func destroyProjectilesTargetingSelfUD(world cardinal.WorldContext, targetFilter cardinal.FilterFn, p1 *comp.Player1, p2 *comp.Player2) error {
	//for each targetting projectile
	err := cardinal.NewSearch().Entity(
		filter.Contains(filter.Component[comp.ProjectileTag]())).
		Where(targetFilter).Each(world, func(projectileID types.EntityID) bool {

		//get projectile uid
		projectileUID, err := cardinal.GetComponent[comp.UID](world, projectileID)
		if err != nil {
			fmt.Printf("error retrieving projectile UID component (unit_destroyer.go): %s \n", err)
			return false
		}

		p1.RemovalList[projectileUID.UID] = true // add to players removal lists
		p2.RemovalList[projectileUID.UID] = true

		//remove entity
		if err := cardinal.Remove(world, projectileID); err != nil {
			fmt.Println("Error removing entity projectile (unit_destroyer.go):", err)
			return false
		}
		return true
	})
	if err != nil {
		return fmt.Errorf("error retrieving projectile entities (destroyProjectilesTargetingSelfUD): %s", err)
	}
	return nil
}

// for each Special power targeting targetFilter, destroy
func destroySPTargetingSelfUD(world cardinal.WorldContext, targetFilter cardinal.FilterFn) error {
	//for each targetting sp
	err := cardinal.NewSearch().Entity(
		filter.Contains(filter.Component[comp.SpEntity]())).
		Where(targetFilter).Each(world, func(spID types.EntityID) bool {

		//remove entity
		if err := cardinal.Remove(world, spID); err != nil {
			fmt.Printf("Error removing entity sp (unit_destroyer.go): %v \n", err)
			return false
		}
		return true
	})
	if err != nil {
		return fmt.Errorf("error retrieving sp entities (destroyProjectilesTargetingSelfUD): %s", err)
	}
	return nil
}
