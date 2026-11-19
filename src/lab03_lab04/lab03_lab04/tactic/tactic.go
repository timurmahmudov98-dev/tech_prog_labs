package tactic

// Для выбора действия достаточно здоровья участников: пакет не зависит от реализации боя.
type Combatant interface {
	Health() int
	MaxHealth() int
}
type Tactic interface {
	Code() string
	DisplayName() string
	Description() string
	ChooseAction(round int, hero, enemy Combatant) Action
}
