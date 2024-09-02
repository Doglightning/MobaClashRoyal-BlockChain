package query

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"

	comp "MobaClashRoyal/component"
	"MobaClashRoyal/system"

	"pkg.world.dev/world-engine/cardinal"
)

type UnitMatchIdRequest struct {
	MatchId string
}

type UnitStateResponse struct {
	Units         []UnitDetails
	Structures    []StructureDetails
	Projectiles   []ProjectileDetails
	SpecialPowers []SpDetails
}

type UnitDetails struct {
	Team            string
	CurrentHP       float32
	MaxHP           float32
	PositionVectorX float32
	PositionVectorY float32
	PositionVectorZ float32
	RotationVectorX float32
	RotationVectorY float32
	RotationVectorZ float32
	Movespeed       float32
	UID             int
	UnitName        string
	Combat          bool
	AttackFrame     int
	CurrentSp       int
	Animation       string
}

type ProjectileDetails struct {
	UID             int
	Name            string
	PositionVectorX float32
	PositionVectorY float32
	PositionVectorZ float32
	RotationVectorX float32
	RotationVectorY float32
	RotationVectorZ float32
}

type SpDetails struct {
	UID             int
	Name            string
	PositionVectorX float32
	PositionVectorY float32
	PositionVectorZ float32
	RotationVectorX float32
	RotationVectorY float32
	RotationVectorZ float32
}

type StructureDetails struct {
	UID             int
	CurrentHP       float32
	StructureName   string
	Team            string
	State           string
	PositionVectorX float32
	PositionVectorY float32
	PositionVectorZ float32
}

func GameState(world cardinal.WorldContext, req *UnitMatchIdRequest) (*UnitStateResponse, error) {
	var response UnitStateResponse

	//filter for entities with matchID
	matchFilter := cardinal.ComponentFilter(func(m comp.MatchId) bool {
		return m.MatchId == req.MatchId
	})

	//get unit state updates
	response, err := unitStateGS(world, matchFilter, response)
	if err != nil {
		return nil, err
	}
	//get structure state updates
	response, err = structureStateGS(world, matchFilter, response)
	if err != nil {
		return nil, err
	}

	//get projectile state updates
	response, err = projectileStateGS(world, matchFilter, response)
	if err != nil {
		return nil, err
	}

	//get special power state updates
	response, err = SpStateGS(world, matchFilter, response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// get all units states
func unitStateGS(world cardinal.WorldContext, matchFilter cardinal.FilterFn, response UnitStateResponse) (UnitStateResponse, error) {
	// get each unit
	unitSearch := cardinal.NewSearch().Entity(
		filter.Contains(system.UnitFilters())).
		Where(matchFilter)

	err := unitSearch.Each(world, func(id types.EntityID) bool {
		unit := UnitDetails{}

		// Fetch Team component
		team, err := cardinal.GetComponent[comp.Team](world, id)
		if err != nil {
			return false
		}
		unit.Team = team.Team

		// Fetch Health component
		health, err := cardinal.GetComponent[comp.Health](world, id)
		if err != nil {
			return false
		}
		unit.CurrentHP = health.CurrentHP
		unit.MaxHP = health.MaxHP

		// Fetch Position component
		position, err := cardinal.GetComponent[comp.Position](world, id)
		if err != nil {
			return false
		}
		unit.PositionVectorX = position.PositionVectorX
		unit.PositionVectorY = position.PositionVectorY
		unit.PositionVectorZ = position.PositionVectorZ

		unit.RotationVectorX = position.RotationVectorX
		unit.RotationVectorY = position.RotationVectorY
		unit.RotationVectorZ = position.RotationVectorZ

		// Fetch Movespeed component
		movespeed, err := cardinal.GetComponent[comp.Movespeed](world, id)
		if err != nil {
			return false
		}
		unit.Movespeed = movespeed.CurrentMS

		// Fetch UID component
		uid, err := cardinal.GetComponent[comp.UID](world, id) // Assuming UID is a component; adjust as needed
		if err != nil {
			return false
		}
		unit.UID = uid.UID

		// Fetch UnitName component
		unitName, err := cardinal.GetComponent[comp.UnitName](world, id)
		if err != nil {
			return false
		}
		unit.UnitName = unitName.UnitName

		// Fetch UnitAttack component
		unitAttack, err := cardinal.GetComponent[comp.Attack](world, id)
		if err != nil {
			return false
		}
		unit.Combat = unitAttack.Combat
		unit.AttackFrame = unitAttack.Frame

		// Fetch SP component
		unitSp, err := cardinal.GetComponent[comp.Sp](world, id)
		if err != nil {
			return false
		}
		unit.CurrentSp = unitSp.CurrentSp
		unit.Animation = unitSp.Animation

		// Append the gathered data to the response
		response.Units = append(response.Units, unit)
		return true // Continue iteration
	})
	if err != nil {
		return response, fmt.Errorf("error during unit search (unit state/unitStateGS): %w", err)
	}

	return response, nil
}

// get all projectile states
func structureStateGS(world cardinal.WorldContext, matchFilter cardinal.FilterFn, response UnitStateResponse) (UnitStateResponse, error) {
	// get the structure id's
	StructureSearch := cardinal.NewSearch().Entity(
		filter.Exact(system.StructureFilters())).
		Where(matchFilter)

	err := StructureSearch.Each(world, func(id types.EntityID) bool {
		structure := StructureDetails{}

		// Fetch Team component
		team, err := cardinal.GetComponent[comp.Team](world, id)
		if err != nil {
			return false
		}
		structure.Team = team.Team

		// Fetch uid component
		uid, err := cardinal.GetComponent[comp.UID](world, id)
		if err != nil {
			return false
		}
		structure.UID = uid.UID

		// Fetch health component
		health, err := cardinal.GetComponent[comp.Health](world, id)
		if err != nil {
			return false
		}
		structure.CurrentHP = health.CurrentHP

		// Fetch Team component
		name, err := cardinal.GetComponent[comp.UnitName](world, id)
		if err != nil {
			return false
		}
		structure.StructureName = name.UnitName

		// Fetch State component
		state, err := cardinal.GetComponent[comp.State](world, id)
		if err != nil {
			return false
		}
		structure.State = state.State

		// Fetch Position component
		position, err := cardinal.GetComponent[comp.Position](world, id)
		if err != nil {
			return false
		}
		structure.PositionVectorX = position.PositionVectorX
		structure.PositionVectorY = position.PositionVectorY
		structure.PositionVectorZ = position.PositionVectorZ

		response.Structures = append(response.Structures, structure)
		return true
	})
	if err != nil {
		return response, fmt.Errorf("error during StructureSearch (unit state/structureStateGS): %w", err)
	}
	return response, nil
}

// get all structure states
func projectileStateGS(world cardinal.WorldContext, matchFilter cardinal.FilterFn, response UnitStateResponse) (UnitStateResponse, error) {
	// get the projectile id's
	ProjectileSearch := cardinal.NewSearch().Entity(
		filter.Exact(system.ProjectileFilters())).
		Where(matchFilter)

	err := ProjectileSearch.Each(world, func(id types.EntityID) bool {
		projectile := ProjectileDetails{}

		// Fetch uid component
		uid, err := cardinal.GetComponent[comp.UID](world, id)
		if err != nil {
			return false
		}
		projectile.UID = uid.UID

		// Fetch name component
		name, err := cardinal.GetComponent[comp.UnitName](world, id)
		if err != nil {
			return false
		}
		projectile.Name = name.UnitName

		// Fetch Position component
		position, err := cardinal.GetComponent[comp.Position](world, id)
		if err != nil {
			return false
		}
		projectile.PositionVectorX = position.PositionVectorX
		projectile.PositionVectorY = position.PositionVectorY
		projectile.PositionVectorZ = position.PositionVectorZ

		projectile.RotationVectorX = position.RotationVectorX
		projectile.RotationVectorY = position.RotationVectorY
		projectile.RotationVectorZ = position.RotationVectorZ

		response.Projectiles = append(response.Projectiles, projectile)
		return true
	})
	if err != nil {
		return response, fmt.Errorf("error during StructureSearch (unit state/projectileStateGS): %w", err)
	}
	return response, nil
}

// get all Special power states
func SpStateGS(world cardinal.WorldContext, matchFilter cardinal.FilterFn, response UnitStateResponse) (UnitStateResponse, error) {
	// get the projectile id's
	SpSearch := cardinal.NewSearch().Entity(
		filter.Contains(filter.Component[comp.SpName]())).
		Where(matchFilter)

	err := SpSearch.Each(world, func(id types.EntityID) bool {
		sp := SpDetails{}

		// Fetch uid component
		uid, err := cardinal.GetComponent[comp.UID](world, id)
		if err != nil {
			return false
		}
		sp.UID = uid.UID

		// Fetch name component
		name, err := cardinal.GetComponent[comp.SpName](world, id)
		if err != nil {
			return false
		}
		sp.Name = name.SpName

		// Fetch Position component
		position, err := cardinal.GetComponent[comp.Position](world, id)
		if err != nil {
			return false
		}
		sp.PositionVectorX = position.PositionVectorX
		sp.PositionVectorY = position.PositionVectorY
		sp.PositionVectorZ = position.PositionVectorZ

		sp.RotationVectorX = position.RotationVectorX
		sp.RotationVectorY = position.RotationVectorY
		sp.RotationVectorZ = position.RotationVectorZ

		response.SpecialPowers = append(response.SpecialPowers, sp)
		return true
	})
	if err != nil {
		return response, fmt.Errorf("error during  spSearch (unit state/SpStateGS): %w", err)
	}
	return response, nil
}
