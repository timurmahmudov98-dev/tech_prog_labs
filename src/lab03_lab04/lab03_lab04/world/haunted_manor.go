package world

import (
	"lab03_lab04/core"
	"lab03_lab04/enemy"
	"lab03_lab04/item"
)

const HauntedManorCode = "haunted-manor"

type HauntedManor struct{ baseLocation }

func NewHauntedManor(catalog *item.ItemCatalog) *HauntedManor {
	return &HauntedManor{baseLocation{HauntedManorCode, "Проклятый особняк", 2, catalog}}
}
func (l *HauntedManor) LootTable() []string {
	return []string{item.IronSword, item.ElvenBow, item.EmberStaff}
}
func (l *HauntedManor) CreateEnemy() enemy.Enemy {
	table := l.LootTable()
	weapon, err := l.catalog.CreateWeapon(table[core.GetInstance().Index(len(table))])
	if err != nil {
		panic(err)
	}
	return enemy.NewHauntedWeapon(weapon)
}
