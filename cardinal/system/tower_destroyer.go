package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

// destroy towers with no health
func TowerDestroyerSystem(world cardinal.WorldContext) error {
	// Filter for no HP
	healthFilter := cardinal.ComponentFilter(func(m comp.Health) bool {
		return m.CurrentHP <= 0
	})
	//for each unit with no hp's ids
	err := cardinal.NewSearch().Entity(
		filter.Contains(StructureFilters())).
		Where(healthFilter).Each(world, func(id types.EntityID) bool {

		//get needed compoenents
		MatchID, state, UnitPosition, UnitRadius, team, health, unitName, err := getStructComponentsTD(world, id)
		if err != nil {
			fmt.Printf("(tower destroyer.go): %v \n", err)
			return false
		}

		//get game state
		gameState, err := getGameStateGSS(world, MatchID)
		if err != nil {
			fmt.Printf("(tower destroyer.go) %v \n", err)
			return false
		}
		//get player components
		p1, p2, err := getPlayerComponentsGSS(world, gameState)
		if err != nil {
			fmt.Printf("(tower destroyer.go) %v \n", err)
			return false
		}

		//filter for units targeting self
		targetFilter := cardinal.ComponentFilter(func(m comp.Attack) bool {
			return m.Target == id
		})

		//for units targetting self, reset combat
		err = resetUnitsTargetingSelfUD(world, targetFilter)
		if err != nil {
			fmt.Printf("(tower destroyer.go) %v \n", err)
			return false
		}

		//for projectiles targetting self destroy
		err = destroyProjectilesTargetingSelfUD(world, targetFilter, p1, p2)
		if err != nil {
			fmt.Printf("(tower destroyer.go) %v \n", err)
			return false
		}

		//filter for sp targeting self
		targetFilter = cardinal.ComponentFilter(func(m comp.Target) bool {
			return m.Target == id
		})
		//for app special powers targettting self
		err = destroySPTargetingSelfUD(world, targetFilter)
		if err != nil {
			fmt.Printf("(tower destroyer.go) %v \n", err)
			return false
		}

		//get collision Hash
		CollisionSpartialHash, err := cardinal.GetComponent[comp.SpatialHash](world, gameState)
		if err != nil {
			fmt.Printf("error retrieving SpartialHash component on tempSpartialHash (tower destroyer.go): %s \n", err)
			return false
		}
		RemoveObjectFromSpatialHash(CollisionSpartialHash, id, UnitPosition.PositionVectorX, UnitPosition.PositionVectorY, UnitRadius.UnitRadius)

		if unitName.UnitName != "Base" {
			if team.Team == "Blue" {
				//change tower team
				team.Team = "Red"
				if err := cardinal.SetComponent(world, id, team); err != nil {
					fmt.Printf("error updating team component (tower destroyer.go): %s \n", err)
					return false
				}
				//add structure to spatial hash collision map
				AddObjectSpatialHash(CollisionSpartialHash, id, UnitPosition.PositionVectorX, UnitPosition.PositionVectorY, UnitRadius.UnitRadius, "Red", "structure")
			} else {
				//change tower team
				team.Team = "Blue"
				if err := cardinal.SetComponent(world, id, team); err != nil {
					fmt.Printf("error updating team component (tower destroyer.go): %s \n", err)
					return false
				}
				//add structure to spatial hash collision map
				AddObjectSpatialHash(CollisionSpartialHash, id, UnitPosition.PositionVectorX, UnitPosition.PositionVectorY, UnitRadius.UnitRadius, "Blue", "structure")
			}

			state.State = "Converting"
			if err := cardinal.SetComponent(world, id, state); err != nil {
				fmt.Printf("error updating state component (tower destroyer.go): %s \n", err)
				return false
			}

			health.CurrentHP = health.MaxHP / 4
			if err := cardinal.SetComponent(world, id, health); err != nil {
				fmt.Printf("error updating health component (tower destroyer.go): %s \n", err)
				return false
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
			fmt.Printf("error on vampire attack (tower destroyer.go): %s \n", err)
			return false
		}

		//add removed unit to player1 removal list component
		if err := cardinal.SetComponent(world, gameState, p1); err != nil {
			fmt.Printf("error updating player1 component (tower destroyer.go): %s \n", err)
			return false
		}

		//add removed unit to player2 removal list component
		if err := cardinal.SetComponent(world, gameState, p2); err != nil {
			fmt.Printf("error updating player2 component (tower destroyer.go): %s \n", err)
			return false
		}
		return true
	})
	if err != nil {
		return fmt.Errorf("error retrieving unit entities (tower destroyer.go): %w ", err)
	}
	return nil
}

// fetches unit components needed for spatial hash removal
func getStructComponentsTD(world cardinal.WorldContext, id types.EntityID) (matchID *comp.MatchId, state *comp.State, unitPosition *comp.Position, unitRadius *comp.UnitRadius, team *comp.Team, health *comp.Health, unitName *comp.UnitName, err error) {
	unitPosition, err = cardinal.GetComponent[comp.Position](world, id)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving enemy Position component (tower destroyer.go): %v", err)
	}
	unitRadius, err = cardinal.GetComponent[comp.UnitRadius](world, id)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving enemy Radius component (tower destroyer.go): %v", err)
	}
	matchID, err = cardinal.GetComponent[comp.MatchId](world, id)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving MatchID component (tower destroyer.go): %v", err)
	}
	state, err = cardinal.GetComponent[comp.State](world, id)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving state component (tower destroyer.go): %v", err)
	}
	team, err = cardinal.GetComponent[comp.Team](world, id)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Team component (tower destroyer.go): %v", err)
	}
	health, err = cardinal.GetComponent[comp.Health](world, id)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving health  component (tower destroyer.go): %v", err)
	}
	unitName, err = cardinal.GetComponent[comp.UnitName](world, id)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving unit name  component (tower destroyer.go)): %v", err)
	}
	return matchID, state, unitPosition, unitRadius, team, health, unitName, nil
}
