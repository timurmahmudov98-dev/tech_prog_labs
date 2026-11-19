package tactic

type BalancedTactic struct{}

func (BalancedTactic) Code() string                                         { return "balanced" }
func (BalancedTactic) DisplayName() string                                  { return "Баланс" }
func (BalancedTactic) Description() string                                  { return "обычный удар каждый раунд" }
func (BalancedTactic) ChooseAction(round int, hero, enemy Combatant) Action { return Attack }
