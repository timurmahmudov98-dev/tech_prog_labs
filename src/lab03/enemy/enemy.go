package enemy

import "lab03/combat"

// Enemy - враг: боец, за победу над которым дают награду
type Enemy interface {
	combat.Combatant
	ExperienceReward() int
	RollGoldReward() int
}
