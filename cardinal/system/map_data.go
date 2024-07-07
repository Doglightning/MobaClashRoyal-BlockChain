package system

type MapData struct {
	StartX    int `json:"startX"`
	StartY    int `json:"startY"`
	EndX      int `json:"endX"`
	EndY      int `json:"endY"`
	Increment int `json:"increment"`

	BlueBase []int `json:"bluebase"`
	RedBase  []int `json:"redbase"`
}

var MapDataRegistry = map[string]MapData{
	"ProtoType": {StartX: -5440, StartY: -3660, EndX: 5260, EndY: 4640, Increment: 100, BlueBase: []int{3860, 500}, RedBase: []int{-3680, 700}},
}

// func getMapData(mapName string) (mapData MapData, exists bool) {
// 	mapData, exists = MapDataRegistry[mapName]
// 	return mapData, exists
// }
