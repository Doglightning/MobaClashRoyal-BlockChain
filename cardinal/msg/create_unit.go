package msg

type CreateUnitMsg struct {
	MatchID string
	MapName string
	Team    string

	UnitType  string
	PositionX float32
	PositionY float32
	PositionZ float32
	RotationX float32
	RotationY float32
	RotationZ float32
}

type CreateUnitResult struct {
	Success bool `json:"success"`
}
