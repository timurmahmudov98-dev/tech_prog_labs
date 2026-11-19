package world

import "lab03_lab04/combat"

// ExpeditionReport - отчет об экспедиции: с каким сидом, куда, как прошел бой и что получил герой
type ExpeditionReport struct {
	Seed         uint64
	LocationCode string
	TacticCode   string
	Battle       combat.BattleReport
	Rewards      Rewards
}
