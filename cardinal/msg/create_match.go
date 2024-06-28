package msg

type CreateMatchMsg struct {
	MatchID string
	MapName string
}

type CreateMatchResult struct {
	Success bool `json:"success"`
}
