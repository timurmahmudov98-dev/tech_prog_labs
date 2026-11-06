package enemy

// Skeleton - кости плохо поддаются ударам, поэтому часть урона проходит мимо
type Skeleton struct {
	baseEnemy
}

func NewSkeleton() *Skeleton {
	return &Skeleton{newBaseEnemy("skeleton", "Скелет", 110, 10, 24, 0.3, 90, 70)}
}
