package system

import (
	"fmt"
	"math"
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
	priorityUnitIDs, err := priorityUnitMovement(world)
	if err != nil {
		return fmt.Errorf("(unit_movement.go) -  %v", err)
	}

	//go through all Unit ID's
	for _, id := range priorityUnitIDs {
		//get Unit CC component
		cc, err := cardinal.GetComponent[comp.CC](world, id)
		if err != nil {
			fmt.Printf("error getting unit cc component (unit_movement.go): %v \n", err)
			continue
		}

		if cc.Stun > 0 { //if unit stunned cannot move
			continue
		}

		//get Unit Components
		uPos, uRadius, uAtk, uTeam, uMs, MatchID, mapName, class, err := GetComponents8[comp.Position, comp.UnitRadius, comp.Attack, comp.Team, comp.Movespeed, comp.MatchId, comp.MapName, comp.Class](world, id)
		if err != nil {
			fmt.Printf("unit components (unit_movement.go) %v \n", err)
			continue
		}

		if uAtk.State == "Channeling" { //if unit chenneling cannot move
			continue
		}

		gameState, err := getGameStateGSS(world, MatchID)
		if err != nil {
			fmt.Printf("error retrieving gamestate id (unit_movement.go): %s \n", err)
			continue
		}
		//get collision Hash
		collisionHash, err := getCollisionHashGSS(world, MatchID)
		if err != nil {
			fmt.Printf("error retrieving SpartialHash component on tempSpartialHash (unit_movement.go): %s \n", err)
			continue
		}

		secondIfCondition := true

		//if units in combat
		if uAtk.Combat {
			//get enemy position and radius components
			ePos, eRadius, err := GetComponents2[comp.Position, comp.UnitRadius](world, uAtk.Target)
			if err != nil {
				fmt.Printf("combat compoenents (unit_movement.go): %s \n", err)
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
					fmt.Printf("error setting attack component (unit_movement.go): %v \n", err)
					continue
				}
				//move towards enemy in combat with
				if uMs.CurrentMS > 0 {
					tempX := uPos.PositionVectorX //Store Original X and Y
					tempY := uPos.PositionVectorY
					//move towards enemy
					uPos = moveUnitTowardsEnemy(uPos, ePos.PositionVectorX, ePos.PositionVectorY, eRadius.UnitRadius, uMs.CurrentMS, uRadius.UnitRadius)
					//check that unit isnt walking through out of bounds towards a found unit
					exists := moveDirectionExsist(uPos.PositionVectorX, uPos.PositionVectorY, mapName.MapName)
					if exists {
						//attempt to push blocking units
						pushBlockingUnit(world, collisionHash, id, uPos.PositionVectorX, uPos.PositionVectorY, uRadius.UnitRadius, class.Class, uTeam.Team, uMs.CurrentMS, mapName)
						//move unit.  walk around blocking units
						uPos.PositionVectorX, uPos.PositionVectorY = moveFreeSpace(collisionHash, id, tempX, tempY, uPos.PositionVectorX, uPos.PositionVectorY, uRadius.UnitRadius, uTeam.Team, class.Class, mapName)
						// Set the new position component
						if err := cardinal.SetComponent(world, id, uPos); err != nil {
							fmt.Printf("error set component on position (unit movement.go): %v \n", err)
							continue
						}
						// Update units new distance from enemy base
						if err = updateUnitDistance(world, id, uTeam, uPos, mapName); err != nil {
							fmt.Printf("%v", err)
							continue
						}
					} else {
						//no unit found because requires walking through map
						uPos.PositionVectorX = tempX
						uPos.PositionVectorY = tempY
						//move with direction map
						secondIfCondition = true
					}
				}
				//if out of both attack and aggro range
			} else if adjustedDistance > float32(uAtk.AggroRadius) {
				uAtk.Combat = false
				uAtk.Frame = 0
				//set attack component
				if err = cardinal.SetComponent(world, id, uAtk); err != nil {
					fmt.Printf("error setting attack component (unit_movement.go): %v \n", err)
					continue
				}
				//in attack range just rotate towards enemy
			} else {
				// Compute direction vector towards the enemy
				uPos.RotationVectorX, uPos.RotationVectorY = directionVectorBetweenTwoPoints(uPos.PositionVectorX, uPos.PositionVectorY, ePos.PositionVectorX, ePos.PositionVectorY)
				// Set the new position component
				if err := cardinal.SetComponent(world, id, uPos); err != nil {
					fmt.Printf("error set component on tempPosition (unit movement/MoveUnitTowardsEnemyUM): %v \n", err)
					continue
				}
			}

		}
		//if units not in combat
		if !uAtk.Combat && secondIfCondition {
			//Check for in range Enemies
			eID, eX, eY, eRadius, found := findClosestEnemy(collisionHash, id, uPos.PositionVectorX, uPos.PositionVectorY, uAtk.AggroRadius, uTeam.Team, class.Class)
			if found { //found enemy
				// Calculate squared distance between the unit and the enemy, minus their radii
				adjustedDistance := distanceBetweenTwoPoints(uPos.PositionVectorX, uPos.PositionVectorY, eX, eY) - float32(eRadius) - float32(uRadius.UnitRadius)
				//if within attack range
				if adjustedDistance <= float32(uAtk.AttackRadius) {
					uAtk.Combat = true
					uAtk.Target = eID
					//set attack component
					if err = cardinal.SetComponent(world, id, uAtk); err != nil {
						fmt.Printf("error setting attack component (unit_movement.go): %v \n", err)
						continue
					}
					//not within attack range
				} else {
					if uMs.CurrentMS > 0 { // move towards enemy
						// //Store Original X and Y
						tempX := uPos.PositionVectorX
						tempY := uPos.PositionVectorY
						uPos = moveUnitTowardsEnemy(uPos, eX, eY, eRadius, uMs.CurrentMS, uRadius.UnitRadius)
						//check that unit isnt walking through out of bounds towards a found unit
						exists := moveDirectionExsist(uPos.PositionVectorX, uPos.PositionVectorY, mapName.MapName)
						if exists {
							//attempt to push blocking units
							pushBlockingUnit(world, collisionHash, id, uPos.PositionVectorX, uPos.PositionVectorY, uRadius.UnitRadius, uTeam.Team, class.Class, uMs.CurrentMS, mapName)
							//move unit.  walk around blocking units
							uPos.PositionVectorX, uPos.PositionVectorY = moveFreeSpace(collisionHash, id, tempX, tempY, uPos.PositionVectorX, uPos.PositionVectorY, uRadius.UnitRadius, uTeam.Team, class.Class, mapName)
							// Set the new position component
							err := cardinal.SetComponent(world, id, uPos)
							if err != nil {
								fmt.Printf("error set component on tempPosition (unit movement.go): %v \n", err)
								continue
							}
							if err = updateUnitDistance(world, id, uTeam, uPos, mapName); err != nil {
								fmt.Printf("%v", err)
								continue
							}
						} else {
							//no unit found because requires walking through map
							found = false
							uPos.PositionVectorX = tempX
							uPos.PositionVectorY = tempY
						}
					}
				}
			}
			if !found {
				//no enemies found and not in combat, move with direction map.
				if uMs.CurrentMS > 0 {
					// //Store Original X and Y
					tempX := uPos.PositionVectorX
					tempY := uPos.PositionVectorY
					uPos, err = moveUnitDirectionMap(uPos, uTeam, uMs.CurrentMS, mapName)
					if err != nil {
						fmt.Printf("(unit_movement.go): %v \n", err)
						continue
					}
					//attempt to push blocking units
					pushBlockingUnit(world, collisionHash, id, uPos.PositionVectorX, uPos.PositionVectorY, uRadius.UnitRadius, uTeam.Team, class.Class, uMs.CurrentMS, mapName)
					//move unit.  walk around blocking units
					uPos.PositionVectorX, uPos.PositionVectorY = moveFreeSpace(collisionHash, id, tempX, tempY, uPos.PositionVectorX, uPos.PositionVectorY, uRadius.UnitRadius, uTeam.Team, class.Class, mapName)
					//set updated position component
					err := cardinal.SetComponent(world, id, uPos)
					if err != nil {
						fmt.Printf("error set component on tempPosition (unit movement/MoveUnitDirectionMapUM): %v \n", err)
						continue
					}

					//update units new distance from enemy base
					if err = updateUnitDistance(world, id, uTeam, uPos, mapName); err != nil {
						fmt.Printf("(unit_movement.go): %v \n", err)
						continue
					}
				}
			}
		}
		//set collision Hash
		err = cardinal.SetComponent(world, gameState, collisionHash)
		if err != nil {
			fmt.Printf("error setting SpartialHash component (unit_movement.go): %s \n", err)
			continue
		}
	}
	return err
}

// orders all units from distance from enemy base
// this way units infront move before units behind them to avoid walking around units that should not be in the way if they moved first
func priorityUnitMovement(world cardinal.WorldContext) ([]types.EntityID, error) {
	// UnitData struct to store both the EntityID and its Distance for sorting
	type UnitData struct {
		ID       types.EntityID
		Distance float32
	}

	// Search all units
	unitList, err := cardinal.NewSearch().Entity(
		filter.Contains(filter.Component[comp.UnitTag]())).Collect(world)
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
func moveUnitDirectionMap(position *comp.Position, team *comp.Team, movespeed float32, mapName *comp.MapName) (*comp.Position, error) {

	//check if mapName exsists and if direction vector exsists at (x, y) location
	directionVector, err := getMapDirection(position.PositionVectorX, position.PositionVectorY, mapName.MapName)
	if err != nil {
		return nil, fmt.Errorf("(MoveUnitDirectionMap): %w", err)
	}

	var dirX, dirY float32
	//updated rotation based on team
	if team.Team == "Blue" {
		dirX = directionVector[0]
		dirY = directionVector[1]
	} else {
		dirX = directionVector[0] * -1 //reverse direction for red
		dirY = directionVector[1] * -1
	}

	//update new x,y based on movespeed
	tempX := position.PositionVectorX + (dirX * movespeed)
	tempY := position.PositionVectorY + (dirY * movespeed)
	tempDirX := dirX
	tempDirY := dirY

	angle := 5.0    // Degrees to rotate each step
	sumAngle := 0.0 // Sum of rotated angles

	for i := 0; ; i++ {
		if moveDirectionExsist(tempX, tempY, mapName.MapName) {
			//update new x,y based on movespeed
			position.PositionVectorX = tempX
			position.PositionVectorY = tempY
			position.RotationVectorX = tempDirX
			position.RotationVectorY = tempDirY
			break

		}
		// Rotate clockwise for even i, counterclockwise for odd i
		if i%2 == 0 {
			dirX, dirY = rotateVectorDegrees(tempDirX, tempDirY, angle)
			angle += angle
		} else {
			dirX, dirY = rotateVectorDegrees(tempDirX, tempDirY, -angle)
			angle += angle
		}
		tempX = position.PositionVectorX + (dirX * movespeed)
		tempY = position.PositionVectorY + (dirY * movespeed)

		// Check if the rotation has reached 180 degrees
		if math.Abs(sumAngle) >= 180 {
			break
		}
	}

	return position, nil
}

// Moves Unit towards enemy position
func moveUnitTowardsEnemy(position *comp.Position, enemyX float32, enemyY float32, enemyRadius int, movespeed float32, radius int) *comp.Position {
	// Compute direction vector towards the enemy
	position.RotationVectorX, position.RotationVectorY = directionVectorBetweenTwoPoints(position.PositionVectorX, position.PositionVectorY, enemyX, enemyY)

	// Compute new position based on movespeed and direction
	position.PositionVectorX = position.PositionVectorX + position.RotationVectorX*movespeed
	position.PositionVectorY = position.PositionVectorY + position.RotationVectorY*movespeed

	// Calculate the stopping distance (combined radii of the unit and enemy plus 1 pixel for separation)
	stoppingDistance := radius + enemyRadius + 1
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
func pushBlockingUnit(world cardinal.WorldContext, hash *comp.SpatialHash, objID types.EntityID, targetX, targetY float32, radius int, team, class string, distance float32, mapName *comp.MapName) {
	//list of all units colliding with at target position
	collisionList := CheckCollisionSpatialHashList(hash, targetX, targetY, radius, class, true)

	//try to push blocking units
	if len(collisionList) > 0 {
		for _, collisionID := range collisionList {
			//skip if collides with self
			if collisionID == objID {
				continue
			}

			//get targets attack
			tClass, err := cardinal.GetComponent[comp.Class](world, collisionID)
			if err != nil {
				fmt.Printf("error getting targets attack compoenent (pushBlockingUnit): %v \n", err)
				continue
			}

			if class == "air" && tClass.Class != "air" { //air can only push air
				continue
			}

			// get target components
			targetTeam, targetName, err := GetComponents2[comp.Team, comp.UnitName](world, collisionID)
			if err != nil {
				fmt.Printf("error getting targets compoenents (pushBlockingUnit): %v \n", err)
				continue
			}

			//if unit is ally push
			if targetTeam.Team == team && targetName.UnitName != "Base" && targetName.UnitName != "Tower" {
				//get targets posisiton and radius components
				targetPos, targetRadius, err := GetComponents2[comp.Position, comp.UnitRadius](world, collisionID)
				if err != nil {
					fmt.Printf("target compoenents (pushBlockingUnit): %v \n", err)
					continue
				}

				// Remove the object from its current position in collision hash
				RemoveObjectFromSpatialHash(hash, collisionID, targetPos.PositionVectorX, targetPos.PositionVectorY, targetRadius.UnitRadius)
				//location to push unit to
				newTargetX, newTargetY := pushUnitDirection(targetX, targetY, targetPos.PositionVectorX, targetPos.PositionVectorY, targetPos.RotationVectorX, targetPos.RotationVectorY, distance)
				//find closest non occupide location between target location and currnt location
				targetPos.PositionVectorX, targetPos.PositionVectorY = pushFromPtBtoA(world, hash, collisionID, targetPos.PositionVectorX, targetPos.PositionVectorY, newTargetX, newTargetY, targetRadius.UnitRadius, mapName)
				// Add the objects position to collosion hash

				AddObjectSpatialHash(hash, collisionID, targetPos.PositionVectorX, targetPos.PositionVectorY, targetRadius.UnitRadius, targetTeam.Team, tClass.Class)
				//set collided units new position component
				if err = cardinal.SetComponent(world, collisionID, targetPos); err != nil {
					fmt.Printf("error setting target pos component (pushBlockingUnit): %v \n", err)
					continue
				}
			}
		}
	}
}

// push a unit when they collide
// simulates if position 1 hits position 2 like a billards ball and bounces the position 2 to our new target returns
// Pos1 is hitting pos2
func pushUnitDirection(posX1, posY1, posX2, posY2, dirX2, dirY2, distance float32) (targetX, targetY float32) {
	//angle to move fowards if hit
	//(angle) * math.Pi / 180
	middleWidth := 0.26179938779 //angle of 15 degrees

	// Calculate the vector from the incoming ball to the main ball
	dirToBallX := posX2 - posX1
	dirToBallY := posY2 - posY1

	//normalize ball 2 and direction of impact
	normBallHitX, normBallHitY := normalize(dirX2, dirY2)
	dirToBallX, dirToBallY = normalize(dirToBallX, dirToBallY)

	// Calculate the dot product to find the angle
	dotProduct := dotProduct(dirToBallX, dirToBallY, normBallHitX, normBallHitY)

	// Calculate the cross product (for determining left/right side)
	crossProduct := crossProduct(dirToBallX, dirToBallY, normBallHitX, normBallHitY)

	// Calculate angle in degrees
	angle := math.Acos(float64(dotProduct))

	// Determine the direction based on the angle and the width for forward motion
	var dirX, dirY float32
	if math.Abs(angle) <= middleWidth/2 {
		// Forward push
		dirX, dirY = normBallHitX, normBallHitY
	} else if crossProduct < 0 {
		// Perpendicular push to the left
		dirX, dirY = -normBallHitY, normBallHitX
	} else {
		// Perpendicular push to the right
		dirX, dirY = normBallHitY, -normBallHitX
	}
	targetX = posX2 + dirX*(distance/4)
	targetY = posY2 + dirY*(distance/4)
	return targetX, targetY
}

// find closest free point between points B to A
func pushFromPtBtoA(world cardinal.WorldContext, hash *comp.SpatialHash, id types.EntityID, startX, startY, targetX, targetY float32, radius int, mapName *comp.MapName) (float32, float32) {
	//get attack component
	atk, class, err := GetComponents2[comp.Attack, comp.Class](world, id)
	if err != nil {
		fmt.Printf("(pushFromPtBtoA): %v \n", err)
		return startX, startY
	}

	//get length
	deltaX := targetX - startX
	deltaY := targetY - startY
	length := distanceBetweenTwoPoints(startX, startY, targetX, targetY)

	if length == 0 {
		fmt.Printf("length Dividing by 0 (pushFromPtBtoA)\n")
		return startX, startY
	}
	// Normalize direction vector
	dirX := deltaX / length
	dirY := deltaY / length

	// Step size, which can be adjusted as needed
	step := length / 8

	// if in combat must be a posisiton still in attack range of enemy
	if atk.Combat {
		//get target components
		targetPos, targetRadius, err := GetComponents2[comp.Position, comp.UnitRadius](world, atk.Target)
		if err != nil {
			fmt.Printf("target components (pushFromPtBtoA): %v \n", err)
			return startX, startY
		}

		// Search along the line from target to start (reverse)
		for d := float32(0); d <= length; d += float32(step) {
			testX := targetX + dirX*float32(d) // Start from target position
			testY := targetY + dirY*float32(d) // Start from target position
			//move unit towards even outside radius
			test := moveUnitTowardsEnemy(&comp.Position{PositionVectorX: testX, PositionVectorY: testY}, targetPos.PositionVectorX, targetPos.PositionVectorY, targetRadius.UnitRadius, length, radius)
			// Check if the position is free of collisions
			if !CheckCollisionSpatialHash(hash, test.PositionVectorX, test.PositionVectorY, radius, class.Class, true) && moveDirectionExsist(test.PositionVectorX, test.PositionVectorY, mapName.MapName) {
				return test.PositionVectorX, test.PositionVectorY // Return the first free spot found
			}
		}
	} else {
		//not in combat
		// Search along the line from target to start (reverse)
		for d := float32(0); d <= length; d += float32(step) {
			testX := targetX + dirX*float32(d) // Start from target position
			testY := targetY + dirY*float32(d) // Start from target position
			// Check if the position is free of collisions
			if !CheckCollisionSpatialHash(hash, testX, testY, radius, class.Class, true) {
				return testX, testY // Return the first free spot found
			}
		}
	}
	return startX, startY // Stay at the current position if no free spot is found
}

// walks around blocking unit if exsists to closest free space
func moveFreeSpace(hash *comp.SpatialHash, objID types.EntityID, startX, startY, targetX, targetY float32, radius int, team string, class string, mapName *comp.MapName) (float32, float32) {
	// Remove the object from its current position
	RemoveObjectFromSpatialHash(hash, objID, startX, startY, radius)
	// Find an alternative position if the target is occupied
	if CheckCollisionSpatialHash(hash, targetX, targetY, radius, class, true) {
		//walk around blocking unit
		targetX, targetY = moveToNearestFreeSpaceBox(hash, startX, startY, targetX, targetY, float32(radius), mapName, class)
	}
	// Add the object to the new position
	AddObjectSpatialHash(hash, objID, targetX, targetY, radius, team, class)
	return targetX, targetY
}

// creates a box between start and target point and checks all points from target to start finding first open position.
// box length is distance from start to target
// box width is radius*2
func moveToNearestFreeSpaceBox(hash *comp.SpatialHash, startX, startY, targetX, targetY, radius float32, mapName *comp.MapName, class string) (newX float32, newY float32) {
	//get length
	deltaX := targetX - startX
	deltaY := targetY - startY
	length := distanceBetweenTwoPoints(startX, startY, targetX, targetY)

	// Normalize direction vector
	dirX := deltaX / length
	dirY := deltaY / length

	// Perpendicular vector (normalized)
	perpX := -dirY
	perpY := dirX

	// Step size, which can be adjusted as needed
	step := float32(2) // or another division factor
	//length/8
	// Half the unit's radius
	halfWidth := radius / 2

	// Center to edge zigzag pattern
	maxOffset := int(halfWidth / step) // Number of steps from center to edge

	// Search within the square around the line from A to B
	for d := length; d >= -length; d -= step {
		// Alternate checking right and left of the center line
		for offset := 0; offset <= maxOffset; offset++ {
			offsets := []int{offset, -offset} // Check positive and negative offsets
			for _, w := range offsets {
				testX := startX + dirX*d + perpX*float32(w)*step
				testY := startY + dirY*d + perpY*float32(w)*step

				// Check if the position is free of collisions
				if !CheckCollisionSpatialHash(hash, testX, testY, int(radius), class, true) && moveDirectionExsist(testX, testY, mapName.MapName) {
					return testX, testY
				}
			}
		}
	}

	return startX, startY // Stay at the current position if no free spot is found
}

// Update units distance from enemy base to help with movement priority queue
func updateUnitDistance(world cardinal.WorldContext, id types.EntityID, team *comp.Team, position *comp.Position, mapName *comp.MapName) error {
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
