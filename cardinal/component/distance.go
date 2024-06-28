package component

type Distance struct {
	Distance float64 `json:"distnace"`
}

func (Distance) Name() string {
	return "Distance"
}
