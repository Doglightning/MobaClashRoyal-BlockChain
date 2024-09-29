package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

// destroy units with no health
func UnitDestroyerSystem(world cardinal.WorldContext) error {
	// Filter for no HP
	healthFilter := cardinal.ComponentFilter(func(m comp.Health) bool {
		return m.CurrentHP <= 0
	})
	//for each unit with no hp's ids
	err := cardinal.NewSearch().Entity(
		filter.Contains(filter.Component[comp.UnitTag]())).
		Where(healthFilter).Each(world, func(id types.EntityID) bool {

		//get needed compoenents
		MatchID, uid, UnitPosition, UnitRadius, err := getUnitComponentsUD(world, id)
		if err != nil {
			fmt.Printf("(unit_destroyer.go): %v \n", err)
			return false
		}

		//get game state
		gameState, err := getGameStateGSS(world, MatchID)
		if err != nil {
			fmt.Printf("(unit_destroyer.go) %v \n", err)
			return false
		}
		//get player components
		p1, p2, err := getPlayerComponentsGSS(world, gameState)
		if err != nil {
			fmt.Printf("(unit_destroyer.go) %v \n", err)
			return false
		}

		//filter for units targeting self
		targetFilter := cardinal.ComponentFilter(func(m comp.Attack) bool {
			return m.Target == id
		})

		//for units targetting self, reset combat
		err = resetUnitsTargetingSelfUD(world, targetFilter)
		if err != nil {
			fmt.Printf("(unit_destroyer.go) %v \n", err)
			return false
		}

		//for Structures targetting self, reset combat
		err = resetStructuresTargetingSelfUD(world, targetFilter)
		if err != nil {
			fmt.Printf("(unit_destroyer.go) %v \n", err)
			return false
		}

		//for projectiles targetting self destroy
		err = destroyProjectilesTargetingSelfUD(world, targetFilter, p1, p2)
		if err != nil {
			fmt.Printf("(unit_destroyer.go) %v \n", err)
			return false
		}

		//filter for units targeting self
		targetFilter = cardinal.ComponentFilter(func(m comp.Target) bool {
			return m.Target == id
		})
		//for app special powers targettting self
		err = destroySPTargetingSelfUD(world, targetFilter)
		if err != nil {
			fmt.Printf("(unit_destroyer.go) %v \n", err)
			return false
		}

		//get collision Hash
		CollisionSpartialHash, err := cardinal.GetComponent[comp.SpatialHash](world, gameState)
		if err != nil {
			fmt.Printf("error retrieving SpartialHash component on tempSpartialHash (unit_destroyer.go): %s \n", err)
			return false
		}
		RemoveObjectFromSpatialHash(CollisionSpartialHash, id, UnitPosition.PositionVectorX, UnitPosition.PositionVectorY, UnitRadius.UnitRadius)

		//remove entity
		if err := cardinal.Remove(world, id); err != nil {
			fmt.Println("Error removing entity (unit_destroyer.go): \n", err) // Log error if any
			return false                                                      // Stop iteration on error
		}

		p1.RemovalList[uid.UID] = true //add removed units to players removal list
		p2.RemovalList[uid.UID] = true

		//add removed unit to player1 removal list component
		if err := cardinal.SetComponent(world, gameState, p1); err != nil {
			fmt.Printf("error updating player1 component (unit_destroyer.go): %s \n", err)
			return false
		}

		//add removed unit to player2 removal list component
		if err := cardinal.SetComponent(world, gameState, p2); err != nil {
			fmt.Printf("error updating player2 component (unit_destroyer.go): %s \n", err)
			return false
		}
		return true
	})
	if err != nil {
		return fmt.Errorf("error retrieving unit entities (unit_destroyer.go): %w", err)
	}
	return nil
}

// for each unit targeting targetFilter, reset combat to false and attack frame to 0
func resetUnitsTargetingSelfUD(world cardinal.WorldContext, targetFilter cardinal.FilterFn) error {
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
		filter.Contains(StructureFilters())).
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
		filter.Exact(ProjectileFilters())).
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

// fetches unit components needed for spatial hash removal
func getUnitComponentsUD(world cardinal.WorldContext, id types.EntityID) (matchID *comp.MatchId, uid *comp.UID, unitPosition *comp.Position, unitRadius *comp.UnitRadius, err error) {
	unitPosition, err = cardinal.GetComponent[comp.Position](world, id)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error retrieving enemy Position component (unit_destroyer.go): %v", err)
	}
	unitRadius, err = cardinal.GetComponent[comp.UnitRadius](world, id)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error retrieving enemy Radius component (unit_destroyer.go): %v", err)
	}
	matchID, err = cardinal.GetComponent[comp.MatchId](world, id)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error retrieving MatchID component (unit_destroyer.go): %v", err)
	}
	uid, err = cardinal.GetComponent[comp.UID](world, id)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error retrieving UID component (unit_destroyer.go): %v", err)
	}
	return matchID, uid, unitPosition, unitRadius, nil
}
