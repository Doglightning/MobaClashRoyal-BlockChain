package component

type GridUtils struct {
	StartX    int `json:"startX"`
	StartY    int `json:"startY"`
	EndX      int `json:"endX"`
	EndY      int `json:"endY"`
	Increment int `json:"increment"`

	BlueX int `json:"blueX"`
	BlueY int `json:"blueY"`
	RedX  int `json:"redX"`
	RedY  int `json:"redY"`
}

func (GridUtils) Name() string {
	return "GridUtils"
}
