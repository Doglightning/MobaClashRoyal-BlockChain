package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

// moves projectiles towards target
func ProjectileMovementSystem(world cardinal.WorldContext) error {
	//filter for class type projectile
	classFilter := cardinal.ComponentFilter(func(m comp.Class) bool {
		return m.Class == "projectile"
	})
	//for each projectile id
	err := cardinal.NewSearch().Entity(
		filter.Contains(filter.Component[comp.ProjectileTag]())).
		Where(classFilter).Each(world, func(projectileID types.EntityID) bool {
		//get needed projectile components
		projectileAtk, projectileMs, projectilePos, err := GetComponents3[comp.Attack, comp.Movespeed, comp.Position](world, projectileID)
		if err != nil {
			fmt.Printf("projectile components (Projectile_movement.go): %v \n", err)
			return false
		}
		//copy starting position
		oldPos := &comp.Position{
			PositionVectorX: projectilePos.PositionVectorX,
			PositionVectorY: projectilePos.PositionVectorY,
			PositionVectorZ: projectilePos.PositionVectorZ,
		}

		//get projectiles targets position component
		enemyPos, eCenOffset, err := GetComponents2[comp.Position, comp.CenterOffset](world, projectileAtk.Target)
		if err != nil {
			fmt.Printf("projectile target components (Projectile_movement.go): %v \n", err)
			return false
		}

		// // Compute direction vector towards the enemy
		// deltaX := enemyPos.PositionVectorX - projectilePos.PositionVectorX
		// deltaY := enemyPos.PositionVectorY - projectilePos.PositionVectorY
		// deltaZ := enemyPos.PositionVectorZ + eCenOffset.CenterOffset - projectilePos.PositionVectorZ
		// magnitude := float32(math.Sqrt(float64(deltaX*deltaX + deltaY*deltaY + deltaZ*deltaZ)))

		// // Normalize the direction vector
		// projectilePos.RotationVectorX = deltaX / magnitude
		// projectilePos.RotationVectorY = deltaY / magnitude
		// projectilePos.RotationVectorZ = deltaZ / magnitude

		projectilePos.RotationVectorX, projectilePos.RotationVectorY, projectilePos.RotationVectorZ = directionVectorBetweenTwoPoints3D(projectilePos.PositionVectorX, projectilePos.PositionVectorY, projectilePos.PositionVectorZ, enemyPos.PositionVectorX, enemyPos.PositionVectorY, enemyPos.PositionVectorZ+eCenOffset.CenterOffset)

		// Compute new position based on movespeed and direction
		projectilePos.PositionVectorX += projectilePos.RotationVectorX * projectileMs.CurrentMS
		projectilePos.PositionVectorY += projectilePos.RotationVectorY * projectileMs.CurrentMS
		projectilePos.PositionVectorZ += projectilePos.RotationVectorZ * projectileMs.CurrentMS

		// Set the new position component
		err = cardinal.SetComponent(world, projectileID, projectilePos)
		if err != nil {
			fmt.Printf("error set posisiton compoenent on projectileID (Projectile_movement.go): %v \n", err)
			return false
		}

		//if projectile has passed enemy
		if hasPassedEnemyPM(oldPos, projectilePos, enemyPos) {
			//update attack component in combat
			projectileAtk.Combat = true
			err = cardinal.SetComponent(world, projectileID, projectileAtk)
			if err != nil {
				fmt.Printf("error set attack compoenent on projectileID (Projectile_movement.go): %v \n", err)
				return false
			}

		}

		return true
	})

	if err != nil {

		return fmt.Errorf("error retrieving projectile entities (Projectile_movement.go): %s ", err)
	}

	return nil
}

// hasPassedEnemy checks if the projectile has passed directly over the enemy's position.
func hasPassedEnemyPM(oldPos *comp.Position, newPos *comp.Position, enemyPos *comp.Position) bool {
	// Vectors from old and new positions to the enemy's position
	oldToEnemy := comp.Position{PositionVectorX: enemyPos.PositionVectorX - oldPos.PositionVectorX, PositionVectorY: enemyPos.PositionVectorY - oldPos.PositionVectorY}
	newToEnemy := comp.Position{PositionVectorX: enemyPos.PositionVectorX - newPos.PositionVectorX, PositionVectorY: enemyPos.PositionVectorY - newPos.PositionVectorY}

	// Dot product of vectors from old and new positions to enemy's position
	dotProduct := dotProduct(oldToEnemy.PositionVectorX, oldToEnemy.PositionVectorY, newToEnemy.PositionVectorX, newToEnemy.PositionVectorY)

	// If the dot product is negative, the direction relative to the enemy has changed, meaning the projectile has passed the enemy
	return dotProduct < 0
}
