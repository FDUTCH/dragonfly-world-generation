package features

import (
	"slices"

	"github.com/Ikarolyi/dragonfly-world-generation/biomemap"
	"github.com/Ikarolyi/dragonfly-world-generation/terrain"
	"github.com/Ikarolyi/dragonfly-world-generation/wgrandom"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

type FeatureConditionType int

const(
	TYPE_BIOME_FILTER FeatureConditionType = iota
)

type Operator int

const(
	OP_EQ Operator = iota
	OP_NOT_EQ
)

type FeatureCondition struct {
	Type FeatureConditionType 
	Test string
	Operator Operator
	Value string
}

func EvalConditions(ChunkPos world.ChunkPos, Origin FeaturePos, Conditions []FeatureCondition, WGRandom *wgrandom.WGRandom) bool{
	chunkWorldPos := [2]float64{float64(ChunkPos[0]) * 16, float64(ChunkPos[1]) * 16}

	originBiome := biomemap.Column(chunkWorldPos, float64(Origin[0]), float64(Origin[2]), WGRandom)
	b, _ := world.BiomeByID(int(originBiome))
	bName := b.String()

	var result bool = false
	for _, c := range Conditions{
		switch c.Type{
			case TYPE_BIOME_FILTER:
				if c.Test == "has_biome_tag"{
					if slices.Contains(biomemap.BiomeTags[bName], c.Value){
						result = true
					}
				}
			default:
				result = true
		}
	}

	if len(Conditions) == 0{
		return true
	}

	return result
}

func Compare(Op Operator, X, Y any) bool{
	switch Op{
		case OP_EQ:
			return X == Y
		case OP_NOT_EQ:
			return X != Y
		default:
			return true
	}
}

type PlacementPassType int

const(
	SURFACE_PASS PlacementPassType = iota
	UNDERGROUND_PASS
)

func EvalPlacementPass(fr FeatureRule, Origin FeaturePos, ChunkPos world.ChunkPos, r *wgrandom.WGRandom) bool{
	chunkWorldPos := [2]int{int(ChunkPos[0]) * 16, int(ChunkPos[1]) * 16}
	globalPos := []int{chunkWorldPos[0] + Origin[0], chunkWorldPos[1] + Origin[2]}
	b0 := terrain.GetBlock(cube.Pos{globalPos[0], Origin[1], globalPos[1]}, r)
	b5 := terrain.GetBlock(cube.Pos{globalPos[0], Origin[1] + 5, globalPos[1]}, r)
	bn1 := terrain.GetBlock(cube.Pos{globalPos[0], Origin[1] - 1, globalPos[1]}, r)

	if fr.PlacementPass == SURFACE_PASS{
		return !b0 && !b5 && bn1
	}else if fr.PlacementPass == UNDERGROUND_PASS{
		return b0 && b5 && bn1
		// return false
	}

	return false
}