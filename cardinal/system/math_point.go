package system

type Point struct {
	X, Y float32
}

// Rotate a vector clockwise by 90 degrees to get the perpendicular vector
func rotate90Clockwise(v Point) Point {
	return Point{-v.Y, v.X}
}

// Scale a vector by a scalar
func scaleVector(v Point, scale float32) Point {
	return Point{v.X * scale, v.Y * scale}
}

// Add two points
func addPoints(p1, p2 Point) Point {
	return Point{p1.X + p2.X, p1.Y + p2.Y}
}
