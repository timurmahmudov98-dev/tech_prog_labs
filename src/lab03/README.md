# Лаба 3 — полное описание структуры проекта

Расположение: `tech_prog_labs/src/lab03/` (Go-модуль `lab03`, импорты вида `lab03/combat`).

Проект — консольная RPG. Игрок создаёт героя, ходит в экспедиции (бой с врагом), получает опыт, золото и добычу, меняет снаряжение и торгует. Лаба демонстрирует пять паттернов: **Одиночка, Прототип, Абстрактная фабрика, Строитель, Фабричный метод**.

Как читать документ:
- В разделе 1 — общая картина и дерево проекта.
- В разделе 2 — иерархия **папка → файл → тип → методы** в виде дерева.
- В разделе 3 — подробное описание каждой папки, файла, типа и метода: за что отвечает и **от чего зависит** (какой тип встраивает, какой интерфейс реализует, какие методы вызывает).
- В разделе 4 — паттерны, в разделе 5 — сквозной сценарий одной экспедиции.

Обозначения в колонке «Зависит от»: `→` означает «вызывает / использует», `⊂` — «встраивает в себя», `⊨` — «реализует интерфейс», `⇒` — «возвращает ошибку этого типа».

---

## 1. Общая картина

### Слои и направление зависимостей

Пакеты зависят друг от друга только в одну сторону, циклов нет:

```
core  →  combat  →  item  →  kit  →  hero  →  enemy  →  world  →  camp  →  console  →  main
```

Точнее (реальные импорты):

| Пакет | Импортирует |
|---|---|
| `core` | ничего своего |
| `combat` | ничего своего |
| `item` | `combat`, `core` |
| `kit` | `item` |
| `hero` | `combat`, `core`, `item`, `kit` |
| `enemy` | `combat`, `core` |
| `world` | `combat`, `core`, `enemy`, `hero`, `item` |
| `camp` | `core`, `hero`, `item` |
| `console` | `combat`, `hero`, `item`, `world` |
| `main` | `camp`, `combat`, `console`, `hero`, `item`, `world` |
| `cmd/demo` | `console`, `core`, `hero`, `item`, `world` |

### Дерево проекта

```
lab03/
├── go.mod                    модуль lab03, Go 1.22
├── main.go                   игра: меню и игровой цикл
├── cmd/demo/main.go          демонстрация всех паттернов
├── core/                     общая инфраструктура
│   ├── dice.go
│   └── errors.go
├── combat/                   боевая система
│   ├── combatant.go  hit.go  battle.go
│   ├── battle_event.go  battle_report.go  participant.go
│   └── event_type.go  side.go  outcome.go
├── item/                     предметы
│   ├── item.go  weapon.go  sword.go  bow.go  staff.go  armor.go
│   └── enchantment.go  catalog.go  codes.go
├── kit/                      стартовое снаряжение
│   ├── starter_kit.go
│   └── warrior.go  ranger.go  mage.go
├── hero/                     герой
│   ├── hero_class.go
│   └── hero.go
├── enemy/                    враги
│   ├── enemy.go  base_enemy.go
│   └── goblin.go  skeleton.go  dragon.go
├── world/                    локации и экспедиции
│   ├── location.go  registry.go
│   ├── forest.go  crypt.go  dragon_lair.go
│   └── rewards.go  expedition_report.go
├── camp/                     лагерь
│   └── merchant.go
└── console/                  вывод в консоль
    └── view.go
```

### Сколько типов в проекте

Всего 43 типа: 34 структуры, 6 интерфейсов (`Combatant`, `Enemy`, `Item`, `Weapon`, `StarterKitFactory`, `Location`) и 3 перечисления (`EventType`, `Outcome`, `Side`). Не в каждом файле есть одноимённая структура: например, в `item/codes.go` типов нет совсем, а в `hero/hero.go`, `item/item.go`, `item/weapon.go`, `world/location.go` их по два.

---

## 2. Иерархия: папка → файл → тип → методы

Методы с маленькой буквы неэкспортируемые (доступны только внутри пакета).

```
core/
├─ dice.go
│   └─ dice (struct)            GetInstance(), Reseed(), NextSeed(), Roll(), Chance(), Index()
└─ errors.go
    └─ UnknownCodeError (struct) Error()

combat/
├─ combatant.go      └─ Combatant (interface)   Code, Name, Health, MaxHealth, IsAlive, RollDamage, ReceiveHit
├─ hit.go            └─ Hit (struct)            — только поля
├─ battle.go         └─ Battle (struct)         NewBattle(), Run(), strike(), emit()
├─ battle_event.go   └─ BattleEvent (struct)    — только поля
├─ battle_report.go  └─ BattleReport (struct)   — только поля
├─ participant.go    └─ Participant (struct)    ParticipantOf() (функция)
├─ event_type.go     └─ EventType (int)         EventHit, EventCriticalHit, EventDeath
├─ side.go           └─ Side (int)              SideHero, SideEnemy
└─ outcome.go        └─ Outcome (int)           Victory, Defeat

item/
├─ item.go           ├─ Item (interface)        Clone, Code, Name, Price, Enchantments, PossibleEnchantments, Enchant, Description
│                    └─ baseItem (struct)       cloneBase(), enchant(), Code(), Name(), Price(), Enchantments()
├─ weapon.go         ├─ Weapon (interface)      Item + RollDamage, Damage
│                    └─ baseWeapon (struct)     cloneWeapon(), Damage(), PossibleEnchantments(), Enchant()
├─ sword.go          └─ Sword (struct)          NewSword(), RollDamage(), Clone(), Description()
├─ bow.go            └─ Bow (struct)            NewBow(), RollDamage(), Clone(), Description()
├─ staff.go          └─ Staff (struct)          NewStaff(), RollDamage(), Clone(), Description(), minDamage(), maxDamage()
├─ armor.go          └─ Armor (struct)          NewArmor(), Clone(), Defense(), PossibleEnchantments(), Enchant(), Description()
├─ enchantment.go    └─ Enchantment (struct)    Code, Name, BonusDamage, BonusDefense, BonusPrice  + Sharp, Fiery, Sturdy, Warding, ForWeapons, ForArmor
├─ catalog.go        └─ ItemCatalog (struct)    Add(), Create(), CreateWeapon(), CreateArmor(), Codes()
└─ codes.go          (типов нет)                константы кодов предметов + NewDefaultCatalog()

kit/
├─ starter_kit.go    └─ StarterKitFactory (interface)  CreateWeapon, CreateArmor
├─ warrior.go        └─ WarriorKitFactory (struct)     NewWarriorKitFactory(), CreateWeapon(), CreateArmor()
├─ ranger.go         └─ RangerKitFactory (struct)      NewRangerKitFactory(), CreateWeapon(), CreateArmor()
└─ mage.go           └─ MageKitFactory (struct)        NewMageKitFactory(), CreateWeapon(), CreateArmor()

hero/
├─ hero_class.go     └─ HeroClass (struct)      HeroClasses(), HeroClassByCode(), StarterKit(), Code(), DisplayName(), BaseHealth()  + Warrior, Ranger, Mage
└─ hero.go           ├─ Builder (struct)        NewBuilder(), SetName/SetHeroClass/SetWeapon/SetArmor/SetStarterKit/SetLevel/SetGold, Build(), fail()
                     └─ Hero (struct)           RollDamage, ReceiveHit, IsAlive, RestoreHealth, GainExperience, AddGold, SpendGold,
                                                AddItem, RemoveItem, Equip + геттеры

enemy/
├─ enemy.go          └─ Enemy (interface)       Combatant + ExperienceReward, RollGoldReward
├─ base_enemy.go     └─ baseEnemy (struct)      newBaseEnemy(), RollDamage, ReceiveHit, IsAlive, RollGoldReward, ExperienceReward, Code, Name, Health, MaxHealth
├─ goblin.go         └─ Goblin (struct)         NewGoblin()
├─ skeleton.go       └─ Skeleton (struct)       NewSkeleton()
└─ dragon.go         └─ Dragon (struct)         NewDragon(), RollDamage()

world/
├─ location.go       ├─ Location (interface)    Code, DisplayName, RecommendedLevel, Catalog, CreateEnemy, LootTable
│                    └─ baseLocation (struct)   Code, DisplayName, RecommendedLevel, Catalog
│                    (функции пакета)           Explore(), ExploreWithSeed(), reward(), rollLoot(), punish()
├─ forest.go         └─ Forest (struct)         NewForest(), CreateEnemy(), LootTable()
├─ crypt.go          └─ Crypt (struct)          NewCrypt(), CreateEnemy(), LootTable()
├─ dragon_lair.go    └─ DragonLair (struct)     NewDragonLair(), CreateEnemy(), LootTable()
├─ registry.go       └─ LocationRegistry (struct) NewLocationRegistry(), Get(), All()
├─ rewards.go        └─ Rewards (struct)        — только поля
└─ expedition_report.go └─ ExpeditionReport (struct) — только поля

camp/
└─ merchant.go       └─ Merchant (struct)       NewMerchant(), Showcase(), Buy(), Sell(), SellPrice()

console/
└─ view.go           └─ ConsoleView (struct)    Describe(), PrintHero(), PrintReport(), printEvent(), printRewards()  + функция participant()

main.go              (типов нет)  main(), createHero(), goOnExpedition(), manageEquipment(), visitMerchant(), buy(), sell(), printItems(), readLine(), readChoice()
cmd/demo/main.go     (типов нет)  main(), must[T]()
```

---

## 3. Подробное описание по папкам

---

###  `core/` — общая инфраструктура

**За что отвечает папка.** Здесь лежит то, что нужно всем остальным пакетам и что не относится к игровым правилам: единый источник случайности и общая ошибка «не найдено по коду». Пакет ничего не импортирует из проекта, поэтому на него могут ссылаться все.

####  `core/dice.go`

##### Тип `dice` (struct) — паттерн **Одиночка**
Единственный на всю игру генератор случайных чисел. Все броски (урон посоха, шанс крита, дроп, выбор предмета) идут через него. Тип неэкспортируемый (с маленькой буквы), поэтому создать вторую «кость» вне пакета `core` нельзя — получить её можно только через `GetInstance()`.

Поле: `random *rand.Rand` — генератор из `math/rand/v2` (алгоритм PCG).
Переменные пакета: `instance *dice` (единственный экземпляр) и `once sync.Once` (гарантия однократного создания).

| Метод | Что делает | Зависит от |
|---|---|---|
| `GetInstance() *dice` | Возвращает единственный экземпляр. При первом вызове создаёт его со случайным сидом. `sync.Once` гарантирует, что это произойдёт ровно один раз, даже при вызове из нескольких горутин. | → `sync.Once.Do`, `rand.New`, `rand.NewPCG` |
| `Reseed(seed uint64)` | Пересоздаёт генератор с заданным сидом. После этого все последующие броски повторяются в точности. Нужен для воспроизводимых экспедиций. | → `rand.New`, `rand.NewPCG` |
| `NextSeed() uint64` | Возвращает случайный сид (число до 1 000 000), который потом можно показать игроку и использовать повторно. | → `random.Uint64N` |
| `Roll(low, high int) int` | Случайное целое от `low` до `high` включительно. | → `random.IntN` |
| `Chance(probability float64) bool` | `true` с вероятностью `probability` (от 0 до 1). Так срабатывают крит лука, огненное дыхание дракона, выпадение добычи. | → `random.Float64` |
| `Index(size int) int` | Случайный индекс от 0 до `size-1`, чтобы выбрать случайный элемент слайса. | → `random.IntN` |

####  `core/errors.go`

##### Тип `UnknownCodeError` (struct)
Общая ошибка «по такому коду ничего не нашлось». Поля: `What` (что искали: «Предмет», «Локация», «Класс героя», «Товар») и `Code` (искомый код).

| Метод | Что делает | Зависит от |
|---|---|---|
| `Error() string` | Формирует текст ошибки: «<What> с кодом '<Code>' не существует». Благодаря этому методу тип реализует встроенный интерфейс `error`. | ⊨ `error`; → `fmt.Sprintf` |

Кто её возвращает: `ItemCatalog.Create`, `Merchant.Buy`, `LocationRegistry.Get`, `HeroClassByCode`.

---

###  `combat/` — боевая система

**За что отвечает папка.** Правила боя в чистом виде. Пакет **ничего не знает про героев, врагов и предметы** — только про абстрактного бойца (`Combatant`). Поэтому один и тот же `Battle` может провести бой между любыми двумя сущностями, которые реализуют этот интерфейс. Сюда же входят все структуры, описывающие результат боя.

####  `combat/combatant.go`

##### Тип `Combatant` (interface)
Всё, что может участвовать в бою. Его реализуют `Hero` (через свои методы) и все враги (через `baseEnemy`).

| Метод интерфейса | Смысл |
|---|---|
| `Code() string` | Код вида сущности: класс героя или тип врага |
| `Name() string` | Отображаемое имя |
| `Health() int` / `MaxHealth() int` | Текущее и максимальное здоровье |
| `IsAlive() bool` | Жив ли боец |
| `RollDamage() Hit` | Удар, который боец нанесёт |
| `ReceiveHit(damage int) int` | Получение удара; возвращает урон, который реально прошёл |

####  `combat/hit.go`
##### Тип `Hit` (struct)
Один удар: `Damage int` (сколько урона) и `Critical bool` (был ли он критическим). Результат `RollDamage()`. Методов нет.

####  `combat/battle.go`
##### Тип `Battle` (struct)
Бой, идущий по раундам: сначала бьёт герой, потом враг, и так пока кто-то не погибнет. Поля: `hero`, `enemy` (оба типа `Combatant`), `events []BattleEvent` (накопленный журнал), `round int`.

| Метод | Что делает | Зависит от |
|---|---|---|
| `NewBattle(hero, enemy Combatant) *Battle` | Создаёт бой между двумя бойцами. | — |
| `Run() BattleReport` | Крутит раунды, пока оба живы: `strike` героя по врагу, затем, если враг жив, `strike` врага по герою. В конце определяет исход (герой жив → `Victory`, иначе `Defeat`) и собирает `BattleReport`. | → `strike`, `IsAlive`, `ParticipantOf` |
| `strike(attacker, source, defender, target)` | Один удар: берёт удар у атакующего, передаёт защитнику, пишет событие (`EventHit` или `EventCriticalHit`), а если защитник погиб — ещё и `EventDeath`. | → `Combatant.RollDamage`, `Combatant.ReceiveHit`, `Combatant.Health`, `emit` |
| `emit(event BattleEvent)` | Добавляет событие в журнал. | — |

####  `combat/battle_event.go`
##### Тип `BattleEvent` (struct)
Одно событие боя. Поля: `Round` (раунд), `Type` (`EventType`), `Source` и `Target` (`Side` — кто и по кому), `Amount` (реальный нанесённый урон), `TargetHealthAfter` (сколько здоровья осталось у цели). Методов нет. Зависит от `EventType` и `Side`.

####  `combat/battle_report.go`
##### Тип `BattleReport` (struct)
Итог боя: `Hero` и `Enemy` (`Participant`), `Outcome` и `Events []BattleEvent`. Именно этот отчёт получает `console` для печати. Методов нет. Зависит от `Participant`, `Outcome`, `BattleEvent`.

####  `combat/participant.go`
##### Тип `Participant` (struct)
Снимок участника боя для отчёта: `Code`, `Name`, `MaxHealth`. Нужен, чтобы отчёт не держал ссылку на живого бойца.

| Функция | Что делает | Зависит от |
|---|---|---|
| `ParticipantOf(combatant Combatant) Participant` | Делает снимок бойца. | → `Code()`, `Name()`, `MaxHealth()` |

####  `combat/event_type.go`, `side.go`, `outcome.go`
Три перечисления (`int` + константы через `iota`), методов нет:
- `EventType`: `EventHit` (обычный удар), `EventCriticalHit` (критический), `EventDeath` (смерть).
- `Side`: `SideHero`, `SideEnemy` — на чьей стороне участник события.
- `Outcome`: `Victory`, `Defeat` — исход боя с точки зрения героя.

---

###  `item/` — предметы (паттерн **Прототип**)

**За что отвечает папка.** Всё, что герой может носить, покупать, продавать: оружие, броня и зачарования. Ключевая идея: каталог хранит **шаблоны** (прототипы), а наружу отдаёт только их **копии**. Поэтому зачарование, наложенное на предмет героя, никогда не портит шаблон в каталоге.

Иерархия предметов (встраивание):

```
Item (interface)
 ├─ Weapon (interface = Item + RollDamage + Damage)
 │    └─ Sword / Bow / Staff   ⊂ baseWeapon ⊂ baseItem
 └─ Armor                      ⊂ baseItem
```

####  `item/item.go`

##### Тип `Item` (interface)
Любой предмет. Главный метод — `Clone() Item`: предмет сам создаёт свою копию того же конкретного типа (в этом суть Прототипа). Остальные: `Code`, `Name`, `Price`, `Enchantments`, `PossibleEnchantments` (какие зачарования подходят), `Enchant`, `Description` (например, «урон 12»).

##### Тип `baseItem` (struct)
Общая часть всех предметов, неэкспортируемая. Поля: `code`, `name`, `basePrice`, `enchantments []*Enchantment`. Встраивается в `baseWeapon` и `Armor`.

| Метод | Что делает | Зависит от |
|---|---|---|
| `cloneBase() baseItem` | Копирует общую часть. Присваивание структуры копирует поля, но слайс — это ссылка, поэтому для копии создаётся свой слайс зачарований, иначе копия и оригинал делили бы один список. | → `slices.Clone` |
| `enchant(possible, enchantment) error` | Накладывает зачарование, если оно входит в список допустимых; иначе возвращает ошибку. | → `slices.Contains`, `Enchantment.Name()` |
| `Code()`, `Name()` | Геттеры. | — |
| `Price() int` | Базовая цена плюс бонусы цены всех зачарований. | → `Enchantment.BonusPrice()` |
| `Enchantments() []*Enchantment` | Возвращает копию списка зачарований (снаружи его не изменить). | → `slices.Clone` |

####  `item/weapon.go`

##### Тип `Weapon` (interface)
Оружие: `Item` плюс `RollDamage() combat.Hit` (как оружие бьёт — решает конкретный тип) и `Damage() int` (урон с зачарованиями). Зависит от `Item`, `combat.Hit`.

##### Тип `baseWeapon` (struct) ⊂ `baseItem`
Общая часть всего оружия. Поле: `damage int`.

| Метод | Что делает | Зависит от |
|---|---|---|
| `cloneWeapon() baseWeapon` | Копирует оружие вместе с общей частью. | → `baseItem.cloneBase` |
| `Damage() int` | Базовый урон + бонусы урона всех зачарований. | → `Enchantment.BonusDamage()` |
| `PossibleEnchantments()` | Для оружия подходят `ForWeapons` (Острое, Огненное). | → `ForWeapons` |
| `Enchant(e)` | Накладывает зачарование из списка для оружия. | → `baseItem.enchant`, `ForWeapons` |

####  `item/sword.go` — тип `Sword` ⊂ `baseWeapon`
Меч бьёт стабильно.

| Метод | Что делает | Зависит от |
|---|---|---|
| `NewSword(code, name, basePrice, damage)` | Конструктор. | — |
| `RollDamage() combat.Hit` | Всегда `Damage()`, без крита. | → `baseWeapon.Damage` |
| `Clone() Item` | Копия меча. | → `cloneWeapon` |
| `Description()` | «урон N». | → `Damage` |

####  `item/bow.go` — тип `Bow` ⊂ `baseWeapon`
Лук с некоторым шансом наносит **двойной** урон (константа `criticalMultiplier = 2`). Доп. поле: `criticalChance`.

| Метод | Что делает | Зависит от |
|---|---|---|
| `NewBow(code, name, basePrice, damage, criticalChance)` | Конструктор. | — |
| `RollDamage()` | С вероятностью `criticalChance` урон ×2 и `Critical: true`, иначе обычный. | → `core.GetInstance().Chance`, `Damage` |
| `Clone()` | Копия лука (копируется вся структура, общая часть — отдельно). | → `cloneWeapon` |
| `Description()` | «урон N, шанс крита P%». | → `Damage` |

####  `item/staff.go` — тип `Staff` ⊂ `baseWeapon`
Урон посоха непредсказуем: может быть как больше, так и меньше базового. Доп. поле: `scatter` (разброс, доля).

| Метод | Что делает | Зависит от |
|---|---|---|
| `NewStaff(code, name, basePrice, damage, scatter)` | Конструктор. | — |
| `RollDamage()` | Случайный урон между `minDamage()` и `maxDamage()`. | → `core.GetInstance().Roll` |
| `Clone()` | Копия посоха. | → `cloneWeapon` |
| `Description()` | «урон min-max». | → `minDamage`, `maxDamage` |
| `minDamage()` / `maxDamage()` | Нижняя/верхняя граница: `Damage × (1 ∓ scatter)`. | → `Damage` |

####  `item/armor.go` — тип `Armor` ⊂ `baseItem`
Броня поглощает долю входящего урона. Поле: `defense float64`. Константа `maxDefense = 0.75` — больше 75% не поглощает никакая броня.

| Метод | Что делает | Зависит от |
|---|---|---|
| `NewArmor(code, name, basePrice, defense)` | Конструктор. | — |
| `Clone()` | Копия брони. | → `cloneBase` |
| `Defense() float64` | Защита + бонусы зачарований, но не выше `maxDefense`. | → `Enchantment.BonusDefense()` |
| `PossibleEnchantments()` / `Enchant(e)` | Для брони подходят `ForArmor` (Прочное, Охранное). | → `baseItem.enchant`, `ForArmor` |
| `Description()` | «защита P%». | → `Defense` |

####  `item/enchantment.go` — тип `Enchantment`
Зачарование. Все поля неэкспортируемые, меняются только при создании, поэтому зачарование неизменяемо, и одно и то же можно безопасно вешать на разные предметы.

Готовые значения (пакетные переменные):

| Переменная | Название | Бонус урона | Бонус защиты | Бонус цены |
|---|---|---|---|---|
| `Sharp` | Острое | +3 | — | +20 |
| `Fiery` | Огненное | +6 | — | +50 |
| `Sturdy` | Прочное | — | +5% | +20 |
| `Warding` | Охранное | — | +10% | +50 |

Списки: `ForWeapons = {Sharp, Fiery}`, `ForArmor = {Sturdy, Warding}`.
Методы-геттеры: `Code()`, `Name()`, `BonusDamage()`, `BonusDefense()`, `BonusPrice()`.

####  `item/catalog.go` — тип `ItemCatalog`
Каталог шаблонов. Поле: `prototypes []Item`. Каталог **не знает конкретных типов** предметов: копию создаёт сам шаблон.

| Метод | Что делает | Зависит от |
|---|---|---|
| `Add(prototype Item)` | Добавляет шаблон. | — |
| `Create(code) (Item, error)` | Находит шаблон по коду и возвращает его копию; если не найден — `UnknownCodeError`. | → `Item.Code`, `Item.Clone`, ⇒ `core.UnknownCodeError` |
| `CreateWeapon(code) (Weapon, error)` | То же, но проверяет, что предмет — оружие. | → `Create`, приведение типа к `Weapon` |
| `CreateArmor(code) (*Armor, error)` | То же, но проверяет, что предмет — броня. | → `Create`, приведение типа к `*Armor` |
| `Codes() []string` | Коды всех шаблонов. | → `Item.Code` |

####  `item/codes.go` (типов нет)
Константы с кодами предметов и функция `NewDefaultCatalog()`, наполняющая каталог всеми предметами игры.

| Группа | Предметы (код → название: параметры) |
|---|---|
| Стартовое снаряжение | `TrainingSword` Тренировочный меч (урон 12, цена 20); `HuntingBow` Охотничий лук (10, крит 30%, 20); `ApprenticeStaff` Посох ученика (14, разброс 30%, 20); `Chainmail` Кольчуга (защита 30%, 40); `LeatherArmor` Кожаная броня (20%, 30); `LinenRobe` Льняная роба (10%, 15) |
| Товары торговца | `IronSword` Железный меч (17, 80); `ElvenBow` Эльфийский лук (14, крит 40%, 90); `EmberStaff` Тлеющий посох (20, 90); `PlateArmor` Латы (40%, 120) |
| Добыча | `GoblinCleaver` Гоблинский тесак (15, 40); `BoneBow` Костяной лук (17, 140); `NecromancerStaff` Посох некроманта (25, 150); `KnightBlade` Клинок рыцаря (23, 150); `CryptMail` Кольчуга стража склепа (45%, 160); `DawnBlade` Клинок рассвета (35, 500); `DragonscaleArmor` Драконья чешуя (55%, 500) |

---

###  `kit/` — стартовое снаряжение (паттерн **Абстрактная фабрика**)

**За что отвечает папка.** Выдаёт герою **согласованный набор** «оружие + броня» своего класса. Код, который создаёт героя, работает только с интерфейсом `StarterKitFactory` и не знает, какой именно набор получит.

####  `kit/starter_kit.go`
##### Тип `StarterKitFactory` (interface)
Абстрактная фабрика с двумя методами: `CreateWeapon() (item.Weapon, error)` и `CreateArmor() (*item.Armor, error)`. Зависит от `item.Weapon`, `item.Armor`.

#### `kit/warrior.go`, `ranger.go`, `mage.go`
Три конкретные фабрики. Каждая хранит `catalog *item.ItemCatalog` и берёт предметы копированием из каталога.

| Тип | Оружие | Броня | Конструктор |
|---|---|---|---|
| `WarriorKitFactory` | Тренировочный меч | Кольчуга | `NewWarriorKitFactory(catalog)` |
| `RangerKitFactory` | Охотничий лук | Кожаная броня | `NewRangerKitFactory(catalog)` |
| `MageKitFactory` | Посох ученика | Льняная роба | `NewMageKitFactory(catalog)` |

| Метод (у всех трёх) | Что делает | Зависит от |
|---|---|---|
| `CreateWeapon()` | Создаёт копию стартового оружия класса. | ⊨ `StarterKitFactory`; → `ItemCatalog.CreateWeapon` |
| `CreateArmor()` | Создаёт копию стартовой брони класса. | → `ItemCatalog.CreateArmor` |

Конструкторы возвращают интерфейс `StarterKitFactory`, а не конкретный тип.

---

###  `hero/` — герой (паттерн **Строитель**)

**За что отвечает папка.** Модель игрока: класс героя, характеристики, прокачка, инвентарь, экипировка, и строитель, через который героя можно создать.

####  `hero/hero_class.go`

##### Тип `HeroClass` (struct)
Описание класса героя. Поля неэкспортируемые, поэтому классы — это только готовые значения, заранее созданные в пакете. Поля: `code`, `displayName`, `baseHealth` и `starterKit` — функция, создающая фабрику снаряжения этого класса.

| Значение | Код | Название | Базовое здоровье | Фабрика |
|---|---|---|---|---|
| `Warrior` | warrior | Воин | 120 | `kit.NewWarriorKitFactory` |
| `Ranger` | ranger | Следопыт | 100 | `kit.NewRangerKitFactory` |
| `Mage` | mage | Маг | 85 | `kit.NewMageKitFactory` |

| Метод / функция | Что делает | Зависит от |
|---|---|---|
| `HeroClasses() []*HeroClass` | Все классы по порядку. | — |
| `HeroClassByCode(code)` | Находит класс по коду, иначе `UnknownCodeError`. | → `HeroClasses`, ⇒ `core.UnknownCodeError` |
| `(c) StarterKit(catalog)` | Каждый класс сам знает свою фабрику стартового снаряжения. | → поле `starterKit` → `kit.New…KitFactory` |
| `Code()`, `DisplayName()`, `BaseHealth()` | Геттеры. | — |

####  `hero/hero.go`

Константы: `maxNameLength = 20`, `startingGold = 50`, `healthPerLevel = 15`, `damagePerLevel = 3`, `experiencePerLevel = 100`.

##### Тип `Builder` (struct)
Строитель героя. Поля: `name`, `heroClass`, `weapon`, `armor`, `level`, `gold` и `err` (первая ошибка при заполнении, её вернёт `Build`). Необязательные поля сразу получают значения по умолчанию: уровень 1, золото 50.

| Метод | Что делает | Зависит от |
|---|---|---|
| `NewBuilder()` | Создаёт строителя с уровнем 1 и золотом 50. | — |
| `SetName`, `SetHeroClass`, `SetWeapon`, `SetArmor`, `SetLevel`, `SetGold` | Заполняют поля и возвращают самого строителя, поэтому вызовы можно писать цепочкой. | — |
| `SetStarterKit(kit)` | Берёт оружие и броню из стартового набора; первая ошибка запоминается. | → `StarterKitFactory.CreateWeapon/CreateArmor`, `fail` |
| `Build() (*Hero, error)` | Проверяет все обязательные поля: имя не пустое и ≤ 20 символов, класс выбран, оружие и броня заданы, уровень ≥ 1, золото ≥ 0. Только после этого создаёт героя и выставляет полное здоровье. | → `MaxHealth` |
| `fail(err)` | Запоминает только первую ошибку. | — |

##### Тип `Hero` (struct) ⊨ `combat.Combatant`
Игрок. Поля неэкспортируемые: `name`, `heroClass`, `level`, `experience`, `gold`, `health`, `weapon`, `armor`, `inventory`. Создать героя вне пакета можно только через `Builder`.

**Бой** (методы `Combatant`):

| Метод | Что делает | Зависит от |
|---|---|---|
| `RollDamage()` | Удар оружия + бонус за уровень `(level−1)×3`; крит наследуется от оружия. | → `Weapon.RollDamage` |
| `ReceiveHit(damage)` | Урон уменьшается бронёй: `damage × (1 − защита)`; броня всегда пропускает минимум 1 единицу. Возвращает реально прошедший урон. | → `Armor.Defense` |
| `IsAlive()` | `health > 0`. | — |

**Поведение между экспедициями:**

| Метод | Что делает | Зависит от |
|---|---|---|
| `RestoreHealth()` | Полное восстановление здоровья перед экспедицией. | → `MaxHealth` |
| `GainExperience(amount) int` | Начисляет опыт, повышает уровень, пока хватает опыта (порог `уровень × 100`), возвращает число полученных уровней. | → `ExperienceToNextLevel` |
| `AddGold(amount)` | Добавляет золото. | — |
| `SpendGold(amount) error` | Списывает золото; если не хватает — ошибка. | — |
| `AddItem(it)` | Кладёт предмет в инвентарь. | — |
| `RemoveItem(it) error` | Убирает предмет из инвентаря; если его там нет — ошибка. | → `slices.Index`, `slices.Delete` |
| `Equip(it) error` | Экипирует предмет из инвентаря; снятый предмет возвращается в инвентарь. Через `switch` по типу различает оружие и броню; всё остальное экипировать нельзя. | → `RemoveItem`, типы `item.Weapon`, `*item.Armor` |

**Геттеры:** `Code()` (код класса), `Name()`, `Health()`, `MaxHealth()` (базовое здоровье класса + `(level−1)×15`), `HeroClass()`, `Level()`, `Experience()`, `ExperienceToNextLevel()` (`level × 100`), `Gold()`, `Weapon()`, `Armor()`, `Inventory()` (возвращает **копию** списка, чтобы инвентарь нельзя было изменить снаружи).

---

### `enemy/` — враги

**За что отвечает папка.** Противники героя. Все они — бойцы (`Combatant`) с наградой за победу.

####  `enemy/enemy.go`
##### Тип `Enemy` (interface)
Враг: `combat.Combatant` плюс `ExperienceReward() int` и `RollGoldReward() int`. Этот интерфейс возвращает фабричный метод локаций.

####  `enemy/base_enemy.go`
##### Тип `baseEnemy` (struct) ⊨ `Enemy`
Общая логика врагов. Конкретные враги встраивают её в себя. Поля: `code`, `name`, `maxHealth`, `minDamage`, `maxDamage`, `resistance` (доля урона, которую враг игнорирует), `experienceReward`, `goldReward`, `health`.

| Метод | Что делает | Зависит от |
|---|---|---|
| `newBaseEnemy(...)` | Создаёт заготовку; здоровье = максимальному. | — |
| `RollDamage()` | Случайный урон между `minDamage` и `maxDamage`, без крита. | → `core.GetInstance().Roll` |
| `ReceiveHit(damage)` | Урон уменьшается сопротивлением; проходит минимум 1. | — |
| `IsAlive()` | `health > 0`. | — |
| `RollGoldReward()` | Случайное золото от половины до полной награды. | → `core.GetInstance().Roll` |
| `ExperienceReward()` | Фиксированный опыт. | — |
| `Code()`, `Name()`, `Health()`, `MaxHealth()` | Геттеры. | — |

####  `enemy/goblin.go`, `skeleton.go`, `dragon.go`
Каждый встраивает `baseEnemy` (`⊂ baseEnemy`) и отличается только параметрами.

| Враг | Здоровье | Урон | Сопротивление | Опыт | Золото (макс.) | Особенность |
|---|---|---|---|---|---|---|
| `Goblin` | 50 | 6–14 | 0% | 40 | 30 | Простой враг без особенностей |
| `Skeleton` | 110 | 10–24 | 30% | 90 | 70 | Кости плохо поддаются ударам |
| `Dragon` | 260 | 16–22 | 25% | 500 | 400 | С шансом 25% дышит огнём: урон ×2, удар считается критическим |

| Метод | Что делает | Зависит от |
|---|---|---|
| `NewGoblin()`, `NewSkeleton()`, `NewDragon()` | Конструкторы. | → `newBaseEnemy` |
| `Dragon.RollDamage()` | **Переопределяет** метод встроенной структуры: бой вызывает его через интерфейс, поэтому срабатывает именно он. Сначала бросает обычный укус, затем с шансом 25% удваивает урон. | → `baseEnemy.RollDamage`, `core.GetInstance().Chance` |

---

###  `world/` — локации и экспедиции (паттерн **Фабричный метод**)

**За что отвечает папка.** Места, куда герой ходит в экспедиции, и общая логика самой экспедиции. Каждая локация сама решает, какого врага создать, а вся остальная логика (бой, награда, штраф) общая для всех и написана один раз.

####  `world/location.go`

##### Тип `Location` (interface)
Локация. Методы: `Code`, `DisplayName`, `RecommendedLevel`, `Catalog`, **`CreateEnemy() enemy.Enemy`** (фабричный метод — создание врага делегировано конкретным локациям) и `LootTable() []string` (коды предметов, которые можно найти).

##### Тип `baseLocation` (struct)
Общие данные всех локаций. Поля: `code`, `displayName`, `recommendedLevel`, `catalog`. Методы: `Code()`, `DisplayName()`, `RecommendedLevel()`, `Catalog()`. Встраивается в `Forest`, `Crypt`, `DragonLair`.

##### Функции пакета (общая бизнес-логика экспедиции)
Константы: `lootChance = 0.6` (шанс найти добычу), `enchantChance = 0.3` (шанс, что добыча зачарована).

| Функция | Что делает | Зависит от |
|---|---|---|
| `Explore(location, hero)` | Экспедиция со случайным сидом. | → `core.NextSeed`, `ExploreWithSeed` |
| `ExploreWithSeed(location, hero, seed)` | Главная функция: задаёт сид, восстанавливает здоровье, создаёт врага **через `location.CreateEnemy()`**, проводит бой, при победе начисляет награду, при поражении штрафует. Возвращает `ExpeditionReport`. | → `core.Reseed`, `Hero.RestoreHealth`, `Location.CreateEnemy`, `combat.NewBattle().Run`, `reward`, `punish` |
| `reward(location, hero, foe)` | Начисляет опыт, золото и добычу. | → `Enemy.ExperienceReward`, `Enemy.RollGoldReward`, `rollLoot`, `Hero.GainExperience/AddGold/AddItem` |
| `rollLoot(location)` | С шансом 60% берёт случайный код из `LootTable`, создаёт **копию** предмета из каталога и с шансом 30% накладывает на неё случайное подходящее зачарование. | → `core.Chance/Index`, `Catalog().Create`, `Item.PossibleEnchantments`, `Item.Enchant` |
| `punish(hero)` | Проигравший теряет половину золота. | → `Hero.Gold`, `Hero.SpendGold` |

####  `world/forest.go`, `crypt.go`, `dragon_lair.go`
Конкретные локации, встраивают `baseLocation`.

| Тип | Код | Название | Рек. уровень | Враг (`CreateEnemy`) | Добыча (`LootTable`) |
|---|---|---|---|---|---|
| `Forest` | forest | Мистический лес | 1 | `enemy.NewGoblin()` | Гоблинский тесак, Охотничий лук, Кожаная броня |
| `Crypt` | crypt | Заброшенный склеп | 3 | `enemy.NewSkeleton()` | Клинок рыцаря, Костяной лук, Посох некроманта, Кольчуга стража склепа |
| `DragonLair` | dragon-lair | Логово дракона | 5 | `enemy.NewDragon()` | Клинок рассвета, Драконья чешуя |

Методы у всех: конструктор (`NewForest`, `NewCrypt`, `NewDragonLair`), `CreateEnemy()` и `LootTable()`. Константы кодов: `ForestCode`, `CryptCode`, `DragonLairCode`.

#### `world/registry.go`
##### Тип `LocationRegistry` (struct)
Реестр всех локаций игры. Поле: `locations []Location`.

| Метод | Что делает | Зависит от |
|---|---|---|
| `NewLocationRegistry(locations...)` | Создаёт реестр. | — |
| `Get(code)` | Находит локацию по коду, иначе `UnknownCodeError`. | → `Location.Code`, ⇒ `core.UnknownCodeError` |
| `All()` | Копия списка всех локаций. | → `slices.Clone` |

####  `world/rewards.go`
##### Тип `Rewards` (struct)
Итоги экспедиции для героя: `Experience`, `Gold` (**отрицательное**, если герой проиграл и потерял золото), `Loot item.Item` (`nil`, если добычи нет), `LevelsGained`.

#### 📄 `world/expedition_report.go`
##### Тип `ExpeditionReport` (struct)
Полный отчёт: `Seed`, `LocationCode`, `Battle combat.BattleReport`, `Rewards`.

---

### `camp/` — лагерь

**За что отвечает папка.** Торговля между экспедициями.

#### 📄 `camp/merchant.go`
##### Тип `Merchant` (struct)
Торговец. Поля: `catalog` и `assortment []string` (коды товаров). Покупатель всегда получает **свою копию** предмета из каталога. Константа `sellPricePercent = 50` — торговец выкупает предметы за половину цены.

| Метод | Что делает | Зависит от |
|---|---|---|
| `NewMerchant(catalog, assortment...)` | Создаёт торговца с ассортиментом. | — |
| `Showcase()` | Образцы товаров для показа. | → `ItemCatalog.Create` |
| `Buy(hero, code)` | Проверяет, что товар есть в ассортименте, создаёт копию, списывает золото, кладёт предмет в инвентарь. | → `ItemCatalog.Create`, `Hero.SpendGold`, `Hero.AddItem`, ⇒ `core.UnknownCodeError` |
| `Sell(hero, item)` | Убирает предмет из инвентаря и платит герою. | → `SellPrice`, `Hero.RemoveItem`, `Hero.AddGold` |
| `SellPrice(item)` | Цена выкупа: 50% от стоимости (с учётом зачарований). | → `Item.Price` |

Ассортимент в игре (задаётся в `main.go`): Железный меч, Эльфийский лук, Тлеющий посох, Кольчуга, Латы.

---

###  `console/` — вывод в консоль

**За что отвечает папка.** Весь вывод собран в одном месте и не содержит игровых правил: он только красиво печатает то, что ему передали.

####  `console/view.go`
##### Тип `ConsoleView` (struct{})
Пустая структура, только методы.

| Метод | Что делает | Зависит от |
|---|---|---|
| `Describe(item) string` | Строка «Название [зачарования] (характеристики)». | → `Item.Name/Enchantments/Description` |
| `PrintHero(hero)` | Карточка героя: имя, класс, уровень, здоровье, опыт, золото, оружие, броня. | → геттеры `Hero`, `Describe` |
| `PrintReport(report, locationName)` | Печатает всю экспедицию: заголовок с сидом, встречу с врагом, события боя, награды. | → `printEvent`, `printRewards` |
| `printEvent(event, battle)` | Строка одного события боя (удар, крит или смерть). | → `participant` |
| `printRewards(battle, rewards)` | Итог: победа (опыт, золото, добыча, новый уровень) или поражение (потеря золота). | → `Describe` |
| `participant(side, battle)` | Функция пакета: выбирает героя или врага по стороне. | — |

---

###  `main.go` — игра

Точка входа. Типов нет. Глобальные переменные: `input` (сканер ввода) и `view` (`console.ConsoleView`).

| Функция | Что делает | Зависит от |
|---|---|---|
| `main()` | Собирает каталог, реестр локаций (лес, склеп, логово) и торговца; создаёт героя; запускает цикл меню: 1 — экспедиция, 2 — снаряжение, 3 — торговец, 0 — выход. | → `item.NewDefaultCatalog`, `world.NewLocationRegistry`, `camp.NewMerchant`, `createHero` |
| `createHero(catalog)` | Меню выбора класса и ввода имени; героя собирает строителем, при ошибке спрашивает заново. | → `HeroClasses`, `StarterKit`, `hero.NewBuilder()…Build()` |
| `goOnExpedition(hero, locations)` | Выбор локации и запуск экспедиции; печатает отчёт; при победе над драконом — финальная фраза. | → `world.Explore`, `view.PrintReport` |
| `manageEquipment(hero)` | Выбор предмета из инвентаря и экипировка. | → `Hero.Inventory`, `Hero.Equip` |
| `visitMerchant(hero, merchant)` | Меню торговца: купить / продать. | → `buy`, `sell` |
| `buy(hero, merchant)` | Показывает витрину и покупает выбранное. | → `Merchant.Showcase`, `Merchant.Buy` |
| `sell(hero, merchant)` | Показывает инвентарь с ценами выкупа и продаёт выбранное. | → `Merchant.SellPrice`, `Merchant.Sell` |
| `printItems(items)` | Печатает пронумерованный список предметов. | → `view.Describe` |
| `readLine()` | Читает строку без пробелов по краям; при конце ввода завершает программу. | — |
| `readChoice(count)` | Читает номер от 1 до `count`, возвращает индекс или −1, если ввод некорректный. | → `readLine` |

###  `cmd/demo/main.go` — демонстрация паттернов

Отдельная программа (`go run ./cmd/demo`), по очереди показывает работу каждого паттерна: Одиночка (два вызова `GetInstance` дают один объект), Прототип (две копии меча независимы, зачарование одной не влияет на другую), Абстрактная фабрика (наборы всех трёх классов), Строитель (успешная сборка героя и две ошибочные), Фабричный метод (экспедиции в лес и склеп с одинаковым сидом 7). Вспомогательная функция `must[T any](value T, err error) T` останавливает демо паникой, если случилась ошибка (данные заданы заранее, поэтому ошибка означала бы баг в коде).

---

## 4. Паттерны: где что находится

| Паттерн | Где | Как работает |
|---|---|---|
| **Одиночка** | `core/dice.go` | `GetInstance()` + `sync.Once` + неэкспортируемый тип `dice`: один генератор случайности на всю игру |
| **Прототип** | `item/` (`Item.Clone`, `ItemCatalog`) | Каталог хранит шаблоны и отдаёт только копии; копию делает сам предмет, поэтому каталог не знает конкретных типов |
| **Абстрактная фабрика** | `kit/` (`StarterKitFactory` + три фабрики) | Согласованная пара «оружие + броня» для каждого класса |
| **Строитель** | `hero/hero.go` (`Builder`) | Пошаговая сборка героя с единой проверкой в `Build()` |
| **Фабричный метод** | `world/` (`Location.CreateEnemy`) | Локация сама создаёт врага, общая логика экспедиции написана один раз в `ExploreWithSeed` |

Дополнительные приёмы Go, которые используются в проекте: встраивание структур вместо наследования (`baseItem`, `baseWeapon`, `baseEnemy`, `baseLocation`), переопределение метода встроенной структуры (`Dragon.RollDamage`), `switch` по типу (`Hero.Equip`), приведение типа (`ItemCatalog.CreateWeapon/CreateArmor`), неэкспортируемые поля для защиты инвариантов, дженерик-функция (`must[T]`).

---

## 5. Сквозной сценарий: что происходит при экспедиции

1. `main` вызывает `world.Explore(location, hero)`.
2. `Explore` берёт новый сид у `core.GetInstance().NextSeed()` и передаёт его в `ExploreWithSeed`.
3. `ExploreWithSeed` вызывает `Reseed(seed)`: дальше все броски определены сидом и воспроизводимы.
4. Герой восстанавливает здоровье (`Hero.RestoreHealth`).
5. **Фабричный метод:** `location.CreateEnemy()` создаёт врага конкретной локации (гоблин, скелет или дракон).
6. `combat.NewBattle(hero, enemy).Run()` проводит бой. Для `Battle` и герой, и враг — просто `Combatant`.
   - Удар героя: `Hero.RollDamage` → `Weapon.RollDamage` (меч — стабильно, лук — шанс крита, посох — разброс) + бонус уровня; цель получает его через `ReceiveHit` (у врага минус сопротивление).
   - Удар врага: `baseEnemy.RollDamage` (у дракона — `Dragon.RollDamage` с огненным дыханием); герой получает его через `Hero.ReceiveHit` (минус защита брони).
   - Каждый удар и смерть пишутся как `BattleEvent` в журнал.
7. По итогу боя создаётся `BattleReport` с исходом `Victory` или `Defeat`.
8. **Победа:** `reward` начисляет опыт (`GainExperience` может повысить уровень), золото и с шансом 60% добычу (`rollLoot`: копия шаблона из каталога, с шансом 30% зачарованная).
   **Поражение:** `punish` отбирает половину золота.
9. Всё собирается в `ExpeditionReport` (сид, локация, бой, награды), и `console.ConsoleView.PrintReport` печатает его.
