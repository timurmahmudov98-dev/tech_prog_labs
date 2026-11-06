package world

import (
	"lab03/enemy"
	"lab03/item"
)

const ForestCode = "forest"

type Forest struct {
	baseLocation
}

func NewForest(catalog *item.ItemCatalog) *Forest {
	return &Forest{baseLocation{ForestCode, "Мистический лес", 1, catalog}}
}

func (f *Forest) CreateEnemy() enemy.Enemy {
	return enemy.NewGoblin()
}

func (f *Forest) LootTable() []string {
	return []string{item.GoblinCleaver, item.HuntingBow, item.LeatherArmor}
}
