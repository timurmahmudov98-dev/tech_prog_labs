package event

import (
	"lab03_lab04/progress"
)

type AchievementTracker struct{}

func (*AchievementTracker) OnEvent(e GameEvent) {
	if e.Type == Victory {
		e.Hero.Unlock(progress.FirstBlood)
		if e.Expedition.Battle.Enemy.Code == "dragon" {
			e.Hero.Unlock(progress.DragonSlayer)
		}
		if loot := e.Expedition.Rewards.Loot; loot != nil && len(loot.Enchantments()) > 0 {
			e.Hero.Unlock(progress.Lucky)
		}
	} else if e.Type == LevelUp && e.Hero.Level() >= 5 {
		e.Hero.Unlock(progress.Veteran)
	}
}
