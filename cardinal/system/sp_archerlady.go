package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

var numArrows int = 5
var arrowSeperationDegree float64 = 25
var distance float32 = 1600
var radiusArrows = 150

func archerLadySpawn(world cardinal.WorldContext, id types.EntityID) error {
	uPos, matchID, mapName, team, err := GetSpComponentsAL(world, id)
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
			comp.SpName{SpName: "ArcherLady"},
			comp.Movespeed{CurrentMS: ProjectileRegistry["ArcherLady"].Speed},
			comp.Position{PositionVectorX: uPos.PositionVectorX, PositionVectorY: uPos.PositionVectorY, PositionVectorZ: uPos.PositionVectorZ, RotationVectorX: vectors[i][0], RotationVectorY: vectors[i][1], RotationVectorZ: uPos.RotationVectorZ},
			comp.MapName{MapName: mapName.MapName},
			comp.Damage{Damage: ProjectileRegistry["ArcherLady"].Damage},
			comp.Destroyed{Destroyed: false},
			comp.Distance{Distance: distance},
			comp.Team{Team: team.Team},
		)
	}

	return err
}

func archerLadyUpdate(world cardinal.WorldContext, id types.EntityID) error {

	//update Sp location
	cardinal.UpdateComponent(world, id, func(pos *comp.Position) *comp.Position {
		if pos == nil {
			fmt.Printf("error retrieving enemy position component (sp_archerlady.go): ")
			return nil
		}

		ms, dist, mID, team, dmg, err := GetUpdateComponentsAL(world, id)
		if err != nil {
			fmt.Printf("%v", err)
			return nil
		}
		collisionHash, err := getCollisionHashGSS(world, mID)
		if err != nil {
			fmt.Printf("(sp_archerlady.go) - ")
			return nil
		}

		cID := CheckCollisionSpatialHashList(collisionHash, pos.PositionVectorX, pos.PositionVectorY, radiusArrows)

		for _, value := range cID {
			cTeam, err := cardinal.GetComponent[comp.Team](world, value)
			if err != nil {
				fmt.Printf("(error getting team component (sp_archerlady.go): %v", err)
				return nil
			}

			if team.Team != cTeam.Team {
				cPos, err := cardinal.GetComponent[comp.Position](world, value)
				if err != nil {
					fmt.Printf("(error getting team component (sp_archerlady.go): %v", err)
					return nil
				}

				cRad, err := cardinal.GetComponent[comp.UnitRadius](world, value)
				if err != nil {
					fmt.Printf("(error getting team component (sp_archerlady.go): %v", err)
					return nil
				}

				if checkLineIntersectionSpatialHash(pos.PositionVectorX, pos.PositionVectorY, pos.PositionVectorX+ms.CurrentMS*pos.RotationVectorX, pos.PositionVectorY+ms.CurrentMS*pos.RotationVectorY, cPos.PositionVectorX, cPos.PositionVectorY, cRad.UnitRadius) {
					cardinal.UpdateComponent(world, id, func(destroyed *comp.Destroyed) *comp.Destroyed {
						if destroyed == nil {
							fmt.Printf("error retrieving enemy destroyed component (sp_archerlady.go): ")
							return nil
						}
						destroyed.Destroyed = true
						return destroyed
					})
					//update projectiles destroyed component to True
					cardinal.UpdateComponent(world, value, func(health *comp.Health) *comp.Health {
						if health == nil {
							fmt.Printf("error retrieving collision health component (sp_archerlady.go): ")
							return nil
						}
						health.CurrentHP -= float32(dmg.Damage)
						if health.CurrentHP < 0 {
							health.CurrentHP = 0
						}
						return health
					})
				}
			}
		}

		//updated position and distance travelled
		pos.PositionVectorX += ms.CurrentMS * pos.RotationVectorX
		pos.PositionVectorY += ms.CurrentMS * pos.RotationVectorY
		dist.Distance -= ms.CurrentMS

		if dist.Distance <= 0 { //if reached max range
			//update projectiles destroyed component to True
			cardinal.UpdateComponent(world, id, func(destroyed *comp.Destroyed) *comp.Destroyed {
				if destroyed == nil {
					fmt.Printf("error retrieving enemy destroyed component (sp_archerlady.go): ")
					return nil
				}
				destroyed.Destroyed = true
				return destroyed
			})
		} else { // else update distance component
			if err = cardinal.SetComponent(world, id, dist); err != nil {
				fmt.Printf("error setting distance component (sp_archerlady.go): ")
				return nil
			}
		}
		return pos
	})

	return nil
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
func GetSpComponentsAL(world cardinal.WorldContext, unitID types.EntityID) (*comp.Position, *comp.MatchId, *comp.MapName, *comp.Team, error) {
	position, err := cardinal.GetComponent[comp.Position](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error retrieving Position component (sp_archerlay.go): %v", err)
	}

	matchId, err := cardinal.GetComponent[comp.MatchId](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error retrieving MatchId component (sp_archerlay.go): %v", err)
	}

	mapName, err := cardinal.GetComponent[comp.MapName](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error retrieving mapname component (sp_archerlay.go): %v", err)
	}

	team, err := cardinal.GetComponent[comp.Team](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error retrieving team component (sp_archerlay.go): %v", err)
	}

	return position, matchId, mapName, team, nil
}

// GetSpComponentsAL fetches all necessary components related to a sp entity.
func GetUpdateComponentsAL(world cardinal.WorldContext, unitID types.EntityID) (*comp.Movespeed, *comp.Distance, *comp.MatchId, *comp.Team, *comp.Damage, error) {
	ms, err := cardinal.GetComponent[comp.Movespeed](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error retrieving movespeed component (sp_archerlay.go): %v", err)
	}
	dist, err := cardinal.GetComponent[comp.Distance](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error retrieving distance component (sp_archerlay.go): %v", err)
	}
	mID, err := cardinal.GetComponent[comp.MatchId](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error retrieving matchID component (sp_archerlay.go): %v", err)
	}
	team, err := cardinal.GetComponent[comp.Team](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error retrieving team component (sp_archerlay.go): %v", err)
	}
	dmg, err := cardinal.GetComponent[comp.Damage](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error retrieving damage component (sp_archerlay.go): %v", err)
	}
	return ms, dist, mID, team, dmg, nil
}
