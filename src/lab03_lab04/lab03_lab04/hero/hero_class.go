package hero

import (
	"lab03_lab04/defense"

	"lab03_lab04/core"
	"lab03_lab04/item"
	"lab03_lab04/kit"
)

// HeroClass - класс героя. Поля неэкспортируемые, поэтому классы героя - это только значения,
// заранее созданные в этом пакете (Warrior, Ranger, Mage)
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

// HeroClasses возвращает все классы героя по порядку
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

// StarterKit - каждый класс сам знает свою фабрику стартового снаряжения
func (c *HeroClass) StarterKit(catalog *item.ItemCatalog) kit.StarterKitFactory {
	return c.starterKit(catalog)
}

func (c *HeroClass) Code() string { return c.code }

func (c *HeroClass) DisplayName() string { return c.displayName }

func (c *HeroClass) BaseHealth() int { return c.baseHealth }

func (c *HeroClass) CreateDefense(armor *item.Armor, level int) defense.DamageHandler {
	evasion, barrier := 0.0, 0
	if c == Ranger {
		evasion = 0.15
	}
	if c == Mage {
		barrier = 30
	}
	if c == Warrior && level >= 3 {
		barrier = 20
	} else if c == Ranger && level >= 5 {
		barrier = 15
	} else if c == Mage && level >= 3 {
		evasion = 0.10
	}
	armorHandler := defense.NewArmorHandler(armor.Defense())
	var first defense.DamageHandler = armorHandler
	if evasion > 0 {
		first = defense.NewEvasionHandler(evasion)
		first.SetNext(armorHandler)
	}
	if barrier > 0 {
		armorHandler.SetNext(defense.NewBarrierHandler(barrier))
	}
	return first
}
