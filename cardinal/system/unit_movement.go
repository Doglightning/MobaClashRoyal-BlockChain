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
		uPos, uRadius, uAtk, uTeam, uMs, MatchID, uDist, mapName, err := GetUnitComponentsUM(world, id)
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
			ePos, eRadius, errr := getTargetComponentsUM(world, enemyID) //get enemy position and radius components
			if errr != nil {
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
				err = cardinal.SetComponent(world, id, uAtk)
				if err != nil {
					fmt.Printf("error setting Distance component on tempDistance (unit_movement.go): %v", err)
					continue
				}
				//move towards enemy in combat with
				if uMs.CurrentMS > 0 {
					err = MoveUnitTowardsEnemyUM(world, id, uPos, ePos.PositionVectorX, ePos.PositionVectorY, eRadius.UnitRadius, uTeam, uMs, collisionHash, uRadius, uDist, mapName)
				}
				//if out of both attack and aggro range
			} else if adjustedDistance > float32(uAtk.AggroRadius) {
				uAtk.Combat = false
				uAtk.Frame = 0
				err = cardinal.SetComponent(world, id, uAtk)
				if err != nil {
					fmt.Printf("error setting Distance component on tempDistance (unit_movement.go): %v", err)
					continue
				}

				//in attack range just rotate towards enemy
			} else {
				err = RotateUnitTowardsEnemyUM(world, id, uPos, ePos.PositionVectorX, ePos.PositionVectorY)
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
					err = cardinal.SetComponent(world, id, uAtk)
					if err != nil {
						fmt.Printf("error setting Distance component on tempDistance (unit_movement.go): %v", err)
						continue
					}
					//not within attack range
				} else {
					if uMs.CurrentMS > 0 { // move towards enemy
						err = MoveUnitTowardsEnemyUM(world, id, uPos, eX, eY, eRadius, uTeam, uMs, collisionHash, uRadius, uDist, mapName)
					}
				}
			} else {
				//no enemies found and not in combat, move with direction map.
				if uMs.CurrentMS > 0 {
					err = MoveUnitDirectionMapUM(world, id, uPos, uTeam, uMs, collisionHash, uRadius, uDist, mapName)
				}
			}
		}
		if err != nil {
			fmt.Printf("(unit_movement.go): %s", err)
			continue
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
func MoveUnitDirectionMapUM(world cardinal.WorldContext, id types.EntityID, position *comp.Position, team *comp.Team, movespeed *comp.Movespeed, collisionSpartialHash *comp.SpatialHash, radius *comp.UnitRadius, distance *comp.Distance, mapName *comp.MapName) error {
	//check map data exsists
	mapData, exists := MapDataRegistry[mapName.MapName]
	if !exists {
		return fmt.Errorf("error key for MapDataRegistry does not exsist (unit_movement.go)")
	}
	//check direction map exsists
	mapDir, ok := MapRegistry[mapName.MapName]
	if !ok {
		return fmt.Errorf("error key for MapRegistry does not exsist (unit_movement.go)")
	}

	//normalize the units position to the maps grid increments.
	normalizedX := int(((int(position.PositionVectorX)-mapData.StartX)/mapData.Increment))*mapData.Increment + mapData.StartX
	normalizedY := int(((int(position.PositionVectorY)-mapData.StartY)/mapData.Increment))*mapData.Increment + mapData.StartY
	//The units (x,y) coordinates normalized and turned into proper key(string) format for seaching map
	coordKey := fmt.Sprintf("%d,%d", normalizedX, normalizedY)

	// Retrieve direction vector using coordinate key
	directionVector, exists := mapDir.DMap[coordKey]
	if !exists {
		return fmt.Errorf("no direction vector found for the given coordinates (unit_movement.go)")
	}
	//updated rotation based on team
	if team.Team == "Blue" {
		position.RotationVectorX = directionVector[0]
		position.RotationVectorY = directionVector[1]
	} else {
		position.RotationVectorX = directionVector[0] * -1 //reverse direction for red
		position.RotationVectorY = directionVector[1] * -1
	}

	// //Store Original X and Y
	tempX := position.PositionVectorX
	tempY := position.PositionVectorY
	//update new x,y based on movespeed
	position.PositionVectorX = position.PositionVectorX + (position.RotationVectorX * movespeed.CurrentMS)
	position.PositionVectorY = position.PositionVectorY + (position.RotationVectorY * movespeed.CurrentMS)
	//push or walk around blocking units
	position.PositionVectorX, position.PositionVectorY = UpdateUnitPositionPushSpatialHash(world, collisionSpartialHash, id, tempX, tempY, position.PositionVectorX, position.PositionVectorY, radius.UnitRadius, team.Team, movespeed.CurrentMS)
	//set updated position component
	err := cardinal.SetComponent(world, id, position)
	if err != nil {
		return fmt.Errorf("error set component on tempPosition (unit movement/MoveUnitDirectionMapUM): %v", err)
	}

	//update units new distance from enemy base
	if err = UpdateUnitDistance(world, id, team, distance, position, mapName); err != nil {
		return err
	}
	return nil
}

// Moves Unit towards enemy position
func MoveUnitTowardsEnemyUM(world cardinal.WorldContext, id types.EntityID, position *comp.Position, enemyX float32, enemyY float32, enemyRadius int, team *comp.Team, movespeed *comp.Movespeed, collisionSpartialHash *comp.SpatialHash, radius *comp.UnitRadius, distance *comp.Distance, mapName *comp.MapName) error {
	// Compute direction vector towards the enemy
	position.RotationVectorX, position.RotationVectorY = directionVectorBetweenTwoPoints(position.PositionVectorX, position.PositionVectorY, enemyX, enemyY)

	// Compute new position based on movespeed and direction
	newPosX := position.PositionVectorX + position.RotationVectorX*movespeed.CurrentMS
	newPosY := position.PositionVectorY + position.RotationVectorY*movespeed.CurrentMS

	// Calculate the stopping distance (combined radii of the unit and enemy plus 1 pixel for separation)
	stoppingDistance := radius.UnitRadius + enemyRadius + 1
	// Calculate the target position to move towards, stopping 1 pixel outside the enemy's radius
	targetX := enemyX - position.RotationVectorX*float32(stoppingDistance)
	targetY := enemyY - position.RotationVectorY*float32(stoppingDistance)

	// Ensure the unit does not overshoot the target position
	if (position.RotationVectorX > 0 && newPosX > targetX) || (position.RotationVectorX < 0 && newPosX < targetX) {
		newPosX = targetX
	}
	if (position.RotationVectorY > 0 && newPosY > targetY) || (position.RotationVectorY < 0 && newPosY < targetY) {
		newPosY = targetY
	}

	//push or walk around blocking units
	position.PositionVectorX, position.PositionVectorY = UpdateUnitPositionPushSpatialHash(world, collisionSpartialHash, id, position.PositionVectorX, position.PositionVectorY, newPosX, newPosY, radius.UnitRadius, team.Team, movespeed.CurrentMS)

	// Set the new position component
	err := cardinal.SetComponent(world, id, position)
	if err != nil {
		return fmt.Errorf("error set component on tempPosition (unit movement/MoveUnitTowardsEnemyUM): %v", err)
	}

	// Update units new distance from enemy base
	if err = UpdateUnitDistance(world, id, team, distance, position, mapName); err != nil {
		return err
	}
	return nil
}

// rotate Unit towards enemy position
func RotateUnitTowardsEnemyUM(world cardinal.WorldContext, id types.EntityID, position *comp.Position, enemyX float32, enemyY float32) error {
	// Compute direction vector towards the enemy
	position.RotationVectorX, position.RotationVectorY = directionVectorBetweenTwoPoints(position.PositionVectorX, position.PositionVectorY, enemyX, enemyY)

	// Set the new position component
	err := cardinal.SetComponent(world, id, position)
	if err != nil {
		return fmt.Errorf("error set component on tempPosition (unit movement/MoveUnitTowardsEnemyUM): %v", err)
	}
	return nil
}

// Update units distance from enemy base to help with movement priority queue
func UpdateUnitDistance(world cardinal.WorldContext, id types.EntityID, team *comp.Team, distance *comp.Distance, position *comp.Position, mapName *comp.MapName) error {
	//check map exsists in registy
	mapData, exists := MapDataRegistry[mapName.MapName]
	if !exists {
		return fmt.Errorf("error key for MapDataRegistry does not exsist (unit_movement.go)")
	}

	// calculate distance from enemy spawn
	if team.Team == "Blue" {
		distance.Distance = distanceBetweenTwoPoints(float32(mapData.Bases[1][0]), float32(mapData.Bases[1][1]), position.PositionVectorX, position.PositionVectorY)
	} else {
		distance.Distance = distanceBetweenTwoPoints(float32(mapData.Bases[0][0]), float32(mapData.Bases[0][1]), position.PositionVectorX, position.PositionVectorY)
	}
	// set distance
	err := cardinal.SetComponent(world, id, distance)
	if err != nil {
		return fmt.Errorf("error setting Distance component on UpdateUnitDistance (unit_movement.go): %v", err)
	}
	return nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////UTILITY FUNCTIONS//////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// GetUnitComponents fetches all necessary components related to a unit entity.
func GetUnitComponentsUM(world cardinal.WorldContext, unitID types.EntityID) (*comp.Position, *comp.UnitRadius, *comp.Attack, *comp.Team, *comp.Movespeed, *comp.MatchId, *comp.Distance, *comp.MapName, error) {
	position, err := cardinal.GetComponent[comp.Position](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Position component (unit_movement.go): %v", err)
	}
	unitRadius, err := cardinal.GetComponent[comp.UnitRadius](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Unit Radius component (unit_movement.go): %v", err)
	}
	unitAttack, err := cardinal.GetComponent[comp.Attack](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Unit attack component (unit_movement.go): %v", err)
	}
	team, err := cardinal.GetComponent[comp.Team](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Team component (unit_movement.go): %v", err)
	}
	movespeed, err := cardinal.GetComponent[comp.Movespeed](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Movespeed component (unit_movement.go): %v", err)
	}
	matchId, err := cardinal.GetComponent[comp.MatchId](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving MatchId component (unit_movement.go): %v", err)
	}
	distance, err := cardinal.GetComponent[comp.Distance](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Distance component (unit_movement.go): %v", err)
	}
	mapName, err := cardinal.GetComponent[comp.MapName](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Distance component (unit_movement.go): %v", err)
	}
	return position, unitRadius, unitAttack, team, movespeed, matchId, distance, mapName, nil
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
