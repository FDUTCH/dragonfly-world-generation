package wgrandom

import (
	"math"
	"github.com/aquilax/go-perlin"
	"github.com/cnkei/gospline"
)

const ALPHA = 2.0
const BETA = 2.0

const ITERATIONS int32 = 1

const OVERWORLD_SCALE = 171.103
const OVERWORLD_HEIGHT_SCALE = 85.5515
const DEPTH_SCALE = 50

const SURFACE_SCALE = 2.138

const (
	SEED_CONTINENTALNESS = iota
	SEED_EROSION
	SEED_TEMPERATURE
	SEED_HUMIDITY
	SEED_DENSITY
	SEED_WEIRDNESS
	SEED_SURFACE
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

	Surface *perlin.Perlin
}

func New(Seed int64) *WGRandom {
	Continentalness := perlin.NewPerlin(2, 2, 1, Seed+SEED_CONTINENTALNESS)
	Erosion := perlin.NewPerlin(ALPHA, BETA, ITERATIONS, Seed+SEED_EROSION)
	Temperature := perlin.NewPerlin(ALPHA, BETA, ITERATIONS, Seed+SEED_TEMPERATURE)
	Humidity := perlin.NewPerlin(ALPHA, BETA, ITERATIONS, Seed+SEED_HUMIDITY)
	Density := perlin.NewPerlin(2, BETA, ITERATIONS, Seed+SEED_DENSITY)
	Weirdness := perlin.NewPerlin(2, BETA, ITERATIONS, Seed+SEED_WEIRDNESS)
	Surface := perlin.NewPerlin(2, BETA, ITERATIONS, Seed+SEED_SURFACE)

	return &WGRandom{Seed, Continentalness, Erosion, Temperature, Humidity, Density, Weirdness, Surface}
}

var ContinentalSpline = gospline.NewMonotoneSpline(
	[]float64{-1, -0.9, -0.47, -0.43, -0.19, -0.034, 0, 0.03, 1},
	[]float64{1, 0.1, 0.03, 0.3, 0.41, 0.9, 0.9, 0.93, 0.98},
)

var ErosionSpline = gospline.NewMonotoneSpline(
	[]float64{-1, -0.93, -0.76, -0.32, -0.12, 0.3, 0.5, 0.53, 0.67, 0.72, 1},
	[]float64{1, 0.75, 0.24, 0.26, 0.13, 0.11, 0.08, 0.12, 0.13, 0.06, 0},
)

var PVSpline = gospline.NewMonotoneSpline(
	[]float64{-1, -0.9, -0.85, 0, 0.4, 1},
	[]float64{0, 0.3, 0.53, 0.61, 0.9, 0.9},
)

func PeaksValleys(Weirdness float64) float64 {
	return 1 - math.Abs((3*math.Abs(Weirdness))-2)
}
