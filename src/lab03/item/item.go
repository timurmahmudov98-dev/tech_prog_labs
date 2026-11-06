package item

import (
	"fmt"
	"slices"
)

// Item - прототип: любой предмет умеет создавать свою копию
type Item interface {
	// Clone возвращает копию предмета того же конкретного типа
	Clone() Item
	Code() string
	Name() string
	Price() int
	Enchantments() []*Enchantment
	// PossibleEnchantments - какие зачарования подходят этому предмету
	PossibleEnchantments() []*Enchantment
	Enchant(enchantment *Enchantment) error
	// Description - описание характеристик предмета, например "урон 12"
	Description() string
}

// baseItem - общая часть всех предметов, ее встраивают в себя конкретные типы предметов
type baseItem struct {
	// код шаблона, из которого получен предмет
	code         string
	name         string
	basePrice    int
	enchantments []*Enchantment
}

// cloneBase копирует общую часть предмета
func (b baseItem) cloneBase() baseItem {
	// присваивание структуры копирует все ее поля
	clone := b
	// но слайс - это ссылка на общий массив, поэтому для копии создаем свой слайс,
	// иначе копия и оригинал делили бы один и тот же объект списка зачарований
	clone.enchantments = slices.Clone(b.enchantments)
	return clone
}

func (b *baseItem) enchant(possible []*Enchantment, enchantment *Enchantment) error {
	if !slices.Contains(possible, enchantment) {
		return fmt.Errorf("зачарование '%s' не подходит для предмета '%s'", enchantment.Name(), b.name)
	}

	b.enchantments = append(b.enchantments, enchantment)
	return nil
}

func (b *baseItem) Code() string { return b.code }

func (b *baseItem) Name() string { return b.name }

func (b *baseItem) Price() int {
	price := b.basePrice
	for _, enchantment := range b.enchantments {
		price += enchantment.BonusPrice()
	}

	return price
}

func (b *baseItem) Enchantments() []*Enchantment {
	return slices.Clone(b.enchantments)
}
