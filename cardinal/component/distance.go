package component

type Distance struct {
	Distance float32 `json:"distnace"`
}

func (Distance) Name() string {
	return "Distance"
}
