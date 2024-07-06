package component

type Player2 struct {
	Nickname2   string       `json:"player2"`
	RemovalList map[int]bool `json:"removallist"`
}

func (Player2) Name() string {
	return "Player2"
}
