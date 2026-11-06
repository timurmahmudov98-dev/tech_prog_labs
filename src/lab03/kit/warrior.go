package kit

import "lab03/item"

// WarriorKitFactory - воин начинает с мечом и кольчугой
type WarriorKitFactory struct {
	catalog *item.ItemCatalog
}

func NewWarriorKitFactory(catalog *item.ItemCatalog) StarterKitFactory {
	return &WarriorKitFactory{catalog: catalog}
}

func (f *WarriorKitFactory) CreateWeapon() (item.Weapon, error) {
	return f.catalog.CreateWeapon(item.TrainingSword)
}

func (f *WarriorKitFactory) CreateArmor() (*item.Armor, error) {
	return f.catalog.CreateArmor(item.Chainmail)
}
