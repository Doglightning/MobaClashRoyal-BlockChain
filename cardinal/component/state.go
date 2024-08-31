package component

type State struct {
	State string `json:"State"`
}

func (State) Name() string {
	return "State"
}
