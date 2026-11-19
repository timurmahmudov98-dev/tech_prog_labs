package main

import (
	"fmt"
	"lab03_lab04/combat"
	"lab03_lab04/console"
	"lab03_lab04/core"
	"lab03_lab04/defense"
	"lab03_lab04/enemy"
	"lab03_lab04/event"
	"lab03_lab04/game"
	"lab03_lab04/hero"
	"lab03_lab04/item"
	"lab03_lab04/storage"
	"lab03_lab04/tactic"
	"lab03_lab04/world"
	"os"
	"path/filepath"
	"strings"
)

func must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}
func check(err error) {
	if err != nil {
		panic(err)
	}
}
func newHero(name string, hc *hero.HeroClass, level int, catalog *item.ItemCatalog) *hero.Hero {
	return must(hero.NewBuilder().SetName(name).SetHeroClass(hc).SetStarterKit(hc.StarterKit(catalog)).SetLevel(level).Build())
}
func printHits(title string, chain defense.DamageHandler, count int) {
	results := []string{}
	for i := 0; i < count; i++ {
		d := defense.Damage{Amount: 30}
		chain.Handle(&d)
		if d.Evaded {
			results = append(results, "уклонение")
		} else if d.Absorbed > 0 {
			results = append(results, fmt.Sprintf("поглощено %d, прошло %d", d.Absorbed, d.Amount))
		} else {
			results = append(results, fmt.Sprintf("прошло %d", d.Amount))
		}
	}
	fmt.Println(title + ": " + strings.Join(results, "; "))
}
func main() {
	view := console.ConsoleView{}
	catalog := item.NewDefaultCatalog()
	dice := core.GetInstance()
	fmt.Println("===== Цепочка обязанностей =====")
	printHits("Броня 30%", defense.NewArmorHandler(0.3), 2)
	armor := defense.NewArmorHandler(0.3)
	armor.SetNext(defense.NewBarrierHandler(30))
	printHits("Броня 30% -> барьер 30", armor, 3)
	dice.Reseed(3)
	chain := defense.NewEvasionHandler(0.5)
	chain.SetNext(defense.NewArmorHandler(0.3)).SetNext(defense.NewBarrierHandler(30))
	printHits("Уклонение 50% -> броня 30% -> барьер 30", chain, 3)
	fmt.Println("\n===== Стратегия =====")
	for _, t := range []tactic.Tactic{tactic.BalancedTactic{}, tactic.AggressiveTactic{}, tactic.CautiousTactic{}} {
		h := newHero("Мерлин", hero.Mage, 2, catalog)
		report := must(world.ExploreWithTactic(world.NewCrypt(catalog), h, t, 9))
		result := "поражение"
		if report.Battle.Outcome == combat.Victory {
			result = "победа"
		}
		fmt.Printf("%s: %s, раундов %d, у героя осталось %d/%d здоровья\n", t.DisplayName(), result, report.Battle.Events[len(report.Battle.Events)-1].Round, h.Health(), h.MaxHealth())
	}
	fmt.Println("\n===== Декоратор =====")
	for _, e := range []enemy.Enemy{enemy.NewGoblin(), enemy.NewLegendaryEnemy(enemy.NewGoblin()), enemy.NewWindfuryEnemy(enemy.NewGoblin()), enemy.NewLegendaryEnemy(enemy.NewWindfuryEnemy(enemy.NewGoblin())), enemy.NewArmoredEnemy(enemy.NewGoblin()), enemy.NewArmoredEnemy(enemy.NewLegendaryEnemy(enemy.NewWindfuryEnemy(enemy.NewDragon())))} {
		dice.Reseed(5)
		hits := []string{}
		for i := 0; i < 3; i++ {
			hits = append(hits, fmt.Sprint(e.RollDamage().Damage))
		}
		d := e.ReceiveHit(20)
		fmt.Printf("%s: удары %s, награда %d опыта; входящий урон 20 -> %d, здоровье %d/%d\n", e.Name(), strings.Join(hits, ", "), e.ExperienceReward(), d.Amount, e.Health(), e.MaxHealth())
	}
	fmt.Println("\n===== Адаптер =====")
	bow := must(catalog.CreateWeapon(item.ElvenBow))
	haunted := enemy.NewHauntedWeapon(bow)
	fmt.Println("Предмет: " + view.Describe(bow))
	fmt.Printf("Враг:    %s (%d здоровья, награда %d опыта)\n", haunted.Name(), haunted.MaxHealth(), haunted.ExperienceReward())
	manor := world.NewHauntedManor(catalog)
	view.PrintReport(must(world.ExploreWithTactic(manor, newHero("Леголас", hero.Ranger, 2, catalog), tactic.CautiousTactic{}, 55)), manor.DisplayName())
	fmt.Println("\n===== Наблюдатель =====")
	events := &event.GameEventPublisher{}
	journal := &event.ExpeditionJournal{}
	events.Subscribe(event.Victory, &event.AchievementTracker{})
	events.Subscribe(event.Victory, journal)
	events.Subscribe(event.Defeat, journal)
	events.Subscribe(event.Victory, event.ListenerFunc(func(e event.GameEvent) { fmt.Printf("Событие %s: герой %s\n", e.Type, e.Hero.Name()) }))
	warrior := newHero("Артур", hero.Warrior, 1, catalog)
	report := must(world.ExploreWithSeed(world.NewForest(catalog), warrior, 7))
	result := event.Defeat
	if report.Battle.Outcome == combat.Victory {
		result = event.Victory
	}
	events.Publish(event.GameEvent{Type: result, Hero: warrior, Expedition: report})
	names := []string{}
	for _, a := range warrior.Achievements() {
		names = append(names, a.DisplayName())
	}
	fmt.Println("Достижения героя: " + strings.Join(names, ", "))
	fmt.Println("Записей в журнале:", len(journal.Entries(warrior.Name())))
	directory := must(os.MkdirTemp("", "lab-demo-"))
	defer os.RemoveAll(directory)
	fmt.Println("\n===== Заместитель =====")
	path := filepath.Join(directory, "heroes.json")
	file := storage.NewFileHeroRepository(path, catalog)
	check(file.Save(warrior))
	fmt.Println("Содержимое файла:")
	fmt.Print(string(must(os.ReadFile(path))))
	fmt.Println("Без заместителя два поиска дают один и тот же объект:", must(file.FindByName("Артур")) == must(file.FindByName("Артур")))
	cached := storage.NewCachingHeroRepository(file)
	h := must(cached.FindByName("Артур"))
	fmt.Println("С заместителем два поиска дают один и тот же объект:", h == must(cached.FindByName("Артур")))
	h.AddGold(100)
	check(cached.Save(h))
	fmt.Println("После сохранения через заместитель золото героя в файле:", must(file.FindByName("Артур")).Gold())
	fmt.Println("\n===== Фасад =====")
	g := game.NewGameFacade(filepath.Join(directory, "game.json"))
	must(g.CreateHero("Гэндальф", "mage"))
	dice.Reseed(11)
	expedition := must(g.Explore("Гэндальф", "forest", "cautious"))
	fmt.Println("Экспедиция через фасад: " + view.Summarize(expedition, must(g.Location(expedition.LocationCode)).DisplayName(), must(g.Tactic(expedition.TacticCode)).DisplayName()))
	view.PrintHero(must(g.Hero("Гэндальф")))
}
