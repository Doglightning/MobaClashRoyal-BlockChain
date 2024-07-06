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
	).Each(world, func(mapID types.EntityID) bool {

		//Get All Components Needed From Map ID
		dirMap, gridUtils, nameMap, err := GetMapComponentsUM(world, mapID)
		if err != nil {
			fmt.Printf("%v", err)
			return false
		}

		//get all Unit Id's to current map name in priority of distance to base
		priorityUnitIDs, err := PriorityUnitMovement(world, nameMap)
		if err != nil {
			fmt.Printf("error in the priorityUnitIDs function (unit movement): %v", err)
			return false
		}

		//go through all Unit ID's
		for _, id := range priorityUnitIDs {
			//get Unit Components
			UnitPosition, UnitRadius, UnitAttackRadius, UnitAttack, UnitTeam, UnitMovespeed, MatchID, UnitDistance, err := GetUnitComponentsUM(world, id)
			if err != nil {
				fmt.Printf("%v", err)
				continue
			}

			//get team state
			foundTeam, err := getTeamFilterUM(world, MatchID)
			if err != nil {
				fmt.Printf("%v", err)
				continue
			}

			//get Spatial Hash
			CollisionSpartialHash, err := cardinal.GetComponent[comp.SpatialHash](world, foundTeam)
			if err != nil {
				fmt.Printf("error retrieving SpartialHash component on tempSpartialHash (unit movement): %s", err)
				continue
			}

			//if units in combat
			if UnitAttack.Combat {
				//get enemyID
				enemyID := UnitAttack.Target
				enemyPosition, enemyRadius, errr := getEnemyComponentsUM(world, enemyID)
				if errr != nil {
					fmt.Printf("(unit movement): %s", err)
					continue
				}
				err = MoveUnitTowardsEnemyUM(world, id, UnitPosition, float64(enemyPosition.PositionVectorX), float64(enemyPosition.PositionVectorY), enemyRadius.UnitRadius, UnitTeam, UnitMovespeed, CollisionSpartialHash, UnitRadius, UnitDistance, gridUtils)

			}

			//if units not in combat
			if !UnitAttack.Combat {
				//Check for in range Enemies
				enemyID, enemyX, enemyY, enemyRadius, found := FindClosestEnemySpatialHash(CollisionSpartialHash, id, UnitPosition.PositionVectorX, UnitPosition.PositionVectorY, UnitAttackRadius.AttackRadius, UnitTeam.Team)
				if found {
					// Calculate squared distance between the unit and the enemy, minus their radii
					deltaX := enemyX - UnitPosition.PositionVectorX
					deltaY := enemyY - UnitPosition.PositionVectorY
					squaredDistance := (deltaX * deltaX) + (deltaY * deltaY)
					adjustedDistance := math.Sqrt(float64(squaredDistance)) - float64(enemyRadius) - float64(UnitRadius.UnitRadius)
					if adjustedDistance <= 3 {
						UnitAttack.Combat = true
						UnitAttack.Target = enemyID
						err = cardinal.SetComponent(world, id, UnitAttack)
						if err != nil {
							fmt.Printf("error setting Distance component on tempDistance (unit movement): %v", err)
							continue
						}

					}

					err = MoveUnitTowardsEnemyUM(world, id, UnitPosition, float64(enemyX), float64(enemyY), enemyRadius, UnitTeam, UnitMovespeed, CollisionSpartialHash, UnitRadius, UnitDistance, gridUtils)
				} else {
					//no enemies found and not in combat, move with direction map.
					err = MoveUnitDirectionMapUM(world, id, UnitPosition, gridUtils, dirMap, UnitTeam, UnitMovespeed, CollisionSpartialHash, UnitRadius, UnitDistance)
				}

			}
			if err != nil {
				fmt.Printf("(unit movement): %s", err)
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

// Moves Unit in direction of the map Direction vector
func MoveUnitDirectionMapUM(world cardinal.WorldContext, id types.EntityID, position *comp.Position, gridUtils *comp.GridUtils, dirMap *comp.DirectionMap, team *comp.Team, movespeed *comp.Movespeed, collisionSpartialHash *comp.SpatialHash, radius *comp.UnitRadius, distance *comp.Distance) error {

	//normalize the units position to the maps grid increments.
	normalizedX := int(((int(position.PositionVectorX)-gridUtils.StartX)/gridUtils.Increment))*gridUtils.Increment + gridUtils.StartX
	normalizedY := int(((int(position.PositionVectorY)-gridUtils.StartY)/gridUtils.Increment))*gridUtils.Increment + gridUtils.StartY
	//The units (x,y) coordinates normalized and turned into proper key(string) format for seaching map
	coordKey := fmt.Sprintf("%d,%d", normalizedX, normalizedY)

	maps, _ := getMapData("ProtoType")
	// Retrieve direction vector using coordinate key
	//directionVector, exists := dirMap.Map[coordKey]
	directionVector, exists := maps.DirMap.DMap[coordKey]
	if !exists {
		return fmt.Errorf("no direction vector found for the given coordinates (unit movement)")
	}
	//updated rotation based on team
	if team.Team == "Blue" {
		position.RotationVectorX = directionVector[0]
		position.RotationVectorY = directionVector[1]
	} else {
		position.RotationVectorX = directionVector[0] * -1
		position.RotationVectorY = directionVector[1] * -1
	}

	// //Store Original X and Y
	tempX := position.PositionVectorX
	tempY := position.PositionVectorY
	//update new x,y
	position.PositionVectorX = position.PositionVectorX + (position.RotationVectorX * movespeed.CurrentMS)
	position.PositionVectorY = position.PositionVectorY + (position.RotationVectorY * movespeed.CurrentMS)

	position.PositionVectorX, position.PositionVectorY = UpdateUnitPositionSpatialHash(collisionSpartialHash, id, tempX, tempY, position.PositionVectorX, position.PositionVectorY, radius.UnitRadius, team.Team)

	err := cardinal.SetComponent(world, id, position)
	if err != nil {
		return fmt.Errorf("error set component on tempPosition (unit movement/MoveUnitDirectionMapUM): %v", err)
	}

	//update units new distance from enemy base
	if err = UpdateUnitDistance(world, id, team, distance, position, gridUtils); err != nil {
		return err
	}

	return nil
}

// Moves Unit towards enemy position
func MoveUnitTowardsEnemyUM(world cardinal.WorldContext, id types.EntityID, position *comp.Position, enemyX float64, enemyY float64, enemyRadius int, team *comp.Team, movespeed *comp.Movespeed, collisionSpartialHash *comp.SpatialHash, radius *comp.UnitRadius, distance *comp.Distance, gridUtils *comp.GridUtils) error {
	// Compute direction vector towards the enemy
	deltaX := enemyX - float64(position.PositionVectorX)
	deltaY := enemyY - float64(position.PositionVectorY)
	magnitude := math.Sqrt(deltaX*deltaX + deltaY*deltaY)

	// Calculate the stopping distance (combined radii of the unit and enemy plus 1 pixel for separation)
	stoppingDistance := float64(radius.UnitRadius + enemyRadius + 1)

	// Normalize the direction vector
	directionVectorX := float32(deltaX / magnitude)
	directionVectorY := float32(deltaY / magnitude)

	// Calculate the target position to move towards, stopping 1 pixel outside the enemy's radius
	targetX := float32(enemyX) - directionVectorX*float32(stoppingDistance)
	targetY := float32(enemyY) - directionVectorY*float32(stoppingDistance)

	// Compute new position based on movespeed and direction, but do not exceed the target position
	newPosX := position.PositionVectorX + directionVectorX*movespeed.CurrentMS
	newPosY := position.PositionVectorY + directionVectorY*movespeed.CurrentMS

	// Ensure the unit does not overshoot the target position
	if (directionVectorX > 0 && newPosX > targetX) || (directionVectorX < 0 && newPosX < targetX) {
		newPosX = targetX
	}
	if (directionVectorY > 0 && newPosY > targetY) || (directionVectorY < 0 && newPosY < targetY) {
		newPosY = targetY
	}

	// Update position in spatial hash
	position.PositionVectorX, position.PositionVectorY = UpdateUnitPositionSpatialHash(collisionSpartialHash, id, position.PositionVectorX, position.PositionVectorY, newPosX, newPosY, radius.UnitRadius, team.Team)

	position.RotationVectorX = directionVectorX
	position.RotationVectorY = directionVectorY

	// Set the new position component
	err := cardinal.SetComponent(world, id, position)
	if err != nil {
		return fmt.Errorf("error set component on tempPosition (unit movement/MoveUnitTowardsEnemyUM): %v", err)
	}

	// Update units new distance from enemy base
	if err = UpdateUnitDistance(world, id, team, distance, position, gridUtils); err != nil {
		return err
	}

	return nil
}

// Update units distance from enemy base to help with movement priority queue
func UpdateUnitDistance(world cardinal.WorldContext, id types.EntityID, team *comp.Team, distance *comp.Distance, position *comp.Position, gridUtils *comp.GridUtils) error {
	// calculate distance from enemy spawn
	if team.Team == "Blue" {
		distance.Distance = math.Sqrt(((float64(position.PositionVectorX) - float64(gridUtils.RedX)) * (float64(position.PositionVectorX) - float64(gridUtils.RedX))) + ((float64(position.PositionVectorY) - float64(gridUtils.RedY)) * (float64(position.PositionVectorY) - float64(gridUtils.RedY))))
	} else {
		distance.Distance = math.Sqrt(((float64(position.PositionVectorX) - float64(gridUtils.BlueX)) * (float64(position.PositionVectorX) - float64(gridUtils.BlueX))) + ((float64(position.PositionVectorY) - float64(gridUtils.BlueY)) * (float64(position.PositionVectorY) - float64(gridUtils.BlueY))))
	}
	// set distance
	err := cardinal.SetComponent(world, id, distance)
	if err != nil {
		return fmt.Errorf("error setting Distance component on UpdateUnitDistance (unit movement): %v", err)
	}
	return nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////UTILITY FUNCTIONS//////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// GetMapComponents fetches all necessary components related to a map entity.
func GetMapComponentsUM(world cardinal.WorldContext, mapID types.EntityID) (*comp.DirectionMap, *comp.GridUtils, *comp.MapName, error) {
	dirMap, err := cardinal.GetComponent[comp.DirectionMap](world, mapID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error retrieving Direction Map component (Unit Movement): %v", err)
	}
	gridUtils, err := cardinal.GetComponent[comp.GridUtils](world, mapID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error retrieving Grid Utilities component (Unit Movement): %v", err)
	}
	mapName, err := cardinal.GetComponent[comp.MapName](world, mapID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error retrieving Map Name component (Unit Movement): %v", err)
	}
	return dirMap, gridUtils, mapName, nil
}

// GetUnitComponents fetches all necessary components related to a unit entity.
func GetUnitComponentsUM(world cardinal.WorldContext, unitID types.EntityID) (*comp.Position, *comp.UnitRadius, *comp.AttackRadius, *comp.Attack, *comp.Team, *comp.Movespeed, *comp.MatchId, *comp.Distance, error) {
	position, err := cardinal.GetComponent[comp.Position](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Position component (Unit Movement): %v", err)
	}
	unitRadius, err := cardinal.GetComponent[comp.UnitRadius](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Unit Radius component (Unit Movement): %v", err)
	}
	unitAttackRadius, err := cardinal.GetComponent[comp.AttackRadius](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Unit attack Radius component (Unit Movement): %v", err)
	}
	unitAttack, err := cardinal.GetComponent[comp.Attack](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Unit attack component (Unit Movement): %v", err)
	}
	team, err := cardinal.GetComponent[comp.Team](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Team component (Unit Movement): %v", err)
	}
	movespeed, err := cardinal.GetComponent[comp.Movespeed](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Movespeed component (Unit Movement): %v", err)
	}
	matchId, err := cardinal.GetComponent[comp.MatchId](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving MatchId component (Unit Movement): %v", err)
	}
	distance, err := cardinal.GetComponent[comp.Distance](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Distance component (Unit Movement): %v", err)
	}
	return position, unitRadius, unitAttackRadius, unitAttack, team, movespeed, matchId, distance, nil
}

// fetches enemy components
func getEnemyComponentsUM(world cardinal.WorldContext, enemyID types.EntityID) (enemyPosition *comp.Position, enemyRadius *comp.UnitRadius, err error) {

	enemyPosition, err = cardinal.GetComponent[comp.Position](world, enemyID)
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving enemy Position component (Unit Movement): %v", err)
	}
	enemyRadius, err = cardinal.GetComponent[comp.UnitRadius](world, enemyID)
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving enemy Radius component (Unit Movement): %v", err)
	}
	return enemyPosition, enemyRadius, nil
}

// Returns the Team state
func getTeamFilterUM(world cardinal.WorldContext, mID *comp.MatchId) (types.EntityID, error) {
	//get teamstate to get spatialhash tree
	teamFilter := cardinal.ComponentFilter[comp.MatchId](func(m comp.MatchId) bool {
		return m.MatchId == mID.MatchId
	})
	foundTeam, err := cardinal.NewSearch().Entity(
		filter.Exact(TeamFilters())).
		Where(teamFilter).First(world)

	if err != nil {

		fmt.Printf("error searching for match (unit movement): %s", err)
		return foundTeam, err
	}

	if foundTeam == iterators.BadID { // Assuming cardinal.NoEntity represents no result found
		return foundTeam, fmt.Errorf("no match found with ID or missing components (unit movement): %s", mID.MatchId)
	}
	return foundTeam, nil
}
