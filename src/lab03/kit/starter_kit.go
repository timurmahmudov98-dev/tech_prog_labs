package kit

import "lab03/item"

// StarterKitFactory - абстрактная фабрика: стартовый набор снаряжения
type StarterKitFactory interface {
	// CreateWeapon возвращает некое оружие
	CreateWeapon() (item.Weapon, error)
	// CreateArmor возвращает некую броню
	CreateArmor() (*item.Armor, error)
}
