package system

// CreateRectangle constructs a rectangle given a midpoint on the base, direction vector, half-width, and full length.
// The direction vector should be normalized and point along the length of the rectangle.
func CreateRectangleBase(midBase Point, direction Point, halfWidth, fullLength float32) (Point, Point, Point, Point) {
	// Calculate half-length vector along the direction (length goes along the direction vector)
	fullLengths := scaleVector(direction, fullLength)

	// Calculate the perpendicular vector for the width
	perpendicularVector := rotate90Clockwise(direction)
	halfWidthVector := scaleVector(perpendicularVector, halfWidth)

	// Calculate the bottom corners of the rectangle, starting from the midBase
	bottomLeft := addPoints(midBase, scaleVector(halfWidthVector, -1)) // Midpoint to the left
	bottomRight := addPoints(midBase, halfWidthVector)                 // Midpoint to the right

	// Calculate the top corners by moving along the length
	topLeft := addPoints(bottomLeft, fullLengths)
	topRight := addPoints(bottomRight, fullLengths)

	return topLeft, topRight, bottomLeft, bottomRight
}

func CreateRectangleAroundPoint(center Point, direction Point, halfWidth, halfLength float32) (Point, Point, Point, Point) {
	// Normalize the direction vector and calculate half-length vector along the direction
	halfLengths := scaleVector(direction, halfLength)

	// Calculate the perpendicular vector for the width
	perpendicularVector := rotate90Clockwise(direction)
	halfWidthVector := scaleVector(perpendicularVector, halfWidth)

	// Calculate all corners relative to the center
	topLeft := addPoints(addPoints(center, halfLengths), scaleVector(halfWidthVector, -1))
	topRight := addPoints(addPoints(center, halfLengths), halfWidthVector)
	bottomLeft := addPoints(addPoints(center, scaleVector(halfLengths, -1)), scaleVector(halfWidthVector, -1))
	bottomRight := addPoints(addPoints(center, scaleVector(halfLengths, -1)), halfWidthVector)

	return topLeft, topRight, bottomLeft, bottomRight
}

// FindAABB calculates the axis-aligned bounding box that contains the given rectangle.
// It takes the four corners of a rotated rectangle and returns the corners of the AABB.
func FindRectangleAABB(corner1, corner2, corner3, corner4 Point) (Point, Point, Point, Point) {
	// Initialize min and max with the first corner
	minX, maxX := corner1.X, corner1.X
	minY, maxY := corner1.Y, corner1.Y

	// List all corners
	corners := []Point{corner1, corner2, corner3, corner4}

	// Find the minimum and maximum x and y coordinates
	for _, corner := range corners {
		if corner.X < minX {
			minX = corner.X
		}
		if corner.X > maxX {
			maxX = corner.X
		}
		if corner.Y < minY {
			minY = corner.Y
		}
		if corner.Y > maxY {
			maxY = corner.Y
		}
	}

	// Construct the corners of the AABB
	bottomLeft := Point{minX, minY}
	bottomRight := Point{maxX, minY}
	topLeft := Point{minX, maxY}
	topRight := Point{maxX, maxY}

	return topLeft, topRight, bottomLeft, bottomRight
}

// IsPointInRectangle checks if a point is inside a rectangle defined by four corners in order.
func IsPointInRectangle(p Point, topLeft, topRight, bottomRight, bottomLeft Point) bool {
	// Cross-product function to determine if a point is to the left of a line segment
	isLeft := func(a, b, p Point) bool {
		return (b.X-a.X)*(p.Y-a.Y)-(b.Y-a.Y)*(p.X-a.X) >= 0
	}

	// Check if the point is inside the rectangle (counter-clockwise order)
	return isLeft(topLeft, topRight, p) &&
		isLeft(topRight, bottomRight, p) &&
		isLeft(bottomRight, bottomLeft, p) &&
		isLeft(bottomLeft, topLeft, p)
}

// CircleIntersectsRectangle checks if a circle intersects a rectangle defined by its corners.
// circleCenter is the center of the circle, and radius is its radius.
func CircleIntersectsRectangle(circleCenter Point, radius float32, topLeft, topRight, bottomRight, bottomLeft Point) bool {
	// Step 1: Check if the circle's center is inside the rectangle
	if IsPointInRectangle(circleCenter, topLeft, topRight, bottomRight, bottomLeft) {
		return true
	}

	// Step 2: Check for intersection with rectangle edges
	rectangleEdges := [][2]Point{
		{topLeft, topRight},
		{topRight, bottomRight},
		{bottomRight, bottomLeft},
		{bottomLeft, topLeft},
	}

	// Check each edge of the rectangle
	for _, edge := range rectangleEdges {
		if CircleIntersectsEdge(circleCenter, radius, edge[0], edge[1]) {
			return true
		}
	}

	return false
}

// CircleIntersectsEdge checks if a circle intersects a line segment defined by points a and b.
func CircleIntersectsEdge(circleCenter Point, radius float32, a, b Point) bool {
	// Find the closest point on the segment to the circle's center
	closestPoint := ClosestPointOnSegment(circleCenter, a, b)

	// Calculate the distance from the circle's center to this closest point
	distSquared := (circleCenter.X-closestPoint.X)*(circleCenter.X-closestPoint.X) +
		(circleCenter.Y-closestPoint.Y)*(circleCenter.Y-closestPoint.Y)

	// If the distance is less than or equal to the circle's radius squared, they intersect
	return distSquared <= radius*radius
}

// ClosestPointOnSegment finds the point on the line segment between a and b that is closest to p.
func ClosestPointOnSegment(p, a, b Point) Point {
	// Compute the line segment direction vector
	segmentDirX := b.X - a.X
	segmentDirY := b.Y - a.Y

	// Compute vector from a to p
	apX := p.X - a.X
	apY := p.Y - a.Y

	// Calculate projection of ap onto the segment direction
	segmentLengthSquared := segmentDirX*segmentDirX + segmentDirY*segmentDirY
	if segmentLengthSquared == 0 {
		return a // a and b are the same point
	}

	t := (apX*segmentDirX + apY*segmentDirY) / segmentLengthSquared

	// Clamp t to [0, 1] to stay within the segment
	t = max(0, min(1, t))

	// Return the closest point on the segment
	return Point{
		X: a.X + t*segmentDirX,
		Y: a.Y + t*segmentDirY,
	}
}
