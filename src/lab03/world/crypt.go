package world

import (
	"lab03/enemy"
	"lab03/item"
)

const CryptCode = "crypt"

type Crypt struct {
	baseLocation
}

func NewCrypt(catalog *item.ItemCatalog) *Crypt {
	return &Crypt{baseLocation{CryptCode, "Заброшенный склеп", 3, catalog}}
}

func (c *Crypt) CreateEnemy() enemy.Enemy {
	return enemy.NewSkeleton()
}

func (c *Crypt) LootTable() []string {
	return []string{item.KnightBlade, item.BoneBow, item.NecromancerStaff, item.CryptMail}
}
