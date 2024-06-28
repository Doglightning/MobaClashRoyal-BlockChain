package system

import (
	"fmt"
	"math"
	"sort"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/iterators"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"

	comp "MobaClashRoyal/component"
)

// This function is called every tick automatically
// It updates the units position based on the direction vector stored in the games map for every unit entity
func UnitMovementSystem(world cardinal.WorldContext) error {
	// search through every map
	err := cardinal.NewSearch().Entity(
		filter.Exact(MapFilters()),
	).Each(world, func(id types.EntityID) bool {
		//store map component
		dirMap, err := cardinal.GetComponent[comp.DirectionMap](world, id)
		if err != nil {
			fmt.Printf("error retrieving map direction component: %v", err)
			return false
		}

		//store Grid Utilities component
		gridUtils, err := cardinal.GetComponent[comp.GridUtils](world, id)
		if err != nil {
			fmt.Printf("error retrieving map GridUtils component: %v", err)
			return false
		}

		//store map name component
		nameMap, err := cardinal.GetComponent[comp.MapName](world, id)
		if err != nil {
			fmt.Printf("error retrieving map name component (unit movement): %v", err)
			return false
		}

		//get all Unit Id's to current map name in priority of distance to base
		priorityUnitIDs, err := PriorityUnitMovement(world, nameMap)
		if err != nil {
			fmt.Printf("error in the priorityUnitIDs function (unit movement): %v", err)
			return false
		}

		//go through all ID's
		for _, id := range priorityUnitIDs {
			tempPosition, err := cardinal.GetComponent[comp.Position](world, id)
			if err != nil {
				fmt.Printf("error retrieving position component on tempPosition (unit movement): %v", err)
				continue
			}
			tempX := tempPosition.PositionVectorX
			tempY := tempPosition.PositionVectorY

			tempRadius, err := cardinal.GetComponent[comp.SizeCircle](world, id)
			if err != nil {
				fmt.Printf("error retrieving SizeCircle component on tempradius (unit movement): %v", err)
				continue
			}

			tempMatchID, err := cardinal.GetComponent[comp.MatchId](world, id)
			if err != nil {
				fmt.Printf("error retrieving matchid component on tempMatchID (unit movement): %v", err)
				continue
			}

			//get teamstate to get spatialhash tree
			teamFilter := cardinal.ComponentFilter[comp.MatchId](func(m comp.MatchId) bool {
				return m.MatchId == tempMatchID.MatchId
			})
			teamSearch := cardinal.NewSearch().Entity(
				filter.Exact(TeamFilters())).
				Where(teamFilter)

			foundTeam, err := teamSearch.First(world)

			if err != nil {

				fmt.Printf("error searching for match (unit movement): %s", err)
				continue
			}

			if foundTeam == iterators.BadID { // Assuming cardinal.NoEntity represents no result found
				fmt.Printf("no match found with ID or missing components (unit movement): %s", tempMatchID.MatchId)
				continue
			}

			tempSpartialHash, err := cardinal.GetComponent[comp.SpatialHash](world, foundTeam)
			if err != nil {
				fmt.Printf("error retrieving SpartialHash component on tempSpartialHash (unit movement): %s", err)
				continue
			}

			//RemoveObjectFromSpatialHash(tempSpartialHash, id, tempPosition.PositionVectorX, tempPosition.PositionVectorY, tempRadius.Radius)

			//normalize the units position to the maps grid increments.
			normalizedX := int(((int(tempPosition.PositionVectorX)-gridUtils.StartX)/gridUtils.Increment))*gridUtils.Increment + gridUtils.StartX
			normalizedY := int(((int(tempPosition.PositionVectorY)-gridUtils.StartY)/gridUtils.Increment))*gridUtils.Increment + gridUtils.StartY
			//The units (x,y) coordinates normalized and turned into proper key(string) format for seaching map
			coordKey := fmt.Sprintf("%d,%d", normalizedX, normalizedY)

			// Retrieve direction vector using coordinate key
			directionVector, exists := dirMap.Map[coordKey]
			if !exists {
				fmt.Println("No direction vector found for the given coordinates (unit movement)")
				continue // continue processing       <<---------------------------------- later need to create handling for units not found.  They outside map and bugged won't move.
			}
			//get units movespeed
			tempMovespeed, err := cardinal.GetComponent[comp.Movespeed](world, id)
			if err != nil {
				fmt.Printf("error retrieving Movespeed component on tempMovespeed (unit movement): %v", err)
				continue
			}
			//get team
			tempTeam, err := cardinal.GetComponent[comp.Team](world, id)
			if err != nil {
				fmt.Printf("error retrieving Team component on tempTeam (unit movement): %v", err)
				continue
			}

			//updated rotation based on team
			if tempTeam.Team == "Blue" {
				tempPosition.RotationVectorX = directionVector[0]
				tempPosition.RotationVectorY = directionVector[1]
			} else {
				tempPosition.RotationVectorX = directionVector[0] * -1
				tempPosition.RotationVectorY = directionVector[1] * -1
			}

			//update new x,y
			tempPosition.PositionVectorX = tempPosition.PositionVectorX + (tempPosition.RotationVectorX * tempMovespeed.CurrentMS)
			tempPosition.PositionVectorY = tempPosition.PositionVectorY + (tempPosition.RotationVectorY * tempMovespeed.CurrentMS)

			err = cardinal.SetComponent(world, id, tempPosition)
			if err != nil {
				fmt.Printf("error set component on tempPosition (unit movement): %v", err)
				continue
			}
			tempPosition.PositionVectorX, tempPosition.PositionVectorY = UpdateUnitPositionSpatialHash(tempSpartialHash, id, tempX, tempY, tempPosition.PositionVectorX, tempPosition.PositionVectorY, tempRadius.Radius)
			//AddObjectSpatialHash(tempSpartialHash, id, tempPosition.PositionVectorX, tempPosition.PositionVectorY, tempRadius.Radius)

			tempDistance, err := cardinal.GetComponent[comp.Distance](world, id)
			if err != nil {
				fmt.Printf("error retrieving distance component on tempDistance (unit movement): %v", err)
				continue
			}
			//calculate distance from enemy spawn
			if tempTeam.Team == "Blue" {
				tempDistance.Distance = math.Sqrt(((float64(tempPosition.PositionVectorX) - float64(gridUtils.RedX)) * (float64(tempPosition.PositionVectorX) - float64(gridUtils.RedX))) + ((float64(tempPosition.PositionVectorY) - float64(gridUtils.RedY)) * (float64(tempPosition.PositionVectorY) - float64(gridUtils.RedY))))
			} else {
				tempDistance.Distance = math.Sqrt(((float64(tempPosition.PositionVectorX) - float64(gridUtils.BlueX)) * (float64(tempPosition.PositionVectorX) - float64(gridUtils.BlueX))) + ((float64(tempPosition.PositionVectorY) - float64(gridUtils.BlueY)) * (float64(tempPosition.PositionVectorY) - float64(gridUtils.BlueY))))
			}
			//set distance
			err = cardinal.SetComponent(world, id, tempDistance)
			if err != nil {
				fmt.Printf("error setting Distance component on tempDistance (unit movement): %v", err)
				continue
			}
		}

		return true
	})
	if err != nil {
		return fmt.Errorf("error retrieving map entity (unit movement): %w", err)
	}

	return err

}

// NOTE THIS COULD GET REALLY NOT OPTIMAL AT LARGE SCALE WHEN POTENTIALLY HUNDREDS OF THOUSANDS OF UNITS COULD BE EVERY TICK GETTING REORDERED!!!
// Input:  MapNameComponent - Component containing the maps name.
// Return: A list of all the Units in order of closest distance Priority from closests to farthest
func PriorityUnitMovement(world cardinal.WorldContext, mapNameComponent *comp.MapName) ([]types.EntityID, error) {
	// UnitData struct to store both the EntityID and its Distance for sorting
	type UnitData struct {
		ID       types.EntityID
		Distance float64
	}
	// Filter for current map
	mapFilter := cardinal.ComponentFilter[comp.MapName](func(m comp.MapName) bool {
		return m.MapName == mapNameComponent.MapName
	})
	// Search all units with the filtered map into a list
	unitList, err := cardinal.NewSearch().Entity(
		filter.Exact(UnitFilters())).
		Where(mapFilter).Collect(world)
	if err != nil {
		return nil, fmt.Errorf("PriorityUnitMovement error searching for unit with map (unit movement): %w", err)
	}

	// Create a slice to store the units with their distances
	unitsData := make([]UnitData, 0, len(unitList))

	// Fetch distances for each unit and handle errors
	for _, unit := range unitList {
		distanceComp, err := cardinal.GetComponent[comp.Distance](world, unit)
		if err != nil {
			// Optionally handle the error differently or log it
			return nil, fmt.Errorf("error fetching distance for unit %v: %w (unit movement)", unit, err)
		}
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
