package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

var numArrows int = 5
var arrowSeperationDegree float64 = 30

func archerLadyAttack(world cardinal.WorldContext, id types.EntityID) error {
	uPos, matchID, mapName, err := GetSpComponentsAL(world, id)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	vectors := generateVectors(uPos.RotationVectorX, uPos.RotationVectorY, arrowSeperationDegree, numArrows)

	for i := 0; i < numArrows; i++ {

		//get next uid
		UID, err := getNextUID(world, matchID.MatchId)
		if err != nil {
			return fmt.Errorf("(archerLadyAttack): %v - ", err)
		}
		//create projectile entity
		cardinal.Create(world,
			comp.MatchId{MatchId: matchID.MatchId},
			comp.UID{UID: UID},
			comp.SpName{SpName: "ArcherLadySp"},
			comp.Movespeed{CurrentMS: ProjectileRegistry["ArcherLadySp"].Speed},
			comp.Position{PositionVectorX: uPos.PositionVectorX, PositionVectorY: uPos.PositionVectorY, PositionVectorZ: uPos.PositionVectorZ, RotationVectorX: vectors[i][0], RotationVectorY: vectors[i][1], RotationVectorZ: uPos.RotationVectorZ},
			comp.MapName{MapName: mapName.MapName},
			comp.Damage{Damage: ProjectileRegistry["ArcherLadySp"].Damage},
			comp.Destroyed{Destroyed: false},
		)
	}

	return err
}

// generateVectors generates evenly distributed vectors within a given angle around the central vector
// dirX, dirY: central Direction vector (normalized)
// angle: angle between vectors
// count: Number of vectors
func generateVectors(dirX, dirY float32, angle float64, count int) [][]float32 {
	halfAngle := angle / 2
	stepAngle := angle / float64(count-1)

	vectors := make([][]float32, count)
	for i := 0; i < count; i++ {
		currentAngle := -halfAngle + stepAngle*float64(i)
		rotatedX, rotatedY := rotateVectorDegrees(dirX, dirY, currentAngle)
		vectors[i] = []float32{rotatedX, rotatedY}
	}
	return vectors
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
