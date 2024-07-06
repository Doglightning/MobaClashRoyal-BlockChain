package query

import (
	comp "MobaClashRoyal/component"
	"MobaClashRoyal/system"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
)

type RemovalMatchIdRequest struct {
	MatchId string
	Team    string
}

type RemovalStateResponse struct {
	Units []int
}

func RemovalState(world cardinal.WorldContext, req *RemovalMatchIdRequest) (*RemovalStateResponse, error) {
	var response RemovalStateResponse
	var removeList = []int{}

	teamFilter := cardinal.ComponentFilter[comp.MatchId](func(m comp.MatchId) bool {
		return m.MatchId == req.MatchId
	})

	teamSearch, err := cardinal.NewSearch().Entity(
		filter.Exact(system.TeamFilters())).
		Where(teamFilter).First(world)

	if err != nil {
		return nil, fmt.Errorf("error searching for team (Removal State Query): %w", err)
	}

	if req.Team == "Blue" {
		// Get Player1 component
		player1, err := cardinal.GetComponent[comp.Player1](world, teamSearch)
		if err != nil {
			return nil, fmt.Errorf("error retrieving Player1 component (Removal State Query): %w", err)
		}

		for key := range player1.RemovalList {
			removeList = append(removeList, key)
		}

	} else {
		// Get Player1 component
		player2, err := cardinal.GetComponent[comp.Player2](world, teamSearch)
		if err != nil {
			return nil, fmt.Errorf("error retrieving Player2 component (Removal State Query): %w", err)
		}

		for key := range player2.RemovalList {
			removeList = append(removeList, key)
		}
	}
	response.Units = removeList
	return &response, nil
}
