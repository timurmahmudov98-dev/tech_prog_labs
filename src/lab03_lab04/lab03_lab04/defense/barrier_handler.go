package defense

type BarrierHandler struct {
	baseHandler
	strength int
}

func NewBarrierHandler(strength int) *BarrierHandler { return &BarrierHandler{strength: strength} }
func (h *BarrierHandler) Handle(d *Damage) {
	absorbed := min(h.strength, d.Amount)
	if absorbed > 0 {
		h.strength -= absorbed
		d.Absorb(absorbed)
	}
	h.passToNext(d)
}
