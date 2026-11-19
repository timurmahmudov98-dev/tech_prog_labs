package enemy

import (
	"lab03_lab04/defense"

	"lab03_lab04/combat"
	"lab03_lab04/core"
)

// baseEnemy - общая логика для врагов: здоровье, урон в диапазоне, защита и награда.
// Конкретные враги встраивают эту структуру в себя
type baseEnemy struct {
	code      string
	name      string
	maxHealth int
	minDamage int
	maxDamage int
	// цепочка защиты
	defense          defense.DamageHandler
	experienceReward int
	goldReward       int
	health           int
}

func newBaseEnemy(code, name string, maxHealth, minDamage, maxDamage int, protection defense.DamageHandler, experienceReward, goldReward int) baseEnemy {
	return baseEnemy{code, name, maxHealth, minDamage, maxDamage, protection, experienceReward, goldReward, maxHealth}
}

func (e *baseEnemy) RollDamage() combat.Hit {
	return combat.Hit{Damage: core.GetInstance().Roll(e.minDamage, e.maxDamage), Critical: false}
}

func (e *baseEnemy) ReceiveHit(amount int) defense.Damage {
	d := defense.Damage{Amount: amount}
	if e.defense != nil {
		e.defense.Handle(&d)
	}
	e.health = max(0, e.health-d.Amount)
	return d
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
