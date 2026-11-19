package event

type GameEventListener interface{ OnEvent(GameEvent) }
type ListenerFunc func(GameEvent)

func (f ListenerFunc) OnEvent(e GameEvent) { f(e) }
