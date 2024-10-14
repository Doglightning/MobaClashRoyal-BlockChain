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

	return nil
}

// check if unit not in combat can find a unit to be in combat with
func unitCombatSearch(world cardinal.WorldContext) error {

	//for each unit not in combat
	err := cardinal.NewSearch().Entity(
		filter.Contains(filter.Component[comp.UnitTag]())).
		Each(world, func(id types.EntityID) bool {

			//get Unit CC component
			cc, err := cardinal.GetComponent[comp.CC](world, id)
			if err != nil {
				fmt.Printf("error getting unit cc component (unitCombatSearch - check_combat.go): %v \n", err)
				return false
			}

			if cc.Stun > 0 { //if unit stunned cannot attack
				return true
			}

			//get Unit Components
			uPos, uRadius, uAtk, uSp, uTeam, MatchID, class, err := GetComponents7[comp.Position, comp.UnitRadius, comp.Attack, comp.Sp, comp.Team, comp.MatchId, comp.Class](world, id)
			if err != nil {
				fmt.Printf("unit comps (unitCombatSearch - check_combat.go) -%v \n", err)
				return false
			}

			// get collision Hash
			collisionHash, err := getCollisionHashGSS(world, MatchID)
			if err != nil {
				fmt.Printf("error retrieving SpartialHash component on tempSpartialHash (unitCombatSearch - check_combat.go): %s \n", err)
				return false
			}

			var atkRadius int
			if uSp.Charged {
				atkRadius = uSp.AttackRadius
			} else {
				atkRadius = uAtk.AttackRadius
			}

			if (uAtk.Combat || uSp.Combat) && uAtk.Target != 0 {
				var targetID types.EntityID

				if uSp.Combat { // get target id if targeting Sp or normal attack
					targetID = uSp.Target
				} else if uAtk.Combat {
					targetID = uAtk.Target
				}
				//get enemy position and radius components
				ePos, eRadius, err := GetComponents2[comp.Position, comp.UnitRadius](world, targetID)
				if err != nil {
					fmt.Printf("combat compoenents (unitCombatSearch - check_combat.go): %s \n", err)
					return false
				}

				//distance between unit and enemy minus their radius
				adjustedDistance := distanceBetweenTwoPoints(uPos.PositionVectorX, uPos.PositionVectorY, ePos.PositionVectorX, ePos.PositionVectorY) - float32(eRadius.UnitRadius) - float32(uRadius.UnitRadius)
				//if out of attack range

				if adjustedDistance > float32(atkRadius)+10 {

					if uSp.Combat {
						uSp.Combat = false //stop special because its not in range anymore
						uAtk.Frame = 0
						uSp.Target = 0

					} else {

						//reset combat
						err = ClassResetCombat(world, id, uAtk)
						if err != nil {
							fmt.Printf("2 (unitCombatSearch - check_combat.go): %v \n", err)
							return false
						}
					}

				}

			}

			if !uAtk.Combat && !uSp.Combat {
				//find closest enemy
				eID, eX, eY, eRadius, found := findClosestEnemy(collisionHash, id, uPos.PositionVectorX, uPos.PositionVectorY, uAtk.AggroRadius, uTeam.Team, class.Class, true)
				if found { //found enemy
					// Calculate squared distance between the unit and the enemy, minus their radii
					adjustedDistance := distanceBetweenTwoPoints(uPos.PositionVectorX, uPos.PositionVectorY, eX, eY) - float32(eRadius) - float32(uRadius.UnitRadius)
					//if within attack range
					if adjustedDistance <= float32(atkRadius)+20 {

						uAtk.Combat = true
						uAtk.Target = eID
						uAtk.Frame = 0
					} else if adjustedDistance > float32(atkRadius) && adjustedDistance <= float32(uAtk.AggroRadius) {
						uAtk.Target = eID
					}

				} else {
					//reset combat
					err = ClassResetCombat(world, id, uAtk)
					if err != nil {
						fmt.Printf("2 (unitCombatSearch - check_combat.go): %v \n", err)
						return false
					}
				}
			}

			//check if in a SP animation or a regular attack
			if uAtk.Frame == 0 && uSp.CurrentSp >= uSp.MaxSp && uAtk.Target != 0 { //In special power
				uSp.Charged = true

				if !uSp.StructureTargetable {
					//get Target name
					tarName, err := cardinal.GetComponent[comp.UnitName](world, uAtk.Target)
					if err != nil {
						fmt.Printf("error retrieving target name component (unitCombatSearch - check_combat.go): %v \n", err)
						return false
					}
					if tarName.UnitName == "Base" || tarName.UnitName == "Tower" {

						found, err := findClosestEnemySP(world, id, uSp) //setup sp target and combat if unit is in charge but cannot target structures
						if err != nil {
							fmt.Printf("(unitCombatSearch - check_combat.go): %v \n", err)
							return false
						}

						if !found {
							uSp.Charged = false
						}

					}
				}

			} else if uAtk.Frame == 0 && uSp.CurrentSp < uSp.MaxSp { // in regular attack
				uSp.Charged = false
			}

			//set combat check comps
			if err = SetComponents2(world, id, uAtk, uSp); err != nil {
				fmt.Printf("main (unitCombatSearch - check_combat.go): %v \n", err)
				return false
			}

			return true
		})
	return err
}

// check if a structure not in combat can find a unit in range to attack
func structureCombatSearch(world cardinal.WorldContext) error {
	//for each structure not in combat
	err := cardinal.NewSearch().Entity(
		filter.Contains(filter.Component[comp.StructureTag]())).
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
					uPos, uRadius, uAtk, err := GetComponents3[comp.Position, comp.UnitRadius, comp.Attack](world, id)
					if err != nil {
						fmt.Printf("3 (structureCombatSearch - check_combat.go) -%v \n", err)
						return false
					}

					//get Unit Components
					ePos, eRadius, err := GetComponents2[comp.Position, comp.UnitRadius](world, uAtk.Target)
					if err != nil {
						fmt.Printf("2 (structureCombatSearch - check_combat.go) -%v \n", err)
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
					uPos, uRadius, uAtk, uTeam, MatchID, class, err := GetComponents6[comp.Position, comp.UnitRadius, comp.Attack, comp.Team, comp.MatchId, comp.Class](world, id)
					if err != nil {
						fmt.Printf("5 not in combat (structureCombatSearch - check_combat.go) -%v \n", err)
						return false
					}

					// get collision Hash
					collisionHash, err := getCollisionHashGSS(world, MatchID)
					if err != nil {
						fmt.Printf("error retrieving SpartialHash component on tempSpartialHash (structureCombatSearch - check_combat.go): %s  \n", err)
						return false
					}
					//find closest enemy
					eID, eX, eY, eRadius, found := findClosestEnemy(collisionHash, id, uPos.PositionVectorX, uPos.PositionVectorY, uAtk.AggroRadius, uTeam.Team, class.Class, true)
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
func findClosestEnemy(hash *comp.SpatialHash, objID types.EntityID, startX, startY float32, attackRadius int, team, class string, targetStruct bool) (types.EntityID, float32, float32, int, bool) {
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

						if !targetStruct && cell.Type[i] == "structure" { //if cannot target structures and target is a strutures
							continue
						}

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

// /for attack phase.  unit is attacking structure but cannot use sp on structures
// this functions finds closest targetable enemy and sets sp combat and target
func findClosestEnemySP(world cardinal.WorldContext, id types.EntityID, sp *comp.Sp) (bool, error) {
	//get Unit Components
	uPos, uRadius, uAtk, uTeam, MatchID, class, err := GetComponents6[comp.Position, comp.UnitRadius, comp.Attack, comp.Team, comp.MatchId, comp.Class](world, id)
	if err != nil {
		return false, fmt.Errorf("(unitCombatSearch - check_combat.go) -%v ", err)
	}

	// get collision Hash
	collisionHash, err := getCollisionHashGSS(world, MatchID)
	if err != nil {
		return false, fmt.Errorf("error retrieving SpartialHash component on tempSpartialHash (cunitCombatSearch - check_combat.go): %s ", err)
	}

	//find closest enemy
	eID, eX, eY, eRadius, found := findClosestEnemy(collisionHash, id, uPos.PositionVectorX, uPos.PositionVectorY, uAtk.AggroRadius, uTeam.Team, class.Class, false)
	if found { //found enemy
		// Calculate squared distance between the unit and the enemy, minus their radii
		adjustedDistance := distanceBetweenTwoPoints(uPos.PositionVectorX, uPos.PositionVectorY, eX, eY) - float32(eRadius) - float32(uRadius.UnitRadius)
		//if within attack range
		if adjustedDistance <= float32(sp.AttackRadius) {
			sp.Combat = true
			sp.Target = eID
			return true, nil
		}
	}
	return false, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////Default Combat Reset Functions////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func resetCombat(world cardinal.WorldContext, id types.EntityID, attack *comp.Attack) error {

	//get special power component
	sp, err := cardinal.GetComponent[comp.Sp](world, id)
	if err != nil {
		fmt.Printf("error retrieving special power comp (ResetCombat/check combat.go): \n")
		return nil
	}

	if !sp.Combat { //if unit is in sp dont reset attack
		attack.Frame = 0
	}
	attack.Target = 0
	attack.Combat = false

	return nil

}

// overwrite base destruction.
// if unit being attacked by channel dies don't cancel attack.
func channelingResetCombat(world cardinal.WorldContext, id types.EntityID, attack *comp.Attack) error {
	//reset attack component

	//get special power component
	sp, err := cardinal.GetComponent[comp.Sp](world, id)
	if err != nil {
		fmt.Printf("error retrieving special power comp (channelingResetCombat/check combat.go): \n")
		return nil
	}

	if attack.Frame < sp.DamageFrame && sp.Charged || !sp.Charged { //if target dies b4 fire attack goes off
		//reset units combat
		attack.Frame = 0
		attack.Combat = false
		attack.Target = 0
		attack.State = "Default"
	} else { //if unit started channeling fire
		attack.State = "Channeling"
		attack.Target = 0
		// attack.Target = id //set target to self to not get errors if triggering functions that ref this but unit is dead
	}

	return nil
}
