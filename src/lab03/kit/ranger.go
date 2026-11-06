package kit

import "lab03/item"

// RangerKitFactory - следопыт начинает с луком и кожаной броней
type RangerKitFactory struct {
	catalog *item.ItemCatalog
}

func NewRangerKitFactory(catalog *item.ItemCatalog) StarterKitFactory {
	return &RangerKitFactory{catalog: catalog}
}

func (f *RangerKitFactory) CreateWeapon() (item.Weapon, error) {
	return f.catalog.CreateWeapon(item.HuntingBow)
}

func (f *RangerKitFactory) CreateArmor() (*item.Armor, error) {
	return f.catalog.CreateArmor(item.LeatherArmor)
}
