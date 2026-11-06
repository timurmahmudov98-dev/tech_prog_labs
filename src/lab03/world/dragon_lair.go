package world

import (
	"lab03/enemy"
	"lab03/item"
)

const DragonLairCode = "dragon-lair"

type DragonLair struct {
	baseLocation
}

func NewDragonLair(catalog *item.ItemCatalog) *DragonLair {
	return &DragonLair{baseLocation{DragonLairCode, "Логово дракона", 5, catalog}}
}

func (d *DragonLair) CreateEnemy() enemy.Enemy {
	return enemy.NewDragon()
}

func (d *DragonLair) LootTable() []string {
	return []string{item.DawnBlade, item.DragonscaleArmor}
}
