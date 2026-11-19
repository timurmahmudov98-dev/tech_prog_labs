package world

import (
	"lab03_lab04/tactic"

	"lab03_lab04/combat"
	"lab03_lab04/core"
	"lab03_lab04/enemy"
	"lab03_lab04/hero"
	"lab03_lab04/item"
)

const (
	lootChance    = 0.6
	enchantChance = 0.3
)

// Location - локация. Каждая конкретная локация сама решает, какого врага создать (фабричный метод CreateEnemy)
type Location interface {
	Code() string
	DisplayName() string
	RecommendedLevel() int
	Catalog() *item.ItemCatalog
	// CreateEnemy - создание врага делегируется конкретным локациям для повышения гибкости
	CreateEnemy() enemy.Enemy
	// LootTable - коды предметов из каталога, которые можно найти в локации
	LootTable() []string
}

// baseLocation - общие данные всех локаций, конкретные локации встраивают эту структуру в себя
type baseLocation struct {
	code             string
	displayName      string
	recommendedLevel int
	catalog          *item.ItemCatalog
}

func (l *baseLocation) Code() string { return l.code }

func (l *baseLocation) DisplayName() string { return l.displayName }

func (l *baseLocation) RecommendedLevel() int { return l.recommendedLevel }

func (l *baseLocation) Catalog() *item.ItemCatalog { return l.catalog }

// Explore отправляет героя в экспедицию со случайным сидом
func Explore(location Location, h *hero.Hero) (ExpeditionReport, error) {
	return ExploreWithSeed(location, h, core.GetInstance().NextSeed())
}

// ExploreWithSeed - бизнес-логика, общая для всех локаций:
// герой встречает врага, сражается с ним и получает награду, если победил
func ExploreWithSeed(location Location, h *hero.Hero, seed uint64) (ExpeditionReport, error) {
	return ExploreWithTactic(location, h, tactic.BalancedTactic{}, seed)
}

func ExploreWithTactic(location Location, h *hero.Hero, strategy tactic.Tactic, seed uint64) (ExpeditionReport, error) {
	core.GetInstance().Reseed(seed)

	h.Rest()
	foe := location.CreateEnemy()
	if core.GetInstance().Chance(0.2) {
		foe = makeElite(foe)
	}
	battle := combat.NewBattle(h, foe, strategy).Run()

	var rewards Rewards
	if battle.Outcome == combat.Victory {
		var err error
		rewards, err = reward(location, h, foe)
		if err != nil {
			return ExpeditionReport{}, err
		}
	} else {
		rewards = punish(h)
	}

	return ExpeditionReport{Seed: seed, LocationCode: location.Code(), TacticCode: strategy.Code(), Battle: battle, Rewards: rewards}, nil
}

func reward(location Location, h *hero.Hero, foe enemy.Enemy) (Rewards, error) {
	experience := foe.ExperienceReward()
	gold := foe.RollGoldReward()
	loot, err := rollLoot(location)
	if err != nil {
		return Rewards{}, err
	}

	levelsGained := h.GainExperience(experience)
	h.AddGold(gold)
	if loot != nil {
		h.AddItem(loot)
	}

	return Rewards{Experience: experience, Gold: gold, Loot: loot, LevelsGained: levelsGained}, nil
}

func rollLoot(location Location) (item.Item, error) {
	dice := core.GetInstance()
	if !dice.Chance(lootChance) {
		return nil, nil
	}

	// добыча - это копия шаблона из каталога, ее можно менять, не боясь испортить шаблон
	table := location.LootTable()
	loot, err := location.Catalog().Create(table[dice.Index(len(table))])
	if err != nil {
		return nil, err
	}
	if dice.Chance(enchantChance) {
		possible := loot.PossibleEnchantments()
		if err := loot.Enchant(possible[dice.Index(len(possible))]); err != nil {
			return nil, err
		}
	}

	return loot, nil
}

// проигравший герой теряет половину золота
func punish(h *hero.Hero) Rewards {
	lostGold := h.Gold() / 2
	// потратить половину имеющегося золота можно всегда, поэтому ошибку здесь не проверяем
	_ = h.SpendGold(lostGold)
	return Rewards{Experience: 0, Gold: -lostGold, Loot: nil, LevelsGained: 0}
}

func makeElite(foe enemy.Enemy) enemy.Enemy {
	dice := core.GetInstance()
	both := dice.Chance(0.25)
	if dice.Chance(0.5) {
		foe = enemy.NewLegendaryEnemy(foe)
		if both {
			foe = enemy.NewWindfuryEnemy(foe)
		}
	} else {
		foe = enemy.NewWindfuryEnemy(foe)
		if both {
			foe = enemy.NewLegendaryEnemy(foe)
		}
	}
	if dice.Chance(0.3) {
		foe = enemy.NewArmoredEnemy(foe)
	}
	return foe
}
