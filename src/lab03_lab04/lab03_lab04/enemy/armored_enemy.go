package enemy

import (
	"lab03_lab04/defense"
)

type ArmoredEnemy struct{ EnemyDecorator }

func NewArmoredEnemy(wrapped Enemy) *ArmoredEnemy { return &ArmoredEnemy{NewEnemyDecorator(wrapped)} }
func (e *ArmoredEnemy) Name() string              { return "Бронированный " + e.Enemy.Name() }
func (e *ArmoredEnemy) ReceiveHit(amount int) defense.Damage {
	if amount <= 0 {
		return defense.Damage{}
	}
	return e.Enemy.ReceiveHit(amount - amount/2)
}
