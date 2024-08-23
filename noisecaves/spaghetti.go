package noisecaves

import "math"

func Spaghetti(Density float64) bool{
	return math.Abs(Density) < 0.1
}