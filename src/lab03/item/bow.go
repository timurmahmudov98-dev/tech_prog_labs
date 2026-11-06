package item

import (
	"fmt"
	"math"

	"lab03/combat"
	"lab03/core"
)

const criticalMultiplier = 2

// Bow - лук с некоторым шансом наносит двойной урон
type Bow struct {
	baseWeapon
	criticalChance float64
}

func NewBow(code, name string, basePrice, damage int, criticalChance float64) *Bow {
	return &Bow{baseWeapon{baseItem{code: code, name: name, basePrice: basePrice}, damage}, criticalChance}
}

func (b *Bow) RollDamage() combat.Hit {
	if core.GetInstance().Chance(b.criticalChance) {
		return combat.Hit{Damage: b.Damage() * criticalMultiplier, Critical: true}
	}

	return combat.Hit{Damage: b.Damage(), Critical: false}
}

func (b *Bow) Clone() Item {
	// копируем всю структуру, а общую часть оружия - отдельно, вместе с ее слайсом зачарований
	clone := *b
	clone.baseWeapon = b.cloneWeapon()
	return &clone
}

func (b *Bow) Description() string {
	return fmt.Sprintf("урон %d, шанс крита %d%%", b.Damage(), int(math.Round(b.criticalChance*100)))
}
