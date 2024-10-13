package system

import "math"

// func subtract(x1, y1, x2, y2 float32) (x, y float32) {
// 	return x1 - x2, y1 - y2
// }

func addVector2D(x1, y1, x2, y2 float32) (x, y float32) {
	return x1 + x2, y1 + y2
}

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

	if magnitude == 0 {
		return 0, 0
	}

	return deltaX / magnitude, deltaY / magnitude
}

func directionVectorBetweenTwoPoints3D(x1, y1, z1, x2, y2, z2 float32) (float32, float32, float32) {
	// Compute direction vector towards the enemy
	deltaX := x2 - x1
	deltaY := y2 - y1
	deltaZ := z2 - z1
	magnitude := float32(math.Sqrt(float64(deltaX*deltaX + deltaY*deltaY + deltaZ*deltaZ)))

	if magnitude == 0 {
		return 0, 0, 0
	}

	return deltaX / magnitude, deltaY / magnitude, deltaZ / magnitude
}

// rotate returns a new vector, rotated by the given angle (in degrees)
func rotateVectorDegrees(dirX, dirY float32, angle float64) (float32, float32) {
	radian := angle * math.Pi / 180
	cosTheta := float32(math.Cos(radian))
	sinTheta := float32(math.Sin(radian))
	return dirX*cosTheta - dirY*sinTheta, dirX*sinTheta + dirY*cosTheta

}

// Compute the new position given a rotation vector, forward offset, and lateral offset
func RelativeOffsetXY(x, y, rX, rY, lateralOffset, forwardOffset float32) (float32, float32) {
	// Rotate the rotation vector 90 degrees to get the lateral direction (right)
	lateralX, lateralY := rotateVectorDegrees(rX, rY, 90)

	// Move forward by multiplying the rotation vector by the forward offset
	movedForwardX := rX * forwardOffset
	movedForwardY := rY * forwardOffset //Vector2{X: rot.X * forwardOffset, Y: rot.Y * forwardOffset}

	// Move laterally by multiplying the lateral vector by the lateral offset
	movedLateralX := lateralX * lateralOffset // Vector2{X: lateralVector.X * lateralOffset, Y: lateralVector.Y * lateralOffset}
	movedLateralY := lateralY * lateralOffset

	// Combine both movements
	newPositionX, newPositionY := addVector2D(movedForwardX, movedForwardY, movedLateralX, movedLateralY)
	return addVector2D(x, y, newPositionX, newPositionY)
}
