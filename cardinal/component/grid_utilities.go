package component

type GridUtils struct {
	StartX    int `json:"startX"`
	StartY    int `json:"startY"`
	EndX      int `json:"endX"`
	EndY      int `json:"endY"`
	Increment int `json:"increment"`

	BlueX        int `json:"blueX"`
	BlueY        int `json:"blueY"`
	BlueLength   int `json:"bluelength"`
	BlueWidth    int `json:"bluewidth"`
	BlueRotation int `json:"bluerotation"`

	RedX        int `json:"redX"`
	RedY        int `json:"redY"`
	RedLength   int `json:"redlength"`
	RedWidth    int `json:"redwidth"`
	RedRotation int `json:"redrotation"`
}

func (GridUtils) Name() string {
	return "GridUtils"
}
