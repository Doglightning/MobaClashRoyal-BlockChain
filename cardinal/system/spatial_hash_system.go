package system

import (
	"fmt"
	"math"

	comp "MobaClashRoyal/component"

	"pkg.world.dev/world-engine/cardinal/types"
)

// AddObject adds an object with a radius to the spatial hash grid, considering all cells it may intersect.
func AddObjectSpatialHash(hash *comp.SpatialHash, objID types.EntityID, x, y float32, radius int) {
	startX := x - float32(radius)
	endX := x + float32(radius)
	startY := y - float32(radius)
	endY := y + float32(radius)

	// Calculate the range of cells that the object might occupy
	startCellX, startCellY := calculateSpatialHash(hash, startX, startY)
	endCellX, endCellY := calculateSpatialHash(hash, endX, endY)

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
				}
			}
			cell.UnitIDs = append(cell.UnitIDs, objID)
			cell.PositionsX = append(cell.PositionsX, x)
			cell.PositionsY = append(cell.PositionsY, y)
			cell.Radii = append(cell.Radii, radius)
			hash.Cells[hashKey] = cell
		}
	}
}

// RemoveObjectFromSpatialHash removes an object based on its position, radius, and ID from the spatial hash grid.
func RemoveObjectFromSpatialHash(hash *comp.SpatialHash, objID types.EntityID, x, y float32, radius int) {
	startX := x - float32(radius)
	endX := x + float32(radius)
	startY := y - float32(radius)
	endY := y + float32(radius)

	// Calculate the range of cells that the object might occupy
	startCellX, startCellY := calculateSpatialHash(hash, startX, startY)
	endCellX, endCellY := calculateSpatialHash(hash, endX, endY)

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
// It returns a list of collided object IDs.
func CheckCollisionSpatialHash(hash *comp.SpatialHash, x, y float32, radius int) bool {
	startX := x - float32(radius)
	endX := x + float32(radius)
	startY := y - float32(radius)
	endY := y + float32(radius)

	// Calculate the range of cells that the object might touch
	startCellX, startCellY := calculateSpatialHash(hash, startX, startY)
	endCellX, endCellY := calculateSpatialHash(hash, endX, endY)

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

func moveToNearestFreeSpaceSpatialHash(hash *comp.SpatialHash, startX, startY, targetX, targetY, radius float32) (float32, float32) {
	movementRadius := float32(math.Hypot(float64(targetX-startX), float64(targetY-startY)))
	minX := targetX - movementRadius
	maxX := targetX + movementRadius
	minY := targetY - movementRadius
	maxY := targetY + movementRadius

	// Step size could be a fraction of the radius for finer granularity
	step := radius / 4

	// Iterate over a rectangle surrounding the target position
	for x := minX; x <= maxX; x += step {
		for y := minY; y <= maxY; y += step {
			// Check if the position is free of collisions
			if !CheckCollisionSpatialHash(hash, x, y, int(radius)) {
				return x, y
			}
		}
	}

	return startX, startY // Stay at the current position if no free spot is found
}

// UpdateUnitPosition attempts to move the unit to a new position or finds an alternative nearby spot.
func UpdateUnitPositionSpatialHash(hash *comp.SpatialHash, objID types.EntityID, startX, startY, targetX, targetY float32, radius int) (newtargetX, newtargetY float32) {

	// Remove the object from its current position
	RemoveObjectFromSpatialHash(hash, objID, startX, startY, radius)

	if CheckCollisionSpatialHash(hash, targetX, targetY, radius) {
		// Find an alternative position if the target is occupied
		targetX, targetY = moveToNearestFreeSpaceSpatialHash(hash, startX, startY, targetX, targetY, float32(radius))
	}

	// Add the object to the new position
	AddObjectSpatialHash(hash, objID, targetX, targetY, radius)
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
