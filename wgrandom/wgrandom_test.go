package wgrandom

import (
	"testing"
)

func TestInterpolation(t *testing.T) {
	for i := float64(-1); i < 20; i += 0.1 {
		println("X:%v Y:%v", i, ContinentalSpline.At(i))
	}
}