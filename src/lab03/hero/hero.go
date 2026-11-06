package hero

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"strings"
	"unicode/utf8"

	"lab03/combat"
	"lab03/item"
	"lab03/kit"
)

const (
	maxNameLength      = 20
	startingGold       = 50
	healthPerLevel     = 15
	damagePerLevel     = 3
	experiencePerLevel = 100
)

// Hero - герой. Все поля неэкспортируемые, поэтому вне пакета hero собрать героя можно только через строителя
type Hero struct {
	name       string
	heroClass  *HeroClass
	level      int
	experience int
	gold       int
	health     int
	weapon     item.Weapon
	armor      *item.Armor
	inventory  []item.Item
}

// Builder - строитель, чтобы удобно и безопасно создавать героев
type Builder struct {
	name      string
	heroClass *HeroClass
	weapon    item.Weapon
	armor     *item.Armor
	level     int
	gold      int
	// первая ошибка, случившаяся при заполнении полей; ее вернет Build
	err error
}

// NewBuilder создает строителя: необязательные поля сразу получают значения по умолчанию
func NewBuilder() *Builder {
	return &Builder{level: 1, gold: startingGold}
}

func (b *Builder) SetName(name string) *Builder {
	b.name = name
	return b
}

func (b *Builder) SetHeroClass(heroClass *HeroClass) *Builder {
	b.heroClass = heroClass
	return b
}

func (b *Builder) SetWeapon(weapon item.Weapon) *Builder {
	b.weapon = weapon
	return b
}

func (b *Builder) SetArmor(armor *item.Armor) *Builder {
	b.armor = armor
	return b
}

// SetStarterKit берет оружие и броню из стартового набора
func (b *Builder) SetStarterKit(starterKit kit.StarterKitFactory) *Builder {
	weapon, err := starterKit.CreateWeapon()
	if err != nil {
		b.fail(err)
		return b
	}
	armor, err := starterKit.CreateArmor()
	if err != nil {
		b.fail(err)
		return b
	}

	b.weapon = weapon
	b.armor = armor
	return b
}

func (b *Builder) SetLevel(level int) *Builder {
	b.level = level
	return b
}

func (b *Builder) SetGold(gold int) *Builder {
	b.gold = gold
	return b
}

// Build проверяет, что все обязательные поля заполнены корректно, и только тогда создает героя
func (b *Builder) Build() (*Hero, error) {
	if b.err != nil {
		return nil, b.err
	}
	if strings.TrimSpace(b.name) == "" {
		return nil, errors.New("у героя должно быть имя")
	}
	if utf8.RuneCountInString(b.name) > maxNameLength {
		return nil, fmt.Errorf("имя героя не может быть длиннее %d символов", maxNameLength)
	}
	if b.heroClass == nil {
		return nil, errors.New("не выбран класс героя")
	}
	if b.weapon == nil || b.armor == nil {
		return nil, errors.New("герою нужны оружие и броня")
	}
	if b.level < 1 {
		return nil, errors.New("уровень героя не может быть меньше 1")
	}
	if b.gold < 0 {
		return nil, errors.New("золота не может быть меньше нуля")
	}

	hero := &Hero{
		name:      b.name,
		heroClass: b.heroClass,
		level:     b.level,
		gold:      b.gold,
		weapon:    b.weapon,
		armor:     b.armor,
	}
	hero.health = hero.MaxHealth()
	return hero, nil
}

func (b *Builder) fail(err error) {
	if b.err == nil {
		b.err = err
	}
}

func (h *Hero) RollDamage() combat.Hit {
	hit := h.weapon.RollDamage()
	return combat.Hit{Damage: hit.Damage + (h.level-1)*damagePerLevel, Critical: hit.Critical}
}

func (h *Hero) ReceiveHit(damage int) int {
	reducedDamage := int(math.Round(float64(damage) * (1 - h.armor.Defense())))
	// хотя бы единицу урона броня пропускает всегда
	dealt := max(1, reducedDamage)
	h.health = max(0, h.health-dealt)
	return dealt
}

func (h *Hero) IsAlive() bool {
	return h.health > 0
}

// поведение героя между экспедициями

// RestoreHealth - перед каждой экспедицией герой успевает отдохнуть
func (h *Hero) RestoreHealth() {
	h.health = h.MaxHealth()
}

// GainExperience начисляет опыт и возвращает количество полученных уровней
func (h *Hero) GainExperience(amount int) int {
	h.experience += amount
	levelsGained := 0
	for h.experience >= h.ExperienceToNextLevel() {
		h.experience -= h.ExperienceToNextLevel()
		h.level++
		levelsGained++
	}

	return levelsGained
}

func (h *Hero) AddGold(amount int) {
	h.gold += amount
}

func (h *Hero) SpendGold(amount int) error {
	if amount > h.gold {
		return fmt.Errorf("не хватает золота: нужно %d, есть %d", amount, h.gold)
	}

	h.gold -= amount
	return nil
}

func (h *Hero) AddItem(it item.Item) {
	h.inventory = append(h.inventory, it)
}

func (h *Hero) RemoveItem(it item.Item) error {
	index := slices.Index(h.inventory, it)
	if index < 0 {
		return fmt.Errorf("в инвентаре нет предмета %s", it.Name())
	}

	h.inventory = slices.Delete(h.inventory, index, index+1)
	return nil
}

// Equip экипирует предмет из инвентаря, снятый предмет возвращается в инвентарь
func (h *Hero) Equip(it item.Item) error {
	// switch по типу: в каждой ветке equipment имеет конкретный тип
	switch equipment := it.(type) {
	case item.Weapon:
		if err := h.RemoveItem(equipment); err != nil {
			return err
		}
		h.inventory = append(h.inventory, h.weapon)
		h.weapon = equipment
	case *item.Armor:
		if err := h.RemoveItem(equipment); err != nil {
			return err
		}
		h.inventory = append(h.inventory, h.armor)
		h.armor = equipment
	default:
		return fmt.Errorf("предмет нельзя экипировать: %s", it.Name())
	}

	return nil
}

// геттеры

func (h *Hero) Code() string { return h.heroClass.Code() }

func (h *Hero) Name() string { return h.name }

func (h *Hero) Health() int { return h.health }

func (h *Hero) MaxHealth() int {
	return h.heroClass.BaseHealth() + (h.level-1)*healthPerLevel
}

func (h *Hero) HeroClass() *HeroClass { return h.heroClass }

func (h *Hero) Level() int { return h.level }

func (h *Hero) Experience() int { return h.experience }

func (h *Hero) ExperienceToNextLevel() int { return h.level * experiencePerLevel }

func (h *Hero) Gold() int { return h.gold }

func (h *Hero) Weapon() item.Weapon { return h.weapon }

func (h *Hero) Armor() *item.Armor { return h.armor }

// Inventory возвращает копию списка предметов, чтобы инвентарь нельзя было поменять снаружи
func (h *Hero) Inventory() []item.Item {
	return slices.Clone(h.inventory)
}
