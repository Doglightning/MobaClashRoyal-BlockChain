package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

// archerLadySP struct contains configuration for an archer lady in terms of her shooting properties.
type archerLadySP struct {
	Name                  string
	NumArrows             int
	ArrowSeparationDegree float64
	Distance              float32
	Speed                 float32
	RadiusArrows          int
	Damage                int
}

// NewArcherLadySP creates a new instance of archerLadySP with default settings.
func NewArcherLadySP() *archerLadySP {
	return &archerLadySP{
		Name:                  "ArcherLadySP",
		NumArrows:             6,
		ArrowSeparationDegree: 20,
		Distance:              1600,
		Speed:                 150,
		RadiusArrows:          150,
		Damage:                30,
	}
}

// spawns the archer ladies villy special power
func archerLadySpawn(world cardinal.WorldContext, id types.EntityID) error {
	archerLady := NewArcherLadySP()

	//get needed components
	uPos, matchID, mapName, team, err := GetSpComponentsAL(world, id)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	//generate the vectors the volly will follow
	vectors := generateVectors(uPos.RotationVectorX, uPos.RotationVectorY, archerLady.ArrowSeparationDegree, archerLady.NumArrows)

	for i := 0; i < archerLady.NumArrows; i++ { //for each arrow
		//get next uid
		UID, err := getNextUID(world, matchID.MatchId)
		if err != nil {
			return fmt.Errorf("(archerLadyAttack): %v - ", err)
		}
		//create projectile entity
		cardinal.Create(world,
			comp.MatchId{MatchId: matchID.MatchId},
			comp.UID{UID: UID},
			comp.SpName{SpName: archerLady.Name},
			comp.Movespeed{CurrentMS: archerLady.Speed},
			comp.Position{PositionVectorX: uPos.PositionVectorX, PositionVectorY: uPos.PositionVectorY, PositionVectorZ: uPos.PositionVectorZ, RotationVectorX: vectors[i][0], RotationVectorY: vectors[i][1], RotationVectorZ: uPos.RotationVectorZ},
			comp.MapName{MapName: mapName.MapName},
			comp.Damage{Damage: archerLady.Damage},
			comp.Destroyed{Destroyed: false},
			comp.Distance{Distance: archerLady.Distance},
			comp.Team{Team: team.Team},
			comp.UnitRadius{UnitRadius: archerLady.RadiusArrows},
			comp.SpEntity{SpName: archerLady.Name},
		)
	}

	return err
}

// called every tick to updated the archerladies arrows
func archerLadyUpdate(world cardinal.WorldContext, id types.EntityID) error {
	//update Sp location
	err := cardinal.UpdateComponent(world, id, func(pos *comp.Position) *comp.Position {
		if pos == nil {
			fmt.Printf("error retrieving enemy position component (sp_archerlady.go): ")
			return nil
		}
		//get components
		ms, dist, mID, team, dmg, radi, err := GetUpdateComponentsAL(world, id)
		if err != nil {
			fmt.Printf("%v", err)
			return nil
		}
		//get collision hash
		collisionHash, err := getCollisionHashGSS(world, mID)
		if err != nil {
			fmt.Printf("(sp_archerlady.go) - ")
			return nil
		}

		//updated position end point
		endX := pos.PositionVectorX + ms.CurrentMS*pos.RotationVectorX
		endY := pos.PositionVectorY + ms.CurrentMS*pos.RotationVectorY
		//find close units that arrow could have possibly crossed
		cID := CheckCollisionSpatialHashList(collisionHash, pos.PositionVectorX, pos.PositionVectorY, radi.UnitRadius)

		var closestUnit types.EntityID
		var closestDistance float32 = -1
		for _, value := range cID { //for each collision
			//get collision team component
			cTeam, err := cardinal.GetComponent[comp.Team](world, value)
			if err != nil {
				fmt.Printf("(error getting team component (sp_archerlady.go): %v", err)
				continue
			}

			if team.Team != cTeam.Team { // if different teams
				//get colision position and radius components
				cPos, cRad, err := getTargetComponentsUM(world, value)
				if err != nil {
					fmt.Printf("(sp_archerlady.go) -  %v", err)
					continue
				}
				//check if passed over a enemy
				if checkLineIntersectionSpatialHash(pos.PositionVectorX, pos.PositionVectorY, endX, endY, cPos.PositionVectorX, cPos.PositionVectorY, cRad.UnitRadius) {
					if closestDistance == -1 { //first unit found
						closestUnit = value
						closestDistance = distanceBetweenTwoPoints(pos.PositionVectorX, pos.PositionVectorY, endX, endY)
					}
					tempDistance := distanceBetweenTwoPoints(pos.PositionVectorX, pos.PositionVectorY, endX, endY)
					if tempDistance < closestDistance { // if distance is less then stored closest
						closestUnit = value
						closestDistance = tempDistance
					}

				}
			}
		}

		if closestDistance != -1 { // if units were found
			//destroy arrow
			cardinal.UpdateComponent(world, id, func(destroyed *comp.Destroyed) *comp.Destroyed {
				if destroyed == nil {
					fmt.Printf("error retrieving enemy destroyed component (sp_archerlady.go): ")
					return nil
				}
				destroyed.Destroyed = true
				return destroyed
			})
			//reduce enemy current health
			cardinal.UpdateComponent(world, closestUnit, func(health *comp.Health) *comp.Health {
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

		//updated position and distance travelled
		pos.PositionVectorX = endX
		pos.PositionVectorY = endY
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
func GetUpdateComponentsAL(world cardinal.WorldContext, unitID types.EntityID) (*comp.Movespeed, *comp.Distance, *comp.MatchId, *comp.Team, *comp.Damage, *comp.UnitRadius, error) {
	ms, err := cardinal.GetComponent[comp.Movespeed](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving movespeed component (sp_archerlay.go): %v", err)
	}
	dist, err := cardinal.GetComponent[comp.Distance](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving distance component (sp_archerlay.go): %v", err)
	}
	mID, err := cardinal.GetComponent[comp.MatchId](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving matchID component (sp_archerlay.go): %v", err)
	}
	team, err := cardinal.GetComponent[comp.Team](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving team component (sp_archerlay.go): %v", err)
	}
	dmg, err := cardinal.GetComponent[comp.Damage](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving damage component (sp_archerlay.go): %v", err)
	}
	rad, err := cardinal.GetComponent[comp.UnitRadius](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving radius component (sp_archerlay.go): %v", err)
	}
	return ms, dist, mID, team, dmg, rad, nil
}
