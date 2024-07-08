package system

type MapData struct {
	StartX    int `json:"startX"`
	StartY    int `json:"startY"`
	EndX      int `json:"endX"`
	EndY      int `json:"endY"`
	Increment int `json:"increment"`

	Bases  [][]int `json:"bases"` //[0=Blue 1= red][x, y]
	Towers [][]int `json:"towers"`
}

var MapDataRegistry = map[string]MapData{
	"ProtoType": {StartX: -5440, StartY: -3660, EndX: 5260, EndY: 4640, Increment: 100, Bases: [][]int{{3860, 500}, {-3680, 700}}},
}

type StructureData struct {
	Health float32 `json:"health"`
	Radius int
}

var StructureDataRegistry = map[string]StructureData{
	"Base": {Health: 200, Radius: 240},
}

// func getMapData(mapName string) (mapData MapData, exists bool) {
// 	mapData, exists = MapDataRegistry[mapName]
// 	return mapData, exists
// }
