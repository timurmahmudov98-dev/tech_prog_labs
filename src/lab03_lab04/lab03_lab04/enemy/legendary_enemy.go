package enemy

import (
	"lab03_lab04/combat"
)

type LegendaryEnemy struct{ EnemyDecorator }

func NewLegendaryEnemy(wrapped Enemy) *LegendaryEnemy {
	return &LegendaryEnemy{NewEnemyDecorator(wrapped)}
}
func (e *LegendaryEnemy) Name() string           { return "Легендарный " + e.Enemy.Name() }
func (e *LegendaryEnemy) RollDamage() combat.Hit { h := e.Enemy.RollDamage(); h.Damage *= 2; return h }
func (e *LegendaryEnemy) ExperienceReward() int  { return e.Enemy.ExperienceReward() * 2 }
func (e *LegendaryEnemy) RollGoldReward() int    { return e.Enemy.RollGoldReward() * 2 }
