package world

import "lab03/combat"

// ExpeditionReport - отчет об экспедиции: с каким сидом, куда, как прошел бой и что получил герой
type ExpeditionReport struct {
	Seed         uint64
	LocationCode string
	Battle       combat.BattleReport
	Rewards      Rewards
}
