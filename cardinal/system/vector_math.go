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

// func normalize(x, y float32) (float32, float32) {
// 	mag := math.Sqrt(float64(x*x + y*y))
// 	if mag > 0 {
// 		return x / float32(mag), y / float32(mag)
// 	}
// 	return 0, 0 // Return zero vector if magnitude is zero to prevent division by zero
// }

func dotProduct(x1, y1, x2, y2 float32) float32 {
	return float32(x1*x2 + y1*y2)
}

// Function to calculate distance between two points
func distanceBetweenTwoPoints(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}
