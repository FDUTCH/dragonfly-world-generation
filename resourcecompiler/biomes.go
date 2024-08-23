package resourcecompiler

import (
	"os"
	"path/filepath"

	"github.com/yosuke-furukawa/json5/encoding/json5"

	"github.com/Ikarolyi/dragonfly-world-generation/surface"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

var blockIdAlias = map[string]world.Block{
	"minecraft:water": block.Water{Still: false, Depth: 8, Falling: false},
	"minecraft:grass": block.Grass{},
	"minecraft:lava": block.Lava{Still: false, Depth: 8, Falling: true},
	"minecraft:ice": block.BlueIce{},
	"minecraft:red_sand": block.Sand{Red: true},
	"minecraft:stained_hardened_clay": block.StainedTerracotta{Colour: item.ColourBlack()},

	// Yet to be implemented:
	"minecraft:crimson_nylium": block.Netherrack{},
	"minecraft:mycelium": block.Podzol{},
	"minecraft:warped_nylium": block.Netherrack{},
}


func findBlock(id string)uint32{
	// Correct blocks with alias names
	b, isAlias := blockIdAlias[id]
	if isAlias{
		return world.BlockRuntimeID(b)
	}
	
	// Try to find the base variation of the block
	it, ok := world.ItemByName(id, 0)
	if ok{
		b := it.(world.Block)
		return world.BlockRuntimeID(b)
	}

	// Try to find a block with a matching name
	for _, it := range world.Items(){
		name, _ := it.EncodeItem()
		if id == name {
			return world.BlockRuntimeID(it.(world.Block))
		}
	}

	println(id + " not found, fallback to air")

	// Fallback to air
	return world.BlockRuntimeID(block.Air{})
}

func CompileBiomes() {
	surface.SurfaceTable = make(map[int]surface.SurfaceParams)
	surface.MesaTable = make(map[int]surface.MesaSurface)

	biomesPath := filepath.Join(definitionPath, "biomes")

	biomesDir, err := os.ReadDir(biomesPath)
	if err != nil {
		panic(err)
	}

	for _, b := range biomesDir {
		biomeFile, err := os.ReadFile(filepath.Join(biomesPath, b.Name()))
		if err != err {
			panic(err)
		}

		data := make(map[string]interface{})
		

		json5.Unmarshal(biomeFile, &data)
		// version := data["format_version"].(string)
		biomeGroup, sup := data["minecraft:biome"].(map[string]interface{})
		if !sup{
			panic(b.Name())
		}
		identifier := biomeGroup["description"].(map[string]interface{})["identifier"].(string)

		b, _ := world.BiomeByName(identifier)
		components := biomeGroup["components"].(map[string]interface{})
		surfaceParams, surfaceOk := components["minecraft:surface_parameters"].(map[string]interface{})

		mesaSurface, mesaOk := components["minecraft:mesa_surface"].(map[string]interface{})

		if mesaOk{
			surfaceParams = mesaSurface
			surfaceOk = true
			brycePillars := mesaSurface["bryce_pillars"].(bool)
			clayMaterial := mesaSurface["clay_material"].(string)
			hardClayMaterial := mesaSurface["hard_clay_material"].(string)
			hasForest := mesaSurface["has_forest"].(bool)
			
			surface.MesaTable[b.EncodeBiome()] = surface.MesaSurface{
				BrycePillars: brycePillars,
				ClayMaterial: findBlock(clayMaterial),
				HardClayMaterial: findBlock(hardClayMaterial),
				HasForest: hasForest,
			}
		}

		if surfaceOk{
			SeaFloorDepth := surfaceParams["sea_floor_depth"].(float64)
			SeaFloorMaterial := surfaceParams["sea_floor_material"].(string)
			FoundationMaterial := surfaceParams["foundation_material"].(string)
			MidMaterial := surfaceParams["mid_material"].(string)
			TopMaterial := surfaceParams["top_material"].(string)
			SeaMaterial := surfaceParams["sea_material"].(string)


			surface.SurfaceTable[b.EncodeBiome()] = surface.SurfaceParams{
				SeaFloorDepth: int16(SeaFloorDepth),
				SeaFloorMaterial: findBlock(SeaFloorMaterial),
				FoundationMaterial: findBlock(FoundationMaterial),
				MidMaterial: findBlock(MidMaterial),
				TopMaterial: findBlock(TopMaterial),
				SeaMaterial: findBlock(SeaMaterial),
			}
		}
	}
}
