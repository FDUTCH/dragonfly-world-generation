package features

import (
	"math"
	"math/rand/v2"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world/chunk"
	"github.com/go-gl/mathgl/mgl64"
)

type OreFeature struct {
	ID          string
	ReplaceRule []ReplaceRule
	Count       int
}

func (ore OreFeature) Identifier() string {
	return ore.ID
}

func (ore OreFeature) Type() FeatureType {
	return TYPE_ORE
}

func (ore OreFeature) Build(c *chunk.Chunk, f FeaturePos, r *rand.Rand) {
	_, _ = c, f
	chunkBBox := ChunkBBox(c)

	count := rand.IntN(ore.Count)

	cubeDiameter := math.Cbrt(float64(count))
	fBBox := cube.Box(float64(f.X()), float64(f.Y()), float64(f.Z()), float64(f.X())+cubeDiameter, float64(f.Y())+cubeDiameter, float64(f.Z())+cubeDiameter)

	for x := fBBox.Min().X(); x < fBBox.Max().X(); x++ {
		for y := fBBox.Min().Y(); y < fBBox.Max().Y(); y++ {
			for z := fBBox.Min().Z(); z < fBBox.Max().Z(); z++ {
				pos := mgl64.Vec3{x, y, z}

				// Skip if point is outside the chunk
				if chunkBBox.Vec3Within(pos) {
					if count > 0 {
						TryAndReplace(ore.ReplaceRule, cube.PosFromVec3(pos), c)
						count--
					} else {
						return
					}
				}
			}
		}
	}
}

// func oreBBox
