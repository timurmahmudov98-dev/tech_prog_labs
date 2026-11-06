package item

import (
	"fmt"
	"math"

	"lab03/combat"
	"lab03/core"
)

// Staff - урон посоха непредсказуем: он может быть как больше, так и меньше базового
type Staff struct {
	baseWeapon
	scatter float64
}

func NewStaff(code, name string, basePrice, damage int, scatter float64) *Staff {
	return &Staff{baseWeapon{baseItem{code: code, name: name, basePrice: basePrice}, damage}, scatter}
}

func (s *Staff) RollDamage() combat.Hit {
	return combat.Hit{Damage: core.GetInstance().Roll(s.minDamage(), s.maxDamage()), Critical: false}
}

func (s *Staff) Clone() Item {
	clone := *s
	clone.baseWeapon = s.cloneWeapon()
	return &clone
}

func (s *Staff) Description() string {
	return fmt.Sprintf("урон %d-%d", s.minDamage(), s.maxDamage())
}

func (s *Staff) minDamage() int {
	return int(math.Round(float64(s.Damage()) * (1 - s.scatter)))
}

func (s *Staff) maxDamage() int {
	return int(math.Round(float64(s.Damage()) * (1 + s.scatter)))
}
