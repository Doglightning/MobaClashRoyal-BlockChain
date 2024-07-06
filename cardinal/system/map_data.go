package system

type MapData struct {
	DirMap    DMap `json:"dirmap"`
	StartX    int  `json:"startX"`
	StartY    int  `json:"startY"`
	EndX      int  `json:"endX"`
	EndY      int  `json:"endY"`
	Increment int  `json:"increment"`

	BlueBase []int `json:"bluebase"`
	RedBase  []int `json:"redbase"`
}

var MapDataRegistry = map[string]MapData{
	"ProtoType": {DirMap: getDirMaps("ProtoType"), StartX: -5440, StartY: -3660, EndX: 5260, EndY: 4640, Increment: 100, BlueBase: []int{3860, 500}, RedBase: []int{-3680, 700}},
}

func getMapData(mapName string) (mapData MapData, found bool) {
	mapData, found = MapDataRegistry[mapName]
	if !found {
		return MapData{}, false
	}

	return mapData, true
}
