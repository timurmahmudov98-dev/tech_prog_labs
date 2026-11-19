package enemy

import (
	"lab03_lab04/combat"
)

type WindfuryEnemy struct{ EnemyDecorator }

func NewWindfuryEnemy(wrapped Enemy) *WindfuryEnemy {
	return &WindfuryEnemy{NewEnemyDecorator(wrapped)}
}
func (e *WindfuryEnemy) Name() string { return "Неистовый " + e.Enemy.Name() }
func (e *WindfuryEnemy) RollDamage() combat.Hit {
	first, second := e.Enemy.RollDamage(), e.Enemy.RollDamage()
	return combat.Hit{Damage: first.Damage + second.Damage, Critical: first.Critical || second.Critical}
}
func (e *WindfuryEnemy) ExperienceReward() int { r := e.Enemy.ExperienceReward(); return r + r*50/100 }
func (e *WindfuryEnemy) RollGoldReward() int   { r := e.Enemy.RollGoldReward(); return r + r*50/100 }
