package component

type CenterOffset struct {
	CenterOffset float32 `json:"CenterOffset"`
}

func (CenterOffset) Name() string {
	return "CenterOffset"
}
