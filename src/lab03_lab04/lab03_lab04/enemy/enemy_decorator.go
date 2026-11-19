package enemy

type EnemyDecorator struct{ Enemy }

func NewEnemyDecorator(wrapped Enemy) EnemyDecorator { return EnemyDecorator{wrapped} }
