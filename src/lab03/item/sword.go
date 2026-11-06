package item

import (
	"fmt"

	"lab03/combat"
)

// Sword - меч бьет стабильно
type Sword struct {
	baseWeapon
}

func NewSword(code, name string, basePrice, damage int) *Sword {
	return &Sword{baseWeapon{baseItem{code: code, name: name, basePrice: basePrice}, damage}}
}

func (s *Sword) RollDamage() combat.Hit {
	return combat.Hit{Damage: s.Damage(), Critical: false}
}

func (s *Sword) Clone() Item {
	return &Sword{baseWeapon: s.cloneWeapon()}
}

func (s *Sword) Description() string {
	return fmt.Sprintf("урон %d", s.Damage())
}
