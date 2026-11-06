package item

import (
	"fmt"

	"lab03/core"
)

// ItemCatalog - каталог хранит шаблоны предметов, а наружу отдает только их копии
type ItemCatalog struct {
	// шаблоны хранятся в порядке добавления, предметов немного, поэтому ищем простым перебором
	prototypes []Item
}

func (c *ItemCatalog) Add(prototype Item) {
	c.prototypes = append(c.prototypes, prototype)
}

// Create - каталог не знает, какой конкретно тип у шаблона: копию создает сам шаблон
func (c *ItemCatalog) Create(code string) (Item, error) {
	for _, prototype := range c.prototypes {
		if prototype.Code() == code {
			return prototype.Clone(), nil
		}
	}

	return nil, &core.UnknownCodeError{What: "Предмет", Code: code}
}

func (c *ItemCatalog) CreateWeapon(code string) (Weapon, error) {
	item, err := c.Create(code)
	if err != nil {
		return nil, err
	}

	// проверка типа: является ли предмет оружием
	weapon, ok := item.(Weapon)
	if !ok {
		return nil, fmt.Errorf("предмет '%s' не является оружием", code)
	}

	return weapon, nil
}

func (c *ItemCatalog) CreateArmor(code string) (*Armor, error) {
	item, err := c.Create(code)
	if err != nil {
		return nil, err
	}

	armor, ok := item.(*Armor)
	if !ok {
		return nil, fmt.Errorf("предмет '%s' не является броней", code)
	}

	return armor, nil
}

func (c *ItemCatalog) Codes() []string {
	codes := make([]string, 0, len(c.prototypes))
	for _, prototype := range c.prototypes {
		codes = append(codes, prototype.Code())
	}

	return codes
}
