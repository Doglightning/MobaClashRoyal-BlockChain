package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

// fireSpiritSpawnSP struct contains configuration for an fire spirit in terms of their shooting properties.
type leafBirdSpawnSP struct {
	Hieght    float32 //triangle hieght
	BaseWidth float32 //triangle base width
	Push      float32
	Damage    float32
}

// NewFireSpiritSpawnSP creates a new instance of NewFireSpiritSP with default settings.
func NewLeafBirdSpawnSP() *leafBirdSpawnSP {
	return &leafBirdSpawnSP{
		Hieght:    925,
		BaseWidth: 125,
		Push:      50,
		Damage:    1.4,
	}
}

// shoots the fire attack every frame
func leafBirdSp(world cardinal.WorldContext, id types.EntityID) error {
	//get fire spirit vars
	leafBird := NewLeafBirdSpawnSP()

	//get team comp
	team, matchID, mapName, pos, err := GetComponents4[comp.Team, comp.MatchId, comp.MapName, comp.Position](world, id)
	if err != nil {
		return fmt.Errorf("error getting unit comps (leafBirdSp): %v", err)
	}

	//get collision hash
	gameStateID, hash, err := getCollisionHashAndGameState(world, matchID)
	if err != nil {
		return fmt.Errorf("error getting spatial hash compoenent (leafBirdSp): %v", err)
	}

	//find the 3 points of the fire spirit AoE triangle attack
	topLeft, topRight, botLeft, botRight := CreateRectangleBase(Point{X: pos.PositionVectorX, Y: pos.PositionVectorY}, Point{X: pos.RotationVectorX, Y: pos.RotationVectorY}, leafBird.BaseWidth, leafBird.Hieght)

	_, topRightA, botLeftB, _ := FindRectangleAABB(topLeft, topRight, botLeft, botRight)

	// Define a map to track unique collisions
	collidedEntities := make(map[types.EntityID]bool)

	if SpatialGridCellSize <= 0 {
		return fmt.Errorf("invalid SpatialGridCellSize (leafBirdSp)")
	}

	// Loop through all `x` values from the min to max x, stepping by `stepSize`.
	for x := botLeftB.X; x <= topRightA.X; x += float32(SpatialGridCellSize) {
		// Loop through all `y` values from the min to max y, stepping by `stepSize`.
		for y := botLeftB.Y; y <= topRightA.Y; y += float32(SpatialGridCellSize) {

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

		//get targets team
		targetTeam, targetClass, err := GetComponents2[comp.Team, comp.Class](world, collID)
		if err != nil {
			fmt.Printf("error getting targets compoenents (leafBirdSp): %v \n", err)
			continue
		}

		if team.Team != targetTeam.Team { //dont attack friendlies soilder!!

			targetPos, targetRad, err := GetComponents2[comp.Position, comp.UnitRadius](world, collID)
			if err != nil {
				fmt.Printf("error getting targets compoenents (leafBirdSp): %v \n", err)
				continue
			}

			if CircleIntersectsRectangle(Point{X: targetPos.PositionVectorX, Y: targetPos.PositionVectorY}, float32(targetRad.UnitRadius), topLeft, topRight, botRight, botLeft) {

				if targetClass.Class != "structure" { // cant push structures

					targetCC, err := cardinal.GetComponent[comp.CC](world, collID)
					if err != nil {
						return fmt.Errorf("(leafBirdSp) -  %s ", err)
					}

					// apply knock back
					if err := applyKnockBack(world, collID, hash, targetPos, pos, targetRad, targetTeam, targetClass, mapName, targetCC, leafBird.Push); err != nil {
						return fmt.Errorf("(leafBirdSp) -  %s ", err)
					}

				}

				//apply damage
				if err = applyDamage(world, collID, leafBird.Damage); err != nil {
					return fmt.Errorf("(leafBirdSp) - %v", err)
				}

			}

		}
	}
	// update hash
	if err := cardinal.SetComponent(world, gameStateID, hash); err != nil {
		return fmt.Errorf("error setting hash (leafBirdSp): %s ", err)
	}

	return nil
}

// overwrite phase_attack.go logic to support canneling
func leafBirdAttackSystem(world cardinal.WorldContext, id types.EntityID, atk *comp.Attack) error {

	//get Unit CC component
	cc, err := cardinal.GetComponent[comp.CC](world, id)
	if err != nil {
		fmt.Printf("error getting unit cc component ( leafBirdAttackSystem): %v", err)
	}

	if cc.Stun > 0 { //if unit stunned cannot attack
		return nil
	}

	//get special power component
	unitSp, err := cardinal.GetComponent[comp.Sp](world, id)
	if err != nil {
		return fmt.Errorf("error retrieving special power component ( leafBirdAttackSystem): %v", err)
	}

	//check if in a SP animation or a regular attack
	if atk.Frame == 0 && unitSp.CurrentSp >= unitSp.MaxSp { //In special power

		unitSp.Charged = true

	} else if atk.Frame == 0 && unitSp.CurrentSp < unitSp.MaxSp { // in regular attack
		unitSp.Charged = false
	}

	//if unit is in its damage frame and not charged
	if atk.Frame == atk.DamageFrame && !unitSp.Charged {

		//peck em >:D
		err = applyDamage(world, atk.Target, atk.Damage)
		if err != nil {
			return fmt.Errorf("(leafBirdAttackSystem): %v", err)
		}

		unitSp.CurrentSp += unitSp.SpRate //increase sp after attack
		// make sure we are not over MaxSp
		if unitSp.CurrentSp >= unitSp.MaxSp {
			unitSp.CurrentSp = unitSp.MaxSp
		}
	}

	//if unit is in damage frames when charged
	if unitSp.DamageFrame <= atk.Frame && atk.Frame <= unitSp.DamageEndFrame && unitSp.Charged {
		atk.State = "Channeling"

		//SHOT AIR BIOTCH >:D
		err = leafBirdSp(world, id)
		if err != nil {
			return err
		}

		//return Sp to 0
		unitSp.CurrentSp = 0

	}

	//if target died in cast (self target) and attack frame is at end of animation or start (don't interupt the fire strike once its going even if target died)
	if (atk.Target == id && atk.Frame >= unitSp.Rate) || (atk.Target == id && atk.Frame < unitSp.DamageFrame) || (atk.Frame >= unitSp.Rate && atk.State == "Channeling") {
		atk.State = "Default"
		atk.Combat = false
	}

	//if attack frame is at max and not sp charged  OR attack fram at sp max and charged
	if (atk.Frame >= atk.Rate && !unitSp.Charged) || (atk.Frame >= unitSp.Rate && unitSp.Charged) {
		atk.Frame = -1 //reset attack
	}
	atk.Frame++

	// set updated attack component
	if err := cardinal.SetComponent(world, id, atk); err != nil {
		return fmt.Errorf("error updating attack component ( leafBirdAttackSystem): %s ", err)
	}
	// set updated sp component
	if err := cardinal.SetComponent(world, id, unitSp); err != nil {
		return fmt.Errorf("error updating special power component ( leafBirdAttackSystem): %s ", err)
	}

	return nil
}
