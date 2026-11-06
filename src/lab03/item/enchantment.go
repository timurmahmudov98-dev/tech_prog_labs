package item

// Enchantment - зачарование. Все его поля неэкспортируемые и меняются только при создании,
// поэтому зачарование неизменяемо, и одно и то же зачарование можно безопасно повесить на несколько предметов
type Enchantment struct {
	code         string
	name         string
	bonusDamage  int
	bonusDefense float64
	bonusPrice   int
}

var (
	Sharp   = &Enchantment{"sharp", "Острое", 3, 0, 20}
	Fiery   = &Enchantment{"fiery", "Огненное", 6, 0, 50}
	Sturdy  = &Enchantment{"sturdy", "Прочное", 0, 0.05, 20}
	Warding = &Enchantment{"warding", "Охранное", 0, 0.10, 50}

	ForWeapons = []*Enchantment{Sharp, Fiery}
	ForArmor   = []*Enchantment{Sturdy, Warding}
)

func (e *Enchantment) Code() string          { return e.code }
func (e *Enchantment) Name() string          { return e.name }
func (e *Enchantment) BonusDamage() int      { return e.bonusDamage }
func (e *Enchantment) BonusDefense() float64 { return e.bonusDefense }
func (e *Enchantment) BonusPrice() int       { return e.bonusPrice }
