package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

// fireSpiritSpawnSP struct contains configuration for an fire spirit in terms of their shooting properties.
type fireSpiritSpawnSP struct {
	Hieght    float32
	BaseWidth float32
	Damage    float32
}

// NewArcherLadySP creates a new instance of archerLadySP with default settings.
func NewFireSpiritSpawnSP() *fireSpiritSpawnSP {
	return &fireSpiritSpawnSP{
		Hieght:    570,
		BaseWidth: 385,
		Damage:    2.5,
	}
}

func fireSpiritSpawn(world cardinal.WorldContext, id types.EntityID) error {
	//get fire spirit vars
	fireSprit := NewFireSpiritSpawnSP()

	//get team comp
	team, err := cardinal.GetComponent[comp.Team](world, id)
	if err != nil {
		return fmt.Errorf("error getting team component (class fireSpirit.go): %v", err)
	}

	matchID, err := cardinal.GetComponent[comp.MatchId](world, id)
	if err != nil {
		return fmt.Errorf("error getting position component (class fireSpirit.go): %v", err)
	}

	gameState, err := getGameStateGSS(world, matchID)
	if err != nil {
		return fmt.Errorf("(class fireSpirit.go): %v", err)
	}

	//get position comp
	hash, err := cardinal.GetComponent[comp.SpatialHash](world, gameState)
	if err != nil {
		return fmt.Errorf("error getting spatial hash compoenent(class fireSpirit.go): %v", err)
	}

	//get position comp
	pos, err := cardinal.GetComponent[comp.Position](world, id)
	if err != nil {
		return fmt.Errorf("error getting position component (class fireSpirit.go): %v", err)
	}

	//find the 3 points of the fire spirit AoE triangle attack
	apex, baseLeft, baseRight := CreateIsoscelesTriangle(Point{X: pos.PositionVectorX, Y: pos.PositionVectorY}, Point{X: pos.RotationVectorX, Y: pos.RotationVectorY}, fireSprit.Hieght, fireSprit.BaseWidth)

	//list of every point within the triangle
	points := RasterizeIsoscelesTriangle(apex, baseLeft, baseRight)

	// Define a map to track unique collisions
	collidedEntities := make(map[types.EntityID]bool)

	for _, p := range points {

		collList := CheckCollisionSpatialHashList(hash, p.X, p.Y, 1)
		for _, collID := range collList {
			collidedEntities[collID] = true
		}
	}

	// Iterate over each key in the map
	for collID := range collidedEntities {
		//get targets team
		targetTeam, err := cardinal.GetComponent[comp.Team](world, collID)
		if err != nil {
			fmt.Printf("error getting targets team compoenent (class fireSpirit.go): %v \n", err)
			continue
		}

		fmt.Printf("%s \n", team.Team)
		fmt.Printf("%s \n", targetTeam.Team)
		if team.Team != targetTeam.Team {

			// reduce health by units attack damage
			err = cardinal.UpdateComponent(world, collID, func(health *comp.Health) *comp.Health {
				if health == nil {
					fmt.Printf("error retrieving Health component (class fireSpirit.go) \n")
					return nil
				}
				health.CurrentHP -= fireSprit.Damage
				if health.CurrentHP < 0 {
					health.CurrentHP = 0 //never have negative health
				}
				return health
			})
			if err != nil {
				fmt.Printf("error updating health (class fireSpirit.go): %v \n", err)
				continue
			}

		}
	}

	return nil
}
