package kit

import "lab03/item"

// MageKitFactory - маг начинает с посохом и робой
type MageKitFactory struct {
	catalog *item.ItemCatalog
}

func NewMageKitFactory(catalog *item.ItemCatalog) StarterKitFactory {
	return &MageKitFactory{catalog: catalog}
}

func (f *MageKitFactory) CreateWeapon() (item.Weapon, error) {
	return f.catalog.CreateWeapon(item.ApprenticeStaff)
}

func (f *MageKitFactory) CreateArmor() (*item.Armor, error) {
	return f.catalog.CreateArmor(item.LinenRobe)
}
