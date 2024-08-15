package worldgen

import (
	"github.com/Ikarolyi/dragonfly-world-generation/biomemap"
	"github.com/Ikarolyi/dragonfly-world-generation/terrain"
	"github.com/Ikarolyi/dragonfly-world-generation/wgrandom"
	"github.com/Ikarolyi/dragonfly-world-generation/worldgenconfig"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/chunk"
)

type WorldGenerator struct {
	WGConfig worldgenconfig.WGConfig
	WGRandom           *wgrandom.WGRandom
}

func (gen WorldGenerator) GenerateChunk(pos world.ChunkPos, chunk *chunk.Chunk) {
	biomemap.FillChunk(pos, chunk, gen.WGRandom, gen.WGConfig)
	terrain.GenerateTerrain(pos, chunk, gen.WGRandom)
}

func NewWorldGenerator(Seed int64) func(world.Dimension) world.Generator {
	return func(d world.Dimension) world.Generator {
		return WorldGenerator{
			WGConfig: worldgenconfig.DefaultConfig(d),
			WGRandom: wgrandom.New(Seed),
		}
	}
}
