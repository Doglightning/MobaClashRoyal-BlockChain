package system

import (
	"container/list"
	"fmt"
	"math"

	comp "MobaClashRoyal/component"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

// AddObject adds an object with a radius to the spatial hash grid, considering all cells it may intersect.
func AddObjectSpatialHash(hash *comp.SpatialHash, objID types.EntityID, x, y float32, radius int, team string) {

	startCellX, endCellX, startCellY, endCellY := calculateCellRangeSpatialHash(hash, x, y, radius)

	// Loop over all cells the object might touch
	for cx := startCellX; cx <= endCellX; cx++ {
		for cy := startCellY; cy <= endCellY; cy++ {
			hashKey := fmt.Sprintf("%d,%d", cx, cy)
			cell, exists := hash.Cells[hashKey]
			if !exists {
				cell = comp.SpatialCell{
					UnitIDs:    []types.EntityID{},
					PositionsX: []float32{},
					PositionsY: []float32{},
					Radii:      []int{},
					Team:       []string{},
				}
			}
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
	startCellX, endCellX, startCellY, endCellY := calculateCellRangeSpatialHash(hash, x, y, radius)

	// Loop over all cells the object might touch
	for cx := startCellX; cx <= endCellX; cx++ {
		for cy := startCellY; cy <= endCellY; cy++ {
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
	startCellX, endCellX, startCellY, endCellY := calculateCellRangeSpatialHash(hash, x, y, radius)

	// Loop over all cells the object might touch
	for cx := startCellX; cx <= endCellX; cx++ {
		for cy := startCellY; cy <= endCellY; cy++ {
			hashKey := fmt.Sprintf("%d,%d", cx, cy)
			if cell, exists := hash.Cells[hashKey]; exists {
				for i := range cell.UnitIDs {
					ux := cell.PositionsX[i]
					uy := cell.PositionsY[i]
					uRadius := cell.Radii[i]

					if intersectSpatialHash(x, y, float32(radius), ux, uy, float32(uRadius)) {

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

	startCellX, endCellX, startCellY, endCellY := calculateCellRangeSpatialHash(hash, x, y, radius)
	collidedUnits := []types.EntityID{}

	// Loop over all cells the object might touch
	for cx := startCellX; cx <= endCellX; cx++ {
		for cy := startCellY; cy <= endCellY; cy++ {
			hashKey := fmt.Sprintf("%d,%d", cx, cy)
			if cell, exists := hash.Cells[hashKey]; exists {
				for i, unitID := range cell.UnitIDs {
					ux := cell.PositionsX[i]
					uy := cell.PositionsY[i]
					uRadius := cell.Radii[i]

					// Check if the unit intersects with the given circle
					if intersectSpatialHash(x, y, float32(radius), ux, uy, float32(uRadius)) {
						collidedUnits = append(collidedUnits, unitID)
					}
				}
			}
		}
	}

	return collidedUnits
}

func moveToNearestFreeSpaceBoxSpatialHash(hash *comp.SpatialHash, startX, startY, targetX, targetY, radius float32) (newX float32, newY float32) {
	deltaX := targetX - startX
	deltaY := targetY - startY
	length := float32(distanceBetweenTwoPointsVectorMath(float64(startX), float64(startY), float64(targetX), float64(targetY)))

	// Normalize direction vector
	dirX := deltaX / length
	dirY := deltaY / length

	// Perpendicular vector (normalized)
	perpX := -dirY
	perpY := dirX

	// Step size, which can be adjusted as needed
	step := length / 8 // or another division factor

	//search in a box the size of the units movement (kinda like a radius but less cpu intensive)

	halfWidth := radius / 2 // Half the unit's radius
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

func moveToNearestFreeSpaceLineSpatialHash(world cardinal.WorldContext, hash *comp.SpatialHash, id types.EntityID, startX, startY, targetX, targetY, radius float32) (newX float32, newY float32) {

	attack, err := cardinal.GetComponent[comp.Attack](world, id)
	if err != nil {
		fmt.Printf("error getting attack compoenent (moveToNearestFreeSpaceLineSpatialHash): %v", err)
		return startX, startY
	}
	deltaX := targetX - startX
	deltaY := targetY - startY
	length := float32(distanceBetweenTwoPointsVectorMath(float64(startX), float64(startY), float64(targetX), float64(targetY)))

	// Normalize direction vector
	dirX := deltaX / length
	dirY := deltaY / length

	// Step size, which can be adjusted as needed
	step := length / 8 // or another division factor

	if attack.Combat {
		attack, err := cardinal.GetComponent[comp.Attack](world, id)
		if err != nil {
			fmt.Printf("error getting attack redius compoenent (moveToNearestFreeSpaceLineSpatialHash): %v", err)
			return startX, startY
		}

		targetRadius, err := cardinal.GetComponent[comp.UnitRadius](world, attack.Target)
		if err != nil {
			fmt.Printf("error getting target attack redius compoenent (moveToNearestFreeSpaceLineSpatialHash): %v", err)
			return startX, startY
		}

		targetPos, err := cardinal.GetComponent[comp.Position](world, attack.Target)
		if err != nil {
			fmt.Printf("error getting target position compoenent (moveToNearestFreeSpaceLineSpatialHash): %v", err)
			return startX, startY
		}

		// Search along the line from target to start (reverse)
		for d := 0.0; d <= float64(length); d += float64(step) {
			testX := targetX + dirX*float32(d) // Start from target position
			testY := targetY + dirY*float32(d) // Start from target position
			// Check if the position is free of collisions
			if !CheckCollisionSpatialHash(hash, testX, testY, int(radius)) {
				adjustedDistance := distanceBetweenTwoPointsVectorMath(float64(targetPos.PositionVectorX), float64(targetPos.PositionVectorY), float64(testX), float64(testY)) - float64(radius) - float64(targetRadius.UnitRadius)
				if adjustedDistance <= float64(attack.AttackRadius) {
					fmt.Printf("adjustedDistance: %f\n", adjustedDistance)
					return testX, testY // Return the first free spot found
				}

			}
		}
	} else {
		// Search along the line from target to start (reverse)
		for d := 0.0; d <= float64(length); d += float64(step) {
			testX := targetX + dirX*float32(d) // Start from target position
			testY := targetY + dirY*float32(d) // Start from target position
			fmt.Printf("adjustedDistance: %t\n", attack.Combat)
			// Check if the position is free of collisions
			if !CheckCollisionSpatialHash(hash, testX, testY, int(radius)) {

				fmt.Printf("adjustedDistance: %t\n", attack.Combat)
				return testX, testY // Return the first free spot found

			}
		}

	}

	return startX, startY // Stay at the current position if no free spot is found
}

// UpdateUnitPosition attempts to move the unit to a new position or finds an alternative nearby spot.
func UpdateUnitPositionPushSpatialHash(world cardinal.WorldContext, hash *comp.SpatialHash, objID types.EntityID, startX, startY, targetX, targetY float32, radius int, team string, distance float32) (newtargetX, newtargetY float32) {

	collisionList := CheckCollisionSpatialHashList(hash, targetX, targetY, radius)

	//try to push blocking units
	if len(collisionList) > 0 {
		for _, collisionID := range collisionList {
			if collisionID == objID {
				continue
			}

			enemyTeam, err := cardinal.GetComponent[comp.Team](world, collisionID)
			if err != nil {
				fmt.Printf("error getting enemy radius compoenent (UpdateUnitPositionPushSpatialHash): %v", err)
				continue
			}
			if enemyTeam.Team == team {

				enemyPos, err := cardinal.GetComponent[comp.Position](world, collisionID)
				if err != nil {
					fmt.Printf("error getting enemy position compoenent (UpdateUnitPositionPushSpatialHash): %v", err)
					continue
				}

				enemyRadius, err := cardinal.GetComponent[comp.UnitRadius](world, collisionID)
				if err != nil {
					fmt.Printf("error getting enemy radius compoenent (UpdateUnitPositionPushSpatialHash): %v", err)
					continue
				}

				// Remove the object from its current position
				RemoveObjectFromSpatialHash(hash, collisionID, enemyPos.PositionVectorX, enemyPos.PositionVectorY, enemyRadius.UnitRadius)
				targetEnemyX, targetEnemyY := PushUnitDirSpatialHash(collisionID, targetX, targetY, enemyPos.PositionVectorX, enemyPos.PositionVectorY, startX-targetX, startY-targetY, float32(distance))

				// Find an alternative position if the target is occupied

				enemyPos.PositionVectorX, enemyPos.PositionVectorY = moveToNearestFreeSpaceLineSpatialHash(world, hash, collisionID, enemyPos.PositionVectorX, enemyPos.PositionVectorY, targetEnemyX, targetEnemyY, float32(enemyRadius.UnitRadius))

				// Add the object to the new position
				AddObjectSpatialHash(hash, collisionID, enemyPos.PositionVectorX, enemyPos.PositionVectorY, enemyRadius.UnitRadius, enemyTeam.Team)
				err = cardinal.SetComponent(world, collisionID, enemyPos)
				if err != nil {
					fmt.Printf("error setting enemy pos component (UpdateUnitPositionPushSpatialHash): %v", err)
					continue
				}
			}
		}
	}

	// Remove the object from its current position
	RemoveObjectFromSpatialHash(hash, objID, startX, startY, radius)
	// Find an alternative position if the target is occupied
	if CheckCollisionSpatialHash(hash, targetX, targetY, radius) {
		targetX, targetY = moveToNearestFreeSpaceBoxSpatialHash(hash, startX, startY, targetX, targetY, float32(radius))
	}

	// Add the object to the new position
	AddObjectSpatialHash(hash, objID, targetX, targetY, radius, team)
	return targetX, targetY
}

func PushUnitDirSpatialHash(enemyID types.EntityID, posX1, posY1, posX2, posY2, dirX, dirY, distance float32) (targetX, targetY float32) {
	// Calculate the normal vector
	deltaX := posX1 - posX2
	deltaY := posY1 - posY2
	length := math.Sqrt(float64(deltaX*deltaX + deltaY*deltaY)) // Magnitude of the vector
	if length == 0 {                                            // Avoid division by zero
		fmt.Println("Collision at the same position, no movement.")
		return posX2, posY2 // Return the current position of the second ball
	}

	normalX := deltaX / float32(length)
	normalY := deltaY / float32(length)

	// Calculate the dot product of the incoming vector and the normal
	dotProduct := dirX*normalX + dirY*normalY

	// Apply the reflection formula
	newDirX := dirX - 2*dotProduct*normalX
	newDirY := dirY - 2*dotProduct*normalY

	// Normalize the resulting direction vector
	finalLength := math.Sqrt(float64(newDirX*newDirX + newDirY*newDirY))
	newDirX /= float32(finalLength)
	newDirY /= float32(finalLength)

	targetX = posX2 + newDirX*distance
	targetY = posY2 + newDirY*distance
	return targetX, targetY
}

// intersect determines if two circles intersect.
func intersectSpatialHash(x1, y1, r1, x2, y2, r2 float32) bool {
	distSq := (x1-x2)*(x1-x2) + (y1-y2)*(y1-y2)
	radiusSumSq := (r1 + r2) * (r1 + r2)
	return distSq <= radiusSumSq
}

// calculateHash computes the cell coordinates based on the position and the start of the grid.
func calculateSpatialHash(hash *comp.SpatialHash, x, y float32) (int, int) {
	cx := int(math.Floor(float64((x - hash.StartX) / float32(hash.CellSize))))
	cy := int(math.Floor(float64((y - hash.StartY) / float32(hash.CellSize))))
	return cx, cy
}

func calculateCellRangeSpatialHash(hash *comp.SpatialHash, x, y float32, radius int) (startCellX, endCellX, startCellY, endCellY int) {
	startX := x - float32(radius)
	endX := x + float32(radius)
	startY := y - float32(radius)
	endY := y + float32(radius)

	startCellX, startCellY = calculateSpatialHash(hash, startX, startY)
	endCellX, endCellY = calculateSpatialHash(hash, endX, endY)

	return startCellX, endCellX, startCellY, endCellY
}

// FindClosestEnemy performs a BFS search from the unit's position outward within the attack radius.
func FindClosestEnemySpatialHash(hash *comp.SpatialHash, objID types.EntityID, startX, startY float32, attackRadius int, team string) (types.EntityID, float32, float32, int, bool) {
	queue := list.New()
	visited := make(map[string]bool)
	queue.PushBack(&comp.Position{PositionVectorX: startX, PositionVectorY: startY})
	minDist := float64(attackRadius * attackRadius) // Using squared distance to avoid sqrt calculations.
	closestEnemy := types.EntityID(0)
	closestX, closestY := float32(0), float32(0)
	closestRadius := int(0)
	foundEnemy := false

	//fmt.Printf("Starting search with attackRadius: %d, cellSize: %d\n", attackRadius, hash.CellSize)

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

		//fmt.Printf("Visiting cell: %s\n", hashKey)

		if cell, exists := hash.Cells[hashKey]; exists {
			for i, id := range cell.UnitIDs {
				if cell.Team[i] != team && id != objID {
					distSq := float64((cell.PositionsX[i]-startX)*(cell.PositionsX[i]-startX)+(cell.PositionsY[i]-startY)*(cell.PositionsY[i]-startY)) - float64(cell.Radii[i]*cell.Radii[i])
					//fmt.Printf("Checking unit %d at (%f, %f) with distSq: %f\n", id, cell.PositionsX[i], cell.PositionsY[i], distSq)
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
					if (nx-startX)*(nx-startX)+(ny-startY)*(ny-startY) <= float32(attackRadius*attackRadius) {
						queue.PushBack(&comp.Position{PositionVectorX: nx, PositionVectorY: ny})
						//fmt.Printf("Adding cell to queue: (%f, %f)\n", nx, ny)
					}
				}
			}
		}
	}

	//fmt.Printf("Search completed. Found enemy: %v\n", foundEnemy)

	return closestEnemy, closestX, closestY, closestRadius, foundEnemy
}
