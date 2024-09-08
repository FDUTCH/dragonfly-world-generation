package wgrandom

import (
	"math/rand/v2"
	"testing"
)

func TestInterpolation(t *testing.T) {
	for i := float64(-1); i < 20; i += 0.1 {
		println("X:%v Y:%v", i, ContinentalSpline.At(i))
	}
}

func TestChances(t *testing.T) {
	OneToTen := NewChance(1, 11)

	const ITERATIONS = 100
	success := 0
	r := rand.New(rand.NewPCG(10, 22))
	for range(ITERATIONS){
		if OneToTen.Eval(r){
			success += 1
		}
	}

	t.Logf("Success rate: %v/%v", success, ITERATIONS)
}