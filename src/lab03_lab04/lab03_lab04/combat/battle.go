package combat

import (
	"lab03_lab04/tactic"
	"math"
)

// Battle - бой идет по раундам: сначала бьет герой, потом враг, пока кто-то не погибнет
type Battle struct {
	hero   Combatant
	enemy  Combatant
	events []BattleEvent
	round  int
	tactic tactic.Tactic
}

func NewBattle(hero, enemy Combatant, strategy tactic.Tactic) *Battle {
	return &Battle{hero: hero, enemy: enemy, tactic: strategy}
}

func (b *Battle) Run() BattleReport {
	for b.hero.IsAlive() && b.enemy.IsAlive() {
		b.round++
		action := b.tactic.ChooseAction(b.round, b.hero, b.enemy)
		if action == tactic.PowerStrike {
			b.emit(BattleEvent{b.round, EventPowerStrike, SideHero, SideEnemy, 0, 0})
		} else if action == tactic.Defend {
			b.emit(BattleEvent{b.round, EventDefend, SideHero, SideHero, 0, 0})
		}
		b.strike(b.hero, SideHero, b.enemy, SideEnemy, action.AttackMultiplier)
		if b.enemy.IsAlive() {
			b.strike(b.enemy, SideEnemy, b.hero, SideHero, action.IncomingMultiplier)
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

func (b *Battle) strike(attacker Combatant, source Side, defender Combatant, target Side, multiplier float64) {
	hit := attacker.RollDamage()
	d := defender.ReceiveHit(max(1, int(math.Round(float64(hit.Damage)*multiplier))))
	if d.Evaded {
		b.emit(BattleEvent{b.round, EventEvade, source, target, 0, defender.Health()})
		return
	}
	if d.Absorbed > 0 {
		b.emit(BattleEvent{b.round, EventAbsorb, source, target, d.Absorbed, defender.Health()})
	}
	if d.Amount > 0 {
		eventType := EventHit
		if hit.Critical {
			eventType = EventCriticalHit
		}
		b.emit(BattleEvent{b.round, eventType, source, target, d.Amount, defender.Health()})
	}

	if !defender.IsAlive() {
		b.emit(BattleEvent{b.round, EventDeath, source, target, 0, 0})
	}
}

func (b *Battle) emit(event BattleEvent) {
	b.events = append(b.events, event)
}
