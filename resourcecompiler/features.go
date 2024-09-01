package resourcecompiler

import (
	"os"
	"path/filepath"

	"github.com/Ikarolyi/dragonfly-world-generation/features"
	"github.com/Ikarolyi/dragonfly-world-generation/wgrandom"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/yosuke-furukawa/json5/encoding/json5"
)

var featureRootTypes = []string{
	"minecraft:ore_feature",
	"minecraft:aggregate_feature",
	"minecraft:scatter_feature",
	"minecraft:single_block_feature",
	"minecraft:search_feature",
	"minecraft:tree_feature",
	"minecraft:fossil_feature",
	"minecraft:sequence_feature",
	"minecraft:weighted_random_feature",
	"minecraft:nether_cave_carver_feature",
	"minecraft:cave_carver_feature",
	"minecraft:underwater_cave_carver_feature",
}

var coordinateEvalOrderStrings = []string{
	"xyz",
	"xzy",
	"yxz",
	"yzx",
	"xzy",
	"xyz",
}

var distributionTypeStrings = []string{
	"uniform",
	"gaussian",
	"inverse_gaussian",
	"triangle",
	"fixed_grid",
	"jittered_grid",
}

func woodBlockID(b map[string]interface{}, Stripped bool, Axis cube.Axis) uint32 {
	// name := b["name"].(string)
	states := b["states"].(map[string]interface{})
	oldLogType := states["old_log_type"].(string)

	var logType block.WoodType

	for _, t := range block.WoodTypes() {
		if t.String() == oldLogType {
			logType = t
		}
	}

	return world.BlockRuntimeID(block.Log{Wood: logType, Stripped: Stripped, Axis: Axis})
}

func CompileFeatures() {
	featuresPath := filepath.Join(definitionPath, "features")
	features.FeatureTable = make(map[string]features.Feature)

	FeaturesDir, err := os.ReadDir(featuresPath)
	if err != nil {
		panic(err)
	}

	for _, f := range FeaturesDir {
		featureFile, err := os.ReadFile(filepath.Join(featuresPath, f.Name()))
		if err != err {
			panic(err)
		}

		data := make(map[string]interface{})

		json5.Unmarshal(featureFile, &data)

		var root map[string]interface{}
		var featureType features.FeatureType
		for newFeatureType, rootType := range featureRootTypes {
			newRoot, ok := data[rootType]
			if ok {
				root = newRoot.(map[string]interface{})
			}
			featureType = features.FeatureType(newFeatureType)
		}

		if root == nil {
			println("RootType not found in %s", f.Name())
			continue
		}

		id := root["description"].(map[string]interface{})["identifier"].(string)

		var newFeature features.Feature

		switch featureType {
		case features.TYPE_TREE:
			trunk := root["trunk"].(map[string]interface{})
			newFeature = features.TreeFeature{
				ID:          id,
				TrunkHeight: [2]int{},
				TrunkBlock:  woodBlockID(trunk["trunk_block"].(map[string]interface{}), false, cube.Y),
			}
		default:
			newFeature = features.NopFeature{
				ID: id,
			}
		}

		// newFeature := features.Feature{
		// 	Identifier: id,
		// 	Count: 0,
		// 	ReplaceRules: features.ReplaceRule{
		// 		PlacesBlock: 0,
		// 		MayReplace:  []uint32{},
		// 	},
		// }

		features.FeatureTable[id] = newFeature
	}
}

func CompileFeatureRules() {
	featuresPath := filepath.Join(definitionPath, "feature_rules")

	FeaturesDir, err := os.ReadDir(featuresPath)
	if err != nil {
		panic(err)
	}

	fileloop:
	for _, fr := range FeaturesDir {
		featureRuleFile, err := os.ReadFile(filepath.Join(featuresPath, fr.Name()))
		if err != err {
			panic(err)
		}

		data := make(map[string]interface{})

		json5.Unmarshal(featureRuleFile, &data)

		root := data["minecraft:feature_rules"].(map[string]interface{})
		_ = root

		description := root["description"].(map[string]interface{})
		id := description["identifier"].(string)
		places := description["places_feature"].(string)

		distribution, ok := root["distribution"].(map[string]interface{})
		if !ok{
			continue fileloop
		}

		iterations, ok := distribution["iterations"].(float64)
		if !ok{
			continue fileloop
		}

		scatterChanceRaw, scatterChanceOk := distribution["scatter_chance"]

		var scatterChance wgrandom.Chance
		if scatterChanceOk {
			s, ok := scatterChanceRaw.(map[string]interface{})
			if !ok{
				continue fileloop
			}
			scatterChanceNumerator, scatterChanceDenumerator := s["denominator"].(float64), s["numerator"].(float64)
			scatterChance = wgrandom.NewChance(scatterChanceNumerator, scatterChanceDenumerator)
		} else {
			scatterChance = wgrandom.NewChance(1, 1)
		}

		var dimensions = []string{"x", "y", "z"}

		var xyz [3]features.CoordRange
		for i, dimension := range dimensions {
			value := distribution[dimension]

			ranged, isRanged := value.(map[string]interface{})
			if isRanged {
				distribution := ranged["distribution"].(string)
				var d features.DistributionType
				for i, distrTypeStr := range distributionTypeStrings {
					if distrTypeStr == distribution {
						d = features.DistributionType(i)
					}
				}

				extent := ranged["extent"].([]interface{})

				//TODO: grid_size, step_size

				min, ok := extent[0].(float64)
				if !ok{
					continue fileloop
				}

				max, ok := extent[1].(float64)
				if !ok{
					continue fileloop
				}


				xyz[i] = features.CoordRange{
					DistributionType: d,
					Min:              int(min),
					Max:              int(max),
				}
				continue
			} else {
				simple, ok := value.(float64)
				if !ok{
					continue fileloop
				}
				xyz[i] = features.ConstCoordRange(int(simple))
			}
		}

		var coordinateEvalOrder features.CoordinateEvalOrder
		rawCoordinateEvalOrder, ok := distribution["coordinate_eval_order"]
		if ok{
			for i, str := range coordinateEvalOrderStrings {
				if str == rawCoordinateEvalOrder.(string){
					coordinateEvalOrder = features.CoordinateEvalOrder(i)
				}
			}
		}else{
			coordinateEvalOrder = features.XZY
		}

		features.FeatureRuleTable = append(features.FeatureRuleTable, features.FeatureRule{
			Identifier: id,
			Places:     places,
			Conditions: []features.FeatureCondition{},
			Distribution: features.FeatureDistribution{
				Iterations:          int(iterations),
				ScatterChance:       scatterChance,
				CoordinateEvalOrder: coordinateEvalOrder,
				X:                   xyz[0],
				Y:                   xyz[1],
				Z:                   xyz[2],
			},
		})
	}
}
