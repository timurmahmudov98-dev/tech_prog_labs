package enemy

import (
	"math"

	"lab03/combat"
	"lab03/core"
)

// baseEnemy - общая логика для врагов: здоровье, урон в диапазоне, сопротивление и награда.
// Конкретные враги встраивают эту структуру в себя
type baseEnemy struct {
	code      string
	name      string
	maxHealth int
	minDamage int
	maxDamage int
	// доля урона, которую враг игнорирует
	resistance       float64
	experienceReward int
	goldReward       int
	health           int
}

func newBaseEnemy(code, name string, maxHealth, minDamage, maxDamage int, resistance float64, experienceReward, goldReward int) baseEnemy {
	return baseEnemy{code, name, maxHealth, minDamage, maxDamage, resistance, experienceReward, goldReward, maxHealth}
}

func (e *baseEnemy) RollDamage() combat.Hit {
	return combat.Hit{Damage: core.GetInstance().Roll(e.minDamage, e.maxDamage), Critical: false}
}

func (e *baseEnemy) ReceiveHit(damage int) int {
	reducedDamage := int(math.Round(float64(damage) * (1 - e.resistance)))
	dealt := max(1, reducedDamage)
	e.health = max(0, e.health-dealt)
	return dealt
}

func (e *baseEnemy) IsAlive() bool { return e.health > 0 }

func (e *baseEnemy) RollGoldReward() int {
	return core.GetInstance().Roll(e.goldReward/2, e.goldReward)
}

func (e *baseEnemy) ExperienceReward() int { return e.experienceReward }

func (e *baseEnemy) Code() string { return e.code }

func (e *baseEnemy) Name() string { return e.name }

func (e *baseEnemy) Health() int { return e.health }

func (e *baseEnemy) MaxHealth() int { return e.maxHealth }
