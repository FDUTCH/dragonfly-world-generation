package noisecaves

import (
	"github.com/Ikarolyi/dragonfly-world-generation/internal"
	"math"

	"github.com/Ikarolyi/dragonfly-world-generation/terrain"
	"github.com/Ikarolyi/dragonfly-world-generation/wgrandom"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/chunk"
)

func invSqrtAbs(x float64) float64 {
	return 1 / math.Sqrt(math.Abs(x))
}

func GenerateNoiseCaves(
	chunkPos world.ChunkPos,
	chunk *chunk.Chunk,
	WGRand *wgrandom.WGRandom,
) {
	globalX, globalZ := float64(chunkPos.X())*16, float64(chunkPos.Z())*16
	globalZ += internal.NoiseOffset
	globalZ += internal.NoiseOffset
	min := int16(chunk.Range().Min())
	for x := uint8(0); x < 16; x++ {
		for z := uint8(0); z < 16; z++ {

			TopBlock := chunk.HighestBlock(x, z)
			SurfaceDetail2D := WGRand.Surface.Noise2D((float64(x)+globalX)/6.7, (float64(z)+globalZ)/6.7)

			var cave = false
			for y := min; y <= TopBlock; y++ {
				if chunk.Block(x, y, z, 0) == world.BlockRuntimeID(terrain.NORMAL_WATER) {
					continue
				}
				density := WGRand.Density.Noise3D((float64(x)+globalX)/wgrandom.DEPTH_SCALE, float64(y)/wgrandom.DEPTH_HEIGHT_SCALE, (float64(z)+globalZ)/wgrandom.DEPTH_SCALE)
				ridges := WGRand.Weirdness.Noise3D((float64(x)+globalX)/wgrandom.DEPTH_SCALE, float64(y)/wgrandom.DEPTH_HEIGHT_SCALE, (float64(z)+globalZ)/wgrandom.DEPTH_SCALE)
				surface := WGRand.Surface.Noise3D((float64(x)+globalX)/wgrandom.DEPTH_SCALE, float64(y)/wgrandom.DEPTH_HEIGHT_SCALE, (float64(z)+globalZ)/wgrandom.DEPTH_SCALE)

				depth := float64(TopBlock - y)
				entrance := depth < 3 && cave

				_, _, _ = density, ridges, entrance

				// Only cheese caves have pillars, otherwise the caves are filled
				if Spaghetti(density, ridges, surface, depth) || (Cheese(density, ridges, surface, depth) && !Pillar(surface, SurfaceDetail2D)) || entrance {
					chunk.SetBlock(x, y, z, 0, world.BlockRuntimeID(block.Air{}))
					cave = true
				} else {
					cave = false
				}
			}
		}
	}
}

func Cheese(Density, Ridges, Surface, Depth float64) bool {
	return Surface+(invSqrtAbs(Depth)*2.53) < -0.04
}

const SPAGHETTI_TRESHOLD = 0.0792614
const NOODLE_TRESHOLD = 0.051365
const PILLAR_TRESHOLD = 0.4056

func Spaghetti(Density, Ridges, Surface, Depth float64) bool {
	var treshold float64
	if Surface < 0 {
		treshold = SPAGHETTI_TRESHOLD
	} else {
		treshold = NOODLE_TRESHOLD
	}
	return math.Abs(Density)+math.Abs(Ridges)+(invSqrtAbs(Depth)/17) < treshold
}

func Pillar(Surface, SurfaceDetail2D float64) bool {
	return math.Abs(SurfaceDetail2D)*1.1+math.Max(math.Min(Surface, 0)+0.41, 0)*1.5 > PILLAR_TRESHOLD
	// return math.Abs(SurfaceDetail2D) - math.Abs(Surface)*2 < PILLAR_TRESHOLD
	// return Surface - Density2D < -0.7
}
