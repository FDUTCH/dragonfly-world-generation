package worldgen

import (
	"github.com/Ikarolyi/dragonfly-world-generation/biomemap"
	"github.com/Ikarolyi/dragonfly-world-generation/noisecaves"
	"github.com/Ikarolyi/dragonfly-world-generation/resourcecompiler"
	"github.com/Ikarolyi/dragonfly-world-generation/surface"
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

// Loads resources needed
func LoadResources(Definitions string, Structures string){	
	resourcecompiler.SetPaths(Definitions, Structures)
	resourcecompiler.CompileAll()
}

func (gen WorldGenerator) GenerateChunk(pos world.ChunkPos, chunk *chunk.Chunk) {
	resourcecompiler.CheckForResources()
	terrain.GenerateTerrain(pos, chunk, gen.WGRandom)
	biomemap.FillChunk(pos, chunk, gen.WGRandom, gen.WGConfig)
	surface.GenerateSurface(pos, chunk, gen.WGRandom)
	noisecaves.GenerateNoiseCaves(pos, chunk, gen.WGRandom)
}

func NewDefaultGenerator(Seed int64) func(world.Dimension) world.Generator {
	return func(d world.Dimension) world.Generator {
		return WorldGenerator{
			WGConfig: worldgenconfig.DefaultConfig(d),
			WGRandom: wgrandom.New(Seed),
		}
	}
}
