package component

type IntTracker struct {
	Num int `json:"num"`
}

func (IntTracker) Name() string {
	return "IntTracker"
}
