package query

import (
	comp "MobaClashRoyal/component"
	"MobaClashRoyal/system"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/iterators"
	"pkg.world.dev/world-engine/cardinal/search/filter"
)

type PSMatchIdRequest struct {
	MatchId string
	Team    string
}

type PlayerStateResponse struct {
	Units []int
	Hand  []string
	Deck  []string
	Gold  float32
}

// get a list of all units to be removed for a player to maintian replication
func PlayerState(world cardinal.WorldContext, req *PSMatchIdRequest) (*PlayerStateResponse, error) {
	var response PlayerStateResponse
	var removeList = []int{}

	//find gameState using matchID
	matchFilter := cardinal.ComponentFilter(func(m comp.MatchId) bool {
		return m.MatchId == req.MatchId
	})

	gameState, err := cardinal.NewSearch().Entity(
		filter.Exact(system.GameStateFilters())).
		Where(matchFilter).First(world)

	if err != nil {
		return nil, fmt.Errorf("error searching for team (Removal State Query): %w", err)
	}
	if gameState == iterators.BadID {
		return nil, fmt.Errorf("no match found with ID or missing components: %s", req.MatchId)
	}

	if req.Team == "Blue" {
		// Get Player1 component
		player1, err := cardinal.GetComponent[comp.Player1](world, gameState)
		if err != nil {
			return nil, fmt.Errorf("error retrieving Player1 component (Removal State Query): %w", err)
		}
		//get the remove list component for Player1
		for key := range player1.RemovalList {
			removeList = append(removeList, key)
		}

		response.Hand = player1.Hand
		response.Deck = player1.Deck
		//player1 gold
		response.Gold = player1.Gold

	} else {
		// Get Player2 component
		player2, err := cardinal.GetComponent[comp.Player2](world, gameState)
		if err != nil {
			return nil, fmt.Errorf("error retrieving Player2 component (Removal State Query): %w", err)
		}
		//get the remove list component for Player2
		for key := range player2.RemovalList {
			removeList = append(removeList, key)
		}

		response.Hand = player2.Hand
		response.Deck = player2.Deck
		//player2 gold
		response.Gold = player2.Gold
	}
	response.Units = removeList

	return &response, nil
}
