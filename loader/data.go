package loader

import (
	"fmt"
	"strings"
)

type Rarity string
type Table string

const (
	Common    Rarity = "Common"
	Uncommon  Rarity = "Uncommon"
	Rare      Rarity = "Rare"
	VeryRare  Rarity = "Very Rare"
	Legendary Rarity = "Legendary"
	Artifact  Rarity = "Artifact"
)

const (
	Arcana     Table = "Arcana"
	Armaments  Table = "Armaments"
	Implements Table = "Implements"
	Relics     Table = "Relics"
)

type MagicItem struct {
	Name        string `json:"name"`
	Rarity      Rarity `json:"rarity"`
	CostGP      int    `json:"cost_gp"`
	PageNumber  int    `json:"page_number"`
	Table       Table  `json:"table"`
	TableNumber []int  `json:"table_number"`
}

func ParseRarity(s string) (Rarity, error) {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case "common":
		return Common, nil
	case "uncommon":
		return Uncommon, nil
	case "rare":
		return Rare, nil
	case "very rare":
		return VeryRare, nil
	case "legendary":
		return Legendary, nil
	case "artifact":
		return Artifact, nil
	default:
		return "", fmt.Errorf("invalid rarity: %v", s)
	}
}

func ParseTable(s string) (Table, error) {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case "arcana":
		return Arcana, nil
	case "armaments":
		return Armaments, nil
	case "implements":
		return Implements, nil
	case "relics":
		return Relics, nil
	default:
		return "", fmt.Errorf("invaid table: %v", s)
	}
}

func DefaultCostForRarity(r Rarity) int {
	switch r {
	case Common:
		return 100
	case Uncommon:
		return 400
	case Rare:
		return 4000
	case VeryRare:
		return 40000
	case Legendary:
		return 200000
	case Artifact:
		return 1
	default:
		return 1
	}
}
