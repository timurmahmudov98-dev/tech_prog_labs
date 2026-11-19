package enemy

import (
	"lab03_lab04/combat"
	"lab03_lab04/core"
	"lab03_lab04/defense"
	"lab03_lab04/item"
)

type HauntedWeapon struct {
	weapon            item.Weapon
	health, maxHealth int
}

func NewHauntedWeapon(weapon item.Weapon) *HauntedWeapon {
	hp := weapon.Damage() * 5
	return &HauntedWeapon{weapon, hp, hp}
}
func (e *HauntedWeapon) Code() string           { return "haunted-weapon" }
func (e *HauntedWeapon) Name() string           { return "Одержимый " + e.weapon.Name() }
func (e *HauntedWeapon) Health() int            { return e.health }
func (e *HauntedWeapon) MaxHealth() int         { return e.maxHealth }
func (e *HauntedWeapon) IsAlive() bool          { return e.health > 0 }
func (e *HauntedWeapon) RollDamage() combat.Hit { return e.weapon.RollDamage() }
func (e *HauntedWeapon) ExperienceReward() int  { return e.weapon.Price() / 2 }
func (e *HauntedWeapon) RollGoldReward() int {
	return core.GetInstance().Roll(e.weapon.Price()/4, e.weapon.Price()/2)
}
func (e *HauntedWeapon) ReceiveHit(amount int) defense.Damage {
	e.health = max(0, e.health-amount)
	return defense.Damage{Amount: amount}
}
