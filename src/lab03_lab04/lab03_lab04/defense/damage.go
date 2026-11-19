package defense

type Damage struct {
	Amount   int
	Evaded   bool
	Absorbed int
}

func (d *Damage) ReduceTo(amount int) { d.Amount = amount }
func (d *Damage) Evade()              { d.Amount = 0; d.Evaded = true }
func (d *Damage) Absorb(value int)    { d.Amount -= value; d.Absorbed += value }
