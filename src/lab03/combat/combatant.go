package combat

// Combatant - всё, что может участвовать в бою: и герой, и враги
type Combatant interface {
	// Code - код вида сущности: класс героя или тип врага
	Code() string
	Name() string
	Health() int
	MaxHealth() int
	IsAlive() bool
	// RollDamage - удар, который будет нанесен данной сущностью
	RollDamage() Hit
	// ReceiveHit - получение удара, возвращает урон, который прошел на самом деле
	ReceiveHit(damage int) int
}
