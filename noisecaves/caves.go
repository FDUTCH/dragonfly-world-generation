package noisecaves

import (
	"github.com/Ikarolyi/dragonfly-world-generation/wgrandom"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/chunk"
)

func GenerateSurface(
	chunkPos world.ChunkPos,
	chunk *chunk.Chunk,
	WGRand *wgrandom.WGRandom,
) {
	globalX, globalZ := float64(chunkPos.X()), float64(chunkPos.Z())
	min, max := int16(chunk.Range().Min()), int16(chunk.Range().Max())
	for x := uint8(0); x < 16; x++ {
		for z := uint8(0); z < 16; z++ {
			for y := min; y < max; y++ {
				density := WGRand.Density.Noise3D((float64(x) + globalX)/wgrandom.DEPTH_SCALE, float64(y) / wgrandom.DEPTH_SCALE, (float64(z) + globalZ) / wgrandom.DEPTH_SCALE)
				if Spaghetti(density){
					chunk.SetBlock(x,y,z, 0 , world.BlockRuntimeID(block.Air{}))
				}
			}
		}
	}
}
