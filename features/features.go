package features

import (
	"math/rand/v2"
	"slices"

	"github.com/Ikarolyi/dragonfly-world-generation/wgrandom"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/chunk"
)

var FeatureTable map[string]Feature
var FeatureRuleTable []FeatureRule

type Feature interface {
	Identifier() string
	Type() FeatureType
	Build(*chunk.Chunk, FeaturePos, *rand.Rand)
}

type ReplaceRule struct {
	PlacesBlock uint32
	MayReplace  []uint32
}

func TryAndReplace(allR []ReplaceRule, Pos cube.Pos, c *chunk.Chunk) {
	bBefore := c.Block(uint8(Pos.X()), int16(Pos.Y()), uint8(Pos.Z()),0)
	for _, r := range allR{
		if slices.Contains(r.MayReplace, bBefore){
			c.SetBlock(uint8(Pos.X()), int16(Pos.Y()), uint8(Pos.Z()), 0, r.PlacesBlock)
			break
		}
	}
}

type FeatureType int

const (
	TYPE_ORE FeatureType = iota
	TYPE_AGGREGATE
	TYPE_SCATTER
	TYPE_SINGLE_BLOCK
	TYPE_SEARCH
	TYPE_TREE
	TYPE_FOSSIL
	TYPE_SEQUENCE
	TYPE_WEIGHED_RANDOM
	TYPE_NETHER_CARVER
	TYPE_CAVE_CARVER
	TYPE_UNDERWATER_CARVER
	TYPE_NOOP
)

func ChunkBBox(c *chunk.Chunk) cube.BBox{
	min, max := c.Range().Min(), c.Range().Max()

	return cube.Box(-0.1, float64(min), -0.1, 15.1, float64(max), 15.1)
}

type TreeFeature struct {
	ID          string
	TrunkHeight [2]int
	TrunkBlock  uint32
}

func (tree TreeFeature) Identifier() string {
	return tree.ID
}

func (tree TreeFeature) Type() FeatureType {
	return TYPE_TREE
}

func (tree TreeFeature) Build(c *chunk.Chunk, f FeaturePos, _ *rand.Rand){
	_, _ = c, f
}

type NopFeature struct {
	ID string
}

func (f NopFeature) Identifier() string {
	return f.ID
}

func (f NopFeature) Build(c *chunk.Chunk, fp FeaturePos, _ *rand.Rand){
	if fp[0] > 0 && fp[0] < 16{
		if fp[2] > 0 && fp[2] < 16{
			c.SetBlock(uint8(fp[0]), int16(fp[1]), uint8(fp[2]), 0, world.BlockRuntimeID(block.Log{}))
		} 
	}
	_, _ = c, fp
}


func (f NopFeature) Type() FeatureType {
	return TYPE_NOOP
}

type FeaturePos cube.Pos

func (f FeaturePos) X() int{
	return f[0]
}

func (f FeaturePos) Y() int{
	return f[1]
}

func (f FeaturePos) Z() int{
	return f[2]
}

type CoordinateEvalOrder int

const (
	XYZ CoordinateEvalOrder = iota
	XZY
	YXZ
	YZX
	ZXY
	ZYX
)

type DistributionType int

const (
	UNIFORM DistributionType = iota
	GAUSSIAN
	INVERSE_GAUSSIAN
	TRIANGLE
	FIXED_GRID
	JITTERED_GRID
)

type CoordRange struct {
	DistributionType DistributionType
	Min, Max         int
}

func ConstCoordRange(Val int) CoordRange {
	return CoordRange{DistributionType: UNIFORM, Min: Val, Max: Val}
}

type FeatureDistribution struct {
	Iterations          int
	ScatterChance       wgrandom.Chance
	CoordinateEvalOrder CoordinateEvalOrder
	X,
	Y,
	Z CoordRange
}

type FeatureRule struct {
	Identifier,
	Places string
	Conditions   []FeatureCondition
	Distribution FeatureDistribution
	PlacementPass PlacementPassType
}

func FeatureRandom(x, z int64, ft FeatureType, seed int64) rand.Source {
	// X and Z has to be shuffled to make sure there's no symmetry around the origin

	x2 := (x >> 1) | (x << 1)
	z2 := (z >> 1) | (z << 1)
	// xz := rand.NewPCG(uint64(x2), uint64(z2))

	s := seed * (1 + int64(ft))
	s2 := (s >> 1) | (s << 1)

	return rand.NewPCG(uint64(s2 & x2), uint64(z2))
	// return nil
}

func (fr FeatureRule) PopulateChunk(
	chunkPos world.ChunkPos,
	chunk *chunk.Chunk,
	WGRand *wgrandom.WGRandom,
) {
	// if fr.Distribution.X.DistributionType == UNIFORM{
	// TODO other distibution types and coord eval orders

	var origins []FeaturePos
	var randSources []*rand.Rand

	places := FeatureTable[fr.Places]
	if places == nil{
		return
	}

	// Loop trough every chunk in the 3x3 square around the chunk
	// to be able to generate features placed on two or more chunks
	for cX := int32(-1); cX <= 1; cX++ {
		for cZ := int32(-1); cZ <= 1; cZ++ {
			gx, gz := int64(cX+chunkPos.X()), int64(cZ+chunkPos.Z())

			_, _ = gx, gz

			source := FeatureRandom(gx, gz, places.Type(), (*WGRand).Seed)

			R := rand.New(source)
			

			for range fr.Distribution.Iterations {
				pos := FeaturePos{
					(R.IntN(max(fr.Distribution.X.Max-fr.Distribution.X.Min, 1)) + fr.Distribution.X.Min),
					(R.IntN(max(fr.Distribution.Y.Max-fr.Distribution.Y.Min, 1)) + fr.Distribution.Y.Min),
					(R.IntN(max(fr.Distribution.Z.Max-fr.Distribution.Z.Min, 1)) + fr.Distribution.Z.Min),
				}

				scatter := fr.Distribution.ScatterChance.Eval(R)
				if scatter{
					origins = append(origins, pos)
					randSources = append(randSources, R)
				}
			}
		}
	}

	for i, pos := range origins{
		if EvalConditions(chunkPos, pos, fr.Conditions, WGRand){
			if EvalPlacementPass(fr, pos, chunkPos, WGRand){
				places.Build(chunk, pos, randSources[i])
			}
		}
	}
}

func PopulateChunk(
	chunkPos world.ChunkPos,
	chunk *chunk.Chunk,
	WGRand *wgrandom.WGRandom,
) {
	for _, fr := range FeatureRuleTable {
		fr.PopulateChunk(chunkPos, chunk, WGRand)
	}
}
