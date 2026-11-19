package tactic

type CautiousTactic struct{}

func (CautiousTactic) Code() string        { return "cautious" }
func (CautiousTactic) DisplayName() string { return "Осторожность" }
func (CautiousTactic) Description() string {
	return "при здоровье ниже 40% через раунд уходить в оборону"
}
func (CautiousTactic) ChooseAction(round int, hero, enemy Combatant) Action {
	if float64(hero.Health()) < float64(hero.MaxHealth())*0.4 && round%2 == 0 {
		return Defend
	}
	return Attack
}
