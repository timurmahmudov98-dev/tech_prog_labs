package defense

import (
	"lab03_lab04/core"
)

type EvasionHandler struct {
	baseHandler
	chance float64
}

func NewEvasionHandler(chance float64) *EvasionHandler { return &EvasionHandler{chance: chance} }
func (h *EvasionHandler) Handle(d *Damage) {
	if core.GetInstance().Chance(h.chance) {
		d.Evade()
		return
	}
	h.passToNext(d)
}
