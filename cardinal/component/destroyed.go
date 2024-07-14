package component

type Destroyed struct {
	Destroyed bool `json:"Destroyed"`
}

func (Destroyed) Name() string {
	return "Destroyed"
}
