package msg

type RemoveUnitMsg struct {
	MatchId     string
	Team        string
	RemovalList []int
}

type RemoveUnitResult struct {
	Succsess bool `json:"Succsess"`
}
