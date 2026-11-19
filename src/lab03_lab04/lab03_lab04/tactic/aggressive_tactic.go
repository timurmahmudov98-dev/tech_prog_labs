package tactic

type AggressiveTactic struct{}

func (AggressiveTactic) Code() string        { return "aggressive" }
func (AggressiveTactic) DisplayName() string { return "Натиск" }
func (AggressiveTactic) Description() string {
	return "совершать безрассудные удары, пока у врага больше здоровья, чем у героя"
}
func (AggressiveTactic) ChooseAction(round int, hero, enemy Combatant) Action {
	if enemy.Health() > hero.Health() {
		return PowerStrike
	}
	return Attack
}
