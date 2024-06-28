package component

// DirectionMap is the component that holds mapping from coordinates to vectors directly within.
type DirectionMap struct {
	Map map[string][]float32 `json:"map"`
}

func (DirectionMap) Name() string {
	return "DirectionMap"
}
