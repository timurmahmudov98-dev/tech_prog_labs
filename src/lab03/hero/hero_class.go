package hero

import (
	"lab03/core"
	"lab03/item"
	"lab03/kit"
)


type HeroClass struct {
	code        string
	displayName string
	baseHealth  int
	// функция, которая создает фабрику стартового снаряжения
	starterKit func(catalog *item.ItemCatalog) kit.StarterKitFactory
}

var (
	Warrior = &HeroClass{"warrior", "Воин", 120, kit.NewWarriorKitFactory}
	Ranger  = &HeroClass{"ranger", "Следопыт", 100, kit.NewRangerKitFactory}
	Mage    = &HeroClass{"mage", "Маг", 85, kit.NewMageKitFactory}
)


func HeroClasses() []*HeroClass {
	return []*HeroClass{Warrior, Ranger, Mage}
}

func HeroClassByCode(code string) (*HeroClass, error) {
	for _, heroClass := range HeroClasses() {
		if heroClass.code == code {
			return heroClass, nil
		}
	}

	return nil, &core.UnknownCodeError{What: "Класс героя", Code: code}
}


func (c *HeroClass) StarterKit(catalog *item.ItemCatalog) kit.StarterKitFactory {
	return c.starterKit(catalog)
}

func (c *HeroClass) Code() string { return c.code }

func (c *HeroClass) DisplayName() string { return c.displayName }

func (c *HeroClass) BaseHealth() int { return c.baseHealth }
