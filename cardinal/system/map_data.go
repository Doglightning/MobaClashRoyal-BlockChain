package system

import "fmt"

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
	Bases      [][]int `json:"bases"` //[0=Blue 1= red][x, y, z]
	TowersBlue [][]int `json:"TowersBlue"`
	TowersRed  [][]int `json:"TowersRed"`
	numTowers  int
}

// Maps
var MapDataRegistry = map[string]MapData{
	"ProtoType": {StartX: -5440, StartY: -3660, EndX: 5260, EndY: 4640, Increment: 100, Bases: [][]int{{3860, 500, 100}, {-3680, 700, 100}}, TowersBlue: [][]int{{1920, -1140, 100}}, TowersRed: [][]int{{-2150, 2310, 100}}, numTowers: 1},
}

// Normalize coords to the key required to acsess map data
func normalizeMapCoords(x, y float32, startX, startY, increment int) string {
	// normalize the units position to the maps grid increments.
	normalizedX := int(((int(x)-startX)/increment))*increment + startX
	normalizedY := int(((int(y)-startY)/increment))*increment + startY

	// The units (x,y) coordinates normalized and turned into proper key(string) format for seaching map
	return fmt.Sprintf("%d,%d", normalizedX, normalizedY)
}

func getMapDirection(x, y float32, mapName string) ([]float32, error) {
	//check map data exsists
	mapData, exists := MapDataRegistry[mapName]
	if !exists {
		return nil, fmt.Errorf("error key for MapDataRegistry does not exsist (getMapDirection)")
	}
	//check direction map exsists
	mapDir, ok := MapRegistry[mapName]
	if !ok {
		return nil, fmt.Errorf("error key for MapRegistry does not exsist (getMapDirection)")
	}

	//The units (x,y) coordinates normalized and turned into proper key(string) format for seaching map
	coordKey := normalizeMapCoords(x, y, mapData.StartX, mapData.StartY, mapData.Increment)

	// Retrieve direction vector using coordinate key
	directionVector, exists := mapDir.DMap[coordKey]
	if !exists {
		return nil, fmt.Errorf("no direction vector found for the given coordinates (getMapDirection)")
	}
	return directionVector, nil
}
