package core

import (
	"math/rand/v2"
	"sync"
)

// dice - одиночка: общий источник случайности на всю игру.
// Тип неэкспортируемый (с маленькой буквы), поэтому вне пакета core создать кость нельзя.
type dice struct {
	random *rand.Rand
}

var (
	// единственный экземпляр
	instance *dice
	once     sync.Once
)

// GetInstance возвращает единственный экземпляр кости
func GetInstance() *dice {
	// sync.Once выполнит функцию ровно один раз, даже если GetInstance одновременно вызовут несколько горутин
	once.Do(func() {
		instance = &dice{random: rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))}
	})

	return instance
}

// Reseed задает сид: после этого все последующие броски повторяются в точности
func (d *dice) Reseed(seed uint64) {
	d.random = rand.New(rand.NewPCG(seed, seed))
}

// NextSeed возвращает случайный сид
func (d *dice) NextSeed() uint64 {
	return d.random.Uint64N(1_000_000)
}

// Roll возвращает случайное целое число от low до high включительно
func (d *dice) Roll(low, high int) int {
	return low + d.random.IntN(high-low+1)
}

// Chance возвращает true с вероятностью probability (от 0 до 1)
func (d *dice) Chance(probability float64) bool {
	return d.random.Float64() < probability
}

// Index возвращает случайный индекс от 0 до size-1, например чтобы выбрать элемент слайса
func (d *dice) Index(size int) int {
	return d.random.IntN(size)
}
