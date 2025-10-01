package entity

import (
	"math"
)

func almostEqual(a, b, epsilon float64) bool {
	return math.Abs(a-b) <= epsilon
}
