package terrain

import (
	"github.com/Ikarolyi/dragonfly-world-generation/wgrandom"
	// "github.com/df-mc/dragonfly/server/world"
	// "github.com/df-mc/dragonfly/server/world/chunk"
)

const TERRAIN_TRESHOLD = 0

func Terrain2D(WGRand *wgrandom.WGRandom) int{
	return 60 + 1
}

// func GenerateTerrain(
// 	chunkPos world.ChunkPos,
// 	chunk *chunk.Chunk,
// 	WGRand *wgrandom.WGRandom,
// ) {
// 	chunkWorldPos := []float64{float64(chunkPos[0]) * 16, float64(chunkPos[1]) * 16}

// 	min, max := int16(chunk.Range().Min()), int16(chunk.Range().Max())

// 	for x := uint8(0); x < 16; x++ {
// 		for z := uint8(0); z < 16; z++ {
// 			NoiseX, NoiseZ := chunkWorldPos[0], chunkWorldPos[1]
// 			t := Terrain2D(WGRand)

// 			for y := int16(0); y <= max; y++ {
// 				NoiseY := float64(y)


// 				if y < f.n {
// 					chunk.SetBlock(x, min+y, z, 0, f.layers[f.n-y-1])
// 				}
// 			}
// 		}
// 	}
// }
