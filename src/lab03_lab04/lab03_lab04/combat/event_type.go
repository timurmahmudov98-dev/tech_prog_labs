package combat

type EventType int

const (
	EventHit EventType = iota
	EventCriticalHit
	EventDeath
)

const (
	EventPowerStrike EventType = iota + 3
	EventDefend
	EventEvade
	EventAbsorb
)
