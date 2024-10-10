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
	err := cardinal.NewSearch().
		Entity(
			filter.Or(
				filter.Contains(filter.Component[comp.UnitTag]()),
				filter.Contains(filter.Component[comp.ProjectileTag]()),
				filter.Contains(filter.Component[comp.StructureTag]()),
				filter.Contains(filter.Component[comp.SpEntity]()),
			),
		).
		Where(cardinal.OrFilter(healthFilter, destroyedFilter)).
		Each(world, func(id types.EntityID) bool {

			// get attack component
			class, err := cardinal.GetComponent[comp.Class](world, id)
			if err != nil {
				fmt.Printf("error retrieving unit Attack component (phase_Destroyer.go): %v \n", err)
				return false
			}

			// projectile destroyer
			if class.Class == "projectile" {
				err = projectileDestroyerDefault(world, id)
				if err != nil {
					fmt.Printf("%v \n", err)
					return false
				}

				// unit destroyer
			} else if class.Class == "melee" || class.Class == "range" || class.Class == "air" {
				err = ClassDestroySystem(world, id)
				if err != nil {
					fmt.Printf("%v \n", err)
					return false
				}
				//structure destroyer
			} else if class.Class == "structure" {
				err = structureDestroyerDefault(world, id)
				if err != nil {
					fmt.Printf("%v \n", err)
					return false
				}
			} else if class.Class == "sp" {
				err = projectileDestroyerDefault(world, id)
				if err != nil {
					fmt.Printf("%v \n", err)
					return false
				}
			}

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

	//filter for units targeting self
	targetSpFilter := cardinal.ComponentFilter(func(m comp.Sp) bool {
		return m.Target == id
	})

	//for unit Sp targettign self
	err = resetUnitSpTargetingSelf(world, targetSpFilter)
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

	//set collision hash, player1 and player2
	if err = SetComponents3(world, gameState, CollisionSpartialHash, p1, p2); err != nil {
		return fmt.Errorf("(unit_destroyer): %v", err)
	}

	return nil
}

func structureDestroyerDefault(world cardinal.WorldContext, id types.EntityID) error {
	//get needed compoenents
	MatchID, state, UnitPosition, UnitRadius, team, health, unitName, err := GetComponents7[comp.MatchId, comp.State, comp.Position, comp.UnitRadius, comp.Team, comp.Health, comp.UnitName](world, id)
	if err != nil {
		return fmt.Errorf("tower components (tower destroyer.go): %v", err)
	}

	//get game state
	gameState, err := getGameStateGSS(world, MatchID)
	if err != nil {
		return fmt.Errorf("(tower destroyer.go) %v", err)
	}
	//get player components
	p1, p2, err := getPlayerComponentsGSS(world, gameState)
	if err != nil {
		return fmt.Errorf("(tower destroyer.go) %v", err)
	}

	//filter for units targeting self
	targetFilter := cardinal.ComponentFilter(func(m comp.Attack) bool {
		return m.Target == id
	})

	//for units targetting self, reset combat
	err = resetUnitsTargetingSelf(world, targetFilter)
	if err != nil {
		return fmt.Errorf("(tower destroyer.go) %v", err)
	}

	//filter for units targeting self
	destroyedFilter := cardinal.ComponentFilter(func(m comp.Destroyed) bool {
		return !m.Destroyed
	})

	projectileFilter := cardinal.AndFilter(targetFilter, destroyedFilter)

	//for projectiles targetting self destroy
	err = destroyProjectilesTargetingSelfUD(world, projectileFilter, p1, p2)
	if err != nil {
		return fmt.Errorf("(tower destroyer.go) %v", err)
	}

	//filter for sp targeting self
	targetFilter = cardinal.ComponentFilter(func(m comp.Target) bool {
		return m.Target == id
	})
	//for app special powers targettting self
	err = destroySPTargetingSelfUD(world, targetFilter)
	if err != nil {
		return fmt.Errorf("(tower destroyer.go) %v", err)
	}

	//get collision Hash
	CollisionSpartialHash, err := cardinal.GetComponent[comp.SpatialHash](world, gameState)
	if err != nil {
		return fmt.Errorf("error retrieving SpartialHash component on tempSpartialHash (tower destroyer.go): %v", err)
	}
	RemoveObjectFromSpatialHash(CollisionSpartialHash, id, UnitPosition.PositionVectorX, UnitPosition.PositionVectorY, UnitRadius.UnitRadius)

	if unitName.UnitName != "Base" { // if a tower change teams
		if team.Team == "Blue" {
			//change tower team
			team.Team = "Red"
			AddObjectSpatialHash(CollisionSpartialHash, id, UnitPosition.PositionVectorX, UnitPosition.PositionVectorY, UnitRadius.UnitRadius, "Red", "structure")
		} else {
			//change tower team
			team.Team = "Blue"
			AddObjectSpatialHash(CollisionSpartialHash, id, UnitPosition.PositionVectorX, UnitPosition.PositionVectorY, UnitRadius.UnitRadius, "Blue", "structure")
		}

		state.State = "Converting"
		health.CurrentHP = health.MaxHP / 4

		//set state and health and team
		if err = SetComponents3(world, id, state, health, team); err != nil {
			return fmt.Errorf("(tower destroyer.go): %v", err)
		}

	}

	//set combat to false
	err = cardinal.UpdateComponent(world, id, func(atk *comp.Attack) *comp.Attack {
		if atk == nil {
			fmt.Printf("error retrieving attack component (tower destroyer.go): \n")
			return nil
		}
		atk.Combat = false
		return atk
	})
	if err != nil {
		return fmt.Errorf("error on vampire attack (tower destroyer.go): %v", err)
	}

	//set collision hash, player1 and player2
	if err = SetComponents3(world, gameState, CollisionSpartialHash, p1, p2); err != nil {
		return fmt.Errorf("(tower destroyer.go): %v", err)
	}

	return nil
}

func projectileDestroyerDefault(world cardinal.WorldContext, id types.EntityID) error {

	//get matchid and uid of projectile
	MatchID, uid, err := GetComponents2[comp.MatchId, comp.UID](world, id)
	if err != nil {
		return fmt.Errorf("get projectile components (projectile_destroyer): %v", err)
	}

	//get game state
	gameState, err := getGameStateGSS(world, MatchID)
	if err != nil {
		return fmt.Errorf("(projectile_destroyer.go) - %v", err)
	}

	//add projectile id to player1 removal list
	cardinal.UpdateComponent(world, gameState, func(player1 *comp.Player1) *comp.Player1 {
		if player1 == nil {
			fmt.Printf("error retrieving player1 component (projectile_destroyer)")
			return nil
		}
		//player1.RemovalList = append(player1.RemovalList, uid.UID)
		player1.RemovalList[uid.UID] = true
		return player1
	})

	//add projectile id to player2 removal list
	cardinal.UpdateComponent(world, gameState, func(player2 *comp.Player2) *comp.Player2 {
		if player2 == nil {
			fmt.Printf("error retrieving player2 component (projectile_destroyer)")
			return nil
		}
		//player1.RemovalList = append(player1.RemovalList, uid.UID)
		player2.RemovalList[uid.UID] = true
		return player2
	})

	//remove projectile
	if err := cardinal.Remove(world, id); err != nil {
		return fmt.Errorf("error removing entity (projectile_destroyer): %v", err)
	}

	return nil
}

// for each unit targeting targetFilter, reset combat to false and attack frame to 0
func resetUnitsTargetingSelf(world cardinal.WorldContext, targetFilter cardinal.FilterFn) error {
	//for each targetting unit
	err := cardinal.NewSearch().Entity(
		filter.Contains(filter.Component[comp.UnitTag]())).
		Where(targetFilter).Each(world, func(enemyID types.EntityID) bool {

		// reset attack component
		err := cardinal.UpdateComponent(world, enemyID, func(attack *comp.Attack) *comp.Attack {
			if attack == nil {
				fmt.Printf("error retrieving enemy attack component (resetUnitsTargetingSelf/phase destroyer.go): ")
				return nil
			}

			err := ClassResetCombat(world, enemyID, attack)
			if err != nil {
				fmt.Printf("error running ResetCombat (resetUnitsTargetingSelf/phase destroyer.go): %v \n", err)
				return nil
			}

			return attack
		})
		if err != nil {
			fmt.Printf("error updating attack comp (resetUnitsTargetingSelf/phase destroyer.go): %v", err)
			return false
		}

		return true
	})
	if err != nil {
		return fmt.Errorf("error retrieving unit entities (resetUnitsTargetingSelfUD): %s", err)
	}
	return nil
}

// for each unit sp targeting targetFilter, attack frame to 0 and resest sp combat and target
func resetUnitSpTargetingSelf(world cardinal.WorldContext, targetFilter cardinal.FilterFn) error {
	//for each targetting unit
	err := cardinal.NewSearch().Entity(
		filter.Contains(filter.Component[comp.UnitTag]())).
		Where(targetFilter).Each(world, func(enemyID types.EntityID) bool {

		// reset attack component
		err := cardinal.UpdateComponent(world, enemyID, func(attack *comp.Attack) *comp.Attack {
			if attack == nil {
				fmt.Printf("error retrieving enemy attack component (resetUnitSpTargetingSelf/phase destroyer.go): ")
				return nil
			}

			attack.Frame = 0

			return attack
		})
		if err != nil {
			fmt.Printf("error updating attack comp (resetUnitSpTargetingSelf/phase destroyer.go): %v", err)
			return false
		}

		// reset sp component
		err = cardinal.UpdateComponent(world, enemyID, func(sp *comp.Sp) *comp.Sp {
			if sp == nil {
				fmt.Printf("error retrieving enemy sp component (resetUnitSpTargetingSelf/phase destroyer.go): ")
				return nil
			}

			sp.Target = 0
			sp.Combat = false

			return sp
		})
		if err != nil {
			fmt.Printf("error updating sp comp (resetUnitSpTargetingSelf/phase destroyer.go): %v", err)
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
				fmt.Printf("error retrieving enemy attack component (resetStructuresTargetingSelfUD): \n")
				return nil
			}
			attack.Combat = false
			attack.Frame = 0
			return attack
		})
		return true
	})
	if err != nil {
		return fmt.Errorf("error retrieving unit entities (resetStructuresTargetingSelfUD): %s", err)
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
			fmt.Printf("error retrieving projectile UID component (destroyProjectilesTargetingSelfUD): %s \n", err)
			return false
		}

		p1.RemovalList[projectileUID.UID] = true // add to players removal lists
		p2.RemovalList[projectileUID.UID] = true

		//remove entity
		if err := cardinal.Remove(world, projectileID); err != nil {
			fmt.Println("Error removing entity projectile (destroyProjectilesTargetingSelfUD):", err)
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
			fmt.Printf("Error removing entity sp (destroySPTargetingSelfUD): %v \n", err)
			return false
		}
		return true
	})
	if err != nil {
		return fmt.Errorf("error retrieving sp entities (destroySPTargetingSelfUD): %s", err)
	}
	return nil
}
