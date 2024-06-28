package component

type Team struct {
	Team string `json:"team"`
}

func (Team) Name() string {
	return "Team"
}
