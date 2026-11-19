package game

import (
	"fmt"
	"lab03_lab04/camp"
	"lab03_lab04/combat"
	"lab03_lab04/core"
	"lab03_lab04/event"
	"lab03_lab04/hero"
	"lab03_lab04/item"
	"lab03_lab04/storage"
	"lab03_lab04/tactic"
	"lab03_lab04/world"
	"slices"
)

type GameFacade struct {
	catalog   *item.ItemCatalog
	locations *world.LocationRegistry
	merchant  *camp.Merchant
	tactics   []tactic.Tactic
	heroes    storage.HeroRepository
	events    event.GameEventPublisher
	journal   *event.ExpeditionJournal
}

func NewGameFacade(saveFile string) *GameFacade {
	catalog := item.NewDefaultCatalog()
	g := &GameFacade{catalog: catalog, locations: world.NewLocationRegistry(world.NewForest(catalog), world.NewHauntedManor(catalog), world.NewCrypt(catalog), world.NewDragonLair(catalog)),
		merchant: camp.NewMerchant(catalog, item.IronSword, item.ElvenBow, item.EmberStaff, item.Chainmail, item.PlateArmor),
		tactics:  []tactic.Tactic{tactic.BalancedTactic{}, tactic.AggressiveTactic{}, tactic.CautiousTactic{}},
		heroes:   storage.NewCachingHeroRepository(storage.NewFileHeroRepository(saveFile, catalog)), journal: &event.ExpeditionJournal{}}
	a := &event.AchievementTracker{}
	g.events.Subscribe(event.Victory, a)
	g.events.Subscribe(event.LevelUp, a)
	g.events.Subscribe(event.Victory, g.journal)
	g.events.Subscribe(event.Defeat, g.journal)
	return g
}
func (g *GameFacade) HeroClasses() []*hero.HeroClass { return hero.HeroClasses() }
func (g *GameFacade) StarterKit(code string) ([]item.Item, error) {
	hc, err := hero.HeroClassByCode(code)
	if err != nil {
		return nil, err
	}
	kit := hc.StarterKit(g.catalog)
	w, err := kit.CreateWeapon()
	if err != nil {
		return nil, err
	}
	a, err := kit.CreateArmor()
	if err != nil {
		return nil, err
	}
	return []item.Item{w, a}, nil
}
func (g *GameFacade) Locations() []world.Location                  { return g.locations.All() }
func (g *GameFacade) Location(code string) (world.Location, error) { return g.locations.Get(code) }
func (g *GameFacade) Tactics() []tactic.Tactic                     { return slices.Clone(g.tactics) }
func (g *GameFacade) Tactic(code string) (tactic.Tactic, error) {
	for _, t := range g.tactics {
		if t.Code() == code {
			return t, nil
		}
	}
	return nil, &core.UnknownCodeError{What: "Тактика", Code: code}
}
func (g *GameFacade) Showcase() ([]item.Item, error) { return g.merchant.Showcase() }
func (g *GameFacade) SellPrice(it item.Item) int     { return g.merchant.SellPrice(it) }
func (g *GameFacade) Heroes() ([]*hero.Hero, error)  { return g.heroes.FindAll() }
func (g *GameFacade) Hero(name string) (*hero.Hero, error) {
	h, err := g.heroes.FindByName(name)
	if err != nil {
		return nil, err
	}
	if h == nil {
		return nil, fmt.Errorf("героя с именем %s нет", name)
	}
	return h, nil
}
func (g *GameFacade) CreateHero(name, code string) (*hero.Hero, error) {
	found, err := g.heroes.FindByName(name)
	if err != nil {
		return nil, err
	}
	if found != nil {
		return nil, fmt.Errorf("герой с именем %s уже есть", name)
	}
	hc, err := hero.HeroClassByCode(code)
	if err != nil {
		return nil, err
	}
	h, err := hero.NewBuilder().SetName(name).SetHeroClass(hc).SetStarterKit(hc.StarterKit(g.catalog)).Build()
	if err != nil {
		return nil, err
	}
	return h, g.heroes.Save(h)
}
func (g *GameFacade) Explore(name, locationCode, tacticCode string) (world.ExpeditionReport, error) {
	h, err := g.Hero(name)
	if err != nil {
		return world.ExpeditionReport{}, err
	}
	l, err := g.Location(locationCode)
	if err != nil {
		return world.ExpeditionReport{}, err
	}
	t, err := g.Tactic(tacticCode)
	if err != nil {
		return world.ExpeditionReport{}, err
	}
	report, err := world.ExploreWithTactic(l, h, t, core.GetInstance().NextSeed())
	if err != nil {
		return report, err
	}
	result := event.Defeat
	if report.Battle.Outcome == combat.Victory {
		result = event.Victory
	}
	g.events.Publish(event.GameEvent{Type: result, Hero: h, Expedition: report})
	if report.Rewards.LevelsGained > 0 {
		g.events.Publish(event.GameEvent{Type: event.LevelUp, Hero: h, Expedition: report})
	}
	return report, g.heroes.Save(h)
}
func (g *GameFacade) Equip(name string, index int) error {
	h, err := g.Hero(name)
	if err != nil {
		return err
	}
	inv := h.Inventory()
	if index < 0 || index >= len(inv) {
		return fmt.Errorf("нет предмета с таким номером")
	}
	if err = h.Equip(inv[index]); err != nil {
		return err
	}
	return g.heroes.Save(h)
}
func (g *GameFacade) Buy(name, code string) (item.Item, error) {
	h, err := g.Hero(name)
	if err != nil {
		return nil, err
	}
	it, err := g.merchant.Buy(h, code)
	if err != nil {
		return nil, err
	}
	return it, g.heroes.Save(h)
}
func (g *GameFacade) Sell(name string, index int) (int, error) {
	h, err := g.Hero(name)
	if err != nil {
		return 0, err
	}
	inv := h.Inventory()
	if index < 0 || index >= len(inv) {
		return 0, fmt.Errorf("нет предмета с таким номером")
	}
	price, err := g.merchant.Sell(h, inv[index])
	if err != nil {
		return 0, err
	}
	return price, g.heroes.Save(h)
}
func (g *GameFacade) Journal(name string) []world.ExpeditionReport { return g.journal.Entries(name) }
