package surface

import (
	"github.com/Ikarolyi/dragonfly-world-generation/terrain"
	"github.com/Ikarolyi/dragonfly-world-generation/wgrandom"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/chunk"
)

var SurfaceTable map[int]SurfaceParams

type SurfaceParams struct {
	SeaFloorDepth int16
	SeaFloorMaterial,
	FoundationMaterial,
	MidMaterial,
	TopMaterial,
	SeaMaterial uint32
}

var MesaTable map[int]MesaSurface

type MesaSurface struct {
	BrycePillars bool
	ClayMaterial,
	HardClayMaterial uint32
	HasForest bool
}

func floor(x uint8, y int16, z uint8, chnk *chunk.Chunk) bool {
	if chnk.Range().Max()-1 == int(y) {
		return false
	}
	return isEmpty(x, y+1, z, chnk)
}

func ceiling(x uint8, y int16, z uint8, chnk *chunk.Chunk) bool {
	if chnk.Range().Min() == int(y) {
		return false
	}
	return isEmpty(x, y-1, z, chnk)
}

func bedrock(x uint8, y int16, z uint8) bool {
	_, _ = x, z
	return y == -64
}

func isBiome(b uint32, biome world.Biome) bool {
	return b == uint32(biome.EncodeBiome())
}

func noWaterAbove(x uint8, y int16, z uint8, chnk *chunk.Chunk) bool {
	if chnk.Range().Max()-1 == int(y) {
		return false
	}

	blockAbove := chnk.Block(x, y+1, z, 0)
	return blockAbove == world.BlockRuntimeID(terrain.NORMAL_WATER)
}

func isEmpty(x uint8, y int16, z uint8, chnk *chunk.Chunk) bool {
	return chnk.Block(x, y, z, 0) == world.BlockRuntimeID(block.Air{})
}

func isWater(x uint8, y int16, z uint8, chnk *chunk.Chunk) bool {
	return chnk.Block(x, y, z, 0) == world.BlockRuntimeID(terrain.NORMAL_WATER)
}

func placeHodoo(x uint8, y int16, z uint8, chnk *chunk.Chunk) {
	chnk.SetBlock(x, y, z, 0, world.BlockRuntimeID(block.StainedTerracotta{Colour: item.ColourBlack()}))
}

func noiseS3(n float64) bool {
	if n > -0.909 && n < -0.5454 {
		return false
	} else if n > 0.5454 && n < 0.909 {
		return true
	} else {
		return false
	}
}

func fillColumn(x, z uint8, chnk *chunk.Chunk, chunkPos world.ChunkPos, r *wgrandom.WGRandom) {
	min, max := int16(chnk.Range().Min()), int16(chnk.Range().Max())
	TopBlock := chnk.HighestBlock(x, z)

	SeaFloor := max
	Surface := TopBlock

	for y := TopBlock; y >= min; y-- {
		if isWater(x, y, z, chnk) {
			SeaFloor = y
		}

		if !isEmpty(x, y, z, chnk) {
			b := pickBlock(x, y, z, Surface, SeaFloor, chnk, r)
			chnk.SetBlock(x, y, z, 0, b)
		} else {
			Surface = y
		}
	}
}

func pickBlock(x uint8, y int16, z uint8, Surface, SeaFloor int16, chnk *chunk.Chunk, r *wgrandom.WGRandom) uint32 {
	// ceiling := ceiling(x, y , z, chnk)
	b := chnk.Biome(x, y, z)
	_ = r
	// noise := r.Surface.Noise3D(float64(chunkPos.X() * 16 + int32(x)) / wgrandom.SURFACE_SCALE, float64(y) / wgrandom.SURFACE_SCALE * 2, float64(chunkPos.X() * 16 + int32(x)) / wgrandom.SURFACE_SCALE)

	// swampNoise := false
	// erosion := false

	// maxWaterDepth := min(terrain.SEA_LEVEL-y, 0)

	floor := floor(x, y, z, chnk)
	surface := Surface-y < 5

	surfaceParams, normalSurface := SurfaceTable[int(b)]
	if normalSurface {
		if isWater(x, y, z, chnk) {
			return surfaceParams.SeaMaterial
		} else if (SeaFloor - y) < int16(surfaceParams.SeaFloorDepth) {
			return surfaceParams.SeaFloorMaterial
		} else if floor {
			return surfaceParams.TopMaterial
		} else if surface {
			return surfaceParams.MidMaterial
		} else {
			return surfaceParams.FoundationMaterial
		}
	}

	return 0
}

func GenerateSurface(
	chunkPos world.ChunkPos,
	chunk *chunk.Chunk,
	WGRand *wgrandom.WGRandom,
) {

	for x := uint8(0); x < 16; x++ {
		for z := uint8(0); z < 16; z++ {
			fillColumn(x, z, chunk, chunkPos, WGRand)
		}
	}
}
