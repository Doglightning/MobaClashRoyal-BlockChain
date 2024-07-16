package system

import (
	"container/list"
	"fmt"
	"math"

	comp "MobaClashRoyal/component"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

// adds an object with a radius to the spatial hash grid, considering all cells it may intersect.
func AddObjectSpatialHash(hash *comp.SpatialHash, objID types.EntityID, x, y float32, radius int, team string) {
	//get range of cells covered
	startCellX, endCellX, startCellY, endCellY := calculateCellRangeSpatialHash(hash, x, y, radius)

	// Loop over all cells the object might touch
	for cx := startCellX; cx <= endCellX; cx++ {
		for cy := startCellY; cy <= endCellY; cy++ {
			//get x,y position key
			hashKey := fmt.Sprintf("%d,%d", cx, cy)
			cell, exists := hash.Cells[hashKey]
			//if key doesn't exist
			if !exists {
				//create cell data lists
				cell = comp.SpatialCell{
					UnitIDs:    []types.EntityID{},
					PositionsX: []float32{},
					PositionsY: []float32{},
					Radii:      []int{},
					Team:       []string{},
				}
			}
			//add to the cell data list
			cell.UnitIDs = append(cell.UnitIDs, objID)
			cell.PositionsX = append(cell.PositionsX, x)
			cell.PositionsY = append(cell.PositionsY, y)
			cell.Radii = append(cell.Radii, radius)
			cell.Team = append(cell.Team, team)
			hash.Cells[hashKey] = cell
		}
	}
}

// RemoveObjectFromSpatialHash removes an object based on its position, radius, and ID from the spatial hash grid.
func RemoveObjectFromSpatialHash(hash *comp.SpatialHash, objID types.EntityID, x, y float32, radius int) {
	//get range of cells covered
	startCellX, endCellX, startCellY, endCellY := calculateCellRangeSpatialHash(hash, x, y, radius)

	// Loop over all cells the object might touch
	for cx := startCellX; cx <= endCellX; cx++ {
		for cy := startCellY; cy <= endCellY; cy++ {
			//get x,y position key
			hashKey := fmt.Sprintf("%d,%d", cx, cy)
			if cell, exists := hash.Cells[hashKey]; exists {
				// Find and remove the object ID in the cell
				for i := len(cell.UnitIDs) - 1; i >= 0; i-- {
					if cell.UnitIDs[i] == objID {
						// Remove the object from the cell lists
						cell.UnitIDs = append(cell.UnitIDs[:i], cell.UnitIDs[i+1:]...)
						cell.PositionsX = append(cell.PositionsX[:i], cell.PositionsX[i+1:]...)
						cell.PositionsY = append(cell.PositionsY[:i], cell.PositionsY[i+1:]...)
						cell.Radii = append(cell.Radii[:i], cell.Radii[i+1:]...)
						cell.Team = append(cell.Team[:i], cell.Team[i+1:]...)
					}
				}
				// Update the cell in the map or delete it if empty
				if len(cell.UnitIDs) == 0 {
					delete(hash.Cells, hashKey)
				} else {
					hash.Cells[hashKey] = cell
				}
			}
		}
	}
}

// CheckCollision checks for collisions given an object's position and radius.
// It returns a true if collosion
func CheckCollisionSpatialHash(hash *comp.SpatialHash, x, y float32, radius int) bool {
	//get range of cells covered
	startCellX, endCellX, startCellY, endCellY := calculateCellRangeSpatialHash(hash, x, y, radius)

	// Loop over all cells the object might touch
	for cx := startCellX; cx <= endCellX; cx++ {
		for cy := startCellY; cy <= endCellY; cy++ {
			//get x,y position key
			hashKey := fmt.Sprintf("%d,%d", cx, cy)
			if cell, exists := hash.Cells[hashKey]; exists {
				//check all units in that cell
				for i := range cell.UnitIDs {
					//get unit data
					ux := cell.PositionsX[i]
					uy := cell.PositionsY[i]
					uRadius := cell.Radii[i]

					//check if intersection occurs
					if intersectSpatialHash(x, y, radius, ux, uy, uRadius) {
						return true
					}
				}
			}
		}
	}
	return false
}

// CheckCollisionSpatialHash checks for collisions given an object's position and radius.
// It returns a list of collided unit IDs.
func CheckCollisionSpatialHashList(hash *comp.SpatialHash, x, y float32, radius int) []types.EntityID {
	//get range of cells covered
	startCellX, endCellX, startCellY, endCellY := calculateCellRangeSpatialHash(hash, x, y, radius)
	collidedUnits := []types.EntityID{}

	// Loop over all cells the object might touch
	for cx := startCellX; cx <= endCellX; cx++ {
		for cy := startCellY; cy <= endCellY; cy++ {
			//get x,y position key
			hashKey := fmt.Sprintf("%d,%d", cx, cy)
			if cell, exists := hash.Cells[hashKey]; exists {
				//check all units in that cell
				for i, unitID := range cell.UnitIDs {
					//get unit data
					ux := cell.PositionsX[i]
					uy := cell.PositionsY[i]
					uRadius := cell.Radii[i]

					//check if intersection occurs
					if intersectSpatialHash(x, y, radius, ux, uy, uRadius) {
						collidedUnits = append(collidedUnits, unitID)
					}
				}
			}
		}
	}

	return collidedUnits
}

// creates a box between start and target point and checks all points from target to start finding first open position.
// box length is distance from start to target
// box width is radius*2
func moveToNearestFreeSpaceBoxSpatialHash(hash *comp.SpatialHash, startX, startY, targetX, targetY, radius float32) (newX float32, newY float32) {
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

// find closest free point between points B to A
func moveToNearestFreeSpaceLineSpatialHash(world cardinal.WorldContext, hash *comp.SpatialHash, id types.EntityID, startX, startY, targetX, targetY, radius float32) (newX float32, newY float32) {
	//get attack component
	atk, err := cardinal.GetComponent[comp.Attack](world, id)
	if err != nil {
		fmt.Printf("error getting attack compoenent (moveToNearestFreeSpaceLineSpatialHash): %v", err)
		return startX, startY
	}
	//get length
	deltaX := targetX - startX
	deltaY := targetY - startY
	length := distanceBetweenTwoPoints(startX, startY, targetX, targetY)

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
			fmt.Printf("(moveToNearestFreeSpaceLineSpatialHash): %v", err)
			return startX, startY
		}

		// Search along the line from target to start (reverse)
		for d := float32(0); d <= length; d += float32(step) {
			testX := targetX + dirX*float32(d) // Start from target position
			testY := targetY + dirY*float32(d) // Start from target position
			// Check if the position is free of collisions
			if !CheckCollisionSpatialHash(hash, testX, testY, int(radius)) {
				//distance - unit and target radius'
				adjustedDistance := distanceBetweenTwoPoints(targetPos.PositionVectorX, targetPos.PositionVectorY, testX, testY) - radius - float32(targetRadius.UnitRadius)
				//if within attack range
				if adjustedDistance <= float32(atk.AttackRadius) {
					return testX, testY // Return the first free spot found
				}
			}
		}
	} else {
		//not in cobat
		// Search along the line from target to start (reverse)
		for d := float32(0); d <= length; d += float32(step) {
			testX := targetX + dirX*float32(d) // Start from target position
			testY := targetY + dirY*float32(d) // Start from target position
			// Check if the position is free of collisions
			if !CheckCollisionSpatialHash(hash, testX, testY, int(radius)) {
				return testX, testY // Return the first free spot found
			}
		}
	}
	return startX, startY // Stay at the current position if no free spot is found
}

// attempts to move the unit to a new position.
// Can push units or walk around.
func UpdateUnitPositionPushSpatialHash(world cardinal.WorldContext, hash *comp.SpatialHash, objID types.EntityID, startX, startY, targetX, targetY float32, radius int, team string, distance float32) (newtargetX, newtargetY float32) {
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
				newTargetX, newTargetY := PushUnitDirSpatialHash(collisionID, targetX, targetY, targetPos.PositionVectorX, targetPos.PositionVectorY, startX-targetX, startY-targetY, float32(distance))

				// Find an alternative position if the target is occupied
				targetPos.PositionVectorX, targetPos.PositionVectorY = moveToNearestFreeSpaceLineSpatialHash(world, hash, collisionID, targetPos.PositionVectorX, targetPos.PositionVectorY, newTargetX, newTargetY, float32(targetRadius.UnitRadius))

				// Add the objects position to collosion hash
				AddObjectSpatialHash(hash, collisionID, targetPos.PositionVectorX, targetPos.PositionVectorY, targetRadius.UnitRadius, targetTeam.Team)
				//set collided units new position component
				err = cardinal.SetComponent(world, collisionID, targetPos)
				if err != nil {
					fmt.Printf("error setting target pos component (UpdateUnitPositionPushSpatialHash): %v", err)
					continue
				}
			}
		}
	}

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

// push a unit when they collide
// simulates if position 1 hits position 2 like a billards ball and bounces the position 2 to our new target returns
func PushUnitDirSpatialHash(enemyID types.EntityID, posX1, posY1, posX2, posY2, dirX, dirY, distance float32) (targetX, targetY float32) {
	deltaX := posX1 - posX2
	deltaY := posY1 - posY2
	length := float32(math.Sqrt(float64(deltaX*deltaX) + float64(deltaY*deltaY))) // Magnitude of the vector
	if length == 0 {                                                              // Avoid division by zero
		fmt.Println("Collision at the same position, no movement. (PushUnitDirSpatialHash)")
		return posX2, posY2 // Return the current position of the second ball
	}
	// Calculate the normal vector
	normalX := deltaX / length
	normalY := deltaY / length

	// Calculate the dot product of the incoming vector and the normal
	dotProduct := dotProduct(dirX, dirY, normalX, normalY)

	// Apply the reflection formula
	newDirX := dirX - 2*dotProduct*normalX
	newDirY := dirY - 2*dotProduct*normalY

	// Normalize the resulting direction vector
	finalLength := float32(math.Sqrt(float64(newDirX*newDirX + newDirY*newDirY)))
	newDirX /= finalLength
	newDirY /= finalLength

	//move position 2 in the direction of newDir by the input distance
	targetX = posX2 + newDirX*distance
	targetY = posY2 + newDirY*distance
	return targetX, targetY
}

// intersect determines if two circles intersect.
func intersectSpatialHash(x1, y1 float32, r1 int, x2, y2 float32, r2 int) bool {
	distSq := (x1-x2)*(x1-x2) + (y1-y2)*(y1-y2)
	radiusSumSq := float32(r1+r2) * float32(r1+r2)
	return distSq <= radiusSumSq
}

// normalizes a coord to the maps partition size
// used for hash key
func calculateSpatialHash(hash *comp.SpatialHash, x, y float32) (int, int) {
	cx := int(math.Floor(float64((x - hash.StartX) / float32(hash.CellSize))))
	cy := int(math.Floor(float64((y - hash.StartY) / float32(hash.CellSize))))
	return cx, cy
}

// calculates the range of coverered cells from a position of a radius in size
func calculateCellRangeSpatialHash(hash *comp.SpatialHash, x, y float32, radius int) (startCellX, endCellX, startCellY, endCellY int) {
	startX := x - float32(radius)
	endX := x + float32(radius)
	startY := y - float32(radius)
	endY := y + float32(radius)

	//get normalized cells
	startCellX, startCellY = calculateSpatialHash(hash, startX, startY)
	endCellX, endCellY = calculateSpatialHash(hash, endX, endY)

	return startCellX, endCellX, startCellY, endCellY
}

// FindClosestEnemy performs a BFS search from the unit's position outward within the attack radius.
func FindClosestEnemySpatialHash(hash *comp.SpatialHash, objID types.EntityID, startX, startY float32, attackRadius int, team string) (types.EntityID, float32, float32, int, bool) {
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
