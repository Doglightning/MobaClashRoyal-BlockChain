package component

type Player1 struct {
	Nickname string `json:"player1"`
}

func (Player1) Name() string {
	return "Player1"
}
