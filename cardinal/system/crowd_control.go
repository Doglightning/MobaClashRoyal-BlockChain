package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

func applyKnockBack(world cardinal.WorldContext, id types.EntityID, hash *comp.SpatialHash, pos *comp.Position, dir *comp.Position, rad *comp.UnitRadius, team *comp.Team, class *comp.Class, mapName *comp.MapName, cc *comp.CC, push float32, angle, rotation float64) error {
	//can't push structures
	if class.Class != "structure" {

		tempX := pos.PositionVectorX + dir.RotationVectorX*push
		tempY := pos.PositionVectorY + dir.RotationVectorY*push
		//check that unit isnt walking through out of bounds towards a found unit
		if !moveDirectionExsist(tempX, tempY, mapName.MapName) {
			//find closest occupiable lication from target location to current
			tempX, tempY = pushFromPtBtoA(world, hash, id, pos.PositionVectorX, pos.PositionVectorY, tempX, tempY, rad.UnitRadius, mapName)

			if distanceBetweenTwoPoints(pos.PositionVectorX, pos.PositionVectorY, tempX, tempY) <= 0 {

				position := &comp.Position{PositionVectorX: pos.PositionVectorX, PositionVectorY: pos.PositionVectorY}

				findInboundsRotation(angle, rotation, position, dir.RotationVectorX, dir.RotationVectorY, push, mapName)

				tempX = position.PositionVectorX
				tempY = position.PositionVectorY

			}
		}
		RemoveObjectFromSpatialHash(hash, id, pos.PositionVectorX, pos.PositionVectorX, rad.UnitRadius)

		//attempt to push blocking units
		pushBlockingUnit(world, hash, id, tempX, tempY, rad.UnitRadius, team.Team, class.Class, push, mapName)
		//move unit.  walk around blocking units
		pos.PositionVectorX, pos.PositionVectorY = moveFreeSpace(hash, id, pos.PositionVectorX, pos.PositionVectorY, tempX, tempY, rad.UnitRadius, team.Team, class.Class, mapName)
		AddObjectSpatialHash(hash, id, pos.PositionVectorX, pos.PositionVectorY, rad.UnitRadius, team.Team, class.Class)

		// Update units new distance from enemy base
		if err := updateUnitDistance(world, id, team, pos, mapName); err != nil {
			return fmt.Errorf("(applyKnockBack): %v", err)
		}

		cc.KnockBack = true
	}
	return nil
}

func applyKnockUp(world cardinal.WorldContext, id types.EntityID, matchID *comp.MatchId, targetHeight, speed, damage float32) error {

	//get new uid
	UID, err := getNextUID(world, matchID.MatchId)
	if err != nil {
		return fmt.Errorf("(applyKknockUp - crowd control.go): %v - ", err)
	}
	//create stun entity attached to target
	_, err = cardinal.Create(world,
		comp.MatchId{MatchId: matchID.MatchId},
		comp.UID{UID: UID},
		comp.SpEntity{SpName: "KnockUp"},
		comp.KnockUp{CurrentHieght: 0, TargetHieght: targetHeight, Speed: speed, Damage: damage, ApexReached: false},
		comp.Target{Target: id},
	)
	if err != nil {
		return fmt.Errorf("error creating knockup entity (applyKknockUp - crowd control.go): %v", err)
	}

	return nil
}

func knockUpUpdate(world cardinal.WorldContext, id types.EntityID) error {

	tarID, knockUp, err := GetComponents2[comp.Target, comp.KnockUp](world, id)
	if err != nil {
		return fmt.Errorf("knock up entity (knockUpUpdate - crowd control.go): %w", err)
	}

	pos, err := cardinal.GetComponent[comp.Position](world, tarID.Target)
	if err != nil {
		return fmt.Errorf("knock up entity (knockUpUpdate - crowd control.go): %w", err)
	}

	if !knockUp.ApexReached { // hasn't reached apex
		knockUp.CurrentHieght += knockUp.Speed // increase current hieght

		if knockUp.CurrentHieght < knockUp.TargetHieght { // if hieght hasn't reached apex
			pos.PositionVectorZ += knockUp.Speed

		} else if knockUp.CurrentHieght > knockUp.TargetHieght { // passed apex
			difference := knockUp.Speed - (knockUp.CurrentHieght - knockUp.TargetHieght) // find amount to increment so we don't over shoot
			pos.PositionVectorZ += difference
			knockUp.CurrentHieght = knockUp.TargetHieght
			knockUp.ApexReached = true
		} else { // at apex
			pos.PositionVectorZ += knockUp.Speed
			knockUp.ApexReached = true
		}
	} else { // has reached apex
		knockUp.CurrentHieght -= knockUp.Speed // decrease current hieght

		if knockUp.CurrentHieght > 0 { // if hieght hasn't reached bottom
			pos.PositionVectorZ -= knockUp.Speed

		} else if knockUp.CurrentHieght <= 0 { // passed bottom
			difference := knockUp.Speed + (knockUp.CurrentHieght) // find amount to increment so we don't over shoot
			pos.PositionVectorZ -= difference

			if err := cardinal.SetComponent(world, tarID.Target, pos); err != nil {
				return fmt.Errorf("error setting pos comp (knockUpUpdate - crowd control.go): %w", err)
			}

			if err := applyDamage(world, tarID.Target, knockUp.Damage); err != nil {
				return fmt.Errorf(" (knockUpUpdate - crowd control.go): %w", err)
			}

			// delete entity
			if err := cardinal.Remove(world, id); err != nil {
				return fmt.Errorf("error removing entity sp (knockUpUpdate - crowd control.go): %w", err)
			}

		}
	}
	if err := cardinal.SetComponent(world, id, knockUp); err != nil {
		return fmt.Errorf("error setting knock up comp (knockUpUpdate - crowd control.go): %w", err)
	}

	if err := cardinal.SetComponent(world, tarID.Target, pos); err != nil {
		return fmt.Errorf("error setting pos comp (knockUpUpdate - crowd control.go): %w", err)
	}

	return nil
}
