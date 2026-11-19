package defense

type DamageHandler interface {
	Handle(*Damage)
	SetNext(DamageHandler) DamageHandler
}
type baseHandler struct{ next DamageHandler }

func (h *baseHandler) SetNext(next DamageHandler) DamageHandler { h.next = next; return next }
func (h *baseHandler) passToNext(d *Damage) {
	if h.next != nil {
		h.next.Handle(d)
	}
}
