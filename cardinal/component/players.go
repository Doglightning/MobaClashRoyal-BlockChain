package component

type Player1 struct {
	Nickname    string       `json:"player1"`
	RemovalList map[int]bool `json:"removallist"`
}

type Player2 struct {
	Nickname    string       `json:"player2"`
	RemovalList map[int]bool `json:"removallist"`
}

func (Player1) Name() string {
	return "Player1"
}

func (Player2) Name() string {
	return "Player2"
}
