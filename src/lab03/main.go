package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"lab03/camp"
	"lab03/combat"
	"lab03/console"
	"lab03/hero"
	"lab03/item"
	"lab03/world"
)

var (
	input = bufio.NewScanner(os.Stdin)
	view  = console.ConsoleView{}
)

func main() {
	catalog := item.NewDefaultCatalog()
	locations := world.NewLocationRegistry(
		world.NewForest(catalog),
		world.NewCrypt(catalog),
		world.NewDragonLair(catalog),
	)
	merchant := camp.NewMerchant(catalog,
		item.IronSword, item.ElvenBow, item.EmberStaff, item.Chainmail, item.PlateArmor)

	h := createHero(catalog)
	fmt.Printf("%s разбивает лагерь на опушке мистического леса.\n", h.Name())

	for {
		fmt.Println()
		view.PrintHero(h)
		fmt.Println("1 - отправиться в экспедицию, 2 - снаряжение, 3 - торговец, 0 - выйти из игры")
		switch readLine() {
		case "1":
			goOnExpedition(h, locations)
		case "2":
			manageEquipment(h)
		case "3":
			visitMerchant(h, merchant)
		case "0":
			return
		default:
			fmt.Println("Такого варианта нет")
		}
	}
}

func createHero(catalog *item.ItemCatalog) *hero.Hero {
	fmt.Println("Создайте своего героя!")
	for {
		fmt.Println("Выберите класс:")
		heroClasses := hero.HeroClasses()
		for i, heroClass := range heroClasses {
			kit := heroClass.StarterKit(catalog)
			weapon, _ := kit.CreateWeapon()
			armor, _ := kit.CreateArmor()
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

		h, err := hero.NewBuilder().
			SetName(name).
			SetHeroClass(heroClass).
			SetStarterKit(heroClass.StarterKit(catalog)).
			Build()
		if err != nil {
			fmt.Println("Не получилось создать героя: " + err.Error())
			continue
		}

		return h
	}
}

func goOnExpedition(h *hero.Hero, locations *world.LocationRegistry) {
	all := locations.All()
	fmt.Println("Куда отправимся?")
	for i, location := range all {
		fmt.Printf("%d - %s (рекомендуемый уровень %d)\n", i+1, location.DisplayName(), location.RecommendedLevel())
	}
	index := readChoice(len(all))
	if index < 0 {
		return
	}

	location := all[index]
	report, err := world.Explore(location, h)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println()
	view.PrintReport(report, location.DisplayName())

	if report.Battle.Outcome == combat.Victory && location.Code() == world.DragonLairCode {
		fmt.Printf("\nДракон повержен! Барды еще долго будут петь о подвиге героя по имени %s.\n", h.Name())
	}
}

func manageEquipment(h *hero.Hero) {
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
	if err := h.Equip(it); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Экипировано: " + view.Describe(it))
}

func visitMerchant(h *hero.Hero, merchant *camp.Merchant) {
	fmt.Printf("Торговец: \"Чего желаешь?\" (золото: %d)\n", h.Gold())
	fmt.Println("1 - купить, 2 - продать, 0 - назад")
	switch readLine() {
	case "1":
		buy(h, merchant)
	case "2":
		sell(h, merchant)
	}
}

func buy(h *hero.Hero, merchant *camp.Merchant) {
	showcase, err := merchant.Showcase()
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

	it, err := merchant.Buy(h, showcase[index].Code())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Куплено: " + view.Describe(it))
}

func sell(h *hero.Hero, merchant *camp.Merchant) {
	inventory := h.Inventory()
	if len(inventory) == 0 {
		fmt.Println("Продавать нечего (экипированные предметы продать нельзя)")
		return
	}

	for i, it := range inventory {
		fmt.Printf("%d - %s: %d золота\n", i+1, view.Describe(it), merchant.SellPrice(it))
	}
	index := readChoice(len(inventory))
	if index < 0 {
		return
	}

	it := inventory[index]
	price, err := merchant.Sell(h, it)
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
