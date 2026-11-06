package combat

type EventType int

const (
	EventHit EventType = iota
	EventCriticalHit
	EventDeath
)
