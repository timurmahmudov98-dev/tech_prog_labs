package event

type GameEventType string

const (
	Victory GameEventType = "victory"
	Defeat  GameEventType = "defeat"
	LevelUp GameEventType = "level-up"
)
