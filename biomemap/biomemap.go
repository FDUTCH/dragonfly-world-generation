package biomemap

import (
	"github.com/Ikarolyi/dragonfly-world-generation/wgrandom"
	"github.com/Ikarolyi/dragonfly-world-generation/worldgenconfig"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/chunk"
)

// 0 ~ 6
func ErosionIndex(Erosion float64) int {
	if Erosion < -0.78 {
		return 0
	} else if Erosion < -0.375 {
		return 1
	} else if Erosion < -0.2225 {
		return 2
	} else if Erosion < 0.05 {
		return 3
	} else if Erosion < 0.45 {
		return 4
	} else if Erosion < 0.55 {
		return 5
	} else {
		return 6
	}
}

type ContinentalnessEnum int

const (
	MUSHROOM_ISLAND ContinentalnessEnum = iota
	DEEP_OCEAN
	OCEAN
	COAST
	NEAR_INLAND
	MID_INLAND
	FAR_INLAND
)

func ContinentalnessIndex(Continentalness float64) ContinentalnessEnum {
	// Biome continentalness ranges from -1.2 to +1.0 so it has to be scaled first (from the range -1.0 -> +1.0)
	c := (Continentalness / 2.2 * 2) - 0.1

	if c < -1.05 {
		return MUSHROOM_ISLAND
	} else if c < -0.455 {
		return DEEP_OCEAN
	} else if c < -0.19 {
		return OCEAN
	} else if c < -0.11 {
		return COAST
	} else if c < 0.03 {
		return NEAR_INLAND
	} else if c < 0.3 {
		return MID_INLAND
	} else {
		return FAR_INLAND
	}
}

// 0 ~ 4
func TemperatureIndex(Temperature float64) int {
	if Temperature < -0.45 {
		return 0
	} else if Temperature < -0.15 {
		return 1
	} else if Temperature < 0.2 {
		return 2
	} else if Temperature < 0.55 {
		return 3
	} else {
		return 4
	}
}

// 0 ~ 4
func HumidityIndex(Humidity float64) int {
	if Humidity < -0.35 {
		return 0
	} else if Humidity < -0.1 {
		return 1
	} else if Humidity < 0.1 {
		return 2
	} else if Humidity < 0.3 {
		return 3
	} else {
		return 4
	}
}

type PVEnum int

const (
	VALLEYS PVEnum = iota
	PV_LOW
	PV_MID
	PV_HIGH
	PEAKS
)

func PVIndex(Weirdness float64) PVEnum {
	PeaksValleys := wgrandom.PeaksValleys(Weirdness)

	if PeaksValleys < -0.85 {
		return VALLEYS
	} else if PeaksValleys < -0.6 {
		return PV_LOW
	} else if PeaksValleys < 0.2 {
		return PV_MID
	} else if PeaksValleys < 0.7 {
		return PV_HIGH
	} else {
		return PEAKS
	}
}

func WeirdnessIndex(Weirdness float64) int {
	if Weirdness < 0 {
		return 1
	} else {
		return 0
	}
}

func SelectBiome(Continentalness, Erosion, Temperature, Humidity, Weirdness float64) uint32 {
	ContinentalIdx := ContinentalnessIndex(Continentalness)
	ErosionIdx := ErosionIndex(Erosion)
	TemperatureIdx := TemperatureIndex(Temperature)
	HumidityIdx := HumidityIndex(Humidity)
	PVIdx := PVIndex(Weirdness)
	WIdx := WeirdnessIndex(Weirdness)

	biome := BiomeTable[ContinentalIdx][ErosionIdx][TemperatureIdx][HumidityIdx][PVIdx][WIdx]
	return uint32(biome.EncodeBiome())
}

func FillChunk(
	chunkPos world.ChunkPos,
	chunk *chunk.Chunk,
	WGRand *wgrandom.WGRandom,
	WGConfig worldgenconfig.WGConfig,
) {
	min, max := chunk.SubIndex(int16(chunk.Range().Min())), chunk.SubIndex(int16(chunk.Range().Max()))
	chunkWorldPos := []float64{float64(chunkPos[0]) * 16, float64(chunkPos[1]) * 16}

	for x := uint8(0); x < 16; x++ {
		for z := uint8(0); z < 16; z++ {
			NoiseX, NoiseY := chunkWorldPos[0]+float64(x), chunkWorldPos[1]+float64(z)

			Continentalness := WGRand.Continentalness.Noise2D(NoiseX, NoiseY)
			Erosion := WGRand.Erosion.Noise2D(NoiseX, NoiseY)
			Temperature := WGRand.Temperature.Noise2D(NoiseX, NoiseY)
			Humidity := WGRand.Humidity.Noise2D(NoiseX, NoiseY)
			Weirdness := WGRand.Weirdness.Noise2D(NoiseX, NoiseY)

			b := SelectBiome(Continentalness, Erosion, Temperature, Humidity, Weirdness)

			// Apply Biome only on every subchunk
			for y := min; y < max; y += 16 {
				chunk.SetBiome(x, y, z, b)
			}
		}
	}
}
