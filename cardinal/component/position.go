package component

type Position struct {
	PositionVectorX float32 `json:"positionvectorx"`
	PositionVectorY float32 `json:"positionvectory"`
	PositionVectorZ float32 `json:"positionvectorz"`

	RotationVectorX float32 `json:"rotationvectorx"`
	RotationVectorY float32 `json:"rotationvectory"`
	RotationVectorZ float32 `json:"rotationvectorz"`
}

func (Position) Name() string {
	return "Position"
}
