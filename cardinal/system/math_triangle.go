package system

import "math"

// CreateIsoscelesTriangle constructs an isosceles triangle given an apex, direction, height, and base width.
// The direction vector should point from the apex towards the base.
func CreateIsoscelesTriangle(apex Point, direction Point, height, baseWidth float32) (Point, Point, Point) {

	direction.X = -direction.X
	direction.Y = -direction.Y
	// Calculate the base center by moving from the apex in the opposite direction of the provided vector.
	baseCenter := addPoints(apex, scaleVector(direction, -height))

	// Calculate the perpendicular vector for the base endpoints.
	perpendicularVector := rotate90Clockwise(direction)

	// Calculate the half-width vector for the base.
	halfBase := scaleVector(perpendicularVector, baseWidth/2)

	// Calculate the left and right points of the base.
	baseLeft := addPoints(baseCenter, halfBase)
	baseRight := addPoints(baseCenter, scaleVector(halfBase, -1))

	return apex, baseLeft, baseRight
}

// CalculateBoundingBox finds the minimal bounding box for a triangle defined by three points.
func CalculateBoundingBox(a, b, c Point) (Point, Point) {
	minX := math.Min(float64(a.X), math.Min(float64(b.X), float64(c.X)))
	maxX := math.Max(float64(a.X), math.Max(float64(b.X), float64(c.X)))
	minY := math.Min(float64(a.Y), math.Min(float64(b.Y), float64(c.Y)))
	maxY := math.Max(float64(a.Y), math.Max(float64(b.Y), float64(c.Y)))
	return Point{float32(minX), float32(minY)}, Point{float32(maxX), float32(maxY)}
}

// PointInTriangle checks if a point P is inside the triangle formed by A, B, C
func PointInTriangle(p, a, b, c Point) bool {
	// Using barycentric coordinates to determine if point is inside the triangle
	det := (b.Y-c.Y)*(a.X-c.X) + (c.X-b.X)*(a.Y-c.Y)
	lambda1 := ((b.Y-c.Y)*(p.X-c.X) + (c.X-b.X)*(p.Y-c.Y)) / det
	lambda2 := ((c.Y-a.Y)*(p.X-c.X) + (a.X-c.X)*(p.Y-c.Y)) / det
	lambda3 := 1 - lambda1 - lambda2
	return lambda1 >= 0 && lambda2 >= 0 && lambda3 >= 0
}

func RasterizeIsoscelesTriangle(apex, baseLeft, baseRight Point) []Point {
	topLeft, bottomRight := CalculateBoundingBox(apex, baseLeft, baseRight)

	var points []Point

	// Iterate over all points in the bounding box
	for x := math.Floor(float64(topLeft.X)); x <= math.Ceil(float64(bottomRight.X)); x++ {
		for y := math.Floor(float64(topLeft.Y)); y <= math.Ceil(float64(bottomRight.Y)); y++ {
			point := Point{float32(x), float32(y)}
			if PointInTriangle(point, apex, baseLeft, baseRight) {
				points = append(points, point)
			}
		}
	}
	return points
}
