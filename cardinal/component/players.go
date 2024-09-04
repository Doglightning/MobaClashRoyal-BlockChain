package component

type Player1 struct {
	Nickname    string       `json:"player1"`
	Hand        []string     `json:"Hand"`
	Deck        []string     `json:"Deck"`
	RemovalList map[int]bool `json:"removallist"`
	Gold        float32      `json:"Gold"`
}

type Player2 struct {
	Nickname    string       `json:"player2"`
	Hand        []string     `json:"Hand"`
	Deck        []string     `json:"Deck"`
	RemovalList map[int]bool `json:"removallist"`
	Gold        float32      `json:"Gold"`
}

func (Player1) Name() string {
	return "Player1"
}

func (Player2) Name() string {
	return "Player2"
}
