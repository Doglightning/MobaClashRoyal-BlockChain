package component

type Class struct {
	Class string `json:"Class"`
}

func (Class) Name() string {
	return "Class"
}
