package item

// коды шаблонов: по ним предметы ищутся в каталоге
// стартовое снаряжение
const TrainingSword = "training-sword"
const HuntingBow = "hunting-bow"
const ApprenticeStaff = "apprentice-staff"
const Chainmail = "chainmail"
const LeatherArmor = "leather-armor"
const LinenRobe = "linen-robe"

// товары торговца
const IronSword = "iron-sword"
const ElvenBow = "elven-bow"
const EmberStaff = "ember-staff"
const PlateArmor = "plate-armor"

// добыча
const GoblinCleaver = "goblin-cleaver"
const BoneBow = "bone-bow"
const NecromancerStaff = "necromancer-staff"
const KnightBlade = "knight-blade"
const CryptMail = "crypt-mail"
const DawnBlade = "dawn-blade"
const DragonscaleArmor = "dragonscale-armor"

// NewDefaultCatalog возвращает каталог со всеми предметами, которые есть в игре
func NewDefaultCatalog() *ItemCatalog {
	catalog := &ItemCatalog{}

	// стартовое снаряжение
	catalog.Add(NewSword(TrainingSword, "Тренировочный меч", 20, 12))
	catalog.Add(NewBow(HuntingBow, "Охотничий лук", 20, 10, 0.3))
	catalog.Add(NewStaff(ApprenticeStaff, "Посох ученика", 20, 14, 0.3))
	catalog.Add(NewArmor(Chainmail, "Кольчуга", 40, 0.3))
	catalog.Add(NewArmor(LeatherArmor, "Кожаная броня", 30, 0.2))
	catalog.Add(NewArmor(LinenRobe, "Льняная роба", 15, 0.1))

	// товары торговца
	catalog.Add(NewSword(IronSword, "Железный меч", 80, 17))
	catalog.Add(NewBow(ElvenBow, "Эльфийский лук", 90, 14, 0.4))
	catalog.Add(NewStaff(EmberStaff, "Тлеющий посох", 90, 20, 0.3))
	catalog.Add(NewArmor(PlateArmor, "Латы", 120, 0.4))

	// добыча
	catalog.Add(NewSword(GoblinCleaver, "Гоблинский тесак", 40, 15))
	catalog.Add(NewBow(BoneBow, "Костяной лук", 140, 17, 0.4))
	catalog.Add(NewStaff(NecromancerStaff, "Посох некроманта", 150, 25, 0.35))
	catalog.Add(NewSword(KnightBlade, "Клинок рыцаря", 150, 23))
	catalog.Add(NewArmor(CryptMail, "Кольчуга стража склепа", 160, 0.45))
	catalog.Add(NewSword(DawnBlade, "Клинок рассвета", 500, 35))
	catalog.Add(NewArmor(DragonscaleArmor, "Драконья чешуя", 500, 0.55))

	return catalog
}
