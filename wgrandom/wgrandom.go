package wgrandom

import (
	"math"

	"github.com/aquilax/go-perlin"
	"github.com/cnkei/gospline"
)

const ALPHA = 2.0
const BETA = 2.0

const ITERATIONS int32 = 3

const (
	SEED_CONTINENTALNESS = iota
	SEED_EROSION         = iota
	SEED_TEMPERATURE     = iota
	SEED_HUMIDITY        = iota
	SEED_DENSITY         = iota
	SEED_WEIRDNESS       = iota
)

type SubSeed int64

type WGRandom struct {
	Seed int64

	Continentalness *perlin.Perlin
	Erosion         *perlin.Perlin

	Temperature *perlin.Perlin
	Humidity    *perlin.Perlin

	Density   *perlin.Perlin
	Weirdness *perlin.Perlin
}

func New(Seed int64) *WGRandom {
	Continentalness := perlin.NewPerlin(ALPHA, BETA, ITERATIONS, Seed+SEED_CONTINENTALNESS)
	Erosion := perlin.NewPerlin(ALPHA, BETA, ITERATIONS, Seed+SEED_EROSION)
	Temperature := perlin.NewPerlin(ALPHA, BETA, ITERATIONS, Seed+SEED_TEMPERATURE)
	Humidity := perlin.NewPerlin(ALPHA, BETA, ITERATIONS, Seed+SEED_HUMIDITY)
	Density := perlin.NewPerlin(ALPHA, BETA, ITERATIONS, Seed+SEED_DENSITY)
	Weirdness := perlin.NewPerlin(ALPHA, BETA, ITERATIONS, Seed+SEED_WEIRDNESS)

	return &WGRandom{Seed, Continentalness, Erosion, Temperature, Humidity, Density, Weirdness}
}

var ContinentalSpline = gospline.NewCubicSpline([]float64{-1, -0.9, 1}, []float64{1, 0.1, 0.98})

func PeaksValleys(Weirdness float64) float64 {
	return 1 - math.Abs((3*math.Abs(Weirdness))-2)
}
