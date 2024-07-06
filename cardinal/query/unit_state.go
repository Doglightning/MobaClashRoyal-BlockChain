package query

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"

	comp "MobaClashRoyal/component"
	"MobaClashRoyal/system"

	"pkg.world.dev/world-engine/cardinal"
)

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
}

type UnitMatchIdRequest struct {
	MatchId string
}

type UnitStateResponse struct {
	Units []UnitDetails
}

func UnitState(world cardinal.WorldContext, req *UnitMatchIdRequest) (*UnitStateResponse, error) {
	matchFilter := cardinal.ComponentFilter[comp.MatchId](func(m comp.MatchId) bool {
		return m.MatchId == req.MatchId
	})

	unitSearch := cardinal.NewSearch().Entity(
		filter.Exact(system.UnitFilters())).
		Where(matchFilter)

	var response UnitStateResponse
	err := unitSearch.Each(world, func(id types.EntityID) bool {
		unit := UnitDetails{}

		// Fetch Team component
		team, err := cardinal.GetComponent[comp.Team](world, id)
		if err != nil {
			return false
		}
		unit.Team = team.Team

		// Fetch Health component
		health, err := cardinal.GetComponent[comp.UnitHealth](world, id)
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
			return false // Stop iteration on error
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

		// Append the gathered data to the response
		response.Units = append(response.Units, unit)
		return true // Continue iteration
	})
	if err != nil {
		return nil, fmt.Errorf("error during search: %w", err)
	}

	return &response, nil
}
