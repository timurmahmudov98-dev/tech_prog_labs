package enemy

import (
	"lab03/combat"
	"lab03/core"
)

const (
	fireBreathChance     = 0.25
	fireBreathMultiplier = 2
)

// Dragon - дракон защищен чешуей и время от времени дышит огнем
type Dragon struct {
	baseEnemy
}

func NewDragon() *Dragon {
	return &Dragon{newBaseEnemy("dragon", "Дракон", 260, 16, 22, 0.25, 500, 400)}
}

// RollDamage заменяет метод встроенной структуры: бой вызывает его через интерфейс, поэтому сработает именно этот метод
func (d *Dragon) RollDamage() combat.Hit {
	bite := d.baseEnemy.RollDamage()
	if core.GetInstance().Chance(fireBreathChance) {
		return combat.Hit{Damage: bite.Damage * fireBreathMultiplier, Critical: true}
	}

	return bite
}
