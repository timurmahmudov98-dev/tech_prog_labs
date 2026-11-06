package enemy

import "lab03/combat"


type Enemy interface {
	combat.Combatant
	ExperienceReward() int
	RollGoldReward() int
}
