package msg

type RemoveAllEntitiesMsg struct {
	MatchID string
}

type RemoveAllEntitiesResult struct {
	Success bool `json:"success"`
}
