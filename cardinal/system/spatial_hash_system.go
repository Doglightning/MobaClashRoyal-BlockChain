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

// find closest free point between points B to A
func pushTowardsEnemySpatialHash(world cardinal.WorldContext, hash *comp.SpatialHash, id types.EntityID, startX, startY, targetX, targetY float32, radius int, distance float32, team *comp.Team) (float32, float32) {
	//get attack component
	atk, err := cardinal.GetComponent[comp.Attack](world, id)
	if err != nil {
		fmt.Printf("error getting attack compoenent (moveToNearestFreeSpaceLineSpatialHash): %v", err)
		return startX, startY
	}
	// //get mapname  component
	// mapName, err := cardinal.GetComponent[comp.MapName](world, id)
	// if err != nil {
	// 	fmt.Printf("error getting map name compoenent (moveToNearestFreeSpaceLineSpatialHash): %v", err)
	// 	return startX, startY
	// }
	//get length
	deltaX := targetX - startX
	deltaY := targetY - startY
	length := distanceBetweenTwoPoints(startX, startY, targetX, targetY)

	if length == 0 {
		fmt.Printf("length Dividing by 0 (pushTowardsEnemySpatialHash)\n")
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
			fmt.Printf("(moveToNearestFreeSpaceLineSpatialHash): %v", err)
			return startX, startY
		}

		// Search along the line from target to start (reverse)
		for d := float32(0); d <= length; d += float32(step) {
			testX := targetX + dirX*float32(d) // Start from target position
			testY := targetY + dirY*float32(d) // Start from target position
			// Check if the position is free of collisions
			if !CheckCollisionSpatialHash(hash, testX, testY, radius) {
				//distance - unit and target radius'
				adjustedDistance := distanceBetweenTwoPoints(targetPos.PositionVectorX, targetPos.PositionVectorY, testX, testY) - float32(radius) - float32(targetRadius.UnitRadius)
				//if within attack range
				if adjustedDistance <= float32(atk.AttackRadius) {
					return testX, testY // Return the first free spot found
				}
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

// find closest free point between points B to A
func closestFreeSpaceBetweenTwoPointsSpatialHash(hash *comp.SpatialHash, startX, startY, targetX, targetY float32, radius int) (newX float32, newY float32) {

	//get length
	deltaX := targetX - startX
	deltaY := targetY - startY
	length := distanceBetweenTwoPoints(startX, startY, targetX, targetY)

	// Normalize direction vector
	dirX := deltaX / length
	dirY := deltaY / length

	// Step size, which can be adjusted as needed
	step := length / 8

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

	return startX, startY // Stay at the current position if no free spot is found
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
