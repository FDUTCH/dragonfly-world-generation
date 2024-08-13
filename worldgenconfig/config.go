package worldgenconfig

import "github.com/df-mc/dragonfly/server/world"

type WGConfig struct{
	GenerateStructures bool
	Dimension          world.Dimension
	WorldSize WorldSize
	WorldPreset WorldPreset
}

type WorldSize int

const(
	SIZE_INFINITE WorldSize = iota
	SIZE_BORDERED
	SIZE_OLD
)

type WorldPreset int

const(
	PRESET_DEFAULT WorldPreset = iota
	PRESET_SUPERFLAT
	AMPLIFIED
	LARGE_BIOMES
	SINGLE_BIOME
)

func DefaultConfig(Dim world.Dimension) WGConfig{
	return WGConfig{
		GenerateStructures: true,
		Dimension: Dim,
		WorldSize: SIZE_INFINITE,
		WorldPreset: PRESET_DEFAULT,
	}
}