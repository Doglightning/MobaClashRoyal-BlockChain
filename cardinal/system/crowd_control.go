package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

func applyKnockBack(world cardinal.WorldContext, id types.EntityID, hash *comp.SpatialHash, pos *comp.Position, dir *comp.Position, rad *comp.UnitRadius, team *comp.Team, class *comp.Class, mapName *comp.MapName, cc *comp.CC, push float32) error {
	//can't push structures
	if class.Class != "structure" {

		tempX := pos.PositionVectorX + dir.RotationVectorX*push
		tempY := pos.PositionVectorY + dir.RotationVectorY*push
		//check that unit isnt walking through out of bounds towards a found unit
		if !moveDirectionExsist(tempX, tempY, mapName.MapName) {
			//find closest occupiable lication from target location to current
			tempX, tempY = pushFromPtBtoA(world, hash, id, pos.PositionVectorX, pos.PositionVectorY, tempX, tempY, rad.UnitRadius, mapName)
		}
		RemoveObjectFromSpatialHash(hash, id, pos.PositionVectorX, pos.PositionVectorX, rad.UnitRadius)

		//attempt to push blocking units
		pushBlockingUnit(world, hash, id, tempX, tempY, rad.UnitRadius, class.Class, team.Team, push, mapName)
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
