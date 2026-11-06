package console

import (
	"fmt"
	"strings"

	"lab03/combat"
	"lab03/hero"
	"lab03/item"
	"lab03/world"
)

// ConsoleView - весь вывод в консоль собран здесь
type ConsoleView struct{}

func (v ConsoleView) Describe(it item.Item) string {
	description := it.Name()

	enchantments := it.Enchantments()
	if len(enchantments) > 0 {
		names := make([]string, 0, len(enchantments))
		for _, enchantment := range enchantments {
			names = append(names, enchantment.Name())
		}
		description += " [" + strings.Join(names, ", ") + "]"
	}

	return description + " (" + it.Description() + ")"
}

func (v ConsoleView) PrintHero(h *hero.Hero) {
	fmt.Printf(
		"%s, %s %d ур. | здоровье %d | опыт %d/%d | золото %d\n",
		h.Name(), h.HeroClass().DisplayName(), h.Level(),
		h.MaxHealth(), h.Experience(), h.ExperienceToNextLevel(), h.Gold(),
	)
	fmt.Println("  Оружие: " + v.Describe(h.Weapon()))
	fmt.Println("  Броня:  " + v.Describe(h.Armor()))
}

func (v ConsoleView) PrintReport(report world.ExpeditionReport, locationName string) {
	battle := report.Battle
	fmt.Printf("=== %s (сид %d) ===\n", locationName, report.Seed)
	fmt.Printf("%s встречает врага: %s (%d здоровья)\n", battle.Hero.Name, battle.Enemy.Name, battle.Enemy.MaxHealth)

	for _, event := range battle.Events {
		v.printEvent(event, battle)
	}

	v.printRewards(battle, report.Rewards)
}

func (v ConsoleView) printEvent(event combat.BattleEvent, battle combat.BattleReport) {
	source := participant(event.Source, battle)
	target := participant(event.Target, battle)

	switch event.Type {
	case combat.EventHit:
		fmt.Printf(
			"[раунд %d] %s -> %s: %d урона, осталось %d/%d\n",
			event.Round, source.Name, target.Name, event.Amount, event.TargetHealthAfter, target.MaxHealth,
		)
	case combat.EventCriticalHit:
		fmt.Printf(
			"[раунд %d] %s -> %s: КРИТ! %d урона, осталось %d/%d\n",
			event.Round, source.Name, target.Name, event.Amount, event.TargetHealthAfter, target.MaxHealth,
		)
	case combat.EventDeath:
		fmt.Printf("%s повержен!\n", target.Name)
	}
}

func (v ConsoleView) printRewards(battle combat.BattleReport, rewards world.Rewards) {
	if battle.Outcome == combat.Defeat {
		fmt.Printf("Поражение... %s едва уносит ноги и теряет %d золота\n", battle.Hero.Name, -rewards.Gold)
		return
	}

	fmt.Printf("Победа! +%d опыта, +%d золота\n", rewards.Experience, rewards.Gold)
	if rewards.Loot != nil {
		fmt.Println("Добыча: " + v.Describe(rewards.Loot))
	}
	if rewards.LevelsGained > 0 {
		fmt.Println("Новый уровень!")
	}
}

func participant(side combat.Side, battle combat.BattleReport) combat.Participant {
	if side == combat.SideHero {
		return battle.Hero
	}

	return battle.Enemy
}
