package component

type SizeCircle struct {
	Radius int `json:"radius"`
}

func (SizeCircle) Name() string {
	return "SizeCircle"
}
