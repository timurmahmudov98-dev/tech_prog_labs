package event

type subscription struct {
	id       int
	listener GameEventListener
}
type GameEventPublisher struct {
	listeners map[GameEventType][]subscription
	nextID    int
}

func (p *GameEventPublisher) Subscribe(t GameEventType, l GameEventListener) int {
	if p.listeners == nil {
		p.listeners = make(map[GameEventType][]subscription)
	}
	p.nextID++
	p.listeners[t] = append(p.listeners[t], subscription{p.nextID, l})
	return p.nextID
}
func (p *GameEventPublisher) Unsubscribe(t GameEventType, id int) {
	for i, s := range p.listeners[t] {
		if s.id == id {
			p.listeners[t] = append(p.listeners[t][:i], p.listeners[t][i+1:]...)
			return
		}
	}
}
func (p *GameEventPublisher) Publish(e GameEvent) {
	for _, s := range append([]subscription(nil), p.listeners[e.Type]...) {
		s.listener.OnEvent(e)
	}
}
