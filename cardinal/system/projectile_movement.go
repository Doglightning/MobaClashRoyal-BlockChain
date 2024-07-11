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

		//get enemy Position compoenent
		enemyPosition, err := cardinal.GetComponent[comp.Position](world, projectileAttack.Target)
		if err != nil {
			fmt.Printf("error retrieving enemy Position component (projectile movement): %s", err)
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

		// // Ensure the unit does not overshoot the target position
		// if (directionVectorX > 0 && newPosX > enemyPosition.PositionVectorX) || (directionVectorX < 0 && newPosX < enemyPosition.PositionVectorX) {
		// 	newPosX = enemyPosition.PositionVectorX
		// }
		// if (directionVectorY > 0 && newPosY > enemyPosition.PositionVectorY) || (directionVectorY < 0 && newPosY < enemyPosition.PositionVectorY) {
		// 	newPosY = enemyPosition.PositionVectorY
		// }

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

		return true
	})

	if err != nil {

		return fmt.Errorf("error retrieving projectile entities (unit destroyer): %s", err)
	}

	return nil
}
