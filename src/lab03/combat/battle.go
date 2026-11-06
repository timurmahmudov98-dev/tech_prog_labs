package combat

// Battle - бой идет по раундам: сначала бьет герой, потом враг, пока кто-то не погибнет
type Battle struct {
	hero   Combatant
	enemy  Combatant
	events []BattleEvent
	round  int
}

func NewBattle(hero, enemy Combatant) *Battle {
	return &Battle{hero: hero, enemy: enemy}
}

func (b *Battle) Run() BattleReport {
	for b.hero.IsAlive() && b.enemy.IsAlive() {
		b.round++
		b.strike(b.hero, SideHero, b.enemy, SideEnemy)
		if b.enemy.IsAlive() {
			b.strike(b.enemy, SideEnemy, b.hero, SideHero)
		}
	}

	outcome := Defeat
	if b.hero.IsAlive() {
		outcome = Victory
	}

	return BattleReport{
		Hero:    ParticipantOf(b.hero),
		Enemy:   ParticipantOf(b.enemy),
		Outcome: outcome,
		Events:  b.events,
	}
}

func (b *Battle) strike(attacker Combatant, source Side, defender Combatant, target Side) {
	hit := attacker.RollDamage()
	dealt := defender.ReceiveHit(hit.Damage)
	eventType := EventHit
	if hit.Critical {
		eventType = EventCriticalHit
	}
	b.emit(BattleEvent{b.round, eventType, source, target, dealt, defender.Health()})

	if !defender.IsAlive() {
		b.emit(BattleEvent{b.round, EventDeath, source, target, 0, 0})
	}
}

func (b *Battle) emit(event BattleEvent) {
	b.events = append(b.events, event)
}
