package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

// spawns projectile for tower basic attack
func towerAttack(world cardinal.WorldContext, id types.EntityID, atk *comp.Attack) error {
	//get units component
	unitPosition, matchID, mapName, unitName, err := GetComponents4[comp.Position, comp.MatchId, comp.MapName, comp.UnitName](world, id) //reusing
	if err != nil {
		return fmt.Errorf("tower components (class towerAttack.go): %v ", err)
	}
	//get next uid
	UID, err := getNextUID(world, matchID.MatchId)
	if err != nil {
		return fmt.Errorf("(class towerAttack.go): %v ", err)
	}
	//create projectile entity
	_, err = cardinal.Create(world,
		comp.MatchId{MatchId: matchID.MatchId},
		comp.UID{UID: UID},
		comp.UnitName{UnitName: ProjectileRegistry[unitName.UnitName].Name},
		comp.Movespeed{CurrentMS: ProjectileRegistry[unitName.UnitName].Speed},
		comp.Position{
			PositionVectorX: unitPosition.PositionVectorX,
			PositionVectorY: unitPosition.PositionVectorY,
			PositionVectorZ: unitPosition.PositionVectorZ + ProjectileRegistry[unitName.UnitName].offSetZ,
			RotationVectorX: unitPosition.RotationVectorX,
			RotationVectorY: unitPosition.RotationVectorY,
			RotationVectorZ: unitPosition.RotationVectorZ},
		comp.MapName{MapName: mapName.MapName},
		comp.Attack{Target: atk.Target, Class: "projectile", Damage: StructureDataRegistry[unitName.UnitName].Damage},
		comp.Destroyed{Destroyed: false},
	)

	if err != nil {
		return fmt.Errorf("error spawning tower basic attack (class towerAttack.go): %v ", err)
	}

	return nil
}
