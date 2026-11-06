package item

import "lab03/combat"

// Weapon - оружие: предмет, который умеет наносить урон
type Weapon interface {
	Item
	// RollDamage - как именно оружие наносит урон, решают конкретные типы
	RollDamage() combat.Hit
	Damage() int
}

// baseWeapon - общая часть всего оружия
type baseWeapon struct {
	baseItem
	damage int
}

func (w baseWeapon) cloneWeapon() baseWeapon {
	clone := w
	clone.baseItem = w.cloneBase()
	return clone
}

// Damage - урон с учетом зачарований
func (w *baseWeapon) Damage() int {
	total := w.damage
	for _, enchantment := range w.enchantments {
		total += enchantment.BonusDamage()
	}

	return total
}

func (w *baseWeapon) PossibleEnchantments() []*Enchantment {
	return ForWeapons
}

func (w *baseWeapon) Enchant(enchantment *Enchantment) error {
	return w.enchant(ForWeapons, enchantment)
}
