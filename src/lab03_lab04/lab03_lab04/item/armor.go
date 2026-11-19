package item

import (
	"fmt"
	"math"
)

// больше этой доли урона никакая броня не поглощает
const maxDefense = 0.75

// Armor - броня поглощает долю входящего урона
type Armor struct {
	baseItem
	defense float64
}

func NewArmor(code, name string, basePrice int, defense float64) *Armor {
	return &Armor{baseItem{code: code, name: name, basePrice: basePrice}, defense}
}

func (a *Armor) Clone() Item {
	clone := *a
	clone.baseItem = a.cloneBase()
	return &clone
}

// Defense - доля поглощаемого урона с учетом зачарований
func (a *Armor) Defense() float64 {
	total := a.defense
	for _, enchantment := range a.enchantments {
		total += enchantment.BonusDefense()
	}

	return math.Min(total, maxDefense)
}

func (a *Armor) PossibleEnchantments() []*Enchantment {
	return ForArmor
}

func (a *Armor) Enchant(enchantment *Enchantment) error {
	return a.enchant(ForArmor, enchantment)
}

func (a *Armor) Description() string {
	return fmt.Sprintf("защита %d%%", int(math.Round(a.Defense()*100)))
}
