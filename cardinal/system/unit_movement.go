package system

import (
	"container/list"
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
		//get Unit Components
		uPos, uRadius, uAtk, uTeam, uMs, MatchID, mapName, err := getUnitComponentsUM(world, id)
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
					uPos = moveUnitTowardsEnemy(uPos, ePos.PositionVectorX, ePos.PositionVectorY, eRadius.UnitRadius, uMs.CurrentMS, uRadius.UnitRadius)
					//attempt to push blocking units
					pushBlockingUnit(world, collisionHash, id, uPos.PositionVectorX, uPos.PositionVectorY, uRadius.UnitRadius, uTeam.Team, uMs.CurrentMS)
					//move unit.  walk around blocking units
					uPos.PositionVectorX, uPos.PositionVectorY = moveFreeSpace(collisionHash, id, tempX, tempY, uPos.PositionVectorX, uPos.PositionVectorY, uRadius.UnitRadius, uTeam.Team)
					// Set the new position component
					if err := cardinal.SetComponent(world, id, uPos); err != nil {
						fmt.Printf("error set component on position (unit movement.go): %v", err)
						continue
					}
					// Update units new distance from enemy base
					if err = updateUnitDistance(world, id, uTeam, uPos, mapName); err != nil {
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
			eID, eX, eY, eRadius, found := findClosestEnemy(collisionHash, id, uPos.PositionVectorX, uPos.PositionVectorY, uAtk.AggroRadius, uTeam.Team)
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
						uPos = moveUnitTowardsEnemy(uPos, eX, eY, eRadius, uMs.CurrentMS, uRadius.UnitRadius)
						//attempt to push blocking units
						pushBlockingUnit(world, collisionHash, id, uPos.PositionVectorX, uPos.PositionVectorY, uRadius.UnitRadius, uTeam.Team, uMs.CurrentMS)
						//move unit.  walk around blocking units
						uPos.PositionVectorX, uPos.PositionVectorY = moveFreeSpace(collisionHash, id, tempX, tempY, uPos.PositionVectorX, uPos.PositionVectorY, uRadius.UnitRadius, uTeam.Team)
						// Set the new position component
						err := cardinal.SetComponent(world, id, uPos)
						if err != nil {
							fmt.Printf("error set component on tempPosition (unit movement.go): %v", err)
							continue
						}
						if err = updateUnitDistance(world, id, uTeam, uPos, mapName); err != nil {
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
					uPos, err = moveUnitDirectionMap(uPos, uTeam, uMs.CurrentMS, mapName)
					if err != nil {
						fmt.Printf("(unit_movement.go): %v", err)
						continue
					}
					//attempt to push blocking units
					pushBlockingUnit(world, collisionHash, id, uPos.PositionVectorX, uPos.PositionVectorY, uRadius.UnitRadius, uTeam.Team, uMs.CurrentMS)
					//move unit.  walk around blocking units
					uPos.PositionVectorX, uPos.PositionVectorY = moveFreeSpace(collisionHash, id, tempX, tempY, uPos.PositionVectorX, uPos.PositionVectorY, uRadius.UnitRadius, uTeam.Team)
					//set updated position component
					err := cardinal.SetComponent(world, id, uPos)
					if err != nil {
						fmt.Printf("error set component on tempPosition (unit movement/MoveUnitDirectionMapUM): %v", err)
						continue
					}

					//update units new distance from enemy base
					if err = updateUnitDistance(world, id, uTeam, uPos, mapName); err != nil {
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
func priorityUnitMovement(world cardinal.WorldContext) ([]types.EntityID, error) {
	// UnitData struct to store both the EntityID and its Distance for sorting
	type UnitData struct {
		ID       types.EntityID
		Distance float32
	}

	// Search all units
	unitList, err := cardinal.NewSearch().Entity(
		filter.Contains(UnitFilters())).Collect(world)
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
	//check map data exsists
	mapData, exists := MapDataRegistry[mapName.MapName]
	if !exists {
		return nil, fmt.Errorf("error key for MapDataRegistry does not exsist (MoveUnitDirectionMap)")
	}
	//check direction map exsists
	mapDir, ok := MapRegistry[mapName.MapName]
	if !ok {
		return nil, fmt.Errorf("error key for MapRegistry does not exsist (MoveUnitDirectionMap)")
	}

	//normalize the units position to the maps grid increments.
	normalizedX := int(((int(position.PositionVectorX)-mapData.StartX)/mapData.Increment))*mapData.Increment + mapData.StartX
	normalizedY := int(((int(position.PositionVectorY)-mapData.StartY)/mapData.Increment))*mapData.Increment + mapData.StartY
	//The units (x,y) coordinates normalized and turned into proper key(string) format for seaching map
	coordKey := fmt.Sprintf("%d,%d", normalizedX, normalizedY)

	// Retrieve direction vector using coordinate key
	directionVector, exists := mapDir.DMap[coordKey]
	if !exists {
		return nil, fmt.Errorf("no direction vector found for the given coordinates (MoveUnitDirectionMap)")
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
		if moveDirectionExsist(tempX, tempY, mapName) {
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

// Moves Unit in direction of the map Direction vector
func moveDirectionExsist(x, y float32, mapName *comp.MapName) bool {
	//check map data exsists
	mapData, exists := MapDataRegistry[mapName.MapName]
	if !exists {
		fmt.Printf("error key for MapDataRegistry does not exsist (MoveUnitDirectionMap)")
		return false
	}
	//check direction map exsists
	mapDir, ok := MapRegistry[mapName.MapName]
	if !ok {
		fmt.Printf("error key for MapRegistry does not exsist (MoveUnitDirectionMap)")
		return false
	}

	//normalize the units position to the maps grid increments.
	normalizedX := int(((int(x)-mapData.StartX)/mapData.Increment))*mapData.Increment + mapData.StartX
	normalizedY := int(((int(y)-mapData.StartY)/mapData.Increment))*mapData.Increment + mapData.StartY
	//The units (x,y) coordinates normalized and turned into proper key(string) format for seaching map
	coordKey := fmt.Sprintf("%d,%d", normalizedX, normalizedY)

	// Retrieve direction vector using coordinate key
	_, exists = mapDir.DMap[coordKey]
	if !exists {
		fmt.Printf("no direction vector found for the given coordinates (MoveUnitDirectionMap)")
		return false
	}

	return true
}

// FindClosestEnemy performs a BFS search from the unit's position outward within the attack radius.
func findClosestEnemy(hash *comp.SpatialHash, objID types.EntityID, startX, startY float32, attackRadius int, team string) (types.EntityID, float32, float32, int, bool) {
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
func pushBlockingUnit(world cardinal.WorldContext, hash *comp.SpatialHash, objID types.EntityID, targetX, targetY float32, radius int, team string, distance float32) {
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
				fmt.Printf("error getting targets team compoenent (pushBlockingUnit): %v", err)
				continue
			}
			//if unit is ally push
			if targetTeam.Team == team {
				//get targets posisiton and radius components
				targetPos, targetRadius, err := getTargetComponentsUM(world, collisionID)
				if err != nil {
					fmt.Printf("(pushBlockingUnit): %v", err)
					continue
				}

				// Remove the object from its current position in collision hash
				RemoveObjectFromSpatialHash(hash, collisionID, targetPos.PositionVectorX, targetPos.PositionVectorY, targetRadius.UnitRadius)
				//location to push unit to
				//newTargetX, newTargetY := pushUnitDirection(targetX, targetY, targetPos.PositionVectorX, targetPos.PositionVectorY, startX-targetX, startY-targetY, distance)
				newTargetX, newTargetY := pushUnitDirection(targetX, targetY, targetPos.PositionVectorX, targetPos.PositionVectorY, targetPos.RotationVectorX, targetPos.RotationVectorY, distance)
				//find closest non occupide location between target location and currnt location
				targetPos.PositionVectorX, targetPos.PositionVectorY = pushFromPtBtoA(world, hash, collisionID, targetPos.PositionVectorX, targetPos.PositionVectorY, newTargetX, newTargetY, targetRadius.UnitRadius)
				// Add the objects position to collosion hash

				AddObjectSpatialHash(hash, collisionID, targetPos.PositionVectorX, targetPos.PositionVectorY, targetRadius.UnitRadius, targetTeam.Team)
				//set collided units new position component
				if err = cardinal.SetComponent(world, collisionID, targetPos); err != nil {
					fmt.Printf("error setting target pos component (pushBlockingUnit): %v", err)
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
	middleWidth := 0.3490658503988659 //angle of 20 degrees

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
	targetX = posX2 + dirX*distance
	targetY = posY2 + dirY*distance
	return targetX, targetY
}

// find closest free point between points B to A
func pushFromPtBtoA(world cardinal.WorldContext, hash *comp.SpatialHash, id types.EntityID, startX, startY, targetX, targetY float32, radius int) (float32, float32) {
	//get attack component
	atk, err := cardinal.GetComponent[comp.Attack](world, id)
	if err != nil {
		fmt.Printf("error getting attack compoenent (pushFromPtBtoA): %v", err)
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
		targetPos, targetRadius, err := getTargetComponentsUM(world, atk.Target)
		if err != nil {
			fmt.Printf("(pushFromPtBtoA): %v", err)
			return startX, startY
		}

		// Search along the line from target to start (reverse)
		for d := float32(0); d <= length; d += float32(step) {
			testX := targetX + dirX*float32(d) // Start from target position
			testY := targetY + dirY*float32(d) // Start from target position
			//move unit towards even outside radius
			test := moveUnitTowardsEnemy(&comp.Position{PositionVectorX: testX, PositionVectorY: testY}, targetPos.PositionVectorX, targetPos.PositionVectorY, targetRadius.UnitRadius, length, radius)
			// Check if the position is free of collisions
			if !CheckCollisionSpatialHash(hash, test.PositionVectorX, test.PositionVectorY, radius) {
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
			if !CheckCollisionSpatialHash(hash, testX, testY, radius) {
				return testX, testY // Return the first free spot found
			}
		}
	}
	return startX, startY // Stay at the current position if no free spot is found
}

// walks around blocking unit if exsists to closest free space
func moveFreeSpace(hash *comp.SpatialHash, objID types.EntityID, startX, startY, targetX, targetY float32, radius int, team string) (float32, float32) {
	// Remove the object from its current position
	RemoveObjectFromSpatialHash(hash, objID, startX, startY, radius)
	// Find an alternative position if the target is occupied
	if CheckCollisionSpatialHash(hash, targetX, targetY, radius) {
		//walk around blocking unit
		targetX, targetY = moveToNearestFreeSpaceBox(hash, startX, startY, targetX, targetY, float32(radius))
	}
	// Add the object to the new position
	AddObjectSpatialHash(hash, objID, targetX, targetY, radius, team)
	return targetX, targetY
}

// creates a box between start and target point and checks all points from target to start finding first open position.
// box length is distance from start to target
// box width is radius*2
func moveToNearestFreeSpaceBox(hash *comp.SpatialHash, startX, startY, targetX, targetY, radius float32) (newX float32, newY float32) {
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
	step := length / 8 // or another division factor

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
				if !CheckCollisionSpatialHash(hash, testX, testY, int(radius)) {
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

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////UTILITY FUNCTIONS//////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// GetUnitComponents fetches all necessary components related to a unit entity.
func getUnitComponentsUM(world cardinal.WorldContext, unitID types.EntityID) (*comp.Position, *comp.UnitRadius, *comp.Attack, *comp.Team, *comp.Movespeed, *comp.MatchId, *comp.MapName, error) {
	position, err := cardinal.GetComponent[comp.Position](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Position component (getUnitComponentsUM): %v", err)
	}
	unitRadius, err := cardinal.GetComponent[comp.UnitRadius](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Unit Radius component (getUnitComponentsUM): %v", err)
	}
	unitAttack, err := cardinal.GetComponent[comp.Attack](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Unit attack component (getUnitComponentsUM): %v", err)
	}
	team, err := cardinal.GetComponent[comp.Team](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Team component (getUnitComponentsUM): %v", err)
	}
	movespeed, err := cardinal.GetComponent[comp.Movespeed](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Movespeed component (getUnitComponentsUM): %v", err)
	}
	matchId, err := cardinal.GetComponent[comp.MatchId](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving MatchId component (getUnitComponentsUM): %v", err)
	}
	mapName, err := cardinal.GetComponent[comp.MapName](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Distance component (getUnitComponentsUM): %v", err)
	}
	return position, unitRadius, unitAttack, team, movespeed, matchId, mapName, nil
}

// fetches target components
func getTargetComponentsUM(world cardinal.WorldContext, enemyID types.EntityID) (enemyPosition *comp.Position, enemyRadius *comp.UnitRadius, err error) {

	enemyPosition, err = cardinal.GetComponent[comp.Position](world, enemyID)
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving enemy Position component (getTargetComponentsUM): %v", err)
	}
	enemyRadius, err = cardinal.GetComponent[comp.UnitRadius](world, enemyID)
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving enemy Radius component (getTargetComponentsUM): %v", err)
	}
	return enemyPosition, enemyRadius, nil
}

//BILLARD BALL COLLISION
// // push a unit when they collide
// // simulates if position 1 hits position 2 like a billards ball and bounces the position 2 to our new target returns
// func pushUnitDirection(posX1, posY1, posX2, posY2, dirX, dirY, distance float32) (targetX, targetY float32) {
// 	deltaX := posX1 - posX2
// 	deltaY := posY1 - posY2
// 	length := float32(math.Sqrt(float64(deltaX*deltaX) + float64(deltaY*deltaY))) // Magnitude of the vector
// 	if length == 0 {                                                              // Avoid division by zero
// 		fmt.Println("Collision at the same position, no movement. (PushUnitDirection)")
// 		return posX2, posY2 // Return the current position of the second ball
// 	}
// 	// Calculate the normal vector
// 	normalX := deltaX / length
// 	normalY := deltaY / length

// 	// Calculate the dot product of the incoming vector and the normal
// 	dotProduct := dotProduct(dirX, dirY, normalX, normalY)

// 	// Apply the reflection formula
// 	newDirX := dirX - 2*dotProduct*normalX
// 	newDirY := dirY - 2*dotProduct*normalY

// 	// Normalize the resulting direction vector
// 	finalLength := float32(math.Sqrt(float64(newDirX*newDirX + newDirY*newDirY)))
// 	newDirX /= finalLength
// 	newDirY /= finalLength

// 	//move position 2 in the direction of newDir by the input distance
// 	targetX = posX2 + newDirX*distance
// 	targetY = posY2 + newDirY*distance
// 	return targetX, targetY
// }
