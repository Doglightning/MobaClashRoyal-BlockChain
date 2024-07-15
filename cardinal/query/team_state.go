package query

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal/iterators"
	"pkg.world.dev/world-engine/cardinal/search/filter"

	comp "MobaClashRoyal/component"
	"MobaClashRoyal/system"

	"pkg.world.dev/world-engine/cardinal"
)

type MatchIdRequest struct {
	MatchId string
}

type TeamStateResponse struct {
	Player1 string
	Player2 string
}

func TeamState(world cardinal.WorldContext, req *MatchIdRequest) (*TeamStateResponse, error) {
	matchFilter := cardinal.ComponentFilter[comp.MatchId](func(m comp.MatchId) bool {
		return m.MatchId == req.MatchId
	})

	var response TeamStateResponse
	gameSearch := cardinal.NewSearch().Entity(
		filter.Exact(system.GameStateFilters())).
		Where(matchFilter)

	gameState, err := gameSearch.First(world)

	if err != nil {
		return nil, fmt.Errorf("error searching for match: %w", err)
	}

	if gameState == iterators.BadID { // Assuming cardinal.NoEntity represents no result found
		return nil, fmt.Errorf("no match found with ID or missing components: %s", req.MatchId)
	}

	// Get Player1 component
	player1, err := cardinal.GetComponent[comp.Player1](world, gameState)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Player1 component: %w", err)
	}
	response.Player1 = player1.Nickname

	// Get Player2 component
	player2, err := cardinal.GetComponent[comp.Player2](world, gameState)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Player2 component: %w", err)
	}
	response.Player2 = player2.Nickname

	return &response, nil
}
