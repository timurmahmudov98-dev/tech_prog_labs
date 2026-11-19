package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"lab03_lab04/combat"
	"lab03_lab04/console"
	"lab03_lab04/game"
	"lab03_lab04/hero"
	"lab03_lab04/item"
	"path/filepath"
)

var (
	input = bufio.NewScanner(os.Stdin)
	view  = console.ConsoleView{}
)

func main() {
	saveFile := filepath.Join("saves", "heroes.json")
	absolute, err := filepath.Abs(saveFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println("Файл сохранения: " + absolute)
	g := game.NewGameFacade(saveFile)
	h := chooseHero(g)
	fmt.Printf("%s разбивает лагерь на опушке мистического леса.\n", h.Name())

	for {
		fmt.Println()
		view.PrintHero(h)
		fmt.Println("1 - отправиться в экспедицию, 2 - снаряжение, 3 - торговец, 4 - журнал, 0 - выйти из игры")
		switch readLine() {
		case "1":
			goOnExpedition(h, g)
		case "2":
			manageEquipment(h, g)
		case "3":
			visitMerchant(h, g)
		case "4":
			for _, report := range g.Journal(h.Name()) {
				l, err := g.Location(report.LocationCode)
				if err != nil {
					fmt.Println(err)
					continue
				}
				t, err := g.Tactic(report.TacticCode)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println(view.Summarize(report, l.DisplayName(), t.DisplayName()))
			}
		case "0":
			return
		default:
			fmt.Println("Такого варианта нет")
		}
	}
}

func chooseHero(g *game.GameFacade) *hero.Hero {
	heroes, err := g.Heroes()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if len(heroes) == 0 {
		return createHero(g)
	}
	fmt.Println("Выберите сохраненного героя (0 - создать нового):")
	for i, h := range heroes {
		fmt.Printf("%d - %s, %s %d ур.\n", i+1, h.Name(), h.HeroClass().DisplayName(), h.Level())
	}
	index := readChoice(len(heroes))
	if index >= 0 {
		return heroes[index]
	}
	return createHero(g)
}

func createHero(g *game.GameFacade) *hero.Hero {
	fmt.Println("Создайте своего героя!")
	for {
		fmt.Println("Выберите класс:")
		heroClasses := hero.HeroClasses()
		for i, heroClass := range heroClasses {
			equipment, err := g.StarterKit(heroClass.Code())
			if err != nil {
				fmt.Println(err)
				continue
			}
			weapon, armor := equipment[0], equipment[1]
			fmt.Printf(
				"%d - %s (здоровье %d): %s, %s\n",
				i+1, heroClass.DisplayName(), heroClass.BaseHealth(), view.Describe(weapon), view.Describe(armor),
			)
		}
		classIndex := readChoice(len(heroClasses))
		if classIndex < 0 {
			continue
		}
		heroClass := heroClasses[classIndex]

		fmt.Println("Введите имя:")
		name := readLine()

		h, err := g.CreateHero(name, heroClass.Code())
		if err != nil {
			fmt.Println("Не получилось создать героя: " + err.Error())
			continue
		}

		return h
	}
}

func goOnExpedition(h *hero.Hero, g *game.GameFacade) {
	all := g.Locations()
	fmt.Println("Куда отправимся?")
	for i, location := range all {
		fmt.Printf("%d - %s (рекомендуемый уровень %d)\n", i+1, location.DisplayName(), location.RecommendedLevel())
	}
	index := readChoice(len(all))
	if index < 0 {
		return
	}

	location := all[index]
	fmt.Println("Выберите тактику:")
	tactics := g.Tactics()
	for i, t := range tactics {
		fmt.Printf("%d - %s: %s\n", i+1, t.DisplayName(), t.Description())
	}
	tacticIndex := readChoice(len(tactics))
	if tacticIndex < 0 {
		return
	}
	report, err := g.Explore(h.Name(), location.Code(), tactics[tacticIndex].Code())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println()
	view.PrintReport(report, location.DisplayName())

	if report.Battle.Outcome == combat.Victory && report.Battle.Enemy.Code == "dragon" {
		fmt.Printf("\nДракон повержен! Барды еще долго будут петь о подвиге героя по имени %s.\n", h.Name())
	}
}

func manageEquipment(h *hero.Hero, g *game.GameFacade) {
	inventory := h.Inventory()
	if len(inventory) == 0 {
		fmt.Println("Инвентарь пуст")
		return
	}

	fmt.Println("Что экипировать? (0 - назад)")
	printItems(inventory)
	index := readChoice(len(inventory))
	if index < 0 {
		return
	}

	it := inventory[index]
	if err := g.Equip(h.Name(), index); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Экипировано: " + view.Describe(it))
}

func visitMerchant(h *hero.Hero, g *game.GameFacade) {
	fmt.Printf("Торговец: \"Чего желаешь?\" (золото: %d)\n", h.Gold())
	fmt.Println("1 - купить, 2 - продать, 0 - назад")
	switch readLine() {
	case "1":
		buy(h, g)
	case "2":
		sell(h, g)
	}
}

func buy(h *hero.Hero, g *game.GameFacade) {
	showcase, err := g.Showcase()
	if err != nil {
		fmt.Println(err)
		return
	}
	for i, it := range showcase {
		fmt.Printf("%d - %s: %d золота\n", i+1, view.Describe(it), it.Price())
	}
	index := readChoice(len(showcase))
	if index < 0 {
		return
	}

	it, err := g.Buy(h.Name(), showcase[index].Code())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Куплено: " + view.Describe(it))
}

func sell(h *hero.Hero, g *game.GameFacade) {
	inventory := h.Inventory()
	if len(inventory) == 0 {
		fmt.Println("Продавать нечего (экипированные предметы продать нельзя)")
		return
	}

	for i, it := range inventory {
		fmt.Printf("%d - %s: %d золота\n", i+1, view.Describe(it), g.SellPrice(it))
	}
	index := readChoice(len(inventory))
	if index < 0 {
		return
	}

	it := inventory[index]
	price, err := g.Sell(h.Name(), index)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Продано: %s за %d золота\n", it.Name(), price)
}

func printItems(items []item.Item) {
	for i, it := range items {
		fmt.Printf("%d - %s\n", i+1, view.Describe(it))
	}
}

// readLine читает строку без пробелов по краям; если ввод закончился, завершает программу
func readLine() string {
	if !input.Scan() {
		os.Exit(0)
	}

	return strings.TrimSpace(input.Text())
}

// readChoice читает номер варианта от 1 до count, возвращает индекс или -1, если ввели что-то другое
func readChoice(count int) int {
	choice, err := strconv.Atoi(readLine())
	if err != nil || choice < 1 || choice > count {
		return -1
	}

	return choice - 1
}
