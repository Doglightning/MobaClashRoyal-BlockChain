package component

type Player2 struct {
	Nickname2 string `json:"player2"`
}

func (Player2) Name() string {
	return "Player2"
}
