package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

// leafBirdSP struct contains configuration for an leafBirdSP in terms of their shooting properties.
type lavaGolemDataSP struct {
	Radius float32
	Push   float32
	Speed  float32
	Damage float32
}

// NewleafBirdSPSP creates a new instance of leafBirdSPSP with default settings.
func NewLavaGolemDataSP() *lavaGolemDataSP {
	return &lavaGolemDataSP{
		Radius: float32(UnitRegistry["LavaGolem"].Radius) + float32(SpRegistry["LavaGolem"].AttackRadius),
		Push:   550,
		Speed:  150,
		Damage: 30,
	}
}

// updates SP entity per tick
func vampireUpdateSPs(world cardinal.WorldContext, id types.EntityID) error {

	return nil
}

// spawning the vampire special power
func lavaGolemSpawnSP(world cardinal.WorldContext, id types.EntityID) error {

	//get fire spirit vars
	lavaGolem := NewLavaGolemDataSP()

	//get unit comps
	team, matchID, pos, err := GetComponents3[comp.Team, comp.MatchId, comp.Position](world, id)
	if err != nil {
		return fmt.Errorf("error getting unit comps (lavaGolemSp): %v", err)
	}

	//get collision hash
	gameStateID, hash, err := getCollisionHashAndGameState(world, matchID)
	if err != nil {
		return fmt.Errorf("(lavaGolemSp): %v", err)
	}

	//find the 4 corners of the AoE rectangle
	topLeft, topRight, botLeft, botRight := CreateRectangleAroundPoint(Point{X: pos.PositionVectorX, Y: pos.PositionVectorY}, Point{X: pos.RotationVectorX, Y: pos.RotationVectorY}, lavaGolem.Radius, lavaGolem.Radius)
	//find the rectangle that contains our AoE normalized to the (x, y) coord system
	_, topRightA, botLeftB, _ := FindRectangleAABB(topLeft, topRight, botLeft, botRight)

	// Define a map to track unique collisions
	collidedEntities := make(map[types.EntityID]bool)

	if SpatialGridCellSize <= 0 {
		return fmt.Errorf("invalid SpatialGridCellSize (lavaGolemSp)")
	}

	// Loop through all `x` values from the min to max x, stepping by `stepSize`.
	for x := botLeftB.X; x <= topRightA.X+float32(SpatialGridCellSize); x += float32(SpatialGridCellSize) {
		// Loop through all `y` values from the min to max y, stepping by `stepSize`.
		for y := botLeftB.Y; y <= topRightA.Y+float32(SpatialGridCellSize); y += float32(SpatialGridCellSize) {

			collList := GetEntitiesInCell(hash, x, y) //list of all units in cell
			for _, collID := range collList {         //for each collision
				collidedEntities[collID] = true //add to map
			}
		}
	}

	// Iterate over each key in the map
	for collID := range collidedEntities {
		if collID == id {
			continue
		}

		//get target team and class components
		targetTeam, targetClass, err := GetComponents2[comp.Team, comp.Class](world, collID)
		if err != nil {
			fmt.Printf("error getting targets compoenents (lavaGolemSp): %v \n", err)
			continue
		}

		if team.Team != targetTeam.Team && targetClass.Class != "structure" { //dont attack friendlies soilder!!

			//get target position and radius components
			targetPos, targetRad, err := GetComponents2[comp.Position, comp.UnitRadius](world, collID)
			if err != nil {
				fmt.Printf("error getting targets compoenents (lavaGolemSp): %v \n", err)
				continue
			}

			//does our unit intersect the AoE
			if CircleIntersectsRectangle(Point{X: targetPos.PositionVectorX, Y: targetPos.PositionVectorY}, float32(targetRad.UnitRadius), topLeft, topRight, botRight, botLeft) {

				if err := applyKnockUp(world, collID, matchID, lavaGolem.Push, lavaGolem.Speed, lavaGolem.Damage); err != nil {
					return fmt.Errorf("(lavaGolemSp): %s ", err)
				}

			}
		}
	}
	// update hash
	if err := cardinal.SetComponent(world, gameStateID, hash); err != nil {
		return fmt.Errorf("error setting hash (lavaGolemSp): %s ", err)
	}

	return nil
}

func lavaGolemAttack(world cardinal.WorldContext, atk *comp.Attack) error {
	// reduce health by units attack damage
	err := cardinal.UpdateComponent(world, atk.Target, func(health *comp.Health) *comp.Health {
		if health == nil {
			fmt.Printf("error retrieving Health component (Unit_Attack.go) \n")
			return nil
		}
		health.CurrentHP -= float32(atk.Damage)
		if health.CurrentHP < 0 {
			health.CurrentHP = 0 //never have negative health
		}
		return health
	})
	if err != nil {
		return fmt.Errorf("error on vampire attack (class vampire.go): %v", err)
	}

	return nil
}
