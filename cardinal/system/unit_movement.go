package system

import (
	"fmt"
	"sort"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"

	comp "MobaClashRoyal/component"
)

// This function is called every tick automatically
// It updates the units position based many factors
// out of combat it follows a direction map
// in aggro range moves towards enemy
// can push and walk around units in the way
func UnitMovementSystem(world cardinal.WorldContext) error {
	//get all Unit Id's in priority of distance to base
	priorityUnitIDs, err := PriorityUnitMovement(world)
	if err != nil {
		return fmt.Errorf("(unit_movement.go) -  %v", err)
	}

	//go through all Unit ID's
	for _, id := range priorityUnitIDs {
		//get Unit Components
		uPos, uRadius, uAtk, uTeam, uMs, MatchID, mapName, err := GetUnitComponentsUM(world, id)
		if err != nil {
			fmt.Printf("%v", err)
			continue
		}

		//get game state
		gameState, err := getGameStateGSS(world, MatchID)
		if err != nil {
			fmt.Printf("%v", err)
			continue
		}

		//get collision Hash
		collisionHash, err := cardinal.GetComponent[comp.SpatialHash](world, gameState)
		if err != nil {
			fmt.Printf("error retrieving SpartialHash component on tempSpartialHash (unit_movement.go): %s", err)
			continue
		}

		secondIfCondition := true

		//if units in combat
		if uAtk.Combat {
			//get enemyID  from unit target
			enemyID := uAtk.Target
			ePos, eRadius, err := getTargetComponentsUM(world, enemyID) //get enemy position and radius components
			if err != nil {
				fmt.Printf("(unit_movement.go): %s\n", err)
				continue
			}
			//distance between unit and enemy minus their radius
			adjustedDistance := distanceBetweenTwoPoints(uPos.PositionVectorX, uPos.PositionVectorY, ePos.PositionVectorX, ePos.PositionVectorY) - float32(eRadius.UnitRadius) - float32(uRadius.UnitRadius)

			//if out of attack range but in aggro range
			if adjustedDistance > float32(uAtk.AttackRadius) && adjustedDistance <= float32(uAtk.AggroRadius) {
				uAtk.Combat = false
				uAtk.Frame = 0
				secondIfCondition = false //not in combat but need to make sure not moving with direction map
				//set attack component
				if err = cardinal.SetComponent(world, id, uAtk); err != nil {
					fmt.Printf("error setting attack component (unit_movement.go): %v", err)
					continue
				}
				//move towards enemy in combat with
				if uMs.CurrentMS > 0 {
					tempX := uPos.PositionVectorX //Store Original X and Y
					tempY := uPos.PositionVectorY
					//move towards enemy
					uPos = MoveUnitTowardsEnemyUM(uPos, ePos.PositionVectorX, ePos.PositionVectorY, eRadius.UnitRadius, uMs, uRadius)
					//attempt to push blocking units
					pushBlockingUnit(world, collisionHash, id, tempX, tempY, uPos.PositionVectorX, uPos.PositionVectorY, uRadius.UnitRadius, uTeam.Team, uMs.CurrentMS)
					uPos.PositionVectorX, uPos.PositionVectorY = walkAroundUnit(collisionHash, id, tempX, tempY, uPos.PositionVectorX, uPos.PositionVectorY, uRadius.UnitRadius, uTeam.Team)
					// Set the new position component
					if err := cardinal.SetComponent(world, id, uPos); err != nil {
						fmt.Printf("error set component on position (unit movement.go): %v", err)
						continue
					}
					// Update units new distance from enemy base
					if err = UpdateUnitDistance(world, id, uTeam, uPos, mapName); err != nil {
						fmt.Printf("%v", err)
						continue
					}
				}
				//if out of both attack and aggro range
			} else if adjustedDistance > float32(uAtk.AggroRadius) {
				uAtk.Combat = false
				uAtk.Frame = 0
				//set attack component
				if err = cardinal.SetComponent(world, id, uAtk); err != nil {
					fmt.Printf("error setting attack component (unit_movement.go): %v", err)
					continue
				}
				//in attack range just rotate towards enemy
			} else {
				// Compute direction vector towards the enemy
				uPos.RotationVectorX, uPos.RotationVectorY = directionVectorBetweenTwoPoints(uPos.PositionVectorX, uPos.PositionVectorY, ePos.PositionVectorX, ePos.PositionVectorY)
				// Set the new position component
				if err := cardinal.SetComponent(world, id, uPos); err != nil {
					fmt.Printf("error set component on tempPosition (unit movement/MoveUnitTowardsEnemyUM): %v", err)
					continue
				}
			}

		}
		//if units not in combat
		if !uAtk.Combat && secondIfCondition {
			//Check for in range Enemies
			eID, eX, eY, eRadius, found := FindClosestEnemySpatialHash(collisionHash, id, uPos.PositionVectorX, uPos.PositionVectorY, uAtk.AggroRadius, uTeam.Team)
			if found { //found enemy
				// Calculate squared distance between the unit and the enemy, minus their radii
				adjustedDistance := distanceBetweenTwoPoints(uPos.PositionVectorX, uPos.PositionVectorY, eX, eY) - float32(eRadius) - float32(uRadius.UnitRadius)
				//if within attack range
				if adjustedDistance <= float32(uAtk.AttackRadius) {
					uAtk.Combat = true
					uAtk.Target = eID
					//set attack component
					if err = cardinal.SetComponent(world, id, uAtk); err != nil {
						fmt.Printf("error setting attack component (unit_movement.go): %v", err)
						continue
					}
					//not within attack range
				} else {
					if uMs.CurrentMS > 0 { // move towards enemy
						// //Store Original X and Y
						tempX := uPos.PositionVectorX
						tempY := uPos.PositionVectorY
						uPos = MoveUnitTowardsEnemyUM(uPos, eX, eY, eRadius, uMs, uRadius)
						//attempt to push blocking units
						pushBlockingUnit(world, collisionHash, id, tempX, tempY, uPos.PositionVectorX, uPos.PositionVectorY, uRadius.UnitRadius, uTeam.Team, uMs.CurrentMS)
						uPos.PositionVectorX, uPos.PositionVectorY = walkAroundUnit(collisionHash, id, tempX, tempY, uPos.PositionVectorX, uPos.PositionVectorY, uRadius.UnitRadius, uTeam.Team)
						// Set the new position component
						err := cardinal.SetComponent(world, id, uPos)
						if err != nil {
							fmt.Printf("error set component on tempPosition (unit movement.go): %v", err)
							continue
						}
						if err = UpdateUnitDistance(world, id, uTeam, uPos, mapName); err != nil {
							fmt.Printf("%v", err)
							continue
						}
					}
				}
			} else {
				//no enemies found and not in combat, move with direction map.
				if uMs.CurrentMS > 0 {
					// //Store Original X and Y
					tempX := uPos.PositionVectorX
					tempY := uPos.PositionVectorY
					uPos, err = MoveUnitDirectionMapUM(uPos, uTeam, uMs.CurrentMS, mapName)
					if err != nil {
						fmt.Printf("(unit_movement.go): %v", err)
						continue
					}
					//attempt to push blocking units
					pushBlockingUnit(world, collisionHash, id, tempX, tempY, uPos.PositionVectorX, uPos.PositionVectorY, uRadius.UnitRadius, uTeam.Team, uMs.CurrentMS)
					uPos.PositionVectorX, uPos.PositionVectorY = walkAroundUnit(collisionHash, id, tempX, tempY, uPos.PositionVectorX, uPos.PositionVectorY, uRadius.UnitRadius, uTeam.Team)
					//set updated position component
					err := cardinal.SetComponent(world, id, uPos)
					if err != nil {
						fmt.Printf("error set component on tempPosition (unit movement/MoveUnitDirectionMapUM): %v", err)
						continue
					}

					//update units new distance from enemy base
					if err = UpdateUnitDistance(world, id, uTeam, uPos, mapName); err != nil {
						fmt.Printf("(unit_movement.go): %v", err)
						continue
					}
				}
			}
		}
	}
	return err
}

// orders all units from distance from enemy base
// this way units infront move before units behind them to avoid walking around units that should not be in the way if they moved first
func PriorityUnitMovement(world cardinal.WorldContext) ([]types.EntityID, error) {
	// UnitData struct to store both the EntityID and its Distance for sorting
	type UnitData struct {
		ID       types.EntityID
		Distance float32
	}

	// Search all units
	unitList, err := cardinal.NewSearch().Entity(
		filter.Exact(UnitFilters())).Collect(world)
	if err != nil {
		return nil, fmt.Errorf("PriorityUnitMovement error searching for unit with map (priorityUnitMovement): %w", err)
	}

	// Create a slice to store the units with their distances
	unitsData := make([]UnitData, 0, len(unitList))

	// Fetch distances for each unit
	for _, unit := range unitList {
		distanceComp, err := cardinal.GetComponent[comp.Distance](world, unit)
		if err != nil {
			return nil, fmt.Errorf("error fetching distance for unit %v: %w (priorityUnitMovement)", unit, err)
		}
		// add to list
		unitsData = append(unitsData, UnitData{ID: unit, Distance: distanceComp.Distance})
	}

	// Sort units by Distance component
	sort.Slice(unitsData, func(i, j int) bool {
		return unitsData[i].Distance < unitsData[j].Distance
	})

	// Extract sorted IDs to return
	sortedIDs := make([]types.EntityID, len(unitsData))
	for i, data := range unitsData {
		sortedIDs[i] = data.ID
	}
	return sortedIDs, nil
}

// Moves Unit in direction of the map Direction vector
func MoveUnitDirectionMapUM(position *comp.Position, team *comp.Team, movespeed float32, mapName *comp.MapName) (*comp.Position, error) {
	//check map data exsists
	mapData, exists := MapDataRegistry[mapName.MapName]
	if !exists {
		return nil, fmt.Errorf("error key for MapDataRegistry does not exsist (MoveUnitDirectionMapUM)")
	}
	//check direction map exsists
	mapDir, ok := MapRegistry[mapName.MapName]
	if !ok {
		return nil, fmt.Errorf("error key for MapRegistry does not exsist (MoveUnitDirectionMapUM)")
	}

	//normalize the units position to the maps grid increments.
	normalizedX := int(((int(position.PositionVectorX)-mapData.StartX)/mapData.Increment))*mapData.Increment + mapData.StartX
	normalizedY := int(((int(position.PositionVectorY)-mapData.StartY)/mapData.Increment))*mapData.Increment + mapData.StartY
	//The units (x,y) coordinates normalized and turned into proper key(string) format for seaching map
	coordKey := fmt.Sprintf("%d,%d", normalizedX, normalizedY)

	// Retrieve direction vector using coordinate key
	directionVector, exists := mapDir.DMap[coordKey]
	if !exists {
		return nil, fmt.Errorf("no direction vector found for the given coordinates (MoveUnitDirectionMapUM)")
	}
	//updated rotation based on team
	if team.Team == "Blue" {
		position.RotationVectorX = directionVector[0]
		position.RotationVectorY = directionVector[1]
	} else {
		position.RotationVectorX = directionVector[0] * -1 //reverse direction for red
		position.RotationVectorY = directionVector[1] * -1
	}

	//update new x,y based on movespeed
	position.PositionVectorX = position.PositionVectorX + (position.RotationVectorX * movespeed)
	position.PositionVectorY = position.PositionVectorY + (position.RotationVectorY * movespeed)

	return position, nil
}

// Moves Unit towards enemy position
func MoveUnitTowardsEnemyUM(position *comp.Position, enemyX float32, enemyY float32, enemyRadius int, movespeed *comp.Movespeed, radius *comp.UnitRadius) *comp.Position {
	// Compute direction vector towards the enemy
	position.RotationVectorX, position.RotationVectorY = directionVectorBetweenTwoPoints(position.PositionVectorX, position.PositionVectorY, enemyX, enemyY)

	// Compute new position based on movespeed and direction
	position.PositionVectorX = position.PositionVectorX + position.RotationVectorX*movespeed.CurrentMS
	position.PositionVectorY = position.PositionVectorY + position.RotationVectorY*movespeed.CurrentMS

	// Calculate the stopping distance (combined radii of the unit and enemy plus 1 pixel for separation)
	stoppingDistance := radius.UnitRadius + enemyRadius + 1
	// Calculate the target position to move towards, stopping 1 pixel outside the enemy's radius
	targetX := enemyX - position.RotationVectorX*float32(stoppingDistance)
	targetY := enemyY - position.RotationVectorY*float32(stoppingDistance)

	// Ensure the unit does not overshoot the target position
	if (position.RotationVectorX > 0 && position.PositionVectorX > targetX) || (position.RotationVectorX < 0 && position.PositionVectorX < targetX) {
		position.PositionVectorX = targetX
	}
	if (position.RotationVectorY > 0 && position.PositionVectorY > targetY) || (position.RotationVectorY < 0 && position.PositionVectorY < targetY) {
		position.PositionVectorY = targetY
	}

	return position
}

// attempts to push the blocking unit
func pushBlockingUnit(world cardinal.WorldContext, hash *comp.SpatialHash, objID types.EntityID, startX, startY, targetX, targetY float32, radius int, team string, distance float32) {
	//list of all units colliding with at target position
	collisionList := CheckCollisionSpatialHashList(hash, targetX, targetY, radius)

	//try to push blocking units
	if len(collisionList) > 0 {
		for _, collisionID := range collisionList {
			//skip if collides with self
			if collisionID == objID {
				continue
			}
			//get targets team
			targetTeam, err := cardinal.GetComponent[comp.Team](world, collisionID)
			if err != nil {
				fmt.Printf("error getting targets team compoenent (UpdateUnitPositionPushSpatialHash): %v", err)
				continue
			}
			//if unit is ally push
			if targetTeam.Team == team {
				//get targets posisiton and radius components
				targetPos, targetRadius, err := getTargetComponentsUM(world, collisionID)
				if err != nil {
					fmt.Printf("(UpdateUnitPositionPushSpatialHash): %v", err)
					continue
				}

				// Remove the object from its current position in collision hash
				RemoveObjectFromSpatialHash(hash, collisionID, targetPos.PositionVectorX, targetPos.PositionVectorY, targetRadius.UnitRadius)
				//location to push unit to
				newTargetX, newTargetY := PushUnitDirSpatialHash(collisionID, targetX, targetY, targetPos.PositionVectorX, targetPos.PositionVectorY, startX-targetX, startY-targetY, distance)

				targetPos.PositionVectorX, targetPos.PositionVectorY = pushTowardsEnemySpatialHash(world, hash, collisionID, targetPos.PositionVectorX, targetPos.PositionVectorY, newTargetX, newTargetY, targetRadius.UnitRadius, distance, targetTeam)

				// Add the objects position to collosion hash
				AddObjectSpatialHash(hash, collisionID, targetPos.PositionVectorX, targetPos.PositionVectorY, targetRadius.UnitRadius, targetTeam.Team)
				//set collided units new position component
				if err = cardinal.SetComponent(world, collisionID, targetPos); err != nil {
					fmt.Printf("error setting target pos component (UpdateUnitPositionPushSpatialHash): %v", err)
					continue
				}
			}
		}
	}

}

// walks around blocking unit if exsists
func walkAroundUnit(hash *comp.SpatialHash, objID types.EntityID, startX, startY, targetX, targetY float32, radius int, team string) (newtargetX, newtargetY float32) {
	// Remove the object from its current position
	RemoveObjectFromSpatialHash(hash, objID, startX, startY, radius)
	// Find an alternative position if the target is occupied
	if CheckCollisionSpatialHash(hash, targetX, targetY, radius) {
		//walk around blocking unit
		targetX, targetY = moveToNearestFreeSpaceBoxSpatialHash(hash, startX, startY, targetX, targetY, float32(radius))
	}
	// Add the object to the new position
	AddObjectSpatialHash(hash, objID, targetX, targetY, radius, team)
	return targetX, targetY
}

// Update units distance from enemy base to help with movement priority queue
func UpdateUnitDistance(world cardinal.WorldContext, id types.EntityID, team *comp.Team, position *comp.Position, mapName *comp.MapName) error {
	//check map exsists in registy
	mapData, exists := MapDataRegistry[mapName.MapName]
	if !exists {
		return fmt.Errorf("error key for MapDataRegistry does not exsist (UpdateUnitDistance)")
	}

	//find distance from enemy base and update component
	err := cardinal.UpdateComponent(world, id, func(distance *comp.Distance) *comp.Distance {
		if distance == nil {
			fmt.Printf("error retrieving distance component (UpdateUnitDistance)")
			return nil
		}
		// calculate distance from enemy spawn
		if team.Team == "Blue" {
			distance.Distance = distanceBetweenTwoPoints(float32(mapData.Bases[1][0]), float32(mapData.Bases[1][1]), position.PositionVectorX, position.PositionVectorY)
		} else {
			distance.Distance = distanceBetweenTwoPoints(float32(mapData.Bases[0][0]), float32(mapData.Bases[0][1]), position.PositionVectorX, position.PositionVectorY)
		}
		return distance
	})

	if err != nil {
		return fmt.Errorf("error updating distance (UpdateUnitDistance): %w", err)
	}
	return nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////UTILITY FUNCTIONS//////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// GetUnitComponents fetches all necessary components related to a unit entity.
func GetUnitComponentsUM(world cardinal.WorldContext, unitID types.EntityID) (*comp.Position, *comp.UnitRadius, *comp.Attack, *comp.Team, *comp.Movespeed, *comp.MatchId, *comp.MapName, error) {
	position, err := cardinal.GetComponent[comp.Position](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Position component (unit_movement.go): %v", err)
	}
	unitRadius, err := cardinal.GetComponent[comp.UnitRadius](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Unit Radius component (unit_movement.go): %v", err)
	}
	unitAttack, err := cardinal.GetComponent[comp.Attack](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Unit attack component (unit_movement.go): %v", err)
	}
	team, err := cardinal.GetComponent[comp.Team](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Team component (unit_movement.go): %v", err)
	}
	movespeed, err := cardinal.GetComponent[comp.Movespeed](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Movespeed component (unit_movement.go): %v", err)
	}
	matchId, err := cardinal.GetComponent[comp.MatchId](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving MatchId component (unit_movement.go): %v", err)
	}
	mapName, err := cardinal.GetComponent[comp.MapName](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Distance component (unit_movement.go): %v", err)
	}
	return position, unitRadius, unitAttack, team, movespeed, matchId, mapName, nil
}

// fetches target components
func getTargetComponentsUM(world cardinal.WorldContext, enemyID types.EntityID) (enemyPosition *comp.Position, enemyRadius *comp.UnitRadius, err error) {

	enemyPosition, err = cardinal.GetComponent[comp.Position](world, enemyID)
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving enemy Position component (unit_movement.go): %v", err)
	}
	enemyRadius, err = cardinal.GetComponent[comp.UnitRadius](world, enemyID)
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving enemy Radius component (unit_movement.go): %v", err)
	}
	return enemyPosition, enemyRadius, nil
}
