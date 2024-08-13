package biomemap

import (
	"os"
)

func GenerateTable() {
	var result string = "package biomemap\n"
	result += "\n"
	result += "import (\n"
	result += "\"github.com/df-mc/dragonfly/server/world\"\n"
	result += "\"github.com/df-mc/dragonfly/server/world/biome\"\n"
	result += ")\n"
	result += "\n"
	result += "\n"
	result += "// C E T H PV W"
	result += "\n"
	result += "var BiomeTable = [7][7][5][5][5][2]world.Biome{"
	for C := MUSHROOM_ISLAND; C <= FAR_INLAND; C++ {
		result += "{"
		for E := int(0); E <= 6; E++ {
			result += "{"
			for T := int(0); T <= 4; T++ {
				result += "{"
				for H := int(0); H <= 4; H++ {
					result += "{"
					for PV := VALLEYS; PV <= PEAKS; PV++ {
						result += "{"
						for W := int(0); W <= 1; W++ {
							result += gen1_21_1(C, E, T, H, PV, W) + ","
						}
						result += "},"
					}
					result += "},"
				}
				result += "},"
			}
			result += "},"
		}
		result += "},"
	}
	result += "}"

	data := []byte(result)

	os.WriteFile("./overworldtable.go", data, 0644)
}

func gen1_21_1(C ContinentalnessEnum, E, T, H int, PV PVEnum, W int) string {
	if C == OCEAN {
		switch T {
		case 0:
			return "biome.FrozenOcean{}"
		case 1:
			return "biome.ColdOcean{}"
		case 2:
			return "biome.Ocean{}"
		case 3:
			return "biome.LukewarmOcean{}"
		case 4:
			return "biome.WarmOcean{}"
		}
	} else if C == DEEP_OCEAN {
		switch T {
		case 0:
			return "biome.DeepFrozenOcean{}"
		case 1:
			return "biome.DeepColdOcean{}"
		case 2:
			return "biome.DeepOcean{}"
		case 3:
			return "biome.DeepLukewarmOcean{}"
		case 4:
			return "biome.WarmOcean{}"
		}
	} else if C == MUSHROOM_ISLAND {
		return "biome.MushroomFields{}"
	}

	// Inland Biomes

	if PV == VALLEYS {
		if C == COAST || (C == NEAR_INLAND && E <= 5) || (E >= 2 && E <= 5) {
			goto river
		} else {
			if E == 6 {
				if T == 0 {
					return "biome.FrozenRiver{}"
				} else if T <= 3 { // 1; 2
					return "biome.Swamp{}"
				} else { // 3; 4
					return "biome.MangroveSwamp{}"
				}
			} else {
				if T < 4 {
					goto middle_biomes
				} else {
					goto badland_biomes
				}
			}
		}
	} else if PV == PV_LOW {
		if C == COAST {
			if E <= 2 {
				return "biome.StonyShore{}"
			} else if E <= 4 {
				goto beach_biomes
			} else if E == 5 {
				if W == 1 {
					goto beach_biomes
				} else {
					if T <= 1 || H == 4 {
						goto middle_biomes
					} else {
						return "biome.WindsweptSavanna{}"
					}
				}
			} else {
				goto beach_biomes
			}
		}

		if E <= 1 {
			if C != NEAR_INLAND && T == 0 {
				if H <= 1 {
					return "biome.SnowySlopes{}"
				} else {
					return "biome.Grove{}"
				}
			} else if T < 4 {
				goto middle_biomes
			} else {
				goto badland_biomes
			}
		}

		if C == NEAR_INLAND {
			if E >= 2 && E <= 4 {
				goto middle_biomes
			} else if E == 5 {
				if T <= 1 || H == 4 {
					goto middle_biomes
				} else {
					return "biome.WindsweptSavanna{}"
				}
			}
		}

		if (C == MID_INLAND || C == FAR_INLAND) && (E == 4 || E == 5) {
			goto middle_biomes
		}

		if E == 6 && C != COAST {
			if T == 0 {
				goto middle_biomes
			} else if T <= 2 {
				return "biome.Swamp{}"
			} else {
				return "biome.MangroveSwamp{}"
			}
		}
	} else if PV == PV_MID {
		if C == COAST {
			if E <= 2 {
				return "biome.StonyShore{}"
			} else if E == 3 {
				goto middle_biomes
			} else if E == 4 || E == 6 {
				if W == 1 {
					goto beach_biomes
				} else {
					goto middle_biomes
				}
			} else { // E = 5
				if W == 1 {
					goto beach_biomes
				} else {
					if T <= 1 || H == 4 {
						goto middle_biomes
					} else {
						return "biome.WindsweptSavanna{}"
					}
				}
			}
		} else {
			if E == 0 {
				if T < 3 {
					if H < 2 {
						return "biome.SnowySlopes{}"
					} else {
						return "biome.Grove{}"
					}
				} else {
					goto plateau_biomes
				}
			} else if E == 1 {
				if C != FAR_INLAND {
					if T == 0 {
						if H < 2 {
							return "biome.SnowySlopes{}"
						} else {
							return "biome.Grove{}"
						}
					} else if T < 4 {
						goto middle_biomes
					} else {
						goto badland_biomes
					}
				} else {
					if T == 0 {
						if H < 2 {
							return "biome.SnowySlopes{}"
						} else {
							return "biome.Grove{}"
						}
					} else {
						goto plateau_biomes
					}
				}
			} else if E <= 4 {
				if C == NEAR_INLAND {
					goto middle_biomes
				} else if E < 4 {
					if C == MID_INLAND {
						if T < 4 {
							goto middle_biomes
						} else {
							goto badland_biomes
						}
					} else if E == 2 {
						goto plateau_biomes
					} else {
						if T < 4 {
							goto middle_biomes
						} else {
							goto badland_biomes
						}
					}
				}
			} else if E == 5 {
				if C == NEAR_INLAND {
					if W == 1 || T <= 1 || H == 4 {
						goto middle_biomes
					} else {
						return "biome.WindsweptSavanna{}"
					}
				} else {
					goto shattered_biomes
				}
			} else {
				if T == 0 {
					goto middle_biomes
				} else if T <= 2 {
					return "biome.Swamp{}"
				} else {
					return "biome.MangroveSwamp{}"
				}
			}
		}
	} else if PV == PV_HIGH {
		if E <= 4 {
			if C == COAST {
				goto middle_biomes
			} else if C == NEAR_INLAND {
				if E == 0 {
					if T < 3 {
						if H <= 1 {
							return "biome.SnowySlopes{}"
						} else {
							return "biome.Grove{}"
						}
					} else {
						goto plateau_biomes
					}
				} else if E == 1 {
					if T == 0 {
						if H <= 1 {
							return "biome.SnowySlopes{}"
						} else {
							return "biome.Grove{}"
						}
					} else if T <= 3 {
						goto middle_biomes
					} else {
						goto badland_biomes
					}
				} else {
					goto middle_biomes
				}
			} else {
				if E == 0 {
					if T <= 2 {
						if W == 1 {
							return "biome.JaggedPeaks{}"
						} else {
							return "biome.FrozenPeaks{}"
						}
					} else if T == 3 {
						return "biome.StonyPeaks{}"
					} else {
						goto badland_biomes
					}
				} else if E == 1 {
					if T < 3 {
						if H < 2 {
							return "biome.SnowySlopes{}"
						} else {
							return "biome.Grove{}"
						}
					} else {
						goto plateau_biomes
					}
				} else if E <= 3 {
					if C == MID_INLAND {
						if E == 2 {
							goto plateau_biomes
						} else {
							if T < 4 {
								goto middle_biomes
							} else {
								goto badland_biomes
							}
						}
					} else {
						goto badland_biomes
					}
				} else {
					goto middle_biomes
				}
			}
		} else if E == 5 {
			if C == COAST || C == NEAR_INLAND {
				if W == 1 || T <= 1 || H == 4 {
					goto middle_biomes
				} else {
					return "biome.WindsweptSavanna{}"
				}
			} else {
				goto shattered_biomes
			}
		} else {
			goto middle_biomes
		}
	} else { // PEAKS
		if E == 6 {
			goto middle_biomes
		} else {
			if C == COAST || C == NEAR_INLAND {
				if E == 0 {
					if T < 3 {
						if W == 1 {
							return "biome.JaggedPeaks{}"
						} else {
							return "biome.FrozenPeaks{}"
						}
					} else if T == 3 {
						return "biome.StonyPeaks{}"
					} else {
						goto badland_biomes
					}
				} else if E == 1 {
					if T == 0 {
						if H <= 1 {
							return "biome.SnowySlopes{}"
						} else {
							return "biome.Grove{}"
						}
					} else if T <= 3 {
						goto middle_biomes
					} else {
						goto badland_biomes
					}
				} else if E <= 4 {
					goto middle_biomes
				} else {
					if W == 1 || T <= 1 || H == 4 {
						goto shattered_biomes
					} else {
						return "biome.WindsweptSavanna{}"
					}
				}
			} else {
				if E <= 1 {
					if T <= 2 {
						if W == 1 {
							return "biome.JaggedPeaks{}"
						} else {
							return "biome.FrozenPeaks{}"
						}
					} else if T == 3 {
						return "biome.StonyPeaks{}"
					} else {
						goto badland_biomes
					}
				} else if E <= 3 {
					if C == MID_INLAND {
						if E == 2 {
							goto plateau_biomes
						} else {
							if T < 4 {
								goto middle_biomes
							} else {
								goto badland_biomes
							}
						}
					} else {
						goto badland_biomes
					}
				} else if E == 4 {
					goto middle_biomes
				}
			}
		}
	}

river:
	if T == 0 {
		return "biome.FrozenRiver{}"
	} else {
		return "biome.River{}"
	}

middle_biomes:
	if T == 0 {
		if H == 0 {
			if W == 1 {
				return "biome.SnowyPlains{}"
			} else {
				return "biome.IceSpikes{}"
			}
		} else if H == 1 {
			return "biome.SnowyPlains{}"
		} else if H == 2 {
			if W == 1 {
				return "biome.SnowyPlains{}"
			} else {
				return "biome.SnowyTaiga{}"
			}
		} else if H == 3 {
			return "biome.SnowyTaiga{}"
		} else {
			return "biome.Taiga{}"
		}
	} else if T == 1 {
		if H <= 1 {
			return "biome.Plains{}"
		} else if H == 2 {
			return "biome.Forest{}"
		} else if H == 3 {
			return "biome.Taiga{}"
		} else {
			if W == 1 {
				return "biome.OldGrowthSpruceTaiga{}"
			} else {
				return "biome.OldGrowthPineTaiga{}"
			}
		}
	} else if T == 2 {
		if H == 0 {
			if W == 1 {
				return "biome.FlowerForest{}"
			} else {
				return "biome.Plains{}"
			}
		} else if H == 1 {
			return "biome.Plains{}"
		} else if H == 2 {
			return "biome.Forest{}"
		} else if H == 3 {
			if W == 1 {
				return "biome.BirchForest{}"
			} else {
				return "biome.OldGrowthBirchForest{}"
			}
		} else {
			return "biome.DarkForest{}"
		}
	} else if T == 3 {
		if H <= 1 {
			return "biome.Savanna{}"
		} else if H == 2 {
			if W == 1 {
				return "biome.Forest{}"
			} else {
				return "biome.Plains{}"
			}
		} else if H == 3 {
			if W == 1 {
				return "biome.Jungle{}"
			} else {
				return "biome.JungleEdge{}"
			}
		} else {
			if W == 1 {
				return "biome.Jungle{}"
			} else {
				return "biome.BambooJungle{}"
			}
		}
	} else {
		return "biome.Desert{}"
	}

badland_biomes:
	if H <= 1 {
		if W == 1 {
			return "biome.Badlands{}"
		} else {
			return "biome.ErodedBadlands{}"
		}
	} else if H == 2 {
		return "biome.Badlands{}"
	} else {
		return "biome.WoodedBadlandsPlateau{}"
	}

beach_biomes:
	if T == 0 {
		return "biome.SnowyBeach{}"
	} else if T <= 3 {
		return "biome.Beach{}"
	}else{
		return "biome.Desert{}"
	}

plateau_biomes:
	if T == 0{
		if H == 0{
			if W == 1{
				return "biome.SnowyPlains{}"
			}else{
				return "biome.IceSpikes{}"
			}
		}else if H <= 2{
			return "biome.SnowyPlains{}"
		}else{
			return "biome.SnowyTaiga{}"
		}
	}else if T == 1{
		if H == 0{
			if W == 1{
				return "biome.Meadow{}"
			}else{
				return "biome.CherryGrove{}"
			}
		}else if H == 1{
			return "biome.Meadow{}"
		}else if H == 2{
			if W == 1{
				return "biome.Forest{}"
			}else{
				return "biome.Meadow{}"
			}
		}else if H == 3{
			if W == 1{
				return "biome.Taiga{}"
			}else{
				return "biome.Meadow{}"
			}
		}else{
			if W == 1{
				return "biome.OldGrowthSpruceTaiga{}"
			}else{
				return "biome.OldGrowthPineTaiga{}"
			}
		}
	}else if T == 2{
		if H <= 1{
			if W == 1{
				return "biome.Meadow{}"
			}else{
				return "biome.CherryGrove{}"
			}
		}else if H == 2{
			if W == 1{
				return "biome.Meadow{}"
			}else{
				return "biome.Forest{}"
			}
		}else if H == 3{
			if W == 1{
				return "biome.Meadow{}"
			}else{
				return "biome.BirchForest{}"
			}
		}else{
			return "biome.DarkForest{}"
		}
	}else if T == 3{
		if H <= 1{
			return "biome.SavannaPlateau{}"
		}else if H <= 3{
			return "biome.Forest{}"
		}else{
			return "biome.Jungle{}"
		}
	}else{
		goto badland_biomes
	}

shattered_biomes:
	if T <= 2{
		if H >= 3{
			return "biome.WindsweptForest{}"
		}else if T <= 1 && H <= 1{
			return "biome.WindsweptGravellyHills{}"
		}else{
			return "biome.WindsweptHills{}"
		}
	}else if T == 3{
		if H <= 1{
			return "biome.Savanna{}"
		} else if H == 2 {
			if W == 1 {
				return "biome.Forest{}"
			} else {
				return "biome.Plains{}"
			}
		} else if H == 3 {
			if W == 1 {
				return "biome.Jungle{}"
			} else {
				return "biome.JungleEdge{}"
			}
		} else {
			if W == 1 {
				return "biome.Jungle{}"
			} else {
				return "biome.BambooJungle{}"
			}
		}
	}else{
		return "biome.Desert{}"
	}
}
