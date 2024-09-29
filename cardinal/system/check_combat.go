package system

import (
	comp "MobaClashRoyal/component"
	"container/list"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

// make sure any units who are in range to attack get set for attack b4 attack phase
func CombatCheckSystem(world cardinal.WorldContext) error {

	err := unitCombatSearch(world)
	if err != nil {
		fmt.Printf("error searching for unit combat (check_combat.go): %v \n", err)
	}

	err = structureCombatSearch(world)
	if err != nil {
		fmt.Printf("error searching for structure combat (check_combat.go): %v \n", err)
	}

	return err
}

// check if unit not in combat can find a unit to be in combat with
func unitCombatSearch(world cardinal.WorldContext) error {
	// filter not in combat units
	combatFilter := cardinal.ComponentFilter(func(m comp.Attack) bool {
		return !m.Combat
	})
	//for each unit not in combat
	err := cardinal.NewSearch().Entity(
		filter.Contains(UnitFilters())).
		Where(combatFilter).Each(world, func(id types.EntityID) bool {

		//get Unit CC component
		cc, err := cardinal.GetComponent[comp.CC](world, id)
		if err != nil {
			fmt.Printf("error getting unit cc component (unit_movement.go): %v \n", err)
		}

		if cc.Stun > 0 { //if unit stunned cannot attack
			return false
		}

		//get attack component
		uAtk, err := cardinal.GetComponent[comp.Attack](world, id)
		if err != nil {
			fmt.Printf("failed to get attack comp (check_combat.go): %v \n", err)
			return false
		}

		if !uAtk.Combat { //not in combat
			//get Unit Components
			uPos, uRadius, uAtk, uTeam, MatchID, err := getUnitComponentsCC(world, id)
			if err != nil {
				fmt.Printf("(check_combat.go) -%v \n", err)
				return false
			}

			// get collision Hash
			collisionHash, err := getCollisionHashGSS(world, MatchID)
			if err != nil {
				fmt.Printf("error retrieving SpartialHash component on tempSpartialHash (check_combat.go): %s \n", err)
				return false
			}
			//find closest enemy
			eID, eX, eY, eRadius, found := findClosestEnemy(collisionHash, id, uPos.PositionVectorX, uPos.PositionVectorY, uAtk.AggroRadius, uTeam.Team, uAtk.Class)
			if found { //found enemy
				// Calculate squared distance between the unit and the enemy, minus their radii
				adjustedDistance := distanceBetweenTwoPoints(uPos.PositionVectorX, uPos.PositionVectorY, eX, eY) - float32(eRadius) - float32(uRadius.UnitRadius)
				//if within attack range
				if adjustedDistance <= float32(uAtk.AttackRadius) {
					uAtk.Combat = true
					uAtk.Target = eID
					//set attack component
					if err = cardinal.SetComponent(world, id, uAtk); err != nil {
						fmt.Printf("error setting attack component (check_combat.go): %v \n", err)
						return false
					}
				}

			}
		}
		return true
	})
	return err
}

// check if a structure not in combat can find a unit in range to attack
func structureCombatSearch(world cardinal.WorldContext) error {
	//for each structure not in combat
	err := cardinal.NewSearch().Entity(
		filter.Contains(StructureFilters())).
		Each(world, func(id types.EntityID) bool {
			//get attack component
			uAtk, err := cardinal.GetComponent[comp.Attack](world, id)
			if err != nil {
				fmt.Printf("failed to get attack comp (structureCombatSearch - check_combat.go): %v \n", err)
				return false
			}

			state, err := cardinal.GetComponent[comp.State](world, id)
			if err != nil {
				fmt.Printf("failed to get state comp (structureCombatSearch - check_combat.go): %v \n", err)
				return false
			}

			if state.State != "Converting" { //if tower is not converting teams

				if uAtk.Combat { // in combat make sure target still in range
					//get Unit Components
					uPos, uRadius, uAtk, err := getTowerComponentsCC(world, id)
					if err != nil {
						fmt.Printf("(structureCombatSearch - check_combat.go) -%v \n", err)
						return false
					}

					//get Unit Components
					ePos, eRadius, err := getTowerTargetComponentsCC(world, uAtk.Target)
					if err != nil {
						fmt.Printf("(structureCombatSearch - check_combat.go) -%v \n", err)
						return false
					}

					// Calculate squared distance between the unit and the enemy, minus their radii
					adjustedDistance := distanceBetweenTwoPoints(uPos.PositionVectorX, uPos.PositionVectorY, ePos.PositionVectorX, ePos.PositionVectorY) - float32(eRadius.UnitRadius) - float32(uRadius.UnitRadius)
					//if within attack range
					if adjustedDistance > float32(uAtk.AttackRadius) {
						uAtk.Combat = false
						uAtk.Frame = 0
						//set attack component
						if err = cardinal.SetComponent(world, id, uAtk); err != nil {
							fmt.Printf("error setting attack component (structureCombatSearch - check_combat.go): %v \n", err)
							return false
						}
					}

				}

				if !uAtk.Combat { //not in combat
					//get Unit Components
					uPos, uRadius, uAtk, uTeam, MatchID, err := getUnitComponentsCC(world, id)
					if err != nil {
						fmt.Printf("(structureCombatSearch - check_combat.go) -%v \n", err)
						return false
					}

					// get collision Hash
					collisionHash, err := getCollisionHashGSS(world, MatchID)
					if err != nil {
						fmt.Printf("error retrieving SpartialHash component on tempSpartialHash (structureCombatSearch - check_combat.go): %s  \n", err)
						return false
					}
					//find closest enemy
					eID, eX, eY, eRadius, found := findClosestEnemy(collisionHash, id, uPos.PositionVectorX, uPos.PositionVectorY, uAtk.AggroRadius, uTeam.Team, uAtk.Class)
					if found { //found enemy
						// Calculate squared distance between the unit and the enemy, minus their radii
						adjustedDistance := distanceBetweenTwoPoints(uPos.PositionVectorX, uPos.PositionVectorY, eX, eY) - float32(eRadius) - float32(uRadius.UnitRadius)
						//if within attack range
						if adjustedDistance <= float32(uAtk.AttackRadius) {
							uAtk.Combat = true
							uAtk.Target = eID
							//set attack component
							if err = cardinal.SetComponent(world, id, uAtk); err != nil {
								fmt.Printf("error setting attack component (structureCombatSearch - check_combat.go): %v \n", err)
								return false
							}
						}

					}
				}
			}
			return true
		})
	return err
}

// FindClosestEnemy performs a BFS search from the unit's position outward within the attack radius.
func findClosestEnemy(hash *comp.SpatialHash, objID types.EntityID, startX, startY float32, attackRadius int, team, class string) (types.EntityID, float32, float32, int, bool) {
	queue := list.New()                                                              //queue of cells to check
	visited := make(map[string]bool)                                                 //cells checked
	queue.PushBack(&comp.Position{PositionVectorX: startX, PositionVectorY: startY}) //insert starting position to queue
	minDist := float32(attackRadius * attackRadius)                                  // Using squared distance to avoid sqrt calculations.
	closestEnemy := types.EntityID(0)
	closestX, closestY := float32(0), float32(0)
	closestRadius := int(0)
	foundEnemy := false

	//while units in queue
	for queue.Len() > 0 {
		pos := queue.Remove(queue.Front()).(*comp.Position) // remove first Item
		x, y := pos.PositionVectorX, pos.PositionVectorY
		cellX, cellY := calculateSpatialHash(hash, x, y) //Find the hash key for grid size
		hashKey := fmt.Sprintf("%d,%d", cellX, cellY)    //create key

		// Prevent re-checking the same cell
		if _, found := visited[hashKey]; found {
			continue
		}
		visited[hashKey] = true

		if cell, exists := hash.Cells[hashKey]; exists { //if unit found in cell
			for i, id := range cell.UnitIDs { //go over each unit in cell
				if cell.Team[i] != team && id != objID { //if unit in cell is enemy and not self

					if (class == "melee" && cell.Type[i] != "air") || class == "range" || class == "air" || class == "structure" { // make sure melee cannot attack air

						distSq := (cell.PositionsX[i]-startX)*(cell.PositionsX[i]-startX) + (cell.PositionsY[i]-startY)*(cell.PositionsY[i]-startY) - float32(cell.Radii[i]*cell.Radii[i])
						//if distance is smaller then closest unit found so far
						if distSq < minDist {
							minDist = distSq
							closestEnemy = id
							closestX, closestY = cell.PositionsX[i], cell.PositionsY[i]
							closestRadius = cell.Radii[i]
							foundEnemy = true
						}

					}

				}
			}
		}

		// Add neighboring cells to the queue if within range
		if !foundEnemy {
			for dx := -hash.CellSize; dx <= hash.CellSize; dx += hash.CellSize {
				for dy := -hash.CellSize; dy <= hash.CellSize; dy += hash.CellSize {
					nx, ny := x+float32(dx), y+float32(dy)
					//check if new cell being added is still within attack range
					if (nx-startX)*(nx-startX)+(ny-startY)*(ny-startY) <= float32(attackRadius*attackRadius) {
						queue.PushBack(&comp.Position{PositionVectorX: nx, PositionVectorY: ny}) // add to queue
					}
				}
			}
		}
	}
	return closestEnemy, closestX, closestY, closestRadius, foundEnemy
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////Default Combat Reset Functions////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func resetCombat(world cardinal.WorldContext, id types.EntityID) error {
	// reset attack component
	err := cardinal.UpdateComponent(world, id, func(attack *comp.Attack) *comp.Attack {
		if attack == nil {
			fmt.Printf("error retrieving enemy attack component (resetCombat/check combat.go): ")
			return nil
		}
		attack.Combat = false
		attack.Frame = 0
		return attack
	})
	if err != nil {
		return fmt.Errorf("error updating attack comp (resetCombat/check combat.go): %v", err)
	}
	return nil
}

// overwrite base destruction.
// if unit being attacked by channel dies don't cancel attack.
func channelingResetCombat(world cardinal.WorldContext, id types.EntityID) error {
	//reset attack component
	err := cardinal.UpdateComponent(world, id, func(attack *comp.Attack) *comp.Attack {
		if attack == nil {
			fmt.Printf("error retrieving enemy attack component (channelingResetCombat/check combat.go): \n")
			return nil
		}
		//get special power component
		sp, err := cardinal.GetComponent[comp.Sp](world, id)
		if err != nil {
			fmt.Printf("error retrieving special power comp (channelingResetCombat/check combat.go): \n")
			return nil
		}

		if attack.Frame < sp.DamageFrame && sp.Charged { //if target dies b4 fire attack goes off
			//reset units combat
			attack.Frame = 0
			attack.Combat = false
			attack.State = "Default"
		} else { //if unit started channeling fire
			attack.State = "Channeling"
			attack.Target = id //set target to self to not get errors if triggering functions that ref this but unit is dead
		}
		return attack
	})
	if err != nil {
		return fmt.Errorf("error updating attack comp (channelingResetCombat/check combat.go): %v", err)
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////Default get Component Functions///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// GetUnitComponents fetches all necessary components related to a unit entity.
func getUnitComponentsCC(world cardinal.WorldContext, unitID types.EntityID) (*comp.Position, *comp.UnitRadius, *comp.Attack, *comp.Team, *comp.MatchId, error) {
	position, err := cardinal.GetComponent[comp.Position](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Position component (getUnitComponentsCC): %v", err)
	}
	unitRadius, err := cardinal.GetComponent[comp.UnitRadius](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Unit Radius component (getUnitComponentsCC): %v", err)
	}
	unitAttack, err := cardinal.GetComponent[comp.Attack](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Unit attack component (getUnitComponentsCC): %v", err)
	}
	team, err := cardinal.GetComponent[comp.Team](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Team component (getUnitComponentsCC): %v", err)
	}
	matchId, err := cardinal.GetComponent[comp.MatchId](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error retrieving MatchId component (getUnitComponentsCC): %v", err)
	}
	return position, unitRadius, unitAttack, team, matchId, nil
}

// GetTowerComponents fetches all necessary components related to a Tower entity.
func getTowerComponentsCC(world cardinal.WorldContext, unitID types.EntityID) (*comp.Position, *comp.UnitRadius, *comp.Attack, error) {
	position, err := cardinal.GetComponent[comp.Position](world, unitID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error retrieving Position component (getTowerComponentsCC): %v", err)
	}
	unitRadius, err := cardinal.GetComponent[comp.UnitRadius](world, unitID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error retrieving Unit Radius component (getTowerComponentsCC): %v", err)
	}
	unitAttack, err := cardinal.GetComponent[comp.Attack](world, unitID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error retrieving Unit attack component (getTowerComponentsCC): %v", err)
	}
	return position, unitRadius, unitAttack, nil
}

// GetTowerComponents fetches all necessary components related to a Tower entity.
func getTowerTargetComponentsCC(world cardinal.WorldContext, unitID types.EntityID) (*comp.Position, *comp.UnitRadius, error) {
	position, err := cardinal.GetComponent[comp.Position](world, unitID)
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving Position component (getTowerTargetComponentsCC): %v", err)
	}
	unitRadius, err := cardinal.GetComponent[comp.UnitRadius](world, unitID)
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving Unit Radius component (getTowerTargetComponentsCC): %v", err)
	}
	return position, unitRadius, nil
}
