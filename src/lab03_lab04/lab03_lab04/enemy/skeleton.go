package enemy

import (
	"lab03_lab04/defense"
)

// Skeleton - кости плохо поддаются ударам, поэтому часть урона проходит мимо
type Skeleton struct {
	baseEnemy
}

func NewSkeleton() *Skeleton {
	return &Skeleton{newBaseEnemy("skeleton", "Скелет", 110, 10, 24, defense.NewArmorHandler(0.3), 90, 70)}
}
