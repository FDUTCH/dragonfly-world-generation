package terrain

import (
	"github.com/Ikarolyi/dragonfly-world-generation/internal"
	"github.com/Ikarolyi/dragonfly-world-generation/wgrandom"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/chunk"
)

const TERRAIN_TRESHOLD = 0
const SEA_LEVEL = 62
const BASE_HEIGHT = 30

var NORMAL_WATER = block.Water{Still: false, Depth: 8, Falling: false}

func GenerateTerrain(
	chunkPos world.ChunkPos,
	chunk *chunk.Chunk,
	WGRand *wgrandom.WGRandom,
) {
	chunkWorldPos := []float64{float64(chunkPos[0]) * 16, float64(chunkPos[1]) * 16}

	min, max := int16(chunk.Range().Min()), int16(chunk.Range().Max())

	for x := uint8(0); x < 16; x++ {
		for z := uint8(0); z < 16; z++ {
			NoiseX, NoiseZ := (chunkWorldPos[0]+float64(x))/wgrandom.OVERWORLD_SCALE, (chunkWorldPos[1]+float64(z))/wgrandom.OVERWORLD_SCALE
			C := wgrandom.ContinentalSpline.At(WGRand.Continentalness.Noise2D(NoiseX, NoiseZ))
			E := wgrandom.ErosionSpline.At(WGRand.Erosion.Noise2D(NoiseX, NoiseZ))
			PV := wgrandom.PVSpline.At(wgrandom.PeaksValleys(WGRand.Weirdness.Noise2D(NoiseX, NoiseZ)))

			height := (E * 175) + (C * 39) + (PV * 15)

			for y := min; y < max; y++ {
				NoiseY := float64(y) / wgrandom.OVERWORLD_HEIGHT_SCALE
				D := WGRand.Density.Noise3D(NoiseX*3, NoiseY*3, NoiseZ*3)
				Squish := height - float64(y)

				// solid := y < int16(height)
				solid := D*30+Squish > 0
				if solid {
					chunk.SetBlock(x, y, z, 0, world.BlockRuntimeID(block.Stone{}))
				} else if y < SEA_LEVEL {
					chunk.SetBlock(x, y, z, 0, world.BlockRuntimeID(NORMAL_WATER))
				}
			}
		}
	}
}

// Check if a block is filled or not.
// This is not the method used to make the terrain, because
func GetBlock(Pos cube.Pos, WGRand *wgrandom.WGRandom) bool {
	x, y, z := Pos.X(), Pos.Y(), Pos.Z()

	x += internal.NoiseOffset
	z += internal.NoiseOffset

	NoiseX, NoiseZ := (float64(x))/wgrandom.OVERWORLD_SCALE, (float64(z))/wgrandom.OVERWORLD_SCALE
	C := wgrandom.ContinentalSpline.At(WGRand.Continentalness.Noise2D(NoiseX, NoiseZ))
	E := wgrandom.ErosionSpline.At(WGRand.Erosion.Noise2D(NoiseX, NoiseZ))
	PV := wgrandom.PVSpline.At(wgrandom.PeaksValleys(WGRand.Weirdness.Noise2D(NoiseX, NoiseZ)))

	height := (E * 175) + (C * 39) + (PV * 15)

	NoiseY := float64(y) / wgrandom.OVERWORLD_HEIGHT_SCALE
	D := WGRand.Density.Noise3D(NoiseX*3, NoiseY*3, NoiseZ*3)
	Squish := height - float64(y)

	solid := D*30+Squish > 0

	return solid
}
