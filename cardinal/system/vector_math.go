package system

import "math"

// func subtract(x1, y1, x2, y2 float32) (x, y float32) {
// 	return x1 - x2, y1 - y2
// }

// func add(x1, y1, x2, y2 float32) (x, y float32) {
// 	return x1 + x2, y1 + y2
// }

// func multiply(x, y float32, scalar float32) (float32, float32) {
// 	return x * scalar, y * scalar
// }

func normalize(x, y float32) (float32, float32) {
	mag := math.Sqrt(float64(x*x + y*y))
	if mag > 0 {
		return x / float32(mag), y / float32(mag)
	}
	return 0, 0 // Return zero vector if magnitude is zero to prevent division by zero
}

func dotProduct(x1, y1, x2, y2 float32) float32 {
	return float32(x1*x2 + y1*y2)
}

func crossProduct(x1, y1, x2, y2 float32) float32 {
	return x1*y2 - y1*x2
}

// Function to calculate distance between two points
func distanceBetweenTwoPoints(x1, y1, x2, y2 float32) float32 {
	return float32(math.Sqrt(math.Pow(float64(x2-x1), 2) + math.Pow(float64(y2-y1), 2)))
}

func directionVectorBetweenTwoPoints(x1, y1, x2, y2 float32) (float32, float32) {
	// Compute direction vector towards the enemy
	deltaX := x2 - x1
	deltaY := y2 - y1
	magnitude := float32(math.Sqrt(math.Pow(float64(deltaX), 2) + math.Pow(float64(deltaY), 2)))

	return deltaX / magnitude, deltaY / magnitude
}

// rotate returns a new vector, rotated by the given angle (in degrees)
func rotateVectorDegrees(dirX, dirY float32, angle float64) (float32, float32) {
	radian := angle * math.Pi / 180
	cosTheta := float32(math.Cos(radian))
	sinTheta := float32(math.Sin(radian))
	return dirX*cosTheta - dirY*sinTheta, dirX*sinTheta + dirY*cosTheta

}
