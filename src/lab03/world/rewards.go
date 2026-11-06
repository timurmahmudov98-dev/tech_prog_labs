package world

import "lab03/item"

// Rewards - итоги экспедиции для героя;
// Gold отрицательный, если герой проиграл и потерял золото, Loot равен nil, если добычи нет
type Rewards struct {
	Experience   int
	Gold         int
	Loot         item.Item
	LevelsGained int
}
