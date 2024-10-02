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
	push      float32
}

// NewFireSpiritSpawnSP creates a new instance of NewFireSpiritSP with default settings.
func NewLeafBirdSpawnSP() *leafBirdSpawnSP {
	return &leafBirdSpawnSP{
		Hieght:    925,
		BaseWidth: 125,
		push:      50,
	}
}

// shoots the fire attack every frame
func leafBirdSpawn(world cardinal.WorldContext, id types.EntityID) error {
	//get fire spirit vars
	leafBird := NewLeafBirdSpawnSP()

	//get team comp
	team, err := cardinal.GetComponent[comp.Team](world, id)
	if err != nil {
		return fmt.Errorf("error getting team component (class fireSpirit.go): %v", err)
	}

	//get matchID component
	matchID, err := cardinal.GetComponent[comp.MatchId](world, id)
	if err != nil {
		return fmt.Errorf("error getting position component (class fireSpirit.go): %v", err)
	}

	// //get matchID component
	// mapName, err := cardinal.GetComponent[comp.MapName](world, id)
	// if err != nil {
	// 	return fmt.Errorf("error getting map component (class fireSpirit.go): %v", err)
	// }

	//get collision hash
	hash, err := getCollisionHashGSS(world, matchID)
	if err != nil {
		return fmt.Errorf("error getting spatial hash compoenent(class fireSpirit.go): %v", err)
	}

	//get position comp
	pos, err := cardinal.GetComponent[comp.Position](world, id)
	if err != nil {
		return fmt.Errorf("error getting position component (class fireSpirit.go): %v", err)
	}

	//find the 3 points of the fire spirit AoE triangle attack
	topLeft, topRight, botLeft, botRight := CreateRectangleBase(Point{X: pos.PositionVectorX, Y: pos.PositionVectorY}, Point{X: pos.RotationVectorX, Y: pos.RotationVectorY}, leafBird.BaseWidth, leafBird.Hieght)

	_, topRightA, botLeftB, _ := FindRectangleAABB(topLeft, topRight, botLeft, botRight)

	// Define a map to track unique collisions
	collidedEntities := make(map[types.EntityID]bool)

	if SpatialGridCellSize <= 0 {
		return fmt.Errorf("invalid SpatialGridCellSize (class fireSpirit.go)")
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
		targetTeam, err := cardinal.GetComponent[comp.Team](world, collID)
		if err != nil {
			fmt.Printf("error getting targets team compoenent (class fireSpirit.go): %v \n", err)
			continue
		}

		targetClass, err := cardinal.GetComponent[comp.Class](world, collID)
		if err != nil {
			fmt.Printf("error getting targets attack compoenent (class fireSpirit.go): %v \n", err)
			continue
		}

		if team.Team != targetTeam.Team && targetClass.Class != "structure" { //dont attack friendlies soilder!!

			targetPos, err := cardinal.GetComponent[comp.Position](world, collID)
			if err != nil {
				fmt.Printf("error getting targets Position compoenent (class fireSpirit.go): %v \n", err)
				continue
			}

			targetRad, err := cardinal.GetComponent[comp.UnitRadius](world, collID)
			if err != nil {
				fmt.Printf("error getting targets radius compoenent (class fireSpirit.go): %v \n", err)
				continue
			}

			if CircleIntersectsRectangle(Point{X: targetPos.PositionVectorX, Y: targetPos.PositionVectorY}, float32(targetRad.UnitRadius), topLeft, topRight, botRight, botLeft) {

				// // tempX := targetPos.PositionVectorX
				// // tempY := targetPos.PositionVectorY

				RemoveObjectFromSpatialHash(hash, collID, targetPos.PositionVectorX, targetPos.PositionVectorX, targetRad.UnitRadius)

				targetPos.PositionVectorX += pos.RotationVectorX * leafBird.push
				targetPos.PositionVectorY += pos.RotationVectorY * leafBird.push

				// // // //check that unit isnt walking through out of bounds towards a found unit
				// // exists := moveDirectionExsist(targetPos.PositionVectorX, targetPos.PositionVectorY, mapName)
				// // if exists {
				// // set updated sp component
				// if err := cardinal.SetComponent(world, collID, targetPos); err != nil {
				// 	return fmt.Errorf("error updating special power component (Fire Spirit Attack): %s ", err)
				// }
				// set updated sp component
				if err := cardinal.SetComponent(world, collID, targetPos); err != nil {
					return fmt.Errorf("error 11updating special power component (Fire Spirit Attack): %s ", err)
				}
				AddObjectSpatialHash(hash, collID, targetPos.PositionVectorX, targetPos.PositionVectorX, targetRad.UnitRadius, targetTeam.Team, targetClass.Class)
				// // 	// 	//attempt to push blocking units
				// // 	// 	pushBlockingUnit(world, hash, collID, targetPos.PositionVectorX, targetPos.PositionVectorY, targetRad.UnitRadius, targetAtk.Class, targetTeam.Team, leafBird.push, mapName)
				// // 	// 	//move unit.  walk around blocking units
				// // 	// 	targetPos.PositionVectorX, targetPos.PositionVectorY = moveFreeSpace(hash, collID, tempX, tempY, targetPos.PositionVectorX, targetPos.PositionVectorY, targetRad.UnitRadius, targetTeam.Team, targetAtk.Class, mapName)

				// // 	// 	// Update units new distance from enemy base
				// // 	// 	if err = updateUnitDistance(world, collID, targetTeam, targetPos, mapName); err != nil {
				// // 	// 		fmt.Printf("%v \n", err)
				// // 	// 		return nil
				// // 	// 	}

				// // } else {
				// // AddObjectSpatialHash(hash, collID, tempX, tempY, targetRad.UnitRadius, targetTeam.Team, targetAtk.Class)
				// // }

				// // set updated sp component
				// if err := cardinal.SetComponent(world, collID, targetPos); err != nil {
				// 	return fmt.Errorf("error 11updating special power component (Fire Spirit Attack): %s ", err)
				// }

				// // reduce health by units attack damage
				// err := cardinal.UpdateComponent(world, collID, func(health *comp.Health) *comp.Health {
				// 	if health == nil {
				// 		fmt.Printf("error retrieving Health component (leafBirdAttack) \n")
				// 		return nil
				// 	}
				// 	health.CurrentHP -= float32(1)
				// 	if health.CurrentHP < 0 {
				// 		health.CurrentHP = 0 //never have negative health
				// 	}
				// 	return health
				// })
				// if err != nil {
				// 	return fmt.Errorf("error on leafbird attack (leafBirdAttack): %v", err)
				// }

			}

		}
	}
	gs, _ := getGameStateGSS(world, matchID)
	// set updated sp component
	if err := cardinal.SetComponent(world, gs, hash); err != nil {
		return fmt.Errorf("error updating special power component (Fire Spirit Attack): %s ", err)
	}

	return nil
}

func leafBirdAttack(world cardinal.WorldContext, atk *comp.Attack) error {
	// reduce health by units attack damage
	err := cardinal.UpdateComponent(world, atk.Target, func(health *comp.Health) *comp.Health {
		if health == nil {
			fmt.Printf("error retrieving Health component (leafBirdAttack) \n")
			return nil
		}
		health.CurrentHP -= float32(atk.Damage)
		if health.CurrentHP < 0 {
			health.CurrentHP = 0 //never have negative health
		}
		return health
	})
	if err != nil {
		return fmt.Errorf("error on leafbird attack (leafBirdAttack): %v", err)
	}

	return nil
}

// overwrite phase_attack.go logic to support canneling
func leafBirdAttackSystem(world cardinal.WorldContext, id types.EntityID, atk *comp.Attack) error {

	//get Unit CC component
	cc, err := cardinal.GetComponent[comp.CC](world, id)
	if err != nil {
		fmt.Printf("error getting unit cc component (Fire Spirit Attack): %v", err)
	}

	if cc.Stun > 0 { //if unit stunned cannot attack
		return nil
	}

	//get special power component
	unitSp, err := cardinal.GetComponent[comp.Sp](world, id)
	if err != nil {
		return fmt.Errorf("error retrieving special power component (Fire Spirit Attack): %v", err)
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
		err = leafBirdAttack(world, atk)
		if err != nil {
			return err
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
		err = leafBirdSpawn(world, id)
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
		return fmt.Errorf("error updating attack component (Fire Spirit Attack): %s ", err)
	}
	// set updated sp component
	if err := cardinal.SetComponent(world, id, unitSp); err != nil {
		return fmt.Errorf("error updating special power component (Fire Spirit Attack): %s ", err)
	}

	return nil
}
