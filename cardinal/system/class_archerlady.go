package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

// archerLadySP struct contains configuration for an archer lady in terms of her shooting properties.
type archerLadySpawnSP struct {
	Name                  string
	NumArrows             int
	ArrowSeparationDegree float64 //distance between each arrow
	Distance              float32
	Speed                 float32
	offSet                float32
	RadiusArrows          int
	Damage                int
}

// NewArcherLadySP creates a new instance of archerLadySP with default settings.
func NewArcherLadySpawnSP() *archerLadySpawnSP {
	return &archerLadySpawnSP{
		Name:                  "ArcherLadySP",
		NumArrows:             6,
		ArrowSeparationDegree: 20,
		Distance:              1600,
		Speed:                 150,
		offSet:                190,
		RadiusArrows:          150,
		Damage:                20,
	}
}

// update struct
type archerLadyUpdateSP struct {
	BaseDmgReductionFactor int
}

// update vars
func NewArcherLadyUpdateSP() *archerLadyUpdateSP {
	return &archerLadyUpdateSP{
		BaseDmgReductionFactor: 3,
	}
}

// spawns the archer ladies volly special power
func archerLadySpawn(world cardinal.WorldContext, id types.EntityID) error {
	archerLady := NewArcherLadySpawnSP()

	//get needed components
	uPos, matchID, mapName, team, err := GetComponents4[comp.Position, comp.MatchId, comp.MapName, comp.Team](world, id)
	if err != nil {
		return fmt.Errorf("get components (class archerladySpawn): %v", err)
	}
	//generate the vectors the volly will follow
	vectors := generateVectors(uPos.RotationVectorX, uPos.RotationVectorY, archerLady.ArrowSeparationDegree, archerLady.NumArrows)

	for i := 0; i < archerLady.NumArrows; i++ { //for each arrow
		//get next uid
		UID, err := getNextUID(world, matchID.MatchId)
		if err != nil {
			return fmt.Errorf("(class archerladySpawn): %v - ", err)
		}
		//create projectile entity
		cardinal.Create(world,
			comp.MatchId{MatchId: matchID.MatchId},
			comp.UID{UID: UID},
			comp.SpName{SpName: archerLady.Name},
			comp.Movespeed{CurrentMS: archerLady.Speed},
			comp.Position{
				PositionVectorX: uPos.PositionVectorX,
				PositionVectorY: uPos.PositionVectorY,
				PositionVectorZ: uPos.PositionVectorZ + archerLady.offSet,
				RotationVectorX: vectors[i][0],
				RotationVectorY: vectors[i][1],
				RotationVectorZ: uPos.RotationVectorZ,
			},
			comp.MapName{MapName: mapName.MapName},
			comp.Damage{Damage: archerLady.Damage},
			comp.Destroyed{Destroyed: false},
			comp.Distance{Distance: archerLady.Distance},
			comp.Team{Team: team.Team},
			comp.UnitRadius{UnitRadius: archerLady.RadiusArrows},
			comp.SpEntity{SpName: archerLady.Name},
			comp.Class{Class: "sp"},
		)
	}

	return err
}

// called every tick to updated the archerladies arrows
func archerLadyUpdate(world cardinal.WorldContext, id types.EntityID) error {
	//update Sp location
	err := cardinal.UpdateComponent(world, id, func(pos *comp.Position) *comp.Position {
		if pos == nil {
			fmt.Printf("error retrieving enemy position component (class archerladyUpdate): \n")
			return nil
		}
		//get components
		ms, dist, mID, team, dmg, radi, err := GetComponents6[comp.Movespeed, comp.Distance, comp.MatchId, comp.Team, comp.Damage, comp.UnitRadius](world, id)
		if err != nil {
			fmt.Printf("SP Components (class archerladyUpdate): %v", err)
			return nil
		}
		//get collision hash
		collisionHash, err := getCollisionHashGSS(world, mID)
		if err != nil {
			fmt.Printf("(class archerladyUpdate): - \n")
			return nil
		}

		//updated position end point
		endX := pos.PositionVectorX + ms.CurrentMS*pos.RotationVectorX
		endY := pos.PositionVectorY + ms.CurrentMS*pos.RotationVectorY
		//find close units that arrow could have possibly crossed
		cID := CheckCollisionSpatialHashList(collisionHash, pos.PositionVectorX, pos.PositionVectorY, radi.UnitRadius, "range", false)

		var closestUnit types.EntityID //keeping track of closest unit over for loop
		var closestDistance float32 = -1
		for _, value := range cID { //for each collision
			//get collision team component
			cTeam, err := cardinal.GetComponent[comp.Team](world, value)
			if err != nil {
				fmt.Printf("(error getting team component (class archerladyUpdate):: %v \n", err)
				continue
			}

			if team.Team != cTeam.Team { // if different teams
				//get colision position and radius components
				cPos, cRad, err := GetComponents2[comp.Position, comp.UnitRadius](world, value)
				if err != nil {
					fmt.Printf("collision compoenents (class archerladyUpdate): -  %v \n", err)
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
					fmt.Printf("error retrieving enemy destroyed component (class archerladyUpdate):: \n")
					return nil
				}
				destroyed.Destroyed = true
				return destroyed
			})
			//reduce enemy current health
			cardinal.UpdateComponent(world, closestUnit, func(health *comp.Health) *comp.Health {
				if health == nil {
					fmt.Printf("error retrieving collision health component ((class archerladyUpdate): \n")
					return nil
				}
				///get target name
				targetName, err := cardinal.GetComponent[comp.UnitName](world, closestUnit)
				if err != nil {
					fmt.Printf("error retrieving target unit name component (class archerladyUpdate): \n")
					return nil
				}

				if targetName.UnitName == "Base" || targetName.UnitName == "Tower" { // reduce damage to structures
					archerLady := NewArcherLadyUpdateSP() // get reduction var
					health.CurrentHP -= float32(dmg.Damage / archerLady.BaseDmgReductionFactor)
				} else {
					health.CurrentHP -= float32(dmg.Damage) // damage to non towers
				}

				if health.CurrentHP < 0 { //dont get health get into neg
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
					fmt.Printf("error retrieving enemy destroyed component (class archerladyUpdate): \n")
					return nil
				}
				destroyed.Destroyed = true
				return destroyed
			})
		} else { // else update distance component
			if err = cardinal.SetComponent(world, id, dist); err != nil {
				fmt.Printf("error setting distance component (class archerladyUpdate): \n")
				return nil
			}
		}
		return pos
	})
	return err
}

// spawns projectile for archer basic attack
func archerLadyAttack(world cardinal.WorldContext, id types.EntityID, atk *comp.Attack) error {
	//get units component
	unitPosition, matchID, mapName, unitName, err := GetComponents4[comp.Position, comp.MatchId, comp.MapName, comp.UnitName](world, id)
	if err != nil {
		return fmt.Errorf("unit components (class archerladyAttack.go): %v ", err)
	}
	//get next uid
	UID, err := getNextUID(world, matchID.MatchId)
	if err != nil {
		return fmt.Errorf("(class archerladyAttack.go): %v ", err)
	}

	//off set units arrow spawn location to match model on client
	newX, newY := RelativeOffsetXY(unitPosition.PositionVectorX, unitPosition.PositionVectorY, unitPosition.RotationVectorX, unitPosition.RotationVectorY, ProjectileRegistry[unitName.UnitName].offSetX, ProjectileRegistry[unitName.UnitName].offSetY)
	//create projectile entity
	_, err = cardinal.Create(world,
		comp.MatchId{MatchId: matchID.MatchId},
		comp.UID{UID: UID},
		comp.UnitName{UnitName: ProjectileRegistry[unitName.UnitName].Name},
		comp.Movespeed{CurrentMS: ProjectileRegistry[unitName.UnitName].Speed},
		comp.Position{
			PositionVectorX: newX,
			PositionVectorY: newY,
			PositionVectorZ: unitPosition.PositionVectorZ + ProjectileRegistry[unitName.UnitName].offSetZ,
			RotationVectorX: unitPosition.RotationVectorX,
			RotationVectorY: unitPosition.RotationVectorY,
			RotationVectorZ: unitPosition.RotationVectorZ,
		},
		comp.MapName{MapName: mapName.MapName},
		comp.Class{Class: "projectile"},
		comp.Attack{Target: atk.Target, Damage: UnitRegistry[unitName.UnitName].Damage},
		comp.Destroyed{Destroyed: false},
		comp.ProjectileTag{},
	)

	if err != nil {
		return fmt.Errorf("error spawning archer lady basic attack (class archerladyAttack.go): %v ", err)
	}

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
