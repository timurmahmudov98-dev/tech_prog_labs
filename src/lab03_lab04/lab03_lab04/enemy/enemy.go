package enemy

import "lab03_lab04/combat"

type Enemy interface {
	combat.Combatant
	ExperienceReward() int
	RollGoldReward() int
}
