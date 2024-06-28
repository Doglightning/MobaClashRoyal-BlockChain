package component

type MatchId struct {
	MatchId string `json:"MatchId"`
}

func (MatchId) Name() string {
	return "MatchId"
}
