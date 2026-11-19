package defense

import (
	"math"
)

type ArmorHandler struct {
	baseHandler
	defense float64
}

func NewArmorHandler(defense float64) *ArmorHandler { return &ArmorHandler{defense: defense} }
func (h *ArmorHandler) Handle(d *Damage) {
	d.ReduceTo(max(1, int(math.Round(float64(d.Amount)*(1-h.defense)))))
	h.passToNext(d)
}
