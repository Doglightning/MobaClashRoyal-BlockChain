package system

// size of collision hash map. (should be atleast 1.5x the size of the largest unit)
var SpatialGridCellSize = 300

type MapData struct {
	//Direction Map Data
	StartX    int `json:"startX"`
	StartY    int `json:"startY"`
	EndX      int `json:"endX"`
	EndY      int `json:"endY"`
	Increment int `json:"increment"`

	//Sturcture spawn points
	Bases      [][]int `json:"bases"` //[0=Blue 1= red][x, y]
	TowersBlue [][]int `json:"TowersBlue"`
	TowersRed  [][]int `json:"TowersRed"`
	numTowers  int
}

// Maps
var MapDataRegistry = map[string]MapData{
	"ProtoType": {StartX: -5440, StartY: -3660, EndX: 5260, EndY: 4640, Increment: 100, Bases: [][]int{{3860, 500}, {-3680, 700}}, TowersBlue: [][]int{{1920, -1140}}, TowersRed: [][]int{{-2150, 2310}}, numTowers: 1},
}

type StructureData struct {
	Health float32 `json:"health"`
	Radius int
}

// structures
var StructureDataRegistry = map[string]StructureData{
	"Base":  {Health: 200, Radius: 240},
	"Tower": {Health: 150, Radius: 150},
}
