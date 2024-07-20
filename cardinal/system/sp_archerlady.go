package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

func archerLadyAttack(world cardinal.WorldContext, id types.EntityID) error {
	unitPosition, matchID, mapName, err := GetSpComponentsAL(world, id)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	//get next uid
	UID, err := getNextUID(world, matchID.MatchId)
	if err != nil {
		return fmt.Errorf("(archerLadyAttack): %v - ", err)
	}
	//create projectile entity
	cardinal.Create(world,
		comp.MatchId{MatchId: matchID.MatchId},
		comp.UID{UID: UID},
		comp.SpName{SpName: "ArcherLadyArrow"},
		comp.Movespeed{CurrentMS: ProjectileRegistry["ArcherLadyArrow"].Speed},
		comp.Position{PositionVectorX: unitPosition.PositionVectorX, PositionVectorY: unitPosition.PositionVectorY, PositionVectorZ: unitPosition.PositionVectorZ, RotationVectorX: unitPosition.RotationVectorX, RotationVectorY: unitPosition.RotationVectorY, RotationVectorZ: unitPosition.RotationVectorZ},
		comp.MapName{MapName: mapName.MapName},
		comp.Damage{Damage: ProjectileRegistry["ArcherLadyArrow"].Damage},
		comp.Destroyed{Destroyed: false},
	)

	return err
}

// GetSpComponentsAL fetches all necessary components related to a sp entity.
func GetSpComponentsAL(world cardinal.WorldContext, unitID types.EntityID) (*comp.Position, *comp.MatchId, *comp.MapName, error) {
	position, err := cardinal.GetComponent[comp.Position](world, unitID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error retrieving Position component (sp_archerlay.go): %v", err)
	}

	matchId, err := cardinal.GetComponent[comp.MatchId](world, unitID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error retrieving MatchId component (sp_archerlay.go): %v", err)
	}

	mapName, err := cardinal.GetComponent[comp.MapName](world, unitID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error retrieving mapname component (sp_archerlay.go): %v", err)
	}

	return position, matchId, mapName, nil
}
