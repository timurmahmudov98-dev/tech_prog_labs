package combat

// BattleReport - результат боя: кто сражался, чем все закончилось и что происходило по ходу дела
type BattleReport struct {
	Hero    Participant
	Enemy   Participant
	Outcome Outcome
	Events  []BattleEvent
}
