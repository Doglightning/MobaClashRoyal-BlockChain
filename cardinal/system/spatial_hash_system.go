package system

import (
	"fmt"
	"math"

	comp "MobaClashRoyal/component"

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

// intersect determines if two circles intersect.
func intersectSpatialHash(x1, y1 float32, r1 int, x2, y2 float32, r2 int) bool {
	distSq := (x1-x2)*(x1-x2) + (y1-y2)*(y1-y2)
	radiusSumSq := float32(r1+r2) * float32(r1+r2)
	return distSq <= radiusSumSq
}

// checkLineIntersectionSpatialHash calculates if the moving line intersects with the circle.
func checkLineIntersectionSpatialHash(startX, startY, endX, endY, circleCenterX, circleCenterY float32, radius int) bool {
	// Step 1: Calculate direction vector of the line
	dirX := endX - startX
	dirY := endY - startY

	// Step 2: Calculate the vector from start point to the circle's center
	toCircleX := circleCenterX - startX
	toCircleY := circleCenterY - startY

	// Prevent division by zero if the line is actually a point
	dirLengthSquared := dirX*dirX + dirY*dirY
	if dirLengthSquared == 0 {
		return math.Sqrt(float64(toCircleX*toCircleX+toCircleY*toCircleY)) <= float64(radius)
	}

	// Step 3: Project toCircle onto direction
	t := (toCircleX*dirX + toCircleY*dirY) / dirLengthSquared
	closestX := startX + t*dirX
	closestY := startY + t*dirY

	// Step 4: Calculate the distance from the closest point on the line to the circle's center
	distance := math.Sqrt(float64((closestX-circleCenterX)*(closestX-circleCenterX) + (closestY-circleCenterY)*(closestY-circleCenterY)))

	// Step 5: Check if the distance is less than the radius and the projection falls within the segment
	return distance <= float64(radius) && t >= 0 && t <= 1
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
