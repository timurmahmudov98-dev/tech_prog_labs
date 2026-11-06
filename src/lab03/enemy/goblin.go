package enemy

// Goblin - самый простой враг, без особенностей
type Goblin struct {
	baseEnemy
}

func NewGoblin() *Goblin {
	return &Goblin{newBaseEnemy("goblin", "Гоблин", 50, 6, 14, 0, 40, 30)}
}
