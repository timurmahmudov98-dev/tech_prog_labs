package main

import (
	"fmt"

	"lab03/console"
	"lab03/core"
	"lab03/hero"
	"lab03/item"
	"lab03/world"
)

func main() {
	catalog := item.NewDefaultCatalog()
	view := console.ConsoleView{}

	fmt.Println("===== Одиночка =====")
	first := core.GetInstance()
	second := core.GetInstance()
	fmt.Printf("core.GetInstance() возвращает один и тот же объект: %t\n", first == second)

	fmt.Println()
	fmt.Println("===== Прототип =====")
	plainSword := must(catalog.Create(item.IronSword))
	fierySword := must(catalog.Create(item.IronSword))
	if err := fierySword.Enchant(item.Fiery); err != nil {
		panic(err)
	}
	fmt.Println("Первая копия:        " + view.Describe(plainSword))
	fmt.Println("Вторая копия:        " + view.Describe(fierySword))
	fmt.Println("Новая копия шаблона: " + view.Describe(must(catalog.Create(item.IronSword))))
	fmt.Printf("Копии - разные объекты: %t\n", plainSword != fierySword)

	fmt.Println()
	fmt.Println("===== Абстрактная фабрика =====")
	for _, heroClass := range hero.HeroClasses() {
		kit := heroClass.StarterKit(catalog)
		weapon := must(kit.CreateWeapon())
		armor := must(kit.CreateArmor())
		fmt.Printf("%s: %s, %s\n", heroClass.DisplayName(), view.Describe(weapon), view.Describe(armor))
	}

	fmt.Println()
	fmt.Println("===== Строитель =====")
	h := must(hero.NewBuilder().
		SetName("Артур").
		SetHeroClass(hero.Warrior).
		SetStarterKit(hero.Warrior.StarterKit(catalog)).
		SetLevel(3).
		Build())
	view.PrintHero(h)
	_, err := hero.NewBuilder().
		SetHeroClass(hero.Mage).
		SetStarterKit(hero.Mage.StarterKit(catalog)).
		Build()
	fmt.Println("Герой без имени не создается: " + err.Error())
	_, err = hero.NewBuilder().SetName("Мерлин").Build()
	fmt.Println("Герой без класса не создается: " + err.Error())

	fmt.Println()
	fmt.Println("===== Фабричный метод =====")
	forest := world.NewForest(catalog)
	crypt := world.NewCrypt(catalog)
	view.PrintReport(must(world.ExploreWithSeed(forest, h, 7)), forest.DisplayName())
	fmt.Println()
	view.PrintReport(must(world.ExploreWithSeed(crypt, h, 7)), crypt.DisplayName())

}

// must возвращает значение, а если случилась ошибка, останавливает демонстрацию:
// в этой программе все данные заданы заранее, и ошибка означает ошибку в коде
func must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}

	return value
}
