package system

import (
	comp "MobaClashRoyal/component"
	"fmt"
	"math"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

func ProjectileMovementSystem(world cardinal.WorldContext) error {
	targetFilter := cardinal.ComponentFilter[comp.Attack](func(m comp.Attack) bool {
		return m.Class == "projectile"
	})

	err := cardinal.NewSearch().Entity(
		filter.Exact(ProjectileFilters())).
		Where(targetFilter).Each(world, func(projectileID types.EntityID) bool {

		//get projectile attack compoenent
		projectileAttack, err := cardinal.GetComponent[comp.Attack](world, projectileID)
		if err != nil {
			fmt.Printf("error retrieving projectileattack component (projectile movement): %s", err)
			return false
		}

		//get movespeed compoenent
		projectileMovespeed, err := cardinal.GetComponent[comp.Movespeed](world, projectileID)
		if err != nil {
			fmt.Printf("error retrieving movespeed component (projectile movement): %s", err)
			return false
		}

		//get projectile Position compoenent
		projectilePosition, err := cardinal.GetComponent[comp.Position](world, projectileID)
		if err != nil {
			fmt.Printf("error retrieving enemy Position component (projectile movement): %s", err)
			return false
		}

		oldPos := &comp.Position{PositionVectorX: projectilePosition.PositionVectorX, PositionVectorY: projectilePosition.PositionVectorY}

		enemyPosition, _, err := getEnemyComponentsUM(world, projectileAttack.Target)
		if err != nil {
			fmt.Printf("(projectile Movement): %s\n", err)
			return false
		}

		// Compute direction vector towards the enemy
		deltaX := float64(enemyPosition.PositionVectorX - projectilePosition.PositionVectorX)
		deltaY := float64(enemyPosition.PositionVectorY - projectilePosition.PositionVectorY)
		magnitude := math.Sqrt(deltaX*deltaX + deltaY*deltaY)

		// Normalize the direction vector
		directionVectorX := float32(deltaX / magnitude)
		directionVectorY := float32(deltaY / magnitude)

		// Compute new position based on movespeed and direction, but do not exceed the target position
		newPosX := projectilePosition.PositionVectorX + directionVectorX*projectileMovespeed.CurrentMS
		newPosY := projectilePosition.PositionVectorY + directionVectorY*projectileMovespeed.CurrentMS

		projectilePosition.PositionVectorX = newPosX
		projectilePosition.PositionVectorY = newPosY
		projectilePosition.RotationVectorX = directionVectorX
		projectilePosition.RotationVectorY = directionVectorY

		// Set the new position component
		err = cardinal.SetComponent(world, projectileID, projectilePosition)
		if err != nil {
			fmt.Printf("error set posisiton compoenent on projectileID (projectile movement): %v", err)
			return false
		}

		//has projectile entered enemy radius?
		if hasPassedEnemy(oldPos, projectilePosition, enemyPosition) {
			//update attack component in combat
			projectileAttack.Combat = true
			err = cardinal.SetComponent(world, projectileID, projectileAttack)
			if err != nil {
				fmt.Printf("error set attack compoenent on projectileID (projectile movement): %v", err)
				return false
			}

		}

		return true
	})

	if err != nil {

		return fmt.Errorf("error retrieving projectile entities (unit destroyer): %s", err)
	}

	return nil
}

// hasPassedEnemy checks if the projectile has passed directly over the enemy's position.
func hasPassedEnemy(oldPos *comp.Position, newPos *comp.Position, enemyPos *comp.Position) bool {
	// Vectors from old and new positions to the enemy's position
	oldToEnemy := comp.Position{PositionVectorX: enemyPos.PositionVectorX - oldPos.PositionVectorX, PositionVectorY: enemyPos.PositionVectorY - oldPos.PositionVectorY}
	newToEnemy := comp.Position{PositionVectorX: enemyPos.PositionVectorX - newPos.PositionVectorX, PositionVectorY: enemyPos.PositionVectorY - newPos.PositionVectorY}

	// Dot product of vectors from old and new positions to enemy's position
	oldDot := float64(oldToEnemy.PositionVectorX)*float64(newToEnemy.PositionVectorX) + float64(oldToEnemy.PositionVectorY)*float64(newToEnemy.PositionVectorY)

	// If the dot product is negative, the direction relative to the enemy has changed, meaning the projectile has passed the enemy
	return oldDot < 0
}
