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
		//check that unit isnt walking through out of bounds towards a found unit
		if moveDirectionExsist(pos.PositionVectorX, pos.PositionVectorY, mapName.MapName) {
			tempX := pos.PositionVectorX
			tempY := pos.PositionVectorY

			RemoveObjectFromSpatialHash(hash, id, pos.PositionVectorX, pos.PositionVectorX, rad.UnitRadius)

			pos.PositionVectorX += dir.RotationVectorX * push
			pos.PositionVectorY += dir.RotationVectorY * push

			//attempt to push blocking units
			pushBlockingUnit(world, hash, id, pos.PositionVectorX, pos.PositionVectorY, rad.UnitRadius, class.Class, team.Team, push, mapName)
			//move unit.  walk around blocking units
			pos.PositionVectorX, pos.PositionVectorY = moveFreeSpace(hash, id, tempX, tempY, pos.PositionVectorX, pos.PositionVectorY, rad.UnitRadius, team.Team, class.Class, mapName)
			AddObjectSpatialHash(hash, id, pos.PositionVectorX, pos.PositionVectorY, rad.UnitRadius, team.Team, class.Class)

			// Update units new distance from enemy base
			if err := updateUnitDistance(world, id, team, pos, mapName); err != nil {
				return fmt.Errorf("(applyKnockBack): %v", err)
			}

			cc.KnockBack = true
			// update hash and position
			if err := SetComponents2(world, id, pos, cc); err != nil {
				return fmt.Errorf("(applyKnockBack): %s ", err)
			}

		}
	}
	return nil
}
